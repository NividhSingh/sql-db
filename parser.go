package main

type DataType int

const (
	AST_INT_LITERAL DataType = iota
	AST_VARCHAR_LITERAL
	AST_SELECT
	AST_CREATE
	AST_INSERT
)

type ASTNode struct {
	Type   DataType
	IntVal int
	StrVal string
}

func parseCreateCommand(tokens []*Token) *ASTNode {}

func parseCommands(tokens []*Token) *ASTNode {
	tokenIndex := 0
	if tokens[tokenIndex]._type == TOKEN_EOF {
		return nil
	} else if tokens[tokenIndex]._type == TOKEN_CREATE {
		parseCreateCommand(tokens)
	}
}
