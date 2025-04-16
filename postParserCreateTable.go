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

func groupByAndFunctions(groupBy []string, functionColumns map[string][]string, table Table) Table {
	newTableName := "table_after_group_by_and_functions"
	newColumns := []Column{}

	// Add group-by columns to new table schema.
	for _, groupCol := range groupBy {
		found := false
		for _, column := range table.Columns {
			if column.Name == groupCol {
				newColumns = append(newColumns, Column{
					Name:           column.Name,
					Type:           column.Type,
					Conditions:     column.Conditions,
					varCharLimit:   column.varCharLimit,
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

	// Add function result columns.
	for functionCol, functions := range functionColumns {
		found := false
		for _, column := range table.Columns {
			if column.Name == functionCol {
				var countColumn, sumColumn, avgRequested bool
				for _, function := range functions {
					funcName := strings.ToUpper(function)
					switch funcName {
					case "COUNT":
						countColumn = true
					case "SUM":
						sumColumn = true
					case "AVG":
						avgRequested = true
					}
					newColumns = append(newColumns, Column{
						Name:           column.Name + funcName,
						Type:           "float64",
						Conditions:     []string{},
						functionResult: true,
					})
				}
				// Always add COUNT if it wasn't specified.
				if !countColumn {
					newColumns = append(newColumns, Column{
						Name:           column.Name + "COUNT",
						Type:           "float64",
						Conditions:     []string{},
						functionResult: true,
					})
					functionColumns[functionCol] = append(functionColumns[functionCol], "COUNT")
				}
				// If AVG was requested and no SUM was provided, add a SUM column.
				if avgRequested && !sumColumn {
					newColumns = append(newColumns, Column{
						Name:           column.Name + "SUM",
						Type:           "float64",
						Conditions:     []string{},
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
		Rows:    []map[string]interface{}{},
	}

	// Process each row of the input table for grouping and aggregation.
	for _, row := range table.Rows {
		groupFound := false
		// Look for an existing aggregate row in newTable that matches all group-by columns.
		for i := range newTable.Rows {
			matches := true
			for _, col := range groupBy {
				if newTable.Rows[i][col] != row[col] {
					matches = false
					break
				}
			}
			if matches {
				// Update aggregated values directly in newTable.Rows[i].
				for functionCol, functions := range functionColumns {
					for _, function := range functions {
						key := functionCol + strings.ToUpper(function)
						switch strings.ToUpper(function) {
						case "COUNT":
							newTable.Rows[i][key] = newTable.Rows[i][key].(float64) + 1.0
						case "SUM":
							newTable.Rows[i][key] = newTable.Rows[i][key].(float64) + toFloat64(row[functionCol])
						case "MAX":
							newTable.Rows[i][key] = max(newTable.Rows[i][key].(float64), toFloat64(row[functionCol]))
						case "MIN":
							newTable.Rows[i][key] = min(newTable.Rows[i][key].(float64), toFloat64(row[functionCol]))
						}
					}
				}
				groupFound = true
				break
			}
		}
		// Create a new aggregate row if no matching group was found.
		if !groupFound {
			newRow := make(map[string]interface{})
			// Copy over the group-by column values.
			for _, col := range groupBy {
				newRow[col] = row[col]
			}
			// Initialize aggregated values.
			for functionCol, functions := range functionColumns {
				for _, function := range functions {
					key := functionCol + strings.ToUpper(function)
					switch strings.ToUpper(function) {
					case "COUNT":
						newRow[key] = 1.0
					case "SUM":
						newRow[key] = toFloat64(row[functionCol])
					case "MAX":
						newRow[key] = toFloat64(row[functionCol])
					case "MIN":
						newRow[key] = toFloat64(row[functionCol])
					}
				}
			}
			newTable.Rows = append(newTable.Rows, newRow)
		}
	}

	// Compute AVG values for groups.
	for i := range newTable.Rows {
		for functionCol, functions := range functionColumns {
			for _, function := range functions {
				if strings.ToUpper(function) == "AVG" {
					sumKey := functionCol + "SUM"
					countKey := functionCol + "COUNT"
					sumVal := newTable.Rows[i][sumKey].(float64)
					countVal := newTable.Rows[i][countKey].(float64)
					if countVal == 0 {
						newTable.Rows[i][functionCol+"AVG"] = 0.0
					} else {
						newTable.Rows[i][functionCol+"AVG"] = sumVal / countVal
					}
				}
			}
		}
	}

	return newTable

}
