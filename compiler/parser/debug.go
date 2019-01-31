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

	str += ") "
	str += blockToString(exp.MBlock) + " end"

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

func statToString(stat ast.Stat) string {
	switch x := stat.(type) {
	case *ast.EmptyStat:
		return ";"
	case *ast.BreakStat:
		return "break"
	case *ast.LabelStat:
		return "::" + x.Name + "::"
	case *ast.GotoStat:
		return "goto " + x.Name
	case *ast.DoStat:
		return "do " + blockToString(x.MBlock) + " end"
	case *ast.WhileStat:
		return "while " + expToString(x.BExp) + " do " +
			blockToString(x.MBlock) + " end"
	case *ast.RepeatStat:
		return "repeat " + blockToString(x.MBlock) +
			" until " + expToString(x.BExp)
	case *ast.FuncCallStat:
		return fcExpToString(x)
	case *ast.IfStat:
		return ifStatToString(x)
	case *ast.ForNumStat:
		return forNumStatToString(x)
	case *ast.ForInStat:
		return forInStatToString(x)
	case *ast.AssignStat:
		return assignStatToString(x)
	case *ast.LocalVarDeclStat:
		return localVarDefStatToString(x)
	case *ast.LocalFuncDefStat:
		return "local " + fdExpToString(x.Func, x.Name)
	}

	panic("todo")
}

// {stat} [retstat]
func blockToString(block *ast.Block) string {
	str := ""
	if len(block.Stats) > 0 {
		for i, stat := range block.Stats {
			str += statToString(stat)
			if i < len(block.Stats)-1 {
				str += " "
			}
		}
	}

	if block.RetExps != nil {
		if len(block.Stats) == 0 {
			str += "return"
		} else {
			str += " return"
		}

		for _, exp := range block.RetExps {
			str += " " + expToString(exp)
		}
	}

	return str
}

// if exp then block {elseif exp then block} [else block] end
// => else -> elseif true
func ifStatToString(stat *ast.IfStat) string {
	str := "if " + expToString(stat.BExps[0]) +
		" then " + blockToString(stat.Blocks[0])
	for i := 1; i < len(stat.BExps); i++ {
		str += " elseif " + expToString(stat.BExps[i])
		str += " then " + blockToString(stat.Blocks[i])
	}

	str += " end"
	return str
}

// for v = e1, e2, e3 do block end
func forNumStatToString(stat *ast.ForNumStat) string {
	str := "for " + stat.VarName + " = " + expToString(stat.InitExp) +
		", " + expToString(stat.LimitExp)
	if stat.StepExp != nil {
		str += ", " + expToString(stat.StepExp)
	}
	str += " do " + blockToString(stat.MBlock) + " end"
	return str
}

// for namelist in explist do block end
// namelist ::= Name {`,` Name}
func forInStatToString(stat *ast.ForInStat) string {
	str := "for "
	for i, name := range stat.NameList {
		str += name
		if i < len(stat.NameList)-1 {
			str += ", "
		}
	}

	str += " in "
	for i, exp := range stat.ExpList {
		str += expToString(exp)
		if i < len(stat.ExpList)-1 {
			str += ", "
		}
	}
	str += " do " + blockToString(stat.MBlock) + " end"

	return str
}

func assignStatToString(stat *ast.AssignStat) string {
	str := ""
	for i, name := range stat.VarList {
		str += expToString(name)
		if i < len(stat.VarList)-1 {
			str += ", "
		}
	}

	str += " = "

	for i, exp := range stat.ExpList {
		str += expToString(exp)
		if i < len(stat.ExpList)-1 {
			str += ", "
		}
	}

	return str
}

func localVarDefStatToString(stat *ast.LocalVarDeclStat) string {
	str := "local "
	for i, name := range stat.NameList {
		str += name
		if i < len(stat.NameList)-1 {
			str += ", "
		}
	}

	str += " = "
	for i, exp := range stat.ExpList {
		str += expToString(exp)
		if i < len(stat.ExpList)-1 {
			str += ", "
		}
	}

	return str
}
