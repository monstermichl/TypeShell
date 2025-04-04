package lexer

import (
	"fmt"
	"regexp"
	"slices"
)

type TokenType int8
type VarType = string

const (
	UNKNOWN TokenType = iota

	// Commnet.
	COMMENT

	// Brackets.
	OPENING_ROUND_BRACKET
	CLOSING_ROUND_BRACKET
	OPENING_SQUARE_BRACKET
	CLOSING_SQUARE_BRACKET
	OPENING_CURLY_BRACKET
	CLOSING_CURLY_BRACKET

	// Operators.
	ASSIGN_OPERATOR
	BINARY_OPERATOR
	COMPARE_OPERATOR
	LOGICAL_OPERATOR
	SHORT_INIT_OPERATOR

	// Literals.
	BOOL_LITERAL
	NUMBER_LITERAL
	STRING_LITERAL

	// Types.
	DATA_TYPE

	// Separators.
	COMMA
	SEMICOLON
	SPACE
	NEWLINE

	// Identifier.
	IDENTIFIER

	// Keywords.
	VAR_DEFINITION
	FUNCTION_DEFINITION
	RETURN
	IF
	ELSE
	FOR
	RANGE
	LEN
	BREAK
	CONTINUE
	PRINT
	INPUT

	// Pipe.
	PIPE

	// End of file.
	EOF
)

const (
	DATA_TYPE_BOOLEAN VarType = "bool"
	DATA_TYPE_INTEGER VarType = "int"
	DATA_TYPE_STRING  VarType = "string"
)

type Token struct {
	tokenType TokenType
	value     string
	row       int
	column    int
}

func (t Token) Type() TokenType {
	return t.tokenType
}

func (t Token) Value() string {
	return t.value
}

func (t Token) Row() int {
	return t.row
}

func (t Token) Column() int {
	return t.column
}

type tokenMapping struct {
	value     string
	tokenType TokenType
}

var nonAlphabeticTokens = []tokenMapping{
	{"(", OPENING_ROUND_BRACKET},
	{")", CLOSING_ROUND_BRACKET},
	{"[", OPENING_SQUARE_BRACKET},
	{"]", CLOSING_SQUARE_BRACKET},
	{"{", OPENING_CURLY_BRACKET},
	{"}", CLOSING_CURLY_BRACKET},

	{"==", COMPARE_OPERATOR},
	{"!=", COMPARE_OPERATOR},
	{"<=", COMPARE_OPERATOR},
	{">=", COMPARE_OPERATOR},
	{"<", COMPARE_OPERATOR},
	{">", COMPARE_OPERATOR},

	{"&&", LOGICAL_OPERATOR},
	{"||", LOGICAL_OPERATOR},

	{"=", ASSIGN_OPERATOR},

	{":=", SHORT_INIT_OPERATOR},

	{"+", BINARY_OPERATOR},
	{"-", BINARY_OPERATOR},
	{"*", BINARY_OPERATOR},
	{"/", BINARY_OPERATOR},
	{"%", BINARY_OPERATOR},

	{",", COMMA},
	{";", SEMICOLON},
	{" ", SPACE},
	{"\t", SPACE},

	{"|", PIPE},

	{"\n", NEWLINE},
}

var keywords = map[string]TokenType{
	// Common keywords.
	"var":      VAR_DEFINITION,
	"func":     FUNCTION_DEFINITION,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"range":    RANGE,
	"len":      LEN,
	"break":    BREAK,
	"continue": CONTINUE,
	"print":    PRINT,
	"input":    INPUT,

	// Types.
	DATA_TYPE_BOOLEAN: DATA_TYPE,
	DATA_TYPE_INTEGER: DATA_TYPE,
	DATA_TYPE_STRING:  DATA_TYPE,
}

func newToken(value string, tokenType TokenType, row int, column int) Token {
	return Token{
		value:     value,
		tokenType: tokenType,
		row:       row,
		column:    column,
	}
}

