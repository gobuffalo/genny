package transformer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Transformer_AddImports(t *testing.T) {
	r := require.New(t)

	tcases := []struct {
		goldensPrefix string
		adition       string
	}{
		{"imports-1", `import "github.com/wawandco/fako"`},
		{"imports-2", "import \"github.com/wawandco/fako\"\nimport \"other/package\""},
	}

	for _, tcase := range tcases {

		expected, result, err := matchesAfter(tcase.goldensPrefix, func(tr *Transformer) {
			r.NoError(tr.AddImports(tcase.adition))
		})

		r.NoError(err)
		r.Equal(expected, result)
	}

}

func Test_Transformer_Append(t *testing.T) {
	r := require.New(t)

	tcases := []struct {
		goldensPrefix string
		source        string
	}{
		{"append-1", `//Adding comment at the bottom`},
		{"append-2", "func other() {\n\t//does something else\n}"},
	}

	for _, tcase := range tcases {
		expected, result, err := matchesAfter(tcase.goldensPrefix, func(tr *Transformer) {
			r.NoError(tr.Append(tcase.source))
		})

		r.NoError(err)
		r.Equal(expected, result)
	}

}

func Test_Transformer_AppendAfterSource(t *testing.T) {
	r := require.New(t)

	tcases := []struct {
		goldensPrefix string
		after         string
		source        string
	}{
		{"append-after-1", `c.Flash().Add("success", "Welcome to Buffalo!")`, `c.Flash().Add("warning", "[Warning] Once you get into it, will be hard to leave!")`},
		{"append-after-2", `c.Flash().Add("success", "Welcome to Buffalo!")`, "c.Flash().Add(\"warning\", \"[Warning] Once you get into it, will be hard to leave!\")\nc.Logger().Info(\"Hello\")"},
	}

	for _, tcase := range tcases {
		expected, result, err := matchesAfter(tcase.goldensPrefix, func(tr *Transformer) {
			r.NoError(tr.AppendAfter(tcase.after, tcase.source))
		})

		r.NoError(err)
		r.Equal(expected, result)
	}

}

func Test_Transformer_RemoveBlock(t *testing.T) {
	r := require.New(t)

	tcases := []struct {
		goldensPrefix string
		blockStart    string
	}{
		{"remove-block-1", `func UsersNew(c buffalo.Context) error {`},
		{"remove-block-2", `func DontHaveIt(c buffalo.Context) error {`},
	}

	for _, tcase := range tcases {
		expected, result, err := matchesAfter(tcase.goldensPrefix, func(tr *Transformer) {
			r.NoError(tr.RemoveBlock(tcase.blockStart))
		})

		r.NoError(err)
		r.Equal(expected, result)
	}

}

func Test_Transformer_SetBlockBody(t *testing.T) {
	r := require.New(t)

	tcases := []struct {
		goldensPrefix string
		blockStart    string
		content       string
	}{
		{"insert-in-block-1", `func UsersNew(c buffalo.Context) error {`, "fmt.Println(\"Hello\")"},
	}

	for _, tcase := range tcases {
		expected, result, err := matchesAfter(tcase.goldensPrefix, func(tr *Transformer) {
			r.NoError(tr.SetBlockBody(tcase.blockStart, tcase.content))
		})

		r.NoError(err)
		r.Equal(expected, result)
	}

}

func Test_Transformer_AppendToBlock(t *testing.T) {
	r := require.New(t)

	tcases := []struct {
		goldensPrefix string
		blockStart    string
		content       string
	}{
		{"append-to-block-1", `type Admin struct {`, "Email string"},
	}

	for _, tcase := range tcases {
		expected, result, err := matchesAfter(tcase.goldensPrefix, func(tr *Transformer) {
			r.NoError(tr.AppendToBlock(tcase.blockStart, tcase.content))
		})

		r.NoError(err)
		r.Equal(expected, result)
	}

}

func matchesAfter(prefix string, fn func(tr *Transformer)) (string, string, error) {
	base, err := ioutil.ReadFile(filepath.Join("testdata", prefix+"-in.golden"))
	if err != nil {
		return "", "", err
	}

	tmp := os.TempDir()
	path := filepath.Join(tmp, "file.go")
	ioutil.WriteFile(path, []byte(base), 0644)

	tr := NewTransformer(path)
	fn(tr)

	src, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	expected, err := ioutil.ReadFile(filepath.Join("testdata", prefix+"-out.golden"))
	if err != nil {
		return "", "", err
	}

	return string(src), string(expected), nil
}
