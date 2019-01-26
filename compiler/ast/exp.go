package ast

// exp ::=  `nil` | `false` | `true` | Numeral | LiteralString | `...` | functiondef |
//          prefixexp | tableconstructor | exp binop exp | unop exp

// Exp is expression interface
type Exp interface{}

// NilExp is `nil` expression
type NilExp struct {
	Line int
}

// TrueExp is `true` expression
type TrueExp struct {
	Line int
}

// FalseExp is `false` expression
type FalseExp struct {
	Line int
}

// IntegerExp is integer expression
type IntegerExp struct {
	Line int
	Val  int64
}

// FloatExp is floating point expression
type FloatExp struct {
	Line int
	Val  float64
}

// StringExp is string expression
type StringExp struct {
	Line int
	Str  string
}

// VarargExp is `...` expression
type VarargExp struct {
	Line int
}

// NameExp is identifier name expression
type NameExp struct {
	Line int
	Name string
}

// TableConstructionExp is table construction expression
// tableconstructor ::= `{` [fieldlist] `}`
// fieldlist ::= field {fieldsep field} [fieldsep]
// field ::= `[` exp `]` `=` exp | Name `=` exp | exp
// fieldsep ::= `,` | `;`
type TableConstructionExp struct {
	FirstLine int // line of '{'
	LastLine  int // line of '}'
	KeyExps   []Exp
	ValExps   []Exp
}

// FuncDefExp is function define expression
// functiondef ::= `function` funcbody
// funcbody ::= `(` [parlist] `)` block end
// parlist ::= namelist [`,` `...`] | `...`
// namelist ::= Name {`,` Name}
type FuncDefExp struct {
	FirstLine int
	LastLine  int // line of `end`
	ParList   []string
	IsVararg  bool
	MBlock    *Block
}

// prefixexp includes var expression, function call expression
// and parentheses expression
// prefixexp ::= var | functioncall | `(` exp `)`
// var ::=  Name | prefixexp `[` exp `]` | prefixexp `.` Name
// functioncall ::=  prefixexp args | prefixexp `:` Name args
// =>
// prefixexp ::= Name |
//               `(` exp `)`
//               prefixexp `(` exp `)`
//               prefixexp `.` Name
//               prefixexp [`:` Name] args

// TableAccessExp is table access expression
type TableAccessExp struct {
	LastLint  int
	PrefixExp Exp
	Key       Exp
}

// FuncCallExp is functioncall expression
type FuncCallExp struct {
	FirstLine int
	LastLine  int
	PrefixExp Exp
	NameExp   Exp
	Args      []Exp
}

// ParensExp is parentheses expression
type ParensExp struct {
	MExp Exp
}

// UnOpExp is unary expression
// unop ::= `-` | `not` | `#` | `~`
type UnOpExp struct {
	Line int
	Op   int
	MExp Exp
}

// BinOpExp is binary expression
// binop ::=  `+` | `-` | `*` | `/` | `//` | `^` | `%` |
//            `&` | `~` | `|` | `>>` | `<<` | `..` |
//            `<` | `<=` | `>` | `>=` | `==` | `~=` |
//            `and` | `or`
type BinOpExp struct {
	Line int
	Op   int
	Exp1 Exp
	Exp2 Exp
}

// ConcatExp is `..` expression, for optimizing the concatenating operation
type ConcatExp struct {
	Line int
	Exps []Exp
}
