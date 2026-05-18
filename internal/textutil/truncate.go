package textutil

func TruncateRunes(value string, limit int) string {
	if limit <= 0 {
		return ""
	}

	count := 0
	for index := range value {
		if count == limit {
			return value[:index]
		}

		count++
	}

	return value
}
