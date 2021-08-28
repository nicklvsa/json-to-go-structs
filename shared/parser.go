package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type JSONStrExpr string
type StructsType map[string]JSONStruct
type MethodsType map[string]JSONMethod

type JSONStructField struct {
	Type string              `json:"type"`
	JSON bool                `json:"json"`
	Tags map[string][]string `json:"tags"`
}

type JSONStruct struct {
	Fields map[string]JSONStructField `json:"fields"`
}

type JSONFuncArg struct {
	Type JSONStrExpr `json:"type"`
}

type GoOut struct {
	Name string `json:"name"`
	Contents string `json:"contents"`
}

type JSONMethodRef struct {
	Name string `json:"name"`
	Mutable bool `json:"mutable"`
}

type JSONMethod struct {
	Ref  JSONMethodRef      `json:"ref"`
	Args map[string]JSONFuncArg `json:"args"`
	Returns []JSONStrExpr `json:"returns"`
}

type ParserDoc struct {
	inputPath string                `json:"-"`
	outputCode GoOut 				`json:"-"`
	Structs   StructsType `json:"structs"`
	Methods MethodsType   `json:"methods"`
}

func NewParserDoc(inputPath string) *ParserDoc {
	outputFileName := fmt.Sprintf("%s.go", strings.TrimSpace(strings.Split(inputPath, ".")[0]))

	return &ParserDoc{
		Structs:   make(StructsType),
		Methods:   make(MethodsType),
		inputPath: inputPath,
		outputCode: GoOut{
			Name: outputFileName,
		},
	}
}

func (st StructsType) GetFieldByPath(path string) (*JSONStructField, error) {
	tree := strings.Split(path, ".")

	if len(tree) <= 1 {
		if _, ok := st[path]; ok {
			return &JSONStructField{
				Type: "interface{}",
			}, nil
		}
	}

	if st, ok := st[tree[0]]; ok {
		if field, ok := st.Fields[tree[1]]; ok {
			return &field, nil
		}
	}

	return nil, fmt.Errorf("unable to find struct field")
}

func (e *JSONStrExpr) Parse(parser *ParserDoc) error {
	for name, handler := range GetSpecialMethods(parser) {
		if strings.HasPrefix(e.Str(), fmt.Sprintf("%s:", name)) {
			val, err := handler(e)
			if err != nil {
				return err
			}

			*e = JSONStrExpr(val)
		}
	}

	return nil
}

func (e JSONStrExpr) Str() string {
	return string(e)
}

func (p *ParserDoc) Parse() error {
	data, err := ioutil.ReadFile(p.inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, p); err != nil {
		return err
	}

	//fmt.Printf("PARSER_DOC: %+v", *p)

	for name, structure := range p.Structs {
		structStart := fmt.Sprintf(`type %s struct {`, name) + "\n"

		for field, fieldStructure := range structure.Fields {
			var tags string

			if fieldStructure.JSON {
				tags = "`"

				tags += fmt.Sprintf(`json:"%s" `, ToSnake(field))
			}

			if len(fieldStructure.Tags) > 0 {
				if !strings.HasPrefix(tags, "`") {
					tags = "`"
				}

				for tag, tagStructure := range fieldStructure.Tags {
					tags += fmt.Sprintf(`%s:"%s" `, tag, strings.Join(tagStructure, ";"))
				}
			}

			if len(tags) > 0 {
				tags += "`"
			}

			structStart += fmt.Sprintf("%s %s %s", field, fieldStructure.Type, strings.TrimSpace(tags)) + "\n"
		}

		p.outputCode.Contents += structStart + "}\n"
	}

	for method, methodStructure := range p.Methods {
		refName := fmt.Sprintf("*%s", methodStructure.Ref.Name)
		if !methodStructure.Ref.Mutable {
			refName = strings.TrimSpace(strings.ReplaceAll(refName, "*", ""))
		}

		var args []string
		for arg, argStructure := range methodStructure.Args {
			if err := argStructure.Type.Parse(p); err != nil {
				return err
			}

			args = append(args, fmt.Sprintf("%s %s", arg, argStructure.Type.Str()))
		}

		var returns []string
		for _, ret := range methodStructure.Returns {
			if err := ret.Parse(p); err != nil {
				return err
			}

			returns = append(returns, ret.Str())
		}

		methodStart := fmt.Sprintf(`func (%s %s) %s(%s) %s {}`, string(methodStructure.Ref.Name[0]), refName, method, strings.Join(args, ","), strings.Join(returns, ","))
		p.outputCode.Contents += methodStart + "\n"
	}


	if err := ioutil.WriteFile(p.outputCode.Name, []byte(p.outputCode.Contents), 0644); err != nil {
		return err
	}

	return nil
}