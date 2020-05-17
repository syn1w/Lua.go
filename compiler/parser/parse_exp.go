package parser

import (
	"luago/compiler/ast"
	"luago/compiler/lexer"
	"luago/number"
)

// sort by priority
// exp ::= exp12
// exp12 ::= exp11 { `or` exp11 }
// exp11 ::= exp10 { `and` exp10 }
// exp10 ::= exp9 { (`<`|`>`|`<=`|`>=`|`==`|`~=`) exp9 }
// exp9  ::= exp8 { `|` exp8 }
// exp8  ::= exp7 { `~` exp7 }   // xor
// exp7  ::= exp6 { `&` exp6 }
// exp6  ::= exp5 { (`<<`|`>>`) exp5 }
// exp5  ::= exp4 { `..` exp4 }
// exp4  ::= exp3 { (`+`|`-`) exp3 }
// exp3  ::= exp2 { (`*`|`/`|`//`|`%`) exp2 }
// exp2  ::= {(`not`|`#`|`-`|`~`)} exp1
// exp1  ::= exp0 { `^` exp0 }
// exp0  ::= `nil` | `false` | `true` | Numeral | LiteralString |
//           `...` | functiondef | prefixexp | tableconstructor

// exp ::=  `nil` | `false` | `true` | Numeral | LiteralString | `...` | functiondef |
//          prefixexp | tableconstructor | exp binop exp | unop exp
func parseExpList(lex *lexer.Lexer) []ast.Exp {
	exps := make([]ast.Exp, 0, 4)
	exps = append(exps, parseExp(lex))

	for lex.LookAhead() == lexer.TokenSepComma {
		lex.NextToken() // skips the `,`
		exps = append(exps, parseExp(lex))
	}
	return exps
}

func parseExp(lex *lexer.Lexer) ast.Exp {
	return parseExp12(lex)
}

// exp11 {`or` exp11}
func parseExp12(lex *lexer.Lexer) ast.Exp {
	exp := parseExp11(lex)
	for lex.LookAhead() == lexer.TokenOpOr {
		line, op, _ := lex.NextToken() // `or`
		lOr := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp11(lex)}
		exp = optimizeLogicalOr(lOr)
	}

	return exp
}

// exp10 {`and` exp10}
func parseExp11(lex *lexer.Lexer) ast.Exp {
	exp := parseExp10(lex)
	for lex.LookAhead() == lexer.TokenOpAnd {
		line, op, _ := lex.NextToken() // `and`
		lAnd := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp10(lex)}
		exp = optimizeLogicalAnd(lAnd)
	}

	return exp
}

// exp9 {compare exp9}
func parseExp10(lex *lexer.Lexer) ast.Exp {
	exp := parseExp9(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TokenOpLt, lexer.TokenOpGt, lexer.TokenOpLe,
			lexer.TokenOpGe, lexer.TokenOpEq, lexer.TokenOpNe:
			line, op, _ := lex.NextToken() // comp op
			exp = &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp9(lex)}
		default:
			return exp
		}
	}
}

// exp8 {`|` exp8}
func parseExp9(lex *lexer.Lexer) ast.Exp {
	exp := parseExp8(lex)
	for lex.LookAhead() == lexer.TokenOpBor {
		line, op, _ := lex.NextToken() // `|`
		bOr := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp8(lex)}
		exp = optimizeBitwiseBinOp(bOr)
	}

	return exp
}

// exp7 {`~` exp7} // xor
func parseExp8(lex *lexer.Lexer) ast.Exp {
	exp := parseExp7(lex)
	for lex.LookAhead() == lexer.TokenOpBxor {
		line, op, _ := lex.NextToken() // `~`
		bXor := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp7(lex)}
		exp = optimizeBitwiseBinOp(bXor)
	}

	return exp
}

// exp6 {`&` exp6}
func parseExp7(lex *lexer.Lexer) ast.Exp {
	exp := parseExp6(lex)
	for lex.LookAhead() == lexer.TokenOpBand {
		line, op, _ := lex.NextToken() // `&`
		bAnd := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp6(lex)}
		exp = optimizeBitwiseBinOp(bAnd)
	}

	return exp
}

