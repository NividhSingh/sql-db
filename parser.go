package main

import (
	"fmt"
	"strconv"
	"strings"
)

// --- Helper for token type comparison ---

func checkType(token *Token, _token_type TokenType) bool {
	return token._type == _token_type
}

func panicIfWrongType(token *Token, _token_type TokenType) {
	fmt.Println("Token value:", token.value)
	fmt.Println("Token type:", tokenTypeToString(token._type))
	fmt.Println("Expected type:", tokenTypeToString(_token_type))
	if !checkType(token, _token_type) {
		panic("Wrong token type")
	}
}

// --- AST definitions ---

type ASTNodeType int

const (
	AST_INT_LITERAL ASTNodeType = iota
	AST_VARCHAR_LITERAL
	AST_BOOLEAN_LITERAL
	AST_FLOAT_LITERAL
	AST_SELECT
	AST_CREATE
	AST_INSERT
	AST_FUNCTION
	AST_EXPRESSION
	AST_VALUE
	AST_COLUMN_NAME
	AST_BINARY
	AST_COLUMN
)

type columnType int

const (
	COLUMN_TYPE_NOT_INCLUDED columnType = iota
	COLUMN_TYPE_GROUP_BY
	COLUMN_TYPE_NORMAL
	COLUMN_TYPE_MAX
	COLUMN_TYPE_MIN
	COLUMN_TYPE_AVG
	COLUMN_TYPE_COUNT
	COLUMN_TYPE_SUM
)

type ASTNode struct {
	// Common fields
	Type      ASTNodeType
	IntVal    int    // (unused)
	StrVal    string // (unused)
	floatVal  float64
	boolValue bool

	// Expression node (for identifiers and literals)
	columnName string

	left     *ASTNode
	right    *ASTNode
	operator string

	// For literal numbers and strings (alternative fields)
	intVal int64
	strVal string

	// Function node
	functionName       string
	functionArguements []*ASTNode

	// Column node (used in CREATE)
	name         string
	_type        string
	varCharLimit int
	constraints  []string

	// Select node
	columnNames   []string
	columnTypes   []columnType
	columnAliases []string

	containsGroupBy bool
	groupByColumns  []string
	containsLimit   bool
	limit           int
	containsOffset  bool
	offset          int

	// Create node
	tableName string
	columns   []*ASTNode

	// Insert node
	columnValues []string
}

// --- Functions used by parser ---

func endExpressionValues(token *Token) bool {
	return token._type == TOKEN_COMMA ||
		token._type == TOKEN_FROM ||
		token._type == TOKEN_AS ||
		token._type == TOKEN_RPAREN
}

func functionValues(token *Token) bool {
	return token._type == TOKEN_COUNT ||
		token._type == TOKEN_SUM ||
		token._type == TOKEN_AVG ||
		token._type == TOKEN_MIN ||
		token._type == TOKEN_MAX
}

// --- Parsing functions ---

func parseExpression(tokens []*Token, tokenIndex *int) *ASTNode {
	// Handle function expressions first.
	if functionValues(tokens[*tokenIndex]) {
		newFunctionName := tokens[*tokenIndex].value
		(*tokenIndex)++ // Move past function name
		panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
		(*tokenIndex)++ // Move past LPAREN
		var newFunctionArguements []*ASTNode = make([]*ASTNode, 0)
		for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			newFunctionArguements = append(newFunctionArguements, parseExpression(tokens, tokenIndex))
		}
		(*tokenIndex)++ // Move past RPAREN
		expressionNodeForFunction := ASTNode{
			Type:               AST_FUNCTION,
			functionName:       newFunctionName,
			functionArguements: newFunctionArguements,
		}
		if !endExpressionValues(tokens[*tokenIndex]) {
			(*tokenIndex)++ // Move past operator
			return &ASTNode{
				Type:     AST_BINARY,
				left:     &expressionNodeForFunction,
				operator: tokens[(*tokenIndex)-1].value,
				right:    parseExpression(tokens, tokenIndex),
			}
		} else {
			return &expressionNodeForFunction
		}
	}

	var leftNode *ASTNode = nil

	// Check for various literal types.
	if tokens[*tokenIndex]._type == TOKEN_INT_LITERAL {
		val, _ := strconv.Atoi(tokens[*tokenIndex].value)
		leftNode = &ASTNode{
			Type:   AST_INT_LITERAL,
			intVal: int64(val),
		}
	} else if tokens[*tokenIndex]._type == TOKEN_VARCHAR_LITERAL {
		leftNode = &ASTNode{
			Type:   AST_VARCHAR_LITERAL,
			strVal: tokens[*tokenIndex].value,
		}
	} else if tokens[*tokenIndex]._type == TOKEN_FLOAT_LITERAL {
		f, err := strconv.ParseFloat(tokens[*tokenIndex].value, 64)
		if err != nil {
			panic(err)
		}
		leftNode = &ASTNode{
			Type:     AST_FLOAT_LITERAL,
			floatVal: f,
		}
	} else if tokens[*tokenIndex]._type == TOKEN_BOOLEAN_LITERAL {
		b, err := strconv.ParseBool(tokens[*tokenIndex].value)
		if err != nil {
			panic(err)
		}
		leftNode = &ASTNode{
			Type:      AST_BOOLEAN_LITERAL,
			boolValue: b,
		}
	} else {
		// Otherwise, treat it as a column name.
		leftNode = &ASTNode{
			Type:       AST_COLUMN_NAME,
			columnName: tokens[*tokenIndex].value,
		}
	}
	(*tokenIndex)++ // Move past literal or identifier.

	if endExpressionValues(tokens[*tokenIndex]) {
		return leftNode
	} else {
		operator := tokens[*tokenIndex].value
		(*tokenIndex)++ // Move past the operator.
		return &ASTNode{
			Type:     AST_BINARY,
			left:     leftNode,
			operator: operator,
			right:    parseExpression(tokens, tokenIndex),
		}
	}
	// Unreachable, but required for function completeness.
	return nil
}

