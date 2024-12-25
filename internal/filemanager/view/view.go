package view

import (
	eTypes "github.com/gunererd/grease/internal/editor/types"
	"github.com/gunererd/grease/internal/filemanager/types"
)

type View struct {
	editor eTypes.Editor
}

func New(editor eTypes.Editor) types.View {
	return &View{
		editor: editor,
	}
}

func (v *View) Render() string {
	return v.editor.View()
}
