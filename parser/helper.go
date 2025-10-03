package parser

import "unicode"

func isPublic(name string) bool {
	if len(name) > 0 {
		return unicode.IsUpper([]rune(name)[0]) // https://www.reddit.com/r/golang/comments/11cig0a/comment/ja371qd/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
	}
	return false
}
