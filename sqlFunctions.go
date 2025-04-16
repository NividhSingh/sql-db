package main

import (
	"fmt"
	"strings"
)

type Column struct {
	Name       string
	Type       string
	Conditions []string
}

type Table struct {
	Name    string
	Columns []Column
	Rows    []map[string]interface{}
}

var database = make(map[string]Table)

/*
func processSelectFromTable(command string) Table {
	tokens := strings.Fields(command)
	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "SELECT" {
		fmt.Println("Invalid SELECT command")
		return Table{}
	}

	fromIndex := indexOf(tokens, "FROM")
	if fromIndex == -1 {
		fmt.Println("Invalid SELECT command: missing FROM")
		return Table{}
	}

	tableName := tokens[fromIndex+1]
	if !tableExists(tableName) {
		fmt.Printf("Table %s does not exist\n", tableName)
		return Table{}
	}

	// Determine which columns to include.
	var selectedColumns []Column
	if tokens[1] == "*" {
		selectedColumns = database[tableName].Columns
	} else {
		for _, colName := range tokens[1:fromIndex] {
			for _, col := range database[tableName].Columns {
				if col.Name == strings.Trim(colName, ",") {
					selectedColumns = append(selectedColumns, col)
				}
			}
		}
	}

	// Create a temporary result table.
	resultTable := Table{
		Name:    "result",
		Columns: selectedColumns,
	}

	// For simplicity, we'll only implement a basic WHERE condition with the "<" operator.
	var whereCondition func(row map[string]interface{}) bool
	whereIndex := indexOf(tokens, "WHERE")
	if whereIndex != -1 && len(tokens) > whereIndex+3 {
		// Expect format: WHERE column < value
		colName := tokens[whereIndex+1]
		operator := tokens[whereIndex+2]
		valueStr := tokens[whereIndex+3]
		whereCondition = func(row map[string]interface{}) bool {
			// Look up the column in the original table.
			for _, col := range database[tableName].Columns {
				if col.Name == colName {
					switch strings.ToUpper(col.Type) {
					case "INT":
						rowVal, ok := row[colName].(int)
						if !ok {
							return false
						}
						condVal, err := strconv.Atoi(valueStr)
						if err != nil {
							return false
						}
						if operator == "<" {
							return rowVal < condVal
						} else if operator == ">" {
							return rowVal > condVal
						}
					case "FLOAT":
						rowVal, ok := row[colName].(float64)
						if !ok {
							return false
						}
						condVal, err := strconv.ParseFloat(valueStr, 64)
						if err != nil {
							return false
						}
						if operator == "<" {
							return rowVal < condVal
						} else if operator == ">" {
							return rowVal > condVal
						}
					default:
						// For VARCHAR or TEXT, compare as strings.
						rowVal := fmt.Sprintf("%v", row[colName])
						if operator == "<" {
							return rowVal < valueStr
						} else if operator == ">" {
							return rowVal > valueStr
						}
					}
				}
			}
			return false
		}
	} else {
		// No WHERE clause: include all rows.
		whereCondition = func(row map[string]interface{}) bool { return true }
	}

	// Process rows in the original table.
	for _, row := range database[tableName].Rows {
		if whereCondition(row) {
			newRow := make(map[string]interface{})
			for _, col := range selectedColumns {
				newRow[col.Name] = row[col.Name]
			}
			resultTable.Rows = append(resultTable.Rows, newRow)
		}
	}
	return resultTable
}
*/
// printTable prints a single table in grid format.

func printTable(table Table) {
	if len(table.Columns) == 0 {
		fmt.Println("Empty result")
		return
	}
	cols := make([]string, len(table.Columns))
	widths := make([]int, len(table.Columns))
	for i, col := range table.Columns {
		cols[i] = col.Name
		widths[i] = len(col.Name)
	}
	for _, row := range table.Rows {
		for i, col := range table.Columns {
			cellStr := fmt.Sprintf("%v", row[col.Name])
			if len(cellStr) > widths[i] {
				widths[i] = len(cellStr)
			}
		}
	}
	sep := "+"
	for _, w := range widths {
		sep += strings.Repeat("-", w+2) + "+"
	}
	header := "|"
	for i, name := range cols {
		header += " " + fmt.Sprintf("%-*s", widths[i], name) + " |"
	}
	fmt.Println(sep)
	fmt.Println(header)
	fmt.Println(sep)
	for _, row := range table.Rows {
		rowStr := "|"
		for i, col := range table.Columns {
			cellStr := fmt.Sprintf("%v", row[col.Name])
			rowStr += " " + fmt.Sprintf("%-*s", widths[i], cellStr) + " |"
		}
		fmt.Println(rowStr)
	}
	fmt.Println(sep)
	fmt.Println("")
}

// Existing helper functions

