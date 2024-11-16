package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

// Menu is the view for the menu
type Menu struct {
	*gocui.View
	Gui       *gocui.Gui
	ViewModel viewmodel.Menu
	Name      string
}

// NewMenu creates a new menu view
func NewMenu(g *gocui.Gui, vm viewmodel.Menu, name string) (*Menu, error) {
	m := &Menu{
		Gui:       g,
		ViewModel: vm,
		Name:      name,
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, m.selectNextLine); err != nil {
		return nil, err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, m.selectPrevLine); err != nil {
		return nil, err
	}

	g.SetKeybinding(m.Name, gocui.KeyEnter, gocui.ModNone, func(menuGui *gocui.Gui, menuView *gocui.View) error {
		_, line := menuView.Cursor()
		item, _ := menuView.Line(line)

		for _, v := range m.ViewModel.Items {
			if v.Title == item {
				return m.ViewModel.OnSelected(v)
			}
		}

		return nil
	})

	return m, nil
}

// Render renders the menu view
func (m *Menu) Render(v *gocui.View, opts ...RenderOption) error {
	m.View = v

	for _, opt := range opts {
		opt(v)
	}

	for _, item := range m.ViewModel.Items {
		fmt.Fprintln(v, item.Title)
	}

	if (len(m.ViewModel.Items)) == 0 {
		return nil
	}

	return nil
}

// selectNextLine selects the next line in the menu
func (m *Menu) selectNextLine(g *gocui.Gui, v *gocui.View) error {
	return m.selectLine(v, 1)
}

// selectPrevLine selects the previous line in the menu
func (m *Menu) selectPrevLine(g *gocui.Gui, v *gocui.View) error {
	return m.selectLine(v, -1)
}

// selectLine selects a line in the menu
func (m *Menu) selectLine(v *gocui.View, change int) error {
	if v == nil {
		return nil
	}

	x, y := v.Cursor()

	curr := y + change

	ox, oy := v.Origin()
	index := oy + curr

	if index < 0 || index >= len(m.ViewModel.Items)-1 {
		return nil
	}

	if err := v.SetCursor(x, curr); err != nil {
		if err := v.SetOrigin(ox, oy+change); err != nil {
			return err
		}
	}

	m.ViewModel.OnChange(m.ViewModel.Items[index])

	return nil
}
