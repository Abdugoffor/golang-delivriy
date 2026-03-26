package main

import (
	"embed"
	"my-project/helper"
)

//go:embed views
var viewsFS embed.FS

func init() {
	helper.ViewFS = viewsFS
}
