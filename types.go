package goclay

import (
	"github.com/igadmg/goex/image/colorex"
	"github.com/igadmg/raylib-go/raymath/rect2"
	"github.com/igadmg/raylib-go/raymath/vector2"
)

// Primarily created via the ID(), IDI(), ID_LOCAL() and IDI_LOCAL() macros.
// Represents a hashed string ID used for identifying and finding specific clay UI elements, required
// by functions such as PointerOver() and GetElementData().
type ElementId struct {
	id       uint32 // The resulting hash generated from the other fields.
	offset   uint32 // A numerical offset applied after computing the hash from stringId.
	baseId   uint32 // A base hash value to start from, for example the parent element ID is used when calculating ID_LOCAL().
	stringId string // The string id to hash.
}

var default_ElementId ElementId

// Note: If a compile error led you here, you might be trying to use ID with something other than a string literal. To construct an ID with a dynamic string, use SID instead.
func ID(label string) ElementId { return IDI(label, 0) }

func SID(label string) ElementId { return SIDI(label, 0) }

// Note: If a compile error led you here, you might be trying to use IDI with something other than a string literal. To construct an ID with a dynamic string, use SIDI instead.
func IDI(label string, index uint32) ElementId { return SIDI(label, index) }

func SIDI(label string, index uint32) ElementId { return hashString(label, index, 0) }

// Note: If a compile error led you here, you might be trying to use ID_LOCAL with something other than a string literal. To construct an ID with a dynamic string, use SID_LOCAL instead.
func ID_LOCAL(label string) ElementId { return IDI_LOCAL(label, 0) }

func SID_LOCAL(label string) ElementId { return SIDI_LOCAL(label, 0) }

// Note: If a compile error led you here, you might be trying to use IDI_LOCAL with something other than a string literal. To construct an ID with a dynamic string, use SIDI_LOCAL instead.
func IDI_LOCAL(label string, index uint32) ElementId { return SIDI_LOCAL(label, index) }

func SIDI_LOCAL(label string, index uint32) ElementId {
	return hashString(label, index, GetCurrentContext().getParentElementId())
}

// Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
type CornerRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomLeft  float32
	BottomRight float32
}

func (r CornerRadius) IsEmpty() bool {
	return r.TopLeft == 0 &&
		r.TopRight == 0 &&
		r.BottomLeft == 0 &&
		r.BottomRight == 0
}

// Element Configs ---------------------------

// Controls the direction in which child elements will be automatically laid out.
type LayoutDirection uint8

const (
	// (Default) Lays out child elements from left to right with increasing x.
	LEFT_TO_RIGHT LayoutDirection = iota
	// Lays out child elements from top to bottom with increasing y.
	TOP_TO_BOTTOM
)

// Controls the alignment along the x axis (horizontal) of child elements.
type LayoutAlignmentX uint8

const (
	// (Default) Aligns child elements to the left hand side of this element, offset by padding.width.left
	ALIGN_X_LEFT LayoutAlignmentX = iota
	// Aligns child elements to the right hand side of this element, offset by padding.width.right
	ALIGN_X_RIGHT
	// Aligns child elements horizontally to the center of this element
	ALIGN_X_CENTER
)

// Controls the alignment along the y axis (vertical) of child elements.
type LayoutAlignmentY uint8

const (
	// (Default) Aligns child elements to the top of this element, offset by padding.width.top
	ALIGN_Y_TOP LayoutAlignmentY = iota
	// Aligns child elements to the bottom of this element, offset by padding.width.bottom
	ALIGN_Y_BOTTOM
	// Aligns child elements vertically to the center of this element
	ALIGN_Y_CENTER
)

// Controls how the element takes up space inside its parent container.
type SizingType uint8

const (
	// (default) Wraps tightly to the size of the element's contents.
	SIZING_TYPE_FIT SizingType = iota
	// Expands along this axis to fill available space in the parent element, sharing it with other GROW elements.
	SIZING_TYPE_GROW
	// Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
	SIZING_TYPE_PERCENT
	// Clamps the axis size to an exact size in pixels.
	SIZING_TYPE_FIXED
)

