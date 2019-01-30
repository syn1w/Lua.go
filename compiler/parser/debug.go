package parser

import (
	"fmt"
	"vczn/luago/compiler/ast"
	"vczn/luago/compiler/lexer"
)

func expToString(exp ast.Exp) string {
	switch x := exp.(type) {
	case *ast.NilExp:
		return "nil"
	case *ast.TrueExp:
		return "true"
	case *ast.FalseExp:
		return "false"
	case *ast.VarargExp:
		return "..."
	case *ast.IntegerExp:
		return fmt.Sprintf("%d", x.Val)
	case *ast.FloatExp:
		return fmt.Sprintf("%f", x.Val)
	case *ast.StringExp:
		return "'" + x.Str + "'"
	case *ast.ConcatExp:
		var str string
		for _, exp := range x.Exps {
			str += " .. "
			str += expToString(exp)
		}
		return str[4:]
	case *ast.UnOpExp:
		return unopToString(x.Op) + "(" + expToString(x.MExp) + ")"
	case *ast.BinOpExp:
		return "(" + expToString(x.Exp1) + binopToString(x.Op) + expToString(x.Exp2) + ")"
	case *ast.TableConstructionExp:
		return tcExpToString(x)
	case *ast.FuncDefExp:
		return fdExpToString(x, "")
	case *ast.NameExp:
		return x.Name
	case *ast.ParensExp:
		return "(" + expToString(x.MExp) + ")"
	case *ast.TableAccessExp:
		return expToString(x.PrefixExp) + "[" + expToString(x.Key) + "]"
	case *ast.FuncCallExp:
		return fcExpToString(x)
	case int: // index?
		return fmt.Sprintf("%d", x)
	default:
		panic("unreachable")
	}
}

func unopToString(op int) string {
	switch op {
	case lexer.TokenOpUnm:
		return "-"
	case lexer.TokenOpBnot:
		return "~"
	case lexer.TokenOpNot:
		return "not"
	case lexer.TokenOpLen:
		return "#"
	default:
		panic("unreachable")
	}
}

func binopToString(op int) string {
	switch op {
	case lexer.TokenOpAdd:
		return " + "
	case lexer.TokenOpSub:
		return " - "
	case lexer.TokenOpMul:
		return " * "
	case lexer.TokenOpDiv:
		return " / "
	case lexer.TokenOpIDiv:
		return " // "
	case lexer.TokenOpMod:
		return " % "
	case lexer.TokenOpPow:
		return " ^ "
	case lexer.TokenOpBand:
		return " & "
	case lexer.TokenOpBor:
		return " | "
	case lexer.TokenOpBxor:
		return " ~ "
	case lexer.TokenOpShl:
		return " << "
	case lexer.TokenOpShr:
		return " >> "
	case lexer.TokenOpEq:
		return " == "
	case lexer.TokenOpNe:
		return " ~= "
	case lexer.TokenOpLt:
		return " < "
	case lexer.TokenOpLe:
		return " <= "
	case lexer.TokenOpGt:
		return " > "
	case lexer.TokenOpGe:
		return " >= "
	case lexer.TokenOpAnd:
		return " and "
	case lexer.TokenOpOr:
		return " or "
	case lexer.TokenOpConcat:
		return " .. "
	default:
		panic("unreachable")
	}
}

func tcExpToString(exp *ast.TableConstructionExp) string {
	str := "{"
	for i, k := range exp.KeyExps {
		v := exp.ValExps[i]
		if k != nil {
			str += "[" + expToString(k) + "]="
		}
		str += expToString(v) + ","
	}
	str += "}"
	return str
}

func fdExpToString(exp *ast.FuncDefExp, name string) string {
	str := "function"
	if name != "" {
		str += " " + name
	}

	str += "("
	for i, param := range exp.ParList {
		str += param
		if i < len(exp.ParList)-1 {
			str += ", "
		}
	}

	if exp.IsVararg {
		if len(exp.ParList) > 0 {
			str += ", ..."
		} else {
			str += "..."
		}
	}

	str += ") end"
	return str
}

func fcExpToString(exp *ast.FuncCallExp) string {
	str := expToString(exp.PrefixExp)
	if exp.FNameExp != nil {
		str += ":" + exp.FNameExp.Str
	}

	str += "("
	for i, arg := range exp.Args {
		str += expToString(arg)
		if i < len(exp.Args)-1 {
			str += ", "
		}
	}
	str += ")"
	return str
}
