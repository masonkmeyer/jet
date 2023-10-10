package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet"
	"github.com/masonkmeyer/jet/jet/ui/view"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

type Controller struct {
	g          *gocui.Gui
	git        *jet.Git
	currentLog string
}

func NewController(g *gocui.Gui) *Controller {
	c := &Controller{
		g:          g,
		git:        &jet.Git{},
		currentLog: "This is a message",
	}

	return c
}

func (c *Controller) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Controller) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	items := []*viewmodel.MenuItem{}

	branches := c.git.ListBranches("--sort=-committerdate", "--format=%(align:13,left)%(authordate:relative)%(end) %(refname:short)")

	for _, branch := range branches {
		parts := strings.Split(branch, " ")
		value := parts[len(parts)-1]

		items = append(items, &viewmodel.MenuItem{
			Title: branch,
			Value: value,
		})
	}

	menu := viewmodel.Menu{
		Items: items,
		OnSelected: func(item *viewmodel.MenuItem) error {
			return nil
		},
		OnChange: func(item *viewmodel.MenuItem) error {
			c.currentLog = item.Title
			c.g.DeleteView("logs")
			return nil
		},
	}

	if v, err := g.SetView("branches", 0, 0, maxX/2-1, maxY-1); err != nil {
		v.Autoscroll = true
		menuView, _ := view.NewMenu(g, menu, "branches")
		menuView.Render(v)
	}

	if v, err := g.SetView("logs", maxX/2, 0, maxX-1, maxY-1); err != nil {
		v.Autoscroll = true

		fmt.Fprintln(v, c.currentLog)
	}

	g.SetCurrentView("branches")

	return nil
}
