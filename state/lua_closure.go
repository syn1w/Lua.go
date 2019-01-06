package state

import "vczn/luago/binchunk"

type luaClosure struct {
	proto *binchunk.ProtoType
}

func newLuaClosure(proto *binchunk.ProtoType) *luaClosure {
	return &luaClosure{proto: proto}
}