// Controls how child elements are aligned on each axis.
type ChildAlignment struct {
	X LayoutAlignmentX // Controls alignment of children along the x axis.
	Y LayoutAlignmentY // Controls alignment of children along the y axis.
}

// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to,
// overriding sizing types such as FIT or GROW.
type SizingMinMax struct {
	Min float32 // The smallest final size of the element on this axis will be this value in pixels.
	Max float32 // The largest final size of the element on this axis will be this value in pixels.
}

// Controls the sizing of this element along one axis inside its parent container.
//type SizingAxis struct {
//	// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
//	// Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
//	data SizingMinMax
//	Type SizingType // Controls how the element takes up space inside its parent container.
//}

type SizingAxisMinMax interface {
	GetMinMax() SizingMinMax
}

// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
type SizingAxisFixed struct {
	MinMax SizingMinMax
}

func (s SizingAxisFixed) GetMinMax() SizingMinMax {
	return s.MinMax
}

// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
type SizingAxisFit struct {
	MinMax SizingMinMax
}

func (s SizingAxisFit) GetMinMax() SizingMinMax {
	return s.MinMax
}

// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
type SizingAxisGrow struct {
	MinMax SizingMinMax
}

func (s SizingAxisGrow) GetMinMax() SizingMinMax {
	return s.MinMax
}

// Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
type SizingAxisPercent struct {
	Percent float32
}

type AnySizingAxis any

func SIZING_FIT(s ...float32) AnySizingAxis {
	switch len(s) {
	case 0:
		return SizingAxisFit{}
	case 1:
		return SizingAxisFit{MinMax: SizingMinMax{Min: s[0]}}
	default:
		return SizingAxisFit{MinMax: SizingMinMax{Min: s[0], Max: s[1]}}
	}
}

func SIZING_GROW(s ...float32) AnySizingAxis {
	switch len(s) {
	case 0:
		return SizingAxisGrow{}
	case 1:
		return SizingAxisGrow{MinMax: SizingMinMax{Min: s[0]}}
	default:
		return SizingAxisGrow{MinMax: SizingMinMax{Min: s[0], Max: s[1]}}
	}
}

func SIZING_FIXED(fixedSize float32) AnySizingAxis {
	return SizingAxisFixed{MinMax: SizingMinMax{Min: fixedSize, Max: fixedSize}}
}

func SIZING_PERCENT(percentOfParent float32) AnySizingAxis {
	return SizingAxisPercent{Percent: percentOfParent}
}

// Controls the sizing of this element along one axis inside its parent container.
type Sizing struct {
	Width  AnySizingAxis // Controls the width sizing of the element, along the x axis.
	Height AnySizingAxis // Controls the height sizing of the element, along the y axis.
}

func (s Sizing) GetAxis(xAxis bool) AnySizingAxis {
	if xAxis {
		return s.Width
	} else {
		return s.Height
	}
}

// Controls "padding" in pixels, which is a gap between the bounding box of this element and where its children
// will be placed.
type Padding struct {
	Left   uint16
	Right  uint16
	Top    uint16
	Bottom uint16
}

func PADDING_ALL(padding uint16) Padding {
	return Padding{padding, padding, padding, padding}
}

// Controls various settings that affect the size and position of an element, as well as the sizes and positions
// of any child elements.
type LayoutConfig struct {
	Sizing          Sizing          // Controls the sizing of this element inside it's parent container, including FIT, GROW, PERCENT and FIXED sizing.
	Padding         Padding         // Controls "padding" in pixels, which is a gap between the bounding box of this element and where its children will be placed.
	ChildGap        uint16          // Controls the gap in pixels between child elements along the layout axis (horizontal gap for LEFT_TO_RIGHT, vertical gap for TOP_TO_BOTTOM).
	ChildAlignment  ChildAlignment  // Controls how child elements are aligned on each axis.
	LayoutDirection LayoutDirection // Controls the direction in which child elements will be automatically laid out.
}

var default_LayoutConfig LayoutConfig

// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
type TextElementConfigWrapMode uint8

const (
	// (default) breaks on whitespace characters.
	TEXT_WRAP_WORDS TextElementConfigWrapMode = iota
	// Don't break on space characters, only on newlines.
	TEXT_WRAP_NEWLINES
	// Disable text wrapping entirely.
	TEXT_WRAP_NONE
)

// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
type TextAlignment uint8

const (
	// (default) Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
	TEXT_ALIGN_LEFT TextAlignment = iota
	// Horizontally aligns wrapped lines of text to the center of their bounding box.
	TEXT_ALIGN_CENTER
	// Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
	TEXT_ALIGN_RIGHT
)

// Controls various functionality related to text elements.
type TextElementConfig struct {
	// A pointer that will be transparently passed through to the resulting render command.
	userData any
	// The RGBA color of the font to render, conventionally specified as 0-255.
	textColor colorex.RGBA
	// An integer transparently passed to MeasureText to identify the font to use.
	// The debug view will pass fontId = 0 for its internal text.
	fontId uint16
	// Controls the size of the font. Handled by the function provided to MeasureText.
	fontSize uint16
	// Controls extra horizontal spacing between characters. Handled by the function provided to MeasureText.
	letterSpacing uint16
	// Controls additional vertical space between wrapped lines of text.
	lineHeight uint16
	// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
	// TEXT_WRAP_WORDS (default) breaks on whitespace characters.
	// TEXT_WRAP_NEWLINES doesn't break on space characters, only on newlines.
	// TEXT_WRAP_NONE disables wrapping entirely.
	wrapMode TextElementConfigWrapMode
	// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
	// TEXT_ALIGN_LEFT (default) - Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
	// TEXT_ALIGN_CENTER - Horizontally aligns wrapped lines of text to the center of their bounding box.
	// TEXT_ALIGN_RIGHT - Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
	textAlignment TextAlignment
	// When set to true, clay will hash the entire text contents of this string as an identifier for its internal
	// text measurement cache, rather than just the pointer and length. This will incur significant performance cost for
	// long bodies of text.
	hashStringContents bool
}

var default_TextElementConfig TextElementConfig

// Image --------------------------------

// Controls various settings related to image elements.
type ImageElementConfig struct {
	ImageData        any             // A transparent pointer used to pass image data through to the renderer.
	SourceDimensions vector2.Float32 // The original dimensions of the source image, used to control aspect ratio.
}

var default_ImageElementConfig ImageElementConfig

// Floating -----------------------------

// Controls where a floating element is offset relative to its parent element.
// Note: see https://github.com/user-attachments/assets/b8c6dfaa-c1b1-41a4-be55-013473e4a6ce for a visual explanation.
type FloatingAttachPointType uint8

const (
	ATTACH_POINT_LEFT_TOP FloatingAttachPointType = iota
	ATTACH_POINT_LEFT_CENTER
	ATTACH_POINT_LEFT_BOTTOM
	ATTACH_POINT_CENTER_TOP
	ATTACH_POINT_CENTER_CENTER
	ATTACH_POINT_CENTER_BOTTOM
	ATTACH_POINT_RIGHT_TOP
	ATTACH_POINT_RIGHT_CENTER
	ATTACH_POINT_RIGHT_BOTTOM
)

// Controls where a floating element is offset relative to its parent element.
type FloatingAttachPoints struct {
	element FloatingAttachPointType // Controls the origin point on a floating element that attaches to its parent.
	parent  FloatingAttachPointType // Controls the origin point on the parent element that the floating element attaches to.
}

// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
type PointerCaptureMode uint8

const (
	// (default) "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.
	POINTER_CAPTURE_MODE_CAPTURE PointerCaptureMode = iota
	//    POINTER_CAPTURE_MODE_PARENT, TODO pass pointer through to attached parent

	// Transparently pass through pointer events like hover and click to elements underneath the floating element.
	POINTER_CAPTURE_MODE_PASSTHROUGH
)

// Controls which element a floating element is "attached" to (i.e. relative offset from).
type FloatingAttachToElement uint8

const (
	// (default) Disables floating for this element.
	ATTACH_TO_NONE FloatingAttachToElement = iota
	// Attaches this floating element to its parent, positioned based on the .attachPoints and .offset fields.
	ATTACH_TO_PARENT
	// Attaches this floating element to an element with a specific ID, specified with the .parentId field. positioned based on the .attachPoints and .offset fields.
	ATTACH_TO_ELEMENT_WITH_ID
	// Attaches this floating element to the root of the layout, which combined with the .offset field provides functionality similar to "absolute positioning".
	ATTACH_TO_ROOT
)

