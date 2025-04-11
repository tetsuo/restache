package commonpath

import "strings"

func CommonPathUnix(paths []string) string {
	n := len(paths)
	if n == 0 {
		return ""
	}
	split := make([][]string, n)
	for i, p := range paths {
		split[i] = strings.Split(p, "/")
	}
	common := split[0]
	for _, segs := range split[1:] {
		common = intersectUnix(common, segs)
		if len(common) == 0 {
			return "/"
		}
	}
	return reassembleUnix(common)
}

func CommonPathWin(paths []string) string {
	n := len(paths)
	if n == 0 {
		return ""
	}
	split := make([][]string, n)
	for i, p := range paths {
		p = strings.ReplaceAll(p, "\\", "/")
		split[i] = strings.Split(p, "/")
	}
	common := split[0]
	for _, segs := range split[1:] {
		common = intersectWin(common, segs)
		if len(common) == 0 {
			return `\`
		}
	}
	return reassembleWin(common)
}

func intersectUnix(a, b []string) []string {
	n := min(len(a), len(b))
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return a[:i]
}

func intersectWin(a, b []string) []string {
	n := min(len(a), len(b))
	i := 0
	for i < n && strings.EqualFold(a[i], b[i]) {
		i++
	}
	return a[:i]
}

func reassembleUnix(seg []string) string {
	if len(seg) == 0 {
		return "/"
	}
	if seg[0] == "" {
		return "/" + strings.Join(seg[1:], "/")
	}
	return strings.Join(seg, "/")
}

func reassembleWin(seg []string) string {
	if len(seg) == 0 {
		return `\`
	}
	if len(seg) == 1 && strings.HasSuffix(seg[0], ":") {
		return seg[0] + `\`
	}
	return strings.Join(seg, `\`)
}
