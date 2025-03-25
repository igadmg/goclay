package goclay

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"math"

	"github.com/igadmg/goex/slicesex"
	"github.com/igadmg/raylib-go/raymath/vector2"
)

var LAYOUT_DEFAULT LayoutConfig
var Color_DEFAULT Color
var CornerRadius_DEFAULT CornerRadius
var BorderWidth_DEFAULT BorderWidth

var currentContext *Context = nil
var defaultMaxElementCount int32 = 8192
var defaultMaxMeasureTextWordCacheCount int32 = 16384

func errorHandlerFunctionDefault(errorText ErrorData) {
}

var SPACECHAR string = " "
var STRING_DEFAULT string = ""

type BooleanWarnings struct {
	maxElementsExceeded           bool
	maxRenderCommandsExceeded     bool
	maxTextMeasureCacheExceeded   bool
	textMeasurementFunctionNotSet bool
}

type Warning struct {
	baseMessage    string
	dynamicMessage string
}

var default_SharedElementConfig SharedElementConfig

type SharedElementConfig struct {
	backgroundColor Color
	cornerRadius    CornerRadius
	userData        any
}

//	union {
//	   *TextElementConfig
//	   *ImageElementConfig
//	   *FloatingElementConfig
//	   *CustomElementConfig
//	   *ScrollElementConfig
//	   *BorderElementConfig
//	   *SharedElementConfig
//	}
type ElementConfigType interface {
	*TextElementConfig | *ImageElementConfig | *FloatingElementConfig | *CustomElementConfig | *ScrollElementConfig | *BorderElementConfig | *SharedElementConfig
}

type AnyElementConfig any

type WrappedTextLine struct {
	dimensions Dimensions
	line       string
}

type TextElementData struct {
	text                string
	preferredDimensions Dimensions
	elementIndex        int
	wrappedLines        []WrappedTextLine
}

type LayoutElement struct {
	//union {
	children        []int
	textElementData *TextElementData
	//}
	dimensions     Dimensions
	minDimensions  Dimensions
	layoutConfig   *LayoutConfig
	elementConfigs []AnyElementConfig
	id             uint32
}

type ScrollContainerDataInternal struct {
	layoutElement       *LayoutElement
	boundingBox         BoundingBox
	contentSize         Dimensions
	scrollOrigin        Vector2
	pointerOrigin       Vector2
	scrollMomentum      Vector2
	scrollPosition      Vector2
	previousDelta       Vector2
	momentumTime        float32
	elementId           uint32
	openThisFrame       bool
	pointerScrollActive bool
}

type DebugElementData struct {
	collision bool
	collapsed bool
}

type LayoutElementHashMapItem struct { // todo get this struct into a single cache line
	boundingBox           BoundingBox
	elementId             ElementId
	layoutElement         *LayoutElement
	onHoverFunction       func(elementId ElementId, pointerInfo PointerData, userData any)
	hoverFunctionUserData any
	nextIndex             int32
	generation            uint32
	idAlias               uint32
	debugData             DebugElementData
}

var default_LayoutElementHashMapItem LayoutElementHashMapItem

func (i LayoutElementHashMapItem) IsEmpty() bool {
	return i.elementId.id == 0
}

type MeasuredWord struct {
	startOffset int32
	length      int32
	width       float32
	next        int32
}

type MeasureTextCacheItem struct {
	unwrappedDimensions     Dimensions
	measuredWordsStartIndex int32
	minWidth                float32
	containsNewlines        bool
	// Hash map data
	id         uint32
	nextIndex  int32
	generation uint32
}

var default_MeasureTextCacheItem MeasureTextCacheItem

type LayoutElementTreeNode struct {
	layoutElement   *LayoutElement
	position        Vector2
	nextChildOffset Vector2
}

type LayoutElementTreeRoot struct {
	layoutElementIndex int
	parentId           uint32 // This can be zero in the case of the root layout tree
	clipElementId      uint32 // This can be zero if there is no clip element
	zIndex             int16
	pointerOffset      Vector2 // Only used when scroll containers are managed externally
}

type Context struct {
	maxElementCount              int32
	maxMeasureTextCacheWordCount int32
	warningsEnabled              bool
	errorHandler                 ErrorHandler
	booleanWarnings              BooleanWarnings
	warnings                     []Warning

	pointerInfo                   PointerData
	layoutDimensions              Dimensions
	dynamicElementIndexBaseHash   ElementId
	dynamicElementIndex           uint32
	debugModeEnabled              bool
	disableCulling                bool
	externalScrollHandlingEnabled bool
	debugSelectedElementId        uint32
	generation                    uint32
	measureTextUserData           any
	queryScrollOffsetUserData     any

	// Layout Elements / Render Commands
	layoutElements              []LayoutElement
	renderCommands              []RenderCommand
	openLayoutElementStack      []int
	layoutElementChildren       []int
	layoutElementChildrenBuffer []int
	textElementData             []TextElementData
	imageElementPointers        []int
	reusableElementIndexBuffer  []int32
	layoutElementClipElementIds []int

	// Configs
	layoutConfigs          []LayoutConfig
	elementConfigs         []AnyElementConfig
	textElementConfigs     []TextElementConfig
	imageElementConfigs    []ImageElementConfig
	floatingElementConfigs []FloatingElementConfig
	scrollElementConfigs   []ScrollElementConfig
	customElementConfigs   []CustomElementConfig
	borderElementConfigs   []BorderElementConfig
	sharedElementConfigs   []SharedElementConfig

	// Misc Data Structures
	layoutElementIdStrings        []string
	wrappedTextLines              []WrappedTextLine
	layoutElementTreeNodeArray1   []LayoutElementTreeNode
	layoutElementTreeRoots        []LayoutElementTreeRoot
	layoutElementsHashMapInternal []LayoutElementHashMapItem
	layoutElementsHashMap         map[uint32]*LayoutElementHashMapItem
	measureTextHashMapInternal    []MeasureTextCacheItem
	measureTextHashMap            map[string]*MeasureTextCacheItem
	measuredWords                 []MeasuredWord
	measuredWordsFreeList         []int32
	openClipElementStack          []int
	pointerOverIds                []ElementId
	scrollContainerDatas          []ScrollContainerDataInternal
	dynamicStringData             []byte
	debugElementData              []DebugElementData
}

var measureText MeasureTextFn
var queryScrollOffset QueryScrollOffsetFn

func (c *Context) getOpenLayoutElement() *LayoutElement {
	return &c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-1]]
}

func (c *Context) getParentElementId() uint32 {
	return c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-2]].id
}

func (c *Context) storeLayoutConfig(config LayoutConfig) *LayoutConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_LayoutConfig
	}
	c.layoutConfigs = append(c.layoutConfigs, config)
	return &c.layoutConfigs[len(c.layoutConfigs)-1]
}

func (c *Context) storeTextElementConfig(config TextElementConfig) *TextElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_TextElementConfig
	}
	c.textElementConfigs = append(c.textElementConfigs, config)
	return &c.textElementConfigs[len(c.textElementConfigs)-1]
}

func (c *Context) storeImageElementConfig(config ImageElementConfig) *ImageElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_ImageElementConfig
	}
	c.imageElementConfigs = append(c.imageElementConfigs, config)
	return &c.imageElementConfigs[len(c.imageElementConfigs)-1]
}

func (c *Context) storeFloatingElementConfig(config FloatingElementConfig) *FloatingElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_FloatingElementConfig
	}
	c.floatingElementConfigs = append(c.floatingElementConfigs, config)
	return &c.floatingElementConfigs[len(c.floatingElementConfigs)-1]
}

func (c *Context) storeCustomElementConfig(config CustomElementConfig) *CustomElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_CustomElementConfig
	}
	c.customElementConfigs = append(c.customElementConfigs, config)
	return &c.customElementConfigs[len(c.customElementConfigs)-1]
}

func (c *Context) storeScrollElementConfig(config ScrollElementConfig) *ScrollElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_ScrollElementConfig
	}
	c.scrollElementConfigs = append(c.scrollElementConfigs, config)
	return &c.scrollElementConfigs[len(c.scrollElementConfigs)-1]
}

func (c *Context) storeBorderElementConfig(config BorderElementConfig) *BorderElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_BorderElementConfig
	}
	c.borderElementConfigs = append(c.borderElementConfigs, config)
	return &c.borderElementConfigs[len(c.borderElementConfigs)-1]
}

func (c *Context) storeSharedElementConfig(config SharedElementConfig) *SharedElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_SharedElementConfig
	}
	c.sharedElementConfigs = append(c.sharedElementConfigs, config)
	return &c.sharedElementConfigs[len(c.sharedElementConfigs)-1]
}

func (c *Context) attachElementConfig(config AnyElementConfig) AnyElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return config
	}
	openLayoutElement := c.getOpenLayoutElement()
	c.elementConfigs = append(c.elementConfigs, config)
	openLayoutElement.elementConfigs = append(openLayoutElement.elementConfigs, config)
	return config
}

func findElementConfigWithType[T ElementConfigType](element *LayoutElement) (T, bool) {
	for _, config := range element.elementConfigs {
		switch cfg := config.(type) {
		case T:
			return cfg, true
		}
	}
	return nil, false
}

var hasher hash.Hash32 = fnv.New32a()
var hashBuffer []byte = make([]byte, 4+2*4+1)

func hashNumber(offset uint32, seed uint32) ElementId {
	le := binary.LittleEndian
	na := hashBuffer[0:0]
	na = le.AppendUint32(na, offset)
	na = le.AppendUint32(na, seed)

	hasher.Reset()
	hasher.Write(hashBuffer[:2*4])
	hash := hasher.Sum32()

	return ElementId{
		id:       hash + 1, // Reserve the hash result of zero as "null id"
		stringId: STRING_DEFAULT,
	}
}

