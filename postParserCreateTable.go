package main

import (
	"fmt"
	"strconv"
	"strings"
)

func createTableFromAST(createNode *ASTNode) {
	tableName := createNode.tableName
	columns := createNode.columns
	newColumns := make([]Column, len(columns))
	for i, column := range columns {
		newColumns[i] = Column{
			Name:         column.name,
			Type:         column._type,
			Conditions:   []string{},
			VarCharLimit: 0,
		}
		if column._type == "VARCHAR" {
			newColumns[i].VarCharLimit = column.varCharLimit
		}
	}
	createTable(tableName, newColumns)
}

func insertIntoFromAST(insertNode *ASTNode) {
	tableName := insertNode.tableName
	if !tableExists(tableName) {
		fmt.Printf("Table %s does not exist\n", tableName)
		return
	}

	table := database[tableName]
	columnNames := insertNode.columnNames
	columnValues := insertNode.columnValues

	// Support implicit column list
	if len(columnNames) == 0 {
		if len(columnValues) != len(table.Columns) {
			fmt.Println("Value count doesn't match number of table columns")
			return
		}
		for _, col := range table.Columns {
			columnNames = append(columnNames, col.Name)
		}
	}

	if len(columnNames) != len(columnValues) {
		fmt.Println("Mismatch between number of column names and values")
		return
	}

	newRow := make(map[string]interface{})

	for _, col := range table.Columns {
		for j, colName := range columnNames {
			if col.Name != colName {
				continue
			}

			val := strings.TrimSpace(columnValues[j])

			// Constraint checks
			for _, constraint := range col.Conditions {
				if constraint == "PRIMARY KEY" || constraint == "UNIQUE" {
					if !checkUnique(tableName, col.Name, val) {
						fmt.Printf("Column %s must be unique\n", col.Name)
						return
					}
				}
				if strings.HasPrefix(constraint, "DEFAULT") && val == "" {
					val = strings.TrimPrefix(constraint, "DEFAULT ")
				}
				if constraint == "NOT NULL" && val == "" {
					fmt.Printf("Column %s cannot be NULL\n", col.Name)
					return
				}
			}

			// Type casting
			switch strings.ToUpper(col.Type) {
			case "INT":
				if val != "" {
					intVal, err := strconv.Atoi(val)
					if err != nil {
						fmt.Printf("Column %s must be an integer\n", col.Name)
						return
					}
					newRow[col.Name] = intVal
				} else {
					newRow[col.Name] = nil
				}
			case "FLOAT":
				if val != "" {
					floatVal, err := strconv.ParseFloat(val, 64)
					if err != nil {
						fmt.Printf("Column %s must be a float\n", col.Name)
						return
					}
					newRow[col.Name] = floatVal
				} else {
					newRow[col.Name] = nil
				}
			default:
				if strings.HasPrefix(strings.ToUpper(col.Type), "VARCHAR") {
					maxLen := col.VarCharLimit
					if len(val) > maxLen {
						fmt.Printf("Column %s exceeds VARCHAR(%d)\n", col.Name, maxLen)
						return
					}
					newRow[col.Name] = val
				} else {
					newRow[col.Name] = val
				}
			}
		}
	}

	table.Rows = append(table.Rows, newRow)
	database[tableName] = table
	fmt.Printf("Inserted row into %s\n", tableName)
}

func isGroupByColumn(columnTypes []columnType, columnNames []string, columnName string) bool {
	for i, name := range columnNames {
		if name == columnName {
			return columnTypes[i] == COLUMN_TYPE_GROUP_BY
		}
	}
	return false
}

