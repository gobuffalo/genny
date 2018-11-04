package genny

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/gobuffalo/packd"
	"github.com/pkg/errors"
)

// Disk is a virtual file system that works
// with both dry and wet runners. Perfect for seeding
// Files or non-destructively deleting files
type Disk struct {
	Runner   *Runner
	files    *sync.Map
	original *sync.Map
}

func (d *Disk) AddBox(box packd.Walker) error {
	return box.Walk(func(path string, file packd.File) error {
		d.Add(NewFile(path, file))
		return nil
	})
}

// Files returns a sorted list of all the files in the disk
func (d *Disk) Files() []File {
	var files []File

	d.files.Range(func(k, v interface{}) bool {
		f, ok := v.(File)
		if !ok {
			return false
		}
		if s, ok := f.(io.Seeker); ok {
			s.Seek(0, 0)
		}
		files = append(files, f)
		return true
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	return files
}

func newDisk(r *Runner) *Disk {
	return &Disk{
		Runner:   r,
		files:    &sync.Map{},
		original: &sync.Map{},
	}
}

// Remove a file(s) from the virtual disk.
func (d *Disk) Remove(name string) {
	d.files.Delete(name)
}

// Delete calls the Runner#Delete function
func (d *Disk) Delete(name string) error {
	return d.Runner.Delete(name)
}

// Add file to the virtual disk
func (d *Disk) Add(f File) {
	d.files.Store(f.Name(), f)
}

// Find a file from the virtual disk. If the file doesn't
// exist it will try to read the file from the physical disk.
func (d *Disk) Find(name string) (File, error) {
	ff, ok := d.files.Load(name)
	if ok {
		if f, ok := ff.(File); ok {
			return f, nil
		}
		return nil, errors.Errorf("expected File got %T", ff)
	}

	gf := NewFile(name, bytes.NewReader([]byte("")))

	osname := name
	if runtime.GOOS == "windows" {
		osname = strings.Replace(osname, "/", "\\", -1)
	}
	of, err := os.Open(osname)
	if err != nil {
		return gf, errors.WithStack(err)
	}
	defer of.Close()

	bb := &bytes.Buffer{}

	if _, err := io.Copy(bb, of); err != nil {
		return gf, errors.WithStack(err)
	}
	gf = NewFile(name, bb)
	d.original.LoadOrStore(gf.Name(), gf)
	d.Add(gf)
	return gf, nil
}

func (d *Disk) Rollback() error {
	d.Runner.Logger.Debug("rolling back disk")
	var err error
	d.files.Range(func(k, v interface{}) bool {
		s, ok := k.(string)
		if !ok {
			return false
		}
		if err = d.Delete(s); err != nil {
			return false
		}
		return true
	})

	if err != nil {
		return errors.WithStack(err)
	}

	d.original.Range(func(k, v interface{}) bool {
		f, ok := v.(File)
		if !ok {
			return false
		}
		if err = d.Runner.File(f); err != nil {
			return false
		}
		return true
	})

	return err
}
