package trie

func longestCommonPrefix(a, b []byte) int {
	var i = 0
	for ; i < len(b) && i < len(a) && b[i] == a[i]; i++ {
	}
	return i
}
