package parser

import (
	"vczn/luago/compiler/ast"
	"vczn/luago/compiler/lexer"
)

func parseStat(lex *lexer.Lexer) ast.Stat {
	switch lex.LookAhead() {
	case lexer.TokenSepSemi:
		return parseEmptyStat(lex)
	case lexer.TokenSepLabel:
		return parseLabelStat(lex)
	case lexer.TokenKwBreak:
		return parseBreakStat(lex)
	case lexer.TokenKwGoto:
		return parseGotoStat(lex)
	case lexer.TokenKwDo:
		return parseDoStat(lex)
	case lexer.TokenKwWhile:
		return parseWhileStat(lex)
	case lexer.TokenKwRepeat:
		return parseRepeatStat(lex)
	case lexer.TokenKwIf:
		return parseIfStat(lex)
	case lexer.TokenKwFor:
		return parseForStat(lex)
	case lexer.TokenKwFunction:
		return parseFuncDefStat(lex)
	case lexer.TokenKwLocal:
		return parseLocalAssignOrFuncDefStat(lex)
	default:
		return parseAssignOrFuncCallStat(lex)
	}
}

// `;`
func parseEmptyStat(lex *lexer.Lexer) *ast.EmptyStat {
	lex.AssertNextTokenKind(lexer.TokenSepSemi) // `;`
	return &ast.EmptyStat{}
}

// :: Name ::
func parseLabelStat(lex *lexer.Lexer) *ast.LabelStat {
	lex.AssertNextTokenKind(lexer.TokenSepLabel)                 // `::`
	line, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
	lex.AssertNextTokenKind(lexer.TokenSepLabel)                 // `::`
	return &ast.LabelStat{Line: line, Name: name}
}

// break
func parseBreakStat(lex *lexer.Lexer) *ast.BreakStat {
	lex.AssertNextTokenKind(lexer.TokenKwBreak) // `break`
	return &ast.BreakStat{Line: lex.Line()}
}

// goto Name
func parseGotoStat(lex *lexer.Lexer) *ast.GotoStat {
	lex.AssertNextTokenKind(lexer.TokenKwGoto)                   // `goto`
	line, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
	return &ast.GotoStat{Line: line, Name: name}
}

// do block end
func parseDoStat(lex *lexer.Lexer) *ast.DoStat {
	lex.AssertNextTokenKind(lexer.TokenKwDo)  // `do`
	block := parseBlock(lex)                  // block
	lex.AssertNextTokenKind(lexer.TokenKwEnd) // `end`
	return &ast.DoStat{MBlock: block}
}

// while exp do block end
func parseWhileStat(lex *lexer.Lexer) *ast.WhileStat {
	lex.AssertNextTokenKind(lexer.TokenKwWhile) // `while`
	exp := parseExp(lex)                        // exp
	lex.AssertNextTokenKind(lexer.TokenKwDo)    // `do`
	block := parseBlock(lex)                    // block
	lex.AssertNextTokenKind(lexer.TokenKwEnd)   // `end`
	return &ast.WhileStat{
		BExp:   exp,
		MBlock: block,
	}
}

// repeat block until exp
func parseRepeatStat(lex *lexer.Lexer) *ast.RepeatStat {
	lex.AssertNextTokenKind(lexer.TokenKwRepeat) // `repeat`
	block := parseBlock(lex)                     // block
	lex.AssertNextTokenKind(lexer.TokenKwUntil)  // `until`
	exp := parseExp(lex)                         // exp
	return &ast.RepeatStat{
		BExp:   exp,
		MBlock: block,
	}
}

// if exp then block {elseif exp then block} [else block] end
func parseIfStat(lex *lexer.Lexer) *ast.IfStat {
	exps := make([]ast.Exp, 0, 4)
	blocks := make([]*ast.Block, 0, 4)
	lex.AssertNextTokenKind(lexer.TokenKwIf)   // `if`
	exps = append(exps, parseExp(lex))         // exp
	lex.AssertNextTokenKind(lexer.TokenKwThen) // `then`
	blocks = append(blocks, parseBlock(lex))   // block

	// {elseif exp then block}
	for lex.LookAhead() == lexer.TokenKwElseif {
		lex.NextToken()                            // `elseif`
		exps = append(exps, parseExp(lex))         // exp
		lex.AssertNextTokenKind(lexer.TokenKwThen) // `then`
		blocks = append(blocks, parseBlock(lex))   // block
	}

	// [else block]
	if lex.LookAhead() == lexer.TokenKwElse {
		lex.NextToken()                                     // `else`
		exps = append(exps, &ast.TrueExp{Line: lex.Line()}) // -> true exp
		blocks = append(blocks, parseBlock(lex))            // block
	}

	lex.AssertNextTokenKind(lexer.TokenKwEnd)
	return &ast.IfStat{
		BExps:  exps,
		Blocks: blocks,
	}
}

