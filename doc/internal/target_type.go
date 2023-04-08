package internal

import (
	"go/types"
)

// TargetType represents
type TargetType struct {
	name    string
	origObj types.Object
}

func NewTargetType(name string, origObj types.Object) *TargetType {
	return &TargetType{name: name, origObj: origObj}
}

func (t *TargetType) IsValid() bool {
	return t.origObj != nil && t.origObj.Type() != nil && t.origObj.Type().Underlying() != nil
}

func (t *TargetType) IsStruct() bool {
	_, ok := t.origObj.Type().Underlying().(*types.Struct)
	return ok
}

func (t *TargetType) IsNamedType() bool {
	_, ok := t.origObj.Type().(*types.Named)
	return ok
}

func (t *TargetType) toStruct() *types.Struct {
	return t.origObj.Type().Underlying().(*types.Struct)
}

func (t *TargetType) toType() types.Type {
	return t.origObj.Type()
}

func (t *TargetType) ToNamedType() *types.Named {
	return t.origObj.Type().(*types.Named)
}

func (t *TargetType) ToTargetStruct() *TargetStruct {
	if !t.IsStruct() {
		return nil
	}

	return &TargetStruct{
		name:       t.name,
		origType:   t.toType(),
		origStruct: t.toStruct(),
	}
}
