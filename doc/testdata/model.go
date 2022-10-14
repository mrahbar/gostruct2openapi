package testdata

import (
	"fmt"
	"time"
)

type TestBaseInterface interface {
}

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

type testStruct0 struct {
	fieldA string
	FieldB string
	FieldC int
	FieldD bool
}

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

type httpHandler struct {
}

func (g *httpHandler) handleRequest() {
	//@title GCPAsset
	var resp struct {
		Assets []*TestStruct4 `json:"assets"`
	}
	fmt.Println(resp)
}
