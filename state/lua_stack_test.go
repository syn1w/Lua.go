package state

import "testing"

func TestStack1(t *testing.T) {
	stack := newLuaStack(10, nil)

	if !stack.empty() {
		t.Error("1 stack is not emtpy, newLuaStack error")
	}

	for i := 0; i < 10; i++ {
		stack.push(i)
	}

	if !stack.full() {
		t.Errorf("2 stack is not full, push or full error")
	}

	for i := 9; i >= 0; i-- {
		got := stack.pop()
		want := i
		if got != want {
			t.Errorf("pop, got = %d, want= %d", got, want)
		}
	}
}

func TestStack2(t *testing.T) {
	stack := newLuaStack(10, nil)
	for i := 0; i < 10; i++ {
		stack.push(i)
	}

	stack.check(3)
	if stack.full() {
		t.Error("1 `check` error")
	}
}
