package text

func IsUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func IsLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func IsDigit (r rune) bool {
	return r >= '0' && r <= '9'
}

func IsUpperOrLower (r rune) bool {
	return IsLower(r) || IsUpper(r)
}