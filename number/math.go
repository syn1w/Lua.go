package number

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// arithmetic operator
// +, -(u), -(b), *, /, //(IDIV), %, ^

// 除法和乘方运算先转换为 floating point，再进行运算，计算结果也是 floating point
// 其他先判断操作数是否都为整数，如果是，进行整数运算；否则，转换为 floating point
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
// 位运算先把操作数转化为整数再进行运算，计算结果也为 integer

// ShiftLeft is a << n
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	}
	return ShiftRight(a, -n) // n < 0 向相反方向移动 n bytes
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

var reInteger = regexp.MustCompile(`^[+-]?[0-9]+$|^-?0x[0-9a-f]+$`)

// ParseInteger parses str to int64
func ParseInteger(str string) (int64, bool) {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	if !reInteger.MatchString(str) {
		return 0, false
	}

	if str[0] == '+' {
		str = str[1:]
	}

	if strings.Index(str, "0x") < 0 { // decimal
		i, err := strconv.ParseInt(str, 10, 64)
		return i, err == nil
	}

	// hex
	var sign int64
	if str[0] == '-' {
		sign = -1
		str = str[3:] // -0x
		fmt.Println("-")
	} else {
		sign = 1
		str = str[2:] // 0x
	}

	if len(str) > 16 {
		str = str[len(str)-16:] // cut the long hex string
	}

	// ParseInt(0xFFFFFFFFFFFFFFFF) out of range
	i, err := strconv.ParseUint(str, 16, 64)
	return sign * int64(i), err == nil
}

// ParseFloat parses str to float64
func ParseFloat(str string) (float64, bool) {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	if strings.Contains(str, "nan") || strings.Contains(str, "inf") {
		return 0, false
	}

	if strings.HasPrefix(str, "0x") && len(str) > 2 {
		return parseHexFloat(str[2:])
	}

	if strings.HasPrefix(str, "-0x") && len(str) > 3 {
		f, ok := parseHexFloat(str[3:])
		return -f, ok
	}

	if strings.HasPrefix(str, "+0x") && len(str) > 3 {
		return parseHexFloat(str[3:])
	}

	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}

// Group1: ([0-9a-f]+(\.[0-9a-f]*)?|([0-9a-f]*\.[0-9a-f]+))
// Alternative1: [0-9a-f]+(\.[0-9a-f]*)?
//   matched `hh.`
// Alternative2: ([0-9a-f]*\.[0-9a-f]+)
//   matched `.hh`
// Group2: (p[+\-]?[0-9]+)?
//   mathed p[+-]dd
var reHexFloat = regexp.MustCompile(`^([0-9a-f]+(\.[0-9a-f]*)?|([0-9a-f]*\.[0-9a-f]+))(p[+/-]?[0-9]+)`)

// ABC.DEFp10
func parseHexFloat(str string) (float64, bool) {
	var i16, f16, p10 float64

	if !reHexFloat.MatchString(str) {
		return 0, false
	}

	if idxP := strings.Index(str, "p"); idxP > 0 {
		digits := str[idxP+1:]
		str = str[:idxP]

		var sign float64
		sign = 1.0
		if digits[0] == '-' {
			sign = -1
		}
		if digits[0] == '+' || digits[0] == '-' {
			digits = digits[1:]
		}

		// pdd or hhp
		if len(str) == 0 || len(digits) == 0 {
			return 0, false
		}

		for i := 0; i < len(digits); i++ {
			if x, ok := parseDigit(digits[i], 10); ok {
				p10 = p10*10 + x
			} else {
				return 0, false
			}
		}

		p10 = sign * p10
	}

	if idxDot := strings.Index(str, "."); idxDot > 0 {
		digits := str[idxDot+1:]
		str = str[:idxDot]

		// .
		if len(str) == 0 && len(digits) == 0 {
			return 0, false
		}

		for i := len(digits) - 1; i >= 0; i-- {
			if x, ok := parseDigit(digits[i], 16); ok {
				f16 = (f16 + x) / 16
			} else {
				return 0, false
			}
		}
	}

	for i := 0; i < len(str); i++ {
		if x, ok := parseDigit(str[i], 16); ok {
			i16 = i16*16 + x
		} else {
			return 0, false
		}
	}

	f := i16 + f16
	if p10 != 0 {
		// * 2^p
		f *= math.Pow(2, p10)
	}

	return f, true
}

func parseDigit(digit byte, base int) (float64, bool) {
	if base == 10 || base == 16 {
		switch digit {
		case '0':
			return 0, true
		case '1':
			return 1, true
		case '2':
			return 2, true
		case '3':
			return 3, true
		case '4':
			return 4, true
		case '5':
			return 5, true
		case '6':
			return 6, true
		case '7':
			return 7, true
		case '8':
			return 8, true
		case '9':
			return 9, true
		}
	}
	if base == 16 {
		switch digit {
		case 'a':
			return 10, true
		case 'b':
			return 11, true
		case 'c':
			return 12, true
		case 'd':
			return 13, true
		case 'e':
			return 14, true
		case 'f':
			return 15, true
		}
	}
	return -1, false
}
