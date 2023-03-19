package utils

func CompareKeywords(keywords1, keywords2 []string) bool {
	for _, keyword := range keywords1 {
		if Contains(keywords2, keyword) {
			return true
		}
	}
	return false
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
