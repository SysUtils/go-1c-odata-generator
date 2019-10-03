package native

import (
	"github.com/SysUtils/go-1c-odata-generator/shared"
	"os"
)

type Generator struct {
	TypeMap map[string]string
	NameMap map[string]string
	schema  Schema
}

func NewGenerator(schema shared.Schema, typeMap map[string]string, nameMap map[string]string) (*Generator, error) {
	gen := Generator{TypeMap: typeMap, NameMap: nameMap}
	convertedSchema, err := gen.convertSchema(schema)
	if err != nil {
		return nil, err
	}
	gen.schema = *convertedSchema
	return &gen, err
}

func (g *Generator) Generate() {
	f, _ := os.Create("odata/entity.go")
	writePackageHeader(f, nil)
	for _, v := range g.schema.Entities {
		g.writeEntity(f, v)
	}

	f.Close()

	f, _ = os.Create("odata/complex.go")
	writePackageHeader(f, nil)
	for _, v := range g.schema.Complexes {
		g.writeComplexStruct(f, v)
	}

	f.Close()
}
