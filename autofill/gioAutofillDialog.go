package autofill

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sort"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	gops "github.com/mitchellh/go-ps"
)

type AutofillEntry struct {
	Username string
	Name     string
	UUID     string
}

var autofillEntries = []AutofillEntry{}
var onAutofill func(string, chan bool)
var selectedEntry = 0

type scoredAutofillEntry struct {
	autofillEntry AutofillEntry
	score         int
}

func getProcesses() []string {
	var processes []string
	procs, err := gops.Processes()
	if err != nil {
		return []string{}
	}

	for _, proc := range procs {
		processes = append(processes, proc.Executable())
	}

	return processes
}

func GetFilteredAutofillEntries(entries []AutofillEntry, filter string) []AutofillEntry {
	if len(filter) < 2 {
		return []AutofillEntry{}
	}

	processes := getProcesses()

	scoredEntries := []scoredAutofillEntry{}
	for _, entry := range entries {
		score := 0
		if strings.Contains(strings.ToLower(entry.Name), strings.ToLower(filter)) {
			score += 10
		}
		if strings.HasPrefix(strings.ToLower(entry.Name), strings.ToLower(filter)) {
			score += 5
		}
		if strings.Contains(strings.ToLower(entry.Username), strings.ToLower(filter)) {
			score += 3
		}
		if strings.HasPrefix(strings.ToLower(entry.Username), strings.ToLower(filter)) {
			score += 2
		}

		for _, process := range processes {
			if strings.Contains(strings.ToLower(entry.Name), strings.ToLower(process)) {
				score += 5
				break
			}
		}

		scoredEntries = append(scoredEntries, scoredAutofillEntry{entry, score})
	}

	sort.Slice(scoredEntries, func(i, j int) bool {
		return scoredEntries[i].score > scoredEntries[j].score
	})

	var filteredEntries []AutofillEntry
	for _, scoredEntry := range scoredEntries {
		filteredEntries = append(filteredEntries, scoredEntry.autofillEntry)
	}

	return filteredEntries
}

func RunAutofill(entries []AutofillEntry, onAutofillFunc func(string, chan bool)) {
	autofillEntries = entries
	onAutofill = onAutofillFunc

	go func() {
		w := app.NewWindow()
		w.Option(app.Size(unit.Dp(600), unit.Dp(800)))
		w.Option(app.Decorated(false))
		w.Perform(system.ActionCenter)
		w.Perform(system.ActionRaise)
		lineEditor.Focus()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

var lineEditor = &widget.Editor{
	SingleLine: true,
	Submit:     true,
}

var (
	unselected     = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	unselectedText = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	background     = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	selected       = color.NRGBA{R: 0x65, G: 0x1F, B: 0xFF, A: 0xFF}
	selectedText   = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

var th = material.NewTheme(gofont.Collection())
var list = layout.List{Axis: layout.Vertical}

func doLayout(gtx layout.Context) layout.Dimensions {
	var filteredEntries []AutofillEntry = GetFilteredAutofillEntries(autofillEntries, lineEditor.Text())

	if selectedEntry >= 10 || selectedEntry >= len(filteredEntries) {
		selectedEntry = 0
	}

	return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						searchBox := material.Editor(th, lineEditor, "Search query")
						searchBox.Color = selectedText
						border := widget.Border{Color: selectedText, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
						return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, searchBox.Layout)
						})
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {

					return list.Layout(gtx, len(filteredEntries), func(gtx layout.Context, i int) layout.Dimensions {
						entry := filteredEntries[i]

						return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							isSelected := i == selectedEntry
							var color color.NRGBA
							if isSelected {
								color = selected
							} else {
								color = unselected
							}

							return Background{Color: color, CornerRadius: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									dimens := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											t := material.H6(th, entry.Name)
											if isSelected {
												t.Color = selectedText
											} else {
												t.Color = unselectedText
											}
											return t.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											t := material.Body1(th, entry.Username)
											if isSelected {
												t.Color = selectedText
											} else {
												t.Color = unselectedText
											}

											return t.Layout(gtx)
										}),
									)
									return dimens
								})
							})
						})
					})
				})
			}))
	})
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			key.InputOp{
				Keys: key.Set(key.NameReturn + "|" + key.NameEscape + "|" + key.NameDownArrow),
				Tag:  0,
			}.Add(gtx.Ops)
			t := lineEditor.Events()
			for _, ev := range t {
				switch ev.(type) {
				case widget.SubmitEvent:
					entries := GetFilteredAutofillEntries(autofillEntries, lineEditor.Text())
					if len(entries) == 0 {
						fmt.Println("no entries")
						continue
					} else {
						w.Perform(system.ActionMinimize)
						c := make(chan bool)
						go onAutofill(entries[selectedEntry].UUID, c)
						go func() {
							<-c
							os.Exit(0)
						}()
					}
				}
			}

			test := gtx.Events(0)
			for _, ev := range test {
				switch ev := ev.(type) {
				case key.Event:
					switch ev.Name {
					case key.NameReturn:
						fmt.Println("uncaught submit")
						return nil
					case key.NameDownArrow:
						if ev.State == key.Press {
							selectedEntry++
							if selectedEntry >= 10 {
								selectedEntry = 0
							}
						}
					case key.NameEscape:
						os.Exit(0)
					}
				}
			}

			doLayout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

type Background struct {
	Color        color.NRGBA
	CornerRadius unit.Dp
}

func (b Background) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	m := op.Record(gtx.Ops)
	dims := w(gtx)
	size := dims.Size
	call := m.Stop()
	if r := gtx.Dp(b.CornerRadius); r > 0 {
		defer clip.RRect{
			Rect: image.Rect(0, 0, size.X, size.Y),
			NE:   r, NW: r, SE: r, SW: r,
		}.Push(gtx.Ops).Pop()
	}
	fill{b.Color}.Layout(gtx, size)
	call.Add(gtx.Ops)
	return dims
}

type fill struct {
	col color.NRGBA
}

func (f fill) Layout(gtx layout.Context, sz image.Point) layout.Dimensions {
	defer clip.Rect(image.Rectangle{Max: sz}).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: f.col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: sz}
}
