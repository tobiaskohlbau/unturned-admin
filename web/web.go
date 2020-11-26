package web

import (
	"embed"
	"errors"
	"io/fs"
)

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
	file, err := f.fsys.Open(f.extend + "/" + name)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return f.fsys.Open(f.extend + "/" + "index.html")
	}
	return file, err
}

// Web content
var Content fs.FS = extendFS{extend: "dist", fsys: content}
