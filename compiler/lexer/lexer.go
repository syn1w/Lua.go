package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://cloudwu.github.io/lua53doc/manual.html#3.1

// Lexer lua lexer
type Lexer struct {
	chunk     string // source code
	chunkName string // source file name
	line      int    // current line no.

	nextToken     string
	nextTokenKind int
	nextTokenLine int
}

// NewLexer new lua lexer
func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{
		chunk:     chunk,
		chunkName: chunkName,
		line:      1,
	}
}

// NextToken returns next token
func (lex *Lexer) NextToken() (line, kind int, token string) {
	if lex.nextTokenLine > 0 {
		line = lex.nextTokenLine
		kind = lex.nextTokenKind
		token = lex.nextToken
		lex.line = lex.nextTokenLine
		lex.nextTokenLine = 0
		return
	}

	lex.skipWhiteSpaces()
	if len(lex.chunk) == 0 {
		return lex.line, TokenEOF, "EOF"
	}

	switch lex.chunk[0] {
	case ';':
		lex.next(1)
		return lex.line, TokenSepSemi, ";"
	case ',':
		lex.next(1)
		return lex.line, TokenSepComma, ","
	case '(':
		lex.next(1)
		return lex.line, TokenSepLparen, "("
	case ')':
		lex.next(1)
		return lex.line, TokenSepRparen, ")"
	case ']':
		lex.next(1)
		return lex.line, TokenSepRBrack, "]"
	case '{':
		lex.next(1)
		return lex.line, TokenSepLcurly, "{"
	case '}':
		lex.next(1)
		return lex.line, TokenSepRcurly, "}"
	case '+':
		lex.next(1)
		return lex.line, TokenOpAdd, "+"
	case '-':
		lex.next(1)
		return lex.line, TokenOpMinus, "-"
	case '*':
		lex.next(1)
		return lex.line, TokenOpMul, "*"
	case '^':
		lex.next(1)
		return lex.line, TokenOpPow, "^"
	case '%':
		lex.next(1)
		return lex.line, TokenOpMod, "%"
	case '&':
		lex.next(1)
		return lex.line, TokenOpBand, "&"
	case '|':
		lex.next(1)
		return lex.line, TokenOpBor, "|"
	case '#':
		lex.next(1)
		return lex.line, TokenOpLen, "#"
	case ':':
		if lex.test("::") {
			lex.next(2)
			return lex.line, TokenSepLabel, "::"
		}
		lex.next(1)
		return lex.line, TokenSepColon, ":"
	case '/':
		if lex.test("//") {
			lex.next(2)
			return lex.line, TokenOpIDiv, "//"
		}
		lex.next(1)
		return lex.line, TokenOpDiv, "/"
	case '~':
		if lex.test("~=") {
			lex.next(2)
			return lex.line, TokenOpNe, "~="
		}
		return lex.line, TokenOpWave, "~"
	case '=':
		if lex.test("==") {
			lex.next(2)
			return lex.line, TokenOpEq, "=="
		}
		lex.next(1)
		return lex.line, TokenOpAssign, "="
	case '<':
		if lex.test("<<") {
			lex.next(2)
			return lex.line, TokenOpShl, "<<"
		} else if lex.test("<=") {
			lex.next(2)
			return lex.line, TokenOpLe, "<="
		} else {
			lex.next(1)
			return lex.line, TokenOpLt, "<"
		}
	case '>':
		if lex.test(">>") {
			lex.next(2)
			return lex.line, TokenOpShr, ">>"
		} else if lex.test(">=") {
			lex.next(2)
			return lex.line, TokenOpGe, ">="
		} else {
			lex.next(1)
			return lex.line, TokenOpGt, ">"
		}
	case '.':
		if lex.test("...") {
			lex.next(3)
			return lex.line, TokenVararg, "..."
		} else if lex.test("..") {
			lex.next(2)
			return lex.line, TokenOpConcat, ".."
		} else if len(lex.chunk) == 1 || !isDigit(lex.chunk[1]) {
			lex.next(1)
			return lex.line, TokenSepDot, "."
		}
	case '[':
		if lex.test("[[") || lex.test("[=") {
			return lex.line, TokenString, lex.scanLongString()
		}
		lex.next(1)
		return lex.line, TokenSepLBrack, "["
	case '\'', '"':
		return lex.line, TokenString, lex.scanShortString()
	}

	c := lex.chunk[0]
	if c == '.' || isDigit(c) {
		token := lex.scanNumber()
		return lex.line, TokenNumber, token
	}

	if c == '_' || isLetter(c) {
		token := lex.scanIdentifier()
		if kind, found := keywars[token]; found {
			return lex.line, kind, token
		}
		return lex.line, TokenIdentifier, token
	}

	lex.error("unexpected symbol near '%q'", c)
	return
}

