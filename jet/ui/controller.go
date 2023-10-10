package ui

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui/view"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

type Controller struct {
	g *gocui.Gui
}

func NewController(g *gocui.Gui) *Controller {
	c := &Controller{
		g: g,
	}

	return c
}

func (c *Controller) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Controller) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	menu := viewmodel.Menu{
		Items: []*viewmodel.MenuItem{
			{
				Title: "Branches",
				Value: "branches",
			},
			{
				Title: "Tags",
				Value: "tags",
			},
		},
		OnSelected: func(item *viewmodel.MenuItem) error {
			log.Fatal(item.Title)
			return nil
		},
	}

	if v, err := g.SetView("branches", 0, 0, maxX-1, maxY-1); err != nil {
		v.Autoscroll = true
		menuView, _ := view.NewMenu(g, menu, "branches")
		menuView.Render(v)
	}

	g.SetCurrentView("branches")

	return nil
}
