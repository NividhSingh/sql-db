package main

import (
	"fmt"
	"os"
)

func main() {
	// Read SQL commands from file
	data, err := os.ReadFile("input.sql")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	command := string(data)

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

	printDatabase()

	for _, astNode := range astNodes {
		if astNode.Type == AST_CREATE {
			createTableFromAST(astNode)
		} else if astNode.Type == AST_INSERT {
			insertIntoFromAST(astNode)
		} else if astNode.Type == AST_SELECT {
			result := selectFromAST(astNode)
			printTable(result)
		}
	}
}
