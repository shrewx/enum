package enum

import (
	_ "embed"
)

// TplEnum 枚举代码生成模板
//go:embed template/enum.tpl
var TplEnum string
