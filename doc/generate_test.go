package doc

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_OpenapiGenerator_UnknownPackage(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct1"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/unknown")
	assert.Error(t, err)
	assert.Empty(t, specs)
}

func Test_OpenapiGenerator_Struct0(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("testStruct0"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "testStruct0",
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
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct1"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 1)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "TestStruct1",
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
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct2"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 1)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "TestStruct2",
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
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct3"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 2)

	bytes, err := json.Marshal(specs)
	assert.NoError(t, err)
	assert.JSONEq(t, `[
		{
			"id": "TestStruct3",
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
				"FieldA": {
					"format": "RFC3339",
					"type": "string"
				},
				"FieldB": {
					"$ref": "#/components/schemas/TestUnderlyingStruct",
					"description": "FieldB comment"
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
				}
			}
		},
		{
			"id": "TestUnderlyingStruct",
			"type":"object",
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
			}
		}
	]`, string(bytes))
}

func Test_OpenapiGenerator_Struct4(t *testing.T) {
	generator := NewOpenapiGenerator(regexp.MustCompile("TestStruct4"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 1)

	bytes, err := specs[0].MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "TestStruct4",
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
	generator := NewOpenapiGenerator(regexp.MustCompile("httpHandler|resp"))
	specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
	assert.NoError(t, err)
	assert.Len(t, specs, 3)

	bytes, err := json.Marshal(specs)
	assert.NoError(t, err)
	assert.JSONEq(t, `[
	{
		"id": "httpHandler",
		"type":"object"
	},
	{
		"id": "resp",
		"type":"object",
		"properties": {
			"assets": {
				"items": {
					"$ref": "#/components/schemas/TestStruct4"
				},
				"type": "array"
			}
		}
	},
	{
		"id": "TestStruct4",
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
	}
]`, string(bytes))
}
