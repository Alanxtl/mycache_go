package tools

func CloneBytes(bytes []byte) []byte {
	c := make([]byte, len(bytes))
	copy(c, bytes)
	return c
}
