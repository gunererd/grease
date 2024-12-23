package filemanager

type EntryType int

const (
	File EntryType = iota
	Directory
)

type Entry struct {
	Name string
	Type EntryType
}

type DirectoryReader interface {
	Read(path string) ([]Entry, error)
}
