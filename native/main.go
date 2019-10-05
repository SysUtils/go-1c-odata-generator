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

func (g *Generator) Generate() error {
	f, _ := os.Create("odata/entity.go")
	err := writePackageHeader(f, []string{"github.com/SysUtils/go-1c-odata/types", "github.com/SysUtils/go-1c-odata/client"})
	if err != nil {
		return err
	}
	for _, v := range g.schema.Entities {
		err = g.writeEntity(f, v)
		if err != nil {
			return err
		}
	}

	f.Close()

	f, _ = os.Create("odata/complex.go")
	err = writePackageHeader(f, []string{"github.com/SysUtils/go-1c-odata/types"})
	if err != nil {
		return err
	}
	for _, v := range g.schema.Complexes {
		err = g.writeComplexStruct(f, v)
		if err != nil {
			return err
		}
	}

	f.Close()
	return nil
}
