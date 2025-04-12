package restache

func escapeComment(w writer, s string) error {
	if len(s) == 0 {
		return nil
	}

	i := 0
	for j := 0; j < len(s)-1; j++ {
		if s[j] == '*' && s[j+1] == '/' {
			if i < j {
				if _, err := w.WriteString(s[i:j]); err != nil {
					return err
				}
			}
			if _, err := w.WriteString("*\\/"); err != nil { // escape the '/'
				return err
			}
			i = j + 2
			j++ // skip the '/'
		}
	}

	if i < len(s) {
		if _, err := w.WriteString(s[i:]); err != nil {
			return err
		}
	}
	return nil
}
