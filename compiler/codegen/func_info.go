package codegen

type funcInfo struct {
	constants  map[interface{}]int // constants table
	usedRegs   int
	maxRegs    int
	scopeDepth int // start from 0
	locVars    []*locVarInfo
	locNames   map[string]*locVarInfo
	// TODO
}

type locVarInfo struct {
	prev       *locVarInfo // linked list
	name       string
	scopeDepth int // start from 0
	slot       int // index of binding local variable
	captured   bool
}

func (fi *funcInfo) indexOfConstant(k interface{}) int {
	if idx, found := fi.constants[k]; found {
		return idx
	}
	idx := len(fi.constants)
	fi.constants[k] = idx
	return idx
}

// allocate a register and return the index of the register
func (fi *funcInfo) allocReg() int {
	fi.usedRegs++
	if fi.usedRegs >= 255 {
		panic("function or expression needs too many registers")
	}

	if fi.usedRegs > fi.maxRegs {
		fi.maxRegs = fi.usedRegs
	}

	return fi.usedRegs - 1
}

// free a register
func (fi *funcInfo) freeReg() {
	fi.usedRegs--
}

// allocate n registers and return the index of the first register
func (fi *funcInfo) allocRegs(n int) int {
	for i := 0; i < n; i++ {
		fi.allocReg()
	}
	return fi.usedRegs - n
}

// free n registers
func (fi *funcInfo) freeRegs(n int) {
	for i := 0; i < n; i++ {
		fi.freeReg()
	}
}

func (fi *funcInfo) enterScope() {
	fi.scopeDepth++
}

func (fi *funcInfo) exitScope() {
	fi.scopeDepth--
	for _, locVar := range fi.locNames {
		if locVar.scopeDepth > fi.scopeDepth {
			fi.removeLocVar(locVar)
		}
	}
}

// add a local variable and return the index of the variable
func (fi *funcInfo) addLocVar(name string) int {
	locVar := &locVarInfo{
		name:       name,
		prev:       fi.locNames[name],
		scopeDepth: fi.scopeDepth,
		slot:       fi.allocReg(),
	}

	fi.locVars = append(fi.locVars, locVar)
	fi.locNames[name] = locVar
	return locVar.slot
}

func (fi *funcInfo) removeLocVar(locVar *locVarInfo) {
	fi.freeReg()
	if locVar.prev == nil {
		delete(fi.locNames, locVar.name)
	} else if locVar.prev.scopeDepth == locVar.scopeDepth {
		fi.removeLocVar(locVar.prev)
	} else {
		fi.locNames[locVar.name] = locVar.prev
	}
}

// return the slot of local variable if the local variable bind to a register
// else return -1
func (fi *funcInfo) slotOfLocVar(name string) int {
	if locVar, found := fi.locNames[name]; found {
		return locVar.slot
	}
	return -1
}
