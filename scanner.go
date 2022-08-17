package enum

import (
	"fmt"
	"github.com/go-courier/packagesx"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"
)

type EnumScanner struct {
	pkg   *packagesx.Package
	enums map[*types.TypeName][]EnumValue
}

func NewEnumScanner(pkg *packagesx.Package) *EnumScanner {
	return &EnumScanner{
		pkg:   pkg,
		enums: make(map[*types.TypeName][]EnumValue, 0),
	}
}

func (e *EnumScanner) Enum(typeName *types.TypeName) []EnumValue {
	if typeName == nil {
		return nil
	}

	if values, ok := e.enums[typeName]; ok {
		return sortEnumValues(values)
	}
	ev := EnumValue{}
	typeUnderlying := typeName.Type().Underlying().String()
	if strings.Contains(typeUnderlying, "int") {
		ev.IntValue = new(int64)
	} else if strings.Contains(typeUnderlying, "float") {
		ev.FloatValue = new(float64)
	} else if strings.Contains(typeUnderlying, "string") {
		ev.StringValue = new(string)
	} else {
		panic(fmt.Errorf("enum type underlying must be an int or string, but got %s", typeName.Type().Underlying().String()))
	}

	pkgInfo := e.pkg.Pkg(e.pkg.PkgPath)
	if pkgInfo == nil {
		return nil
	}
	for ident, def := range pkgInfo.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type() != typeName.Type() {
			continue
		}

		ev.Key = typeConst.Name()

		if ev.IntValue != nil {
			intValue, _ := strconv.ParseInt(typeConst.Val().String(), 10, 64)
			ev.IntValue = &intValue
			sp := strings.SplitN(ev.Key, "__", 2)
			stringValue := ev.Key
			if len(sp) == 2 {
				stringValue = sp[1]
				ev.StringValue = &stringValue
			}
		} else if ev.FloatValue != nil {
			floatValue, _ := strconv.ParseFloat(typeConst.Val().String(), 64)
			ev.FloatValue = &floatValue
		} else if ev.StringValue != nil {
			stringValue := strings.Trim(typeConst.Val().String(), "\"")
			ev.StringValue = &stringValue
		}

		var label string
		if v, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
			if v.Doc.Text() != "" {
				label = v.Doc.Text()
			}
			if v.Comment.Text() != "" {
				label = v.Comment.Text()
			}
		}
		ev.Label = strings.TrimRight(strings.Replace(label, "\n", " ", -1), " ")
		e.enums[typeName] = append(e.enums[typeName], ev)
	}

	return sortEnumValues(e.enums[typeName])
}

func sortEnumValues(enumValues []EnumValue) []EnumValue {
	sort.Slice(enumValues, func(i, j int) bool {
		if enumValues[i].IntValue == enumValues[j].IntValue {
			return *enumValues[i].StringValue < *enumValues[j].StringValue
		}
		return *enumValues[i].IntValue < *enumValues[j].IntValue
	})

	return enumValues
}
