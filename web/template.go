package web

import _ "embed"

//go:embed template/menuSelect.gtpl
var menuSelect string

//go:embed template/menu.gtpl
var menu string

//go:embed template/index.html
var index string

//go:embed template/listInside.gtpl
var listInside string

//go:embed template/list.gtpl
var list string

//go:embed template/proInside.gtpl
var proInside string

//go:embed template/progress.gtpl
var progress string

//go:embed template/w3.css
var w3 []byte
