package main

import (
	"fmt"
	"strconv"
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

func main() {
	fmt.Println("Running tests for createTable...")

	// Test 1: Create a table without conditions.
	command1 := "CREATE TABLE users (id INT, name VARCHAR(1000))"
	fmt.Println("\nTest 1 command:")
	fmt.Println(command1)
	processCreateTableCommand(command1)
	printDatabase()

	// Test 2: Create a table with conditions and the IF NOT EXISTS clause.
	command2 := "CREATE TABLE IF NOT EXISTS products (product_id INT PRIMARY KEY, description VARCHAR(1000) NOT NULL, price FLOAT DEFAULT 0)"
	fmt.Println("\nTest 2 command:")
	fmt.Println(command2)
	processCreateTableCommand(command2)
	printDatabase()

	// Test 3: Attempt to create a table that already exists without IF NOT EXISTS.
	command3 := "CREATE TABLE users (id INT, email VARCHAR(1000) UNIQUE)"
	fmt.Println("\nTest 3 command:")
	fmt.Println(command3)
	processCreateTableCommand(command3)
	printDatabase()

	// ----------------------------------------------------------------
	// Insert tests for processInsertIntoTable function.
	fmt.Println("\n----- Running tests for insertIntoTable -----")

	// Test Insert 1: Insert a valid row into the "users" table.
	insertCmd1 := "INSERT INTO users (id, name) VALUES (1, Alice)"
	fmt.Println("\nTest Insert 1 command:")
	fmt.Println(insertCmd1)
	processInsertIntoTable(insertCmd1)
	printDatabase()

	// Test Insert 2: Insert another valid row into the "users" table.
	insertCmd2 := "INSERT INTO users (id, name) VALUES (2, Bob)"
	fmt.Println("\nTest Insert 2 command:")
	fmt.Println(insertCmd2)
	processInsertIntoTable(insertCmd2)
	printDatabase()

	// Test Insert 3: Insert a row into the "products" table.
	insertCmd3 := "INSERT INTO products (product_id, description, price) VALUES (101, Widget, 9.99)"
	fmt.Println("\nTest Insert 3 command (first insert into products):")
	fmt.Println(insertCmd3)
	processInsertIntoTable(insertCmd3)
	printDatabase()

	// Test Insert 4: Attempt to insert a duplicate primary key into "products".
	insertCmd4 := "INSERT INTO products (product_id, description, price) VALUES (101, Gadget, 8.99)"
	fmt.Println("\nTest Insert 4 command (duplicate PRIMARY KEY):")
	fmt.Println(insertCmd4)
	processInsertIntoTable(insertCmd4)
	printDatabase()

	// Test Insert 5: Insert a row missing a NOT NULL value.
	insertCmd5 := "INSERT INTO products (product_id, description, price) VALUES (102, , 19.99)"
	fmt.Println("\nTest Insert 5 command (missing NOT NULL value for description):")
	fmt.Println(insertCmd5)
	processInsertIntoTable(insertCmd5)
	printDatabase()

	// Test Insert 6: Insert a row missing a value for price; should use DEFAULT.
	insertCmd6 := "INSERT INTO products (product_id, description, price) VALUES (103, Gadget, )"
	fmt.Println("\nTest Insert 6 command (missing value for price, should apply DEFAULT):")
	fmt.Println(insertCmd6)
	processInsertIntoTable(insertCmd6)
	printDatabase()

	// ----------------------------------------------------------------
	// New tests for the SELECT query function.
	fmt.Println("\n----- Running tests for SELECT query -----")

	// Test Select 1: Select all columns from the users table.
	selectCmd1 := "SELECT * FROM users"
	fmt.Println("\nTest Select 1 command:")
	fmt.Println(selectCmd1)
	resultTable1 := processSelectFromTable(selectCmd1)
	printTable(resultTable1)

	// Test Select 2: Select only the 'id' column from the users table.
	selectCmd2 := "SELECT id FROM users"
	fmt.Println("\nTest Select 2 command:")
	fmt.Println(selectCmd2)
	resultTable2 := processSelectFromTable(selectCmd2)
	printTable(resultTable2)

	// Test Select 3: Select product_id and price from products where price is less than 10.
	selectCmd3 := "SELECT product_id, price FROM products WHERE price < 10"
	fmt.Println("\nTest Select 3 command:")
	fmt.Println(selectCmd3)
	resultTable3 := processSelectFromTable(selectCmd3)
	printTable(resultTable3)
}

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

func processInsertIntoTable(command string) {
	tokens := strings.Fields(command)
	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "INSERT" || strings.ToUpper(tokens[1]) != "INTO" {
		fmt.Println("Invalid command")
		return
	}
	tableName := tokens[2]
	if !tableExists(tableName) {
		fmt.Printf("Table %s does not exist\n", tableName)
		return
	}
	command = strings.Join(tokens[3:], " ")
	command = strings.TrimSpace(command)
	newColumnNames := strings.Split(getNextInParenthis(command), ",")
	command = command[len(getNextInParenthis(command))+3:] // Remove the column names part
	command = strings.TrimSpace(command)
	tokens = strings.Fields(command)
	if len(tokens) < 2 || strings.ToUpper(tokens[0]) != "VALUES" {
		fmt.Println("Invalid command")
		return
	}
	command = getNextInParenthis(strings.Join(tokens[1:], " "))
	newColumnValues := strings.Split(command, ",")
	for i, v := range newColumnValues {
		newColumnValues[i] = strings.TrimSpace(v)
		newColumnNames[i] = strings.TrimSpace(newColumnNames[i])
	}
	table := database[tableName]
	newRow := make(map[string]interface{})
	for _, col := range table.Columns {
		for j, column := range newColumnNames {
			if col.Name == column {
				for _, constraint := range col.Conditions {
					if constraint == "PRIMARY KEY" || constraint == "UNIQUE" {
						if !checkUnique(tableName, col.Name, newColumnValues[j]) {
							fmt.Printf("Column %s has to be unique\n", col.Name)
							return
						}
					}
					if strings.HasPrefix(constraint, "DEFAULT") {
						if newColumnValues[j] == "" {
							newColumnValues[j] = strings.TrimPrefix(constraint, "DEFAULT ")
						}
					}
					if constraint == "NOT NULL" {
						if newColumnValues[j] == "" {
							fmt.Printf("Column %s cannot be NULL\n", col.Name)
							return
						}
					}
				}
				switch strings.ToUpper(col.Type) {
				case "INT":
					if newColumnValues[j] != "" {
						if _, err := strconv.Atoi(newColumnValues[j]); err != nil {
							fmt.Printf("Column %s must be an integer\n", col.Name)
							return
						} else {
							newRow[col.Name], _ = strconv.Atoi(newColumnValues[j])
						}
					}
				case "FLOAT":
					if newColumnValues[j] != "" {
						floatVal, err := strconv.ParseFloat(newColumnValues[j], 64)
						if err != nil {
							fmt.Printf("Column %s must be a float\n", col.Name)
							return
						} else {
							newRow[col.Name] = floatVal
						}
					} else {
						newRow[col.Name] = nil
					}
				default:
					if strings.HasPrefix(strings.ToUpper(col.Type), "VARCHAR") {
						maxLength, _ := strconv.Atoi(getNextInParenthis(col.Type))
						if len(newColumnValues[j]) > maxLength {
							fmt.Printf("Column %s exceeds maximum length of %d\n", col.Name, maxLength)
							return
						}
						newRow[col.Name] = newColumnValues[j]
					}
				}
			}
		}
	}
	table.Rows = append(table.Rows, newRow)
	database[tableName] = table
	fmt.Printf("Inserted row into %s\n", tableName)
}

func checkUnique(tableName string, columnName string, value string) bool {
	for _, row := range database[tableName].Rows {
		if row[columnName] == value {
			return false
		}
	}
	return true
}

func processCreateTableCommand(command string) {
	tokens := strings.Fields(command)
	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "CREATE" || strings.ToUpper(tokens[1]) != "TABLE" {
		fmt.Println("Invalid command")
		return
	}
	var tableName string
	var colDefStartIndex int
	if strings.ToUpper(tokens[2]) == "IF" {
		if len(tokens) < 6 || strings.ToUpper(tokens[3]) != "NOT" || strings.ToUpper(tokens[4]) != "EXISTS" {
			fmt.Println("Invalid syntax for IF NOT EXISTS clause")
			return
		}
		tableName = tokens[5]
		colDefStartIndex = 6
		if tableExists(tableName) {
			fmt.Printf("Table %s already exists\n", tableName)
			return
		}
	} else {
		tableName = tokens[2]
		colDefStartIndex = 3
		if tableExists(tableName) {
			fmt.Printf("Table %s already exists\n", tableName)
			return
		}
	}
	rest := strings.Join(tokens[colDefStartIndex:], " ")
	rest = strings.TrimSpace(rest)
	if !strings.HasPrefix(rest, "(") || !strings.HasSuffix(rest, ")") {
		fmt.Println("Invalid syntax for column definitions")
		return
	}
	rest = getNextInParenthis(rest)
	values := strings.Split(rest, ",")
	var columns []Column
	for _, v := range values {
		v = strings.TrimSpace(v)
		parts := strings.Split(v, " ")
		if len(parts) < 2 {
			fmt.Println("Invalid command")
			return
		}
		name := strings.TrimSpace(parts[0])
		typeName := strings.TrimSpace(parts[1])
		var conditions []string
		parts = parts[2:]
		for len(parts) > 0 {
			if parts[0] == "PRIMARY" {
				if len(parts) < 2 || parts[1] != "KEY" {
					fmt.Println("Invalid command")
					return
				}
				conditions = append(conditions, "PRIMARY KEY")
				parts = parts[2:]
			} else if parts[0] == "NOT" {
				if len(parts) < 2 || parts[1] != "NULL" {
					fmt.Println("Invalid command")
					return
				}
				conditions = append(conditions, "NOT NULL")
				parts = parts[2:]
			} else if parts[0] == "UNIQUE" {
				conditions = append(conditions, "UNIQUE")
				parts = parts[1:]
			} else if parts[0] == "DEFAULT" {
				if len(parts) < 2 {
					fmt.Println("Invalid command")
					return
				}
				conditions = append(conditions, "DEFAULT "+parts[1])
				parts = parts[2:]
			} else {
				parts = parts[1:]
			}
		}
		newColumn := Column{
			Name:       name,
			Type:       typeName,
			Conditions: conditions,
		}
		for _, column := range columns {
			if column.Name == newColumn.Name {
				fmt.Printf("Duplicate column name: %s\n", newColumn.Name)
				return
			}
		}
		columns = append(columns, newColumn)
	}
	database[tableName] = Table{
		Name:    tableName,
		Columns: columns,
	}
}

func getNextInParenthis(s string) string {
	start := strings.Index(s, "(")
	if start == -1 {
		return ""
	}
	count := 1
	for i := start + 1; i < len(s); i++ {
		if s[i] == '(' {
			count++
		} else if s[i] == ')' {
			count--
			if count == 0 {
				return s[start+1 : i]
			}
		}
	}
	return s
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