func hashString(key string) ElementId {
	hasher.Reset()
	hasher.Write([]byte(key))
	hash := hasher.Sum32()

	return ElementId{
		id:       hash + 1, // Reserve the hash result of zero as "null id"
		stringId: key,
	}
}

func hashTextWithConfig(text string, config *TextElementConfig) uint32 {
	hasher.Reset()
	hasher.Write([]byte(text))

	le := binary.LittleEndian
	na := hashBuffer[0:0]
	na = le.AppendUint32(na, uint32(len(text)))
	na = le.AppendUint16(na, config.FontId)
	na = le.AppendUint16(na, config.FontSize)
	na = le.AppendUint16(na, config.LineHeight)
	na = le.AppendUint16(na, config.LetterSpacing)
	na[0] = byte(config.WrapMode)
	hasher.Write(hashBuffer[:4+2*4+1])

	hash := hasher.Sum32()
	return hash + 1 // Reserve the hash result of zero as "null id"
}

func (c *Context) addMeasuredWord(word MeasuredWord, previousWord *MeasuredWord) *MeasuredWord {
	if len(c.measuredWordsFreeList) > 0 {
		newItemIndex := c.measuredWordsFreeList[len(c.measuredWordsFreeList)-1]
		c.measuredWordsFreeList = c.measuredWordsFreeList[:len(c.measuredWordsFreeList)-1]
		c.measuredWords = slicesex.Set(c.measuredWords, int(newItemIndex), word)
		previousWord.next = newItemIndex
		return &c.measuredWords[newItemIndex]
	} else {
		previousWord.next = int32(len(c.measuredWords))
		c.measuredWords = append(c.measuredWords, word)
		return &c.measuredWords[len(c.measuredWords)-1]
	}
}

func (c *Context) measureTextCached(text string, config *TextElementConfig) *MeasureTextCacheItem {
	if measureText == nil {
		if !c.booleanWarnings.textMeasurementFunctionNotSet {
			c.booleanWarnings.textMeasurementFunctionNotSet = true
			c.errorHandler.ErrorHandlerFunction(ErrorData{
				ErrorType: ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED,
				ErrorText: "Clay's internal MeasureText function is null. You may have forgotten to call Clay_SetMeasureTextFunction(), or passed a NULL function pointer by mistake.",
				UserData:  c.errorHandler.UserData,
			})
		}
		return &default_MeasureTextCacheItem
	}

	id := hashTextWithConfig(text, config)
	if hashEntry, ok := c.measureTextHashMap[text]; ok {
		return hashEntry
	}

	newCacheItem := MeasureTextCacheItem{
		measuredWordsStartIndex: -1,
		id:                      id,
		generation:              c.generation,
	}
	measured := (*MeasureTextCacheItem)(nil)

	if len(c.measureTextHashMapInternal) == cap(c.measureTextHashMapInternal)-1 {
		if !c.booleanWarnings.maxTextMeasureCacheExceeded {
			c.errorHandler.ErrorHandlerFunction(ErrorData{
				ErrorType: ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
				ErrorText: "Clay ran out of capacity while attempting to measure text elements. Try using Clay_SetMaxElementCount() with a higher value.",
				UserData:  c.errorHandler.UserData})
			c.booleanWarnings.maxTextMeasureCacheExceeded = true
		}
		return &default_MeasureTextCacheItem
	}
	c.measureTextHashMapInternal = append(c.measureTextHashMapInternal, newCacheItem)
	measured = &c.measureTextHashMapInternal[len(c.measureTextHashMapInternal)-1]

	start := 0
	end := 0
	lineWidth := float32(0)
	measuredWidth := float32(0)
	measuredHeight := float32(0)
	spaceWidth := measureText(SPACECHAR, config, c.measureTextUserData).X
	tempWord := MeasuredWord{next: -1}
	previousWord := &tempWord
	for end < len(text) {
		if len(c.measuredWords) == cap(c.measuredWords)-1 {
			if !c.booleanWarnings.maxTextMeasureCacheExceeded {
				c.errorHandler.ErrorHandlerFunction(ErrorData{
					ErrorType: ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED,
					ErrorText: "Clay has run out of space in it's internal text measurement cache. Try using Clay_SetMaxMeasureTextCacheWordCount() (default 16384, with 1 unit storing 1 measured word).",
					UserData:  c.errorHandler.UserData,
				})
				c.booleanWarnings.maxTextMeasureCacheExceeded = true
			}
			return &default_MeasureTextCacheItem
		}
		current := text[end]
		if current == ' ' || current == '\n' {
			length := end - start
			dimensions := measureText(text[start:end], config, c.measureTextUserData)
			measured.minWidth = max(dimensions.X, measured.minWidth)
			measuredHeight = max(float32(measuredHeight), dimensions.Y)
			if current == ' ' {
				dimensions.X += spaceWidth
				previousWord = c.addMeasuredWord(MeasuredWord{
					startOffset: int32(start),
					length:      int32(length + 1),
					width:       dimensions.X,
					next:        -1},
					previousWord)
				lineWidth += dimensions.X
			}
			if current == '\n' {
				if length > 0 {
					previousWord = c.addMeasuredWord(MeasuredWord{
						startOffset: int32(start),
						length:      int32(length),
						width:       dimensions.X,
						next:        -1},
						previousWord)
				}
				previousWord = c.addMeasuredWord(MeasuredWord{
					startOffset: int32(end + 1),
					length:      0,
					width:       0,
					next:        -1},
					previousWord)
				lineWidth += dimensions.X
				measuredWidth = max(lineWidth, measuredWidth)
				measured.containsNewlines = true
				lineWidth = 0
			}
			start = end + 1
		}
		end++
	}
	if end-start > 0 {
		dimensions := measureText(text[start:end], config, c.measureTextUserData)
		c.addMeasuredWord(MeasuredWord{
			startOffset: int32(start),
			length:      int32(end - start),
			width:       dimensions.X,
			next:        -1,
		},
			previousWord)
		lineWidth += dimensions.X
		measuredHeight = max(measuredHeight, dimensions.Y)
		measured.minWidth = max(dimensions.X, measured.minWidth)
	}
	measuredWidth = max(lineWidth, measuredWidth)

	measured.measuredWordsStartIndex = tempWord.next
	measured.unwrappedDimensions.X = measuredWidth
	measured.unwrappedDimensions.Y = measuredHeight

	return measured
}

func (c *Context) addHashMapItem(elementId ElementId, layoutElement *LayoutElement, idAlias uint32) *LayoutElementHashMapItem {
	if len(c.layoutElementsHashMapInternal) == cap(c.layoutElementsHashMapInternal)-1 {
		return nil
	}

	item := LayoutElementHashMapItem{
		elementId:     elementId,
		layoutElement: layoutElement,
		nextIndex:     -1,
		generation:    c.generation + 1,
		idAlias:       idAlias,
	}

	c.layoutElementsHashMapInternal = append(c.layoutElementsHashMapInternal, item)
	c.layoutElementsHashMap[elementId.id] = &c.layoutElementsHashMapInternal[len(c.layoutElementsHashMapInternal)-1]

	return c.layoutElementsHashMap[elementId.id]
}

func (c *Context) getHashMapItem(id uint32) *LayoutElementHashMapItem {
	r, ok := c.layoutElementsHashMap[id]
	if !ok {
		return &default_LayoutElementHashMapItem
	}

	return r
}

func (c *Context) generateIdForAnonymousElement(openLayoutElement *LayoutElement) ElementId {
	parentElement := c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-2]]
	elementId := hashNumber(uint32(len(parentElement.children)), parentElement.id)
	openLayoutElement.id = elementId.id
	c.addHashMapItem(elementId, openLayoutElement, 0)
	c.layoutElementIdStrings = append(c.layoutElementIdStrings, elementId.stringId)
	return elementId
}

func elementHasConfig[T ElementConfigType](layoutElement *LayoutElement) bool {
	for _, config := range layoutElement.elementConfigs {
		switch config.(type) {
		case T:
			return true
		}
	}
	return false
}

func updateAspectRatioBox(layoutElement *LayoutElement) {
	for _, config := range layoutElement.elementConfigs {
		switch c := config.(type) {
		case *ImageElementConfig:
			if c.SourceDimensions.X == 0 || c.SourceDimensions.Y == 0 {
				break
			}
			aspect := c.SourceDimensions.X / c.SourceDimensions.Y
			if layoutElement.dimensions.X == 0 && layoutElement.dimensions.Y != 0 {
				layoutElement.dimensions.X = layoutElement.dimensions.Y * aspect
			} else if layoutElement.dimensions.X != 0 && layoutElement.dimensions.Y == 0 {
				layoutElement.dimensions.Y = layoutElement.dimensions.Y * (1 / aspect)
			}
		}
	}
}

