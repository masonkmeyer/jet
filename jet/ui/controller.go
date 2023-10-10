package ui

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet"
	"github.com/masonkmeyer/jet/jet/ui/view"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

type Controller struct {
	g   *gocui.Gui
	git *jet.Git
}

func NewController(g *gocui.Gui) *Controller {
	c := &Controller{
		g:   g,
		git: &jet.Git{},
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
		items = append(items, &viewmodel.MenuItem{
			Title: branch,
			Value: branch,
		})
	}

	menu := viewmodel.Menu{
		Items: items,
		OnSelected: func(item *viewmodel.MenuItem) error {
			log.Fatal(item.Title)
			return nil
		},
	}

	if v, err := g.SetView("branches", 0, 0, maxX/2-1, maxY-1); err != nil {
		v.Autoscroll = true
		menuView, _ := view.NewMenu(g, menu, "branches")
		menuView.Render(v)
	}

	g.SetCurrentView("branches")

	return nil
}
