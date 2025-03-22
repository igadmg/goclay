package goclay

import (
	"math"

	"github.com/igadmg/goex/image/colorex"
	"github.com/igadmg/goex/slicesex"
	"github.com/igadmg/raylib-go/raymath/rect2"
	"github.com/igadmg/raylib-go/raymath/vector2"
)

var LAYOUT_DEFAULT LayoutConfig
var Color_DEFAULT colorex.RGBA
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
	backgroundColor colorex.RGBA
	cornerRadius    CornerRadius
	userData        any
}

/*

bool Clay__Array_RangeCheck(int32 index, int32 length);
bool Clay__Array_AddCapacityCheck(int32 length, int32 capacity);

CLAY__ARRAY_DEFINE_FUNCTIONS(Clay_RenderCommand, Clay_RenderCommandArray)
*/

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
	dimensions vector2.Float32
	line       string
}

type TextElementData struct {
	text                string
	preferredDimensions vector2.Float32
	elementIndex        int
	wrappedLines        []WrappedTextLine
}

type LayoutElement struct {
	//union {
	children        []int
	textElementData *TextElementData
	//}
	dimensions     vector2.Float32
	minDimensions  vector2.Float32
	layoutConfig   *LayoutConfig
	elementConfigs []AnyElementConfig
	id             uint32
}

type ScrollContainerDataInternal struct {
	layoutElement       *LayoutElement
	boundingBox         rect2.Float32
	contentSize         vector2.Float32
	scrollOrigin        vector2.Float32
	pointerOrigin       vector2.Float32
	scrollMomentum      vector2.Float32
	scrollPosition      vector2.Float32
	previousDelta       vector2.Float32
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
	boundingBox           rect2.Float32
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
	unwrappedDimensions     vector2.Float32
	measuredWordsStartIndex int32
	containsNewlines        bool
	// Hash map data
	id         uint32
	nextIndex  int32
	generation uint32
}

var default_MeasureTextCacheItem MeasureTextCacheItem

type LayoutElementTreeNode struct {
	layoutElement   *LayoutElement
	position        vector2.Float32
	nextChildOffset vector2.Float32
}

type LayoutElementTreeRoot struct {
	layoutElementIndex int
	parentId           uint32 // This can be zero in the case of the root layout tree
	clipElementId      uint32 // This can be zero if there is no clip element
	zIndex             int16
	pointerOffset      vector2.Float32 // Only used when scroll containers are managed externally
}

type Context struct {
	maxElementCount              int32
	maxMeasureTextCacheWordCount int32
	warningsEnabled              bool
	errorHandler                 ErrorHandler
	booleanWarnings              BooleanWarnings
	warnings                     []Warning

	pointerInfo                   PointerData
	layoutDimensions              vector2.Float32
	dynamicElementIndexBaseHash   ElementId
	dynamicElementIndex           uint32
	debugModeEnabled              bool
	disableCulling                bool
	externalScrollHandlingEnabled bool
	debugSelectedElementId        uint32
	generation                    uint32
	//arenaResetOffset              any ///uintptr_t
	measureTextUserData       any
	queryScrollOffsetUserData any
	//internalArena                 any ///Clay_Arena
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
	treeNodeVisited               []bool
	dynamicStringData             []byte
	debugElementData              []DebugElementData
}

var measureText func(text string, config *TextElementConfig, userData any) vector2.Float32
var Clay__QueryScrollOffset func(elementId uint32, userData any) vector2.Float32

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

func hashNumber(offset uint32, seed uint32) ElementId {
	hash := seed
	hash += (offset + 48)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)

	return ElementId{
		id:       hash + 1, // Reserve the hash result of zero as "null id"
		offset:   offset,
		baseId:   seed,
		stringId: STRING_DEFAULT,
	}
}

func hashString(key string, offset uint32, seed uint32) ElementId {
	hash := uint32(0)
	base := uint32(seed)

	for i := range key {
		base += uint32(key[i])
		base += (base << 10)
		base ^= (base >> 6)
	}
	hash = base
	hash += offset
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += (hash << 3)
	base += (base << 3)
	hash ^= (hash >> 11)
	base ^= (base >> 11)
	hash += (hash << 15)
	base += (base << 15)

	return ElementId{
		id:       hash + 1, // Reserve the hash result of zero as "null id"
		offset:   offset,
		baseId:   base + 1,
		stringId: key,
	}
}

