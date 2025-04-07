// --------------------------------------------------
//

// --------------------------------------------------

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
	createTable(command1)
	printDatabase()

	// Test 2: Create a table with conditions and the IF NOT EXISTS clause.
	command2 := "CREATE TABLE IF NOT EXISTS products (product_id INT PRIMARY KEY, description VARCHAR(1000) NOT NULL, price FLOAT DEFAULT 0)"
	fmt.Println("\nTest 2 command:")
	fmt.Println(command2)
	createTable(command2)
	printDatabase()

	// Test 3: Attempt to create a table that already exists without IF NOT EXISTS.
	command3 := "CREATE TABLE users (id INT, email VARCHAR(1000) UNIQUE)"
	fmt.Println("\nTest 3 command:")
	fmt.Println(command3)
	createTable(command3)
	printDatabase()

	// ----------------------------------------------------------------
	// Insert tests for insertIntoTable function.
	fmt.Println("\n----- Running tests for insertIntoTable -----")

	// Test Insert 1: Insert a valid row into the "users" table.
	insertCmd1 := "INSERT INTO users (id, name) VALUES (1, Alice)"
	fmt.Println("\nTest Insert 1 command:")
	fmt.Println(insertCmd1)
	insertIntoTable(insertCmd1)
	printDatabase()

	// Test Insert 2: Insert another valid row into the "users" table.
	insertCmd2 := "INSERT INTO users (id, name) VALUES (2, Bob)"
	fmt.Println("\nTest Insert 2 command:")
	fmt.Println(insertCmd2)
	insertIntoTable(insertCmd2)
	printDatabase()

	// Test Insert 3: Insert a row with duplicate primary key in "products".
	// "products" table has product_id as PRIMARY KEY.
	insertCmd3 := "INSERT INTO products (product_id, description, price) VALUES (101, Widget, 9.99)"
	fmt.Println("\nTest Insert 3 command (first insert into products):")
	fmt.Println(insertCmd3)
	insertIntoTable(insertCmd3)
	printDatabase()

	// Try inserting a duplicate product_id.
	insertCmd4 := "INSERT INTO products (product_id, description, price) VALUES (101, Gadget, 8.99)"
	fmt.Println("\nTest Insert 4 command (duplicate PRIMARY KEY):")
	fmt.Println(insertCmd4)
	insertIntoTable(insertCmd4)
	printDatabase()

	// Test Insert 5: Insert a row missing a NOT NULL value.
	// description in products is NOT NULL.
	insertCmd5 := "INSERT INTO products (product_id, description, price) VALUES (102, , 19.99)"
	fmt.Println("\nTest Insert 5 command (missing NOT NULL value for description):")
	fmt.Println(insertCmd5)
	insertIntoTable(insertCmd5)
	printDatabase()

	// Test Insert 6: Insert a row missing a value for price; should use DEFAULT.
	insertCmd6 := "INSERT INTO products (product_id, description, price) VALUES (103, Gadget, )"
	fmt.Println("\nTest Insert 6 command (missing value for price, should apply DEFAULT):")
	fmt.Println(insertCmd6)
	insertIntoTable(insertCmd6)
	printDatabase()
}
func insertIntoTable(command string) {

	tokens := strings.Fields(command)
	if len(tokens) < 4 || tokens[0] != "INSERT" || tokens[1] != "INTO" {
		fmt.Println("Invalid command")
		return
	}
	tableName := tokens[2]

	if _, exists := database[tableName]; !exists {
		fmt.Printf("Table %s does not exist\n", tableName)
		return
	}

	command = strings.Join(tokens[3:], " ")
	command = strings.TrimSpace(command)

	new_column_names := strings.Split(getNextInParenthis(command), ",")

	command = command[len(getNextInParenthis(command))+3:] // Remove the column names part

	command = strings.TrimSpace(command)
	tokens = strings.Fields(command)
	if len(tokens) < 2 || tokens[0] != "VALUES" {
		fmt.Println("Invalid command")
		return
	}
	command = getNextInParenthis(strings.Join(tokens[1:], ""))
	new_column_values := strings.Split(command, ",")

	for i, v := range new_column_values {
		new_column_values[i] = strings.TrimSpace(v)
		new_column_names[i] = strings.TrimSpace(new_column_names[i])
	}

	table := database[tableName]

	newRow := make(map[string]interface{})

	for _, col := range table.Columns {
		for j, column := range new_column_names {
			if col.Name == column {

				constraints := col.Conditions
				for _, constraint := range constraints {
					if constraint == "PRIMARY KEY" || constraint == "UNIQUE" {
						if !checkUnique(tableName, col.Name, new_column_values[j]) {
							fmt.Printf("Column %s has to be unique\n", col.Name)
							return
						}
					}
					if strings.HasPrefix(constraint, "DEFAULT") {
						if new_column_values[j] == "" {
							new_column_values[j] = strings.TrimPrefix(constraint, "DEFAULT ")
						}
					}
					if constraint == "NOT NULL" {
						if new_column_values[j] == "" {
							fmt.Printf("Column %s cannot be NULL\n", col.Name)
							return
						}
					}
				}

				switch strings.ToUpper(col.Type) {
				case "INT":
					if new_column_values[j] != "" {
						if _, err := strconv.Atoi(new_column_values[j]); err != nil {
							fmt.Printf("Column %s must be an integer\n", col.Name)
							return
						} else {
							newRow[col.Name], _ = strconv.Atoi(new_column_values[j])
						}
					}
				case "FLOAT":
					if new_column_values[j] != "" {
						floatVal, err := strconv.ParseFloat(new_column_values[j], 64)
						if err != nil {
							fmt.Printf("Column %s must be a float\n", col.Name)
							return
						} else {
							newRow[col.Name] = floatVal
						}
					} else {
						newRow[col.Name] = nil
					}
				// Add boolean and other types as needed
				default:
					if strings.HasPrefix(strings.ToUpper(col.Type), "VARCHAR") {
						maxLength, _ := strconv.Atoi(getNextInParenthis(col.Type))
						if len(new_column_values[j]) > maxLength {
							fmt.Printf("Column %s exceeds maximum length of %d\n", col.Name, maxLength)
							return
						}
						newRow[col.Name] = new_column_values[j]
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

func createTable(command string) {
	tokens := strings.Fields(command)
	if len(tokens) < 4 || tokens[0] != "CREATE" || tokens[1] != "TABLE" {
		fmt.Println("Invalid command")
		return
	}
	// tableName := tokens[2]

	var tableName string
	var colDefStartIndex int
	if strings.ToUpper(tokens[2]) == "IF" {
		// Expect the tokens to be: IF NOT EXISTS table_name
		if len(tokens) < 6 ||
			strings.ToUpper(tokens[3]) != "NOT" ||
			strings.ToUpper(tokens[4]) != "EXISTS" {
			fmt.Printf("invalid syntax for IF NOT EXISTS clause")
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
			fmt.Printf("table %s already exists", tableName)
			return
		}
	}

	rest := strings.Join(tokens[colDefStartIndex:], " ")
	rest = strings.TrimSpace(rest)

	if !strings.HasPrefix(rest, "(") || !strings.HasSuffix(rest, ")") {
		fmt.Printf("invalid syntax for column definitions")
		return
	}

	// rest = rest[1 : len(rest)-1] // Remove the parentheses
	rest = getNextInParenthis(rest)

	values := strings.Split(rest, ",")

	var columns []Column

	for _, v := range values {
		v = strings.TrimSpace(v)
		v := strings.Split(v, " ")
		if len(v) < 2 {
			fmt.Println("Invalid command")
			return
		}

		name := strings.TrimSpace(v[0])
		typeName := strings.TrimSpace(v[1])
		var conditions []string

		v = v[2:]

		for len(v) > 0 {
			if v[0] == "PRIMARY" {
				if v[1] != "KEY" {
					fmt.Println("Invalid command")
					return
				}
				conditions = append(conditions, "PRIMARY KEY")
				v = v[2:]
			} else if v[0] == "NOT" {
				if v[1] != "NULL" {
					fmt.Println("Invalid command")
					return
				}
				conditions = append(conditions, "NOT NULL")
				v = v[2:]
			} else if v[0] == "UNIQUE" {
				conditions = append(conditions, "UNIQUE")
				v = v[1:]
			} else if v[0] == "DEFAULT" {
				if len(v) < 2 {
					fmt.Println("Invalid command")
					return
				}
				conditions = append(conditions, "DEFAULT "+v[1])
				v = v[2:]
			}
		}
		// Add in check and foreign key eventually
		//helper function to get stuff in parenthesis maybe
		newColumn := Column{
			Name:       name,
			Type:       typeName,
			Conditions: conditions,
		}
		for _, column := range columns {
			if column.Name == newColumn.Name {
				fmt.Printf("Duplicate column name: %s", newColumn.Name)
				return
			}
		}

		columns = append(columns, newColumn)
	}
	database[tableName] = Table{
		Name:    tableName,
		Columns: columns,
	}
	return
}

func getNextInParenthis(s string) string {
	// Find the first opening parenthesis
	start := strings.Index(s, "(")
	if start == -1 {
		return ""
	}

	// Find the corresponding closing parenthesis
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

// printDatabase outputs the current state of the in-memory database.
func printDatabase() {
	fmt.Println("Current Database State:")
	for tableName, table := range database {
		fmt.Printf("Table: %s\n", tableName)

		// Build header row based on table.Columns.
		cols := make([]string, len(table.Columns))
		widths := make([]int, len(table.Columns))
		for i, col := range table.Columns {
			cols[i] = col.Name
			widths[i] = len(col.Name)
		}

		// Update widths with the length of each cell value (converted to string).
		for _, row := range table.Rows {
			for i, col := range table.Columns {
				cellStr := fmt.Sprintf("%v", row[col.Name])
				if len(cellStr) > widths[i] {
					widths[i] = len(cellStr)
				}
			}
		}

		// Build horizontal separator.
		sep := "+"
		for _, w := range widths {
			sep += strings.Repeat("-", w+2) + "+"
		}

		// Print header row.
		header := "|"
		for i, name := range cols {
			header += " " + fmt.Sprintf("%-*s", widths[i], name) + " |"
		}

		fmt.Println(sep)
		fmt.Println(header)
		fmt.Println(sep)
		// Print each row.
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
