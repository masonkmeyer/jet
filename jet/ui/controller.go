package ui

import (
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet"
	"github.com/masonkmeyer/jet/jet/ui/view"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

const (
	BRANCHES = "branches"
	LOGS     = "logs"
	GRAPH    = "graph"
)

type Controller struct {
	g          *gocui.Gui
	git        *jet.Git
	currentLog string
	graph      string
}

func NewController(g *gocui.Gui) *Controller {
	c := &Controller{
		g:          g,
		git:        &jet.Git{},
		currentLog: "",
		graph:      "",
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
			results := c.git.Logs(item.Value, "-1")

			c.currentLog = results
			c.g.DeleteView(LOGS)

			graphLog := c.git.Logs(item.Value, "--graph", "--oneline", "--decorate", "--color", "--abbrev-commit", "--date=relative", "--format=format:%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(auto)%d%C(reset)")
			c.graph = graphLog
			c.g.DeleteView(GRAPH)
			return nil
		},
	}

	if v, err := g.SetView(BRANCHES, 0, 0, maxX/2-1, maxY/2-1); err != nil {
		v.Autoscroll = true
		v.Title = "Recent Branches"
		menuView, _ := view.NewMenu(g, menu, BRANCHES)
		menuView.Render(v)
	}

	if v, err := g.SetView(LOGS, maxX/2, 0, maxX-1, maxY/2-1); err != nil {
		v.Title = "Recent Commit Message"
		textView := view.NewText(g, viewmodel.Text{Value: c.currentLog, Autoscroll: true, Wrap: true}, LOGS)
		textView.Render(v)
	}

	if v, err := g.SetView(GRAPH, 0, maxY/2, maxX-1, maxY-1); err != nil {
		textView := view.NewText(g, viewmodel.Text{Value: c.graph, Autoscroll: true, Wrap: true}, LOGS)
		textView.Render(v)
	}

	g.SetCurrentView(BRANCHES)

	return nil
}
