package types

type Entry interface {
	Name() string
	Type() EntryType
}

type EntryType int

const (
	File EntryType = iota
	Directory
)
