package internal

import (
	"fmt"
	"go/types"
)

type TargetStruct struct {
	name       string
	origType   types.Type
	origStruct *types.Struct
}

func (t *TargetStruct) OriginalType() types.Type {
	return t.origType
}

func (t *TargetStruct) OriginalStruct() *types.Struct {
	return t.origStruct
}

func NewTargetStruct(name string, origType types.Type, origStruct *types.Struct) *TargetStruct {
	return &TargetStruct{name: name, origType: origType, origStruct: origStruct}
}

func (t *TargetStruct) Name() string {
	return t.name
}

func (t *TargetStruct) ID() string {
	if t.IsNamedType() {
		obj := t.ToNamedType().Obj()
		return fmt.Sprintf("%s.%s", obj.Pkg().Path(), obj.Name())
	}
	return t.name
}

func (t *TargetStruct) IsNamedType() bool {
	_, ok := t.origType.(*types.Named)
	return ok
}

func (t *TargetStruct) ToNamedType() *types.Named {
	return t.origType.(*types.Named)
}
