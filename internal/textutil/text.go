package textutil

// TruncateRunes returns value truncated to at most max Unicode code points.
func TruncateRunes(value string, max int) string {
	if max <= 0 {
		return ""
	}

	count := 0
	for index := range value {
		if count == max {
			return value[:index]
		}

		count++
	}

	return value
}
