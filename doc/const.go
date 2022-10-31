package doc

const (
	arrayType   = "array"
	objectType  = "object"
	booleanType = "boolean"
	integerType = "integer"
	numberType  = "number"
	stringType  = "string"

	timeFormat = "RFC3339"
)

var structFieldTypeMap = map[string]specField{
	"string":    {baseType: stringType},
	"int":       {baseType: integerType},
	"float32":   {baseType: numberType},
	"float64":   {baseType: numberType},
	"bool":      {baseType: booleanType},
	"time.Time": {baseType: stringType, format: timeFormat},
}
