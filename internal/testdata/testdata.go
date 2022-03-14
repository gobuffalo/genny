package testdata

import (
	"embed"
	"io/fs"
)

//go:embed foo.txt bar/* sky/*
var testdata embed.FS

func Data() fs.FS {
	return testdata
}

//go:embed a.html b.html
var boxdata embed.FS

func BoxData() fs.FS {
	return boxdata
}
