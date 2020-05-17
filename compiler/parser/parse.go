package parser

import (
	"luago/compiler/ast"
	"luago/compiler/lexer"
)

// Parse lua string to lua chunk
func Parse(chunk, chunkName string) *ast.Block {
	lex := lexer.NewLexer(chunk, chunkName)
	block := parseBlock(lex)
	lex.AssertNextTokenKind(lexer.TokenEOF)
	return block
}