// Controls various settings related to "floating" elements, which are elements that "float" above other elements, potentially overlapping their boundaries,
// and not affecting the layout of sibling or parent elements.
type FloatingElementConfig struct {
	// Offsets this floating element by the provided x,y coordinates from its attachPoints.
	offset vector2.Float32
	// Expands the boundaries of the outer floating element without affecting its children.
	expand vector2.Float32
	// When used in conjunction with .attachTo = ATTACH_TO_ELEMENT_WITH_ID, attaches this floating element to the element in the hierarchy with the provided ID.
	// Hint: attach the ID to the other element with .id = ID("yourId"), and specify the id the same way, with .parentId = ID("yourId").id
	parentId uint32
	// Controls the z index of this floating element and all its children. Floating elements are sorted in ascending z order before output.
	// zIndex is also passed to the renderer for all elements contained within this floating element.
	zIndex int16
	// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath / behind a floating element.
	// Enum is of the form ATTACH_POINT_foo_bar. See Clay_FloatingAttachPoints for more details.
	// Note: see <img src="https://github.com/user-attachments/assets/b8c6dfaa-c1b1-41a4-be55-013473e4a6ce />
	// and <img src="https://github.com/user-attachments/assets/ebe75e0d-1904-46b0-982d-418f929d1516 /> for a visual explanation.
	attachPoints FloatingAttachPoints
	// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
	// POINTER_CAPTURE_MODE_CAPTURE (default) - "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.
	// POINTER_CAPTURE_MODE_PASSTHROUGH - Transparently pass through pointer events like hover and click to elements underneath the floating element.
	pointerCaptureMode PointerCaptureMode
	// Controls which element a floating element is "attached" to (i.e. relative offset from).
	// ATTACH_TO_NONE (default) - Disables floating for this element.
	// ATTACH_TO_PARENT - Attaches this floating element to its parent, positioned based on the .attachPoints and .offset fields.
	// ATTACH_TO_ELEMENT_WITH_ID - Attaches this floating element to an element with a specific ID, specified with the .parentId field. positioned based on the .attachPoints and .offset fields.
	// ATTACH_TO_ROOT - Attaches this floating element to the root of the layout, which combined with the .offset field provides functionality similar to "absolute positioning".
	attachTo FloatingAttachToElement
}

var default_FloatingElementConfig FloatingElementConfig

// Custom -----------------------------

// Controls various settings related to custom elements.
type CustomElementConfig struct {
	// A transparent pointer through which you can pass custom data to the renderer.
	// Generates CUSTOM render commands.
	customData any
}

var default_CustomElementConfig CustomElementConfig

// Scroll -----------------------------

// Controls the axis on which an element switches to "scrolling", which clips the contents and allows scrolling in that direction.
type ScrollElementConfig struct {
	horizontal bool // Clip overflowing elements on the X axis and allow scrolling left and right.
	vertical   bool // Clip overflowing elements on the YU axis and allow scrolling up and down.
}

var default_ScrollElementConfig ScrollElementConfig

// Border -----------------------------

// Controls the widths of individual element borders.
type BorderWidth struct {
	left   uint16
	right  uint16
	top    uint16
	bottom uint16
	// Creates borders between each child element, depending on the .layoutDirection.
	// e.g. for LEFT_TO_RIGHT, borders will be vertical lines, and for TOP_TO_BOTTOM borders will be horizontal lines.
	// .betweenChildren borders will result in individual RECTANGLE render commands being generated.
	betweenChildren uint16
}

func (b BorderWidth) IsEmpty() bool {
	return b.left == 0 &&
		b.right == 0 &&
		b.top == 0 &&
		b.bottom == 0 &&
		b.betweenChildren == 0
}

// Controls settings related to element borders.
type BorderElementConfig struct {
	color colorex.RGBA // Controls the color of all borders with width > 0. Conventionally represented as 0-255, but interpretation is up to the renderer.
	width BorderWidth  // Controls the widths of individual borders. At least one of these should be > 0 for a BORDER render command to be generated.
}

var default_BorderElementConfig BorderElementConfig

