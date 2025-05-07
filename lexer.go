package main

import (
	"fmt"
)

// TokenType represents the type of token.
type TokenType int

const (
	// General Tokens
	TOKEN_EOF TokenType = iota
	TOKEN_ILLEGAL
	TOKEN_INT
	TOKEN_DOT
	TOKEN_VARCHAR
	TOKEN_FLOAT
	TOKEN_STAR
	TOKEN_NULL_LITERAL
	TOKEN_COMMA
	TOKEN_SEMICOLON
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_SINGLE_QUOTE
	TOKEN_DOUBLE_QUOTE

	TOKEN_INT_LITERAL
	TOKEN_VARCHAR_LITERAL
	TOKEN_FLOAT_LITERAL
	TOKEN_BOOLEAN_LITERAL

	TOKEN_IDENTIFIER
	TOKEN_FUNCTION

	// Data Types (additional)
	TOKEN_BOOLEAN
	TOKEN_DATE
	TOKEN_TIME
	TOKEN_TIMESTAMP

	// Operators
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_SLASH
	TOKEN_PERCENT
	TOKEN_EQUALS
	TOKEN_NOT_EQUALS
	TOKEN_LESS_THAN
	TOKEN_GREATER_THAN
	TOKEN_LESS_EQUAL
	TOKEN_GREATER_EQUAL

	// Logical Operators
	TOKEN_AND
	TOKEN_OR
	TOKEN_NOT

	// Create Table & DDL
	TOKEN_CREATE
	TOKEN_TABLE
	TOKEN_COLUMN
	TOKEN_ALTER
	TOKEN_DROP
	TOKEN_TO
	TOKEN_RENAME
	TOKEN_IN

	// Constraints (additional)
	TOKEN_PRIMARY
	TOKEN_KEY // Alternatively, you might split into TOKEN_PRIMARY and TOKEN_KEY
	TOKEN_UNIQUE
	TOKEN_FOREIGN // Consider splitting as needed: TOKEN_FOREIGN, TOKEN_KEY, TOKEN_REFERENCES
	TOKEN_REFERENCES
	TOKEN_CHECK
	TOKEN_DEFAULT
	TOKEN_AUTO_INCREMENT

	// Insert
	TOKEN_INSERT
	TOKEN_INTO
	TOKEN_VALUES

	// Update
	TOKEN_UPDATE
	TOKEN_SET

	// Delete
	TOKEN_DELETE

	// Select Query
	TOKEN_SELECT
	TOKEN_FROM
	TOKEN_WHERE
	TOKEN_GROUP // For GROUP BY (use TOKEN_GROUP and TOKEN_BY together)
	TOKEN_BY
	TOKEN_HAVING
	TOKEN_ORDER // For ORDER BY (use TOKEN_ORDER and TOKEN_BY)
	TOKEN_LIMIT
	TOKEN_OFFSET
	TOKEN_DISTINCT

	// Expressions and Aliases
	TOKEN_AS
	TOKEN_CASE
	TOKEN_WHEN
	TOKEN_THEN
	TOKEN_ELSE
	TOKEN_END

	// Joins
	TOKEN_JOIN
	TOKEN_INNER
	TOKEN_LEFT
	TOKEN_RIGHT
	TOKEN_FULL
	TOKEN_ON

	// Additional Conditional Tokens
	TOKEN_IF
	TOKEN_NOT_EXISTS

	// Transaction Control (if needed)
	TOKEN_COMMIT
	TOKEN_ROLLBACK
	TOKEN_SAVEPOINT

	// Data Control (if needed)
	TOKEN_GRANT
	TOKEN_REVOKE

	// The following tokens are referenced in LookupKeyword. They are not defined above.
	// You can define them if you plan to support these functions.
	// TOKEN_TRUNCATE, TOKEN_COUNT, TOKEN_SUM, TOKEN_MAX, TOKEN_MIN, TOKEN_AVG,
	// TOKEN_UNION, TOKEN_EXCEPT, TOKEN_INTERSECT, TOKEN_IS, TOKEN_NULL, TOKEN_TRUE, TOKEN_FALSE
	TOKEN_TRUNCATE
	TOKEN_COUNT
	TOKEN_SUM
	TOKEN_MAX
	TOKEN_MIN
	TOKEN_AVG
	TOKEN_UNION
	TOKEN_EXCEPT
	TOKEN_INTERSECT
	TOKEN_IS
	TOKEN_NULL
	TOKEN_TRUE
	TOKEN_FALSE
)

// Token represents a lexical token.
type Token struct {
	_type TokenType
	value string
}

// Lexer holds information about the lexing process.
type Lexer struct {
	start   int
	current int
	line    int
	input   string
}

