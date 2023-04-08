package testdata

import (
	"fmt"
	"github.com/mrahbar/gostruct2openapi/testdata"
	"time"
)

type TestBaseInterface interface {
}

// @title Test Base Struct
// Test Base description
type TestBaseStruct struct {
	//baseFieldB comment
	baseFieldB string
	//BaseFieldB comment
	BaseFieldB string
	//BaseFieldC comment
	BaseFieldC float64
	//BaseFieldD comment
	BaseFieldD bool
}

// @title Test Underlying Struct
// Test Underlying Struct description
type TestUnderlyingStruct struct {
	//underlyingFieldB comment
	underlyingFieldB string
	//UnderlyingFieldB comment
	UnderlyingFieldB string
	//UnderlyingFieldC comment
	UnderlyingFieldC float32
	//UnderlyingFieldD comment
	UnderlyingFieldD bool
}

// @title Test Struct 0
// Test Struct 0 description
type testStruct0 struct {
	fieldA string
	FieldB string
	FieldC int
	FieldD bool
}

// @title Test Struct 1
// Test Struct 1 description
type TestStruct1 struct {
	//fieldA comment
	fieldA string
	//FieldB comment
	FieldB string
	//FieldC comment
	FieldC int
	//FieldD comment
	FieldD bool
}

// @title Test Struct 2
// Test Struct 2 description
type TestStruct2 struct {
	//TestBaseStruct comment
	TestBaseStruct
	//fieldA comment
	fieldA []string
	//FieldB comment
	FieldB []string
	//FieldC comment
	FieldC []int
	//FieldD comment
	FieldD []bool
}

// MyString description
type MyString string

// @title Test Struct 3
// Test Struct 3 description
type TestStruct3 struct {
	//TestBaseStruct comment
	TestBaseStruct
	//FieldA comment
	FieldA time.Time
	//FieldB comment
	FieldB TestUnderlyingStruct
	//FieldC comment
	FieldC *TestUnderlyingStruct
	//FieldD comment
	FieldD []TestUnderlyingStruct
	//FieldE comment
	FieldE []*TestUnderlyingStruct
	//FieldF comment
	FieldF TestBaseInterface
	//FieldG comment
	FieldG *time.Time
	//FieldH comment
	FieldH map[string]string
	//FieldI comment
	FieldI chan bool
	//FieldJ comment
	FieldJ MyString
	//FieldK comment
	FieldK map[string]*TestUnderlyingStruct
}

// @title Test Struct 4
// Test Struct 4 description
type TestStruct4 struct {
	//FieldA comment
	FieldA []string `json:"otherFieldA"`
	//FieldB comment
	FieldB []string `json:"otherFieldB"`
	//FieldC comment
	FieldC []int `json:"otherFieldC"`
	//FieldD comment
	FieldD []bool `json:"otherFieldD"`
}

// @title MyAsset
// MyAsset description
type httpHandlerResp struct {
	TestStructs  []*TestStruct4               `json:"structs"`
	OtherStructs []*testdata.TestOtherStruct5 `json:"other_structs"`
}

// @title HTTP Handler
type httpHandler struct {
}

func (g *httpHandler) handleRequest() {
	var resp2 httpHandlerResp
	fmt.Println(resp2)
}
