package testdata

//@title Test Other Base Struct
type TestOtherBaseStruct struct {
	//baseFieldB comment
	baseFieldB string
	//BaseFieldB comment
	BaseFieldB string
	//BaseFieldC comment
	BaseFieldC float64
	//BaseFieldD comment
	BaseFieldD bool
}

//@title Test OtherUnderlying Struct
type TestOtherUnderlyingStruct struct {
	//underlyingFieldB comment
	underlyingFieldB string
	//UnderlyingFieldB comment
	UnderlyingFieldB string
	//UnderlyingFieldC comment
	UnderlyingFieldC float32
	//UnderlyingFieldD comment
	UnderlyingFieldD bool
}

//@title Test Other Struct 5
type TestOtherStruct5 struct {
	TestOtherBaseStruct
	//FieldA comment
	FieldA []string `json:"otherFieldA"`
	//FieldB comment
	FieldB TestOtherUnderlyingStruct `json:"otherFieldB"`
	//FieldC comment
	FieldC []int `json:"otherFieldC"`
	//FieldD comment
	FieldD []bool `json:"otherFieldD"`
}
