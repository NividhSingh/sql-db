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
	columnNames []string

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

func parseSelectCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_SELECT)
	(*tokenIndex)++ // Move past SELECT token

	selectNode := ASTNode{Type: AST_SELECT}
	selectNode.columns = make([]*ASTNode, 0)
	selectNode.columnNames = make([]string, 0)

	expectMore := true
	for expectMore {
		selectNode.columns = append(selectNode.columns, parseExpression(tokens, tokenIndex))
		if checkType(tokens[*tokenIndex], TOKEN_AS) {
			(*tokenIndex)++ // Move past AS token
			selectNode.columnNames = append(selectNode.columnNames, tokens[*tokenIndex].value)
			(*tokenIndex)++
		} else {
			selectNode.columnNames = append(selectNode.columnNames, "")
		}

		if checkType(tokens[*tokenIndex], TOKEN_COMMA) {
			(*tokenIndex)++ // Move past comma and keep looping
		} else {
			expectMore = false
		}
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_FROM)
	(*tokenIndex)++ // Move past FROM token
	selectNode.tableName = tokens[*tokenIndex].value

	return &selectNode
}

func parseInsertCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INSERT)
	(*tokenIndex)++ // Move past INSERT token
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INTO)
	(*tokenIndex)++ // Move past INTO token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
	newInsertNode := ASTNode{Type: AST_INSERT}
	newInsertNode.tableName = tokens[*tokenIndex].value
	(*tokenIndex)++ // Move past table name token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
	(*tokenIndex)++ // Move past LPAREN token

	// Parse column names.
	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
		columnName := tokens[*tokenIndex].value
		newInsertNode.columnNames = append(newInsertNode.columnNames, columnName)
		(*tokenIndex)++ // Move past column name token
		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++ // Move past COMMA token
		}
	}
	(*tokenIndex)++ // Move past RPAREN token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_VALUES)
	(*tokenIndex)++ // Move past VALUES token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
	(*tokenIndex)++ // Move past LPAREN token

	// Parse column values.
	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		if checkType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE) {
			(*tokenIndex)++ // Skip opening quote
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			columnValue := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, columnValue)
			(*tokenIndex)++ // Move past identifier
			panicIfWrongType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE)
			(*tokenIndex)++ // Move past closing quote
		} else {
			// If not in quotes, assume an identifier.
			columnValue := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, columnValue)
			(*tokenIndex)++ // Move past token
		}
		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++ // Move past COMMA token
		}
	}
	(*tokenIndex)++ // Move past RPAREN token
	(*tokenIndex)++ // Move past SEMICOLON token

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
