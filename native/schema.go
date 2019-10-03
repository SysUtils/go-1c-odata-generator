package native

import (
	"github.com/SysUtils/go-1c-odata-generator/shared"
	"github.com/pkg/errors"
)

type Schema struct {
	Entities  []Entity
	Complexes []ComplexType
}

type Entity struct {
	ComplexType
	KeyProperties []Property
}

type Property struct {
	Name         string
	OriginalName string
	Type         string
	OriginalType string
}

type ComplexType struct {
	TypeName         string
	OriginalTypeName string
	Properties       []Property
	Navigations      []Navigation
}

type Navigation struct {
	Name         string
	OriginalName string
	Type         string
	OriginalType string
}

func (g *Generator) extractAssociations(source []shared.Association) map[string]map[string]string {
	associationMap := map[string]map[string]string{}
	for _, assoc := range source {
		name := "StandardODATA." + assoc.Name
		if _, ok := associationMap[name]; !ok {
			associationMap[name] = make(map[string]string, len(assoc.Ends))
		}
		for _, end := range assoc.Ends {
			associationMap[name][end.Role] = end.Type
		}
	}
	return associationMap
}

func (g *Generator) getComplexType(src shared.OneCType, associations map[string]map[string]string) (*ComplexType, error) {
	converted := ComplexType{}
	converted.TypeName = g.translateType(src.Name)
	converted.OriginalTypeName = src.Name

	keyProps := map[string]struct{}{}
	for _, k := range src.Keys {
		keyProps[k.Name] = struct{}{}
	}
	for _, p := range src.Properties {
		prop := Property{
			Name:         g.translateName(p.Name),
			OriginalName: p.Name,
			Type:         g.translateType(p.Type),
			OriginalType: p.Type,
		}

		converted.Properties = append(converted.Properties, prop)
	}

	for _, nav := range src.Navigations {
		navprop := Property{
			Name:         g.translateName(nav.Name),
			OriginalName: nav.Name,
		}
		if relation, ok := associations[nav.Relationship]; ok {
			if to, ok := relation[nav.ToRole]; ok {
				navprop.OriginalType = to
				navprop.Type = g.translateType(to)
				continue
			}
		}
		return nil, errors.Errorf(relationshipNotFound, nav.Relationship, nav.ToRole)
	}
	return &converted, nil
}

func removeProperty(s []Property, i int) []Property {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (g *Generator) convertSchema(src shared.Schema) (*Schema, error) {
	convertedSchema := Schema{}
	associaions := g.extractAssociations(src.Association)

	for _, entity := range src.Entities {
		complexType, err := g.getComplexType(entity, associaions)
		if err != nil {
			return nil, err
		}
		convertedEntity := Entity{
			ComplexType: *complexType,
		}

		keyProps := map[string]struct{}{}
		for _, k := range entity.Keys {
			keyProps[k.Name] = struct{}{}
		}

		// Move KeyProps
		for i := len(convertedEntity.Properties) - 1; i >= 0; i-- {
			if _, ok := keyProps[convertedEntity.Properties[i].OriginalName]; ok {
				convertedEntity.KeyProperties = append(convertedEntity.KeyProperties, convertedEntity.Properties[i])
				convertedEntity.Properties = removeProperty(convertedEntity.Properties, i)
			}
		}

		convertedSchema.Entities = append(convertedSchema.Entities, convertedEntity)
	}
	for _, entity := range src.Complexes {
		complexType, err := g.getComplexType(entity, associaions)
		if err != nil {
			return nil, err
		}
		convertedSchema.Complexes = append(convertedSchema.Complexes, *complexType)
	}

	return &convertedSchema, nil
}
