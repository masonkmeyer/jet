package runtime

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui"
)

func Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	controller := ui.NewController(g)

	g.SetManager(controller)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, controller.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
