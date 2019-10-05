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
	if err := g.writeEntityActions(w, entity); err != nil {
		return err
	}
	if err := g.writeEntityNavigations(w, entity); err != nil {
		return err
	}

	return g.writeEntityCommonFunctions(w, entity)
}

func (g *Generator) writeEntityStruct(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("type %s struct { \n", entity.TypeName)
	data += fmt.Sprintf("	_ struct{} `%s` \n", createTag("typename", entity.OriginalTypeName))
	data += fmt.Sprintf("	Client *client.Client \n")
	data += fmt.Sprintf("	%sKey \n", entity.TypeName)
	data += fmt.Sprintf("	%sData \n", entity.TypeName)
	data += "} \n"
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntityKeyStruct(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("type %sKey struct { \n", entity.TypeName)
	data += fmt.Sprintf("	_ struct{} `%s` \n", createTag("typename", entity.OriginalTypeName))
	for _, k := range entity.KeyProperties {
		data += fmt.Sprintf("	%s *%s `%s` \n", k.Name, k.Type, createTag("json", k.OriginalName+",omitempty"))
	}
	data += "} \n"
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntityDataStruct(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("type %sData struct { \n", entity.TypeName)
	data += fmt.Sprintf("	_ struct{} `%s` \n", createTag("typename", entity.OriginalTypeName))
	for _, k := range entity.Properties {
		data += fmt.Sprintf("	%s *%s `%s` \n", k.Name, k.Type, createTag("json", k.OriginalName+",omitempty"))
	}
	data += "} \n"
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntityCommonFunctions(w io.Writer, entity Entity) error {
	if err := g.writeEntityDataFunc(w, entity); err != nil {
		return err
	}
	if err := g.writeEntitySetClientFunc(w, entity); err != nil {
		return err
	}
	return g.writeEntityKeyFunc(w, entity)
}

func (g *Generator) writeEntityKeyFunc(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("func (e *%s) Key__() interface{} { \n", entity.TypeName)
	data += fmt.Sprintf("	return e.%sKey \n", entity.TypeName)
	data += "} \n"
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntityDataFunc(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("func (e *%s) Data__() interface{} { \n", entity.TypeName)
	data += fmt.Sprintf("	return e.%sData \n", entity.TypeName)
	data += "} \n"
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntitySetClientFunc(w io.Writer, entity Entity) error {
	data := fmt.Sprintf("func (e *%s) SetClient__(c *client.Client) { \n", entity.TypeName)
	data += fmt.Sprintf("	e.Client = c \n")
	data += "} \n"
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntityActions(w io.Writer, entity Entity) error {
	data := ""

	for _, f := range entity.Actions {
		params := ""
		for _, p := range f.Parameters {
			params += fmt.Sprintf("%s %s,", p.Name, p.Type)
		}

		params = params[:len(params)-1]

		if f.Type != "" {
			data += fmt.Sprintf("func (e *%s) %s(%s) (result %s, err error) {  \n", entity.TypeName, f.Name, params, f.Type)
		} else {
			data += fmt.Sprintf("func (e *%s) %s(%s) (err error) {  \n", entity.TypeName, f.Name, params)
		}

		data += "	type params struct { \n"
		for _, p := range f.Parameters {
			data += fmt.Sprintf("		%s %s `%s` \n", p.Name, p.Type, createTag("json", p.OriginalName+",omitempty"))
		}
		data += "	} \n"

		params = ""
		for _, p := range f.Parameters {
			params += fmt.Sprintf("%s: %s,", p.Name, p.Name)
		}

		data += fmt.Sprintf("	args := params {%s} \n", params)
		if f.Type != "" {
			data += fmt.Sprintf(`	err = e.Client.ExecuteMethod(e, "%s", args, result)`, f.OriginalName)
		} else {
			data += fmt.Sprintf(`	err = e.Client.ExecuteMethod(e, "%s", args, nil)`, f.OriginalName)
		}
		data += "\n"
		data += "	return \n"
		data += "} \n"
	}
	_, err := w.Write([]byte(data))
	return err
}

func (g *Generator) writeEntityNavigations(w io.Writer, entity Entity) error {
	data := ""

	for _, f := range entity.Navigations {
		data += fmt.Sprintf("func (e *%s) %s() (result %s, err error) {  \n", entity.TypeName, f.Name, f.Type)

		data += fmt.Sprintf(`	err = e.Client.GetNavigation(e, "%s", &result)`, f.OriginalName)
		data += "\n"
		data += "	return \n"
		data += "} \n"
	}
	_, err := w.Write([]byte(data))
	return err
}
