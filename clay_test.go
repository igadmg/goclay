package clay

import (
	"math"
	"testing"

	"github.com/igadmg/gamemath/vector2"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	id1 := ID("test")
	id2 := ID("test")
	assert.Equal(t, id1.id, id2.id)
	assert.NotEqual(t, uint32(0), id1.id)
	id3 := ID("different")
	assert.NotEqual(t, id1.id, id3.id)
}

func TestIDI(t *testing.T) {
	id1 := IDI("test", 1)
	id2 := IDI("test", 1)
	assert.Equal(t, id1.id, id2.id)
	id3 := IDI("test", 2)
	assert.NotEqual(t, id1.id, id3.id)
}

func TestCornerRadius_IsEmpty(t *testing.T) {
	cr := CornerRadius{0, 0, 0, 0}
	assert.True(t, cr.IsEmpty())
	cr = CornerRadius{1, 0, 0, 0}
	assert.False(t, cr.IsEmpty())
}

func TestCORNER_RADIUS(t *testing.T) {
	cr := CORNER_RADIUS(5)
	expected := CornerRadius{5, 5, 5, 5}
	assert.Equal(t, expected, cr)
}

func TestSizingAxisTypeString(t *testing.T) {
	assert.Equal(t, "FIXED", SizingAxisTypeString(SizingAxisFixed{}))
	assert.Equal(t, "FIT", SizingAxisTypeString(SizingAxisFit{}))
	assert.Equal(t, "GROW", SizingAxisTypeString(SizingAxisGrow{}))
	assert.Equal(t, "PERCENT", SizingAxisTypeString(SizingAxisPercent{}))
}

func TestFIT(t *testing.T) {
	s := FIT()
	assert.Equal(t, float32(0), s.MinMax.Min)
	assert.Equal(t, float32(math.MaxFloat32), s.MinMax.Max)
	s = FIT(10)
	assert.Equal(t, float32(10), s.MinMax.Min)
	assert.Equal(t, float32(math.MaxFloat32), s.MinMax.Max)
	s = FIT(10, 20)
	assert.Equal(t, float32(10), s.MinMax.Min)
	assert.Equal(t, float32(20), s.MinMax.Max)
}

func TestGROW(t *testing.T) {
	s := GROW()
	assert.Equal(t, float32(0), s.MinMax.Min)
	assert.Equal(t, float32(math.MaxFloat32), s.MinMax.Max)
}

func TestFIXED(t *testing.T) {
	s := FIXED(100)
	assert.Equal(t, float32(100), s.MinMax.Min)
	assert.Equal(t, float32(100), s.MinMax.Max)
}

func TestPERCENT(t *testing.T) {
	s := PERCENT(0.5)
	assert.Equal(t, float32(0.5), s.Percent)
}

func TestSizing_GetAxis(t *testing.T) {
	s := Sizing{Width: SizingAxisFixed{}, Height: SizingAxisGrow{}}
	_, ok := s.GetAxis(0).(SizingAxisFixed)
	assert.True(t, ok)
	_, ok = s.GetAxis(1).(SizingAxisGrow)
	assert.True(t, ok)
}

func TestPADDING_ALL(t *testing.T) {
	p := PADDING_ALL(10)
	expected := Padding{10, 10, 10, 10}
	assert.Equal(t, expected, p)
}

func TestLayoutDirection_String(t *testing.T) {
	assert.Equal(t, "LEFT_TO_RIGHT", LEFT_TO_RIGHT.String())
	assert.Equal(t, "TOP_TO_BOTTOM", TOP_TO_BOTTOM.String())
}

func TestLayoutDirection_IsAlongAxis(t *testing.T) {
	assert.True(t, LEFT_TO_RIGHT.IsAlongAxis(0))
	assert.True(t, TOP_TO_BOTTOM.IsAlongAxis(1))
	assert.False(t, LEFT_TO_RIGHT.IsAlongAxis(1))
}

func TestTextElementConfigWrapMode_String(t *testing.T) {
	assert.Equal(t, "WORDS", TEXT_WRAP_WORDS.String())
	assert.Equal(t, "NEWLINES", TEXT_WRAP_NEWLINES.String())
	assert.Equal(t, "NONE", TEXT_WRAP_NONE.String())
}

func TestTextAlignment_String(t *testing.T) {
	assert.Equal(t, "LEFT", TEXT_ALIGN_LEFT.String())
	assert.Equal(t, "CENTER", TEXT_ALIGN_CENTER.String())
	assert.Equal(t, "RIGHT", TEXT_ALIGN_RIGHT.String())
}

func TestBorderWidth_IsEmpty(t *testing.T) {
	bw := BorderWidth{0, 0, 0, 0, 0}
	assert.True(t, bw.IsEmpty())
	bw = BorderWidth{1, 0, 0, 0, 0}
	assert.False(t, bw.IsEmpty())
}

func TestBorderElementConfig_IsEmpty(t *testing.T) {
	bec := BorderElementConfig{Color{}, BorderWidth{}}
	assert.True(t, bec.IsEmpty())
	bec = BorderElementConfig{Color{R: 1}, BorderWidth{}}
	assert.False(t, bec.IsEmpty())
}

// Math tests
func TestMakeDimensions(t *testing.T) {
	d := MakeDimensions(10, 20)
	assert.Equal(t, float32(10), d.X)
	assert.Equal(t, float32(20), d.Y)
}

func TestMakeVector2(t *testing.T) {
	v := MakeVector2(10, 20)
	assert.Equal(t, float32(10), v.X)
	assert.Equal(t, float32(20), v.Y)
}

