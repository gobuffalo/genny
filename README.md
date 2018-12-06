<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://godoc.org/github.com/gobuffalo/genny"><img src="https://godoc.org/github.com/gobuffalo/genny?status.svg" alt="GoDoc" /></a>
<a href="https://travis-ci.org/gobuffalo/genny"><img src="https://travis-ci.org/gobuffalo/genny.svg?branch=master" alt="Build Status" /></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/genny"><img src="https://goreportcard.com/badge/github.com/gobuffalo/genny" alt="Go Report Card" /></a>
</p>

# Genny

## What Is Genny?

Genny is a _framework_ for writing modular generators, it however, doesn't actually generate anything. It just makes it easier for you to. :)

## Documentation

For right now the [GoDoc](https://godoc.org/github.com/gobuffalo/genny) and the source/tests are best documentation as the APIs are currently in flux.

## Core Concepts

### Generators

A [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) is used to build a blue print of what you want to generate.

A few of things that can be added to a `Generator` are:

* [`github.com/gobuffalo/genny#File`](https://godoc.org/github.com/gobuffalo/genny#File)
* [`os/exec#Cmd`](https://godoc.org/os/exec#Cmd)
* [`github.com/gobuffalo/packd#Box`](https://godoc.org/github.com/gobuffalo/packd#Box)
* [`net/http#Request`](https://godoc.org/net/http#Request)
* and more

The `Generator` does *not* actually generate anything; a [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) is needed to run the generator.

```go
g := genny.New()
g.File(genny.NewFileS("index.html", "Hello"))
g.Command(exec.Command("go", "env"))
```

### Runners

A [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) is used to run generators and control the environment in which those generators are run.

Genny ships with three implementations of `Runner` that cover _most_ situations. They can also provide good starting points for customized implementations of `Runner`.

* [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner)
* [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner)
* [`github.com/gobuffalo/genny/gentest#NewRunner`](https://godoc.org/github.com/gobuffalo/genny/gentest#NewRunner)

#### Dry Running (**NON-DESTRUCTIVE**)

The idea of "dry" running means that no commands are executed, no files are written to disk, no HTTP requests are made, etc... Instead these steps are run "dry", which in the case of [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner) is the case.

```go
func main() {
	run := genny.DryRunner(context.Background())

	g := genny.New()
	g.File(genny.NewFileS("index.html", "Hello\n"))
	g.Command(exec.Command("go", "env"))
	run.With(g)

	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
```

```plain
// output
DEBU[2018-12-05T11:26:21-05:00] Step: d4da0ef6
DEBU[2018-12-05T11:26:21-05:00] Chdir: /go/src/github.com/gobuffalo/genny/_examples/dry
DEBU[2018-12-05T11:26:21-05:00] File: /go/src/github.com/gobuffalo/genny/_examples/dry/index.html
Hello
DEBU[2018-12-05T11:26:21-05:00] Exec: go env
```

```bash
// file list
.
└── main.go

0 directories, 1 file
```

Using a "dry" runner can make testing easier when you don't have to worry about commands running, files being written, etc... It can also make it easy to provide a "dry-run" flag to your generators to let people see what will be generated when the generator is run for real.

#### Wet Running (**DESTRUCTIVE**)

While "dry" means to not execute commands or write files, "wet" running means the exact opposite; it will write files and execute commands.

Use the [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner) when "wet" running is the desired outcome.

```go
func main() {
	run := genny.WetRunner(context.Background())

	g := genny.New()
	g.File(genny.NewFileS("index.html", "Hello\n"))
	g.Command(exec.Command("go", "env"))
	run.With(g)

	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
```

```plain
// output
GOARCH="amd64"
....
GOGCCFLAGS=""
```

```bash
// file list
.
├── index.html
└── main.go

0 directories, 2 files
```

```bash
$ cat index.html

Hello
```

#### Changing Runner Behavior

The change the way [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner) or [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner) work, or to build your own [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) you need to implement the `*Fn` attributes on the [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner).

```go
type Runner struct {
	// ...
	ExecFn     func(*exec.Cmd) error                                     // function to use when executing files
	FileFn     func(File) (File, error)                                  // function to use when writing files
	ChdirFn    func(string, func() error) error                          // function to use when changing directories
	DeleteFn   func(string) error                                        // function used to delete files/folders
	RequestFn  func(*http.Request, *http.Client) (*http.Response, error) // function used to make http requests
	LookPathFn func(string) (string, error)                              // function used to make exec.LookPath lookups
	// ...
}
```

These `*Fn` functions represent the **FINAL** end-point for the that is trying to be run.

```go
	run.FileFn = func(f packd.SimpleFile) (packd.SimpleFile, error) {
		io.Copy(os.Stdout, f)
		return f, nil
	}

	run.FileFn = func(f genny.File) (genny.File, error) {
		if d, ok := f.(genny.Dir); ok {
			if err := os.MkdirAll(d.Name(), d.Perm); err != nil {
				return f, err
			}
			return d, nil
		}

		name := f.Name()
		if !filepath.IsAbs(name) {
			name = filepath.Join(run.Root, name)
		}
		dir := filepath.Dir(name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return f, err
		}
		ff, err := os.Create(name)
		if err != nil {
			return f, err
		}
		defer ff.Close()
		if _, err := io.Copy(ff, f); err != nil {
			return f, err
		}
		return f, nil
	}
```

#### Runtime Checks

### Files

#### Reading Files

#### Writing Files

#### Transforming Files

### Testing