func (c *Context) closeElement() {
	if c.booleanWarnings.maxElementsExceeded {
		return
	}

	openLayoutElement := c.getOpenLayoutElement()
	layoutConfig := openLayoutElement.layoutConfig
	elementHasScrollHorizontal := false
	elementHasScrollVertical := false

	for _, config := range openLayoutElement.elementConfigs {
		switch cfg := config.(type) {
		case *ScrollElementConfig:
			elementHasScrollHorizontal = cfg.Horizontal
			elementHasScrollVertical = cfg.Vertical
			c.openClipElementStack = c.openClipElementStack[:len(c.openClipElementStack)-1]
			break
		case *FloatingElementConfig:
		}
	}

	leftRightPadding := float32(layoutConfig.Padding.Left + layoutConfig.Padding.Right)
	topBottomPadding := float32(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)

	// Attach children to the current open element
	c.layoutElementChildren = c.layoutElementChildren[0 : len(c.layoutElementChildren)+len(openLayoutElement.children)]
	openLayoutElement.children = c.layoutElementChildren[len(c.layoutElementChildren) : len(c.layoutElementChildren)+len(openLayoutElement.children)]
	if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
		openLayoutElement.dimensions.X = leftRightPadding
		openLayoutElement.minDimensions.X = leftRightPadding
		for i := range openLayoutElement.children {
			childIndex := c.layoutElementChildrenBuffer[len(c.layoutElementChildrenBuffer)-len(openLayoutElement.children)+i]
			child := c.layoutElements[childIndex]
			openLayoutElement.dimensions.X += child.dimensions.X
			openLayoutElement.dimensions.Y = max(
				openLayoutElement.dimensions.Y,
				child.dimensions.Y+topBottomPadding)

			// Minimum size of child elements doesn't matter to scroll containers as they can shrink and hide their contents
			if !elementHasScrollHorizontal {
				openLayoutElement.minDimensions.X += child.minDimensions.X
			}
			if !elementHasScrollVertical {
				openLayoutElement.minDimensions.Y = max(
					openLayoutElement.minDimensions.Y,
					child.minDimensions.Y+topBottomPadding)
			}
			c.layoutElementChildren = append(c.layoutElementChildren, childIndex)
		}

		childGap := float32(max(len(openLayoutElement.children)-1, 0) * int(layoutConfig.ChildGap))
		openLayoutElement.dimensions.X += childGap // TODO this is technically a bug with childgap and scroll containers
		openLayoutElement.minDimensions.X += childGap
	} else if layoutConfig.LayoutDirection == TOP_TO_BOTTOM {
		openLayoutElement.dimensions.Y = topBottomPadding
		openLayoutElement.minDimensions.Y = topBottomPadding
		for i := range openLayoutElement.children {
			childIndex := c.layoutElementChildrenBuffer[len(c.layoutElementChildrenBuffer)-len(openLayoutElement.children)+i]
			child := c.layoutElements[childIndex]
			openLayoutElement.dimensions.Y += child.dimensions.Y
			openLayoutElement.dimensions.X = max(
				openLayoutElement.dimensions.X,
				child.dimensions.X+leftRightPadding)
			// Minimum size of child elements doesn't matter to scroll containers as they can shrink and hide their contents
			if !elementHasScrollVertical {
				openLayoutElement.minDimensions.Y += child.minDimensions.Y
			}
			if !elementHasScrollHorizontal {
				openLayoutElement.minDimensions.X = max(
					openLayoutElement.minDimensions.X,
					child.minDimensions.X+leftRightPadding)
			}
			c.layoutElementChildren = append(c.layoutElementChildren, childIndex)
		}
		childGap := float32(max(len(openLayoutElement.children)-1, 0) * int(layoutConfig.ChildGap))
		openLayoutElement.dimensions.Y += childGap // TODO this is technically a bug with childgap and scroll containers
		openLayoutElement.minDimensions.Y += childGap
	}

	c.layoutElementChildrenBuffer = c.layoutElementChildrenBuffer[:len(c.layoutElementChildrenBuffer)-len(openLayoutElement.children)]

	// Clamp element min and max width to the values configured in the layout
	switch w := layoutConfig.Sizing.Width.(type) {
	case SizingAxisMinMax:
		mm := w.GetMinMax()
		openLayoutElement.dimensions.X = min(max(openLayoutElement.dimensions.X, mm.Min), mm.Max)
		openLayoutElement.minDimensions.X = min(max(openLayoutElement.minDimensions.X, mm.Min), mm.Max)
	default:
		openLayoutElement.dimensions.X = 0
	}

	// Clamp element min and max height to the values configured in the layout
	switch h := layoutConfig.Sizing.Height.(type) {
	case SizingAxisMinMax:
		mm := h.GetMinMax()
		openLayoutElement.dimensions.Y = min(max(openLayoutElement.dimensions.Y, mm.Min), mm.Max)
		openLayoutElement.minDimensions.Y = min(max(openLayoutElement.minDimensions.Y, mm.Min), mm.Max)
	default:
		openLayoutElement.dimensions.Y = 0
	}

	updateAspectRatioBox(openLayoutElement)

	elementIsFloating := elementHasConfig[*FloatingElementConfig](openLayoutElement)
	c.closeCurrentElement(elementIsFloating)
}

func (c *Context) closeCurrentElement(elementIsFloating bool) {
	// Close the currently open element
	var closingElementIndex int
	c.openLayoutElementStack, closingElementIndex = slicesex.RemoveSwapback(c.openLayoutElementStack, len(c.openLayoutElementStack)-1)
	openLayoutElement := c.getOpenLayoutElement()

	if !elementIsFloating && len(c.openLayoutElementStack) > 1 {
		c.layoutElementChildren = append(c.layoutElementChildren, closingElementIndex)
		openLayoutElement.children = append(openLayoutElement.children, closingElementIndex)
		c.layoutElementChildrenBuffer = append(c.layoutElementChildrenBuffer, closingElementIndex)
	}
}

func (c *Context) openElement() bool {
	if len(c.layoutElements) == cap(c.layoutElements)-1 || c.booleanWarnings.maxElementsExceeded {
		c.booleanWarnings.maxElementsExceeded = true
		return false
	}

	c.layoutElements = append(c.layoutElements, LayoutElement{})
	c.openLayoutElementStack = append(c.openLayoutElementStack, len(c.layoutElements)-1)
	if len(c.openClipElementStack) > 0 {
		c.layoutElementClipElementIds = slicesex.Set(
			c.layoutElementClipElementIds,
			len(c.layoutElements)-1,
			c.openClipElementStack[len(c.openClipElementStack)-1])
	} else {
		c.layoutElementClipElementIds = slicesex.Set(
			c.layoutElementClipElementIds,
			len(c.layoutElements)-1,
			0)
	}

	return true
}

func (c *Context) openTextElement(text string, textConfig *TextElementConfig) {
	if len(c.layoutElements) == cap(c.layoutElements)-1 || c.booleanWarnings.maxElementsExceeded {
		c.booleanWarnings.maxElementsExceeded = true
		return
	}
	parentElement := c.getOpenLayoutElement()

	c.layoutElements = append(c.layoutElements, LayoutElement{})
	c.openLayoutElementStack = append(c.openLayoutElementStack, len(c.layoutElements)-1)
	textElement := &c.layoutElements[len(c.layoutElements)-1]
	if len(c.openClipElementStack) > 0 {
		c.layoutElementClipElementIds = slicesex.Set(c.layoutElementClipElementIds, len(c.layoutElements)-1, c.openClipElementStack[len(c.openClipElementStack)-1])
	} else {
		c.layoutElementClipElementIds = slicesex.Set(c.layoutElementClipElementIds, len(c.layoutElements)-1, 0)
	}

	c.layoutElementChildrenBuffer = append(c.layoutElementChildrenBuffer, len(c.layoutElements)-1)
	textMeasured := c.measureTextCached(text, textConfig)
	elementId := hashNumber(uint32(len(parentElement.children)), parentElement.id)
	textElement.id = elementId.id
	c.addHashMapItem(elementId, textElement, 0)
	c.layoutElementIdStrings = append(c.layoutElementIdStrings, elementId.stringId)
	textDimensions := textMeasured.unwrappedDimensions
	if textConfig.LineHeight > 0 {
		textDimensions.Y = float32(textConfig.LineHeight)
	}
	textElement.dimensions = textDimensions
	textElement.minDimensions = MakeDimensions(textMeasured.minWidth, textDimensions.Y)
	c.textElementData = append(c.textElementData, TextElementData{
		text:                text,
		preferredDimensions: textMeasured.unwrappedDimensions,
		elementIndex:        len(c.layoutElements) - 1,
	})
	textElement.textElementData = &c.textElementData[len(c.textElementData)-1]
	c.elementConfigs = append(c.elementConfigs, textConfig)
	textElement.elementConfigs = c.elementConfigs[len(c.elementConfigs)-1 : len(c.elementConfigs)]
	textElement.layoutConfig = &default_LayoutConfig

	c.closeCurrentElement(false)
}

func (c *Context) attachId(elementId ElementId) ElementId {
	if c.booleanWarnings.maxElementsExceeded {
		return default_ElementId
	}
	openLayoutElement := c.getOpenLayoutElement()
	idAlias := openLayoutElement.id
	openLayoutElement.id = elementId.id
	c.addHashMapItem(elementId, openLayoutElement, idAlias)
	c.layoutElementIdStrings = append(c.layoutElementIdStrings, elementId.stringId)
	return elementId
}

