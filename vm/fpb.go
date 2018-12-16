package vm

// similar to minifloat, except sign bit
// https://en.wikipedia.org/wiki/Minifloat

// Int2fb converts a int to a `floating point byte`
func Int2fb(x int) int {
	e := 0
	if x < 8 { // eeeee == 0
		return x
	}

	for x >= (8 << 4) { // 128 coarse steps
		x = (x + 0xf) >> 4 // x = ceil(x / 16)
		e += 4
	}
	for x >= (8 << 1) { // 16
		x = (x + 1) >> 1 // x = ceil(x / 2)
		e++
	}

	return ((e + 1) << 3) | (x - 8)
}

// Fb2int Float byte to int
// x is eeeeexxx `float byte`
// if eeeee == 0 then xxx
// else 1xxx * 2^(eeeee-1)
func Fb2int(x int) int {
	if x < 8 { // eeeee == 0
		return x
	}
	return ((x & 7) + 8) << uint((x>>3)-1)
}
