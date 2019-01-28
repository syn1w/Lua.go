package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	lex := NewLexer(`;,()[]{}+-*^%&|#`, "string")
	testNextTokenKind(t, lex, TokenSepSemi)
	testNextTokenKind(t, lex, TokenSepComma)
	testNextTokenKind(t, lex, TokenSepLparen)
	testNextTokenKind(t, lex, TokenSepRparen)
	testNextTokenKind(t, lex, TokenSepLBrack)
	testNextTokenKind(t, lex, TokenSepRBrack)
	testNextTokenKind(t, lex, TokenSepLcurly)
	testNextTokenKind(t, lex, TokenSepRcurly)
	testNextTokenKind(t, lex, TokenOpAdd)
	testNextTokenKind(t, lex, TokenOpMinus)
	testNextTokenKind(t, lex, TokenOpMul)
	testNextTokenKind(t, lex, TokenOpPow)
	testNextTokenKind(t, lex, TokenOpMod)
	testNextTokenKind(t, lex, TokenOpBand)
	testNextTokenKind(t, lex, TokenOpBor)
	testNextTokenKind(t, lex, TokenOpLen)
	testNextTokenKind(t, lex, TokenEOF)
}

func TestNextToken2(t *testing.T) {
	lex := NewLexer(`... .. . :: : // / ~= ~ == = <<<= <>> >= >`, "string")
	testNextTokenKind(t, lex, TokenVararg)
	testNextTokenKind(t, lex, TokenOpConcat)
	testNextTokenKind(t, lex, TokenSepDot)
	testNextTokenKind(t, lex, TokenSepLabel)
	testNextTokenKind(t, lex, TokenSepColon)
	testNextTokenKind(t, lex, TokenOpIDiv)
	testNextTokenKind(t, lex, TokenOpDiv)
	testNextTokenKind(t, lex, TokenOpNe)
	testNextTokenKind(t, lex, TokenOpWave)
	testNextTokenKind(t, lex, TokenOpEq)
	testNextTokenKind(t, lex, TokenOpAssign)
	testNextTokenKind(t, lex, TokenOpShl)
	testNextTokenKind(t, lex, TokenOpLe)
	testNextTokenKind(t, lex, TokenOpLt)
	testNextTokenKind(t, lex, TokenOpShr)
	testNextTokenKind(t, lex, TokenOpGe)
	testNextTokenKind(t, lex, TokenOpGt)
	testNextTokenKind(t, lex, TokenEOF)
}

func TestNextTokenKeyword(t *testing.T) {
	kws := `and       break     do        else      elseif    end
	false     for       function  goto      if        in
	local     nil       not       or        repeat    return
	then      true      until     while`

	lex := NewLexer(kws, "string")
	testNextTokenKind(t, lex, TokenOpAnd)
	testNextTokenKind(t, lex, TokenKwBreak)
	testNextTokenKind(t, lex, TokenKwDo)
	testNextTokenKind(t, lex, TokenKwElse)
	testNextTokenKind(t, lex, TokenKwElseif)
	testNextTokenKind(t, lex, TokenKwEnd)
	testNextTokenKind(t, lex, TokenKwFalse)
	testNextTokenKind(t, lex, TokenKwFor)
	testNextTokenKind(t, lex, TokenKwFunction)
	testNextTokenKind(t, lex, TokenKwGoto)
	testNextTokenKind(t, lex, TokenKwIf)
	testNextTokenKind(t, lex, TokenKwIn)
	testNextTokenKind(t, lex, TokenKwLocal)
	testNextTokenKind(t, lex, TokenKwNil)
	testNextTokenKind(t, lex, TokenOpNot)
	testNextTokenKind(t, lex, TokenOpOr)
	testNextTokenKind(t, lex, TokenKwRepeat)
	testNextTokenKind(t, lex, TokenKwReturn)
	testNextTokenKind(t, lex, TokenKwThen)
	testNextTokenKind(t, lex, TokenKwTrue)
	testNextTokenKind(t, lex, TokenKwUntil)
	testNextTokenKind(t, lex, TokenKwWhile)
	testNextTokenKind(t, lex, TokenEOF)
}

func TestNextTokenIdentifier(t *testing.T) {
	ids := `_ __ a helloWorld hello42 HELLO`
	lex := NewLexer(ids, "string")
	testNextIdentifier(t, lex, "_")
	testNextIdentifier(t, lex, "__")
	testNextIdentifier(t, lex, "a")
	testNextIdentifier(t, lex, "helloWorld")
	testNextIdentifier(t, lex, "hello42")
	testNextIdentifier(t, lex, "HELLO")
	testNextTokenKind(t, lex, TokenEOF)
}

