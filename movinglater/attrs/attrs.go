package attrs

import (
	"strings"

	"github.com/gobuffalo/flect/name"
)

type Attr struct {
	Original   string
	Name       name.Ident
	commonType string
	goType     string
}

func (a Attr) String() string {
	return a.Original
}

func (a Attr) GoType() string {
	if a.goType != "" {
		return a.goType
	}

	switch strings.ToLower(a.commonType) {
	case "text":
		return "string"
	case "timestamp", "datetime", "date", "time":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	case "decimal", "float":
		return "float64"
	case "[]byte", "blob":
		return "[]byte"
	}

	return a.commonType
}

func (a Attr) CommonType() string {
	if a.commonType != "" {
		return a.commonType
	}
	return a.commonType
}

type Attrs []Attr
