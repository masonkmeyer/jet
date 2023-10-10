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

func NewText(g *gocui.Gui, vm viewmodel.Text, name string) *Text {
	t := &Text{
		Gui:       g,
		ViewModel: vm,
		Name:      name,
	}

	return t
}

func (t *Text) Render(v *gocui.View) error {
	t.View = v
	t.Editable = false

	fmt.Fprintln(v, t.ViewModel.Value)

	return nil
}
