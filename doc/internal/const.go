package internal

type SpecType string

func (s SpecType) String() string {
	return string(s)
}

const (
	ArrayType   SpecType = "array"
	ObjectType  SpecType = "object"
	StructType  SpecType = "struct"
	BooleanType SpecType = "boolean"
	IntegerType SpecType = "integer"
	NumberType  SpecType = "number"
	StringType  SpecType = "string"

	TimeFormat = "RFC3339"
)

var StructFieldTypeMap = map[string]*SpecField{
	"string":    NewSpecField(StringType),
	"int":       NewSpecField(IntegerType),
	"float32":   NewSpecField(NumberType),
	"float64":   NewSpecField(NumberType),
	"bool":      NewSpecField(BooleanType),
	"time.Time": NewSpecFieldWithFormat(StringType, TimeFormat),
}
