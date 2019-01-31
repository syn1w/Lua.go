package parser

import (
	"testing"
	"vczn/luago/compiler/lexer"
)

func TestExpLiteral(t *testing.T) {
	testExp(t, `nil`)
	testExp(t, `true`)
	testExp(t, `false`)
	testExp(t, `123`)
	testExp(t, `foo`)
	testExp(t, `{}`)
	testExp(t, `...`)
	testExp2(t, `0xFF`, `255`)
}

func TestUnopExp(t *testing.T) {
	testExp(t, `-128`)
	testExp2(t, `-a`, `-(a)`)
	testExp2(t, `-0xFF`, `-255`)
	testExp2(t, `~a`, `~(a)`)
	testExp2(t, `#'foo'`, `#('foo')`)
	testExp2(t, `not a`, `not(a)`)
	testExp2(t, `~0xFF`, `-256`)
	testExp2(t, `not true`, `false`)
	testExp2(t, `- - - - - 1`, `-1`)
	testExp2(t, `- - - - - - 1`, `1`)
}

func TestBinopExp(t *testing.T) {
	testExp2(t, `a^b^c`, `(a ^ (b ^ c))`)
	testExp2(t, `1 | -2`, `-1`)
	testExp2(t, `0xF0F ~ 0xF0`, `4095`)
	testExp2(t, `0xF & 0xFF & 0xFFF`, `15`)
	testExp2(t, `0xF | 0xFF | 0x0`, `255`)
	testExp2(t, `a ^ b ^ c ^ d`, `(a ^ (b ^ (c ^ d)))`) // right associative
	testExp2(t, `4^3^2`, `262144.000000`)
	testExp2(t, `a^3^2`, `(a ^ 9.000000)`)
	testExp2(t, `4^3^a`, `(4 ^ (3 ^ a))`)
	testExp2(t, `2^-2`, `0.250000`)
	testExp2(t, `(2^2)^3`, `64.000000`)
	testExp2(t, `x ^ (-1 / 3)`, `(x ^ -0.333333)`)
	testExp2(t, `-x^2`, `-((x ^ 2))`)
	testExp2(t, `-1^a`, `-((1 ^ a))`)
	testExp2(t, `1 + 2 + 3`, `6`)
	testExp2(t, `2.3+1.9`, `4.200000`)
	testExp2(t, `1+2+a`, `(3 + a)`)
	testExp2(t, `a+1+2`, `((a + 1) + 2)`)
	testExp2(t, `n-1`, `(n - 1)`)
	testExp2(t, `a + b - c + d`, `(((a + b) - c) + d)`)
	testExp2(t, `a + b - c*d`, `((a + b) - (c * d))`)
	testExp2(t, `(a + b) // (c + d)`, `((a + b) // (c + d))`)
	testExp2(t, `a or b or c`, `((a or b) or c)`)
	testExp2(t, `a or b and c`, `(a or (b and c))`)
	testExp2(t, `a / b`, `(a / b)`)
	testExp2(t, `42 % 4`, `2`)
	testExp2(t, `a % b`, `(a % b)`)
	testExp2(t, `5+x^2*8`, `(5 + ((x ^ 2) * 8))`)
	testExp2(t, `5 / 2`, `2.500000`)
	testExp2(t, `5 // 2`, `2`)
	testExp2(t, `true or false or 2 or nil or "foo"`, `true`)
	testExp2(t, `false and true and nil and 0 and a`, `false`)
	testExp2(t, `true and 1 and "foo" and a`, `a`)
	testExp2(t, `true and x and true and x and true`,
		`(((x and true) and x) and true)`)
	testExp2(t, `nil and true or false`, `false`)
	testExp2(t, `false or x`, `x`)
	testExp2(t, `nil and true or false`, `false`)
	testExp2(t, `((((a + b))))`, `(a + b)`)
	testExp2(t, `((((a))))`, `(a)`)
	testExp2(t, `a >> b`, `(a >> b)`)
	testExp2(t, `1 << 2`, `4`)
	testExp2(t, `a << b`, `(a << b)`)
	testExp2(t, `a < b`, `(a < b)`)
	testExp2(t, `2 > 1`, `(2 > 1)`)
	testExp2(t, `a == b`, `(a == b)`)
	testExp2(t, `x ~= y`, `(x ~= y)`)
	testExp2(t, `a <= b`, `(a <= b)`)
	testExp2(t, `a >= b`, `(a >= b)`)
	testExp2(t, `a+i < b/2+1`, `((a + i) < ((b / 2) + 1))`)
	testExp2(t, `a<y and y<=z`, `((a < y) and (y <= z))`)
	testExp2(t, `'hello' .. 42`, `'hello' .. 42`)
}

func TestTcExp(t *testing.T) {
	testExp(t, `{}`)
	testExp(t, `{...,}`)
	testExp2(t, `{f(),}`, `{f(),}`)
	testExp2(t, `{f(), nil}`, `{f(),nil,}`)
	testExp2(t, `{[f(1)] = g, 'x', 'y', x = 1, f(x), [30] = 23, 45}`,
		`{[f(1)]=g,'x','y',['x']=1,f(x),[30]=23,45,}`)
	testExp2(t, `{[f(1)] = g; "x", "y"; x = 1, f(x), [30] = 23; 45 }`,
		`{[f(1)]=g,'x','y',['x']=1,f(x),[30]=23,45,}`)
}

