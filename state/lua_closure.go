package state

import (
	"vczn/luago/api"
	"vczn/luago/binchunk"
)

type luaClosure struct {
	proto  *binchunk.ProtoType
	goFunc api.GoFunction
}

func newLuaClosure(proto *binchunk.ProtoType) *luaClosure {
	return &luaClosure{proto: proto}
}

func newGoClosure(gofunc api.GoFunction) *luaClosure {
	return &luaClosure{goFunc: gofunc}
}