func (b BorderElementConfig) IsEmpty() bool {
	return b.color.IsZero() && b.width.IsEmpty()
}

// Render Command Data -----------------------------

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
type TextRenderData struct {
	// A string slice containing the text to be rendered.
	// Note: this is not guaranteed to be null terminated.
	stringContents string
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	textColor colorex.RGBA
	// An integer representing the font to use to render this text, transparently passed through from the text declaration.
	fontId   uint16
	fontSize uint16
	// Specifies the extra whitespace gap in pixels between each character.
	letterSpacing uint16
	// The height of the bounding box for this line of text.
	lineHeight uint16
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
type RectangleRenderData struct {
	// The solid background color to fill this rectangle with. Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	backgroundColor colorex.RGBA
	// Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	cornerRadius CornerRadius
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
type ImageRenderData struct {
	// The tint color for this image. Note that the default value is 0,0,0,0 and should likely be interpreted
	// as "untinted".
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	backgroundColor colorex.RGBA
	// Controls the "radius", or corner rounding of this image.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	cornerRadius CornerRadius
	// The original dimensions of the source image, used to control aspect ratio.
	sourceDimensions vector2.Float32
	// A pointer transparently passed through from the original element definition, typically used to represent image data.
	imageData any
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
type CustomRenderData struct {
	// Passed through from .backgroundColor in the original element declaration.
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	backgroundColor colorex.RGBA
	// Controls the "radius", or corner rounding of this custom element.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	cornerRadius CornerRadius
	// A pointer transparently passed through from the original element definition.
	customData any
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START || commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
type ScrollRenderData struct {
	horizontal bool
	vertical   bool
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
type BorderRenderData struct {
	// Controls a shared color for all this element's borders.
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	color colorex.RGBA
	// Specifies the "radius", or corner rounding of this border element.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	cornerRadius CornerRadius
	// Controls individual border side widths.
	width BorderWidth
}

type ScissorsStartData struct {
	horizontal bool
	vertical   bool
}
type ScissorsEndData struct {
	horizontal bool
	vertical   bool
}

type RenderDataType interface {
	RectangleRenderData | TextRenderData | ImageRenderData | CustomRenderData | BorderRenderData | ScrollRenderData | ScissorsStartData | ScissorsEndData
}

type AnyRenderData any

type RenderCommand struct {
	// A rectangular box that fully encloses this UI element, with the position relative to the root of the layout.
	BoundingBox rect2.Float32
	// A struct union containing data specific to this command's commandType.
	RenderData AnyRenderData
	// A pointer transparently passed through from the original element declaration.
	UserData any
	// The Id of this element, transparently passed through from the original element declaration.
	Id uint32
	// The z order required for drawing this command correctly.
	// Note: the render command array is already sorted in ascending order, and will produce correct results if drawn in naive order.
	// This field is intended for use in batching renderers for improved performance.
	ZIndex int16
}

// Represents the current state of interaction with clay this frame.
type PointerDataInteractionState uint8

const (
	// A left mouse click, or touch occurred this frame.
	POINTER_DATA_PRESSED_THIS_FRAME PointerDataInteractionState = iota
	// The left mouse button click or touch happened at some point in the past, and is still currently held down this frame.
	POINTER_DATA_PRESSED
	// The left mouse button click or touch was released this frame.
	POINTER_DATA_RELEASED_THIS_FRAME
	// The left mouse button click or touch is not currently down / was released at some point in the past.
	POINTER_DATA_RELEASED
)

// Information on the current state of pointer interactions this frame.
type PointerData struct {
	// The position of the mouse / touch / pointer relative to the root of the layout.
	position vector2.Float32
	// Represents the current state of interaction with clay this frame.
	// POINTER_DATA_PRESSED_THIS_FRAME - A left mouse click, or touch occurred this frame.
	// POINTER_DATA_PRESSED - The left mouse button click or touch happened at some point in the past, and is still currently held down this frame.
	// POINTER_DATA_RELEASED_THIS_FRAME - The left mouse button click or touch was released this frame.
	// POINTER_DATA_RELEASED - The left mouse button click or touch is not currently down / was released at some point in the past.
	state PointerDataInteractionState
}

type ElementDeclaration struct {
	// Primarily created via the ID(), IDI(), ID_LOCAL() and IDI_LOCAL() macros.
	// Represents a hashed string ID used for identifying and finding specific clay UI elements, required by functions such as PointerOver() and GetElementData().
	Id ElementId
	// Controls various settings that affect the size and position of an element, as well as the sizes and positions of any child elements.
	Layout LayoutConfig
	// Controls the background color of the resulting element.
	// By convention specified as 0-255, but interpretation is up to the renderer.
	// If no other config is specified, .BackgroundColor will generate a RECTANGLE render command, otherwise it will be passed as a property to IMAGE or CUSTOM render commands.
	BackgroundColor colorex.RGBA
	// Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
	CornerRadius CornerRadius
	// Controls settings related to Image elements.
	Image ImageElementConfig
	// Controls whether and how an element "floats", which means it layers over the top of other elements in z order, and doesn't affect the position and size of siblings or parent elements.
	// Note: in order to activate Floating, .Floating.attachTo must be set to something other than the default value.
	Floating FloatingElementConfig
	// Used to create CUSTOM render commands, usually to render element types not supported by Clay.
	Custom CustomElementConfig
	// Controls whether an element should clip its contents and allow scrolling rather than expanding to contain them.
	Scroll ScrollElementConfig
	// Controls settings related to element borders, and will generate BORDER render commands.
	Border BorderElementConfig
	// A pointer that will be transparently passed through to resulting render commands.
	UserData any
}

// Represents the type of error clay encountered while computing layout.
type ErrorType uint8

const (
	// A text measurement function wasn't provided using Clay_SetMeasureTextFunction(), or the provided function was null.
	ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED ErrorType = iota
	// Clay attempted to allocate its internal data structures but ran out of space.
	// The arena passed to Clay_Initialize was created with a capacity smaller than that required by Clay_MinMemorySize().
	ERROR_TYPE_ARENA_CAPACITY_EXCEEDED
	// Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxElementCount().
	ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED
	// Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxMeasureTextCacheWordCount().
	ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED
	// Two elements were declared with exactly the same ID within one layout.
	ERROR_TYPE_DUPLICATE_ID
	// A floating element was declared using ATTACH_TO_ELEMENT_ID and either an invalid .parentId was provided or no element with the provided .parentId was found.
	ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND
	// An element was declared that using SIZING_PERCENT but the percentage value was over 1. Percentage values are expected to be in the 0-1 range.
	ERROR_TYPE_PERCENTAGE_OVER_1
	// Clay encountered an internal error. It would be wonderful if you could report this so we can fix it!
	ERROR_TYPE_INTERNAL_ERROR
)

// Data to identify the error that clay has encountered.
type ErrorData struct {
	// Represents the type of error clay encountered while computing layout.
	// ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED - A text measurement function wasn't provided using Clay_SetMeasureTextFunction(), or the provided function was null.
	// ERROR_TYPE_ARENA_CAPACITY_EXCEEDED - Clay attempted to allocate its internal data structures but ran out of space. The arena passed to Clay_Initialize was created with a capacity smaller than that required by Clay_MinMemorySize().
	// ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED - Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxElementCount().
	// ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED - Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxMeasureTextCacheWordCount().
	// ERROR_TYPE_DUPLICATE_ID - Two elements were declared with exactly the same ID within one layout.
	// ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND - A floating element was declared using ATTACH_TO_ELEMENT_ID and either an invalid .parentId was provided or no element with the provided .parentId was found.
	// ERROR_TYPE_PERCENTAGE_OVER_1 - An element was declared that using SIZING_PERCENT but the percentage value was over 1. Percentage values are expected to be in the 0-1 range.
	// ERROR_TYPE_INTERNAL_ERROR - Clay encountered an internal error. It would be wonderful if you could report this so we can fix it!
	ErrorType ErrorType
	// A string containing human-readable error text that explains the error in more detail.
	ErrorText string
	// A transparent pointer passed through from when the error handler was first provided.
	UserData any
}

// A wrapper struct around Clay's error handler function.
type ErrorHandler struct {
	// A user provided function to call when Clay encounters an error during layout.
	ErrorHandlerFunction func(errorText ErrorData)
	// A pointer that will be transparently passed through to the error handler when it is called.
	UserData any
}
