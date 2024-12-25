package types

type DirectoryReader interface {
	Read(path string) ([]Entry, error)
}
