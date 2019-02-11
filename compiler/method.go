package compiler

import (
	"fmt"

	"github.com/actgardner/gogen-avro/schema"
	"github.com/actgardner/gogen-avro/vm"
)

type IRMethod struct {
	name    string
	offset  int
	body    []IRInstruction
	program *IRProgram
}

func NewIRMethod(name string, program *IRProgram) *IRMethod {
	return &IRMethod{
		name:    name,
		body:    make([]IRInstruction, 0),
		program: program,
	}
}

func (p *IRMethod) addLiteral(op vm.Op, t vm.Type, f int) {
	p.body = append(p.body, &LiteralIRInstruction{vm.Instruction{op, t, f}})
}

func (p *IRMethod) addMethodCall(method string) {
	p.body = append(p.body, &MethodCallIRInstruction{method})
}

func (p *IRMethod) addBlockStart() int {
	id := len(p.program.blocks)
	p.program.blocks = append(p.program.blocks, &IRBlock{})
	p.body = append(p.body, &BlockStartIRInstruction{id})
	return id
}

func (p *IRMethod) addBlockEnd(id int) {
	p.body = append(p.body, &BlockEndIRInstruction{id})
}

func (p *IRMethod) VMLength() int {
	len := 0
	for _, inst := range p.body {
		len += inst.VMLength()
	}
	return len
}

func (p *IRMethod) compileType(writer, reader schema.AvroType) error {
	log("compileType()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	switch writer.(type) {
	case *schema.Reference:
		if readerRef, ok := reader.(*schema.Reference); ok || reader == nil {
			return p.compileRef(writer.(*schema.Reference), readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.MapField:
		if readerRef, ok := reader.(*schema.MapField); ok || reader == nil {
			return p.compileMap(writer.(*schema.MapField), readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.ArrayField:
		if readerRef, ok := reader.(*schema.ArrayField); ok || reader == nil {
			return p.compileArray(writer.(*schema.ArrayField), readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.UnionField:
		return p.compileUnion(writer.(*schema.UnionField), reader)
	case *schema.IntField:
		p.addLiteral(vm.Read, vm.Int, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Int, vm.NoopField)
		}
		return nil
	case *schema.LongField:
		p.addLiteral(vm.Read, vm.Long, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Long, vm.NoopField)
		}
		return nil
	case *schema.StringField:
		p.addLiteral(vm.Read, vm.String, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.String, vm.NoopField)
		}
		return nil
	case *schema.BytesField:
		p.addLiteral(vm.Read, vm.Bytes, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Bytes, vm.NoopField)
		}
		return nil
	case *schema.FloatField:
		p.addLiteral(vm.Read, vm.Float, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Float, vm.NoopField)
		}
		return nil
	case *schema.DoubleField:
		p.addLiteral(vm.Read, vm.Double, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Double, vm.NoopField)
		}
		return nil
	case *schema.BoolField:
		p.addLiteral(vm.Read, vm.Boolean, vm.NoopField)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Boolean, vm.NoopField)
		}
		return nil
	case *schema.NullField:
		return nil
	}
	return fmt.Errorf("Unsupported type: %t", writer)
}

func (p *IRMethod) compileRef(writer, reader *schema.Reference) error {
	log("compileRef()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	if reader != nil && writer.TypeName != reader.TypeName {
		return fmt.Errorf("Incompatible types by name: %v %v", reader, writer)
	}

	switch writer.Def.(type) {
	case *schema.RecordDefinition:
		var readerDef *schema.RecordDefinition
		var ok bool
		recordMethodName := fmt.Sprintf("record-r-%v", writer.Def.Name())
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.RecordDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
			recordMethodName = fmt.Sprintf("record-rw-%v", writer.Def.Name())
		}

		if _, ok := p.program.methods[recordMethodName]; !ok {
			method := p.program.createMethod(recordMethodName)
			err := method.compileRecord(writer.Def.(*schema.RecordDefinition), readerDef)
			if err != nil {
				return err
			}
		}
		p.addMethodCall(recordMethodName)
		return nil
	case *schema.FixedDefinition:
		var readerDef *schema.FixedDefinition
		var ok bool
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.FixedDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
		}
		return p.compileFixed(writer.Def.(*schema.FixedDefinition), readerDef)
	case *schema.EnumDefinition:
		var readerDef *schema.EnumDefinition
		var ok bool
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.EnumDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
		}
		return p.compileEnum(writer.Def.(*schema.EnumDefinition), readerDef)
	}
	return fmt.Errorf("Unsupported reference type %t", reader)
}

