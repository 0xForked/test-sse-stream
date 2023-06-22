package web

import (
	"embed"
)

//go:embed all:index.html
var UIResource embed.FS
