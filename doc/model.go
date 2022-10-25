package doc

import (
	"fmt"
	"go/types"
)

type targetType struct {
	name    string
	origObj types.Object
}

func newTargetType(name string, origObj types.Object) *targetType {
	return &targetType{name: name, origObj: origObj}
}

func (t *targetType) isValid() bool {
	return t.origObj != nil && t.origObj.Type() != nil && t.origObj.Type().Underlying() != nil
}

func (t *targetType) isStruct() bool {
	_, ok := t.origObj.Type().Underlying().(*types.Struct)
	return ok
}

func (t *targetType) isNamedType() bool {
	_, ok := t.origObj.Type().(*types.Named)
	return ok
}

func (t *targetType) toStruct() *types.Struct {
	return t.origObj.Type().Underlying().(*types.Struct)
}

func (t *targetType) toType() types.Type {
	return t.origObj.Type()
}

func (t *targetType) toNamedType() *types.Named {
	return t.origObj.Type().(*types.Named)
}

func (t *targetType) toTargetStruct() *targetStruct {
	if !t.isStruct() {
		return nil
	}

	return &targetStruct{
		name:       t.name,
		origType:   t.toType(),
		origStruct: t.toStruct(),
	}
}

type targetStruct struct {
	name       string
	origType   types.Type
	origStruct *types.Struct
}

func newTargetStruct(name string, origType types.Type, origStruct *types.Struct) *targetStruct {
	return &targetStruct{name: name, origType: origType, origStruct: origStruct}
}

func (t *targetStruct) isNamedType() bool {
	_, ok := t.origType.(*types.Named)
	return ok
}

func (t *targetStruct) toNamedType() *types.Named {
	return t.origType.(*types.Named)
}

type targetField struct {
	packageID   string
	structName  string
	fieldTag    string
	fieldName   string
	specField   specField
	elem        types.Type
	isArrayType bool
}

func (t *targetField) ID() string {
	return fmt.Sprintf("%s.%s.%s", t.packageID, t.structName, t.fieldName)
}