// atEnd checks if we have reached the end of the input.
func atEnd(lexer *Lexer) bool {
	return lexer.current >= len(lexer.input)
}

// advance consumes the current character and returns it.
func advance(lexer *Lexer) byte {
	lexer.current++
	return lexer.input[lexer.current-1]
}

// peek returns the current character without consuming it.
func peek(lexer *Lexer) byte {
	if lexer.current < len(lexer.input) {
		return lexer.input[lexer.current]
	} else {
		return 0
	}

}

// peekNext returns the character after the current one.
func peekNext(lexer *Lexer) byte {
	if lexer.current+1 >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.current+1]
}

// skipWhitespace advances the lexer past any whitespace characters.
func skipWhitespace(lexer *Lexer) {
	for {
		c := peek(lexer)
		switch c {
		case ' ', '\t', '\n', '\r':
			advance(lexer)
		default:
			return
		}
		if atEnd(lexer) {
			return
		}
	}
}

// makeToken creates a new token from the current lexeme.
func makeToken(lexer *Lexer, _type TokenType, val string) Token {
	return Token{
		_type: _type,
		value: val,
	}
}

// isAlphabetic checks if the byte is an alphabetic character or underscore.
func isAlphabetic(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isDigit checks if the byte is a digit.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// scanNumber scans a numeric literal (integer or float).
func scanNumber(lexer *Lexer) Token {
	for !atEnd(lexer) && (isDigit(peek(lexer)) || peek(lexer) == '.') {
		advance(lexer)
	}
	return makeToken(lexer, TOKEN_INT_LITERAL, lexer.input[lexer.start:lexer.current])
}

// LookupKeyword returns the token type corresponding to a keyword.
func LookupKeyword(value string) TokenType {
	switch value {
	case "SELECT":
		return TOKEN_SELECT
	case "INSERT":
		return TOKEN_INSERT
	case "UPDATE":
		return TOKEN_UPDATE
	case "DELETE":
		return TOKEN_DELETE
	case "CREATE":
		return TOKEN_CREATE
	case "ALTER":
		return TOKEN_ALTER
	case "DROP":
		return TOKEN_DROP
	case "TRUNCATE":
		return TOKEN_TRUNCATE

	case "FROM":
		return TOKEN_FROM
	case "WHERE":
		return TOKEN_WHERE

	case "COUNT":
		return TOKEN_COUNT
	case "SUM":
		return TOKEN_SUM
	case "MAX":
		return TOKEN_MAX
	case "MIN":
		return TOKEN_MIN
	case "AVG":
		return TOKEN_AVG

	case "GROUP":
		return TOKEN_GROUP
	case "BY":
		return TOKEN_BY
	case "HAVING":
		return TOKEN_HAVING

	case "ORDER":
		return TOKEN_ORDER
	case "LIMIT":
		return TOKEN_LIMIT
	case "OFFSET":
		return TOKEN_OFFSET

	case "JOIN":
		return TOKEN_JOIN
	case "INNER":
		return TOKEN_INNER
	case "LEFT":
		return TOKEN_LEFT
	case "RIGHT":
		return TOKEN_RIGHT
	case "FULL":
		return TOKEN_FULL
	case "ON":
		return TOKEN_ON

	case "DISTINCT":
		return TOKEN_DISTINCT
	case "AS":
		return TOKEN_AS

	case "VALUES":
		return TOKEN_VALUES
	case "SET":
		return TOKEN_SET

	case "CASE":
		return TOKEN_CASE
	case "WHEN":
		return TOKEN_WHEN
	case "THEN":
		return TOKEN_THEN
	case "ELSE":
		return TOKEN_ELSE
	case "END":
		return TOKEN_END

	case "UNION":
		return TOKEN_UNION
	case "EXCEPT":
		return TOKEN_EXCEPT
	case "INTERSECT":
		return TOKEN_INTERSECT

	case "AND":
		return TOKEN_AND
	case "OR":
		return TOKEN_OR
	case "NOT":
		return TOKEN_NOT

	case "IN":
		return TOKEN_IN
	case "IS":
		return TOKEN_IS
	case "NULL":
		return TOKEN_NULL

	case "TRUE":
		return TOKEN_TRUE
	case "FALSE":
		return TOKEN_FALSE
	case "PRIMARY":
		return TOKEN_PRIMARY
	case "KEY":
		return TOKEN_KEY
	case "VARCHAR":
		return TOKEN_VARCHAR
	case "INT":
		return TOKEN_INT
	case "FLOAT":
		return TOKEN_FLOAT

	case "TABLE":
		return TOKEN_TABLE
	case "INTO":
		return TOKEN_INTO

	default:
		return TOKEN_IDENTIFIER
	}
}

// scanIdentifier scans an identifier or keyword.
func scanIdentifier(lexer *Lexer) Token {
	for !atEnd(lexer) && (isAlphabetic(peek(lexer)) || isDigit(peek(lexer))) {
		advance(lexer)
	}
	lexeme := lexer.input[lexer.start:lexer.current]
	return makeToken(lexer, LookupKeyword(lexeme), lexeme)
}

// getNextToken returns the next token from the input.
func getNextToken(lexer *Lexer) Token {
	skipWhitespace(lexer)
	if atEnd(lexer) {
		return makeToken(lexer, TOKEN_EOF, "")
	}

	lexer.start = lexer.current
	c := advance(lexer)

	// Check for identifiers
	if isAlphabetic(c) || c == '_' {
		return scanIdentifier(lexer)
	}
	// Check for numbers
	if isDigit(c) {
		return scanNumber(lexer)
	}

	switch c {
	case ';':
		return makeToken(lexer, TOKEN_SEMICOLON, ";")
	case '(':
		return makeToken(lexer, TOKEN_LPAREN, "(")
	case ')':
		return makeToken(lexer, TOKEN_RPAREN, ")")
	case ',':
		return makeToken(lexer, TOKEN_COMMA, ",")
	case '+':
		return makeToken(lexer, TOKEN_PLUS, "+")
	case '-':
		return makeToken(lexer, TOKEN_MINUS, "-")
	case '*':
		return makeToken(lexer, TOKEN_STAR, "*")
	case '/':
		return makeToken(lexer, TOKEN_SLASH, "/")
	case '%':
		return makeToken(lexer, TOKEN_PERCENT, "%")
	case '=':
		return makeToken(lexer, TOKEN_EQUALS, "=")
	case '<':
		if peek(lexer) == '=' {
			advance(lexer)
			return makeToken(lexer, TOKEN_LESS_EQUAL, "<=")
		} else if peek(lexer) == '>' {
			advance(lexer)
			return makeToken(lexer, TOKEN_NOT_EQUALS, "<>")
		}
		return makeToken(lexer, TOKEN_LESS_THAN, "<")
	case '>':
		if peek(lexer) == '=' {
			advance(lexer)
			return makeToken(lexer, TOKEN_GREATER_EQUAL, ">=")
		}
		return makeToken(lexer, TOKEN_GREATER_THAN, ">")
	case '.':
		return makeToken(lexer, TOKEN_DOT, ".")

	case '\'':
		// Start scanning a VARCHAR literal
		for !atEnd(lexer) && peek(lexer) != '\'' {
			advance(lexer)
		}
		if atEnd(lexer) {
			return makeToken(lexer, TOKEN_ILLEGAL, "Unterminated string")
		}
		advance(lexer)                                        // Consume closing quote
		value := lexer.input[lexer.start+1 : lexer.current-1] // Exclude quotes
		return makeToken(lexer, TOKEN_VARCHAR_LITERAL, value)

	case '"':
		return makeToken(lexer, TOKEN_DOUBLE_QUOTE, "\"")
	default:
		return makeToken(lexer, TOKEN_ILLEGAL, string(c))
	}
}

func tokenTypeToString(t TokenType) string {
	switch t {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	case TOKEN_INT_LITERAL:
		return "INT_LITERAL"
	case TOKEN_VARCHAR_LITERAL:
		return "VARCHAR_LITERAL"
	case TOKEN_FLOAT_LITERAL:
		return "FLOAT_LITERAL"
	case TOKEN_BOOLEAN_LITERAL:
		return "BOOLEAN_LITERAL"
	case TOKEN_NULL_LITERAL:
		return "NULL_LITERAL"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_SELECT:
		return "SELECT"
	case TOKEN_CREATE:
		return "CREATE"
	case TOKEN_TABLE:
		return "TABLE"
	case TOKEN_VARCHAR:
		return "VARCHAR"
	case TOKEN_INT:
		return "INT"
	case TOKEN_INSERT:
		return "INSERT"
	case TOKEN_INTO:
		return "INTO"
	case TOKEN_VALUES:
		return "VALUES"
	case TOKEN_FROM:
		return "FROM"
	case TOKEN_AS:
		return "AS"
	case TOKEN_PRIMARY:
		return "PRIMARY"
	case TOKEN_KEY:
		return "KEY"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_COMMA:
		return "COMMA"
	case TOKEN_SEMICOLON:
		return "SEMICOLON"
	case TOKEN_SINGLE_QUOTE:
		return "SINGLE_QUOTE"
	case TOKEN_COUNT:
		return "COUNT"
	case TOKEN_SUM:
		return "SUM"
	case TOKEN_AVG:
		return "AVG"
	case TOKEN_MIN:
		return "MIN"
	case TOKEN_MAX:
		return "MAX"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}