func isTokenSELECTSpliter(tokens []*Token, tokenIndex *int) bool {
	t := tokens[*tokenIndex]._type
	return t == TOKEN_FROM || t == TOKEN_AS || t == TOKEN_RPAREN ||
		t == TOKEN_COMMA || t == TOKEN_SEMICOLON
}

func checkTokenIsFunction(token *Token) bool {
	if token._type == TOKEN_MAX || token._type == TOKEN_MIN || token._type == TOKEN_AVG || token._type == TOKEN_SUM || token._type == TOKEN_COUNT {
		return true
	}
	return false
}

func tokenToColumnType(token *Token) columnType {
	switch token._type {
	case TOKEN_MAX:
		return COLUMN_TYPE_MAX
	case TOKEN_MIN:

		return COLUMN_TYPE_MIN
	case TOKEN_AVG:

		return COLUMN_TYPE_AVG
	case TOKEN_SUM:
		return COLUMN_TYPE_SUM
	case TOKEN_COUNT:
		return COLUMN_TYPE_COUNT
	default:
		panic("Unknown token type")
	}
}

func parseSelectCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	fmt.Println("Starting parseSelectCommand")

	panicIfWrongType(tokens[*tokenIndex], TOKEN_SELECT)
	fmt.Println("Matched SELECT token")
	(*tokenIndex)++

	selectNode := ASTNode{Type: AST_SELECT}
	selectNode.columns = make([]*ASTNode, 0)
	selectNode.columnNames = make([]string, 0)
	selectNode.columnAliases = make([]string, 0)
	selectNode.columnTypes = make([]columnType, 0)

	if tokens[*tokenIndex]._type == TOKEN_STAR {
		fmt.Println("Found * token — SELECT * not yet implemented")
		// Go through each one (not implemented here)
	} else {
		fmt.Println("Parsing SELECT column list...")
		for tokens[*tokenIndex]._type != TOKEN_FROM {
			if checkTokenIsFunction(tokens[*tokenIndex]) {
				fmt.Printf("Found function: %s\n", tokens[*tokenIndex].value)
				selectNode.columnTypes = append(selectNode.columnTypes, tokenToColumnType(tokens[*tokenIndex]))

				(*tokenIndex)++

				panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
				fmt.Println("Matched LPAREN after function")
				(*tokenIndex)++

				selectNode.columnNames = append(selectNode.columnNames, tokens[*tokenIndex].value)
				fmt.Printf("Added function argument column: %s\n", tokens[*tokenIndex].value)
				(*tokenIndex)++

				panicIfWrongType(tokens[*tokenIndex], TOKEN_RPAREN)
				fmt.Println("Matched RPAREN after function argument")
				(*tokenIndex)++
			} else {
				fmt.Printf("Found regular column: %s\n", tokens[*tokenIndex].value)
				selectNode.columnNames = append(selectNode.columnNames, tokens[*tokenIndex].value)
				selectNode.columnTypes = append(selectNode.columnTypes, COLUMN_TYPE_NORMAL)
				(*tokenIndex)++
			}

			if tokens[*tokenIndex]._type == TOKEN_AS {
				(*tokenIndex)++
				selectNode.columnAliases = append(selectNode.columnAliases, tokens[*tokenIndex].value)
				fmt.Printf("Added alias: %s\n", tokens[*tokenIndex].value)
				(*tokenIndex)++
			} else {
				selectNode.columnAliases = append(selectNode.columnAliases, selectNode.columnNames[len(selectNode.columnNames)-1])
				fmt.Printf("No alias, using column name as alias: %s\n", selectNode.columnNames[len(selectNode.columnNames)-1])
			}

			if tokens[*tokenIndex]._type == TOKEN_COMMA {
				fmt.Println("Found comma, moving to next column")
				(*tokenIndex)++
			}
		}
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_FROM)
	fmt.Println("Matched FROM token")
	(*tokenIndex)++

	selectNode.tableName = tokens[*tokenIndex].value
	fmt.Printf("Set table name: %s\n", selectNode.tableName)
	(*tokenIndex)++

	if tokens[*tokenIndex]._type == TOKEN_GROUP {
		fmt.Println("Found GROUP token")
		(*tokenIndex)++

		panicIfWrongType(tokens[*tokenIndex], TOKEN_BY)
		fmt.Println("Matched BY token after GROUP")
		(*tokenIndex)++

		fmt.Println("Parsing GROUP BY columns...")
		for !isTokenSELECTSpliter(tokens, tokenIndex) {
			if tokens[*tokenIndex]._type == TOKEN_COMMA {
				fmt.Println("Skipping comma in GROUP BY")
				(*tokenIndex)++
			}

			col := tokens[*tokenIndex].value
			fmt.Printf("Checking GROUP BY column: %s\n", col)
			matched := false

			for i, name := range selectNode.columnNames {
				if name == col {
					selectNode.columnTypes[i] = COLUMN_TYPE_GROUP_BY
					fmt.Printf("Marked column %s as GROUP BY\n", col)
					matched = true
					break
				}
			}

			if !matched {
				fmt.Println("Current SELECT column names:")
				for _, name := range selectNode.columnNames {
					fmt.Printf(" - %s\n", name)
				}
				panic(fmt.Sprintf("GROUP BY column %q not in SELECT list", col))
			}

			(*tokenIndex)++
		}
	}

	fmt.Println("Finished parseSelectCommand successfully")
	return &selectNode
}