func parseForStat(lex *lexer.Lexer) ast.Stat {
	lineFor := lex.Line()
	lex.AssertNextTokenKind(lexer.TokenKwFor)                 // `for`
	_, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // name1
	if lex.LookAhead() == lexer.TokenOpAssign {               // for number
		return parseForNumStat(lex, lineFor, name)
	}
	return parseForInStat(lex, name)
}

// for Name = exp , exp [ , exp] do block end
func parseForNumStat(lex *lexer.Lexer, lineFor int, name string) *ast.ForNumStat {
	lex.AssertNextTokenKind(lexer.TokenOpAssign) // `=`
	initExp := parseExp(lex)                     // init exp
	lex.AssertNextTokenKind(lexer.TokenSepComma) // `,`
	limitExp := parseExp(lex)                    // limit exp

	var stepExp ast.Exp
	if lex.LookAhead() == lexer.TokenSepComma { // [, step exp]
		lex.NextToken()         // `,`
		stepExp = parseExp(lex) // step exp
	} else {
		stepExp = &ast.IntegerExp{Line: lex.Line(), Val: 1}
	}

	lineDo, _ := lex.AssertNextTokenKind(lexer.TokenKwDo)
	block := parseBlock(lex)
	lex.AssertNextTokenKind(lexer.TokenKwEnd)

	return &ast.ForNumStat{
		LineFor:  lineFor,
		LineDo:   lineDo,
		VarName:  name,
		InitExp:  initExp,
		LimitExp: limitExp,
		StepExp:  stepExp,
		MBlock:   block,
	}
}

// for namelist in explist do block end
func parseForInStat(lex *lexer.Lexer, name1 string) *ast.ForInStat {
	nameList := parseRestNameList(lex, name1)

	lex.AssertNextTokenKind(lexer.TokenKwIn)              // `in`
	expList := parseExpList(lex)                          // explist
	lineDo, _ := lex.AssertNextTokenKind(lexer.TokenKwDo) // `do`
	block := parseBlock(lex)                              // block
	lex.AssertNextTokenKind(lexer.TokenKwEnd)             // `end`

	return &ast.ForInStat{
		LineDo:   lineDo,
		NameList: nameList,
		ExpList:  expList,
		MBlock:   block,
	}
}

// `function` funcname funcbody
// funcbody ::= `(` [parlist] `)` block end
// `function f () body end` => `f = function() body end`
// `function t.a.b.c.f () body end` => `t.a.b.c.f = function () body end`
// `function t.a.b.c:f (params) body end` => `t.a.b.c.f = function (self, params) body end`
func parseFuncDefStat(lex *lexer.Lexer) *ast.AssignStat {
	lex.AssertNextTokenKind(lexer.TokenKwFunction) // `function`
	funcExp, hasColon := parseFuncName(lex)        // funcname
	funcDef := parseFuncDefExp(lex)                // funcbody

	if hasColon { // v:fn(args) => v.fn(self, args)
		selfParam := []string{"self"}
		funcDef.ParList = append(selfParam, funcDef.ParList...)
	}

	return &ast.AssignStat{
		LastLine: funcDef.FirstLine, // ? LastLine
		VarList:  []ast.Exp{funcExp},
		ExpList:  []ast.Exp{funcDef},
	}
}

