package lexer

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
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
	COMPOUND_ASSIGN_OPERATOR
	UNARY_OPERATOR
	BINARY_OPERATOR

	// Literals.
	BOOL_LITERAL
	NUMBER_LITERAL
	STRING_LITERAL
	NIL_LITERAL

	// Separators.
	COMMA
	COLON
	SEMICOLON
	DOT
	SPACE
	NEWLINE

	// Identifier.
	IDENTIFIER

	// Keywords.
	SECTION_KEYWORD
	KEYWORD

	// App operators.
	AT
	PIPE

	// End of file.
	EOF
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

func (t Token) IsKeyword(keyword Keyword) bool {
	return slices.Contains([]TokenType{KEYWORD, SECTION_KEYWORD}, t.Type()) && t.Value() == keyword
}

type tokenMapping struct {
	value     string
	tokenType TokenType
}

type tokenMappings []tokenMapping

func (m tokenMappings) find(value string) (TokenType, bool) {
	for _, entry := range m {
		if entry.value == value {
			return entry.tokenType, true
		}
	}
	return UNKNOWN, false
}

type nonAlphabeticToken = string

const (
	// Operators highest to lowest priority (top to bottom).
	UnaryOperatorNegate nonAlphabeticToken = "!"

	OperatorMultiplication nonAlphabeticToken = "*"
	OperatorDivision       nonAlphabeticToken = "/"
	OperatorModulo         nonAlphabeticToken = "%"

	OperatorAddition    nonAlphabeticToken = "+"
	OperatorSubtraction nonAlphabeticToken = "-"

	ComparisonLessOrEqual    nonAlphabeticToken = "<="
	ComparisonGreaterOrEqual nonAlphabeticToken = ">="
	ComparisonLess           nonAlphabeticToken = "<"
	ComparisonGreater        nonAlphabeticToken = ">"

	ComparisonEqual    nonAlphabeticToken = "=="
	ComparisonNotEqual nonAlphabeticToken = "!="

	OperatorAnd nonAlphabeticToken = "&&"
	OperatorOr  nonAlphabeticToken = "||"

	OperatorAssign      nonAlphabeticToken = "="
	OperatorShortAssign nonAlphabeticToken = ":="
)

var nonAlphabeticTokens = tokenMappings{
	{"(", OPENING_ROUND_BRACKET},
	{")", CLOSING_ROUND_BRACKET},
	{"[", OPENING_SQUARE_BRACKET},
	{"]", CLOSING_SQUARE_BRACKET},
	{"{", OPENING_CURLY_BRACKET},
	{"}", CLOSING_CURLY_BRACKET},

	{ComparisonEqual, BINARY_OPERATOR},
	{ComparisonNotEqual, BINARY_OPERATOR},
	{ComparisonLessOrEqual, BINARY_OPERATOR},
	{ComparisonGreaterOrEqual, BINARY_OPERATOR},
	{ComparisonLess, BINARY_OPERATOR},
	{ComparisonGreater, BINARY_OPERATOR},

	{OperatorAnd, BINARY_OPERATOR},
	{OperatorOr, BINARY_OPERATOR},

	{"+=", COMPOUND_ASSIGN_OPERATOR},
	{"-=", COMPOUND_ASSIGN_OPERATOR},
	{"*=", COMPOUND_ASSIGN_OPERATOR},
	{"/=", COMPOUND_ASSIGN_OPERATOR},
	{"%=", COMPOUND_ASSIGN_OPERATOR},

	{OperatorShortAssign, ASSIGN_OPERATOR},
	{OperatorAssign, ASSIGN_OPERATOR},

	{UnaryOperatorNegate, UNARY_OPERATOR},

	{"++", UNARY_OPERATOR},
	{"--", UNARY_OPERATOR},

	{OperatorAddition, BINARY_OPERATOR},
	{OperatorSubtraction, BINARY_OPERATOR},
	{OperatorMultiplication, BINARY_OPERATOR},
	{OperatorDivision, BINARY_OPERATOR},
	{OperatorModulo, BINARY_OPERATOR},

	{",", COMMA},
	{":", COLON},
	{";", SEMICOLON},
	{".", DOT},
	{" ", SPACE},
	{"\t", SPACE},

	{"@", AT},
	{"|", PIPE},

	{"\n", NEWLINE},
}

type Keyword = string

const (
	// Common keywords.
	KeywordImport   Keyword = "import"
	KeywordType     Keyword = "type"
	KeywordStruct   Keyword = "struct"
	KeywordConst    Keyword = "const"
	KeywordVar      Keyword = "var"
	KeywordFunc     Keyword = "func"
	KeywordReturn   Keyword = "return"
	KeywordIf       Keyword = "if"
	KeywordElse     Keyword = "else"
	KeywordSwitch   Keyword = "switch"
	KeywordCase     Keyword = "case"
	KeywordDefault  Keyword = "default"
	KeywordFor      Keyword = "for"
	KeywordRange    Keyword = "range"
	KeywordBreak    Keyword = "break"
	KeywordContinue Keyword = "continue"
	KeywordIota     Keyword = "iota"
	KeywordNil      Keyword = "nil"

	// Builtin functions.
	KeywordLen    Keyword = "len"
	KeywordPrint  Keyword = "print"
	KeywordInput  Keyword = "input"
	KeywordCopy   Keyword = "copy"
	KeywordItoa   Keyword = "itoa"
	KeywordExists Keyword = "exists"
	KeywordRead   Keyword = "read"
	KeywordWrite  Keyword = "write"
	KeywordPanic  Keyword = "panic"
	KeywordUnsafe Keyword = "unsafe"
)