// ------------------- SELECT with AVG support -------------------
func selectFromAST(selectNode *ASTNode) Table {
	tableName := selectNode.tableName
	srcTable := database[tableName]

	// Build schema: one Column per selectNode.column + a hidden "count" for AVG
	newCols := make([]Column, len(selectNode.columnNames)+1)
	for i, origName := range selectNode.columnNames {
		ct := selectNode.columnTypes[i]

		// decide visibility
		vis := (ct == COLUMN_TYPE_GROUP_BY ||
			ct == COLUMN_TYPE_NORMAL ||
			ct == COLUMN_TYPE_COUNT ||
			ct == COLUMN_TYPE_MAX ||
			ct == COLUMN_TYPE_MIN ||
			ct == COLUMN_TYPE_SUM ||
			ct == COLUMN_TYPE_AVG)

		// decide type
		typ := ""
		if ct == COLUMN_TYPE_GROUP_BY {
			for _, c := range srcTable.Columns {
				if c.Name == origName {
					typ = c.Type
					break
				}
			}
		} else {
			typ = "float64"
		}

		// compute alias: user‐supplied if given, otherwise for aggregates use func_origName
		alias := selectNode.columnAliases[i]
		// fmt.Println("selectFromAST: origName=%s, typ=%s, vis=%t, alias=%s", origName, typ, vis, alias)
		if alias == "" && ct != COLUMN_TYPE_NORMAL && ct != COLUMN_TYPE_GROUP_BY {
			switch ct {
			case COLUMN_TYPE_SUM:
				alias = "sum_" + origName
			case COLUMN_TYPE_AVG:
				alias = "avg_" + origName
			case COLUMN_TYPE_MIN:
				alias = "min_" + origName
			case COLUMN_TYPE_MAX:
				alias = "max_" + origName
			case COLUMN_TYPE_COUNT:
				alias = "count_" + origName
			}
		}

		newCols[i] = Column{
			Name:           origName,
			Type:           typ,
			Conditions:     nil,
			VarCharLimit:   0,
			FunctionResult: ct != COLUMN_TYPE_NORMAL && ct != COLUMN_TYPE_GROUP_BY,
			Visible:        vis,
			Alias:          alias,
		}
	}
	// Hidden global count for AVG denominator
	countIdx := len(selectNode.columnNames)
	newCols[countIdx] = Column{
		Name:           "count",
		Type:           "float64",
		Conditions:     nil,
		VarCharLimit:   0,
		FunctionResult: true,
		Visible:        false,
		Alias:          "count",
	}

	result := Table{
		Name:    "result",
		Columns: newCols,
		Rows:    []map[string]interface{}{},
	}

	// figure out which cols are GROUP_BY
	groupByIdx := []int{}
	for i, ct := range selectNode.columnTypes {
		if ct == COLUMN_TYPE_GROUP_BY {
			groupByIdx = append(groupByIdx, i)
		}
	}

	// aggregate rows
	for _, srcRow := range srcTable.Rows {
		// find matching bucket
		bucket := -1
		for ri, outRow := range result.Rows {
			match := true
			for _, gi := range groupByIdx {
				orig := selectNode.columnNames[gi]
				if outRow[orig] != srcRow[orig] {
					match = false
					break
				}
			}
			if match {
				bucket = ri
				break
			}
		}

		if bucket >= 0 {
			// update aggregates
			outRow := result.Rows[bucket]
			for i, ct := range selectNode.columnTypes {
				alias := selectNode.columnNames[i] // columnAliases[i]
				switch ct {
				case COLUMN_TYPE_SUM, COLUMN_TYPE_AVG:
					outRow[alias] = toFloat64(outRow[alias]) + toFloat64(srcRow[selectNode.columnNames[i]])
				case COLUMN_TYPE_COUNT:
					outRow[alias] = toFloat64(outRow[alias]) + 1
				case COLUMN_TYPE_MIN:
					outRow[alias] = min(outRow[alias], srcRow[selectNode.columnNames[i]])
				case COLUMN_TYPE_MAX:
					outRow[alias] = max(outRow[alias], srcRow[selectNode.columnNames[i]])
				}
			}
			// always bump count
			outRow["count"] = toFloat64(outRow["count"]) + 1
		} else {
			// first time: initialize a new bucket
			newRow := make(map[string]interface{}, len(newCols))
			for i, ct := range selectNode.columnTypes {
				alias := selectNode.columnNames[i]
				val := srcRow[selectNode.columnNames[i]]
				switch ct {
				case COLUMN_TYPE_SUM, COLUMN_TYPE_AVG:
					newRow[alias] = toFloat64(val)
				case COLUMN_TYPE_COUNT:
					newRow[alias] = float64(1)
				case COLUMN_TYPE_MIN, COLUMN_TYPE_MAX:
					newRow[alias] = val
				default: // GROUP_BY or NORMAL
					newRow[alias] = val
				}
			}
			newRow["count"] = float64(1)
			result.Rows = append(result.Rows, newRow)
		}
	}

	// now, post‑process AVG columns: sum/count → avg
	for _, row := range result.Rows {
		for i, ct := range selectNode.columnTypes {
			if ct == COLUMN_TYPE_AVG {
				alias := selectNode.columnNames[i]
				sum := toFloat64(row[alias])
				cnt := toFloat64(row["count"])
				if cnt != 0 {
					row[alias] = sum / cnt
				} else {
					row[alias] = float64(0)
				}
			}
		}
	}

	return result
}

func evalExpression(expr *ASTNode, row map[string]interface{}) interface{} {
	switch expr.Type {
	case AST_INT_LITERAL:
		return expr.intVal
	case AST_FLOAT_LITERAL:
		return expr.floatVal
	case AST_VARCHAR_LITERAL:
		return expr.strVal
	case AST_BOOLEAN_LITERAL:
		return expr.boolValue
	case AST_COLUMN_NAME:
		return row[expr.columnName]
	case AST_BINARY:
		left := evalExpression(expr.left, row)
		right := evalExpression(expr.right, row)

		if expr.operator == "+" || expr.operator == "-" || expr.operator == "*" || expr.operator == "/" { // Assume numeric for now (extend as needed)
			leftVal := toFloat64(left)
			rightVal := toFloat64(right)

			switch expr.operator {
			case "+":
				return leftVal + rightVal
			case "-":
				return leftVal - rightVal
			case "*":
				return leftVal * rightVal
			case "/":
				if rightVal == 0 {
					return float64(0) // Avoid divide-by-zero panic
				}
				return leftVal / rightVal
			default:
				panic("Unsupported operator: " + expr.operator)
			}
		} else {
			panic("Not implemented")
		}
	default:
		panic(fmt.Sprintf("Unsupported AST node type in evalExpression: %d", expr.Type))
	}
}
