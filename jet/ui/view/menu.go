package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

type Menu struct {
	*gocui.View
	Gui       *gocui.Gui
	ViewModel viewmodel.Menu
	Name      string
}

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

func (m *Menu) Render(v *gocui.View) error {
	m.View = v
	m.SelBgColor = gocui.ColorCyan
	m.SelFgColor = gocui.ColorBlack
	m.Highlight = true

	for _, item := range m.ViewModel.Items {
		fmt.Fprintln(v, item.Title)
	}

	m.ViewModel.OnChange(m.ViewModel.Items[0])

	return nil
}

func (m *Menu) selectNextLine(g *gocui.Gui, v *gocui.View) error {
	return m.selectLine(v, 1)
}

func (m *Menu) selectPrevLine(g *gocui.Gui, v *gocui.View) error {
	return m.selectLine(v, -1)
}

func (m *Menu) selectLine(v *gocui.View, change int) error {
	if v == nil {
		return nil
	}

	x, y := v.Cursor()
	curr := y + change

	if curr < 0 || curr >= len(m.ViewModel.Items)-1 {
		return nil
	}

	v.SetCursor(x, curr)
	m.ViewModel.OnChange(m.ViewModel.Items[y+change])

	return nil
}
