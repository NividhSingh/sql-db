package main

import (
	"fmt"
	"math"
	"os"
)

const (
	// total budget of ε across all SELECTs:
	maxEpsilonBudget = 10.0
	// decayRate r: each query’s ε_n is multiplied by r relative to the prior
	decayRate = 0.5
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

	selectCount := 0

	for _, astNode := range astNodes {
		switch astNode.Type {
		case AST_CREATE:
			createTableFromAST(astNode)

		case AST_INSERT:
			insertIntoFromAST(astNode)

		case AST_SELECT:
			selectCount++
			// compute this query’s ε_n
			epsilon := maxEpsilonBudget * (1.0 - decayRate) * math.Pow(decayRate, float64(selectCount-1))
			sensitivity := 1.0

			result := selectFromAST(astNode)

			// add Laplace noise with ε_n
			for _, col := range result.Columns {
				if col.FunctionResult {
					for _, row := range result.Rows {
						if v, ok := row[col.Name].(float64); ok {
							row[col.Name] = addNoise(v, epsilon, sensitivity)
						}
					}
				}
			}

			// your k‑anonymity & l‑diversity calls
			result = enforceKAnonymity(result, []string{"blood_type", "male_or_female"}, 10)
			result = enforceLDiversity(result, []string{"has_diabetes", "sex"}, 3)

			fmt.Printf("\n-- SELECT #%d: ε=%.4f  (cumulative budget used ≈ %.4f)\n",
				selectCount, epsilon,
				maxEpsilonBudget*(1-math.Pow(decayRate, float64(selectCount))))
			printTable(result)
		}
	}
}