func (c *Context) configureOpenElement(declaration *ElementDeclaration) {
	if declaration.Layout.Sizing.Width == nil {
		declaration.Layout.Sizing.Width = SIZING_FIT()
	}
	if declaration.Layout.Sizing.Height == nil {
		declaration.Layout.Sizing.Height = SIZING_FIT()
	}

	openLayoutElement := c.getOpenLayoutElement()
	openLayoutElement.layoutConfig = c.storeLayoutConfig(declaration.Layout)

	checkSizing := func(sizing Sizing) bool {
		switch w := sizing.Width.(type) {
		case SizingAxisPercent:
			if w.Percent > 1 {
				return true
			}
		}
		switch h := sizing.Height.(type) {
		case SizingAxisPercent:
			if h.Percent > 1 {
				return true
			}
		}
		return false
	}

	if checkSizing(declaration.Layout.Sizing) {
		c.errorHandler.ErrorHandlerFunction(ErrorData{
			ErrorType: ERROR_TYPE_PERCENTAGE_OVER_1,
			ErrorText: "An element was configured with SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2.",
			UserData:  c.errorHandler.UserData,
		})
	}

	openLayoutElementId := declaration.Id

	openLayoutElement.elementConfigs = c.elementConfigs[len(c.elementConfigs):len(c.elementConfigs)]
	sharedConfig := (*SharedElementConfig)(nil)
	if declaration.BackgroundColor.A > 0 {
		sharedConfig = c.storeSharedElementConfig(SharedElementConfig{backgroundColor: declaration.BackgroundColor})
		c.attachElementConfig(sharedConfig)
	}
	if !declaration.CornerRadius.IsEmpty() {
		if sharedConfig != nil {
			sharedConfig.cornerRadius = declaration.CornerRadius
		} else {
			sharedConfig = c.storeSharedElementConfig(SharedElementConfig{cornerRadius: declaration.CornerRadius})
			c.attachElementConfig(sharedConfig)
		}
	}
	if declaration.UserData != nil {
		if sharedConfig != nil {
			sharedConfig.userData = declaration.UserData
		} else {
			sharedConfig = c.storeSharedElementConfig(SharedElementConfig{userData: declaration.UserData})
			c.attachElementConfig(sharedConfig)
		}
	}
	if declaration.Image.ImageData != nil {
		c.attachElementConfig(c.storeImageElementConfig(declaration.Image))
		c.imageElementPointers = append(c.imageElementPointers, len(c.layoutElements)-1)
	}

	if declaration.Floating.AttachTo != ATTACH_TO_NONE {
		floatingConfig := declaration.Floating
		// This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
		hierarchicalParent := c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-2]]
		if true /*hierarchicalParent.id != 0*/ {
			clipElementId := 0
			if declaration.Floating.AttachTo == ATTACH_TO_PARENT {
				// Attach to the element's direct hierarchical parent
				floatingConfig.ParentId = hierarchicalParent.id
				if len(c.openClipElementStack) > 0 {
					clipElementId = c.openClipElementStack[len(c.openClipElementStack)-1]
				}
			} else if declaration.Floating.AttachTo == ATTACH_TO_ELEMENT_WITH_ID {
				parentItem := c.getHashMapItem(floatingConfig.ParentId)
				if parentItem == nil {
					c.errorHandler.ErrorHandlerFunction(ErrorData{
						ErrorType: ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
						ErrorText: "A floating element was declared with a parentId, but no element with that ID was found.",
						UserData:  c.errorHandler.UserData,
					})
				} else {
					// TODO: fix
					//clipElementId = c.layoutElementClipElementIds[(int32)(parentItem.layoutElement-c.layoutElements.internalArray))
				}
			} else if declaration.Floating.AttachTo == ATTACH_TO_ROOT {
				floatingConfig.ParentId = hashString("Clay__RootContainer").id
			}
			if openLayoutElementId.id == 0 {
				openLayoutElementId = IDI("Clay__FloatingContainer", uint32(len(c.layoutElementTreeRoots)))
			}
			currentElementIndex := c.openLayoutElementStack[len(c.openLayoutElementStack)-1]
			c.layoutElementClipElementIds = slicesex.Set(c.layoutElementClipElementIds, currentElementIndex, clipElementId)
			c.openClipElementStack = append(c.openClipElementStack, clipElementId)
			c.layoutElementTreeRoots = append(c.layoutElementTreeRoots, LayoutElementTreeRoot{
				layoutElementIndex: c.openLayoutElementStack[len(c.openLayoutElementStack)-1],
				parentId:           floatingConfig.ParentId,
				clipElementId:      uint32(clipElementId),
				zIndex:             floatingConfig.ZIndex,
			})
			c.attachElementConfig(c.storeFloatingElementConfig(floatingConfig))
		}
	}
	if declaration.Custom.CustomData != nil {
		c.attachElementConfig(c.storeCustomElementConfig(declaration.Custom))
	}

	if openLayoutElementId.id != 0 {
		c.attachId(openLayoutElementId)
	} else if openLayoutElement.id == 0 {
		openLayoutElementId = c.generateIdForAnonymousElement(openLayoutElement)
	}

	if declaration.Scroll.Horizontal || declaration.Scroll.Vertical {
		c.attachElementConfig(c.storeScrollElementConfig(declaration.Scroll))
		c.openClipElementStack = append(c.openClipElementStack, (int)(openLayoutElement.id))
		// Retrieve or create cached data to track scroll position across frames
		scrollOffset := (*ScrollContainerDataInternal)(nil)
		for i := range c.scrollContainerDatas {
			mapping := &c.scrollContainerDatas[i]
			if openLayoutElement.id == mapping.elementId {
				scrollOffset = mapping
				scrollOffset.layoutElement = openLayoutElement
				scrollOffset.openThisFrame = true
			}
		}
		if scrollOffset == nil {
			c.scrollContainerDatas = append(c.scrollContainerDatas, ScrollContainerDataInternal{
				layoutElement: openLayoutElement,
				scrollOrigin:  MakeVector2(-1, -1),
				elementId:     openLayoutElement.id,
				openThisFrame: true,
			})
			scrollOffset = &c.scrollContainerDatas[len(c.scrollContainerDatas)-1]
		}
		if c.externalScrollHandlingEnabled {
			// TODO: fix
			//scrollOffset.scrollPosition = queryScrollOffset(scrollOffset.elementId, c.queryScrollOffsetUserData)
		}
	}

	if !declaration.Border.IsEmpty() {
		c.attachElementConfig(c.storeBorderElementConfig(declaration.Border))
	}
}

// Ephemeral Memory - reset every frame
func (c *Context) initializeEphemeralMemory() {
	c.layoutElementChildrenBuffer = c.layoutElementChildrenBuffer[:0]
	c.layoutElements = c.layoutElements[:0]
	c.warnings = c.warnings[:0]

	c.layoutConfigs = c.layoutConfigs[:0]
	c.elementConfigs = c.elementConfigs[:0]
	c.textElementConfigs = c.textElementConfigs[:0]
	c.imageElementConfigs = c.imageElementConfigs[:0]
	c.floatingElementConfigs = c.floatingElementConfigs[:0]
	c.scrollElementConfigs = c.scrollElementConfigs[:0]
	c.customElementConfigs = c.customElementConfigs[:0]
	c.borderElementConfigs = c.borderElementConfigs[:0]
	c.sharedElementConfigs = c.sharedElementConfigs[:0]

	c.layoutElementIdStrings = c.layoutElementIdStrings[:0]
	c.wrappedTextLines = c.wrappedTextLines[:0]
	c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[:0]
	c.layoutElementTreeRoots = c.layoutElementTreeRoots[:0]
	c.layoutElementChildren = c.layoutElementChildren[:0]
	c.openLayoutElementStack = c.openLayoutElementStack[:0]
	c.textElementData = c.textElementData[:0]
	c.imageElementPointers = c.imageElementPointers[:0]
	c.renderCommands = c.renderCommands[:0]
	c.openClipElementStack = c.openClipElementStack[:0]
	c.reusableElementIndexBuffer = c.reusableElementIndexBuffer[:0]
	c.layoutElementClipElementIds = c.layoutElementClipElementIds[:0]
	c.dynamicStringData = c.dynamicStringData[:0]
}

// Persistent memory - initialized once and not reset
func (c *Context) initializePersistentMemory() {
	maxElementCount := c.maxElementCount
	maxMeasureTextCacheWordCount := c.maxMeasureTextCacheWordCount

	c.scrollContainerDatas = make([]ScrollContainerDataInternal, 0, 10)
	c.layoutElementsHashMapInternal = make([]LayoutElementHashMapItem, 0, maxElementCount)
	c.layoutElementsHashMap = map[uint32]*LayoutElementHashMapItem{}
	c.measureTextHashMapInternal = make([]MeasureTextCacheItem, 0, maxElementCount)
	c.measuredWordsFreeList = make([]int32, 0, maxMeasureTextCacheWordCount)
	c.measureTextHashMap = map[string]*MeasureTextCacheItem{}
	c.measuredWords = make([]MeasuredWord, 0, maxMeasureTextCacheWordCount)
	c.pointerOverIds = make([]ElementId, 0, maxElementCount)
	c.debugElementData = make([]DebugElementData, 0, maxElementCount)

	c.layoutElementChildrenBuffer = make([]int, 0, maxElementCount)
	c.layoutElements = make([]LayoutElement, 0, maxElementCount)
	c.warnings = make([]Warning, 0, 100)

	c.layoutConfigs = make([]LayoutConfig, 0, maxElementCount)
	c.elementConfigs = make([]AnyElementConfig, 0, maxElementCount)
	c.textElementConfigs = make([]TextElementConfig, 0, maxElementCount)
	c.imageElementConfigs = make([]ImageElementConfig, 0, maxElementCount)
	c.floatingElementConfigs = make([]FloatingElementConfig, 0, maxElementCount)
	c.scrollElementConfigs = make([]ScrollElementConfig, 0, maxElementCount)
	c.customElementConfigs = make([]CustomElementConfig, 0, maxElementCount)
	c.borderElementConfigs = make([]BorderElementConfig, 0, maxElementCount)
	c.sharedElementConfigs = make([]SharedElementConfig, 0, maxElementCount)

	c.layoutElementIdStrings = make([]string, 0, maxElementCount)
	c.wrappedTextLines = make([]WrappedTextLine, 0, maxElementCount)
	c.layoutElementTreeNodeArray1 = make([]LayoutElementTreeNode, 0, maxElementCount)
	c.layoutElementTreeRoots = make([]LayoutElementTreeRoot, 0, maxElementCount)
	c.layoutElementChildren = make([]int, 0, maxElementCount)
	c.openLayoutElementStack = make([]int, 0, maxElementCount)
	c.textElementData = make([]TextElementData, 0, maxElementCount)
	c.imageElementPointers = make([]int, 0, maxElementCount)
	c.renderCommands = make([]RenderCommand, 0, maxElementCount)
	c.openClipElementStack = make([]int, 0, maxElementCount)
	c.reusableElementIndexBuffer = make([]int32, 0, maxElementCount)
	c.layoutElementClipElementIds = make([]int, 0, maxElementCount)
	c.dynamicStringData = make([]byte, 0, maxElementCount)
}

