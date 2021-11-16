package genny_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packd"
	"github.com/stretchr/testify/require"
)

func Test_Disk_Files(t *testing.T) {
	r := require.New(t)
	run := genny.DryRunner(context.Background())
	d := run.Disk
	d.Add(genny.NewFile("foo.txt", nil))
	d.Add(genny.NewFile("bar.txt", nil))

	files := d.Files()
	r.Len(files, 2)
	r.Equal("bar.txt", files[0].Name())
	r.Equal("foo.txt", files[1].Name())
}

func Test_Disk_Remove(t *testing.T) {
	r := require.New(t)
	run := genny.DryRunner(context.Background())
	d := run.Disk
	d.Add(genny.NewFile("foo.txt", nil))
	d.Add(genny.NewFile("bar.txt", nil))
	d.Remove("foo.txt")

	files := d.Files()
	r.Len(files, 1)
	r.Equal("bar.txt", files[0].Name())
}

func Test_Disk_Delete(t *testing.T) {
	r := require.New(t)
	run := genny.DryRunner(context.Background())
	d := run.Disk
	d.Add(genny.NewFile("foo.txt", nil))
	d.Add(genny.NewFile("bar.txt", nil))
	r.NoError(d.Delete("foo.txt"))

	files := d.Files()
	r.Len(files, 1)
	r.Equal("bar.txt", files[0].Name())
}

func Test_Disk_Find(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	d := run.Disk
	d.Add(genny.NewFile("foo.txt", nil))
	d.Add(genny.NewFile("foo.txt", nil))

	f, err := d.Find("foo.txt")
	r.NoError(err)
	r.Equal("foo.txt", f.Name())
}

func Test_Disk_Find_FromDisk(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())

	d := run.Disk
	f, err := d.Find("internal/testdata/foo.txt")
	r.NoError(err)

	exp, err := ioutil.ReadFile("./internal/testdata/foo.txt")
	r.NoError(err)

	act, err := ioutil.ReadAll(f)
	r.NoError(err)

	r.Equal(string(exp), string(act))
}

func Test_Disk_FindFile_DoesntExist(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())

	_, err := run.Disk.Find("idontexist")
	r.Error(err)
}

func Test_Disk_AddBox(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()

	td := filepath.Join("internal", "testdata")
	err := filepath.Walk(td, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		p := strings.TrimPrefix(path, td+string(filepath.Separator))
		return box.AddBytes(p, b)
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	d := run.Disk
	err = d.AddBox(box)
	r.NoError(err)

	f, err := d.Find("foo.txt")
	r.NoError(err)
	r.Equal("foo.txt", f.Name())

	f, err = d.Find("bar/baz.txt")
	r.NoError(err)
	r.Equal("bar/baz.txt", f.Name())
}
