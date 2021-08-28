package shared

import "strings"

type SpecialFunc = func(*JSONStrExpr) (string, error)

func (p *ParserDoc) TypeOfHandler(expr *JSONStrExpr) (string, error) {
	str := strings.TrimSpace(strings.ReplaceAll(expr.Str(), "typeof:", ""))

	field, err := p.Structs.GetFieldByPath(str)
	if err != nil {
		return "", nil
	}

	return field.Type, nil
}

func (p *ParserDoc) ParentOfHandler(expr *JSONStrExpr) (string, error) {
	//str := strings.TrimSpace(strings.ReplaceAll(expr.Str(), "parentof:", ""))

	return "", nil
}

func (p *ParserDoc) ChildrenOfHandler(expr *JSONStrExpr) (string, error) {
	//str := strings.TrimSpace(strings.ReplaceAll(expr.Str(), "childrenof:", ""))

	return "", nil
}

func GetSpecialMethods(p *ParserDoc) map[string]SpecialFunc {
	return map[string]SpecialFunc{
		"typeof":     p.TypeOfHandler,
		"parentof":   p.ParentOfHandler,
		"childrenof": p.ChildrenOfHandler,
	}
}