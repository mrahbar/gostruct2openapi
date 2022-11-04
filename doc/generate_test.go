package doc

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_OpenapiGenerator_UnknownPackage(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct1"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/unknown")
	assert.Error(t, err)
	assert.Empty(t, specs)
}

func Test_OpenapiGenerator_Struct0(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("testStruct0"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"description":"Test Struct 0 description",
		"id": "Test Struct 0",
		"type":"object",
		"properties": {
			"FieldB": {
				"type": "string"
			},
			"FieldC": {
				"type": "integer"
			},
			"FieldD": {
				"type": "boolean"
			}
		}
	}`, string(bytes))
}

func Test_OpenapiGenerator_Struct1(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct1"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 1)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"description":"Test Struct 1 description",
		"id": "Test Struct 1",
		"type":"object",
		"properties": {
			"FieldB": {
				"description": "FieldB comment",
				"type": "string"
			},
			"FieldC": {
				"description": "FieldC comment",
				"type": "integer"
			},
			"FieldD": {
				"description": "FieldD comment",
				"type": "boolean"
			}
		}
	}`, string(bytes))
}

func Test_OpenapiGenerator_Struct2(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct2"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 1)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"description":"Test Struct 2 description",
		"id": "Test Struct 2",
		"type":"object",
		"properties": {
			"BaseFieldB": {
				"description": "BaseFieldB comment",
				"type": "string"
			},
			"BaseFieldC": {
				"description": "BaseFieldC comment",
				"type": "number"
			},
			"BaseFieldD": {
				"description": "BaseFieldD comment",
				"type": "boolean"
			},
			"FieldB": {
				"description": "FieldB comment",
				"items": {
					"type": "string"
				},
				"type": "array"
			},
			"FieldC": {
				"description": "FieldC comment",
				"items": {
					"type": "integer"
				},
				"type": "array"
			},
			"FieldD": {
				"description": "FieldD comment",
				"items": {
					"type": "boolean"
				},
				"type": "array"
			}
		}
	}`, string(bytes))
}

func Test_OpenapiGenerator_Struct3(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct3"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 2)

	bytes, err := json.Marshal(specs)
	assert.NoError(t, err)
	assert.JSONEq(t, `[
		{
			"description":"Test Struct 3 description",
			"id": "Test Struct 3",
			"properties": {
				"BaseFieldB": {
					"description": "BaseFieldB comment",
					"type": "string"
				},
				"BaseFieldC": {
					"description": "BaseFieldC comment",
					"type": "number"
				},
				"BaseFieldD": {
					"description": "BaseFieldD comment",
					"type": "boolean"
				},
				"FieldA": {
					"format": "RFC3339",
					"type": "string"
				},
				"FieldB": {
					"$ref": "#/components/schemas/TestUnderlyingStruct"
				},
				"FieldC": {
					"$ref": "#/components/schemas/TestUnderlyingStruct"
				},
				"FieldD": {
					"description": "FieldD comment",
					"items": {
						"$ref": "#/components/schemas/TestUnderlyingStruct"
					},
					"type": "array"
				},
				"FieldE": {
					"description": "FieldE comment",
					"items": {
						"$ref": "#/components/schemas/TestUnderlyingStruct"
					},
					"type": "array"
				},
				"FieldF": {
					"description": "FieldF comment",
					"type": "object"
				},
				"FieldG": {
					"description": "FieldG comment",
					"format": "RFC3339",
					"type": "string"
				},
				"FieldH": {
					"description": "FieldH comment",
					"additionalProperties": {
						"type": "string"
					},
					"type": "object"
				},
				"FieldI": {
					"description": "FieldI comment",
					"type": "object"
				},
				"FieldJ": {
					"description": "FieldJ comment",
					"type": "string"
				},
				"FieldK": {
					"description": "FieldK comment",
					"additionalProperties": {
						"$ref": "#/components/schemas/TestUnderlyingStruct"
					},
					"type": "object"
				}
			},
			"type":"object"
		},
		{
			"description":"Test Underlying Struct description",
			"id": "Test Underlying Struct",
			"properties": {
				"UnderlyingFieldB": {
					"description": "UnderlyingFieldB comment",
					"type": "string"
				},
				"UnderlyingFieldC": {
					"description": "UnderlyingFieldC comment",
					"type": "number"
				},
				"UnderlyingFieldD": {
					"description": "UnderlyingFieldD comment",
					"type": "boolean"
				}
			},
			"type":"object"
		}
	]`, string(bytes))
}

func Test_OpenapiGenerator_Struct4(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct4"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 1)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"description":"Test Struct 4 description",
		"id": "Test Struct 4",
		"type":"object",
		"properties": {
			"otherFieldA": {
				"description": "FieldA comment",
				"items": {
					"type": "string"
				},
				"type": "array"
			},
			"otherFieldB": {
				"description": "FieldB comment",
				"items": {
					"type": "string"
				},
				"type": "array"
			},
			"otherFieldC": {
				"description": "FieldC comment",
				"items": {
					"type": "integer"
				},
				"type": "array"
			},
			"otherFieldD": {
				"description": "FieldD comment",
				"items": {
					"type": "boolean"
				},
				"type": "array"
			}
		}
	}`, string(bytes))
}