func TestPrefixExp(t *testing.T) {
	testExp(t, `name`)
	testExp(t, `(name)`)
	testExp(t, `name[key]`)
	testExp2(t, `name.field`, `name['field']`)
	testExp2(t, `a.b.c.d.e`, `a['b']['c']['d']['e']`)
	testExp(t, `a[b][c][d][e]`)
	testExp(t, `a[b[c[d[e]]]]`)
}

func TestFuncCallExp(t *testing.T) {
	testExp2(t, `print ''`, `print('')`)
	testExp2(t, `print 'hello'`, `print('hello')`)
	testExp2(t, `print {}`, `print({})`)
	testExp2(t, `print {1}`, `print({1,})`)
	testExp(t, `f()`)
	testExp(t, `g(f(), x)`)
	testExp(t, `g(x, f())`)
	testExp2(t, `io.read('*n')`, `io['read']('*n')`)
	testExp2(t, `t.a.b.c:f(a, b)`, `t['a']['b']['c']:f(a, b)`)
}

func testExp(t *testing.T, want string) {
	exp := parseExp(lexer.NewLexer(want, "string"))
	got := expToString(exp)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func testExp2(t *testing.T, src, want string) {
	exp := parseExp(lexer.NewLexer(src, "string"))
	got := expToString(exp)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func TestStat(t *testing.T) {
	testStat(t, `;`)
	testStat(t, `break`)
	testStat(t, `::label::`)
	testStat(t, `goto label`)
	testStat2(t, `do ; end`, `do  end`)
	testStat2(t, `repeat ; until true`, `repeat  until true`)
	testStat2(t, `for v = 1, 100, 1 do ; end`, `for v = 1, 100, 1 do  end`)
	testStat2(t, `function foo() end`, `foo = function()  end`)
	testStat(t, `local function foo()  end`)
	testStat(t, `local a = 1`)
}

func TestIfStat(t *testing.T) {
	testStat2(t, `if true then ; end`, `if true then  end`)
	testStat2(t, `if a then ; else ; end`, `if a then  elseif true then  end`)
	testStat2(t, `if a then ; elseif b then ; else ; end`,
		`if a then  elseif b then  elseif true then  end`)
	testStat2(t, `if a then x = 42 else x = 3.14 end`,
		`if a then x = 42 elseif true then x = 3.140000 end`)
}

func TestFuncCallStat(t *testing.T) {
	testStat(t, `print()`)
	testStat(t, `print(i)`)
	testStat2(t, `print("hello!")`, `print('hello!')`)
	testStat2(t, `fact(n-1)`, `fact((n - 1))`)
	testStat(t, `t:f()`)
	testStat2(t, `assert((4 and 5) == 5)`, `assert((5 == 5))`) // 4 is false => 5
	testStat2(t, `assert((4 & 5) == 4)`, `assert((4 == 4))`)
}

func TestStatAssignStat(t *testing.T) {
	testStat(t, `x = 42`)
	testStat(t, `a = x`)
	testStat2(t, `f().a = 1`, `f()['a'] = 1`)
	testStat2(t, `a = io.read('*n')`, `a = io['read']('*n')`)
	testStat2(t, `foo, bar = 3.14, 42`, `foo, bar = 3.140000, 42`)
	testStat(t, `local x = 42`)
	testStat2(t, `local f = function() end`, `local f = function()  end`)
	testStat2(t, `function f() return 1 end`, `f = function() return 1 end`)
}

func TestBlock(t *testing.T) {
	testBlock(t, `return`)
	testBlock(t, `return 1`)
	testBlock2(t, `return n * fact(n-1)`, `return (n * fact((n - 1)))`)
}

func TestFuncDefStat(t *testing.T) {
	testStat2(t, `function f() end`, `f = function()  end`)
	testStat2(t, `function f(a) end`, `f = function(a)  end`)
	testStat2(t, `function f(a, b) end`, `f = function(a, b)  end`)
	testStat2(t, `function f(a, b, ...) end`, `f = function(a, b, ...)  end`)
	testStat2(t, `function t.a.b.c.f() end`, `t['a']['b']['c']['f'] = function()  end`)
	testStat2(t, `function t.a.b.c:f() end`, `t['a']['b']['c']['f'] = function(self)  end`)
	testStat2(t, `local function f(a) end`, `local function f(a)  end`)
}

func TestError(t *testing.T) {
	testError(t, `f()   = 1`, `string:1: syntax error near '='`)
	testError(t, `a,f() = 1`, `string:1: syntax error near '='`)
	testError(t, `(a)   = 1`, `string:1: syntax error near '='`)
	testError(t, `(a+b) = 1`, `string:1: syntax error near '='`)
	testError(t, `(a.b) = 1`, `string:1: syntax error near '='`)
	testError(t, `a + b = 1`, `string:1: syntax error near '+'`)
	testError(t, `f(,)`, `string:1: syntax error near ','`)
}

func testStat(t *testing.T, want string) {
	stat := parseStat(lexer.NewLexer(want, "string"))
	got := statToString(stat)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func testStat2(t *testing.T, src, want string) {
	stat := parseStat(lexer.NewLexer(src, "string"))
	got := statToString(stat)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func testBlock(t *testing.T, want string) {
	block := parseBlock(lexer.NewLexer(want, "string"))
	got := blockToString(block)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func testBlock2(t *testing.T, src, want string) {
	block := parseBlock(lexer.NewLexer(src, "string"))
	got := blockToString(block)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func testError(t *testing.T, src, want string) {
	got := getError(src)
	if want != got {
		t.Errorf("want='%s', got='%s'", want, got)
	}
}

func getError(src string) (err string) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(string)
		}
	}()

	parseBlock(lexer.NewLexer(src, "string"))
	return
}
