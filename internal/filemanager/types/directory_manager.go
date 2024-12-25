package types

type DirectoryManager interface {
	ReadDirectory() ([]Entry, error)
	ChangeDirectory(path string) error
	CurrentPath() string
	GetDirectoryContent() ([]byte, error)
}
