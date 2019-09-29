package clippy

func longestStringLength(a []string) int {
	longestLength := 0
	for _, s := range a {
		currLength := len(s)
		if currLength > longestLength {
			longestLength = currLength
		}
	}
	return longestLength
}
