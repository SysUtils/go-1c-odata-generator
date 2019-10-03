package native

import (
	"fmt"
	"io"
)

func (g *Generator) writeEntity(w io.Writer, entity Entity) error {
	if err := g.writeEntityStruct(w, entity); err != nil {
		return err
	}
	if err := g.writeEntityKeyStruct(w, entity); err != nil {
		return err
	}
	if err := g.writeEntityDataStruct(w, entity); err != nil {
		return err
	}
	return g.writeEntityCommonFunctions(w, entity)
}

func (g *Generator) writeEntityStruct(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("type %s struct { \n", entity.TypeName)
	data += fmt.Sprintf("	_ struct{} `%s` \n", createTag("typename", entity.OriginalTypeName))
	data += fmt.Sprintf("	%sKey \n", entity.TypeName)
	data += fmt.Sprintf("	%sData \n", entity.TypeName)
	data += "} \n"
	w.Write([]byte(data))
	return nil
}

func (g *Generator) writeEntityKeyStruct(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("type %sKey struct { \n", entity.TypeName)
	data += fmt.Sprintf("	_ struct{} `%s` \n", createTag("typename", entity.OriginalTypeName))
	for _, k := range entity.KeyProperties {
		data += fmt.Sprintf("	%s *%s `%s` \n", k.Name, k.Type, createTag("json", k.OriginalName+",omitempty"))
	}
	data += "} \n"
	w.Write([]byte(data))
	return nil
}

func (g *Generator) writeEntityDataStruct(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("type %sData struct { \n", entity.TypeName)
	data += fmt.Sprintf("	_ struct{} `%s` \n", createTag("typename", entity.OriginalTypeName))
	for _, k := range entity.Properties {
		data += fmt.Sprintf("	%s *%s `%s` \n", k.Name, k.Type, createTag("json", k.OriginalName+",omitempty"))
	}
	data += "} \n"
	w.Write([]byte(data))
	return nil
}

func (g *Generator) writeEntityCommonFunctions(w io.Writer, entity Entity) error {
	if err := g.writeEntityDataFunc(w, entity); err != nil {
		return err
	}
	return g.writeEntityKeyFunc(w, entity)
}

func (g *Generator) writeEntityKeyFunc(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("func (e *%s) _Key() %sKey { \n", entity.TypeName, entity.TypeName)
	data += fmt.Sprintf("	return e.%sKey \n", entity.TypeName)
	data += "} \n"
	w.Write([]byte(data))
	return nil
}

func (g *Generator) writeEntityDataFunc(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("func (e *%s) _Data() %sData { \n", entity.TypeName, entity.TypeName)
	data += fmt.Sprintf("	return e.%sData \n", entity.TypeName)
	data += "} \n"
	w.Write([]byte(data))
	return nil
}
