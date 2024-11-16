package runtime

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui"
)

// Run starts the UI
func Run(exitChannel chan string) {
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	controller, err := ui.NewController(g, exitChannel)

	if err != nil {
		log.Panicln(err)
	}

	g.SetManager(controller)

	// wire up all the type keys
	for _, c := range ".-_/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" {
		g.SetKeybinding("", c, gocui.ModNone, controller.OnType(c))
	}

	// wire up all the delete keys
	g.SetKeybinding("", gocui.KeyBackspace2, gocui.ModNone, controller.OnBackspace)
	g.SetKeybinding("", gocui.KeyBackspace, gocui.ModNone, controller.OnBackspace)
	g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, controller.OnBackspace)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, controller.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
