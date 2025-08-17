package web

import "embed"

//go:embed public
var fs embed.FS

func FS() embed.FS {
	return fs
}