func TestNextTokenNumber(t *testing.T) {
	numbers := `
	3   42   0xff   0xBEBADA
	3.0     3.1416     314.16e-2     0.31416E1     34e1
	0x0.1E  0xA23p-4   0X1.921FB54442D18P+1
	3.	.3	00001`

	lex := NewLexer(numbers, "string")
	testNextNumber(t, lex, "3")
	testNextNumber(t, lex, "42")
	testNextNumber(t, lex, "0xff")
	testNextNumber(t, lex, "0xBEBADA")
	testNextNumber(t, lex, "3.0")
	testNextNumber(t, lex, "3.1416")
	testNextNumber(t, lex, "314.16e-2")
	testNextNumber(t, lex, "0.31416E1")
	testNextNumber(t, lex, "34e1")
	testNextNumber(t, lex, "0x0.1E")
	testNextNumber(t, lex, "0xA23p-4")
	testNextNumber(t, lex, "0X1.921FB54442D18P+1")
	testNextNumber(t, lex, "3.")
	testNextNumber(t, lex, ".3")
	testNextNumber(t, lex, "00001")
	testNextTokenKind(t, lex, TokenEOF)
}

func TestNextTokenComments(t *testing.T) {
	comments := `
	--
	--[[]]
	a -- short comment
	+ --[[ long comment]] b --[===[
		long comment
	]===] - c
	--`

	lex := NewLexer(comments, "string")
	testNextIdentifier(t, lex, "a")
	testNextTokenKind(t, lex, TokenOpAdd)
	testNextIdentifier(t, lex, "b")
	testNextTokenKind(t, lex, TokenOpMinus)
	testNextIdentifier(t, lex, "c")
	testNextTokenKind(t, lex, TokenEOF)
}

func TestNextTokenString(t *testing.T) {
	strs := `
	[[]] [[ long string ]]
	[=[
long string]=]
	[===[long\z
	string]===]
	'''"''short string'
	"""'""short string"
	'\a\b\f\n\r\t\v\\\"\''
	"\8 \08 \64 \122 \x08 \x7A \u{6211} zzz"
	'foo \z
	

	bar'
	`
	lex := NewLexer(strs, "string")
	testNextString(t, lex, "")
	testNextString(t, lex, " long string ")
	testNextString(t, lex, "long string")
	testNextString(t, lex, "long\\z\n\tstring")
	testNextString(t, lex, "")
	testNextString(t, lex, "\"")
	testNextString(t, lex, "short string")
	testNextString(t, lex, "")
	testNextString(t, lex, "'")
	testNextString(t, lex, "short string")
	testNextString(t, lex, "\a\b\f\n\r\t\v\\\"'")
	testNextString(t, lex, "\b \b @ z \b z 我 zzz")
	testNextString(t, lex, "foo bar")
	testNextTokenKind(t, lex, TokenEOF)
	if line := lex.Line(); line != 15 {
		t.Errorf("line failed: want='%d', got='%d'", 15, line)
	}
}

func TestHelloWorld(t *testing.T) {
	src := `print("hello world")`
	lex := NewLexer(src, "string")
	testNextIdentifier(t, lex, "print")
	testNextTokenKind(t, lex, TokenSepLparen)
	testNextString(t, lex, "hello world")
	testNextTokenKind(t, lex, TokenSepRparen)
	testNextTokenKind(t, lex, TokenEOF)
}

func TestLookAhead(t *testing.T) {
	src := `print("hello world")`
	lex := NewLexer(src, "string")

	testEqualKind(t, lex.LookAhead(), TokenIdentifier)
	lex.NextToken() // print
	testEqualKind(t, lex.LookAhead(), TokenSepLparen)
	lex.NextToken() // `(`
	testEqualKind(t, lex.LookAhead(), TokenString)
	lex.NextToken() // "hello world"
	testEqualKind(t, lex.LookAhead(), TokenSepRparen)
	lex.NextToken() // `)`
	testEqualKind(t, lex.LookAhead(), TokenEOF)
}

func TestErrors(t *testing.T) {
	testError(t, "?", "src:1: unexpected symbol near '?'")
	testError(t, "[===", "src:1: invalid long string delimiter near '[='")
	testError(t, "[==[xx", "src:1: unfinished long string or comment")
	testError(t, "'abc\\defg", "src:1: unfinished short string")
	testError(t, "'abc\\defg'", "src:1: invalid escape sequence near '\\d'")
	testError(t, "'\\256'", "src:1: decimal escape too large near '\\256'")
	testError(t, "'\\u{11FFFF}'", "src:1: UTF-8 value too large near '\\u{11FFFF}'")
	testError(t, "'\\'", "src:1: unfinished short string")
}

