package gocache

//抽象出一个只读的数据结构来表示缓存值
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

<<<<<<< HEAD
//返回拷贝，防止缓存值被外部修改
=======
>>>>>>> a5b6419c0fb7f1a2857411e5a29c0d2f87b7dd31
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
