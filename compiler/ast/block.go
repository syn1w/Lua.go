package ast

// chun ::= block

// Block is code block
// block ::= {stat} [retstat]
// retstat ::= return [explist] [';']
// explist ::= exp {',' exp}
type Block struct {
	LastLine int
	Stats    []Stat
	RetExps  []Exp
}
