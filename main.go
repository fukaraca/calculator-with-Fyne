package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Knetic/govaluate"
	"log"
)

type calculator struct {
	screen  string
	buttons map[string]*widget.Button
	window  fyne.Window
}

var screen binding.ExternalString
var calc calculator

func main() {
	a := app.New()
	w := a.NewWindow("Calculator")
	w.SetMainMenu(makeMenu(a, w))

	calc.buttons = make(map[string]*widget.Button)
	screen = binding.BindString(&calc.screen)
	textOnScreen := widget.NewLabelWithData(screen)

	toolBar := makeToolbar(a, w)
	builtScreen := container.NewBorder(toolBar, textOnScreen, nil, nil)
	buttons := container.NewGridWithColumns(4, nil...)

	butStrings := []string{"(", ")", "C", "/", "7", "8", "9", "*", "4", "5", "6", "-", "1", "2", "3", "+", "Del", "0", ".", "="}
	for _, butString := range butStrings {
		buttons.Add(makeButton(butString))
	}
	rowy := container.NewGridWithRows(2, builtScreen, buttons)

	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		switch ke.Name {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "/", "*", "-", "+", ".":
			calc.screen += string(ke.Name)
			screen.Reload()
		case "KP_Enter", "Return":
			res, err := parseAndEvaluate(calc.screen)
			if err != nil {
				log.Println("calc err.", err)
				calc.screen = "Error" + err.Error()
				screen.Reload()
			}
			calc.screen = fmt.Sprintf("%.2f", res)
			screen.Reload()
		case "BackSpace":
			if len(calc.screen) >= 1 {
				calc.screen = calc.screen[:len(calc.screen)-1]
				screen.Reload()
			}
		case "Delete":
			calc.screen = ""
			screen.Reload()
		}

	})

	w.SetContent(rowy)

	rowy.Refresh()
	w.Resize(fyne.NewSize(330, 0))
	w.ShowAndRun()
}

//render buttons with their functions
func makeButton(butString string) fyne.CanvasObject {
	funk := func() {}
	switch butString {
	case "C":
		funk = func() {
			calc.screen = ""
			screen.Reload()
		}
	case "Del":
		funk = func() {
			if len(calc.screen) >= 1 {
				calc.screen = calc.screen[:len(calc.screen)-1]
				screen.Reload()
			}
		}

	case "=":
		funk = func() {
			res, err := parseAndEvaluate(calc.screen)
			if err != nil {
				log.Println("calc err.", err)
				calc.screen = "Error"
				screen.Reload()
			}
			calc.screen = fmt.Sprintf("%.2f", res)
			screen.Reload()
		}

	default:
		funk = func() {
			calc.screen += butString
			screen.Reload()
		}
	}
	b := widget.NewButton(butString, funk)

	calc.buttons[butString] = b

	return b
}

//create menu
func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	dark := fyne.NewMenuItem("Dark theme", func() {
		a.Settings().SetTheme(theme.DarkTheme())
	})
	light := fyne.NewMenuItem("Light theme", func() {
		a.Settings().SetTheme(theme.LightTheme())

	})
	themeM := fyne.NewMenu("Theme", dark, light)

	manuelEntry := fyne.NewMenuItem("Enter manually", func() {
		input := widget.NewEntry()
		input.SetPlaceHolder("Paste you expression")

		items := []*widget.FormItem{
			widget.NewFormItem("Expression:", input),
		}
		InputForm := dialog.NewForm("Enter the expression manually", "OK", "Cancel", items, nil, w)
		InputForm.Show()
		InputForm.Resize(fyne.NewSize(250, 150))

		InputForm.SetOnClosed(func() {
			calc.screen = input.Text
			screen.Reload()
		})
	})

	entry := fyne.NewMenu("Enter manually", manuelEntry)
	return fyne.NewMainMenu(themeM, entry)
}

//evaluator
func parseAndEvaluate(exp string) (float64, error) {
	express, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		log.Println("parsing couldn't be done:", err)
		return 0, err
	}
	result, err := express.Eval(nil)
	if err != nil {
		log.Println("evaluation couldn't be done", err)
		return 0, err
	}
	return result.(float64), nil
}

//create toolbar for cut,copy and paste
func makeToolbar(a fyne.App, w fyne.Window) fyne.CanvasObject {
	t := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {
			text, err := screen.Get()
			if err != nil {
				log.Println("text couldn't be copied:", err)
				return
			}
			w.Clipboard().SetContent(text)
			calc.screen = ""
			screen.Reload()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			text, err := screen.Get()
			if err != nil {
				log.Println("text couldn't be copied:", err)
				return
			}
			w.Clipboard().SetContent(text)

		}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {
			calc.screen += w.Clipboard().Content()
			screen.Reload()
		}))

	return t
}
