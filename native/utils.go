package native

import (
	"fmt"
	"io"
	"strings"
)

func createTag(name, value string) string {
	return name + `:"` + value + `"`
}

func writePackageHeader(w io.Writer, imports []string) error {
	data := `package odata

`
	if len(imports) > 0 {
		data += fmt.Sprintf(`import (
"%s"
)
`, strings.Join(imports, `"
"`))
	}
	_, err := w.Write([]byte(data))
	return err
}
