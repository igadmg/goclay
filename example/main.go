package main

import (
	clay "github.com/igadmg/goclay"
	"github.com/igadmg/goex/image/colorex"
	"github.com/igadmg/raylib-go/raymath/vector2"
)

var COLOR_LIGHT colorex.RGBA = colorex.RGBA{R: 224, G: 215, B: 210, A: 255}
var COLOR_RED colorex.RGBA = colorex.RGBA{R: 168, G: 66, B: 28, A: 255}
var COLOR_ORANGE colorex.RGBA = colorex.RGBA{R: 225, G: 138, B: 50, A: 255}

func main() {
	screenSize := vector2.NewFloat32(320, 200)
	mousePosition := vector2.NewFloat32(160, 100)
	mouseWheel := vector2.NewFloat32(0, 0)
	isMouseDown := false
	var profilePicture any
	var deltaTime float32

	// Note: screenWidth and screenHeight will need to come from your environment, Clay doesn't handle window related tasks
	clay.Initialize(nil, screenSize, clay.ErrorHandler{})

	for range 10 {
		// Optional: Update internal layout dimensions to support resizing
		clay.SetLayoutDimensions(screenSize)
		// Optional: Update internal pointer position for handling mouseover / click / touch events - needed for scrolling & debug tools
		clay.SetPointerState(mousePosition, isMouseDown)
		// Optional: Update internal pointer position for handling mouseover / click / touch events - needed for scrolling and debug tools
		clay.UpdateScrollContainers(true, mouseWheel, deltaTime)

		// All clay layouts are declared between clay.BeginLayout and clay.EndLayout
		clay.BeginLayout()

		// An example of laying out a UI with a fixed width sidebar and flexible width main content
		clay.CLAY(clay.ElementDeclaration{
			Id: clay.ID("OuterContainer"),
			Layout: clay.LayoutConfig{
				Sizing:   clay.Sizing{Width: clay.SIZING_GROW(0), Height: clay.SIZING_GROW(0)},
				Padding:  clay.PADDING_ALL(16),
				ChildGap: 16},
			BackgroundColor: colorex.RGBA{R: 250, G: 250, B: 255, A: 255},
		}, func() {
			clay.CLAY(clay.ElementDeclaration{
				Id: clay.ID("SideBar"),
				Layout: clay.LayoutConfig{LayoutDirection: clay.TOP_TO_BOTTOM,
					Sizing:   clay.Sizing{Width: clay.SIZING_FIXED(300), Height: clay.SIZING_GROW(0)},
					Padding:  clay.PADDING_ALL(16),
					ChildGap: 16},
				BackgroundColor: COLOR_LIGHT,
			}, func() {
				clay.CLAY(clay.ElementDeclaration{
					Id: clay.ID("ProfilePictureOuter"),
					Layout: clay.LayoutConfig{Sizing: clay.Sizing{Width: clay.SIZING_GROW(0)},
						Padding:  clay.PADDING_ALL(16),
						ChildGap: 16, ChildAlignment: clay.ChildAlignment{Y: clay.ALIGN_Y_CENTER}},
					BackgroundColor: COLOR_RED,
				}, func() {
					clay.CLAY(clay.ElementDeclaration{
						Id: clay.ID("ProfilePicture"),
						Layout: clay.LayoutConfig{
							Sizing: clay.Sizing{Width: clay.SIZING_FIXED(60), Height: clay.SIZING_FIXED(60)},
						},
						Image: clay.ImageElementConfig{
							ImageData:        profilePicture,
							SourceDimensions: vector2.NewFloat32(60, 60),
						},
					})
					//clay.TEXT(clay.STRING("Clay - UI Library"), clay.TEXT_CONFIG({ .fontSize = 24, .textColor = {255, 255, 255, 255} }));
				})

				// Standard C code like loops etc work inside components
				//for (int i = 0; i < 5; i++) {
				//    SidebarItemComponent();
				//}

				clay.CLAY(clay.ElementDeclaration{
					Id: clay.ID("MainContent"),
					Layout: clay.LayoutConfig{
						Sizing: clay.Sizing{Width: clay.SIZING_GROW(0), Height: clay.SIZING_GROW(0)},
					},
					BackgroundColor: COLOR_LIGHT})
			})
		})

		// All clay layouts are declared between clay.BeginLayout and clay.EndLayout
		renderCommands := clay.EndLayout()
		_ = renderCommands

		// More comprehensive rendering examples can be found in the renderers/ directory
		//for (int i = 0; i < renderCommands.length; i++) {
		//    clay.RenderCommand *renderCommand = &renderCommands.internalArray[i];
		//
		//    switch (renderCommand->commandType) {
		//        case clay.RENDER_COMMAND_TYPE_RECTANGLE: {
		//            DrawRectangle( renderCommand->boundingBox, renderCommand->renderData.rectangle.backgroundColor);
		//        }
		//        // ... Implement handling of other command types
		//    }
		//}
	}
}
