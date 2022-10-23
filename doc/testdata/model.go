package testdata

import (
	"fmt"
	"time"
)

type TestBaseInterface interface {
}

//@title Test Base Struct
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

//@title Test Underlying Struct
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

//@title Test Struct 0
type testStruct0 struct {
	fieldA string
	FieldB string
	FieldC int
	FieldD bool
}

//@title Test Struct 1
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

//@title Test Struct 2
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

//@title Test Struct 3
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
}

//@title Test Struct 4
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

//@title HTTP Handler
type httpHandler struct {
}

func (g *httpHandler) handleRequest() {
	//@title MyAsset
	var resp struct {
		Assets []*TestStruct4 `json:"assets"`
	}
	fmt.Println(resp)
}
