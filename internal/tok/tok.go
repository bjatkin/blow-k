package tok

// Token is a valid token string
type Token string

// Tokenize converts a string into an array of Tokens
func Tokenize(data string) []Token {
	var current []rune
	var ret []Token

	for _, c := range data {
		if !isSplit(c) {
			current = append(current, c)
			continue
		}

		if len(current) > 0 {
			ret = append(ret, Token(current))
			current = []rune{}
		}

		ret = append(ret, Token(c))
	}

	if len(current) > 0 {
		ret = append(ret, Token(current))
	}

	return ret
}

// isSplit returns true if the token is a valid spliter
func isSplit(token rune) bool {
	ref := []rune{
		' ', '\t', '\n', // white space tokens
		'(', ')', '{', '}', '[', ']', // parens etc. tokens
		':', '.', ',', '$', '"', // punctuation tokens
		'+', '-', '*', '/', '&', '|', '=', // math tokens
	}

	for _, r := range ref {
		if r == token {
			return true
		}
	}

	return false
}
