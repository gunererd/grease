package filemanager

import (
	"strings"
)

func EntriesToBuffer(entries []Entry) string {
	var sb strings.Builder

	for i, entry := range entries {
		sb.WriteString(entry.Name)
		if i < len(entries)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