// expectMore := true
// for expectMore {
// 	selectNode.columns = append(selectNode.columns, parseExpression(tokens, tokenIndex))
// 	if checkType(tokens[*tokenIndex], TOKEN_AS) {
// 		(*tokenIndex)++ // Move past AS token
// 		selectNode.columnNames = append(selectNode.columnNames, tokens[*tokenIndex].value)
// 		(*tokenIndex)++
// 	} else {
// 		selectNode.columnNames = append(selectNode.columnNames, "")
// 	}

// 	if checkType(tokens[*tokenIndex], TOKEN_COMMA) {
// 		(*tokenIndex)++ // Move past comma and keep looping
// 	} else {
// 		expectMore = false
// 	}
// }

// panicIfWrongType(tokens[*tokenIndex], TOKEN_FROM)
// (*tokenIndex)++ // Move past FROM token
// selectNode.tableName = tokens[*tokenIndex].value
// (*tokenIndex)++ // Move past name token
// if checkType(tokens[*tokenIndex], TOKEN_WHERE) {

// }

// fmt.Println("HERERWEWRE")
// fmt.Println(tokens[*tokenIndex].value)

// if checkType(tokens[*tokenIndex], TOKEN_GROUP) {
// 	(*tokenIndex)++
// 	fmt.Println("HERERWEWRE")
// 	if checkType(tokens[*tokenIndex], TOKEN_BY) {
// 		fmt.Println("HERERWEWRE")
// 		selectNode.containsGroupBy = true
// 		(*tokenIndex)++
// 		for isTokenSELECTSpliter(tokens, tokenIndex) {
// 			selectNode.groupByColumns = append(selectNode.groupByColumns, tokens[*tokenIndex].value)
// 			(*tokenIndex)++
// 			if isTokenSELECTSpliter(tokens, tokenIndex) {
// 				break
// 			}
// 			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
// 			(*tokenIndex)++ // Move past comma
// 		}
// 	} else {
// 		panic("Expected by")
// 	}
// }

// if checkType(tokens[*tokenIndex], TOKEN_HAVING) {

