package web

import (
	"embed"
	"io/fs"
)

//go:generate npm run build
//go:embed dist/**
var content embed.FS

type extendFS struct {
	fsys   embed.FS
	extend string
}

func (f extendFS) Open(name string) (fs.File, error) {
	if name == "." {
		return f.fsys.Open(f.extend)
	}
	return f.fsys.Open(f.extend + "/" + name)
}

// Web content
var Content fs.FS = extendFS{extend: "dist", fsys: content}
