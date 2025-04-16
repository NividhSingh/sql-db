package main

import (
	"fmt"
	"strconv"
)

type ASTNodeType int

const (
	AST_INT_LITERAL ASTNodeType = iota
	AST_VARCHAR_LITERAL
	AST_SELECT
	AST_CREATE
	AST_INSERT
)

type ASTNode struct {
	Type   ASTNodeType
	IntVal int
	StrVal string

	_selectNode SELECTNode

	_createNode CREATENode

	_insertNode INSERTNode
}

type columnNode struct {
	name         string
	_type        string
	varCharLimit int
	constraints  []string
}

type CREATENode struct {
	tableName string
	columns   []*columnNode
}

type INSERTNode struct {
	tableName    string
	columns      []string
	columnValues []string
}

type SELECTNode struct {
}

func panicIfWrongType(token *Token, _token_type TokenType) {
	fmt.Println(token.value)
	fmt.Println(token._type)

	if !checkType(token, _token_type) {
		panic("Wrong token type")
		// panicf("expected %s, got %s", _token_type.String(), token._type.String())
	}
}

func checkType(token *Token, _token_type TokenType) bool {
	return token._type == _token_type
}

func parseInsertCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INSERT)
	(*tokenIndex)++ // Move past insert token
	panicIfWrongType(tokens[*tokenIndex], TOKEN_INTO)
	(*tokenIndex)++ // Move past into token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)

	newInsertNode := INSERTNode{}
	newInsertNode.tableName = tokens[*tokenIndex].value
	(*tokenIndex)++ // Move past table name token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
	(*tokenIndex)++ // Move past left paren token

	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
		columnName := tokens[*tokenIndex].value
		newInsertNode.columns = append(newInsertNode.columns, columnName)
		(*tokenIndex)++ // Move past column name token
		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++ // Move past comma token
		}
	}
	panicIfWrongType(tokens[*tokenIndex], TOKEN_VALUES)

	for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
		if checkType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE) {
			(*tokenIndex)++
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			columnName := tokens[*tokenIndex].value
			newInsertNode.columns = append(newInsertNode.columns, columnName)
			(*tokenIndex)++ // Move past column name token
			panicIfWrongType(tokens[*tokenIndex], TOKEN_SINGLE_QUOTE)
			(*tokenIndex)++ // Move past single quote token
		} else {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			columnValue := tokens[*tokenIndex].value
			newInsertNode.columnValues = append(newInsertNode.columnValues, columnValue)
			(*tokenIndex)++ // Move past column name token
		}
		if !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			panicIfWrongType(tokens[*tokenIndex], TOKEN_COMMA)
			(*tokenIndex)++ // Move past comma token
		}
	}
	return &ASTNode{Type: AST_INSERT, _insertNode: newInsertNode}
}

