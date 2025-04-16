package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ---------------------------
// AST Definitions
// ---------------------------
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// SelectStatement represents a simple SELECT statement.
type SelectStatement struct {
	Token   Token        // The SELECT token
	Columns []Expression // List of columns, could be identifiers, stars, or function calls
	Table   Expression   // Table name (an identifier)
	Where   Expression   // Optional WHERE clause; can be nil
}

func (ss *SelectStatement) statementNode()       {}
func (ss *SelectStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SelectStatement) String() string {
	var cols []string
	for _, col := range ss.Columns {
		cols = append(cols, col.String())
	}
	out := fmt.Sprintf("%s %s FROM %s", ss.TokenLiteral(), strings.Join(cols, ", "), ss.Table.String())
	if ss.Where != nil {
		out += " WHERE " + ss.Where.String()
	}
	return out
}

// Identifier represents table/column names or even function names.
type Identifier struct {
	Token Token // The IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token Token // The INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token Token // The STRING token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "'" + sl.Value + "'" }

// CallExpression represents a function call (e.g. COUNT(*)).
type CallExpression struct {
	Token     Token        // The function name token
	Function  Expression   // Typically an Identifier
	Arguments []Expression // Arguments passed to the function
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", ce.Function.String(), strings.Join(args, ", "))
}

// InfixExpression represents binary expressions, mainly used for WHERE clauses.
type InfixExpression struct {
	Token    Token      // The operator token, e.g. =
	Left     Expression // Left-hand side
	Operator string
	Right    Expression // Right-hand side
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// ---------------------------
// Parser
// ---------------------------
type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	// Initialize by reading two tokens.
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseStatement parses a SQL statement.
func (p *Parser) ParseStatement() Statement {
	if p.curToken.Type == SELECT {
		return p.parseSelectStatement()
	}
	msg := fmt.Sprintf("unexpected token %s", p.curToken.Literal)
	p.errors = append(p.errors, msg)
	return nil
}

// parseSelectStatement parses a SELECT statement.
func (p *Parser) parseSelectStatement() *SelectStatement {
	stmt := &SelectStatement{Token: p.curToken}

	// Parse column list.
	stmt.Columns = p.parseColumnList()

	// Expect FROM.
	if !p.expectPeek(FROM) {
		return nil
	}
	p.nextToken() // skip FROM token

	// Parse table identifier.
	stmt.Table = p.parseIdentifier()

	// Optionally parse WHERE clause.
	if p.peekToken.Type == WHERE {
		p.nextToken() // move to WHERE
		p.nextToken() // start of condition
		stmt.Where = p.parseExpression(LOWEST)
	}
	return stmt
}

// parseColumnList parses columns in the SELECT clause.
// Supports '*' or a comma-separated list of expressions.
func (p *Parser) parseColumnList() []Expression {
	var columns []Expression

	p.nextToken() // move to the first column/expression
	if p.curToken.Type == ASTERISK {
		columns = append(columns, &Identifier{Token: p.curToken, Value: p.curToken.Literal})
		p.nextToken()
		return columns
	}

	columns = append(columns, p.parseExpression(LOWEST))
	for p.curToken.Type == COMMA {
		p.nextToken() // skip comma
		p.nextToken() // next expression
		columns = append(columns, p.parseExpression(LOWEST))
	}
	return columns
}

// parseIdentifier parses an identifier.
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

const (
	_ int = iota
	LOWEST
	EQUALS // =
)

var precedences = map[TokenType]int{
	EQ: EQUALS,
	// Additional precedences can be added here.
}

// parseExpression parses an expression given a minimum precedence.
func (p *Parser) parseExpression(precedence int) Expression {
	var leftExp Expression

	// Determine the primary expression.
	switch p.curToken.Type {
	case IDENT:
		// Look ahead to see if it's a function call.
		if p.peekToken.Type == LPAREN {
			leftExp = p.parseCallExpression()
		} else {
			leftExp = p.parseIdentifier()
		}
	case INT:
		leftExp = p.parseIntegerLiteral()
	case STRING:
		leftExp = p.parseStringLiteral()
	case ASTERISK:
		// Allow '*' as an expression (e.g. in COUNT(*))
		leftExp = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case LPAREN:
		p.nextToken() // skip '('
		leftExp = p.parseExpression(LOWEST)
		if !p.expectPeek(RPAREN) {
			return nil
		}
		p.nextToken() // skip ')'
	default:
		msg := fmt.Sprintf("unexpected token in expression: %s", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Parse binary (infix) expressions for WHERE clause.
	for p.peekToken.Type != SEMICOLON && precedence < p.peekPrecedence() {
		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse integer %q", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	prec := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(prec)
	return expression
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
	return false
}

// parseCallExpression handles function calls (e.g. COUNT(*)).
func (p *Parser) parseCallExpression() Expression {
	// The current token is the function name (an identifier).
	function := p.parseIdentifier()

	// Expect a '(' after the function name.
	if !p.expectPeek(LPAREN) {
		return nil
	}

	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

// parseExpressionList parses a comma-separated list of expressions, ending with the given token.
func (p *Parser) parseExpressionList(end TokenType) []Expression {
	var list []Expression

	// Special case: empty argument list.
	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekToken.Type == COMMA {
		p.nextToken() // skip comma
		p.nextToken() // move to next expression
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// ---------------------------
// Main Function: Demonstration
// ---------------------------
func main2() {
	// Examples of SQL statements including function calls.
	inputs := []string{
		"SELECT * FROM users;",
		"SELECT id, COUNT(*) FROM employees WHERE age >= 30;",
		"SELECT SUM(salary), department FROM payroll WHERE department = 'Engineering';",
	}

	for _, input := range inputs {
		fmt.Println("Input:", input)
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		stmt := parser.ParseStatement()
		if stmt != nil {
			fmt.Println("AST:", stmt.String())
		} else {
			fmt.Println("Parsing errors:")
			for _, err := range parser.errors {
				fmt.Println("  -", err)
			}
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}