// exp5 {`<<`|`>>` exp5}
func parseExp6(lex *lexer.Lexer) ast.Exp {
	exp := parseExp5(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TokenOpShl, lexer.TokenOpShr:
			line, op, _ := lex.NextToken() // `<<`|`>>`
			shift := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp5(lex)}
			exp = optimizeBitwiseBinOp(shift)
		default:
			return exp
		}
	}
}

// exp4 { `..` exp4 }
// // NOTE: concatenation operation has right combination
func parseExp5(lex *lexer.Lexer) ast.Exp {
	exp := parseExp4(lex)
	if lex.LookAhead() != lexer.TokenOpConcat {
		return exp
	}

	line := lex.Line()
	exps := make([]ast.Exp, 0, 4)
	exps = append(exps, exp)
	for lex.LookAhead() == lexer.TokenOpConcat {
		line, _, _ = lex.NextToken() // `..`
		exps = append(exps, parseExp4(lex))
	}

	return &ast.ConcatExp{Line: line, Exps: exps}
}

// exp3 {`+`|`-` exp3}
func parseExp4(lex *lexer.Lexer) ast.Exp {
	exp := parseExp3(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TokenOpAdd, lexer.TokenOpSub:
			line, op, _ := lex.NextToken() // +|-
			arith := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp3(lex)}
			exp = optimizeArithBinOp(arith)
		default:
			return exp
		}
	}
}

// exp2 { (`*`|`/`|`//`|`%`) exp2 }
func parseExp3(lex *lexer.Lexer) ast.Exp {
	exp := parseExp2(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TokenOpMul, lexer.TokenOpDiv, lexer.TokenOpIDiv, lexer.TokenOpMod:
			line, op, _ := lex.NextToken() // *|/|//|%
			arith := &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp2(lex)}
			exp = optimizeArithBinOp(arith)
		default:
			return exp
		}
	}
}

// {(`not`|`#`|`-`|`~`)} exp1
// =>
// Recursive
// if has not unary operator then exp1
// else u0 parseExp2
func parseExp2(lex *lexer.Lexer) ast.Exp {
	switch lex.LookAhead() {
	case lexer.TokenOpNot, lexer.TokenOpLen, lexer.TokenOpMinus, lexer.TokenOpBnot:
		line, op, _ := lex.NextToken() // unary op
		exp := &ast.UnOpExp{Line: line, Op: op, MExp: parseExp2(lex)}
		return optimizeUnaryOp(exp)
	}

	return parseExp1(lex)
}

// exp0 {`^` exp0}
// NOTE: power operation is right associative
// `a ^ b ^ c <=> a ^ (b ^ c)`
func parseExp1(lex *lexer.Lexer) ast.Exp {
	exp := parseExp0(lex)
	for lex.LookAhead() == lexer.TokenOpPow {
		line, op, _ := lex.NextToken() // `^`
		exp = &ast.BinOpExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp2(lex)}
		// exp0 ^ exp2
	}

	return optimizePow(exp)
}

// `nil` | `false` | `true` | Numeral | LiteralString |
// `...` | functiondef | prefixexp | tableconstructor
func parseExp0(lex *lexer.Lexer) ast.Exp {
	switch lex.LookAhead() {
	case lexer.TokenKwNil:
		line, _, _ := lex.NextToken() // `nil`
		return &ast.NilExp{Line: line}
	case lexer.TokenKwFalse:
		line, _, _ := lex.NextToken() // `false`
		return &ast.FalseExp{Line: line}
	case lexer.TokenKwTrue:
		line, _, _ := lex.NextToken() // `true`
		return &ast.TrueExp{Line: line}
	case lexer.TokenNumber: // Numeral
		return parseNumberExp(lex)
	case lexer.TokenString:
		line, _, token := lex.NextToken() // LiteralString
		return &ast.StringExp{Line: line, Str: token}
	case lexer.TokenVararg: // `...`
		line, _, _ := lex.NextToken()
		return &ast.VarargExp{Line: line}
	case lexer.TokenKwFunction: // functiondef
		lex.NextToken() // skips the `function`
		return parseFuncDefExp(lex)
	case lexer.TokenSepLcurly: // `{` tablecontructor
		return parseTableConstructor(lex)
	default:
		return parsePrefixExp(lex)
	}
}