func char(s string, position int) string {
	c := ""

	if position < len(s) {
		c = string(s[position])
	}
	return c
}

func Tokenize(source string) ([]Token, error) {
	var err error = nil
	tokens := []Token{}
	i := 0
	startIndex := 1
	row := startIndex
	column := startIndex

	sourceLength := len(source)

	for i < sourceLength {
		var token Token
		c0 := char(source, i)
		ogI := i
		ogRow := row
		ogColumn := column

		if c0 == "\"" {
			// Evaluate string.
			str := ""
			i++

			for i < sourceLength {
				c0 = char(source, i)

				if regexp.MustCompile(`^\\.`).MatchString(source[i:]) {
					// Skip escaped character.
					str += char(source, i)
					i++
				} else if c0 == "\"" {
					// Detected string end.
					i++
					token = newToken(str, STRING_LITERAL, ogRow, ogColumn)
					break
				}
				str += c0
				i++
			}

			if token.tokenType == UNKNOWN {
				err = fmt.Errorf("string at row %d, column %d has not been terminated", ogRow, ogColumn)
				break
			}
		} else if match := regexp.MustCompile(`^\/(\*|\/)`).FindString(source[i:]); match != "" {
			// Evaluate comment.
			oneLineComment := true

			// If second char is a asterisk, it's a multiline comment.
			if string(match[1]) == "*" {
				oneLineComment = false
			}
			i += len(match)
			comment := ""

			for i < sourceLength {
				cancel := false
				c0 = char(source, i)
				c1 := char(source, i+1)

				if oneLineComment && c0 == "\n" {
					cancel = true
				} else if !oneLineComment && c0 == "*" && c1 == "/" {
					i++
					cancel = true
				}
				i++

				if cancel {
					break
				}
				comment += c0
			}
			token = newToken(regexp.MustCompile(`\s+`).ReplaceAllString(comment, " "), COMMENT, ogRow, ogColumn)
		} else if match := regexp.MustCompile(`^true|false`).FindString(source[i:]); match != "" {
			// Create bool token.
			token = newToken(match, BOOL_LITERAL, ogRow, ogColumn)
			i += len(match)

		} else if match := regexp.MustCompile(`^\d+(\.\d+)?`).FindString(source[i:]); match != "" {
			// Create number token.
			token = newToken(match, NUMBER_LITERAL, ogRow, ogColumn)
			i += len(match)
		} else if regexp.MustCompile(`[a-zA-Z_]`).MatchString(c0) {
			identifier := ""

			for {
				c0 = char(source, i)

				if !regexp.MustCompile(`[a-zA-Z0-9_]`).MatchString(c0) {
					break
				}
				identifier += c0
				i++
			}

			// Check if identifier is a keyword.
			tokenType, hasKey := keywords[identifier]

			// If it's not a keyword, it's an identifier.
			if !hasKey {
				tokenType = IDENTIFIER
			}
			token = newToken(identifier, tokenType, ogRow, ogColumn)
		}

		// If no complex token has been found, try to find simple tokens.
		if token.tokenType == UNKNOWN {
			// Try to find non-alphabetic token.
			for _, mapping := range nonAlphabeticTokens {
				key := mapping.value
				tokenType := mapping.tokenType
				endIndex := i + len(key)

				if endIndex <= sourceLength && source[i:endIndex] == key {
					token = newToken(key, tokenType, ogRow, ogColumn)
					i = endIndex
					break
				}
			}
		}

		if token.tokenType == NEWLINE {
			row++
			column = startIndex
		} else {
			column = ogColumn + (i - ogI)
		}

		// If still no token has been found, exit with error.
		if token.tokenType == UNKNOWN {
			err = fmt.Errorf("unknown token \"%s\" at position %d", c0, i)
			break
		} else if slices.Contains([]TokenType{SPACE, COMMENT}, token.tokenType) {
			// Ignore spaces and comments for now.
		} else {
			tokens = append(tokens, token)
		}
	}
	return append(tokens, newToken("", EOF, row, column)), err
}
