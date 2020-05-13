package oscutil

import (
	"fmt"
	"log"
	"github.com/E-Health/goscar"
	"github.com/jroimartin/gocui"
)

var RecordCount int
var CsvMap, CsvMapValid []map[string]string

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	getLine(g, v)
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	getLine(g, v)
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	mainOutput(g, &l)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("title", -1, -1, maxX, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorRed
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "OSCAR eForm Export Tool Helper by Bell Eapen")
		fmt.Fprintln(v, "Valid Records: ", RecordCount)
	}
	if _, err := g.SetView("main", 30, 4, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		message := "OSCAR Helper v 1.0.0"
		mainOutput(g, &message)
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("side", -1, 4, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		sideOutput(g)
	}
	return nil
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	return nil
}

func mainOutput(g *gocui.Gui, message *string) {
	if v, err := g.SetCurrentView("main"); err != nil {
		log.Panicln(err)
	} else {
		v.Editable = true
		v.Wrap = true
		v.Clear()
		fmt.Fprintln(v, *message)
		fmt.Fprintln(v, "---------------------")

		toPrint := goscar.GetStats(*message, RecordCount, CsvMapValid)

		if toPrint["num"] > 0 {
			for k, val := range toPrint {
				if k != "num" {
					fmt.Fprintln(v, k , " : ", val)
				}
			}
		} else {
			total := 0.0
			for _, val := range toPrint {
				total = total + val
			}
			for k, val := range toPrint {
				if k != "num"{
					fmt.Fprintln(v, k , " : ", val, "(", (val*100/total), " % )")
				}
			}
		}
		g.SetCurrentView("side")
		recover()
	}
}

func sideOutput(g *gocui.Gui) {
	toIgnore := []string{"id", "fdid", "dateCreated", "eform_link", "StaffSig", "SubmitButton", "efmfid"}
	if v, err := g.SetCurrentView("side"); err != nil {
		log.Panicln(err)
	} else {
		firstRecord := CsvMap[0]
		for key, _ := range firstRecord {
			if !goscar.IsMember(key, toIgnore) {
				fmt.Fprintln(v, key)
			}
		}
	}
}