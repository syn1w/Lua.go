package state

import (
	"vczn/luago/api"
	"vczn/luago/binchunk"
)

type luaClosure struct {
	proto  *binchunk.ProtoType
	goFunc api.GoFunction
	upvals []*upvalue // non-local variables captured by closure
}

type upvalue struct {
	val *LuaValue
}

func newLuaClosure(proto *binchunk.ProtoType) *luaClosure {
	c := &luaClosure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(gofunc api.GoFunction, nUpvals int) *luaClosure {
	c := &luaClosure{goFunc: gofunc}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
