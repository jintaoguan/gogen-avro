package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
	"fmt"
)

/*
  A named Reference to a user-defined type (fixed, enum, record). Just a wrapper with a name around a Definition.
*/

type Reference struct {
	name     string
	typeName QualifiedName
	def      Definition
}

func (s *Reference) Name() string {
	return generator.ToPublicName(s.name)
}

func (s *Reference) GoType() string {
	return s.def.GoType()
}

func (s *Reference) FieldType() string {
	return s.def.FieldType()
}

func (s *Reference) SerializerMethod() string {
	return s.def.SerializerMethod()
}

func (s *Reference) DeserializerMethod() string {
	return s.def.DeserializerMethod()
}

func (s *Reference) AddStruct(p *generator.Package) {
	s.def.AddStruct(p)
}

func (s *Reference) AddSerializer(p *generator.Package) {
	s.def.AddSerializer(p)
}

func (s *Reference) AddDeserializer(p *generator.Package) {
	s.def.AddDeserializer(p)
}

func (s *Reference) ResolveReferences(n *Namespace) error {
	if s.def == nil {
		var ok bool
		if s.def, ok = n.Definitions[s.typeName]; !ok {
			return fmt.Errorf("Unable to resolve definition of type %v", s.typeName)
		}
		return s.def.ResolveReferences(n)
	}
	return nil
}