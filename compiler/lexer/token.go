package lexer

// token kind
const (
	TokenEOF        = iota         // EOF
	TokenVararg                    // ...
	TokenSepSemi                   // ;
	TokenSepComma                  // ,
	TokenSepDot                    // .
	TokenSepColon                  // :
	TokenSepLabel                  // ::
	TokenSepLparen                 // (
	TokenSepRparen                 // )
	TokenSepLBrack                 // [
	TokenSepRBrack                 // ]
	TokenSepLcurly                 // {
	TokenSepRcurly                 // }
	TokenOpAssign                  // =
	TokenOpMinus                   // -(unm or sub)
	TokenOpWave                    // ~(not or bxor)
	TokenOpAdd                     // +
	TokenOpMul                     // *
	TokenOpDiv                     // /
	TokenOpIDiv                    // //
	TokenOpMod                     // %
	TokenOpPow                     // ^
	TokenOpBand                    // &
	TokenOpBor                     // |
	TokenOpShr                     // >>
	TokenOpShl                     // <<
	TokenOpConcat                  // ..
	TokenOpLt                      // <
	TokenOpLe                      // <=
	TokenOpGt                      // >
	TokenOpGe                      // >=
	TokenOpEq                      // ==
	TokenOpNe                      // ~=
	TokenOpLen                     // #
	TokenOpAnd                     // and
	TokenOpOr                      // or
	TokenOpNot                     // not
	TokenKwBreak                   // break
	TokenKwDo                      // do
	TokenKwElse                    // else
	TokenKwElseif                  // elseif
	TokenKwEnd                     // end
	TokenKwFalse                   // false
	TokenKwFor                     // for
	TokenKwFunction                // function
	TokenKwGoto                    // goto
	TokenKwIf                      // if
	TokenKwIn                      // in
	TokenKwLocal                   // local
	TokenKwNil                     // nil
	TokenKwRepeat                  // repeat
	TokenKwReturn                  // return
	TokenKwThen                    // then
	TokenKwTrue                    // true
	TokenKwUntil                   // until
	TokenKwWhile                   // while
	TokenIdentifier                // identifier
	TokenNumber                    // number literal
	TokenString                    // string literal
	TokenOpUnm      = TokenOpMinus // unary minus
	TokenOpSub      = TokenOpMinus // binary minus
	TokenOpBnot     = TokenOpWave  // not
	TokenOpBxor     = TokenOpWave  // xor
)

var keywords = map[string]int{
	"and":      TokenOpAnd,
	"break":    TokenKwBreak,
	"do":       TokenKwDo,
	"else":     TokenKwElse,
	"elseif":   TokenKwElseif,
	"end":      TokenKwEnd,
	"false":    TokenKwFalse,
	"for":      TokenKwFor,
	"function": TokenKwFunction,
	"goto":     TokenKwGoto,
	"if":       TokenKwIf,
	"in":       TokenKwIn,
	"local":    TokenKwLocal,
	"nil":      TokenKwNil,
	"not":      TokenOpNot,
	"or":       TokenOpOr,
	"repeat":   TokenKwRepeat,
	"return":   TokenKwReturn,
	"then":     TokenKwThen,
	"true":     TokenKwTrue,
	"until":    TokenKwUntil,
	"while":    TokenKwWhile,
}

// KindToString coverts kind int to string
func KindToString(kind int) string {
	switch {
	case kind == TokenEOF:
		return "eof"
	case kind == TokenVararg:
		return "vararg"
	case kind <= TokenSepRcurly:
		return "separator"
	case kind <= TokenOpNot:
		return "operator"
	case kind <= TokenKwWhile:
		return "keyword"
	case kind == TokenIdentifier:
		return "identifier"
	case kind == TokenNumber:
		return "number"
	case kind == TokenString:
		return "string"
	default:
		return "other"
	}
}
