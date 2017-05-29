package strings

//Each runs the function 'ef' on each item in the slice
//and returns the resulting slice
func Each(strs []string, ef func(string) string) []string {
	var out []string
	for _, s := range strs {
		out = append(out, ef(s))
	}
	return out
}