// funcname ::= Name {`.` Name} [`:` Name]
func parseFuncName(lex *lexer.Lexer) (exp ast.Exp, hasColon bool) {
	line, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // name
	exp = &ast.NameExp{Line: line, Name: name}
	hasColon = false

	for lex.LookAhead() == lexer.TokenSepDot { // { `.` Name }
		lex.NextToken()                                             // `.`
		line, name = lex.AssertNextTokenKind(lexer.TokenIdentifier) // name
		key := &ast.StringExp{Line: line, Str: name}
		exp = &ast.TableAccessExp{LastLine: line, PrefixExp: exp, Key: key}
	}

	if lex.LookAhead() == lexer.TokenSepColon { // [ `:` Name ]
		lex.NextToken()                                             // `:`
		line, name = lex.AssertNextTokenKind(lexer.TokenIdentifier) // name
		key := &ast.StringExp{Line: line, Str: name}
		exp = &ast.TableAccessExp{LastLine: line, PrefixExp: exp, Key: key}
		hasColon = true
	}

	return
}

func parseLocalAssignOrFuncDefStat(lex *lexer.Lexer) ast.Stat {
	lex.AssertNextTokenKind(lexer.TokenKwLocal) // `local`
	if lex.LookAhead() == lexer.TokenKwFunction {
		return parseLocalFuncDefStat(lex)
	}
	return parseLocalAssignStat(lex)
}

// `local` | `function` Name funcbody
// `local function f() end`  =>  `local f; f = function() end`
func parseLocalFuncDefStat(lex *lexer.Lexer) *ast.LocalFuncDefStat {
	lex.AssertNextTokenKind(lexer.TokenKwFunction)            // `function`
	_, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
	funcDefExp := parseFuncDefExp(lex)                        // funcbody
	return &ast.LocalFuncDefStat{
		Name: name,
		Func: funcDefExp,
	}
}

// `local` | namelist [ `=` explist]
func parseLocalAssignStat(lex *lexer.Lexer) *ast.LocalVarDeclStat {
	_, name1 := lex.AssertNextTokenKind(lexer.TokenIdentifier) // Name
	nameList := parseRestNameList(lex, name1)                  // `{ , Name }`

	var expList []ast.Exp
	if lex.LookAhead() == lexer.TokenOpAssign { // [ `=` explist]
		lex.NextToken() // `=`
		expList = parseExpList(lex)
	}

	lastLine := lex.Line()
	return &ast.LocalVarDeclStat{
		LastLine: lastLine,
		NameList: nameList,
		ExpList:  expList,
	}
}

// varlist = explist | functioncall
func parseAssignOrFuncCallStat(lex *lexer.Lexer) ast.Stat {
	prefixExp := parsePrefixExp(lex)
	if funcCall, ok := prefixExp.(*ast.FuncCallExp); ok {
		return funcCall
	}
	return parseAssignStat(lex, prefixExp)
}

// varlist = explist
func parseAssignStat(lex *lexer.Lexer, prefixExp ast.Exp) *ast.AssignStat {
	varList := parseVarList(lex, prefixExp)
	lex.AssertNextTokenKind(lexer.TokenOpAssign) // `=`
	expList := parseExpList(lex)

	lastLine := lex.Line()
	return &ast.AssignStat{
		LastLine: lastLine,
		VarList:  varList,
		ExpList:  expList,
	}
}

// varlist ::= var {`,` var}
// var ::=  Name | prefixexp `[` exp `]` | prefixexp `.` Name
func parseVarList(lex *lexer.Lexer, prefixExp ast.Exp) []ast.Exp {
	varlist := make([]ast.Exp, 0, 4)
	varlist = append(varlist, checkVar(lex, prefixExp))
	for lex.LookAhead() == lexer.TokenSepComma {
		lex.NextToken() // `,`
		exp := parsePrefixExp(lex)
		varlist = append(varlist, checkVar(lex, exp))
	}

	return varlist
}

func checkVar(lex *lexer.Lexer, exp ast.Exp) ast.Exp {
	switch exp.(type) {
	case *ast.NameExp, *ast.TableAccessExp:
		return exp
	}
	lex.AssertNextTokenKind(-1) // assert error
	panic("unreachable")
}

// parse name list apart from the name1
// namelist ::= Name {`,` Name}
func parseRestNameList(lex *lexer.Lexer, name1 string) []string {
	nameList := make([]string, 0, 4)
	nameList = append(nameList, name1)
	for lex.LookAhead() == lexer.TokenSepComma { // { , Name}
		lex.NextToken()                                           // `,`
		_, name := lex.AssertNextTokenKind(lexer.TokenIdentifier) // `name`
		nameList = append(nameList, name)
	}

	return nameList
}
