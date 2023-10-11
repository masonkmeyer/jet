package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

type Text struct {
	*gocui.View
	Gui       *gocui.Gui
	ViewModel viewmodel.Text
	Name      string
}

// NewText creates a new text view
func NewText(g *gocui.Gui, vm viewmodel.Text, name string) *Text {
	t := &Text{
		Gui:       g,
		ViewModel: vm,
		Name:      name,
	}

	return t
}

// Render renders the text view
func (t *Text) Render(v *gocui.View, opts ...RenderOption) error {
	t.View = v
	t.Editable = false

	for _, opt := range opts {
		opt(v)
	}

	fmt.Fprintln(v, t.ViewModel.Value)

	return nil
}
