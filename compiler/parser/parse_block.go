package parser

import (
	"vczn/luago/compiler/ast"
	"vczn/luago/compiler/lexer"
)

func parseBlock(lex *lexer.Lexer) *ast.Block {
	return &ast.Block{
		Stats:    parseStats(lex),
		RetExps:  parseRetExps(lex),
		LastLine: lex.Line(),
	}
}

func parseStats(lex *lexer.Lexer) []ast.Stat {
	stats := make([]ast.Stat, 0, 8)
	for !isReturnOrBlockEnd(lex.LookAhead()) {
		stat := parseStat(lex)
		if _, ok := stat.(*ast.EmptyStat); !ok { // ignore the empty stat
			stats = append(stats, stat)
		}
	}

	return stats
}

// see https://cloudwu.github.io/lua53doc/manual.html#9
// after the block there are `return`, eof, `end`, `until`, `elseif`, `else`
func isReturnOrBlockEnd(kind int) bool {
	switch kind {
	case lexer.TokenKwReturn, lexer.TokenEOF, lexer.TokenKwEnd, lexer.TokenKwUntil,
		lexer.TokenKwElseif, lexer.TokenKwElse:
		return true
	}
	return false
}

func parseRetExps(lex *lexer.Lexer) []ast.Exp {
	if lex.LookAhead() != lexer.TokenKwReturn {
		return nil
	}

	lex.NextToken() // skips the `return`
	switch lex.LookAhead() {
	case lexer.TokenEOF, lexer.TokenKwEnd, lexer.TokenKwElseif,
		lexer.TokenKwElse, lexer.TokenKwUntil: // return
		return []ast.Exp{}
	case lexer.TokenSepSemi: // return;
		lex.NextToken() // skips the `;`
		return []ast.Exp{}
	default:
		exps := parseExpList(lex)
		if lex.LookAhead() == lexer.TokenSepSemi { // return exps;
			lex.NextToken() // skips the `;`
		}
		return exps
	}
}