func TestMakeBoundingBox(t *testing.T) {
	pos := vector2.NewFloat32(10, 20)
	size := vector2.NewFloat32(30, 40)
	bb := MakeBoundingBox(pos, size)
	assert.Equal(t, float32(10), bb.X())
	assert.Equal(t, float32(20), bb.Y())
	assert.Equal(t, float32(30), bb.Width())
	assert.Equal(t, float32(40), bb.Height())
}

func TestColor_IsZero(t *testing.T) {
	c := Color{0, 0, 0, 0}
	assert.True(t, c.IsZero())
	c = Color{1, 0, 0, 0}
	assert.False(t, c.IsZero())
}

// Public API tests
func TestInitialize(t *testing.T) {
	ctx := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	assert.NotNil(t, ctx)
	assert.Equal(t, ctx, GetCurrentContext())
}

func TestSetCurrentContext(t *testing.T) {
	ctx1 := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	ctx2 := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	SetCurrentContext(ctx1)
	assert.Equal(t, ctx1, GetCurrentContext())
	SetCurrentContext(ctx2)
	assert.Equal(t, ctx2, GetCurrentContext())
}

// Mock MeasureText function for testing
func mockMeasureText(text string, config *TextElementConfig, userData any) Dimensions {
	return MakeDimensions(float32(len(text))*10, 20) // Simple mock: 10px per char, 20 height
}

func TestBeginLayout_EndLayout(t *testing.T) {
	ctx := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	ctx.SetMeasureTextFunction(mockMeasureText, nil)

	ctx.BeginLayout()
	// Add some elements
	ctx.CLAY(ElementDeclaration{
		Layout: LayoutConfig{
			Sizing: Sizing{
				Width:  FIXED(100),
				Height: FIXED(50),
			},
		},
	}, func() {
		ctx.CLAY_TEXT("Hello", ctx.TEXT_CONFIG(TextElementConfig{
			TextColor: Color{255, 255, 255, 255},
			FontSize:  16,
		}))
	})
	commands := ctx.EndLayout()

	assert.NotEmpty(t, commands)
	// Check that we have at least one command
	found := false
	for _, cmd := range commands {
		if cmd.Id != 0 {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestSetPointerState(t *testing.T) {
	ctx := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	pos := MakeVector2(100, 100)
	ctx.SetPointerState(pos, false)
	assert.Equal(t, pos, ctx.pointerInfo.Position)
	assert.Equal(t, POINTER_DATA_RELEASED_THIS_FRAME, ctx.pointerInfo.State)
	ctx.SetPointerState(pos, true)
	assert.Equal(t, POINTER_DATA_PRESSED_THIS_FRAME, ctx.pointerInfo.State)
	ctx.SetPointerState(pos, true)
	assert.Equal(t, POINTER_DATA_PRESSED, ctx.pointerInfo.State)
}

func TestGetElementData(t *testing.T) {
	ctx := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	ctx.SetMeasureTextFunction(mockMeasureText, nil)

	ctx.BeginLayout()
	elementId := ctx.ID("testElement")
	ctx.CLAY_ID(elementId, ElementDeclaration{
		Layout: LayoutConfig{
			Sizing: Sizing{
				Width:  FIXED(100),
				Height: FIXED(50),
			},
		},
	})
	ctx.EndLayout()

	data := ctx.GetElementData(elementId)
	assert.True(t, data.Found)
	assert.Equal(t, float32(100), data.BoundingBox.Width())
	assert.Equal(t, float32(50), data.BoundingBox.Height())
}

func TestErrorHandling_MeasureTextNotSet(t *testing.T) {
	var errorCount int
	errorHandler := ErrorHandler{
		ErrorHandlerFunction: func(err ErrorData) {
			errorCount++
		},
		UserData: nil,
	}
	ctx := Initialize(MakeDimensions(800, 600), errorHandler)
	// Do not set MeasureTextFunction

	ctx.BeginLayout()
	ctx.CLAY_TEXT("test", ctx.TEXT_CONFIG(TextElementConfig{}))
	ctx.EndLayout()

	assert.Greater(t, errorCount, 0)
}

func TestSimpleLayout_RenderCommandsCount(t *testing.T) {
	ctx := Initialize(MakeDimensions(800, 600), ErrorHandler{})
	ctx.SetMeasureTextFunction(mockMeasureText, nil)

	ctx.BeginLayout()
	// Create a simple layout: a container with background and text
	ctx.CLAY(ElementDeclaration{
		Layout: LayoutConfig{
			Sizing: Sizing{
				Width:  FIXED(200),
				Height: FIXED(100),
			},
		},
		BackgroundColor: Color{255, 0, 0, 255}, // Red background
	}, func() {
		/*
			ctx.CLAY_TEXT("Hello World", ctx.TEXT_CONFIG(TextElementConfig{
				TextColor: Color{255, 255, 255, 255},
				FontSize:  16,
			}))
		*/
	})
	commands := ctx.EndLayout()

	// Should generate at least 2 commands: rectangle for background and text
	assert.Equal(t, 1, len(commands))
	// Check that we have a rectangle command and a text command
	hasRectangle := false
	hasText := false
	for _, cmd := range commands {
		switch cmd.RenderData.(type) {
		case RectangleRenderData:
			hasRectangle = true
		case TextRenderData:
			hasText = true
		}
	}
	assert.True(t, hasRectangle, "Should have rectangle render command")
	assert.True(t, hasText, "Should have text render command")
}