var CLAY__EPSILON float32 = 0.01

func floatEqual(left float32, right float32) bool {
	subtracted := left - right
	return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON
}

func (c *Context) sizeContainersAlongAxis(axis int) {
	bfsBuffer := c.layoutElementChildrenBuffer
	resizableContainerBuffer := c.openLayoutElementStack[:]
	for _, root := range c.layoutElementTreeRoots {
		bfsBuffer = bfsBuffer[0:0]
		rootElement := &c.layoutElements[root.layoutElementIndex]
		bfsBuffer = append(bfsBuffer, root.layoutElementIndex)

		// Size floating containers to their parents
		if floatingElementConfig, ok := findElementConfigWithType[*FloatingElementConfig](rootElement); ok {
			parentItem := c.getHashMapItem(floatingElementConfig.ParentId)
			if parentItem != nil && !parentItem.IsEmpty() {
				parentLayoutElement := parentItem.layoutElement
				switch rootElement.layoutConfig.Sizing.Width.(type) {
				case SizingAxisGrow:
					rootElement.dimensions.X = parentLayoutElement.dimensions.X
				}
				switch rootElement.layoutConfig.Sizing.Height.(type) {
				case SizingAxisGrow:
					rootElement.dimensions.Y = parentLayoutElement.dimensions.Y
				}
			}
		}

		if mm, ok := rootElement.layoutConfig.Sizing.Width.(SizingAxisMinMax); ok {
			rootElement.dimensions.X = min(max(rootElement.dimensions.X, mm.GetMinMax().Min), mm.GetMinMax().Max)
		}
		if mm, ok := rootElement.layoutConfig.Sizing.Height.(SizingAxisMinMax); ok {
			rootElement.dimensions.Y = min(max(rootElement.dimensions.Y, mm.GetMinMax().Min), mm.GetMinMax().Max)
		}

		for bi := 0; bi < len(bfsBuffer); bi++ {
			parentIndex := bfsBuffer[bi]
			parent := &c.layoutElements[parentIndex]
			parentStyleConfig := parent.layoutConfig
			var growContainerCount int32
			parentSize := parent.dimensions.Axis(axis)
			var parentPadding float32
			var innerContentSize float32

			if axis == 0 {
				parentPadding = float32(parent.layoutConfig.Padding.Left + parent.layoutConfig.Padding.Right)
			} else {
				parentPadding = float32(parent.layoutConfig.Padding.Top + parent.layoutConfig.Padding.Bottom)
			}

			totalPaddingAndChildGaps := parentPadding
			sizingAlongAxis := parentStyleConfig.LayoutDirection.IsAlongAxis(axis)
			c.openLayoutElementStack = c.openLayoutElementStack[:0]
			resizableContainerBuffer = resizableContainerBuffer[:0]
			parentChildGap := parentStyleConfig.ChildGap

			for childOffset, childElementIndex := range parent.children {
				child := &c.layoutElements[childElementIndex]
				childSizing := child.layoutConfig.Sizing.GetAxis(axis)
				childSize := child.dimensions.Axis(axis)

				if !elementHasConfig[*TextElementConfig](child) && len(child.children) > 0 {
					bfsBuffer = append(bfsBuffer, childElementIndex)
				}

				if func() bool {
					switch childSizing.(type) {
					case SizingAxisFit:
						return true
					case SizingAxisGrow:
						return true
					}

					return false
				}() {
					if func() bool {
						if tc, ok := findElementConfigWithType[*TextElementConfig](child); !ok || tc.WrapMode == TEXT_WRAP_WORDS {
							if axis == 0 || !elementHasConfig[*ImageElementConfig](child) {
								return true
							}
						}
						return false
					}() {
						c.openLayoutElementStack = append(c.openLayoutElementStack, childElementIndex)
						resizableContainerBuffer = append(resizableContainerBuffer, childElementIndex)
					}
				}

				if sizingAlongAxis {
					switch childSizing.(type) {
					case SizingAxisPercent:
					case SizingAxisGrow:
						growContainerCount++
						innerContentSize += childSize
					default:
						innerContentSize += childSize
					}
					if childOffset > 0 {
						innerContentSize += float32(parentChildGap) // For children after index 0, the childAxisOffset is the gap from the previous child
						totalPaddingAndChildGaps += float32(parentChildGap)
					}
				} else {
					innerContentSize = max(childSize, innerContentSize)
				}
			}

			// Expand percentage containers to size
			for _, childElementIndex := range parent.children {
				child := &c.layoutElements[childElementIndex]
				childSizing := child.layoutConfig.Sizing.GetAxis(axis)
				childSize := child.dimensions.Axis(axis)

				switch p := childSizing.(type) {
				case SizingAxisPercent:
					childSize = (parentSize - totalPaddingAndChildGaps) * p.Percent
					child.dimensions = child.dimensions.SetAxis(axis, childSize)
					if sizingAlongAxis {
						innerContentSize += childSize
					}
					updateAspectRatioBox(child)
				}
			}

			if sizingAlongAxis {
				sizeToDistribute := parentSize - parentPadding - innerContentSize
				// The content is too large, compress the children as much as possible
				if sizeToDistribute < 0 {
					// If the parent can scroll in the axis direction in this direction, don't compress children, just leave them alone
					if scrollElementConfig, ok := findElementConfigWithType[*ScrollElementConfig](parent); ok {
						if (axis == 0 && scrollElementConfig.Horizontal) || (axis == 1 && scrollElementConfig.Vertical) {
							continue
						}
					}
					// Scrolling containers preferentially compress before others
					for sizeToDistribute < -CLAY__EPSILON && len(resizableContainerBuffer) > 0 {
						var largest float32
						var secondLargest float32
						widthToAdd := sizeToDistribute
						for _, childIndex := range resizableContainerBuffer {
							child := c.layoutElements[childIndex]
							childSize := child.dimensions.Axis(axis)
							if floatEqual(childSize, largest) {
								continue
							}
							if childSize > largest {
								secondLargest = largest
								largest = childSize
							}
							if childSize < largest {
								secondLargest = max(secondLargest, childSize)
								widthToAdd = secondLargest - largest
							}
						}

						widthToAdd = max(widthToAdd, sizeToDistribute/float32(len(resizableContainerBuffer)))

						for childIndex := range resizableContainerBuffer {
							child := &c.layoutElements[resizableContainerBuffer[childIndex]]
							childSize := child.dimensions.Axis(axis)
							minSize := child.minDimensions.Axis(axis)
							previousWidth := childSize
							if floatEqual(childSize, largest) {
								childSize += widthToAdd
								if childSize <= minSize {
									childSize = minSize
									resizableContainerBuffer, _ = slicesex.RemoveSwapback(resizableContainerBuffer, childIndex)
									childIndex--
								}
								child.dimensions = child.dimensions.SetAxis(axis, childSize)
								sizeToDistribute -= (childSize - previousWidth)
							}
						}
					}
					// The content is too small, allow SIZING_GROW containers to expand
				} else if sizeToDistribute > 0 && growContainerCount > 0 {
					for ci := 0; ci < len(resizableContainerBuffer); ci++ {
						child := &c.layoutElements[resizableContainerBuffer[ci]]
						childSizing := child.layoutConfig.Sizing.GetAxis(axis)
						switch childSizing.(type) {
						case SizingAxisGrow:
						default:
							resizableContainerBuffer, _ = slicesex.RemoveSwapback(resizableContainerBuffer, ci)
							ci--
						}
					}

					for sizeToDistribute > CLAY__EPSILON && len(resizableContainerBuffer) > 0 {
						smallest := float32(math.MaxFloat32)
						secondSmallest := float32(math.MaxFloat32)
						widthToAdd := sizeToDistribute
						for _, rcb := range resizableContainerBuffer {
							child := c.layoutElements[rcb]
							childSize := child.dimensions.Axis(axis)
							if floatEqual(childSize, smallest) {
								continue
							}
							if childSize < smallest {
								secondSmallest = smallest
								smallest = childSize
							}
							if childSize > smallest {
								secondSmallest = min(secondSmallest, childSize)
								widthToAdd = secondSmallest - smallest
							}
						}

						widthToAdd = min(widthToAdd, sizeToDistribute/float32(len(resizableContainerBuffer)))

						for childIndex := range resizableContainerBuffer {
							child := &c.layoutElements[resizableContainerBuffer[childIndex]]
							childSize := child.dimensions.Axis(axis)
							childSizing := child.layoutConfig.Sizing.GetAxis(axis)
							var maxSize float32
							switch mm := childSizing.(type) {
							case SizingAxisMinMax:
								maxSize = mm.GetMinMax().Max
							}
							previousWidth := childSize
							if floatEqual(childSize, smallest) {
								childSize += widthToAdd
								if childSize >= maxSize {
									childSize = maxSize
									resizableContainerBuffer, _ = slicesex.RemoveSwapback(resizableContainerBuffer, childIndex)
									childIndex--
								}
								child.dimensions = child.dimensions.SetAxis(axis, childSize)
								sizeToDistribute -= (childSize - previousWidth)
							}
						}
					}

				}
				// Sizing along the non layout axis ("off axis")
			} else {
				for _, rcb := range resizableContainerBuffer {
					child := &c.layoutElements[rcb]
					childSizing := child.layoutConfig.Sizing.GetAxis(axis)
					minSize := child.minDimensions.Axis(axis)
					childSize := child.dimensions.Axis(axis)

					if axis == 1 && elementHasConfig[*ImageElementConfig](child) {
						continue // Currently we don't support resizing aspect ratio images on the Y axis because it would break the ratio
					}

					// If we're laying out the children of a scroll panel, grow containers expand to the size of the inner content, not the outer container
					maxSize := parentSize - parentPadding
					if scrollElementConfig, ok := findElementConfigWithType[*ScrollElementConfig](parent); ok {
						if (axis == 0 && scrollElementConfig.Horizontal) || (axis == 1 && scrollElementConfig.Vertical) {
							maxSize = max(maxSize, innerContentSize)
						}
					}
					switch s := childSizing.(type) {
					case SizingAxisGrow:
						childSize = min(maxSize, s.MinMax.Max)
					}
					child.dimensions = child.dimensions.SetAxis(axis, max(minSize, min(childSize, maxSize)))
				}
			}
		}
	}
}