// }

// if checkType(tokens[*tokenIndex], TOKEN_ORDER) {
// 	if checkType(tokens[*tokenIndex], TOKEN_BY) {

// 	} else {
// 		panic("Expected by")
// 	}
// }
// if checkType(tokens[*tokenIndex], TOKEN_LIMIT) {

// }
// if checkType(tokens[*tokenIndex], TOKEN_OFFSET) {

// }

func parseInsertCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INSERT)
	(*tokenIndex)++ // INSERT
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INTO)
	(*tokenIndex)++ // INTO

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
	newInsertNode := ASTNode{Type: AST_INSERT}
	newInsertNode.tableName = tokens[*tokenIndex].value
	(*tokenIndex)++ // Table name

	// Optional column list
	if checkType(tokens[*tokenIndex], TOKEN_LPAREN) {
		(*tokenIndex)++ // LPAREN

		for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			columnName := tokens[*tokenIndex].value
			newInsertNode.columnNames = append(newInsertNode.columnNames, columnName)
			(*tokenIndex)++

			if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
				panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
				(*tokenIndex)++
			}
		}
		(*tokenIndex)++ // RPAREN
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_VALUES)
	(*tokenIndex)++ // VALUES

	panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
	(*tokenIndex)++ // LPAREN

	// Values
	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		if checkType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE) {
			(*tokenIndex)++ // Opening quote
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			val := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, val)
			(*tokenIndex)++
			panicIfWrongType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE)
			(*tokenIndex)++ // Closing quote
		} else {
			val := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, val)
			(*tokenIndex)++
		}

		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++
		}
	}

	(*tokenIndex)++ // RPAREN
	if checkType(tokens[*tokenIndex], TOKEN_SEMICOLON) {
		(*tokenIndex)++
	}

	return &newInsertNode
}

func parseCreateCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_CREATE)
	(*tokenIndex)++ // Move past CREATE token
	panicIfWrongType(tokens[*tokenIndex], TOKEN_TABLE)
	(*tokenIndex)++ // Move past TABLE token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
	tableName := tokens[*tokenIndex].value
	(*tokenIndex)++ // Move past table name token

	newColumns := make([]*ASTNode, 0)

	if checkType(tokens[*tokenIndex], TOKEN_SEMICOLON) {
		(*tokenIndex)++ // Ingest semicolon (empty column list)
	} else if checkType(tokens[*tokenIndex], TOKEN_LPAREN) {
		(*tokenIndex)++ // Move past LPAREN token
		for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			newColumn := ASTNode{Type: AST_COLUMN}
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			newColumn.name = tokens[*tokenIndex].value
			(*tokenIndex)++ // Move past column name token

			// Todo: Panic if type is not INT/VARCHAR/...
			if checkType(tokens[*tokenIndex], TOKEN_VARCHAR) {
				newColumn._type = "VARCHAR"
				(*tokenIndex)++ // Move past VARCHAR token
				panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
				(*tokenIndex)++ // Move past LPAREN token
				panicIfWrongType(tokens[*tokenIndex], TOKEN_INT_LITERAL)
				newColumn.varCharLimit, _ = strconv.Atoi(tokens[*tokenIndex].value)
				(*tokenIndex)++ // Move past INT_LITERAL token
				panicIfWrongType(tokens[*tokenIndex], TOKEN_RPAREN)
				(*tokenIndex)++ // Move past RPAREN token
			} else {
				newColumn._type = tokens[*tokenIndex].value
				(*tokenIndex)++ // Move past type token
			}
			newColumn.constraints = make([]string, 0)

			for !checkType(tokens[*tokenIndex], TOKEN_COMMA) && !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
				if checkType(tokens[*tokenIndex], TOKEN_PRIMARY) {
					(*tokenIndex)++ // Move past PRIMARY
					panicIfWrongType(tokens[*tokenIndex], TOKEN_KEY)
					(*tokenIndex)++ // Move past KEY
					newColumn.constraints = append(newColumn.constraints, "PRIMARY KEY")
				} else {
					newColumn.constraints = append(newColumn.constraints, tokens[*tokenIndex].value)
					(*tokenIndex)++
				}
			}
			if checkType(tokens[*tokenIndex], TOKEN_COMMA) {
				(*tokenIndex)++ // Ingest comma
			}
			newColumns = append(newColumns, &newColumn)
		}
		(*tokenIndex)++ // Move past RPAREN token
	} else {
		panic("Invalid type")
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_SEMICOLON)
	(*tokenIndex)++ // Move past SEMICOLON token

	newCreateNode := ASTNode{
		Type:      AST_CREATE,
		tableName: tableName,
		columns:   newColumns,
	}
	return &newCreateNode
}

