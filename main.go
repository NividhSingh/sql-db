package main

func main() {
	// Input command string with CREATE, INSERT, and SELECT commands.
	command := "CREATE TABLE myTable (col1 VARCHAR (255) PRIMARY KEY, col2 INT); " +
		"INSERT INTO myTable (col1, col2) VALUES ('John', 42); " +
		"SELECT col1 AS c1, col2 FROM myTable GROUP BY c1;"

	// Initialize the lexer with the command string.
	lexer := &Lexer{
		input:   command,
		start:   0,
		current: 0,
		line:    1,
	}

	// Tokenize the input.
	var tokens []*Token
	for {
		token := getNextToken(lexer)
		tokens = append(tokens, &token)
		if token._type == TOKEN_EOF {
			break
		}
	}

	// Parse the tokens into AST nodes.
	astNodes := parseCommands(tokens)

	for _, node := range astNodes {
		printAST(node, 0)
	}

	for _, astNode := range astNodes {
		if astNode.Type == AST_CREATE {
			createTableFromAST(astNode)
		} else if astNode.Type == AST_INSERT {
			insertIntoFromAST(astNode)
		} else if astNode.Type == AST_SELECT {
			// selectFromAST(astNode)
		}
	}
	printDatabase()

}
