package FilterLeveldb

import (
	"testgui/internal/filterdatabase"
)

type NameDatabase struct {
	name string
}

func NewFileterLeveldb(database string) filterdatabase.FilterData {
	return &NameDatabase{
		name: database,
	}
}

func (l *NameDatabase) Filters() {

}
