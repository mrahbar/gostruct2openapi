# GO struct to OpenAPI spec converter 

This repository can be used to converts GO struct definitions to OpenAPI schema definitions using AST.
By using the AST all struct fields with comments can be parsed. Additionally the library parses also references of custom type in different packages. For parsing struct level comments ``go/doc`` is being used. 

### Config
- To change the property name struct tags can be used e.g. ``json``.
- To change the struct name the comment directive ``@title`` can be used.
- To only generate for a set of struct regular expression can be used to filtger struct names, e.g. ``*HandlerResponse``.

### Use 

To update the library to the latest version, use ``go get -u github.com/stretchr/testify``.

```
generator := NewOpenapiGenerator(regexp.MustCompile("TestBaseStruct"), "json")
specs, err := generator.DocumentStruct("github.com/mrahbar/gostruct2openapi/doc/testdata")
if err != nil {
    log.Fatal(err)
}
//TODO use specs variable, e.g. by writting it to a file
```

### Example

Given the following struct
```
//@title Test Base Struct
//Test Base description
type TestBaseStruct struct {
	//baseFieldB comment
	baseFieldB string
	//BaseFieldB comment
	BaseFieldB string `json:"otherBaseFieldB"`
	//BaseFieldC comment
	BaseFieldC float64
	//BaseFieldD comment
	BaseFieldD bool
}
```

will result in the following output
```
{
    "description": "Test Base description",
    "id": "Test Base Struct",
    "properties": {
        "otherBaseFieldB": {
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
        }
    }
}
```

For more examples check out the test and the corresponding test data


## Additional sources
- https://github.com/dave/dst
- https://github.com/swaggo/swag
- https://stackoverflow.com/questions/19580688/go-parser-not-detecting-doc-comments-on-struct-type