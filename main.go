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

	for _, astNode := range astNodes {
		if astNode.Type == AST_CREATE {
			createTableFromAST(astNode)
		} else if astNode.Type == AST_INSERT {
			insertIntoFromAST(astNode)
		} else if astNode.Type == AST_SELECT {
			result := selectFromAST(astNode)
			epsilon := 2.0
			sensitivity := 1.0
			// ─── add Laplace noise to all aggregate columns ─────────────────────────
			for _, col := range result.Columns {
				// FunctionResult marks SUM, AVG, COUNT, MIN, MAX, etc.
				if col.FunctionResult {
					for _, row := range result.Rows {
						// pull the raw float64 value
						if v, ok := row[col.Name].(float64); ok {
							row[col.Name] = addNoise(v, epsilon, sensitivity)
						}
					}
				}
			}

			result = enforceKAnonymity(result, []string{"blood_type", "male_or_female"}, 10)
			result = enforceLDiversity(result, []string{"has_diabetes", "sex"}, 3)

			printTable(result)
		}
	}
}
