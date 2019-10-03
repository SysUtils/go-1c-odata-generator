package native

import "io"

func createTag(name, value string) string {
	return name + `:"` + value + `"`
}

func writePackageHeader(w io.Writer, inports []string) {
	data := `package odata

`
	w.Write([]byte(data))
}
