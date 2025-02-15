package cache

import "github.com/Alanxtl/mycache_go/pkg/tools"

type ByteView struct {
	Bytes []byte
}

func (b ByteView) Len() int {
	return len(b.Bytes)
}

func (b ByteView) ByteSlice() []byte {
	return tools.CloneBytes(b.Bytes)
}

func (b ByteView) String() string {
	return string(b.Bytes)
}
