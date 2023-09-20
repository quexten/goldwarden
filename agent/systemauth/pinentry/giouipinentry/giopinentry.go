// go:build windows || darwin
package giouipinentry

import (
	"fmt"
	"image"
	"image/color"
	"log"

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
)

func GetPin(title string, description string, onPinEntered func(string)) {
	go func() {
		w := app.NewWindow()
		w.Option(app.Size(unit.Dp(500), unit.Dp(300)))
		w.Option(app.Decorated(false))
		w.Perform(system.ActionCenter)
		w.Perform(system.ActionRaise)
		if err := runPinentry(w, title, description, onPinEntered); err != nil {
			log.Fatal(err)
		}
	}()
}

func GetApproval(title string, description string, onApproval func(bool)) {
	go func() {
		w := app.NewWindow()
		w.Option(app.Size(unit.Dp(500), unit.Dp(300)))
		w.Option(app.Decorated(false))
		w.Perform(system.ActionCenter)
		w.Perform(system.ActionRaise)
		if err := runApproval(w, title, description, onApproval); err != nil {
			log.Fatal(err)
		}
	}()
}

var (
	// #000000
	unselected = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	// #FFFFFF
	unselectedText = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	// #000000
	background = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	// #651FFF
	selected = color.NRGBA{R: 0x65, G: 0x1F, B: 0xFF, A: 0xFF}
	// #FFFFFF
	selectedText = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	// #4caf50
	buttonOk = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	// #F44336
	buttonCancel = color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}
)

var th = material.NewTheme(gofont.Collection())

func runPinentry(w *app.Window, title string, description string, onPinEntered func(string)) error {
	var lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditor.Focus()
	var ops op.Ops

	var btnOk widget.Clickable
	var btnCancel widget.Clickable

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
					onPinEntered(lineEditor.Text())
					w.Perform(system.ActionClose)
				}
			}
			if btnOk.Clicked() {
				onPinEntered(lineEditor.Text())
				w.Perform(system.ActionClose)
			}
			if btnCancel.Clicked() {
				onPinEntered("")
				w.Perform(system.ActionClose)
			}

			test := gtx.Events(0)
			for _, ev := range test {
				switch ev := ev.(type) {
				case key.Event:
					switch ev.Name {
					case key.NameReturn:
						fmt.Println("uncaught submit")
						return nil
					case key.NameEscape:
						onPinEntered("")
						w.Perform(system.ActionClose)
					}
				}
			}

			Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					//title
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, title)
								title.Color = unselectedText
								return title.Layout(gtx)
							})
						})
					}),
					// description
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								body := material.Body1(th, description)
								body.Color = unselectedText
								return body.Layout(gtx)
							})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								searchBox := material.Editor(th, lineEditor, "Pin")
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
							return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									acceptBtn := material.Button(th, &btnOk, "Accept")
									acceptBtn.Background = buttonOk
									return acceptBtn.Layout(gtx)
								}),
								layout.Rigid(
									layout.Spacer{Height: unit.Dp(10)}.Layout,
								),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									cancel := material.Button(th, &btnCancel, "Cancel")
									cancel.Background = buttonCancel
									return cancel.Layout(gtx)
								}),
							)
						})
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}

	return nil
}

func runApproval(w *app.Window, title string, description string, onApproval func(bool)) error {
	var ops op.Ops

	var btnOk widget.Clickable
	var btnCancel widget.Clickable

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

			if btnOk.Clicked() {
				onApproval(true)
				w.Perform(system.ActionClose)
			}
			if btnCancel.Clicked() {
				onApproval(false)
				w.Perform(system.ActionClose)
			}

			test := gtx.Events(0)
			for _, ev := range test {
				switch ev := ev.(type) {
				case key.Event:
					switch ev.Name {
					case key.NameReturn:
						fmt.Println("uncaught submit")
						return nil
					case key.NameEscape:
						onApproval(false)
						w.Perform(system.ActionClose)
					}
				}
			}

			Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					//title
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, title)
								title.Color = unselectedText
								return title.Layout(gtx)
							})
						})
					}),
					// description
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return Background{Color: background, CornerRadius: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								body := material.Body1(th, description)
								body.Color = unselectedText
								return body.Layout(gtx)
							})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									acceptBtn := material.Button(th, &btnOk, "Accept")
									acceptBtn.Background = buttonOk
									return acceptBtn.Layout(gtx)
								}),
								layout.Rigid(
									layout.Spacer{Height: unit.Dp(10)}.Layout,
								),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									cancel := material.Button(th, &btnCancel, "Cancel")
									cancel.Background = buttonCancel
									return cancel.Layout(gtx)
								}),
							)
						})
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}

	return nil
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
