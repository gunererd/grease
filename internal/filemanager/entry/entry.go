package entry

import "github.com/gunererd/grease/internal/filemanager/types"

type entry struct {
	name      string
	entryType types.EntryType
}

func New(name string, entryType types.EntryType) types.Entry {
	return &entry{name: name, entryType: entryType}
}

func (e *entry) Name() string {
	return e.name
}

func (e *entry) Type() types.EntryType {
	return e.entryType
}
