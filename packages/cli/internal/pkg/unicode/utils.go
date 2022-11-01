package unicode

import "sort"

// DeDuplicateStrings sorts and de-duplicates a slice of strings
func DeDuplicateStrings(strs []string) []string {
	if len(strs) == 0 {
		return strs
	}
	sort.Strings(strs)

	var last string
	firstItem := true
	var dedupped []string
	for _, s := range strs {
		if firstItem {
			dedupped = append(dedupped, s)
			firstItem = false
		} else {
			if s != last {
				dedupped = append(dedupped, s)
			}
		}
		last = s
	}

	return dedupped
}

// SubString performs a unicode aware substring operation on 'str'. Will panic if start or length are out of bounds
func SubString(str string, start int, length int) string {
	runes := []rune(str)
	return string(runes[start : start+length])
}