var keywords = tokenMappings{
	// Common keywords.
	{KeywordImport, SECTION_KEYWORD},
	{KeywordType, SECTION_KEYWORD},
	{KeywordStruct, KEYWORD},
	{KeywordConst, SECTION_KEYWORD},
	{KeywordVar, SECTION_KEYWORD},
	{KeywordFunc, SECTION_KEYWORD},
	{KeywordReturn, KEYWORD},
	{KeywordIf, SECTION_KEYWORD},
	{KeywordElse, SECTION_KEYWORD},
	{KeywordSwitch, SECTION_KEYWORD},
	{KeywordCase, KEYWORD},
	{KeywordDefault, KEYWORD},
	{KeywordFor, SECTION_KEYWORD},
	{KeywordRange, KEYWORD},
	{KeywordBreak, KEYWORD},
	{KeywordContinue, KEYWORD},
	{KeywordIota, KEYWORD},
	{KeywordNil, KEYWORD},

	// Builtin functions.
	{KeywordLen, KEYWORD},
	{KeywordPrint, KEYWORD},
	{KeywordInput, KEYWORD},
	{KeywordCopy, KEYWORD},
	{KeywordItoa, KEYWORD},
	{KeywordExists, KEYWORD},
	{KeywordRead, KEYWORD},
	{KeywordWrite, KEYWORD},
	{KeywordPanic, KEYWORD},
	{KeywordUnsafe, KEYWORD},
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

	source = strings.ReplaceAll(source, "\r\n", "\n") // Replace all \r\n with \n.
	sourceLength := len(source)

	for i < sourceLength {
		var token Token
		c0 := char(source, i)
		ogI := i
		ogRow := row
		ogColumn := column

		if raw := c0 == "`"; raw || c0 == `"` {
			// Evaluate string.
			str := ""
			i++

			for i < sourceLength {
				c0 = char(source, i)
				appended := false

				if match := regexp.MustCompile(`^\\.`).FindString(source[i:]); !raw && match != "" {
					// Convert escaped character to be a control character (https://pkg.go.dev/strconv#Unquote).
					parsed, err := strconv.Unquote(fmt.Sprintf(`"%s"`, match))

					if err != nil {
						return nil, fmt.Errorf(`invalid escape sequence "%s"`, match)
					}
					str += parsed
					i += len(match)
					appended = true
				} else if (raw && c0 == "`") || (!raw && c0 == `"`) {
					// Detected string end.
					i++
					token = newToken(str, STRING_LITERAL, ogRow, ogColumn)
					break
				}

				if !appended {
					if c0 == "\n" {
						row++
						ogColumn = startIndex
					}
					str += c0
					i++
				}
			}

			if token.tokenType == UNKNOWN {
				err = fmt.Errorf("string at row %d, column %d has not been terminated", ogRow, ogColumn)
				break
			}
		} else if matches := regexp.MustCompile(`(?s)^\/\*(.*)\*\/`).FindStringSubmatch(source[i:]); matches != nil {
			// Multiline comment.
			token = newToken(matches[1], COMMENT, ogRow, ogColumn)
			match := matches[0]
			lines := strings.Split(match, "\n")
			lastLinesIndex := len(lines) - 1
			row += lastLinesIndex
			ogColumn = startIndex
			i += len(match)
			ogI = i - len(lines[lastLinesIndex])
		} else if matches := regexp.MustCompile(`^\/\/(.*)`).FindStringSubmatch(source[i:]); matches != nil {
			// Single line comment.
			token = newToken(matches[1], COMMENT, ogRow, ogColumn)
			i += len(matches[0])
		} else if match := regexp.MustCompile(`^(true|false)`).FindString(source[i:]); match != "" {
			// Create bool token.
			token = newToken(match, BOOL_LITERAL, ogRow, ogColumn)
			i += len(match)
		} else if match := regexp.MustCompile(`^(((\d(o|O|x|X|b|B))(\d+))|(-?\d+(\.\d+)?))`).FindString(source[i:]); match != "" {
			parsed, err := strconv.ParseInt(strings.ToLower(match), 0, 32)

			if err == nil {
				// Create number token.
				token = newToken(strconv.Itoa(int(parsed)), NUMBER_LITERAL, ogRow, ogColumn)
				i += len(match)
			}
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
			tokenType, hasKey := keywords.find(identifier)

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
			err = fmt.Errorf(`unknown token "%s" at position %d`, c0, i)
			break
		} else if slices.Contains([]TokenType{SPACE, COMMENT}, token.tokenType) {
			// Ignore spaces and comments for now.
		} else {
			tokens = append(tokens, token)
		}
	}
	return append(tokens, newToken("", EOF, row, column)), err
}