func Test_OpenapiGenerator_Method(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("httpHandler|resp"), "json")
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 5)

	bytes, err := json.Marshal(specs)
	assert.NoError(t, err)
	assert.JSONEq(t, `[
		{
			"id": "HTTP Handler",
			"type": "object"
		},
		{
			"description": "MyAsset description",
			"id": "MyAsset",
			"properties": {
				"other_structs": {
					"items": {
						"$ref": "#/components/schemas/TestOtherStruct5"
					},
					"type": "array"
				},
				"structs": {
					"items": {
						"$ref": "#/components/schemas/TestStruct4"
					},
					"type": "array"
				}
			},
			"type": "object"
		},
		{
			"description": "Test Other Struct 5 description",
			"id": "Test Other Struct 5",
			"properties": {
				"BaseFieldB": {
					"description": "BaseFieldB comment",
					"type": "string"
				},
				"BaseFieldC": {
					"description": "BaseFieldC comment",
					"type": "number"
				},
				"BaseFieldD": {
					"description": "BaseFieldD comment",
					"type": "boolean"
				},
				"otherFieldA": {
					"description": "FieldA comment",
					"items": {
						"type": "string"
					},
					"type": "array"
				},
				"otherFieldB": {
					"$ref": "#/components/schemas/TestOtherUnderlyingStruct"
				},
				"otherFieldC": {
					"description": "FieldC comment",
					"items": {
						"type": "integer"
					},
					"type": "array"
				},
				"otherFieldD": {
					"description": "FieldD comment",
					"items": {
						"type": "boolean"
					},
					"type": "array"
				}
			},
			"type": "object"
		},
		{
			"description": "Test OtherUnderlying description",
			"id": "Test OtherUnderlying Struct",
			"properties": {
				"UnderlyingFieldB": {
					"description": "UnderlyingFieldB comment",
					"type": "string"
				},
				"UnderlyingFieldC": {
					"description": "UnderlyingFieldC comment",
					"type": "number"
				},
				"UnderlyingFieldD": {
					"description": "UnderlyingFieldD comment",
					"type": "boolean"
				}
			},
			"type": "object"
		},
		{
			"description": "Test Struct 4 description",
			"id": "Test Struct 4",
			"properties": {
				"otherFieldA": {
					"description": "FieldA comment",
					"items": {
						"type": "string"
					},
					"type": "array"
				},
				"otherFieldB": {
					"description": "FieldB comment",
					"items": {
						"type": "string"
					},
					"type": "array"
				},
				"otherFieldC": {
					"description": "FieldC comment",
					"items": {
						"type": "integer"
					},
					"type": "array"
				},
				"otherFieldD": {
					"description": "FieldD comment",
					"items": {
						"type": "boolean"
					},
					"type": "array"
				}
			},
			"type": "object"
		}
	]
`, string(bytes))
}