func (p *IRMethod) compileMap(writer, reader *schema.MapField) error {
	log("compileMap()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	blockId := p.addBlockStart()
	p.addLiteral(vm.Read, vm.MapKey, vm.NoopField)
	var readerType schema.AvroType
	if reader != nil {
		p.addLiteral(vm.AppendMap, vm.Unused, vm.NoopField)
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType)
	if err != nil {
		return err
	}
	if reader != nil {
		p.addLiteral(vm.Exit, vm.Unused, vm.NoopField)
	}
	p.addBlockEnd(blockId)
	return nil
}

func (p *IRMethod) compileArray(writer, reader *schema.ArrayField) error {
	log("compileArray()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	blockId := p.addBlockStart()
	var readerType schema.AvroType
	if reader != nil {
		p.addLiteral(vm.AppendArray, vm.Unused, vm.NoopField)
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType)
	if err != nil {
		return err
	}
	if reader != nil {
		p.addLiteral(vm.Exit, vm.Unused, vm.NoopField)
	}
	p.addBlockEnd(blockId)
	return nil
}

func (p *IRMethod) compileRecord(writer, reader *schema.RecordDefinition) error {
	// Look up whether there's a corresonding target field and if so, parse the source field into that target
	log("compileRecord()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	for _, field := range writer.Fields() {
		var readerType schema.AvroType
		var readerField *schema.Field
		if reader != nil {
			readerField = reader.FieldByName(field.Name())
			if readerField != nil {
				readerType = readerField.Type()
				p.addLiteral(vm.Enter, vm.Unused, readerField.Index())
			}
		}
		err := p.compileType(field.Type(), readerType)
		if err != nil {
			return err
		}
		if readerField != nil {
			p.addLiteral(vm.Exit, vm.Unused, vm.NoopField)
		}
	}
	return nil
}

func (p *IRMethod) compileEnum(writer, reader *schema.EnumDefinition) error {
	log("compileEnum()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	p.addLiteral(vm.Read, vm.Int, vm.NoopField)
	if reader != nil {
		p.addLiteral(vm.Set, vm.Int, vm.NoopField)
	}
	return nil
}

func (p *IRMethod) compileFixed(writer, reader *schema.FixedDefinition) error {
	log("compileFixed()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	p.addLiteral(vm.Read, vm.Fixed, writer.SizeBytes())
	if reader != nil {
		p.addLiteral(vm.Set, vm.Bytes, vm.NoopField)
	}
	return nil
}

func (p *IRMethod) compileUnion(writer *schema.UnionField, reader schema.AvroType) error {
	log("compileUnion()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)

	p.addLiteral(vm.Read, vm.UnionElem, vm.NoopField)
	if _, ok := reader.(*schema.UnionField); ok {
		p.addLiteral(vm.Set, vm.UnionElem, vm.NoopField)
	}
	p.addLiteral(vm.SwitchStart, vm.Unused, vm.NoopField)
writer:
	for i, t := range writer.AvroTypes() {
		p.addLiteral(vm.SwitchCase, vm.Unused, i)
		if unionReader, ok := reader.(*schema.UnionField); ok {
			// If there's an exact match between the reader and writer preserve type
			// This avoids weird cases like ["string", "bytes"] which would always resolve to "string"
			if unionReader.Equals(unionReader) {
				p.addLiteral(vm.Enter, vm.Unused, i)
				err := p.compileType(t, writer.AvroTypes()[i])
				if err != nil {
					return err
				}
				p.addLiteral(vm.Exit, vm.Unused, vm.NoopField)
				continue writer
			}
			for readerIndex, r := range unionReader.AvroTypes() {
				if t.IsReadableBy(r) {
					p.addLiteral(vm.Enter, vm.Unused, readerIndex)
					err := p.compileType(t, r)
					if err != nil {
						return err
					}
					p.addLiteral(vm.Exit, vm.Unused, vm.NoopField)
					continue writer
				}
			}
			return fmt.Errorf("Incompatible types, no match for %v in %v", unionReader, writer)
		} else if t.IsReadableBy(reader) {
			err := p.compileType(t, reader)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Incompatible types: %v %v", reader, writer)
		}
	}
	p.addLiteral(vm.SwitchEnd, vm.Unused, vm.NoopField)
	return nil
}