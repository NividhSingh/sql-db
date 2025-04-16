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
			varCharLimit: 0,
		}
		if column._type == "VARCHAR" {
			newColumns[i].varCharLimit = column.varCharLimit
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

	columnNames := insertNode.columnNames
	columnValues := insertNode.columnValues

	if len(columnNames) != len(columnValues) {
		fmt.Println("Mismatch between number of column names and values")
		return
	}

	table := database[tableName]
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
					maxLen := col.varCharLimit
					if len(val) > maxLen {
						fmt.Printf("Column %s exceeds maximum VARCHAR(%d)\n", col.Name, maxLen)
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

func selectFromAST(selectNode *ASTNode) {
	tableName := selectNode.tableName
	selectColumns := selectNode.columns
	selectColumnNames := selectNode.columnName
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

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

func groupByAndFunctions(groupBy []string, functionColumns map[string][]string, table Table) Table {
	newTableName := "table_after_group_by_and_functions"
	newColumns := []Column{}

	for _, groupCol := range groupBy {
		found := false
		for _, column := range table.Columns {
			if column.Name == groupCol {
				newColumns = append(newColumns, Column{
					Name:           column.Name,
					Type:           column.Type,
					Conditions:     column.Conditions,
					functionResult: false,
				})
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("Column %s not found in table", groupCol))
		}
	}

	for functionCol, functions := range functionColumns {
		found := false
		for _, column := range table.Columns {
			if column.Name == functionCol {
				countColumn := false
				averageColumn := false
				sumColumn := false
				for _, function := range functions {
					if function == "COUNT" {
						countColumn = true
					} else if function == "AVG" {
						averageColumn = true
					}else if function == "SUM" {
						averageColumn = true
					}
					newColumns = append(newColumns, Column{
						Name:           column.Name + strings.ToUpper(function),
						Type:           "float64",
						Conditions:     []string{}, //{"CAST(" + functionCol + " AS FLOAT)"},
						functionResult: true,
					})
				}
				if !countColumn {
					newColumns = append(newColumns, Column{
						Name:           column.Name + "COUNT",
						Type:           "float64",
						Conditions:     []string{}, //{"CAST(" + functionCol + " AS FLOAT)"},
						functionResult: true,
					})
					functionColumns[functionCol] = append(functionColumns[functionCol], "COUNT")
				}
				if averageColumn && !sumColumn {
					newColumns = append(newColumns, Column{
						Name:           column.Name + "SUM",
						Type:           "float64",
						Conditions:     []string{}, //{"CAST(" + functionCol + " AS FLOAT)"},
						functionResult: true,
					})
					functionColumns[functionCol] = append(functionColumns[functionCol], "SUM")
				}
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("Column %s not found in table", functionCol))
		}
	}

	newTable := Table{
		Name:    newTableName,
		Columns: newColumns,
		Rows:    []map[string]interface{}{}, // Initialize with an empty slice of maps
	}
	for _, row := range table.Rows {
		newRow := map[string]interface{}{}
		foundMatch := false
		for i, newTableRow := range newTable.Rows {
			matches := true
			for _, column := range groupBy {
				if row[column] != newTableRow[column] {
					matches = false
					break
				}
			}
			if matches {
				for functionCol, functions := range functionColumns {
					for _, function := functions {
						if strings.ToUpper(function) == "COUNT" {
							newRow[functionCol + strings.ToUpper("COUNT")]++
						} else if strings.ToUpper(function) == "SUM" {
							newRow[functionCol + strings.ToUpper("SUM")] += row[functionCol]
						} else if strings.ToUpper(function) == "MAX" {
							newRow[functionCol + strings.ToUpper("MAX")] = max(newRow[functionCol + strings.ToUpper("MIN")], row[functionCol])
						} else if strings.ToUpper(function) == "MIN" {
							newRow[functionCol + strings.ToUpper("MIN")] = min(newRow[functionCol + strings.ToUpper("MIN")], row[functionCol])
						}
					}
				}
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			for functionCol, functions := range functionColumns {
				for _, function := functions {
					if strings.ToUpper(function) == "COUNT" {
						newRow[functionCol + strings.ToUpper("COUNT")] = 1
					} else if strings.ToUpper(function) == "SUM" {
						newRow[functionCol + strings.ToUpper("SUM")] = row[functionCol]
					} else if strings.ToUpper(function) == "MAX" {
						newRow[functionCol + strings.ToUpper("MAX")] = row[functionCol]
					} else if strings.ToUpper(function) == "MIN" {
						newRow[functionCol + strings.ToUpper("MIN")] =  row[functionCol]
					}
				}
			}
		}
	}

	return newTable
}
