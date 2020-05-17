package parser

import (
	"luago/compiler/ast"
	"luago/compiler/lexer"
	"luago/number"
	"math"
)

func optimizeLogicalOr(exp *ast.BinOpExp) ast.Exp {
	if isTrue(exp.Exp1) {
		return exp.Exp1 // true or x => true
	}
	if isFalse(exp.Exp1) && !isVarargOrFuncCall(exp.Exp2) {
		return exp.Exp2 // false or x => x
	}

	return exp
}

func optimizeLogicalAnd(exp *ast.BinOpExp) ast.Exp {
	if isFalse(exp.Exp1) {
		return exp.Exp1 // false and x => false
	}
	if isTrue(exp.Exp1) && !isVarargOrFuncCall(exp.Exp2) {
		return exp.Exp2 // true and x => x
	}

	return exp
}

func optimizeBitwiseBinOp(exp *ast.BinOpExp) ast.Exp {
	if x, ok := castToInt(exp.Exp1); ok {
		if y, ok := castToInt(exp.Exp2); ok {
			switch exp.Op {
			case lexer.TokenOpBand:
				return &ast.IntegerExp{Line: exp.Line, Val: x & y}
			case lexer.TokenOpBor:
				return &ast.IntegerExp{Line: exp.Line, Val: x | y}
			case lexer.TokenOpBxor:
				return &ast.IntegerExp{Line: exp.Line, Val: x ^ y}
			case lexer.TokenOpShl:
				return &ast.IntegerExp{Line: exp.Line, Val: number.ShiftLeft(x, y)}
			case lexer.TokenOpShr:
				return &ast.IntegerExp{Line: exp.Line, Val: number.ShiftRight(x, y)}
			}
		}
	}

	return exp
}

func optimizeArithBinOp(exp *ast.BinOpExp) ast.Exp {
	if x, ok := exp.Exp1.(*ast.IntegerExp); ok {
		if y, ok := exp.Exp2.(*ast.IntegerExp); ok {
			switch exp.Op {
			case lexer.TokenOpAdd:
				return &ast.IntegerExp{Line: exp.Line, Val: x.Val + y.Val}
			case lexer.TokenOpSub:
				return &ast.IntegerExp{Line: exp.Line, Val: x.Val - y.Val}
			case lexer.TokenOpMul:
				return &ast.IntegerExp{Line: exp.Line, Val: x.Val * y.Val}
			case lexer.TokenOpIDiv:
				if y.Val != 0 {
					return &ast.IntegerExp{
						Line: exp.Line,
						Val:  number.IFloorDiv(x.Val, y.Val),
					}
				}
			case lexer.TokenOpMod:
				if y.Val != 0 {
					return &ast.IntegerExp{
						Line: exp.Line,
						Val:  number.IMod(x.Val, y.Val),
					}
				}
			}
		}
	}

	if x, ok := castToFloat(exp.Exp1); ok {
		if y, ok := castToFloat(exp.Exp2); ok {
			switch exp.Op {
			case lexer.TokenOpAdd:
				return &ast.FloatExp{Line: exp.Line, Val: x + y}
			case lexer.TokenOpSub:
				return &ast.FloatExp{Line: exp.Line, Val: x - y}
			case lexer.TokenOpMul:
				return &ast.FloatExp{Line: exp.Line, Val: x * y}
			case lexer.TokenOpDiv:
				if y != 0 {
					return &ast.FloatExp{Line: exp.Line, Val: x / y}
				}
			case lexer.TokenOpIDiv:
				if y != 0 {
					return &ast.FloatExp{Line: exp.Line, Val: number.FFloorDiv(x, y)}
				}
			case lexer.TokenOpMod:
				return &ast.FloatExp{Line: exp.Line, Val: number.FMod(x, y)}
			case lexer.TokenOpPow:
				return &ast.FloatExp{Line: exp.Line, Val: math.Pow(x, y)}
			}
		}
	}

	return exp
}

func optimizeUnaryOp(exp *ast.UnOpExp) ast.Exp {
	switch exp.Op {
	case lexer.TokenOpNot:
		return optimizeNot(exp)
	case lexer.TokenOpUnm:
		return optimizeUnm(exp)
	case lexer.TokenOpBnot:
		return optimizeBnot(exp)
		// NOTE: len?
	default:
		return exp
	}
}

func optimizeNot(exp *ast.UnOpExp) ast.Exp {
	switch exp.MExp.(type) {
	case *ast.NilExp, *ast.FalseExp:
		return &ast.TrueExp{Line: exp.Line}
	case *ast.TrueExp, *ast.IntegerExp, *ast.FloatExp, *ast.StringExp:
		return &ast.FalseExp{Line: exp.Line}
	default:
		return exp
	}
}

func optimizeUnm(exp *ast.UnOpExp) ast.Exp {
	switch x := exp.MExp.(type) {
	case *ast.IntegerExp:
		x.Val = -x.Val
		return x
	case *ast.FloatExp:
		x.Val = -x.Val
		return x
	default:
		return exp
	}
}

func optimizeBnot(exp *ast.UnOpExp) ast.Exp {
	switch x := exp.MExp.(type) {
	case *ast.IntegerExp:
		x.Val = ^x.Val // ^ is bitwise not in golang
		return x
	case *ast.FloatExp:
		if i, ok := number.FloatToInteger(x.Val); ok {
			return &ast.IntegerExp{Line: x.Line, Val: ^i}
		}
	}

	return exp
}

// exp = exp0 ^ exp2 | exp0
func optimizePow(exp ast.Exp) ast.Exp {
	if binOp, ok := exp.(*ast.BinOpExp); ok { // exp0 ^ exp2
		if binOp.Op == lexer.TokenOpPow {
			binOp.Exp2 = optimizePow(binOp.Exp2)
		}

		return optimizeArithBinOp(binOp)
	}

	return exp // exp0
}

// nil and false => false
// other => true
func isTrue(exp ast.Exp) bool {
	switch exp.(type) {
	case *ast.TrueExp, *ast.IntegerExp, *ast.FloatExp, *ast.StringExp:
		return true
	default:
		return false
	}
}

func isFalse(exp ast.Exp) bool {
	switch exp.(type) {
	case *ast.NilExp, *ast.FalseExp:
		return true
	default:
		return false
	}
}

func isVarargOrFuncCall(exp ast.Exp) bool {
	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp:
		return true
	default:
		return false
	}
}

func castToInt(exp ast.Exp) (int64, bool) {
	switch x := exp.(type) {
	case *ast.IntegerExp:
		return x.Val, true
	case *ast.FloatExp:
		return number.FloatToInteger(x.Val)
	default:
		return 0, false
	}
}

func castToFloat(exp ast.Exp) (float64, bool) {
	switch x := exp.(type) {
	case *ast.IntegerExp:
		return float64(x.Val), true
	case *ast.FloatExp:
		return x.Val, true
	default:
		return 0.0, false
	}
}