// func processInsertIntoTable(command string) {
// 	tokens := strings.Fields(command)
// 	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "INSERT" || strings.ToUpper(tokens[1]) != "INTO" {
// 		fmt.Println("Invalid command")
// 		return
// 	}
// 	tableName := tokens[2]
// 	if !tableExists(tableName) {
// 		fmt.Printf("Table %s does not exist\n", tableName)
// 		return
// 	}
// 	command = strings.Join(tokens[3:], " ")
// 	command = strings.TrimSpace(command)
// 	newColumnNames := strings.Split(getNextInParenthis(command), ",")
// 	command = command[len(getNextInParenthis(command))+3:] // Remove the column names part
// 	command = strings.TrimSpace(command)
// 	tokens = strings.Fields(command)
// 	if len(tokens) < 2 || strings.ToUpper(tokens[0]) != "VALUES" {
// 		fmt.Println("Invalid command")
// 		return
// 	}
// 	command = getNextInParenthis(strings.Join(tokens[1:], " "))
// 	newColumnValues := strings.Split(command, ",")
// 	for i, v := range newColumnValues {
// 		newColumnValues[i] = strings.TrimSpace(v)
// 		newColumnNames[i] = strings.TrimSpace(newColumnNames[i])
// 	}
// 	table := database[tableName]
// 	newRow := make(map[string]interface{})
// 	for _, col := range table.Columns {
// 		for j, column := range newColumnNames {
// 			if col.Name == column {
// 				for _, constraint := range col.Conditions {
// 					if constraint == "PRIMARY KEY" || constraint == "UNIQUE" {
// 						if !checkUnique(tableName, col.Name, newColumnValues[j]) {
// 							fmt.Printf("Column %s has to be unique\n", col.Name)
// 							return
// 						}
// 					}
// 					if strings.HasPrefix(constraint, "DEFAULT") {
// 						if newColumnValues[j] == "" {
// 							newColumnValues[j] = strings.TrimPrefix(constraint, "DEFAULT ")
// 						}
// 					}
// 					if constraint == "NOT NULL" {
// 						if newColumnValues[j] == "" {
// 							fmt.Printf("Column %s cannot be NULL\n", col.Name)
// 							return
// 						}
// 					}
// 				}
// 				switch strings.ToUpper(col.Type) {
// 				case "INT":
// 					if newColumnValues[j] != "" {
// 						if _, err := strconv.Atoi(newColumnValues[j]); err != nil {
// 							fmt.Printf("Column %s must be an integer\n", col.Name)
// 							return
// 						} else {
// 							newRow[col.Name], _ = strconv.Atoi(newColumnValues[j])
// 						}
// 					}
// 				case "FLOAT":
// 					if newColumnValues[j] != "" {
// 						floatVal, err := strconv.ParseFloat(newColumnValues[j], 64)
// 						if err != nil {
// 							fmt.Printf("Column %s must be a float\n", col.Name)
// 							return
// 						} else {
// 							newRow[col.Name] = floatVal
// 						}
// 					} else {
// 						newRow[col.Name] = nil
// 					}
// 				default:
// 					if strings.HasPrefix(strings.ToUpper(col.Type), "VARCHAR") {
// 						maxLength, _ := strconv.Atoi(getNextInParenthis(col.Type))
// 						if len(newColumnValues[j]) > maxLength {
// 							fmt.Printf("Column %s exceeds maximum length of %d\n", col.Name, maxLength)
// 							return
// 						}
// 						newRow[col.Name] = newColumnValues[j]
// 					}
// 				}
// 			}
// 		}
// 	}
// 	table.Rows = append(table.Rows, newRow)
// 	database[tableName] = table
// 	fmt.Printf("Inserted row into %s\n", tableName)
// }

func checkUnique(tableName string, columnName string, value string) bool {
	for _, row := range database[tableName].Rows {
		if row[columnName] == value {
			return false
		}
	}
	return true
}

func createColumn(newName string, newType string, newConditions []string) Column {
	return Column{Name: newName, Type: newType, Conditions: newConditions}
}

func processCreateTableCommand(newName string, newColumns []Column) {
	if tableExists(newName) {
		// Error table exists
		return
	}
	database[newName] = Table{
		Name:    newName,
		Columns: newColumns,
	}
}

func printDatabase() {
	fmt.Println("Current Database State:")
	for tableName, table := range database {
		fmt.Printf("Table: %s\n", tableName)
		cols := make([]string, len(table.Columns))
		widths := make([]int, len(table.Columns))
		for i, col := range table.Columns {
			cols[i] = col.Name
			widths[i] = len(col.Name)
		}
		for _, row := range table.Rows {
			for i, col := range table.Columns {
				cellStr := fmt.Sprintf("%v", row[col.Name])
				if len(cellStr) > widths[i] {
					widths[i] = len(cellStr)
				}
			}
		}
		sep := "+"
		for _, w := range widths {
			sep += strings.Repeat("-", w+2) + "+"
		}
		header := "|"
		for i, name := range cols {
			header += " " + fmt.Sprintf("%-*s", widths[i], name) + " |"
		}
		fmt.Println(sep)
		fmt.Println(header)
		fmt.Println(sep)
		for _, row := range table.Rows {
			rowStr := "|"
			for i, col := range table.Columns {
				cellStr := fmt.Sprintf("%v", row[col.Name])
				rowStr += " " + fmt.Sprintf("%-*s", widths[i], cellStr) + " |"
			}
			fmt.Println(rowStr)
		}
		fmt.Println(sep)
		fmt.Println("")
	}
}

func tableExists(tableName string) bool {
	_, exists := database[tableName]
	return exists
}

func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}
