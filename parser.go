package main

import (
	"fmt"
	"strconv"
)

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
	Type      ASTNodeType
	IntVal    int
	StrVal    string
	floatVal  float64
	boolValue bool

	// Expression Node
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

	// Column node
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

// --- Helper functions ---

func panicIfWrongType(token *Token, _token_type TokenType) {
	fmt.Println(token.value)
	fmt.Println(token._type)
	fmt.Println(_token_type)
	if !checkType(token, _token_type) {
		panic("Wrong token type")
		// panicf("expected %s, got %s", _token_type.String(), token._type.String())
	}
}

func endExpressionValues(token *Token) bool {
	return token._type == TOKEN_COMMA || token._type == TOKEN_FROM || token._type == TOKEN_AS || token._type == TOKEN_FROM || token._type == TOKEN_RPAREN
}

func functionValues(token *Token) bool {
	return token._type == TOKEN_COUNT || token._type == TOKEN_SUM || token._type == TOKEN_AVG || token._type == TOKEN_MIN || token._type == TOKEN_MAX
}

func checkType(token *Token, _token_type TokenType) bool {
	return token._type == _token_type
}

// --- Parsing functions ---

func parseExpression(tokens []*Token, tokenIndex *int) *ASTNode {
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
			(*tokenIndex)++
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
		leftNode = &ASTNode{
			Type:       AST_COLUMN_NAME,
			columnName: tokens[*tokenIndex].value,
		}
	}
	(*tokenIndex)++ // Move past the column name or literal.

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
		// TODO: Add parenthesis
	}

	// Return nil if no expression was parsed.
	return nil
}

func isTokenSELECTSpliter(tokens []*Token, tokenIndex *int) bool {
	return checkType(tokens[*tokenIndex], TOKEN_FROM) || checkType(tokens[*tokenIndex], TOKEN_WHERE) || checkType(tokens[*tokenIndex], TOKEN_GROUP) || checkType(tokens[*tokenIndex], TOKEN_HAVING) || checkType(tokens[*tokenIndex], TOKEN_ORDER) || checkType(tokens[*tokenIndex], TOKEN_LIMIT) || checkType(tokens[*tokenIndex], TOKEN_OFFSET) || checkType(tokens[*tokenIndex], TOKEN_LIMIT) || checkType(tokens[*tokenIndex], TOKEN_ORDER)
}

func parseSelectCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_SELECT)
	(*tokenIndex)++ // Move past select token

	selectNode := ASTNode{Type: AST_SELECT}

	selectNode.columns = make([]*ASTNode, 0)
	selectNode.columnNames = make([]string, 0)

	for !isTokenSELECTSpliter(tokens, tokenIndex) {
		selectNode.columns = append(selectNode.columns, parseExpression(tokens, tokenIndex))
		if checkType(tokens[*tokenIndex], TOKEN_AS) {
			(*tokenIndex)++ // Move past as token
			selectNode.columnNames = append(selectNode.columnNames, tokens[*tokenIndex].value)
			(*tokenIndex)++
		}
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_FROM)
	(*tokenIndex)++ // Move past from token
	selectNode.tableName = tokens[*tokenIndex].value

	return &selectNode

	// Go until from
	// GO until where
	// Go until group by
	// Go until having
	// Go until order by
	// Go until limit
	// Go unitl offset

	// Not yet implemented – return nil.
	return nil
}

func parseInsertCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INSERT)
	(*tokenIndex)++ // Move past insert token
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INTO)
	(*tokenIndex)++ // Move past into token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)

	newInsertNode := ASTNode{Type: AST_INSERT}
	newInsertNode.tableName = tokens[*tokenIndex].value
	(*tokenIndex)++ // Move past table name token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
	(*tokenIndex)++ // Move past left paren token

	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
		columnName := tokens[*tokenIndex].value
		newInsertNode.columnNames = append(newInsertNode.columnNames, columnName)
		(*tokenIndex)++ // Move past column name token
		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++ // Move past comma token
		}
	}
	(*tokenIndex)++ // Move past RPAREN token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_VALUES)
	(*tokenIndex)++ // Move past values token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
	(*tokenIndex)++ // Move past left paren token

	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		if checkType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE) {
			(*tokenIndex)++
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			columnValue := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, columnValue)
			(*tokenIndex)++ // Move past column name token
			panicIfWrongType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE)
			(*tokenIndex)++ // Move past single quote token
		} else {
			// If not in quotes, assume an identifier.
			columnValue := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, columnValue)
			(*tokenIndex)++ // Move past column name token
		}
		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++ // Move past comma token
		}
	}
	(*tokenIndex)++ // Move past RPAREN token
	(*tokenIndex)++ // Move past semicolon token

	return &newInsertNode
}

func parseCreateCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_CREATE)
	(*tokenIndex)++ // Move past create token
	panicIfWrongType(tokens[*tokenIndex], TOKEN_TABLE)
	(*tokenIndex)++ // Move past table token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
	tableName := tokens[*tokenIndex].value
	(*tokenIndex)++ // Move past table name token

	newColumnms := make([]*ASTNode, 0)

	if checkType(tokens[*tokenIndex], TOKEN_SEMICOLON) {
		(*tokenIndex)++
	} else if checkType(tokens[*tokenIndex], TOKEN_LPAREN) {
		(*tokenIndex)++
		for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			newColumn := ASTNode{Type: AST_COLUMN}
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			newColumn.name = tokens[*tokenIndex].value
			(*tokenIndex)++ // Move past column name token
			// Todo: Panic if not int/varchar/date etc.
			if checkType(tokens[*tokenIndex], TOKEN_VARCHAR) {
				newColumn._type = "VARCHAR"
				(*tokenIndex)++ // Move past VARCHAR token
				panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
				(*tokenIndex)++ // Move past left paren token
				panicIfWrongType(tokens[*tokenIndex], TOKEN_INT_LITERAL)
				newColumn.varCharLimit, _ = strconv.Atoi(tokens[*tokenIndex].value)
				(*tokenIndex)++ // Move past INT_LITERAL token
				panicIfWrongType(tokens[*tokenIndex], TOKEN_RPAREN)
				(*tokenIndex)++ // Move past right paren token
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

			newColumnms = append(newColumnms, &newColumn)
		}
		(*tokenIndex)++ // Ingest right paren token

	} else {
		panic("Invalid type")
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_SEMICOLON)
	(*tokenIndex)++ // Ingest semicolon

	newCreateNode := ASTNode{
		Type:      AST_CREATE,
		tableName: tableName,
		columns:   newColumnms,
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
			// tokenIndex++ // Skip token for now.
		} else {
			// Skip unhandled token types.
			tokenIndex++
		}
	}
	return retNodes
}

func main() {
	// Your refactored code now compiles.
	// (Tokenization and a lexer are assumed to have been performed to produce the tokens slice.)
	// For demonstration purposes, you would call parseCommands with a valid slice of *Token.
}
