package main

import (
	"log"
	"os"
	"runtime/pprof"

	clay "github.com/igadmg/goclay"
)

var COLOR_LIGHT clay.Color = clay.Color{R: 224, G: 215, B: 210, A: 255}
var COLOR_RED clay.Color = clay.Color{R: 168, G: 66, B: 28, A: 255}
var COLOR_ORANGE clay.Color = clay.Color{R: 225, G: 138, B: 50, A: 255}

// Layout config is just a struct that can be declared statically, or inline
var sidebarItemConfig clay.ElementDeclaration

func gx_Must[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}

func pprofex_WriteCPUProfile(fileName string) (func(), error) {
	f, err := os.Create(fileName + ".prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
		return func() {}, err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
		return func() {}, err
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}

// Re-useable components are just normal functions
func SidebarItemComponent(ctx *clay.Context) {
	ctx.CLAY(sidebarItemConfig, func() {
		// children go here...
	})
}

func main() {
	defer gx_Must(pprofex_WriteCPUProfile("goclay"))()

	screenSize := clay.MakeDimensions(640, 480)
	mousePosition := clay.MakeVector2(160, 100)
	mouseWheel := clay.MakeVector2(0, 0)
	isMouseDown := false
	var profilePicture any = &struct{ ImageData []byte }{ImageData: nil}
	var deltaTime float32 = 0.1

	// Note: screenWidth and screenHeight will need to come from your environment, Clay doesn't handle window related tasks
	ui := clay.Initialize(screenSize, clay.ErrorHandler{})
	ui.SetMeasureTextFunction(func(text string, config *clay.TextElementConfig, userData any) clay.Dimensions {
		width := config.FontSize * 3 / 4
		return clay.MakeDimensions(
			width*uint16(len(text)),
			config.FontSize,
		)
	}, nil)

	sidebarItemConfig = clay.ElementDeclaration{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  ui.SIZING_GROW(0),
				Height: ui.SIZING_FIXED(50),
			},
		},
		BackgroundColor: COLOR_ORANGE,
	}

	for range 100000 {
		// Optional: Update internal layout dimensions to support resizing
		ui.SetLayoutDimensions(screenSize)
		// Optional: Update internal pointer position for handling mouseover / click / touch events - needed for scrolling & debug tools
		ui.SetPointerState(mousePosition, isMouseDown)
		// Optional: Update internal pointer position for handling mouseover / click / touch events - needed for scrolling and debug tools
		ui.UpdateScrollContainers(true, mouseWheel, deltaTime)

		// All clay layouts are declared between clay.BeginLayout and clay.EndLayout
		ui.BeginLayout()

		// An example of laying out a UI with a fixed width sidebar and flexible width main content
		ui.CLAY(clay.ElementDeclaration{
			Id: ui.ID("OuterContainer"),
			Layout: clay.LayoutConfig{
				Sizing:   clay.Sizing{Width: ui.SIZING_GROW(0), Height: ui.SIZING_GROW(0)},
				Padding:  ui.PADDING_ALL(16),
				ChildGap: 16},
			BackgroundColor: clay.Color{R: 250, G: 250, B: 255, A: 255},
		}, func() {
			ui.CLAY(clay.ElementDeclaration{
				Id: ui.ID("SideBar"),
				Layout: clay.LayoutConfig{LayoutDirection: clay.TOP_TO_BOTTOM,
					Sizing:   clay.Sizing{Width: ui.SIZING_FIXED(300), Height: ui.SIZING_GROW(0)},
					Padding:  ui.PADDING_ALL(16),
					ChildGap: 16},
				BackgroundColor: COLOR_LIGHT,
			}, func() {
				ui.CLAY(clay.ElementDeclaration{
					Id: ui.ID("ProfilePictureOuter"),
					Layout: clay.LayoutConfig{Sizing: clay.Sizing{Width: ui.SIZING_GROW(0)},
						Padding:  ui.PADDING_ALL(16),
						ChildGap: 16, ChildAlignment: clay.ChildAlignment{Y: clay.ALIGN_Y_CENTER}},
					BackgroundColor: COLOR_RED,
				}, func() {
					ui.CLAY(clay.ElementDeclaration{
						Id: ui.ID("ProfilePicture"),
						Layout: clay.LayoutConfig{
							Sizing: clay.Sizing{Width: ui.SIZING_FIXED(60), Height: ui.SIZING_FIXED(60)},
						},
						Image: clay.ImageElementConfig{
							ImageData:        profilePicture,
							SourceDimensions: clay.MakeDimensions(60, 60),
						},
					})
					ui.CLAY(clay.ElementDeclaration{
						Id: ui.ID("TextContent"),
						Layout: clay.LayoutConfig{
							Sizing: clay.Sizing{Width: ui.SIZING_GROW(0), Height: ui.SIZING_GROW(0)},
						},
						BackgroundColor: COLOR_LIGHT})
					/*/
					ctx.TEXT("Clay - UI Library", ctx.TEXT_CONFIG(clay.TextElementConfig{
						FontSize:  24,
						TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255},
					}))
					*/
				})

				// Standard C code like loops etc work inside components
				for range 5 {
					SidebarItemComponent(ui)
				}

				ui.CLAY(clay.ElementDeclaration{
					Id: ui.ID("MainContent"),
					Layout: clay.LayoutConfig{
						Sizing: clay.Sizing{Width: ui.SIZING_GROW(0), Height: ui.SIZING_GROW(0)},
					},
					BackgroundColor: COLOR_LIGHT})
			})
		})

		// All clay layouts are declared between clay.BeginLayout and clay.EndLayout
		renderCommands := ui.EndLayout()

		// More comprehensive rendering examples can be found in the renderers/ directory
		for _, renderCommand := range renderCommands {
			_ = renderCommand
			//fmt.Printf("%d: %s\t\t%s\n", renderCommand.Id, renderCommand.BoundingBox, reflect.TypeOf(renderCommand.RenderData))
		}
	}
}
