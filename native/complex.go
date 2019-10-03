package native

import (
	"fmt"
	"io"
)

func (g *Generator) writeComplexStruct(w io.Writer, entity ComplexType) error {
	data := fmt.Sprintf("type %s struct { \n", entity.TypeName)
	for _, k := range entity.Properties {
		data += fmt.Sprintf("	%s *%s `%s` \n", k.Name, k.Type, createTag("json", k.OriginalName+",omitempty"))
	}
	data += "} \n"
	w.Write([]byte(data))
	return nil
}