func Clay__HashTextWithConfig(text string, config *TextElementConfig) uint32 {
	hash := uint32(0)

	if config.hashStringContents {
		maxLengthToHash := min(len(text), 256)
		for i := range maxLengthToHash {
			hash += uint32(text[i])
			hash += (hash << 10)
			hash ^= (hash >> 6)
		}
	} else {
		//pointerAsNumber = uint32(&text[0])
		//hash += pointerAsNumber
		//hash += (hash << 10)
		//hash ^= (hash >> 6)
	}

	hash += uint32(len(text))
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.fontId)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.fontSize)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.lineHeight)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.letterSpacing)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.wrapMode)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
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

	id := Clay__HashTextWithConfig(text, config)
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
			break
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
			elementHasScrollHorizontal = cfg.horizontal
			elementHasScrollVertical = cfg.vertical
			c.openClipElementStack = c.openClipElementStack[:len(c.openClipElementStack)-1]
			break
		case *FloatingElementConfig:
		}
	}

	// Attach children to the current open element
	openLayoutElement.children = c.layoutElementChildren[len(c.layoutElementChildren):len(c.layoutElementChildren)]
	if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
		openLayoutElement.dimensions.X = (float32)(layoutConfig.Padding.Left + layoutConfig.Padding.Right)
		for i := range openLayoutElement.children {
			childIndex := c.layoutElementChildrenBuffer[len(c.layoutElementChildrenBuffer)-(int)(len(openLayoutElement.children)+i)]
			child := c.layoutElements[childIndex]
			openLayoutElement.dimensions.X += child.dimensions.X
			openLayoutElement.dimensions.Y = max(
				openLayoutElement.dimensions.Y,
				child.dimensions.Y+(float32)(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom))

			// Minimum size of child elements doesn't matter to scroll containers as they can shrink and hide their contents
			if !elementHasScrollHorizontal {
				openLayoutElement.minDimensions.X += child.minDimensions.X
			}
			if !elementHasScrollVertical {
				openLayoutElement.minDimensions.Y = max(
					openLayoutElement.minDimensions.Y,
					child.minDimensions.Y+(float32)(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom))
			}
			c.layoutElementChildren = append(c.layoutElementChildren, childIndex)
		}

		childGap := (float32)(max(len(openLayoutElement.children)-1, 0) * int(layoutConfig.ChildGap))
		openLayoutElement.dimensions.X += childGap // TODO this is technically a bug with childgap and scroll containers
		openLayoutElement.minDimensions.X += childGap
	} else if layoutConfig.LayoutDirection == TOP_TO_BOTTOM {
		openLayoutElement.dimensions.Y = (float32)(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)
		for i := range openLayoutElement.children {
			childIndex := c.layoutElementChildrenBuffer[len(c.layoutElementChildrenBuffer)-len(openLayoutElement.children)+i]
			child := c.layoutElements[childIndex]
			openLayoutElement.dimensions.Y += child.dimensions.Y
			openLayoutElement.dimensions.X = max(
				openLayoutElement.dimensions.X,
				child.dimensions.X+(float32)(layoutConfig.Padding.Left+layoutConfig.Padding.Right))
			// Minimum size of child elements doesn't matter to scroll containers as they can shrink and hide their contents
			if !elementHasScrollVertical {
				openLayoutElement.minDimensions.Y += child.minDimensions.Y
			}
			if !elementHasScrollHorizontal {
				openLayoutElement.minDimensions.X = max(openLayoutElement.minDimensions.X, child.minDimensions.X+(float32)(layoutConfig.Padding.Left+layoutConfig.Padding.Right))
			}
			c.layoutElementChildren = append(c.layoutElementChildren, childIndex)
		}
		childGap := (float32)(max(len(openLayoutElement.children)-1, 0) * int(layoutConfig.ChildGap))
		openLayoutElement.dimensions.Y += childGap // TODO this is technically a bug with childgap and scroll containers
		openLayoutElement.minDimensions.Y += childGap
	}

	c.layoutElementChildrenBuffer = c.layoutElementChildrenBuffer[:len(c.layoutElementChildrenBuffer)-len(openLayoutElement.children)]

	// Clamp element min and max width to the values configured in the layout
	switch w := layoutConfig.Sizing.Width.(type) {
	case SizingAxisMinMax:
		mm := w.GetMinMax()
		if mm.Max <= 0 { // Set the max size if the user didn't specify, makes calculations easier
			mm.Max = math.MaxFloat32
		}
		openLayoutElement.dimensions.X = min(max(openLayoutElement.dimensions.X, mm.Min), mm.Max)
		openLayoutElement.minDimensions.X = min(max(openLayoutElement.minDimensions.X, mm.Min), mm.Max)
	default:
		openLayoutElement.dimensions.X = 0
	}

	// Clamp element min and max height to the values configured in the layout
	switch h := layoutConfig.Sizing.Height.(type) {
	case SizingAxisMinMax:
		mm := h.GetMinMax()
		if mm.Max <= 0 { // Set the max size if the user didn't specify, makes calculations easier
			mm.Max = math.MaxFloat32
		}
		openLayoutElement.dimensions.Y = min(max(openLayoutElement.dimensions.Y, mm.Min), mm.Max)
		openLayoutElement.minDimensions.Y = min(max(openLayoutElement.minDimensions.Y, mm.Min), mm.Max)
	default:
		openLayoutElement.dimensions.Y = 0
	}

	updateAspectRatioBox(openLayoutElement)

	elementIsFloating := elementHasConfig[*FloatingElementConfig](openLayoutElement)

	// Close the currently open element
	var closingElementIndex int
	c.openLayoutElementStack, closingElementIndex = slicesex.RemoveSwapback(c.openLayoutElementStack, len(c.openLayoutElementStack)-1)
	openLayoutElement = c.getOpenLayoutElement()

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

	layoutElement := LayoutElement{}
	c.layoutElements = append(c.layoutElements, layoutElement)
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
	if textConfig.lineHeight > 0 {
		textDimensions.Y = (float32)(textConfig.lineHeight)
	}
	textElement.dimensions = textDimensions
	textElement.minDimensions = vector2.NewFloat32(textMeasured.unwrappedDimensions.Y, textDimensions.Y) // TODO not sure this is the best way to decide min width for text
	c.textElementData = append(c.textElementData, TextElementData{
		text:                text,
		preferredDimensions: textMeasured.unwrappedDimensions,
		elementIndex:        len(c.layoutElements) - 1,
	})
	textElement.textElementData = &c.textElementData[len(c.textElementData)-1]
	//textElement.elementConfigs = CLAY__INIT(Clay__ElementConfigArraySlice) {
	//        .length = 1,
	//        .internalArray = Clay__ElementConfigArray_Add(&c.elementConfigs, CLAY__INIT(Clay_ElementConfig) { .Type = CLAY__ELEMENT_CONFIG_TYPE_TEXT, .config = { .textElementConfig = textConfig }})
	//};
	textElement.layoutConfig = &default_LayoutConfig

	// TODO: fix
	//c.layoutElementChildren = append(c.layoutElementChildren, closingElementIndex)
	//parentElement.children.length++
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

	if declaration.Floating.attachTo != ATTACH_TO_NONE {
		floatingConfig := declaration.Floating
		// This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
		hierarchicalParent := c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-2]]
		if true /*hierarchicalParent.id != 0*/ {
			clipElementId := 0
			if declaration.Floating.attachTo == ATTACH_TO_PARENT {
				// Attach to the element's direct hierarchical parent
				floatingConfig.parentId = hierarchicalParent.id
				if len(c.openClipElementStack) > 0 {
					clipElementId = c.openClipElementStack[len(c.openClipElementStack)-1]
				}
			} else if declaration.Floating.attachTo == ATTACH_TO_ELEMENT_WITH_ID {
				parentItem := c.getHashMapItem(floatingConfig.parentId)
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
			} else if declaration.Floating.attachTo == ATTACH_TO_ROOT {
				floatingConfig.parentId = hashString("Clay__RootContainer", 0, 0).id
			}
			if openLayoutElementId.id == 0 {
				openLayoutElementId = hashString("Clay__FloatingContainer", uint32(len(c.layoutElementTreeRoots)), 0)
			}
			currentElementIndex := c.openLayoutElementStack[len(c.openLayoutElementStack)-1]
			c.layoutElementClipElementIds = slicesex.Set(c.layoutElementClipElementIds, currentElementIndex, clipElementId)
			c.openClipElementStack = append(c.openClipElementStack, clipElementId)
			c.layoutElementTreeRoots = append(c.layoutElementTreeRoots, LayoutElementTreeRoot{
				layoutElementIndex: c.openLayoutElementStack[len(c.openLayoutElementStack)-1],
				parentId:           floatingConfig.parentId,
				clipElementId:      uint32(clipElementId),
				zIndex:             floatingConfig.zIndex,
			})
			c.attachElementConfig(c.storeFloatingElementConfig(floatingConfig))
		}
	}
	if declaration.Custom.customData != nil {
		c.attachElementConfig(c.storeCustomElementConfig(declaration.Custom))
	}

	if openLayoutElementId.id != 0 {
		c.attachId(openLayoutElementId)
	} else if openLayoutElement.id == 0 {
		openLayoutElementId = c.generateIdForAnonymousElement(openLayoutElement)
	}

	if declaration.Scroll.horizontal || declaration.Scroll.vertical {
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
				scrollOrigin:  vector2.NewFloat32(-1, -1),
				elementId:     openLayoutElement.id,
				openThisFrame: true,
			})
			scrollOffset = &c.scrollContainerDatas[len(c.scrollContainerDatas)-1]
		}
		if c.externalScrollHandlingEnabled {
			// TODO: fix
			//scrollOffset.scrollPosition = Clay__QueryScrollOffset(scrollOffset.elementId, c.queryScrollOffsetUserData)
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
	for i := range c.treeNodeVisited {
		c.treeNodeVisited[i] = false
	}
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
	c.treeNodeVisited = make([]bool, maxElementCount)
	c.openClipElementStack = make([]int, 0, maxElementCount)
	c.reusableElementIndexBuffer = make([]int32, 0, maxElementCount)
	c.layoutElementClipElementIds = make([]int, 0, maxElementCount)
	c.dynamicStringData = make([]byte, 0, maxElementCount)
}