// LookAhead caches next token and returns the of next token
func (lex *Lexer) LookAhead() int {
	if lex.nextTokenLine > 0 {
		return lex.nextTokenKind
	}

	currentLine := lex.line
	line, kind, token := lex.NextToken()
	lex.line = currentLine
	lex.nextTokenLine = line
	lex.nextTokenKind = kind
	lex.nextToken = token
	return kind
}

// Line return current line of lexer
func (lex *Lexer) Line() int {
	return lex.line
}

// AssertNextTokenKind extracts the token of the specified type
func (lex *Lexer) AssertNextTokenKind(k int) (line int, token string) {
	line, kind, token := lex.NextToken()
	if kind != k {
		lex.error("syntax error near '%s'", token)
	}
	return line, token
}

var reLeftLongBracket = regexp.MustCompile(`^\[=*\[`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")

// (?s): s modifier: single line. Dot matches newline characters
// (Alternative1)|(Alternative2)
//   Alternative1: (^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')
//     ^' matches the character ' literally at start of a line
//     (\\\\|\\'|\\\n|\\z\s*|[^'\n])*
//       \\ \\   matches \ \
//       \\ '    matches \ '
//       \\ \n   matches \ \n
//       \\ z \s* matches \ z space*
//       [^'\n] matches a single character not ' or \n
// Alternative2: (^"(\\\\|\\'|\\\n|\\z\s*|[^'\n])*")
// similar to the former
var reShortString = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)

func (lex *Lexer) skipComment() {
	lex.next(2) // skip --

	// long string comment, mutil-line comment
	if reLeftLongBracket.FindString(lex.chunk) != "" {
		lex.scanLongString()
		return
	}

	// short string comment, single-line comment
	for len(lex.chunk) > 0 && !isNewLine(lex.chunk[0]) {
		lex.next(1)
	}
}

