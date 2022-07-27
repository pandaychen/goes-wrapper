package strhash

func Hash33(src string) int {
	hash := 5381
	for i := 0; i < len(src); i++ {
		hash = hash<<5 + hash + int(src[i])
	}
	return hash & 0x7fffffff
}
