// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package avro

import (
	"io"
	
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)

  
type MapOptionalTestRecord struct { 
	
	
		OptField MapUnionNullStringInt
	

}

func DeserializeMapOptionalTestRecord(r io.Reader) (t MapOptionalTestRecord, err error) {
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err == nil {
		err = vm.Eval(r, deser, &t)
	}
	return
}

func DeserializeMapOptionalTestRecordFromSchema(r io.Reader, schema string) (t MapOptionalTestRecord, err error) {
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err == nil {
		err = vm.Eval(r, deser, &t)
	}
	return
}

func writeMapOptionalTestRecord(r MapOptionalTestRecord, w io.Writer) error {
	var err error
	
	err = writeMapUnionNullStringInt( r.OptField, w)
	if err != nil {
		return err			
	}
	
	return err
}

func (r MapOptionalTestRecord) Serialize(w io.Writer) error {
	return writeMapOptionalTestRecord(r, w)
}

func (r MapOptionalTestRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"OptField\",\"type\":{\"type\":\"map\",\"values\":[\"null\",\"string\",\"int\"]}}],\"name\":\"MapOptionalTestRecord\",\"type\":\"record\"}"
}

func (r MapOptionalTestRecord) SchemaName() string {
	return "MapOptionalTestRecord"
}

func (_ *MapOptionalTestRecord) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetInt(v int32) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetLong(v int64) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetString(v string) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *MapOptionalTestRecord) Get(i int) types.Field {
	switch (i) {
	
	case 0:
		
			r.OptField = NewMapUnionNullStringInt()

		
		
			return &r.OptField
		
	
	default:
		panic("Unknown field index")
	}
}

func (r *MapOptionalTestRecord) SetDefault(i int) {
	switch (i) { 
	default:
		panic("Unknown field index")
	}
}

func (r *MapOptionalTestRecord) Clear(i int) {
	switch (i) { 
	default:
		panic("Non-optional field index")
	}
}

func (_ *MapOptionalTestRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) ClearMap(key string) { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *MapOptionalTestRecord) Finalize() { }