func (lex *Lexer) scanShortString() string {
	if str := reShortString.FindString(lex.chunk); str != "" {
		lex.next(len(str))
		str = str[1 : len(str)-1] // "xxx" -> xxx
		if strings.Index(str, `\`) >= 0 {
			lex.line += len(reNewLine.FindAllString(str, -1))
			str = lex.escape(str)
		}
		return str
	}

	lex.error("unfinished short string")
	return ""
}

func (lex *Lexer) scanLongString() string {
	leftLongBracket := reLeftLongBracket.FindString(lex.chunk)
	if leftLongBracket == "" {
		lex.error("invalid long string delimiter hear '%s'", lex.chunk[0:2])
	}
	rightLongBracket := strings.Replace(leftLongBracket, "[", "]", -1)
	rightLongBracketIdx := strings.Index(lex.chunk, rightLongBracket)
	if rightLongBracketIdx < 0 {
		lex.error("unfinished long string or comment")
	}
	str := lex.chunk[len(leftLongBracket):rightLongBracketIdx]
	lex.next(rightLongBracketIdx + len(rightLongBracket))
	str = reNewLine.ReplaceAllString(str, "\n")
	lex.line += strings.Count(str, "\n")
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}
	return str
}

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

func (lex *Lexer) escape(str string) string {
	var buf bytes.Buffer
	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}
		// str[0] == '\\'
		if len(str) == 1 {
			lex.error("unfinished string: escape character")
		}

		switch str[1] {
		case 'a': // alert bell
			buf.WriteByte('\a')
			str = str[2:]
			continue
		case 'b': // back space
			buf.WriteByte('\b')
			str = str[2:]
			continue
		case 'f': // form feed
			buf.WriteByte('\f')
			str = str[2:]
			continue
		case 'n', '\n': // line feed or newline
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case 'r': // carriage return
			buf.WriteByte('\r')
			str = str[2:]
			continue
		case 't': // horizontal tab
			buf.WriteByte('\t')
			str = str[2:]
			continue
		case 'v': // vertical tab
			buf.WriteByte('\v')
			str = str[2:]
			continue
		case '"':
			buf.WriteByte('"')
			str = str[2:]
			continue
		case '\'':
			buf.WriteByte('\'')
			str = str[2:]
			continue
		case '\\':
			buf.WriteByte('\\')
			str = str[2:]
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
			if found := reDecEscapeSeq.FindString(lex.chunk); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d <= 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				lex.error("decimal escape too large near '%s'", found)
			}
		case 'x': // \xhh
			if found := reHexEscapeSeq.FindString(lex.chunk); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		case 'u': // \u{hhh}
			if found := reUnicodeEscapeSeq.FindString(lex.chunk); found != "" {
				d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				lex.error("UTF-8 value too large near '%s'", found)
			}
		case 'z': // \z skips the following span of white-space characters, including line breaks
			str = str[2:] // skips '\\' and 'z'
			for len(str) > 0 && isWhiteSpace(str[0]) {
				str = str[1:]
			}
			continue
		}
		lex.error("invalid escape sequence near '%q'", str[1])
	}

	return buf.String()
}

// Alternative1 | Alternative2
// Alternative1: ^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?
//   matches floating point number
//   0[xX] matches 0x or 0X
//   [0-9a-fA-F]* matches integer part
//   (\.[0-9a-fA-F]*)? matches decimals part
//   ([pP][+\-]?[0-9]+)? matches power part, such as 0xA23p-4
// Alternative2: ^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?
//   matches integer number
//   [0-9]* matches integer part
//   (\.[0-9]*)? matches decimals part
//   ([eE][+\-]?[0-9]+)? matches exponent part
var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)

func (lex *Lexer) scanNumber() string {
	return lex.scan(reNumber)
}

var reIdentifier = regexp.MustCompile(`^[_a-zA-Z][_\w]*`)

func (lex *Lexer) scanIdentifier() string {
	return lex.scan(reIdentifier)
}

func (lex *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(lex.chunk); token != "" {
		lex.next(len(token))
		return token
	}
	panic("unreachable")
}

func (lex *Lexer) test(str string) bool {
	return strings.HasPrefix(lex.chunk, str)
}

func (lex *Lexer) next(n int) {
	lex.chunk = lex.chunk[n:]
}

func (lex *Lexer) error(format string, args ...interface{}) {
	err := fmt.Sprintf("%s:%d: %s", lex.chunkName, lex.line,
		fmt.Sprintf(format, args...))
	panic(err)
}

func (lex *Lexer) skipWhiteSpaces() {
	for len(lex.chunk) > 0 {
		if lex.test("--") {
			lex.skipComment()
		} else if lex.test("\r\n") || lex.test("\n\r") {
			lex.next(2)
			lex.line++
		} else if isNewLine(lex.chunk[0]) {
			lex.next(1)
			lex.line++
		} else if isWhiteSpace(lex.chunk[0]) {
			lex.next(1)
		} else {
			break
		}
	}
}

func isWhiteSpace(c byte) bool {
	// \v vertical tab
	// \f form feed
	switch c {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}
	return false
}

func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' || c <= 'Z')
}
