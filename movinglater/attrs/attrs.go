package attrs

import (
	"github.com/markbates/inflect"
)

type Attr struct {
	Original   string
	Name       inflect.Name
	commonType inflect.Name
	goType     inflect.Name
}

func (a Attr) GoType() inflect.Name {
	if a.goType != "" {
		return a.goType
	}
	switch a.commonType {
	case "timestamp", "datetime", "date", "time":
		return "time.Time"
	}
	return a.commonType
}

func (a Attr) CommonType() inflect.Name {
	if a.commonType != "" {
		return a.commonType
	}
	return a.commonType
}

type Attrs []Attr
