package tools

func SetLastLash(text string) string {
	last := text[len(text)-1:]

	if last != "/" {
		return text + "/"

	}
	return text

}