func parseCommands(tokens []*Token) []*ASTNode {
	retNodes := make([]*ASTNode, 0)
	tokenIndex := 0
	for tokenIndex < len(tokens) && tokens[tokenIndex]._type != TOKEN_EOF {
		if tokens[tokenIndex]._type == TOKEN_CREATE {
			retNodes = append(retNodes, parseCreateCommand(tokens, &tokenIndex))
		} else if tokens[tokenIndex]._type == TOKEN_INSERT {
			retNodes = append(retNodes, parseInsertCommand(tokens, &tokenIndex))
		} else if tokens[tokenIndex]._type == TOKEN_SELECT {
			retNodes = append(retNodes, parseSelectCommand(tokens, &tokenIndex))
		} else {
			// Skip unhandled tokens.
			tokenIndex++
		}
	}
	return retNodes
}

// --- Print AST function ---
func printAST(node *ASTNode, indent int) {
	indentStr := strings.Repeat("  ", indent)
	switch node.Type {
	case AST_CREATE:
		fmt.Printf("%sCREATE TABLE %s\n", indentStr, node.tableName)
		fmt.Printf("%sColumns:\n", indentStr)
		for _, col := range node.columns {
			fmt.Printf("%s- Column: %s, Type: %s", indentStr+"  ", col.name, col._type)
			if col._type == "VARCHAR" {
				fmt.Printf("(%d)", col.varCharLimit)
			}
			if len(col.constraints) > 0 {
				fmt.Printf(", Constraints: %v", col.constraints)
			}
			fmt.Println()
		}
	case AST_INSERT:
		fmt.Printf("%sINSERT INTO %s\n", indentStr, node.tableName)
		if len(node.columnNames) > 0 {
			fmt.Printf("%sColumns: %s\n", indentStr+"  ", strings.Join(node.columnNames, ", "))
		}
		if len(node.columnValues) > 0 {
			fmt.Printf("%sValues: %s\n", indentStr+"  ", strings.Join(node.columnValues, ", "))
		}
	case AST_SELECT:
		fmt.Printf("%sSELECT statement\n", indentStr)
		if len(node.columns) > 0 {
			fmt.Printf("%sColumns:\n", indentStr+"  ")
			for i, col := range node.columns {
				fmt.Printf("%s- %s", indentStr+"    ", col.columnName)
				if i < len(node.columnNames) && node.columnNames[i] != "" {
					fmt.Printf(" AS %s", node.columnNames[i])
				}
				fmt.Println()
			}
			if node.containsGroupBy {
				fmt.Printf("%sGROUP BY: %s\n", indentStr+"  ", strings.Join(node.groupByColumns, ", "))
			}
		}
		fmt.Printf("%sFROM: %s\n", indentStr+"  ", node.tableName)
	case AST_FUNCTION:
		fmt.Printf("%sFUNCTION: %s\n", indentStr, node.functionName)
		fmt.Printf("%sArguments:\n", indentStr+"  ")
		for _, arg := range node.functionArguements {
			printAST(arg, indent+2)
		}
	case AST_BINARY:
		fmt.Printf("%sBINARY EXPRESSION:\n", indentStr)
		fmt.Printf("%sLeft:\n", indentStr+"  ")
		printAST(node.left, indent+2)
		fmt.Printf("%sOperator: %s\n", indentStr+"  ", node.operator)
		fmt.Printf("%sRight:\n", indentStr+"  ")
		printAST(node.right, indent+2)
	case AST_INT_LITERAL:
		fmt.Printf("%sINT_LITERAL: %d\n", indentStr, node.intVal)
	case AST_VARCHAR_LITERAL:
		fmt.Printf("%sVARCHAR_LITERAL: %s\n", indentStr, node.strVal)
	case AST_FLOAT_LITERAL:
		fmt.Printf("%sFLOAT_LITERAL: %f\n", indentStr, node.floatVal)
	case AST_BOOLEAN_LITERAL:
		fmt.Printf("%sBOOLEAN_LITERAL: %t\n", indentStr, node.boolValue)
	case AST_COLUMN_NAME:
		fmt.Printf("%sCOLUMN_NAME: %s\n", indentStr, node.columnName)
	default:
		fmt.Printf("%sUnknown AST Node type: %d\n", indentStr, node.Type)
	}
}
