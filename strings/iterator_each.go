package strings

//Each runs the function 'do' on each item in the slice
//and returns the resulting slice
func Each(strs []string, do func(string) string) []string {
	var out []string
	for _, s := range strs {
		out = append(out, do(s))
	}
	return out
}