func parseNumberExp(lex *lexer.Lexer) ast.Exp {
	line, _, token := lex.NextToken()
	if i, ok := number.ParseInteger(token); ok {
		return &ast.IntegerExp{Line: line, Val: i}
	} else if f, ok := number.ParseFloat(token); ok {
		return &ast.FloatExp{Line: line, Val: f}
	}

	panic("not a number: " + token)
}

// tableconstructor ::= `{` [fieldlist] `}`
// fieldlist ::= field {fieldsep field} [fieldsep]
func parseTableConstructor(lex *lexer.Lexer) *ast.TableConstructionExp {
	line := lex.Line()
	lex.AssertNextTokenKind(lexer.TokenSepLcurly) // `{`
	keyExps, valExps := parseFieldList(lex)       // `[fieldlist]`
	lex.AssertNextTokenKind(lexer.TokenSepRcurly) // `}`
	lastLine := lex.Line()

	return &ast.TableConstructionExp{
		FirstLine: line,
		LastLine:  lastLine,
		KeyExps:   keyExps,
		ValExps:   valExps,
	}
}

func parseFieldList(lex *lexer.Lexer) (ks, vs []ast.Exp) {
	if lex.LookAhead() != lexer.TokenSepRcurly {
		k, v := parseField(lex) // field
		ks = append(ks, k)
		vs = append(vs, v)
		for isFieldSep(lex.LookAhead()) {
			lex.NextToken() // fieldsep
			if lex.LookAhead() != lexer.TokenSepRcurly {
				k, v := parseField(lex) // field
				ks = append(ks, k)
				vs = append(vs, v)
			} else {
				break
			}
		}
	}

	return
}

// fieldsep ::= `,` | `;`
func isFieldSep(kind int) bool {
	return kind == lexer.TokenSepComma || kind == lexer.TokenSepSemi
}

// // field ::= `[` exp `]` `=` exp | Name `=` exp | exp
func parseField(lex *lexer.Lexer) (k, v ast.Exp) {
	if lex.LookAhead() == lexer.TokenSepLBrack { // `[` exp `]` `=` exp
		lex.NextToken()                               // `[`
		k = parseExp(lex)                             // exp
		lex.AssertNextTokenKind(lexer.TokenSepRBrack) // `]`
		lex.AssertNextTokenKind(lexer.TokenOpAssign)  // `=`
		v = parseExp(lex)
		return
	}

	exp := parseExp(lex)
	if name, ok := exp.(*ast.NameExp); ok { // Name `=` exp
		if lex.LookAhead() == lexer.TokenOpAssign { // Name => LiteralString
			lex.NextToken() // `=`
			k = &ast.StringExp{Line: name.Line, Str: name.Name}
			v = parseExp(lex)
			return
		}
	}

	return nil, exp
}

// funcbody ::= `(` [parlist] `)` block end
func parseFuncDefExp(lex *lexer.Lexer) *ast.FuncDefExp {
	line := lex.Line()
	lex.AssertNextTokenKind(lexer.TokenSepLparen)            // `(`
	parList, isVararg := parseParList(lex)                   // [parlist]
	lex.AssertNextTokenKind(lexer.TokenSepRparen)            // `)`
	block := parseBlock(lex)                                 // block
	lastLine, _ := lex.AssertNextTokenKind(lexer.TokenKwEnd) // `end`
	return &ast.FuncDefExp{
		FirstLine: line,
		LastLine:  lastLine,
		ParList:   parList,
		IsVararg:  isVararg,
		MBlock:    block,
	}
}

// parlist ::= namelist [`,` `...`] | `...`
// namelist ::= Name {`,` Name}
func parseParList(lex *lexer.Lexer) (names []string, isVararg bool) {
	switch lex.LookAhead() {
	case lexer.TokenSepRparen:
		return nil, false
	case lexer.TokenVararg:
		lex.NextToken() // `...`
		return nil, true
	}

	_, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
	names = append(names, name)
	for lex.LookAhead() == lexer.TokenSepComma { // `,`
		lex.NextToken() // skips the `,`
		if lex.LookAhead() == lexer.TokenIdentifier {
			_, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
			names = append(names, name)
		} else {
			lex.AssertNextTokenKind(lexer.TokenVararg)
			isVararg = true
			break
		}
	}

	return
}

