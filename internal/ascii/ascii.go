package ascii

func EqualFold(s, t string) bool {
	for i, n := 0, min(len(s), len(t)); i < n; i++ {
		sr := s[i]
		tr := t[i]

		if tr == sr {
			continue
		}

		if tr < sr {
			tr, sr = sr, tr
		}

		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			continue
		}

		return false
	}

	return len(s) == len(t)
}

// ToLowerBytes lowercases ASCII letters of b in place.
func ToLowerBytes(b []byte) {
	for i, c := range b {
		if c-'A' <= 'Z'-'A' {
			b[i] = c | ('a' - 'A')
		}
	}
}
