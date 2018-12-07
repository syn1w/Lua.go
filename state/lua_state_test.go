package state

import (
	"fmt"
	"log"
	"testing"
	"vczn/luago/api"
)

func printStack(ls *LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LuaTBoolean:
			log.Printf("[%t]", ls.ToBoolean(i))
		case api.LuaTNumber:
			log.Printf("[%g]", ls.ToNumber(i))
		case api.LuaTString:
			log.Printf("[%q]", ls.ToString(i))
		default:
			log.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func TestState1(t *testing.T) {
	ls := NewLuaState()
	ls.PushBoolean(true)
	printStack(ls)
	ls.PushInteger(10)
	printStack(ls)
	ls.PushNil()
	printStack(ls)
	ls.PushString("hello")
	printStack(ls)
	ls.PushValue(-4)
	printStack(ls)
	ls.Replace(3)
	printStack(ls)
	ls.SetTop(6)
	printStack(ls)
	ls.Remove(-3)
	printStack(ls)
	ls.SetTop(-5)
	printStack(ls)
}
