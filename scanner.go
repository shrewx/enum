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

	typeUnderlying := typeName.Type().Underlying().String()
	isInt := false
	if strings.Contains(typeUnderlying, "int") {
		isInt = true
	} else if !strings.Contains(typeUnderlying, "string") {
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

		key := typeConst.Name()
		sp := strings.SplitN(key, "__", 2)
		stringValue := key
		if len(sp) == 2 {
			stringValue = sp[1]
		}

		var intValue int64
		if isInt {
			intValue, _ = strconv.ParseInt(typeConst.Val().String(), 10, 32)
		} else {
			stringValue = strings.Trim(typeConst.Val().String(), "\"")
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

		e.enums[typeName] = append(e.enums[typeName], EnumValue{
			Key:         key,
			StringValue: stringValue,
			IntValue:    int(intValue),
			Label:       strings.TrimRight(strings.Replace(label, "\n", " ", -1), " "),
		})
	}

	return sortEnumValues(e.enums[typeName])
}

func sortEnumValues(enumValues []EnumValue) []EnumValue {
	sort.Slice(enumValues, func(i, j int) bool {
		if enumValues[i].IntValue == enumValues[j].IntValue {
			return enumValues[i].StringValue < enumValues[j].StringValue
		}
		return enumValues[i].IntValue < enumValues[j].IntValue
	})

	return enumValues
}
