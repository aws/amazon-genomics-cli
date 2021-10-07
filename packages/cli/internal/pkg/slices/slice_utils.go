package slices

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