// prefixexp includes var expression, function call expression
// and parentheses expression
// prefixexp ::= var | functioncall | `(` exp `)`
// var ::=  Name | prefixexp `[` exp `]` | prefixexp `.` Name
// functioncall ::=  prefixexp args | prefixexp `:` Name args
// =>
// prefixexp ::= Name |
//               `(` exp `)`
//               prefixexp `[` exp `]`
//               prefixexp `.` Name
//               prefixexp [`:` Name] args
func parsePrefixExp(lex *lexer.Lexer) (exp ast.Exp) {
	if lex.LookAhead() == lexer.TokenIdentifier { // Name
		line, name := lex.AssertNextTokenKind(lexer.TokenIdentifier)
		exp = &ast.NameExp{Line: line, Name: name}
	} else { // `(` exp `)`
		exp = parseParensExp(lex)
	}
	return finishPrefixExp(lex, exp)
}

func parseParensExp(lex *lexer.Lexer) ast.Exp {
	lex.AssertNextTokenKind(lexer.TokenSepLparen) // `(`
	exp := parseExp(lex)                          // exp
	lex.AssertNextTokenKind(lexer.TokenSepRparen) // `)`

	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp, *ast.NameExp, *ast.TableAccessExp:
		return &ast.ParensExp{MExp: exp}
	}

	// no need to keep the parens?
	return exp
}

func finishPrefixExp(lex *lexer.Lexer, exp ast.Exp) ast.Exp {
	for {
		switch lex.LookAhead() {
		case lexer.TokenSepLBrack: // prefixexp `[` exp `]`
			lex.NextToken()                               // `[`
			keyExp := parseExp(lex)                       // exp
			lex.AssertNextTokenKind(lexer.TokenSepRBrack) // `]`
			exp = &ast.TableAccessExp{LastLine: lex.Line(), PrefixExp: exp, Key: keyExp}

		case lexer.TokenSepDot: // prefixexp `.` Name
			lex.NextToken()                                              // `.`
			line, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
			keyExp := &ast.StringExp{Line: line, Str: name}
			exp = &ast.TableAccessExp{LastLine: line, PrefixExp: exp, Key: keyExp}

		case lexer.TokenSepColon, lexer.TokenSepLparen,
			lexer.TokenSepLcurly, lexer.TokenString:
			exp = parseFuncCallExp(lex, exp)
		default:
			return exp
		}
	}
}

// functioncall ::=  prefixexp args | prefixexp `:` Name args
// https://cloudwu.github.io/lua53doc/manual.html#3.4.10
// f(a)
// t:f(a)    => f(t, a)
// f{fields} => f({fields})
// f"string" => f("string") // LiteralString
func parseFuncCallExp(lex *lexer.Lexer, exp ast.Exp) ast.Exp {
	nameExp := parseNameExp(lex)
	line := lex.Line()
	args := parseArgs(lex)
	lastLine := lex.Line()
	return &ast.FuncCallExp{
		FirstLine: line,
		LastLine:  lastLine,
		PrefixExp: exp,
		FNameExp:  nameExp,
		Args:      args,
	}
}

func parseNameExp(lex *lexer.Lexer) *ast.StringExp {
	if lex.LookAhead() == lexer.TokenSepColon {
		lex.NextToken() // `:`
		line, name := lex.AssertNextTokenKind(lexer.TokenIdentifier)
		return &ast.StringExp{Line: line, Str: name}
	}
	return nil
}

// args ::=  `(` [explist] `)` | tableconstructor | LiteralString
func parseArgs(lex *lexer.Lexer) (args []ast.Exp) {
	switch lex.LookAhead() {
	case lexer.TokenSepLparen: // `(` [explist] `)`
		lex.NextToken() // `(`
		if lex.LookAhead() != lexer.TokenSepRparen {
			args = parseExpList(lex)
		}
		lex.AssertNextTokenKind(lexer.TokenSepRparen)
	case lexer.TokenSepLcurly: // tableconstructor
		args = []ast.Exp{parseTableConstructor(lex)}
	case lexer.TokenString:
		line, _, str := lex.NextToken()
		args = []ast.Exp{&ast.StringExp{Line: line, Str: str}}
	}
	return
}