func parseCreateCommand(tokens []*Token, tokenIndex *int) *ASTNode {
	panicIfWrongType(tokens[*tokenIndex], TOKEN_CREATE)
	(*tokenIndex)++ // Move past create token
	panicIfWrongType(tokens[*tokenIndex], TOKEN_TABLE)
	(*tokenIndex)++ // Move past table token

	panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
	tableName := tokens[*tokenIndex].value
	(*tokenIndex)++ // Move past table name token

	newColumnms := make([]*columnNode, 0)

	if checkType(tokens[*tokenIndex], TOKEN_SEMICOLON) {
		(*tokenIndex)++
	} else if checkType(tokens[*tokenIndex], TOKEN_LPAREN) {
		(*tokenIndex)++
		for !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
			newColumn := columnNode{}
			panicIfWrongType(tokens[*tokenIndex], TOKEN_IDENTIFIER)
			newColumn.name = tokens[*tokenIndex].value
			(*tokenIndex)++
			// Todo: Panic if not itn/varchar/date etc
			if checkType(tokens[*tokenIndex], TOKEN_VARCHAR) {
				newColumn._type = "VARCHAR"
				(*tokenIndex)++
				panicIfWrongType(tokens[*tokenIndex], TOKEN_LPAREN)
				(*tokenIndex)++
				panicIfWrongType(tokens[*tokenIndex], TOKEN_INT_LITERAL)
				newColumn.varCharLimit, _ = strconv.Atoi(tokens[*tokenIndex].value)
				(*tokenIndex)++
				panicIfWrongType(tokens[*tokenIndex], TOKEN_RPAREN)
				(*tokenIndex)++
			} else {
				newColumn._type = tokens[*tokenIndex].value
				(*tokenIndex)++
			}
			newColumn.constraints = make([]string, 0)

			for !checkType(tokens[*tokenIndex], TOKEN_COMMA) && !checkType(tokens[*tokenIndex], TOKEN_RPAREN) {
				if checkType(tokens[*tokenIndex], TOKEN_PRIMARY) {
					(*tokenIndex)++
					panicIfWrongType(tokens[*tokenIndex], TOKEN_KEY)
					(*tokenIndex)++
					newColumn.constraints = append(newColumn.constraints, "PRIMARY KEY")
					// Todo: Add other things like not null here
				} else {
					newColumn.constraints = append(newColumn.constraints, tokens[*tokenIndex].value)
					(*tokenIndex)++

				}
			}
			if checkType(tokens[*tokenIndex], TOKEN_COMMA) {
				(*tokenIndex)++ // Injest comma

			}

			newColumnms = append(newColumnms, &newColumn)
		}
		(*tokenIndex)++ // Injest right parenthesis

	} else {
		panic("Invalid type")
	}

	panicIfWrongType(tokens[*tokenIndex], TOKEN_SEMICOLON)

	(*tokenIndex)++ // Injest semicolon

	newCreateNode := CREATENode{
		tableName: tableName,
		columns:   newColumnms,
	}
	return &ASTNode{
		Type:        AST_CREATE,
		_createNode: newCreateNode,
	}
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
			// retNodes = append(retNodes, parseSelectCommand(tokens, &tokenIndex))
		}

	}
	return retNodes
}

func printCreateAST(node *ASTNode) {
	if node.Type != AST_CREATE {
		fmt.Println("Not a CREATE statement.")
		return
	}
	fmt.Printf("CREATE TABLE %s\n", node._createNode.tableName)
	fmt.Println("Columns:")
	for _, col := range node._createNode.columns {
		fmt.Printf("  Name: %s, Type: %s", col.name, col._type)
		if col._type == "VARCHAR" {
			fmt.Printf("(%d)", col.varCharLimit)
		}
		if len(col.constraints) > 0 {
			fmt.Printf(", Constraints: %v", col.constraints)
		}
		fmt.Println()
	}
}
func main() {
	// The input string is designed to produce the following tokens:
	// CREATE, TABLE, IDENTIFIER("myTable"), RPAREN, IDENTIFIER("col1"), VARCHAR,
	// LPAREN, INT_LITERAL("255"), RPAREN, PRIMARY, KEY, COMMA, RPAREN, SEMICOLON, EOF.
	input := "CREATE TABLE myTable (col1 VARCHAR (255) PRIMARY KEY, col2 INT);"

	// Initialize the lexer.
	lexer := &Lexer{
		input:   input,
		start:   0,
		current: 0,
		line:    1,
	}

	// Tokenize the input by repeatedly calling getNextToken.
	var tokens []*Token
	for {
		token := getNextToken(lexer)
		tokens = append(tokens, &token)
		if token._type == TOKEN_EOF {
			break
		}
	}

	// Optionally, you can print all the tokens for debugging.
	fmt.Println("Tokens:")
	for _, t := range tokens {
		fmt.Printf("Type: %-15s Value: %q\n", tokenTypeToString(t._type), t.value)
	}

	// Now, parse the tokens to create an AST.
	// tokenIndex := 0
	ast := parseCommands(tokens)

	// Print the resulting AST.
	printCreateAST(ast[0])
}
