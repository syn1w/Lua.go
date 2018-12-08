package number

import (
	"math"
	"strconv"
)

// arithmetic operator
// +, -(u), -(b), *, /, //(IDIV), %, ^

// 除法和乘方运算先转换为浮点数，再进行运算，计算结果也是浮点数
// 其他先判断操作数是否都为整数，如果是，进行整数运算；否则，转换为浮点数。
// 乘方运算为右结合，比如 4^3^2 == 4^(3^2)

// IFloorDiv is integer floor division
func IFloorDiv(a, b int64) int64 {
	// 整除向下截断(需要注意负数向负无穷截断，其他许多语言向零截断)
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

// FFloorDiv is float floor division
func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

// IMod is integer modular
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

// FMod is float modular
func FMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}

// bitwise operator(lua 5.3)
// & | ~(一元取反，二元异或) << >>
// Lua5.1 used bit.xxx(); Lua5.2 used bit32; Lua5.3 used builtin operators
// 位运算先把操作数转化为整数再进行运算，计算结果也为整数

// ShiftLeft is a << n
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	}
	return ShiftRight(a, -n) // n < 0 向相反方向移动 n 个比特
}

// ShiftRight is a >> n
func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n)) // logical right left
	}
	return ShiftLeft(a, -n)
}

// comparison operator
// ==, ~=(!=), >, >=, <, <=

// logical operator
// and, or, not

// len operator
// #, e.g. #"hello" --> 5; #{1,2,3} --> 3

// string concat
// a .. b

// 隐式类型转换
// https://cloudwu.github.io/lua53doc/manual.html#3.4.3

// FloatToInteger converts float64 to int64,
// and returns true if the float represents an integer, such as 3.0
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}

// ParseInteger parses str to int64
func ParseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

// ParseFloat parses str to float64
func ParseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}
