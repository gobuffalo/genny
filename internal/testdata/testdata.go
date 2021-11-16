package testdata

import (
	"embed"
	"io/fs"
)

//go:embed foo.txt bar/*
var testdata embed.FS

func Data() fs.FS {
	return testdata
}