func testNextTokenKind(t *testing.T, lex *Lexer, want int) {
	line, got, token := lex.NextToken()
	if got != want {
		t.Errorf("next kind failed, near: '%d:%s', want='%s', got='%s'", line, token,
			kindToString(want), kindToString(got))
	}
}

func testNextIdentifier(t *testing.T, lex *Lexer, want string) {
	_, got := lex.AssertNextTokenKind(TokenIdentifier)
	if got != want {
		t.Errorf("next identifier failed: want='%s', got='%s'", want, got)
	}
}

func testNextNumber(t *testing.T, lex *Lexer, want string) {
	_, got := lex.AssertNextTokenKind(TokenNumber)
	if got != want {
		t.Errorf("next number failed: want='%s', got='%s'", want, got)
	}
}

func testNextString(t *testing.T, lex *Lexer, want string) {
	_, got := lex.AssertNextTokenKind(TokenString)
	if got != want {
		t.Errorf("next string failed: want='%s', got='%s'", want, got)
	}
}

func testEqualKind(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("want='%s', got='%s'",
			kindToString(want), kindToString(got))
	}
}

func testEqualString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func testError(t *testing.T, chunk, want string) {
	got := safeNextToken(NewLexer(chunk, "src"))
	testEqualString(t, got, want)
}

func safeNextToken(lex *Lexer) (err string) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(string)
		}
	}()
	_, _, err = lex.NextToken()
	return
}

// detailed token kind information for testing
func kindToString(kind int) string {
	switch kind {
	case TokenVararg:
		return "..."
	case TokenSepSemi:
		return ";"
	case TokenSepComma:
		return ","
	case TokenSepDot:
		return "."
	case TokenSepColon:
		return ":"
	case TokenSepLabel:
		return "::"
	case TokenSepLparen:
		return "("
	case TokenSepRparen:
		return ")"
	case TokenSepLBrack:
		return "["
	case TokenSepRBrack:
		return "]"
	case TokenSepLcurly:
		return "{"
	case TokenSepRcurly:
		return "}"
	case TokenOpAssign:
		return "="
	case TokenOpMinus:
		return "-(unm or sub)"
	case TokenOpWave:
		return "~(not or bxor)"
	case TokenOpAdd:
		return "+"
	case TokenOpMul:
		return "*"
	case TokenOpDiv:
		return "/"
	case TokenOpIDiv:
		return "//"
	case TokenOpMod:
		return "%"
	case TokenOpPow:
		return "^"
	case TokenOpBand:
		return "&"
	case TokenOpBor:
		return "|"
	case TokenOpShr:
		return ">>"
	case TokenOpShl:
		return "<<"
	case TokenOpConcat:
		return ".."
	case TokenOpLt:
		return "<"
	case TokenOpLe:
		return "<="
	case TokenOpGt:
		return ">"
	case TokenOpGe:
		return ">="
	case TokenOpEq:
		return "=="
	case TokenOpNe:
		return "~="
	case TokenOpLen:
		return "#"
	case TokenOpAnd:
		return "and"
	case TokenOpOr:
		return "or"
	case TokenOpNot:
		return "not"
	case TokenKwBreak:
		return "break"
	case TokenKwDo:
		return "do"
	case TokenKwElse:
		return "else"
	case TokenKwElseif:
		return "elseif"
	case TokenKwEnd:
		return "end"
	case TokenKwFalse:
		return "false"
	case TokenKwFor:
		return "for"
	case TokenKwFunction:
		return "function"
	case TokenKwGoto:
		return "goto"
	case TokenKwIf:
		return "if"
	case TokenKwIn:
		return "in"
	case TokenKwLocal:
		return "local"
	case TokenKwNil:
		return "nil"
	case TokenKwRepeat:
		return "repeat"
	case TokenKwReturn:
		return "return"
	case TokenKwThen:
		return "then"
	case TokenKwTrue:
		return "true"
	case TokenKwUntil:
		return "until"
	case TokenKwWhile:
		return "while"
	case TokenIdentifier:
		return "identifier"
	case TokenNumber:
		return "number literal"
	case TokenString:
		return "string literal"
	default:
		return "unknown"
	}
}
