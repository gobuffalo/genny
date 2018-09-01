package genny

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Disk_Files(t *testing.T) {
	r := require.New(t)
	run := DryRunner(context.Background())
	d := run.Disk
	d.Add(NewFile("foo.txt", nil))
	d.Add(NewFile("bar.txt", nil))

	files := d.Files()
	r.Len(files, 2)
	r.Equal("bar.txt", files[0].Name())
	r.Equal("foo.txt", files[1].Name())
}

func Test_Disk_Remove(t *testing.T) {
	r := require.New(t)
	run := DryRunner(context.Background())
	d := run.Disk
	d.Add(NewFile("foo.txt", nil))
	d.Add(NewFile("bar.txt", nil))
	d.Remove("foo.txt")

	files := d.Files()
	r.Len(files, 1)
	r.Equal("bar.txt", files[0].Name())
}

func Test_Disk_Delete(t *testing.T) {
	r := require.New(t)
	run := DryRunner(context.Background())
	d := run.Disk
	d.Add(NewFile("foo.txt", nil))
	d.Add(NewFile("bar.txt", nil))
	d.Delete("foo.txt")

	files := d.Files()
	r.Len(files, 1)
	r.Equal("bar.txt", files[0].Name())
}

func Test_Disk_Find(t *testing.T) {
	r := require.New(t)

	run := DryRunner(context.Background())
	d := run.Disk
	d.Add(NewFile("foo.txt", nil))
	d.Add(NewFile("foo.txt", nil))

	f, err := d.Find("foo.txt")
	r.NoError(err)
	r.Equal("foo.txt", f.Name())
}

func Test_Disk_Find_FromDisk(t *testing.T) {
	r := require.New(t)

	run := DryRunner(context.Background())

	d := run.Disk
	f, err := d.Find("fixtures/foo.txt")
	r.NoError(err)

	exp, err := ioutil.ReadFile("./fixtures/foo.txt")
	r.NoError(err)

	act, err := ioutil.ReadAll(f)
	r.NoError(err)

	r.Equal(string(exp), string(act))
}

func Test_Disk_FindFile_DoesntExist(t *testing.T) {
	r := require.New(t)

	run := DryRunner(context.Background())

	_, err := run.Disk.Find("idontexist")
	r.Error(err)
}