var CLAY__EPSILON float32 = 0.01

func Clay__FloatEqual(left float32, right float32) bool {
	subtracted := left - right
	return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON
}

func (c *Context) Clay__SizeContainersAlongAxis(xAxis bool) {

	getSize := func(d vector2.Float32) float32 {
		if xAxis {
			return d.X
		} else {
			return d.Y
		}
	}
	getSizePtr := func(d *vector2.Float32) *float32 {
		if xAxis {
			return &d.X
		} else {
			return &d.Y
		}
	}

	bfsBuffer := c.layoutElementChildrenBuffer
	resizableContainerBuffer := c.openLayoutElementStack
	for rootIndex := range c.layoutElementTreeRoots {
		bfsBuffer = bfsBuffer[0:0]
		root := c.layoutElementTreeRoots[rootIndex]
		rootElement := c.layoutElements[root.layoutElementIndex]
		bfsBuffer = append(bfsBuffer, root.layoutElementIndex)

		// Size floating containers to their parents
		if floatingElementConfig, ok := findElementConfigWithType[*FloatingElementConfig](&rootElement); ok {
			parentItem := c.getHashMapItem(floatingElementConfig.parentId)
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

		for _, parentIndex := range bfsBuffer {
			parent := c.layoutElements[parentIndex]
			parentStyleConfig := parent.layoutConfig
			var growContainerCount int32
			parentSize := getSize(parent.dimensions)
			var parentPadding float32
			var innerContentSize float32

			if xAxis {
				parentPadding = float32(parent.layoutConfig.Padding.Left + parent.layoutConfig.Padding.Right)
			} else {
				parentPadding = float32(parent.layoutConfig.Padding.Top + parent.layoutConfig.Padding.Bottom)
			}

			totalPaddingAndChildGaps := parentPadding
			sizingAlongAxis := (xAxis && parentStyleConfig.LayoutDirection == LEFT_TO_RIGHT) || (!xAxis && parentStyleConfig.LayoutDirection == TOP_TO_BOTTOM)
			resizableContainerBuffer = resizableContainerBuffer[:0]
			parentChildGap := parentStyleConfig.ChildGap

			for childOffset, childElementIndex := range parent.children {
				childElement := c.layoutElements[childElementIndex]
				childSizing := childElement.layoutConfig.Sizing.GetAxis(xAxis)
				childSize := getSize(childElement.dimensions)

				if !elementHasConfig[*TextElementConfig](&childElement) && len(childElement.children) > 0 {
					bfsBuffer = append(bfsBuffer, childElementIndex)
				}

				switch childSizing.(type) {
				case SizingAxisFit:
					c.openLayoutElementStack = append(c.openLayoutElementStack, childElementIndex)
					resizableContainerBuffer = append(resizableContainerBuffer, childElementIndex)
				case SizingAxisGrow:
					c.openLayoutElementStack = append(c.openLayoutElementStack, childElementIndex)
					resizableContainerBuffer = append(resizableContainerBuffer, childElementIndex)
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
				childElement := c.layoutElements[childElementIndex]
				childSizing := childElement.layoutConfig.Sizing.GetAxis(xAxis)
				childSize := getSizePtr(&childElement.dimensions)

				switch p := childSizing.(type) {
				case SizingAxisPercent:
					*childSize = (parentSize - totalPaddingAndChildGaps) * p.Percent
					if sizingAlongAxis {
						innerContentSize += *childSize
					}
					// TODO: fix Clay__UpdateAspectRatioBox(childElement)
				}
			}

			if sizingAlongAxis {
				sizeToDistribute := parentSize - parentPadding - innerContentSize
				// The content is too large, compress the children as much as possible
				if sizeToDistribute < 0 {
					// If the parent can scroll in the axis direction in this direction, don't compress children, just leave them alone
					if scrollElementConfig, ok := findElementConfigWithType[*ScrollElementConfig](&parent); ok {
						if (xAxis && scrollElementConfig.horizontal) || (!xAxis && scrollElementConfig.vertical) {
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
							childSize := getSize(child.dimensions)
							if Clay__FloatEqual(childSize, largest) {
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
							childSize := getSizePtr(&child.dimensions)
							minSize := getSize(child.minDimensions)
							previousWidth := *childSize
							if Clay__FloatEqual(*childSize, largest) {
								*childSize += widthToAdd
								if *childSize <= minSize {
									*childSize = minSize
									resizableContainerBuffer, _ = slicesex.RemoveSwapback(resizableContainerBuffer, childIndex)
									childIndex--
								}
								sizeToDistribute -= (*childSize - previousWidth)
							}
						}
					}
					// The content is too small, allow SIZING_GROW containers to expand
				} else if sizeToDistribute > 0 && growContainerCount > 0 {
					for childIndex := range resizableContainerBuffer {
						child := &c.layoutElements[resizableContainerBuffer[childIndex]]
						childSizing := child.layoutConfig.Sizing.GetAxis(xAxis)
						switch childSizing.(type) {
						case SizingAxisGrow:
						default:
							resizableContainerBuffer, _ = slicesex.RemoveSwapback(resizableContainerBuffer, childIndex)
							childIndex--
						}
					}

					for sizeToDistribute > CLAY__EPSILON && len(resizableContainerBuffer) > 0 {
						smallest := float32(math.MaxFloat32)
						secondSmallest := float32(math.MaxFloat32)
						widthToAdd := sizeToDistribute
						for _, rcb := range resizableContainerBuffer {
							child := c.layoutElements[rcb]
							childSize := getSize(child.dimensions)
							if Clay__FloatEqual(childSize, smallest) {
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
							child := c.layoutElements[resizableContainerBuffer[childIndex]]
							childSize := getSizePtr(&child.dimensions)
							childSizing := child.layoutConfig.Sizing.GetAxis(xAxis)
							var maxSize float32
							switch mm := childSizing.(type) {
							case SizingMinMax:
								maxSize = mm.Max
							}
							previousWidth := *childSize
							if Clay__FloatEqual(*childSize, smallest) {
								*childSize += widthToAdd
								if *childSize >= maxSize {
									*childSize = maxSize
									resizableContainerBuffer, _ = slicesex.RemoveSwapback(resizableContainerBuffer, childIndex)
									childIndex--
								}
								sizeToDistribute -= (*childSize - previousWidth)
							}
						}
					}

				}
				// Sizing along the non layout axis ("off axis")
			} else {
				for _, rcb := range resizableContainerBuffer {
					childElement := c.layoutElements[rcb]
					childSizing := childElement.layoutConfig.Sizing.GetAxis(xAxis)
					childSize := getSizePtr(&childElement.dimensions)

					if !xAxis && elementHasConfig[*ImageElementConfig](&childElement) {
						continue // Currently we don't support resizing aspect ratio images on the Y axis because it would break the ratio
					}

					// If we're laying out the children of a scroll panel, grow containers expand to the height of the inner content, not the outer container
					maxSize := parentSize - parentPadding
					if scrollElementConfig, ok := findElementConfigWithType[*ScrollElementConfig](&parent); ok {
						if (xAxis && scrollElementConfig.horizontal) || (!xAxis && scrollElementConfig.vertical) {
							maxSize = max(maxSize, innerContentSize)
						}
					}
					switch s := childSizing.(type) {
					case SizingAxisFit:
						*childSize = max(s.MinMax.Min, min(*childSize, maxSize))
					case SizingAxisGrow:
						*childSize = min(maxSize, s.MinMax.Max)
					}
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

func (c *Context) Clay__ElementIsOffscreen(boundingBox rect2.Float32) bool {
	if c.disableCulling {
		return false
	}

	return (boundingBox.X() > c.layoutDimensions.X) ||
		(boundingBox.Y() > c.layoutDimensions.Y) ||
		(boundingBox.X()+boundingBox.Width() < 0) ||
		(boundingBox.Y()+boundingBox.Height() < 0)
}

func (c *Context) Clay__CalculateFinalLayout() {
	// Calculate sizing along the X axis
	c.Clay__SizeContainersAlongAxis(true)

	/*
		    // Wrap text
		    for i := range c.textElementData {
		        textElementData = &c.textElementData[i];
		        textElementData.wrappedLines = WrappedTextLineArraySlice
				 { .length = 0, .internalArray = &c.wrappedTextLines.internalArray[c.wrappedTextLines.length] };
		        Clay_LayoutElement *containerElement = Clay_LayoutElementArray_Get(&c.layoutElements, (int)textElementData.elementIndex);
		        Clay_TextElementConfig *textConfig = findElementConfigWithType(containerElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig;
		        Clay__MeasureTextCacheItem *measureTextCacheItem = Clay__MeasureTextCached(&textElementData.text, textConfig);
		        float lineWidth = 0;
		        float lineHeight = textConfig.lineHeight > 0 ? (float32)textConfig.lineHeight : textElementData.preferredDimensions.Y;
		        int32 lineLengthChars = 0;
		        int32 lineStartOffset = 0;
		        if (!measureTextCacheItem.containsNewlines && textElementData.preferredDimensions.X <= containerElement.dimensions.X) {
		            Clay__WrappedTextLineArray_Add(&c.wrappedTextLines, WrappedTextLine
						 { containerElement.dimensions,  textElementData.text });
		            textElementData.wrappedLines.length++;
		            continue;
		        }
		        float spaceWidth = Clay__MeasureText(tringSlice
					 { .length = 1, .chars = SPACECHAR.chars, .baseChars = SPACECHAR.chars }, textConfig, c.measureTextUserData).X;
		        int32 wordIndex = measureTextCacheItem.measuredWordsStartIndex;
		        while (wordIndex != -1) {
		            if (c.wrappedTextLines.length > c.wrappedTextLines.capacity - 1) {
		                break;
		            }
		            Clay__MeasuredWord *measuredWord = Clay__MeasuredWordArray_Get(&c.measuredWords, wordIndex);
		            // Only word on the line is too large, just render it anyway
		            if (lineLengthChars == 0 && lineWidth + measuredWord.X > containerElement.dimensions.X) {
		                Clay__WrappedTextLineArray_Add(&c.wrappedTextLines, WrappedTextLine
							 { { measuredWord.X, lineHeight }, { .length = measuredWord.length, .chars = &textElementData.text.chars[measuredWord.startOffset] } });
		                textElementData.wrappedLines.length++;
		                wordIndex = measuredWord.next;
		                lineStartOffset = measuredWord.startOffset + measuredWord.length;
		            }
		            // measuredWord.length == 0 means a newline character
		            else if (measuredWord.length == 0 || lineWidth + measuredWord.X > containerElement.dimensions.X) {
		                // Wrapped text lines list has overflowed, just render out the line
		                bool finalCharIsSpace = textElementData.text.chars[lineStartOffset + lineLengthChars - 1] == ' ';
		                Clay__WrappedTextLineArray_Add(&c.wrappedTextLines, WrappedTextLine
							 { { lineWidth + (finalCharIsSpace ? -spaceWidth : 0), lineHeight }, { .length = lineLengthChars + (finalCharIsSpace ? -1 : 0), .chars = &textElementData.text.chars[lineStartOffset] } });
		                textElementData.wrappedLines.length++;
		                if (lineLengthChars == 0 || measuredWord.length == 0) {
		                    wordIndex = measuredWord.next;
		                }
		                lineWidth = 0;
		                lineLengthChars = 0;
		                lineStartOffset = measuredWord.startOffset;
		            } else {
		                lineWidth += measuredWord.X;
		                lineLengthChars += measuredWord.length;
		                wordIndex = measuredWord.next;
		            }
		        }
		        if (lineLengthChars > 0) {
		            Clay__WrappedTextLineArray_Add(&c.wrappedTextLines, WrappedTextLine
						 { { lineWidth, lineHeight }, {.length = lineLengthChars, .chars = &textElementData.text.chars[lineStartOffset] } });
		            textElementData.wrappedLines.length++;
		        }
		        containerElement.dimensions.Y = lineHeight * (float32)textElementData.wrappedLines.length;
		    }
	*/
	// Scale vertical image heights according to aspect ratio
	for _, iep := range c.imageElementPointers {
		imageElement := c.layoutElements[iep]
		if config, ok := findElementConfigWithType[*ImageElementConfig](&imageElement); ok {
			imageElement.dimensions.Y = (config.SourceDimensions.Y / max(config.SourceDimensions.X, 1)) * imageElement.dimensions.X
		}
	}

	// Propagate effect of text wrapping, image aspect scaling etc. on height of parents
	dfsBuffer := c.layoutElementTreeNodeArray1[0:0]
	for _, root := range c.layoutElementTreeRoots {
		c.treeNodeVisited[len(dfsBuffer)] = false
		dfsBuffer = append(dfsBuffer, LayoutElementTreeNode{
			layoutElement: &c.layoutElements[root.layoutElementIndex],
		})
	}
	for len(dfsBuffer) > 0 {
		currentElementTreeNode := dfsBuffer[len(dfsBuffer)-1]
		currentElement := currentElementTreeNode.layoutElement
		if !c.treeNodeVisited[len(dfsBuffer)-1] {
			c.treeNodeVisited[len(dfsBuffer)-1] = true
			// If the element has no children or is the container for a text element, don't bother inspecting it
			if elementHasConfig[*TextElementConfig](currentElement) || len(currentElement.children) == 0 {
				dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]
				continue
			}
			// Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
			for _, child := range currentElement.children {
				c.treeNodeVisited[len(dfsBuffer)] = false
				dfsBuffer = append(dfsBuffer, LayoutElementTreeNode{
					layoutElement: &c.layoutElements[child],
				})
			}
			continue
		}
		dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]

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
			contentHeight := (float32)(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)
			for _, child := range currentElement.children {
				childElement := c.layoutElements[child]
				contentHeight += childElement.dimensions.Y
			}
			contentHeight += (float32)(max(uint16(len(currentElement.children))-1, 0) * layoutConfig.ChildGap)
			switch mm := layoutConfig.Sizing.Height.(type) {
			case SizingAxisMinMax:
				currentElement.dimensions.Y = min(max(contentHeight, mm.GetMinMax().Min), mm.GetMinMax().Max)
			}
		}
	}

	// Calculate sizing along the Y axis
	c.Clay__SizeContainersAlongAxis(false)

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
		dfsBuffer = dfsBuffer[0:0]
		rootElement := &c.layoutElements[root.layoutElementIndex]
		var rootPosition vector2.Float32
		parentHashMapItem := c.getHashMapItem(root.parentId)
		// Position root floating containers
		if config, ok := findElementConfigWithType[*FloatingElementConfig](rootElement); ok && parentHashMapItem != nil {
			rootDimensions := rootElement.dimensions
			parentBoundingBox := parentHashMapItem.boundingBox
			// Set X position
			var targetAttachPosition vector2.Float32
			switch config.attachPoints.parent {
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_LEFT_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X()
			case ATTACH_POINT_CENTER_TOP, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_CENTER_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X() + (parentBoundingBox.Width() / 2)
			case ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_RIGHT_CENTER, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X() + parentBoundingBox.Width()
			}
			switch config.attachPoints.element {
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_LEFT_BOTTOM:
				break
			case ATTACH_POINT_CENTER_TOP, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_CENTER_BOTTOM:
				targetAttachPosition.X -= (rootDimensions.X / 2)
			case ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_RIGHT_CENTER, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.X -= rootDimensions.X
			}
			switch config.attachPoints.parent { // I know I could merge the x and y switch statements, but this is easier to read
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_CENTER_TOP:
				targetAttachPosition.Y = parentBoundingBox.Y()
			case ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_RIGHT_CENTER:
				targetAttachPosition.Y = parentBoundingBox.Y() + (parentBoundingBox.Height() / 2)
			case ATTACH_POINT_LEFT_BOTTOM, ATTACH_POINT_CENTER_BOTTOM, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.Y = parentBoundingBox.Y() + parentBoundingBox.Height()
			}
			switch config.attachPoints.element {
			case ATTACH_POINT_LEFT_TOP, ATTACH_POINT_RIGHT_TOP, ATTACH_POINT_CENTER_TOP:
				break
			case ATTACH_POINT_LEFT_CENTER, ATTACH_POINT_CENTER_CENTER, ATTACH_POINT_RIGHT_CENTER:
				targetAttachPosition.Y -= (rootDimensions.Y / 2)
			case ATTACH_POINT_LEFT_BOTTOM, ATTACH_POINT_CENTER_BOTTOM, ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.Y -= rootDimensions.Y
			}
			targetAttachPosition.X += config.offset.X
			targetAttachPosition.Y += config.offset.Y
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
								if scrollConfig.horizontal {
									rootPosition.X += mapping.scrollPosition.X
								}
								if scrollConfig.vertical {
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
		dfsBuffer = append(dfsBuffer, LayoutElementTreeNode{
			layoutElement:   rootElement,
			position:        rootPosition,
			nextChildOffset: vector2.NewFloat32((float32)(rootElement.layoutConfig.Padding.Left), (float32)(rootElement.layoutConfig.Padding.Top)),
		})

		c.treeNodeVisited[0] = false
		for len(dfsBuffer) > 0 {
			currentElementTreeNode := &dfsBuffer[len(dfsBuffer)-1]
			currentElement := currentElementTreeNode.layoutElement
			layoutConfig := currentElement.layoutConfig
			var scrollOffset vector2.Float32

			// This will only be run a single time for each element in downwards DFS order
			if !c.treeNodeVisited[len(dfsBuffer)-1] {
				c.treeNodeVisited[len(dfsBuffer)-1] = true

				currentElementBoundingBox := rect2.NewFloat32(currentElementTreeNode.position, currentElement.dimensions)
				if floatingElementConfig, ok := findElementConfigWithType[*FloatingElementConfig](currentElement); ok {
					expand := floatingElementConfig.expand
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
							if scrollConfig.horizontal {
								scrollOffset.X = mapping.scrollPosition.X
							}
							if scrollConfig.vertical {
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
					for i := range sortMax {
						current := sortedConfigIndexes[i]
						next := sortedConfigIndexes[i+1]
						_ = current
						_ = next
						/* TODO: fix that
						currentType := currentElement.elementConfigs[current].type;
						nextType := currentElement.elementConfigs[next].type;
						if (nextType == CLAY__ELEMENT_CONFIG_TYPE_SCROLL || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
							sortedConfigIndexes[i] = next;
							sortedConfigIndexes[i + 1] = current;
						}
						*/
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

					offscreen := c.Clay__ElementIsOffscreen(currentElementBoundingBox)
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
							horizontal: cfg.horizontal,
							vertical:   cfg.vertical,
						}
					case *ImageElementConfig:
						renderCommand.RenderData = ImageRenderData{
							backgroundColor:  sharedConfig.backgroundColor,
							cornerRadius:     sharedConfig.cornerRadius,
							sourceDimensions: cfg.SourceDimensions,
							imageData:        cfg.ImageData,
						}
						emitRectangle = false

					case *TextElementConfig:
						if !shouldRender {
							break
						}
						shouldRender = false
						naturalLineHeight := currentElement.textElementData.preferredDimensions.Y
						finalLineHeight := naturalLineHeight
						if cfg.lineHeight > 0 {
							finalLineHeight = (float32)(cfg.lineHeight)
						}
						lineHeightOffset := (finalLineHeight - naturalLineHeight) / 2
						yPosition := lineHeightOffset
						for lineIndex, wrappedLine := range currentElement.textElementData.wrappedLines {
							if len(wrappedLine.line) == 0 {
								yPosition += finalLineHeight
								continue
							}
							offset := (currentElementBoundingBox.Width() - wrappedLine.dimensions.X)
							if cfg.textAlignment == TEXT_ALIGN_LEFT {
								offset = 0
							}
							if cfg.textAlignment == TEXT_ALIGN_CENTER {
								offset /= 2
							}
							c.addRenderCommand(RenderCommand{
								BoundingBox: rect2.NewFloat32(
									currentElementBoundingBox.Position.AddXY(offset, yPosition),
									wrappedLine.dimensions,
								),
								RenderData: TextRenderData{
									stringContents: wrappedLine.line,
									textColor:      cfg.textColor,
									fontId:         cfg.fontId,
									fontSize:       cfg.fontSize,
									letterSpacing:  cfg.letterSpacing,
									lineHeight:     cfg.lineHeight,
								},
								UserData: cfg.userData,
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
								backgroundColor: sharedConfig.backgroundColor,
								cornerRadius:    sharedConfig.cornerRadius,
								customData:      cfg.customData,
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
							backgroundColor: sharedConfig.backgroundColor,
							cornerRadius:    sharedConfig.cornerRadius,
						},
						UserData: sharedConfig.userData,
						Id:       currentElement.id,
						ZIndex:   root.zIndex,
					})
				}

				// Setup initial on-axis alignment
				if !elementHasConfig[*TextElementConfig](currentElementTreeNode.layoutElement) {
					var contentSize vector2.Float32
					if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
						for child := range currentElement.children {
							childElement := c.layoutElements[child]
							contentSize.X += childElement.dimensions.X
							contentSize.Y = max(contentSize.Y, childElement.dimensions.Y)
						}
						contentSize.X += float32(max(len(currentElement.children)-1, 0) * int(layoutConfig.ChildGap))
						extraSpace := currentElement.dimensions.X - (float32)(layoutConfig.Padding.Left+layoutConfig.Padding.Right) - contentSize.X
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
						contentSize.Y += (float32)(max(len(currentElement.children)-1, 0) * int(layoutConfig.ChildGap))
						extraSpace := currentElement.dimensions.Y - (float32)(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom) - contentSize.Y
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
						scrollContainerData.contentSize = vector2.NewFloat32(
							contentSize.X+(float32)(layoutConfig.Padding.Left+layoutConfig.Padding.Right),
							contentSize.Y+(float32)(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom))
					}
				}
			} else {
				// DFS is returning upwards backwards
				closeScrollElement := false
				if scrollConfig, ok := findElementConfigWithType[*ScrollElementConfig](currentElement); ok {
					closeScrollElement = true
					for _, mapping := range c.scrollContainerDatas {
						if mapping.layoutElement == currentElement {
							if scrollConfig.horizontal {
								scrollOffset.X = mapping.scrollPosition.X
							}
							if scrollConfig.vertical {
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
					if !c.Clay__ElementIsOffscreen(currentElementBoundingBox) {
						sharedConfig, ok := findElementConfigWithType[*SharedElementConfig](currentElement)
						if !ok {
							sharedConfig = &default_SharedElementConfig
						}
						renderCommand := RenderCommand{
							BoundingBox: currentElementBoundingBox,
							RenderData: BorderRenderData{
								color:        borderConfig.color,
								cornerRadius: sharedConfig.cornerRadius,
								width:        borderConfig.width,
							},
							UserData: sharedConfig.userData,
							Id:       hashNumber(currentElement.id, uint32(len(currentElement.children))).id,
						}
						c.addRenderCommand(renderCommand)
						if borderConfig.width.betweenChildren > 0 && borderConfig.color.A > 0 {
							halfGap := layoutConfig.ChildGap / 2
							borderOffset := vector2.NewFloat32(float32(layoutConfig.Padding.Left-halfGap), float32(layoutConfig.Padding.Top-halfGap))
							if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
								for i, child := range currentElement.children {
									childElement := c.layoutElements[child]
									if i > 0 {
										c.addRenderCommand(RenderCommand{
											BoundingBox: rect2.NewFloat32(
												currentElementBoundingBox.Position.Add(borderOffset).Add(scrollOffset),
												vector2.NewFloat32(float32(borderConfig.width.betweenChildren), currentElement.dimensions.Y),
											),
											RenderData: RectangleRenderData{
												backgroundColor: borderConfig.color,
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
											BoundingBox: rect2.NewFloat32(
												currentElementBoundingBox.Position.Add(scrollOffset).AddY(borderOffset.Y),
												vector2.NewFloat32(currentElement.dimensions.X, float32(borderConfig.width.betweenChildren)),
											),
											RenderData: RectangleRenderData{
												backgroundColor: borderConfig.color,
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

				dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]
				continue
			}

			// Add children to the DFS buffer
			if !elementHasConfig[*TextElementConfig](currentElement) {
				c.layoutElementTreeNodeArray1 = slicesex.Reserve(c.layoutElementTreeNodeArray1, len(c.layoutElementTreeNodeArray1)+len(currentElement.children))
				dfsBuffer = c.layoutElementTreeNodeArray1[:]
				for i, child := range currentElement.children {
					childElement := &c.layoutElements[child]
					// Alignment along non layout axis
					if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
						currentElementTreeNode.nextChildOffset.Y = float32(currentElement.layoutConfig.Padding.Top)
						whiteSpaceAroundChild := currentElement.dimensions.Y - (float32)(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom) - childElement.dimensions.Y
						switch layoutConfig.ChildAlignment.Y {
						case ALIGN_Y_TOP:
							break
						case ALIGN_Y_CENTER:
							currentElementTreeNode.nextChildOffset.Y += whiteSpaceAroundChild / 2
							break
						case ALIGN_Y_BOTTOM:
							currentElementTreeNode.nextChildOffset.Y += whiteSpaceAroundChild
							break
						}
					} else {
						currentElementTreeNode.nextChildOffset.X = float32(currentElement.layoutConfig.Padding.Left)
						whiteSpaceAroundChild := currentElement.dimensions.X - (float32)(layoutConfig.Padding.Left+layoutConfig.Padding.Right) - childElement.dimensions.X
						switch layoutConfig.ChildAlignment.X {
						case ALIGN_X_LEFT:
							break
						case ALIGN_X_CENTER:
							currentElementTreeNode.nextChildOffset.X += whiteSpaceAroundChild / 2
							break
						case ALIGN_X_RIGHT:
							currentElementTreeNode.nextChildOffset.X += whiteSpaceAroundChild
							break
						}
					}

					childPosition := vector2.NewFloat32(
						currentElementTreeNode.position.X+currentElementTreeNode.nextChildOffset.X+scrollOffset.X,
						currentElementTreeNode.position.Y+currentElementTreeNode.nextChildOffset.Y+scrollOffset.Y,
					)

					// DFS buffer elements need to be added in reverse because stack traversal happens backwards
					newNodeIndex := len(dfsBuffer) - 1 - i
					dfsBuffer[newNodeIndex] = LayoutElementTreeNode{
						layoutElement:   childElement,
						position:        childPosition,
						nextChildOffset: vector2.NewFloat32(float32(childElement.layoutConfig.Padding.Left), float32(childElement.layoutConfig.Padding.Top)),
					}
					c.treeNodeVisited[newNodeIndex] = false

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
