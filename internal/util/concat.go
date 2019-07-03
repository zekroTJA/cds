package util

func ConcatToString(arr ...[]byte) string {
	var ln, cursor int

	for _, a := range arr {
		ln += len(a)
	}

	res := make([]byte, ln)

	for _, a := range arr {
		ln = len(a)
		copy(res[cursor:cursor+ln], a)
		cursor += ln
	}

	return string(res)
}
