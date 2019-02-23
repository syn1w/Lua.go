package ast

// stat ::=  `;` |
// 	varlist `=` explist |
// 	functioncall |
// 	label |
// 	break |
// 	goto Name |
// 	do block end |
// 	while exp do block end |
// 	repeat block until exp |
// 	if exp then block {elseif exp then block} [else block] end |
// 	for Name `=` exp `,` exp [`,` exp] do block end |
// 	for namelist in explist do block end |
// 	function funcname funcbody |
// 	local function Name funcbody |
// 	local namelist [`,` explist]

// Stat is statement interface
type Stat interface{}

// EmptyStat is `;` statement
type EmptyStat struct{}

// AssignStat is `varlist `=` explist` statement
type AssignStat struct {
	LastLine int
	VarList  []Exp
	ExpList  []Exp
}

// FuncCallStat is `func()` statement
type FuncCallStat = FuncCallExp

// LabelStat is `::` statement
type LabelStat struct {
	Line int
	Name string
}

// BreakStat is `break` statement
type BreakStat struct {
	Line int
}

// GotoStat is `goto name` statement
type GotoStat struct {
	Line int
	Name string
}

// DoStat is `do mblock end` statement block
type DoStat struct {
	MBlock *Block
}

// WhileStat is `while bexp do mblock end` statement
type WhileStat struct {
	BExp   Exp
	MBlock *Block
}

// RepeatStat is `repeat block until exp` statement
type RepeatStat struct {
	BExp   Exp
	MBlock *Block
}

// IfStat is `if exp then block {elseif exp then block} [else block] end`
type IfStat struct {
	BExps  []Exp
	Blocks []*Block
}

// ForNumStat is `for Name `=` exp `,` exp [`,` exp] do block end` statement
type ForNumStat struct {
	LineFor  int
	LineDo   int
	VarName  string
	InitExp  Exp
	LimitExp Exp
	StepExp  Exp
	MBlock   *Block
}

// ForInStat is `for namelist in explist do block end` statement
type ForInStat struct {
	LineDo   int
	NameList []string
	ExpList  []Exp
	MBlock   *Block
}

// LocalFuncDefStat is `local function funcname funcbody` statement
type LocalFuncDefStat struct {
	Name string
	Func *FuncDefExp
}

// LocalVarDeclStat is `local namelist [`=` explist]` statement
type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}
