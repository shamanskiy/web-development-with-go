package templates

import "embed"

//go:embed *.gohtml users/*.gohtml galleries/*.gohtml
var FS embed.FS