func (c *Context) addRenderCommand(renderCommand RenderCommand) {
	if len(c.renderCommands) < cap(c.renderCommands)-1 {
		c.renderCommands = append(c.renderCommands, renderCommand)
	} else {
		if !c.booleanWarnings.maxRenderCommandsExceeded {
			c.booleanWarnings.maxRenderCommandsExceeded = true
			c.errorHandler.ErrorHandlerFunction(ErrorData{
				ErrorType: ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
				ErrorText: "Clay ran out of capacity while attempting to create render commands. This is usually caused by a large amount of wrapping text elements while close to the max element capacity. Try using Clay_SetMaxElementCount() with a higher value.",
				UserData:  c.errorHandler.UserData,
			})
		}
	}
}

func (c *Context) elementIsOffscreen(boundingBox BoundingBox) bool {
	if c.disableCulling {
		return false
	}

	return (boundingBox.X() > c.layoutDimensions.X) ||
		(boundingBox.Y() > c.layoutDimensions.Y) ||
		(boundingBox.X()+boundingBox.Width() < 0) ||
		(boundingBox.Y()+boundingBox.Height() < 0)
}

func (c *Context) calculateFinalLayout() {
	treeNodeVisited := make([]bool, len(c.layoutElements))

	// Calculate sizing along the X axis
	c.sizeContainersAlongAxis(0)

	// Wrap text
	for i := range c.textElementData {
		textElementData := &c.textElementData[i]
		textElementData.wrappedLines = c.wrappedTextLines[len(c.wrappedTextLines):]
		containerElement := c.layoutElements[textElementData.elementIndex]
		textConfig, _ := findElementConfigWithType[*TextElementConfig](&containerElement)
		measureTextCacheItem := c.measureTextCached(textElementData.text, textConfig)
		var lineWidth float32
		lineHeight := textElementData.preferredDimensions.Y
		if textConfig.LineHeight > 0 {
			lineHeight = float32(textConfig.LineHeight)
		}
		var lineLengthChars int32
		var lineStartOffset int32

		if !measureTextCacheItem.containsNewlines && textElementData.preferredDimensions.X <= containerElement.dimensions.X {
			c.wrappedTextLines = append(c.wrappedTextLines, WrappedTextLine{})
			textElementData.wrappedLines = append(textElementData.wrappedLines, WrappedTextLine{containerElement.dimensions, textElementData.text})
			continue
		}
		spaceWidth := measureText(SPACECHAR, textConfig, c.measureTextUserData).X
		wordIndex := measureTextCacheItem.measuredWordsStartIndex
		for wordIndex != -1 {
			if len(c.wrappedTextLines) > cap(c.wrappedTextLines)-1 {
				break
			}
			measuredWord := c.measuredWords[wordIndex]
			// Only word on the line is too large, just render it anyway
			if lineLengthChars == 0 && lineWidth+measuredWord.width > containerElement.dimensions.X {
				c.wrappedTextLines = append(c.wrappedTextLines, WrappedTextLine{})
				textElementData.wrappedLines = append(textElementData.wrappedLines, WrappedTextLine{
					MakeDimensions(measuredWord.width, lineHeight),
					textElementData.text[measuredWord.startOffset : measuredWord.startOffset+measuredWord.length]})
				wordIndex = measuredWord.next
				lineStartOffset = measuredWord.startOffset + measuredWord.length
			} else if measuredWord.length == 0 || lineWidth+measuredWord.width > containerElement.dimensions.X {
				// measuredWord.length == 0 means a newline character
				// Wrapped text lines list has overflowed, just render out the line
				var addSpace float32
				if textElementData.text[lineStartOffset+lineLengthChars-1] == ' ' {
					addSpace = -spaceWidth
					lineLengthChars--
				}
				c.wrappedTextLines = append(c.wrappedTextLines, WrappedTextLine{})
				textElementData.wrappedLines = append(textElementData.wrappedLines, WrappedTextLine{
					MakeDimensions(lineWidth+addSpace, lineHeight),
					textElementData.text[lineStartOffset : lineStartOffset+lineLengthChars],
				})
				if lineLengthChars == 0 || measuredWord.length == 0 {
					wordIndex = measuredWord.next
				}
				lineWidth = 0
				lineLengthChars = 0
				lineStartOffset = measuredWord.startOffset
			} else {
				lineWidth += measuredWord.width
				lineLengthChars += measuredWord.length
				wordIndex = measuredWord.next
			}
		}
		if lineLengthChars > 0 {
			c.wrappedTextLines = append(c.wrappedTextLines, WrappedTextLine{})
			textElementData.wrappedLines = append(textElementData.wrappedLines, WrappedTextLine{
				MakeDimensions(lineWidth, lineHeight), textElementData.text[lineStartOffset : lineStartOffset+lineLengthChars]})
		}
		containerElement.dimensions.Y = lineHeight * float32(len(textElementData.wrappedLines))
	}
	// Scale vertical image heights according to aspect ratio
	for _, iep := range c.imageElementPointers {
		imageElement := &c.layoutElements[iep]
		if config, ok := findElementConfigWithType[*ImageElementConfig](imageElement); ok {
			imageElement.dimensions.Y = (config.SourceDimensions.Y / max(config.SourceDimensions.X, 1)) * imageElement.dimensions.X
		}
	}

	// Propagate effect of text wrapping, image aspect scaling etc. on height of parents
	c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[0:0]
	for _, root := range c.layoutElementTreeRoots {
		treeNodeVisited[len(c.layoutElementTreeNodeArray1)] = false
		c.layoutElementTreeNodeArray1 = append(c.layoutElementTreeNodeArray1, LayoutElementTreeNode{
			layoutElement: &c.layoutElements[root.layoutElementIndex],
		})
	}
	for len(c.layoutElementTreeNodeArray1) > 0 {
		currentElementTreeNode := c.layoutElementTreeNodeArray1[len(c.layoutElementTreeNodeArray1)-1]
		currentElement := currentElementTreeNode.layoutElement
		if !treeNodeVisited[len(c.layoutElementTreeNodeArray1)-1] {
			treeNodeVisited[len(c.layoutElementTreeNodeArray1)-1] = true
			// If the element has no children or is the container for a text element, don't bother inspecting it
			if elementHasConfig[*TextElementConfig](currentElement) || len(currentElement.children) == 0 {
				c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[:len(c.layoutElementTreeNodeArray1)-1]
				continue
			}
			// Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
			for _, child := range currentElement.children {
				treeNodeVisited[len(c.layoutElementTreeNodeArray1)] = false
				c.layoutElementTreeNodeArray1 = append(c.layoutElementTreeNodeArray1, LayoutElementTreeNode{
					layoutElement: &c.layoutElements[child],
				})
			}
			continue
		}
		c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[:len(c.layoutElementTreeNodeArray1)-1]

		// DFS node has been visited, this is on the way back up to the root
		layoutConfig := currentElement.layoutConfig
		if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
			// Resize any parent containers that have grown in height along their non layout axis
			for _, child := range currentElement.children {
				childElement := c.layoutElements[child]
				childHeightWithPadding := max(childElement.dimensions.Y+float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom), currentElement.dimensions.Y)
				switch mm := layoutConfig.Sizing.Height.(type) {
				case SizingAxisMinMax:
					currentElement.dimensions.Y = min(max(childHeightWithPadding, mm.GetMinMax().Min), mm.GetMinMax().Max)
				}
			}
		} else if layoutConfig.LayoutDirection == TOP_TO_BOTTOM {
			// Resizing along the layout axis
			contentHeight := float32(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)
			for _, child := range currentElement.children {
				childElement := c.layoutElements[child]
				contentHeight += childElement.dimensions.Y
			}
			contentHeight += float32(max(uint16(len(currentElement.children))-1, 0) * layoutConfig.ChildGap)
			switch mm := layoutConfig.Sizing.Height.(type) {
			case SizingAxisMinMax:
				currentElement.dimensions.Y = min(max(contentHeight, mm.GetMinMax().Min), mm.GetMinMax().Max)
			}
		}
	}

	// Calculate sizing along the Y axis
	c.sizeContainersAlongAxis(1)

	// Sort tree roots by z-index
	sortMax := len(c.layoutElementTreeRoots) - 1
	for sortMax > 0 { // todo dumb bubble sort
		for i := range sortMax {
			current := c.layoutElementTreeRoots[i]
			next := c.layoutElementTreeRoots[i+1]
			if next.zIndex < current.zIndex {
				c.layoutElementTreeRoots[i] = next
				c.layoutElementTreeRoots[i+1] = current
			}
		}
		sortMax--
	}

	// Calculate final positions and generate render commands
	c.renderCommands = c.renderCommands[0:0]
	for iroot := range c.layoutElementTreeRoots {
		root := &c.layoutElementTreeRoots[iroot]
		c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[0:0]
		rootElement := &c.layoutElements[root.layoutElementIndex]
		var rootPosition Vector2
		parentHashMapItem := c.getHashMapItem(root.parentId)
		// Position root floating containers
		if config, ok := findElementConfigWithType[*FloatingElementConfig](rootElement); ok && parentHashMapItem != nil {
			rootDimensions := rootElement.dimensions
			parentBoundingBox := parentHashMapItem.boundingBox
			// Set X position
			var targetAttachPosition Vector2
			switch config.AttachPoints.Parent {
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_LEFT_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X()
			case ATTACH_POINT_CENTER_TOP, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_CENTER_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X() + (parentBoundingBox.Width() / 2)
			case ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_RIGHT_CENTER, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X() + parentBoundingBox.Width()
			}
			switch config.AttachPoints.Element {
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_LEFT_BOTTOM:
				break
			case ATTACH_POINT_CENTER_TOP, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_CENTER_BOTTOM:
				targetAttachPosition.X -= (rootDimensions.X / 2)
			case ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_RIGHT_CENTER, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.X -= rootDimensions.X
			}
			switch config.AttachPoints.Parent { // I know I could merge the x and y switch statements, but this is easier to read
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_CENTER_TOP:
				targetAttachPosition.Y = parentBoundingBox.Y()
			case ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_RIGHT_CENTER:
				targetAttachPosition.Y = parentBoundingBox.Y() + (parentBoundingBox.Height() / 2)
			case ATTACH_POINT_LEFT_BOTTOM, ATTACH_POINT_CENTER_BOTTOM, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.Y = parentBoundingBox.Y() + parentBoundingBox.Height()
			}
			switch config.AttachPoints.Element {
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_CENTER_TOP:
				break
			case ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_RIGHT_CENTER:
				targetAttachPosition.Y -= (rootDimensions.Y / 2)
			case ATTACH_POINT_LEFT_BOTTOM, ATTACH_POINT_CENTER_BOTTOM, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.Y -= rootDimensions.Y
			}
			targetAttachPosition.X += config.Offset.X
			targetAttachPosition.Y += config.Offset.Y
			rootPosition = targetAttachPosition
		}

		if root.clipElementId != 0 {
			clipHashMapItem := c.getHashMapItem(root.clipElementId)
			if clipHashMapItem != nil {
				// Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
				if c.externalScrollHandlingEnabled {
					if scrollConfig, ok := findElementConfigWithType[*ScrollElementConfig](clipHashMapItem.layoutElement); ok {
						for _, mapping := range c.scrollContainerDatas {
							if mapping.layoutElement == clipHashMapItem.layoutElement {
								root.pointerOffset = mapping.scrollPosition
								if scrollConfig.Horizontal {
									rootPosition.X += mapping.scrollPosition.X
								}
								if scrollConfig.Vertical {
									rootPosition.Y += mapping.scrollPosition.Y
								}
								break
							}
						}
					}
				}
				c.addRenderCommand(RenderCommand{
					BoundingBox: clipHashMapItem.boundingBox,
					UserData:    0,
					Id:          hashNumber(rootElement.id, uint32(len(rootElement.children))+10).id, // TODO need a better strategy for managing derived ids
					ZIndex:      root.zIndex,
					RenderData:  ScissorsStartData{},
				})
			}
		}
		c.layoutElementTreeNodeArray1 = append(c.layoutElementTreeNodeArray1, LayoutElementTreeNode{
			layoutElement:   rootElement,
			position:        rootPosition,
			nextChildOffset: MakeVector2(rootElement.layoutConfig.Padding.Left, rootElement.layoutConfig.Padding.Top),
		})

		treeNodeVisited[0] = false
		for len(c.layoutElementTreeNodeArray1) > 0 {
			currentElementTreeNode := &c.layoutElementTreeNodeArray1[len(c.layoutElementTreeNodeArray1)-1]
			currentElement := currentElementTreeNode.layoutElement
			layoutConfig := currentElement.layoutConfig
			var scrollOffset Vector2

			// This will only be run a single time for each element in downwards DFS order
			if !treeNodeVisited[len(c.layoutElementTreeNodeArray1)-1] {
				treeNodeVisited[len(c.layoutElementTreeNodeArray1)-1] = true

				currentElementBoundingBox := MakeBoundingBox(currentElementTreeNode.position, currentElement.dimensions)
				if floatingElementConfig, ok := findElementConfigWithType[*FloatingElementConfig](currentElement); ok {
					expand := floatingElementConfig.Expand
					currentElementBoundingBox = currentElementBoundingBox.AddXYWH(-expand.X, -expand.Y, expand.X*2, expand.Y*2)
				}

				var scrollContainerData *ScrollContainerDataInternal
				// Apply scroll offsets to container
				if scrollConfig, ok := findElementConfigWithType[*ScrollElementConfig](currentElement); ok {
					// This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
					for i := range c.scrollContainerDatas {
						mapping := &c.scrollContainerDatas[i]
						if mapping.layoutElement == currentElement {
							scrollContainerData = mapping
							mapping.boundingBox = currentElementBoundingBox
							if scrollConfig.Horizontal {
								scrollOffset.X = mapping.scrollPosition.X
							}
							if scrollConfig.Vertical {
								scrollOffset.Y = mapping.scrollPosition.Y
							}
							if c.externalScrollHandlingEnabled {
								scrollOffset = vector2.Zero[float32]()
							}
							break
						}
					}
				}

				hashMapItem := c.getHashMapItem(currentElement.id)
				if hashMapItem != nil {
					hashMapItem.boundingBox = currentElementBoundingBox
					if hashMapItem.idAlias != 0 {
						hashMapItemAlias := c.getHashMapItem(hashMapItem.idAlias)
						if hashMapItemAlias != nil {
							hashMapItemAlias.boundingBox = currentElementBoundingBox
						}
					}
				}

				var sortedConfigIndexes [20]int
				for elementConfigIndex := range currentElement.elementConfigs {
					sortedConfigIndexes[elementConfigIndex] = elementConfigIndex
				}
				sortMax = len(currentElement.elementConfigs) - 1
				for sortMax > 0 { // todo dumb bubble sort

					sortOrder := func(ec AnyElementConfig) int {
						switch ec.(type) {
						case *ScrollElementConfig:
							return 1
						case *BorderElementConfig:
							return -1
						default:
							return 0
						}
					}

					for i := range sortMax {
						current := sortedConfigIndexes[i]
						next := sortedConfigIndexes[i+1]
						currentConfig := currentElement.elementConfigs[current]
						nextConfig := currentElement.elementConfigs[next]
						if sortOrder(nextConfig) > 0 || sortOrder(currentConfig) < 0 {
							sortedConfigIndexes[i] = next
							sortedConfigIndexes[i+1] = current
						}
					}
					sortMax--
				}

				emitRectangle := false
				// Create the render commands for this element
				sharedConfig, ok := findElementConfigWithType[*SharedElementConfig](currentElement)
				if ok && sharedConfig.backgroundColor.A > 0 {
					emitRectangle = true
				} else if !ok {
					emitRectangle = false
					sharedConfig = &default_SharedElementConfig
				}

				for elementConfigIndex := range currentElement.elementConfigs {
					elementConfig := currentElement.elementConfigs[sortedConfigIndexes[elementConfigIndex]]
					renderCommand := RenderCommand{
						BoundingBox: currentElementBoundingBox,
						UserData:    sharedConfig.userData,
						Id:          currentElement.id,
					}

					offscreen := c.elementIsOffscreen(currentElementBoundingBox)
					// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
					shouldRender := !offscreen
					switch cfg := elementConfig.(type) {
					case *FloatingElementConfig:
						shouldRender = false
					case *SharedElementConfig:
						shouldRender = false
					case *BorderElementConfig:
						shouldRender = false
					case *ScrollElementConfig:
						renderCommand.RenderData = ScissorsStartData{
							ScrollRenderData: ScrollRenderData{
								Horizontal: cfg.Horizontal,
								Vertical:   cfg.Vertical,
							},
						}
					case *ImageElementConfig:
						renderCommand.RenderData = ImageRenderData{
							BackgroundColor:  sharedConfig.backgroundColor,
							CornerRadius:     sharedConfig.cornerRadius,
							SourceDimensions: cfg.SourceDimensions,
							ImageData:        cfg.ImageData,
						}
						emitRectangle = false

					case *TextElementConfig:
						if !shouldRender {
							break
						}
						shouldRender = false
						naturalLineHeight := currentElement.textElementData.preferredDimensions.Y
						finalLineHeight := naturalLineHeight
						if cfg.LineHeight > 0 {
							finalLineHeight = float32(cfg.LineHeight)
						}
						lineHeightOffset := (finalLineHeight - naturalLineHeight) / 2
						yPosition := lineHeightOffset
						for lineIndex, wrappedLine := range currentElement.textElementData.wrappedLines {
							if len(wrappedLine.line) == 0 {
								yPosition += finalLineHeight
								continue
							}
							offset := (currentElementBoundingBox.Width() - wrappedLine.dimensions.X)
							if cfg.TextAlignment == TEXT_ALIGN_LEFT {
								offset = 0
							}
							if cfg.TextAlignment == TEXT_ALIGN_CENTER {
								offset /= 2
							}
							c.addRenderCommand(RenderCommand{
								BoundingBox: MakeBoundingBox(
									currentElementBoundingBox.Position.AddXY(offset, yPosition),
									wrappedLine.dimensions,
								),
								RenderData: TextRenderData{
									StringContents: wrappedLine.line,
									TextColor:      cfg.TextColor,
									FontId:         cfg.FontId,
									FontSize:       cfg.FontSize,
									LetterSpacing:  cfg.LetterSpacing,
									LineHeight:     cfg.LineHeight,
								},
								UserData: cfg.UserData,
								Id:       hashNumber(uint32(lineIndex), currentElement.id).id,
								ZIndex:   root.zIndex,
							})
							yPosition += finalLineHeight

							if !c.disableCulling && (currentElementBoundingBox.Y()+yPosition > c.layoutDimensions.Y) {
								break
							}
						}
					case *CustomElementConfig:
						{
							renderCommand.RenderData = CustomRenderData{
								BackgroundColor: sharedConfig.backgroundColor,
								CornerRadius:    sharedConfig.cornerRadius,
								CustomData:      cfg.CustomData,
							}
							emitRectangle = false
							break
						}
					default:
						break
					}
					if shouldRender {
						c.addRenderCommand(renderCommand)
					}
					if offscreen {
						// NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
						// Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
						// be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
					}
				}

				if emitRectangle {
					c.addRenderCommand(RenderCommand{
						BoundingBox: currentElementBoundingBox,
						RenderData: RectangleRenderData{
							BackgroundColor: sharedConfig.backgroundColor,
							CornerRadius:    sharedConfig.cornerRadius,
						},
						UserData: sharedConfig.userData,
						Id:       currentElement.id,
						ZIndex:   root.zIndex,
					})
				}

				// Setup initial on-axis alignment
				if !elementHasConfig[*TextElementConfig](currentElementTreeNode.layoutElement) {
					var contentSize Dimensions
					if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
						for child := range currentElement.children {
							childElement := c.layoutElements[child]
							contentSize.X += childElement.dimensions.X
							contentSize.Y = max(contentSize.Y, childElement.dimensions.Y)
						}
						contentSize.X += float32(max(len(currentElement.children)-1, 0) * int(layoutConfig.ChildGap))
						extraSpace := currentElement.dimensions.X - float32(layoutConfig.Padding.Left+layoutConfig.Padding.Right) - contentSize.X
						switch layoutConfig.ChildAlignment.X {
						case ALIGN_X_LEFT:
							extraSpace = 0
						case ALIGN_X_CENTER:
							extraSpace /= 2
						default:
							break
						}
						currentElementTreeNode.nextChildOffset.X += extraSpace
					} else {
						for child := range currentElement.children {
							childElement := c.layoutElements[child]
							contentSize.X = max(contentSize.X, childElement.dimensions.X)
							contentSize.Y += childElement.dimensions.Y
						}
						contentSize.Y += float32(max(len(currentElement.children)-1, 0) * int(layoutConfig.ChildGap))
						extraSpace := currentElement.dimensions.Y - float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom) - contentSize.Y
						switch layoutConfig.ChildAlignment.Y {
						case ALIGN_Y_TOP:
							extraSpace = 0
						case ALIGN_Y_CENTER:
							extraSpace /= 2
						default:
							break
						}
						currentElementTreeNode.nextChildOffset.Y += extraSpace
					}

					if scrollContainerData != nil {
						scrollContainerData.contentSize = MakeDimensions(
							contentSize.X+float32(layoutConfig.Padding.Left+layoutConfig.Padding.Right),
							contentSize.Y+float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom))
					}
				}
			} else {
				// DFS is returning upwards backwards
				closeScrollElement := false
				if scrollConfig, ok := findElementConfigWithType[*ScrollElementConfig](currentElement); ok {
					closeScrollElement = true
					for _, mapping := range c.scrollContainerDatas {
						if mapping.layoutElement == currentElement {
							if scrollConfig.Horizontal {
								scrollOffset.X = mapping.scrollPosition.X
							}
							if scrollConfig.Vertical {
								scrollOffset.Y = mapping.scrollPosition.Y
							}
							if c.externalScrollHandlingEnabled {
								scrollOffset = vector2.Zero[float32]()
							}
							break
						}
					}
				}

				if borderConfig, ok := findElementConfigWithType[*BorderElementConfig](currentElement); ok {
					currentElementData := c.getHashMapItem(currentElement.id)
					currentElementBoundingBox := currentElementData.boundingBox

					// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
					if !c.elementIsOffscreen(currentElementBoundingBox) {
						sharedConfig, ok := findElementConfigWithType[*SharedElementConfig](currentElement)
						if !ok {
							sharedConfig = &default_SharedElementConfig
						}
						renderCommand := RenderCommand{
							BoundingBox: currentElementBoundingBox,
							RenderData: BorderRenderData{
								Color:        borderConfig.Color,
								CornerRadius: sharedConfig.cornerRadius,
								Width:        borderConfig.Width,
							},
							UserData: sharedConfig.userData,
							Id:       hashNumber(currentElement.id, uint32(len(currentElement.children))).id,
						}
						c.addRenderCommand(renderCommand)
						if borderConfig.Width.BetweenChildren > 0 && borderConfig.Color.A > 0 {
							halfGap := layoutConfig.ChildGap / 2
							borderOffset := MakeVector2(layoutConfig.Padding.Left-halfGap, layoutConfig.Padding.Top-halfGap)
							if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
								for i, child := range currentElement.children {
									childElement := c.layoutElements[child]
									if i > 0 {
										c.addRenderCommand(RenderCommand{
											BoundingBox: MakeBoundingBox(
												currentElementBoundingBox.Position.Add(borderOffset).Add(scrollOffset),
												MakeDimensions(borderConfig.Width.BetweenChildren, currentElement.dimensions.Y),
											),
											RenderData: RectangleRenderData{
												BackgroundColor: borderConfig.Color,
											},
											UserData: sharedConfig.userData,
											Id:       hashNumber(currentElement.id, uint32(len(currentElement.children)+1+i)).id,
										})
									}
									borderOffset.X += (childElement.dimensions.X + float32(layoutConfig.ChildGap))
								}
							} else {
								for i, child := range currentElement.children {
									childElement := c.layoutElements[child]
									if i > 0 {
										c.addRenderCommand(RenderCommand{
											BoundingBox: MakeBoundingBox(
												currentElementBoundingBox.Position.Add(scrollOffset).AddY(borderOffset.Y),
												MakeDimensions(currentElement.dimensions.X, borderConfig.Width.BetweenChildren),
											),
											RenderData: RectangleRenderData{
												BackgroundColor: borderConfig.Color,
											},
											UserData: sharedConfig.userData,
											Id:       hashNumber(currentElement.id, uint32(len(currentElement.children)+1+i)).id,
										})
									}
									borderOffset.Y += (childElement.dimensions.Y + float32(layoutConfig.ChildGap))
								}
							}
						}
					}
				}
				// This exists because the scissor needs to end _after_ borders between elements
				if closeScrollElement {
					c.addRenderCommand(RenderCommand{
						Id:         hashNumber(currentElement.id, uint32(len(rootElement.children)+11)).id,
						RenderData: ScissorsEndData{},
					})
				}

				c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[:len(c.layoutElementTreeNodeArray1)-1]
				continue
			}

			// Add children to the DFS buffer
			if !elementHasConfig[*TextElementConfig](currentElement) {
				c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[0 : len(c.layoutElementTreeNodeArray1)+len(currentElement.children)]
				for i, child := range currentElement.children {
					childElement := &c.layoutElements[child]
					// Alignment along non layout axis
					if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
						currentElementTreeNode.nextChildOffset.Y = float32(currentElement.layoutConfig.Padding.Top)
						whiteSpaceAroundChild := currentElement.dimensions.Y - float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom) - childElement.dimensions.Y
						switch layoutConfig.ChildAlignment.Y {
						case ALIGN_Y_TOP:
							break
						case ALIGN_Y_CENTER:
							currentElementTreeNode.nextChildOffset.Y += whiteSpaceAroundChild / 2
						case ALIGN_Y_BOTTOM:
							currentElementTreeNode.nextChildOffset.Y += whiteSpaceAroundChild
						}
					} else {
						currentElementTreeNode.nextChildOffset.X = float32(currentElement.layoutConfig.Padding.Left)
						whiteSpaceAroundChild := currentElement.dimensions.X - float32(layoutConfig.Padding.Left+layoutConfig.Padding.Right) - childElement.dimensions.X
						switch layoutConfig.ChildAlignment.X {
						case ALIGN_X_LEFT:
							break
						case ALIGN_X_CENTER:
							currentElementTreeNode.nextChildOffset.X += whiteSpaceAroundChild / 2
						case ALIGN_X_RIGHT:
							currentElementTreeNode.nextChildOffset.X += whiteSpaceAroundChild
						}
					}

					childPosition := MakeVector2(
						currentElementTreeNode.position.X+currentElementTreeNode.nextChildOffset.X+scrollOffset.X,
						currentElementTreeNode.position.Y+currentElementTreeNode.nextChildOffset.Y+scrollOffset.Y,
					)

					// DFS buffer elements need to be added in reverse because stack traversal happens backwards
					newNodeIndex := len(c.layoutElementTreeNodeArray1) - 1 - i
					c.layoutElementTreeNodeArray1[newNodeIndex] = LayoutElementTreeNode{
						layoutElement:   childElement,
						position:        childPosition,
						nextChildOffset: MakeVector2(childElement.layoutConfig.Padding.Left, childElement.layoutConfig.Padding.Top),
					}
					treeNodeVisited[newNodeIndex] = false

					// Update parent offsets
					if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
						currentElementTreeNode.nextChildOffset.X += childElement.dimensions.X + float32(layoutConfig.ChildGap)
					} else {
						currentElementTreeNode.nextChildOffset.Y += childElement.dimensions.Y + float32(layoutConfig.ChildGap)
					}
				}
			}
		}

		if root.clipElementId != 0 {
			c.addRenderCommand(RenderCommand{
				Id:         hashNumber(rootElement.id, uint32(len(rootElement.children))+11).id,
				RenderData: ScissorsEndData{},
			})
		}
	}
}
