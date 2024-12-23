package filemanager

type DirectorySource struct {
	fm *FileManager
}

func NewDirectorySource(fm *FileManager) *DirectorySource {
	return &DirectorySource{
		fm: fm,
	}
}

func (s *DirectorySource) Read() ([]byte, error) {
	entries, err := s.fm.ReadDirectory()
	if err != nil {
		return nil, err
	}

	content := EntriesToBuffer(entries)
	return []byte(content), nil
}

func (s *DirectorySource) Name() string {
	return s.fm.currentPath
}

func (s *DirectorySource) Close() error {
	return nil
}
