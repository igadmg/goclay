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

/*
Clay_Context* Clay__Context_Allocate_Arena(Clay_Arena *arena) {
    size_t totalSizeBytes = sizeof(Clay_Context);
    uintptr_t memoryAddress = (uintptr_t)arena.memory;
    // Make sure the memory address passed in for clay to use is cache line aligned
    uintptr_t nextAllocOffset = (memoryAddress % 64);
    if (nextAllocOffset + totalSizeBytes > arena.capacity)
    {
        return NULL;
    }
    arena.nextAllocation = nextAllocOffset + totalSizeBytes;
    return (Clay_Context*)(memoryAddress + nextAllocOffset);
}

string Clay__WriteStringToCharBuffer(Clay__charArray *buffer, string string) {
    for (int32 i = 0; i < string.length; i++) {
        buffer.internalArray[buffer.length + i] = string.chars[i];
    }
    buffer.length += string.length;
    return CLAY__INIT(string) { .length = string.length, .chars = (const char *)(buffer.internalArray + buffer.length - string.length) };
}
*/

/*
#ifdef CLAY_WASM
    __attribute__((import_module("clay"), import_name("measureTextFunction"))) vector2.Float32 Clay__MeasureText(Clay_StringSlice text, Clay_TextElementConfig *config, void *userData);
    __attribute__((import_module("clay"), import_name("queryScrollOffsetFunction"))) vector2.Float32 Clay__QueryScrollOffset(uint32 elementId, void *userData);
#else
    vector2.Float32 (*Clay__MeasureText)(Clay_StringSlice text, Clay_TextElementConfig *config, void *userData);
    vector2.Float32 (*Clay__QueryScrollOffset)(uint32 elementId, void *userData);
#endif
*/

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
	//c.treeNodeVisited = c.treeNodeVisited[:0]
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
					/*
					   while (sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.length > 0) {
					       float smallest = math.MaxFloat32;
					       float secondSmallest = math.MaxFloat32;
					       float widthToAdd = sizeToDistribute;
					       for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
					           Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&c.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
					           float childSize = xAxis ? child.dimensions.X : child.dimensions.Y;
					           if (Clay__FloatEqual(childSize, smallest)) { continue; }
					           if (childSize < smallest) {
					               secondSmallest = smallest;
					               smallest = childSize;
					           }
					           if (childSize > smallest) {
					               secondSmallest = min(secondSmallest, childSize);
					               widthToAdd = secondSmallest - smallest;
					           }
					       }

					       widthToAdd = min(widthToAdd, sizeToDistribute / resizableContainerBuffer.length);

					       for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
					           Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&c.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
					           float *childSize = xAxis ? &child.dimensions.X : &child.dimensions.Y;
					           float maxSize = xAxis ? child.layoutConfig.Sizing.Width.size.MinMax.Max : child.layoutConfig.Sizing.Height.size.MinMax.Max;
					           float previousWidth = *childSize;
					           if (Clay__FloatEqual(*childSize, smallest)) {
					               *childSize += widthToAdd;
					               if (*childSize >= maxSize) {
					                   *childSize = maxSize;
					                   Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
					               }
					               sizeToDistribute -= (*childSize - previousWidth);
					           }
					       }
					   }
					*/
				}
				// Sizing along the non layout axis ("off axis")
			} else {
				/*
				   for (int32 childOffset = 0; childOffset < resizableContainerBuffer.length; childOffset++) {
				       Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&c.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childOffset));
				       Clay_SizingAxis childSizing = xAxis ? childElement.layoutConfig.Sizing.Width : childElement.layoutConfig.Sizing.Height;
				       float *childSize = xAxis ? &childElement.dimensions.X : &childElement.dimensions.Y;

				       if (!xAxis && elementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_IMAGE)) {
				           continue; // Currently we don't support resizing aspect ratio images on the Y axis because it would break the ratio
				       }

				       // If we're laying out the children of a scroll panel, grow containers expand to the height of the inner content, not the outer container
				       float maxSize = parentSize - parentPadding;
				       if (elementHasConfig(parent, CLAY__ELEMENT_CONFIG_TYPE_SCROLL)) {
				           Clay_ScrollElementConfig *scrollElementConfig = findElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
				           if (((xAxis && scrollElementConfig.horizontal) || (!xAxis && scrollElementConfig.vertical))) {
				               maxSize = max(maxSize, innerContentSize);
				           }
				       }
				       if (childSizing.Type == CLAY__SIZING_TYPE_FIT) {
				           *childSize = max(childSizing.size.MinMax.Min, min(*childSize, maxSize));
				       } else if (childSizing.Type == CLAY__SIZING_TYPE_GROW) {
				           *childSize = min(maxSize, childSizing.size.MinMax.Max);
				       }
				   }
				*/
			}
		}
	}
}

/*
	string Clay__IntToString(int32 integer) {
	    if (integer == 0) {
	        return CLAY__INIT(string) { .length = 1, .chars = "0" };
	    }
	    context := GetCurrentContext();
	    char *chars = (char *)(context.dynamicStringData.internalArray + context.dynamicStringData.length);
	    int32 length = 0;
	    int32 sign = integer;

	    if (integer < 0) {
	        integer = -integer;
	    }
	    while (integer > 0) {
	        chars[length++] = (char)(integer % 10 + '0');
	        integer /= 10;
	    }

	    if (sign < 0) {
	        chars[length++] = '-';
	    }

	    // Reverse the string to get the correct order
	    for (int32 j = 0, k = length - 1; j < k; j++, k--) {
	        char temp = chars[j];
	        chars[j] = chars[k];
	        chars[k] = temp;
	    }
	    context.dynamicStringData.length += length;
	    return CLAY__INIT(string) { .length = length, .chars = chars };
	}
*/
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
				/*
					// DFS is returning upwards backwards
					bool closeScrollElement = false;
					Clay_ScrollElementConfig *scrollConfig = findElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
					if (scrollConfig) {
						closeScrollElement = true;
						for (int32 i = 0; i < c.scrollContainerDatas.length; i++) {
							Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&c.scrollContainerDatas, i);
							if (mapping.layoutElement == currentElement) {
								if (scrollConfig.horizontal) { scrollOffset.x = mapping.scrollPosition.x; }
								if (scrollConfig.vertical) { scrollOffset.y = mapping.scrollPosition.y; }
								if (c.externalScrollHandlingEnabled) {
									scrollOffset = CLAY__INIT(vector2.Float32) CLAY__DEFAULT_STRUCT;
								}
								break;
							}
						}
					}

					if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER)) {
						Clay_LayoutElementHashMapItem *currentElementData = getHashMapItem(currentElement.id);
						Clay_BoundingBox currentElementBoundingBox = currentElementData.boundingBox;

						// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
						if (!Clay__ElementIsOffscreen(&currentElementBoundingBox)) {
							Clay_SharedElementConfig *sharedConfig = elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) ? findElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig : &Clay_SharedElementConfig_DEFAULT;
							Clay_BorderElementConfig *borderConfig = findElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).borderElementConfig;
							Clay_RenderCommand renderCommand = {
									.boundingBox = currentElementBoundingBox,
									.renderData = { .border = {
										.color = borderConfig.color,
										.cornerRadius = sharedConfig.cornerRadius,
										.X = borderConfig.X
									}},
									.userData = sharedConfig.userData,
									.id = hashNumber(currentElement.id, currentElement.childrenOrTextContent.children.length).id,
									.commandType = CLAY_RENDER_COMMAND_TYPE_BORDER,
							};
							addRenderCommand(renderCommand);
							if (borderConfig.X.betweenChildren > 0 && borderConfig.color.a > 0) {
								float halfGap = layoutConfig.childGap / 2;
								vector2.Float32 borderOffset = { (float32)layoutConfig.padding.Left - halfGap, (float32)layoutConfig.padding.Top - halfGap };
								if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
									for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
										Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&c.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
										if (i > 0) {
											addRenderCommand(CLAY__INIT(Clay_RenderCommand) {
												.boundingBox = { currentElementBoundingBox.x + borderOffset.x + scrollOffset.x, currentElementBoundingBox.y + scrollOffset.y, (float32)borderConfig.X.betweenChildren, currentElement.dimensions.Y },
												.renderData = { .rectangle = {
													.backgroundColor = borderConfig.color,
												} },
												.userData = sharedConfig.userData,
												.id = hashNumber(currentElement.id, currentElement.childrenOrTextContent.children.length + 1 + i).id,
												.commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
											});
										}
										borderOffset.x += (childElement.dimensions.X + (float32)layoutConfig.childGap);
									}
								} else {
									for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
										Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&c.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
										if (i > 0) {
											addRenderCommand(CLAY__INIT(Clay_RenderCommand) {
												.boundingBox = { currentElementBoundingBox.x + scrollOffset.x, currentElementBoundingBox.y + borderOffset.y + scrollOffset.y, currentElement.dimensions.X, (float32)borderConfig.X.betweenChildren },
												.renderData = { .rectangle = {
														.backgroundColor = borderConfig.color,
												} },
												.userData = sharedConfig.userData,
												.id = hashNumber(currentElement.id, currentElement.childrenOrTextContent.children.length + 1 + i).id,
												.commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
											});
										}
										borderOffset.y += (childElement.dimensions.Y + (float32)layoutConfig.childGap);
									}
								}
							}
						}
					}
					// This exists because the scissor needs to end _after_ borders between elements
					if (closeScrollElement) {
						addRenderCommand(CLAY__INIT(Clay_RenderCommand) {
							.id = hashNumber(currentElement.id, rootElement.childrenOrTextContent.children.length + 11).id,
							.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
						});
					}
				*/

				dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]
				continue
			}
			/*
				// Add children to the DFS buffer
				if (!elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
					dfsBuffer.length += currentElement.childrenOrTextContent.children.length;
					for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
						Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&c.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
						// Alignment along non layout axis
						if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
							currentElementTreeNode.nextChildOffset.y = currentElement.layoutConfig.padding.Top;
							float whiteSpaceAroundChild = currentElement.dimensions.Y - (float32)(layoutConfig.padding.Top + layoutConfig.padding.Bottom) - childElement.dimensions.Y;
							switch (layoutConfig.childAlignment.y) {
								case CLAY_ALIGN_Y_TOP: break;
								case CLAY_ALIGN_Y_CENTER: currentElementTreeNode.nextChildOffset.y += whiteSpaceAroundChild / 2; break;
								case CLAY_ALIGN_Y_BOTTOM: currentElementTreeNode.nextChildOffset.y += whiteSpaceAroundChild; break;
							}
						} else {
							currentElementTreeNode.nextChildOffset.x = currentElement.layoutConfig.padding.Left;
							float whiteSpaceAroundChild = currentElement.dimensions.X - (float32)(layoutConfig.padding.Left + layoutConfig.padding.Right) - childElement.dimensions.X;
							switch (layoutConfig.childAlignment.x) {
								case CLAY_ALIGN_X_LEFT: break;
								case CLAY_ALIGN_X_CENTER: currentElementTreeNode.nextChildOffset.x += whiteSpaceAroundChild / 2; break;
								case CLAY_ALIGN_X_RIGHT: currentElementTreeNode.nextChildOffset.x += whiteSpaceAroundChild; break;
							}
						}

						vector2.Float32 childPosition = {
							currentElementTreeNode.position.x + currentElementTreeNode.nextChildOffset.x + scrollOffset.x,
							currentElementTreeNode.position.y + currentElementTreeNode.nextChildOffset.y + scrollOffset.y,
						};

						// DFS buffer elements need to be added in reverse because stack traversal happens backwards
						uint32 newNodeIndex = dfsBuffer.length - 1 - i;
						dfsBuffer.internalArray[newNodeIndex] = CLAY__INIT(Clay__LayoutElementTreeNode) {
							.layoutElement = childElement,
							.position = { childPosition.x, childPosition.y },
							.nextChildOffset = { .x = (float32)childElement.layoutConfig.padding.Left, .y = (float32)childElement.layoutConfig.padding.Top },
						};
						c.treeNodeVisited.internalArray[newNodeIndex] = false;

						// Update parent offsets
						if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
							currentElementTreeNode.nextChildOffset.x += childElement.dimensions.X + (float32)layoutConfig.childGap;
						} else {
							currentElementTreeNode.nextChildOffset.y += childElement.dimensions.Y + (float32)layoutConfig.childGap;
						}
					}
				}
			*/
		}

		if root.clipElementId != 0 {
			c.addRenderCommand(RenderCommand{
				Id:         hashNumber(rootElement.id, uint32(len(rootElement.children))+11).id,
				RenderData: ScissorsEndData{},
			})
		}
	}
}

/*
#pragma region DebugTools
colorex.RGBA CLAY__DEBUGVIEW_COLOR_1 = {58, 56, 52, 255};
colorex.RGBA CLAY__DEBUGVIEW_COLOR_2 = {62, 60, 58, 255};
colorex.RGBA CLAY__DEBUGVIEW_COLOR_3 = {141, 133, 135, 255};
colorex.RGBA CLAY__DEBUGVIEW_COLOR_4 = {238, 226, 231, 255};
colorex.RGBA CLAY__DEBUGVIEW_COLOR_SELECTED_ROW = {102, 80, 78, 255};
const int32 CLAY__DEBUGVIEW_ROW_HEIGHT = 30;
const int32 CLAY__DEBUGVIEW_OUTER_PADDING = 10;
const int32 CLAY__DEBUGVIEW_INDENT_WIDTH = 16;
Clay_TextElementConfig Clay__DebugView_TextNameConfig = {.textColor = {238, 226, 231, 255}, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE };
Clay_LayoutConfig Clay__DebugView_ScrollViewItemLayoutConfig = CLAY__DEFAULT_STRUCT;

	typedef struct {
	    string label;
	    colorex.RGBA color;
	} Clay__DebugElementConfigTypeLabelConfig;

	Clay__DebugElementConfigTypeLabelConfig Clay__DebugGetElementConfigTypeLabel(Clay__ElementConfigType type) {
	    switch (type) {
	        case CLAY__ELEMENT_CONFIG_TYPE_SHARED: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Shared"), {243,134,48,255} };
	        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Text"), {105,210,231,255} };
	        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Image"), {121,189,154,255} };
	        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Floating"), {250,105,0,255} };
	        case CLAY__ELEMENT_CONFIG_TYPE_SCROLL: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) {CLAY_STRING("Scroll"), {242, 196, 90, 255} };
	        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) {CLAY_STRING("Border"), {108, 91, 123, 255} };
	        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Custom"), {11,72,107,255} };
	        default: break;
	    }
	    return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Error"), {0,0,0,255} };
	}

	typedef struct {
	    int32 rowCount;
	    int32 selectedElementRowIndex;
	} Clay__RenderDebugLayoutData;

// Returns row count

	Clay__RenderDebugLayoutData Clay__RenderDebugLayoutElementsList(int32 initialRootsLength, int32 highlightedRowIndex) {
	    context := GetCurrentContext();
	    Clay__int32_tArray dfsBuffer = context.reusableElementIndexBuffer;
	    Clay__DebugView_ScrollViewItemLayoutConfig = CLAY__INIT(Clay_LayoutConfig) { .sizing = { .Y = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT) }, .childGap = 6, .childAlignment = { .y = CLAY_ALIGN_Y_CENTER }};
	    Clay__RenderDebugLayoutData layoutData = CLAY__DEFAULT_STRUCT;

	    uint32 highlightedElementId = 0;

	    for (int32 rootIndex = 0; rootIndex < initialRootsLength; ++rootIndex) {
	        dfsBuffer.length = 0;
	        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex);
	        Clay__int32_tArray_Add(&dfsBuffer, (int32)root.layoutElementIndex);
	        context.treeNodeVisited.internalArray[0] = false;
	        if (rootIndex > 0) {
	            CLAY({ .id = CLAY_IDI("Clay__DebugView_EmptyRowOuter", rootIndex), .layout = { .sizing = {.X = CLAY_SIZING_GROW(0)}, .padding = {CLAY__DEBUGVIEW_INDENT_WIDTH / 2, 0, 0, 0} } }) {
	                CLAY({ .id = CLAY_IDI("Clay__DebugView_EmptyRow", rootIndex), .layout = { .sizing = { .X = CLAY_SIZING_GROW(0), .Y = CLAY_SIZING_FIXED((float32)CLAY__DEBUGVIEW_ROW_HEIGHT) }}, .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .X = { .top = 1 } } }) {}
	            }
	            layoutData.rowCount++;
	        }
	        while (dfsBuffer.length > 0) {
	            int32 currentElementIndex = Clay__int32_tArray_GetValue(&dfsBuffer, (int)dfsBuffer.length - 1);
	            Clay_LayoutElement *currentElement = Clay_LayoutElementArray_Get(&context.layoutElements, (int)currentElementIndex);
	            if (context.treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
	                if (!elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && currentElement.childrenOrTextContent.children.length > 0) {
	                    Clay__CloseElement();
	                    Clay__CloseElement();
	                    Clay__CloseElement();
	                }
	                dfsBuffer.length--;
	                continue;
	            }

	            if (highlightedRowIndex == layoutData.rowCount) {
	                if (context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
	                    context.debugSelectedElementId = currentElement.id;
	                }
	                highlightedElementId = currentElement.id;
	            }

	            context.treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
	            Clay_LayoutElementHashMapItem *currentElementData = getHashMapItem(currentElement.id);
	            bool offscreen = Clay__ElementIsOffscreen(&currentElementData.boundingBox);
	            if (context.debugSelectedElementId == currentElement.id) {
	                layoutData.selectedElementRowIndex = layoutData.rowCount;
	            }
	            CLAY({ .id = CLAY_IDI("Clay__DebugView_ElementOuter", currentElement.id), .layout = Clay__DebugView_ScrollViewItemLayoutConfig }) {
	                // Collapse icon / button
	                if (!(elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement.childrenOrTextContent.children.length == 0)) {
	                    CLAY({
	                        .id = CLAY_IDI("Clay__DebugView_CollapseElement", currentElement.id),
	                        .layout = { .sizing = {CLAY_SIZING_FIXED(16), CLAY_SIZING_FIXED(16)}, .childAlignment = { CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER} },
	                        .cornerRadius = CLAY_CORNER_RADIUS(4),
	                        .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .X = {1, 1, 1, 1, 0} },
	                    }) {
	                        CLAY_TEXT((currentElementData && currentElementData.debugData.collapsed) ? CLAY_STRING("+") : CLAY_STRING("-"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
	                    }
	                } else { // Square dot for empty containers
	                    CLAY({ .layout = { .sizing = {CLAY_SIZING_FIXED(16), CLAY_SIZING_FIXED(16)}, .childAlignment = { CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER } } }) {
	                        CLAY({ .layout = { .sizing = {CLAY_SIZING_FIXED(8), CLAY_SIZING_FIXED(8)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3, .cornerRadius = CLAY_CORNER_RADIUS(2) }) {}
	                    }
	                }
	                // Collisions and offscreen info
	                if (currentElementData) {
	                    if (currentElementData.debugData.collision) {
	                        CLAY({ .layout = { .padding = { 8, 8, 2, 2 }}, .border = { .color = {177, 147, 8, 255}, .X = {1, 1, 1, 1, 0} } }) {
	                            CLAY_TEXT(CLAY_STRING("Duplicate ID"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }));
	                        }
	                    }
	                    if (offscreen) {
	                        CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .border = {  .color = CLAY__DEBUGVIEW_COLOR_3, .X = { 1, 1, 1, 1, 0} } }) {
	                            CLAY_TEXT(CLAY_STRING("Offscreen"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }));
	                        }
	                    }
	                }
	                string idString = context.layoutElementIdStrings.internalArray[currentElementIndex];
	                if (idString.length > 0) {
	                    CLAY_TEXT(idString, offscreen ? CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }) : &Clay__DebugView_TextNameConfig);
	                }
	                for (int32 elementConfigIndex = 0; elementConfigIndex < currentElement.elementConfigs.length; ++elementConfigIndex) {
	                    Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, elementConfigIndex);
	                    if (elementConfig.type == CLAY__ELEMENT_CONFIG_TYPE_SHARED) {
	                        colorex.RGBA labelColor = {243,134,48,90};
	                        labelColor.a = 90;
	                        colorex.RGBA backgroundColor = elementConfig.config.sharedElementConfig.backgroundColor;
	                        Clay_CornerRadius radius = elementConfig.config.sharedElementConfig.cornerRadius;
	                        if (backgroundColor.a > 0) {
	                            CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = labelColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = labelColor, .X = { 1, 1, 1, 1, 0} } }) {
	                                CLAY_TEXT(CLAY_STRING("Color"), CLAY_TEXT_CONFIG({ .textColor = offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
	                            }
	                        }
	                        if (radius.bottomLeft > 0) {
	                            CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = labelColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = labelColor, .X = { 1, 1, 1, 1, 0 } } }) {
	                                CLAY_TEXT(CLAY_STRING("Radius"), CLAY_TEXT_CONFIG({ .textColor = offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
	                            }
	                        }
	                        continue;
	                    }
	                    Clay__DebugElementConfigTypeLabelConfig config = Clay__DebugGetElementConfigTypeLabel(elementConfig.type);
	                    colorex.RGBA backgroundColor = config.color;
	                    backgroundColor.a = 90;
	                    CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = backgroundColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = config.color, .X = { 1, 1, 1, 1, 0 } } }) {
	                        CLAY_TEXT(config.label, CLAY_TEXT_CONFIG({ .textColor = offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
	                    }
	                }
	            }

	            // Render the text contents below the element as a non-interactive row
	            if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
	                layoutData.rowCount++;
	                Clay__TextElementData *textElementData = currentElement.childrenOrTextContent.textElementData;
	                Clay_TextElementConfig *rawTextConfig = offscreen ? CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }) : &Clay__DebugView_TextNameConfig;
	                CLAY({ .layout = { .sizing = { .Y = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .childAlignment = { .y = CLAY_ALIGN_Y_CENTER } } }) {
	                    CLAY({ .layout = { .sizing = {.X = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_INDENT_WIDTH + 16) } } }) {}
	                    CLAY_TEXT(CLAY_STRING("\""), rawTextConfig);
	                    CLAY_TEXT(textElementData.text.length > 40 ? (CLAY__INIT(string) { .length = 40, .chars = textElementData.text.chars }) : textElementData.text, rawTextConfig);
	                    if (textElementData.text.length > 40) {
	                        CLAY_TEXT(CLAY_STRING("..."), rawTextConfig);
	                    }
	                    CLAY_TEXT(CLAY_STRING("\""), rawTextConfig);
	                }
	            } else if (currentElement.childrenOrTextContent.children.length > 0) {
	                Clay__OpenElement();
	                Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) { .layout = { .padding = { .left = 8 } } });
	                Clay__OpenElement();
	                Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) { .layout = { .padding = { .left = CLAY__DEBUGVIEW_INDENT_WIDTH }}, .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .X = { .left = 1 } }});
	                Clay__OpenElement();
	                Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) { .layout = { .layoutDirection = CLAY_TOP_TO_BOTTOM } });
	            }

	            layoutData.rowCount++;
	            if (!(elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (currentElementData && currentElementData.debugData.collapsed))) {
	                for (int32 i = currentElement.childrenOrTextContent.children.length - 1; i >= 0; --i) {
	                    Clay__int32_tArray_Add(&dfsBuffer, currentElement.childrenOrTextContent.children.elements[i]);
	                    context.treeNodeVisited.internalArray[dfsBuffer.length - 1] = false; // TODO needs to be ranged checked
	                }
	            }
	        }
	    }

	    if (context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
	        Clay_ElementId collapseButtonId = hashString(CLAY_STRING("Clay__DebugView_CollapseElement"), 0, 0);
	        for (int32 i = (int)context.pointerOverIds.length - 1; i >= 0; i--) {
	            Clay_ElementId *elementId = Clay__ElementIdArray_Get(&context.pointerOverIds, i);
	            if (elementId.baseId == collapseButtonId.baseId) {
	                Clay_LayoutElementHashMapItem *highlightedItem = getHashMapItem(elementId.offset);
	                highlightedItem.debugData.collapsed = !highlightedItem.debugData.collapsed;
	                break;
	            }
	        }
	    }

	    if (highlightedElementId) {
	        CLAY({ .id = CLAY_ID("Clay__DebugView_ElementHighlight"), .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .floating = { .parentId = highlightedElementId, .zIndex = 32767, .pointerCaptureMode = CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, .attachTo = CLAY_ATTACH_TO_ELEMENT_WITH_ID } }) {
	            CLAY({ .id = CLAY_ID("Clay__DebugView_ElementHighlightRectangle"), .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .backgroundColor = debugViewHighlightColor }) {}
	        }
	    }
	    return layoutData;
	}

	void Clay__RenderDebugLayoutSizing(Clay_SizingAxis sizing, Clay_TextElementConfig *infoTextConfig) {
	    string sizingLabel = CLAY_STRING("GROW");
	    if (sizing.Type == CLAY__SIZING_TYPE_FIT) {
	        sizingLabel = CLAY_STRING("FIT");
	    } else if (sizing.Type == CLAY__SIZING_TYPE_PERCENT) {
	        sizingLabel = CLAY_STRING("PERCENT");
	    }
	    CLAY_TEXT(sizingLabel, infoTextConfig);
	    if (sizing.Type == CLAY__SIZING_TYPE_GROW || sizing.Type == CLAY__SIZING_TYPE_FIT) {
	        CLAY_TEXT(CLAY_STRING("("), infoTextConfig);
	        if (sizing.size.MinMax.Min != 0) {
	            CLAY_TEXT(CLAY_STRING("min: "), infoTextConfig);
	            CLAY_TEXT(Clay__IntToString(sizing.size.MinMax.Min), infoTextConfig);
	            if (sizing.size.MinMax.Max != math.MaxFloat32) {
	                CLAY_TEXT(CLAY_STRING(", "), infoTextConfig);
	            }
	        }
	        if (sizing.size.MinMax.Max != math.MaxFloat32) {
	            CLAY_TEXT(CLAY_STRING("max: "), infoTextConfig);
	            CLAY_TEXT(Clay__IntToString(sizing.size.MinMax.Max), infoTextConfig);
	        }
	        CLAY_TEXT(CLAY_STRING(")"), infoTextConfig);
	    }
	}

	func (c *Context) Clay__RenderDebugViewElementConfigHeader(elementId string, Type Clay__ElementConfigType ) {
		    Clay__DebugElementConfigTypeLabelConfig config = Clay__DebugGetElementConfigTypeLabel(type);
		    colorex.RGBA backgroundColor = config.color;
		    backgroundColor.a = 90;
		    CLAY({ .layout = { .sizing = { .X = CLAY_SIZING_GROW(0) }, .padding = CLAY_PADDING_ALL(CLAY__DEBUGVIEW_OUTER_PADDING), .childAlignment = { .y = CLAY_ALIGN_Y_CENTER } } }) {
		        CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = backgroundColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = config.color, .X = { 1, 1, 1, 1, 0 } } }) {
		            CLAY_TEXT(config.label, CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
		        }
		        CLAY({ .layout = { .sizing = { .X = CLAY_SIZING_GROW(0) } } }) {}
		        CLAY_TEXT(elementId, CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE }));
		    }
		}

		void Clay__RenderDebugViewColor(colorex.RGBA color, Clay_TextElementConfig *textConfig) {
		    CLAY({ .layout = { .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
		        CLAY_TEXT(CLAY_STRING("{ r: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(color.r), textConfig);
		        CLAY_TEXT(CLAY_STRING(", g: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(color.g), textConfig);
		        CLAY_TEXT(CLAY_STRING(", b: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(color.b), textConfig);
		        CLAY_TEXT(CLAY_STRING(", a: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(color.a), textConfig);
		        CLAY_TEXT(CLAY_STRING(" }"), textConfig);
		        CLAY({ .layout = { .sizing = { .X = CLAY_SIZING_FIXED(10) } } }) {}
		        CLAY({ .layout = { .sizing = { CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 8), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 8)} }, .backgroundColor = color, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = CLAY__DEBUGVIEW_COLOR_4, .X = { 1, 1, 1, 1, 0 } } }) {}
		    }
		}

		void Clay__RenderDebugViewCornerRadius(Clay_CornerRadius cornerRadius, Clay_TextElementConfig *textConfig) {
		    CLAY({ .layout = { .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
		        CLAY_TEXT(CLAY_STRING("{ topLeft: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(cornerRadius.topLeft), textConfig);
		        CLAY_TEXT(CLAY_STRING(", topRight: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(cornerRadius.topRight), textConfig);
		        CLAY_TEXT(CLAY_STRING(", bottomLeft: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(cornerRadius.bottomLeft), textConfig);
		        CLAY_TEXT(CLAY_STRING(", bottomRight: "), textConfig);
		        CLAY_TEXT(Clay__IntToString(cornerRadius.bottomRight), textConfig);
		        CLAY_TEXT(CLAY_STRING(" }"), textConfig);
		    }
		}

		void HandleDebugViewCloseButtonInteraction(Clay_ElementId elementId, Clay_PointerData pointerInfo, intptr_t userData) {
		    context := GetCurrentContext();
		    (void) elementId; (void) pointerInfo; (void) userData;
		    if (pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
		        context.debugModeEnabled = false;
		    }
		}
*/
func (c *Context) Clay__RenderDebugView() {
	/*
	   Clay_ElementId closeButtonId = hashString(CLAY_STRING("Clay__DebugViewTopHeaderCloseButtonOuter"), 0, 0);
	   if (c.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
	       for (int32 i = 0; i < c.pointerOverIds.length; ++i) {
	           Clay_ElementId *elementId = Clay__ElementIdArray_Get(&c.pointerOverIds, i);
	           if (elementId.id == closeButtonId.id) {
	               c.debugModeEnabled = false;
	               return;
	           }
	       }
	   }

	   uint32 initialRootsLength = c.layoutElementTreeRoots.length;
	   uint32 initialElementsLength = c.layoutElements.length;
	   Clay_TextElementConfig *infoTextConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
	   Clay_TextElementConfig *infoTitleConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
	   Clay_ElementId scrollId = hashString(CLAY_STRING("Clay__DebugViewOuterScrollPane"), 0, 0);
	   float scrollYOffset = 0;
	   bool pointerInDebugView = c.pointerInfo.position.y < c.layoutDimensions.Y - 300;
	   for (int32 i = 0; i < c.scrollContainerDatas.length; ++i) {
	       Clay__ScrollContainerDataInternal *scrollContainerData = Clay__ScrollContainerDataInternalArray_Get(&c.scrollContainerDatas, i);
	       if (scrollContainerData.elementId == scrollId.id) {
	           if (!c.externalScrollHandlingEnabled) {
	               scrollYOffset = scrollContainerData.scrollPosition.y;
	           } else {
	               pointerInDebugView = c.pointerInfo.position.y + scrollContainerData.scrollPosition.y < c.layoutDimensions.Y - 300;
	           }
	           break;
	       }
	   }
	   int32 highlightedRow = pointerInDebugView
	           ? (int32)((c.pointerInfo.position.y - scrollYOffset) / (float32)CLAY__DEBUGVIEW_ROW_HEIGHT) - 1
	           : -1;
	   if (c.pointerInfo.position.x < c.layoutDimensions.X - (float32)debugViewWidth) {
	       highlightedRow = -1;
	   }
	   Clay__RenderDebugLayoutData layoutData = CLAY__DEFAULT_STRUCT;
	   CLAY({ .id = CLAY_ID("Clay__DebugView"),
	        .layout = { .sizing = { CLAY_SIZING_FIXED((float32)debugViewWidth) , CLAY_SIZING_FIXED(c.layoutDimensions.Y) }, .layoutDirection = CLAY_TOP_TO_BOTTOM },
	       .floating = { .zIndex = 32765, .attachPoints = { .element = CLAY_ATTACH_POINT_LEFT_CENTER, .parent = CLAY_ATTACH_POINT_RIGHT_CENTER }, .attachTo = CLAY_ATTACH_TO_ROOT },
	       .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .X = { .bottom = 1 } }
	   }) {
	       CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_2 }) {
	           CLAY_TEXT(CLAY_STRING("Clay Debug Tools"), infoTextConfig);
	           CLAY({ .layout = { .sizing = { .X = CLAY_SIZING_GROW(0) } } }) {}
	           // Close button
	           CLAY({
	               .layout = { .sizing = {CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10)}, .childAlignment = {CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER} },
	               .backgroundColor = {217,91,67,80},
	               .cornerRadius = CLAY_CORNER_RADIUS(4),
	               .border = { .color = { 217,91,67,255 }, .X = { 1, 1, 1, 1, 0 } },
	           }) {
	               Clay_OnHover(HandleDebugViewCloseButtonInteraction, 0);
	               CLAY_TEXT(CLAY_STRING("x"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
	           }
	       }
	       CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(1)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3 } ) {}
	       CLAY({ .id = scrollId, .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .scroll = { .horizontal = true, .vertical = true } }) {
	           CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)}, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = ((initialElementsLength + initialRootsLength) & 1) == 0 ? CLAY__DEBUGVIEW_COLOR_2 : CLAY__DEBUGVIEW_COLOR_1 }) {
	               Clay_ElementId panelContentsId = hashString(CLAY_STRING("Clay__DebugViewPaneOuter"), 0, 0);
	               // Element list
	               CLAY({ .id = panelContentsId, .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .floating = { .zIndex = 32766, .pointerCaptureMode = CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, .attachTo = CLAY_ATTACH_TO_PARENT } }) {
	                   CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)}, .padding = { CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                       layoutData = Clay__RenderDebugLayoutElementsList((int32)initialRootsLength, highlightedRow);
	                   }
	               }
	               float contentWidth = getHashMapItem(panelContentsId.id).layoutElement.dimensions.X;
	               CLAY({ .layout = { .sizing = {.X = CLAY_SIZING_FIXED(contentWidth) }, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {}
	               for (int32 i = 0; i < layoutData.rowCount; i++) {
	                   colorex.RGBA rowColor = (i & 1) == 0 ? CLAY__DEBUGVIEW_COLOR_2 : CLAY__DEBUGVIEW_COLOR_1;
	                   if (i == layoutData.selectedElementRowIndex) {
	                       rowColor = CLAY__DEBUGVIEW_COLOR_SELECTED_ROW;
	                   }
	                   if (i == highlightedRow) {
	                       rowColor.r *= 1.25f;
	                       rowColor.g *= 1.25f;
	                       rowColor.b *= 1.25f;
	                   }
	                   CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = rowColor } ) {}
	               }
	           }
	       }
	       CLAY({ .layout = { .sizing = {.X = CLAY_SIZING_GROW(0), .Y = CLAY_SIZING_FIXED(1)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3 }) {}
	       if (c.debugSelectedElementId != 0) {
	           Clay_LayoutElementHashMapItem *selectedItem = getHashMapItem(c.debugSelectedElementId);
	           CLAY({
	               .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(300)}, .layoutDirection = CLAY_TOP_TO_BOTTOM },
	               .backgroundColor = CLAY__DEBUGVIEW_COLOR_2 ,
	               .scroll = { .vertical = true },
	               .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .X = { .betweenChildren = 1 } }
	           }) {
	               CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT + 8)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
	                   CLAY_TEXT(CLAY_STRING("Layout Config"), infoTextConfig);
	                   CLAY({ .layout = { .sizing = { .X = CLAY_SIZING_GROW(0) } } }) {}
	                   if (selectedItem.elementId.stringId.length != 0) {
	                       CLAY_TEXT(selectedItem.elementId.stringId, infoTitleConfig);
	                       if (selectedItem.elementId.offset != 0) {
	                           CLAY_TEXT(CLAY_STRING(" ("), infoTitleConfig);
	                           CLAY_TEXT(Clay__IntToString(selectedItem.elementId.offset), infoTitleConfig);
	                           CLAY_TEXT(CLAY_STRING(")"), infoTitleConfig);
	                       }
	                   }
	               }
	               Clay_Padding attributeConfigPadding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 8, 8};
	               // Clay_LayoutConfig debug info
	               CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                   // .boundingBox
	                   CLAY_TEXT(CLAY_STRING("Bounding Box"), infoTitleConfig);
	                   CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                       CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(selectedItem.boundingBox.x), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(selectedItem.boundingBox.y), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", width: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(selectedItem.boundingBox.X), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(selectedItem.boundingBox.Y), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                   }
	                   // .layoutDirection
	                   CLAY_TEXT(CLAY_STRING("Layout Direction"), infoTitleConfig);
	                   Clay_LayoutConfig *layoutConfig = selectedItem.layoutElement.layoutConfig;
	                   CLAY_TEXT(layoutConfig.layoutDirection == CLAY_TOP_TO_BOTTOM ? CLAY_STRING("TOP_TO_BOTTOM") : CLAY_STRING("LEFT_TO_RIGHT"), infoTextConfig);
	                   // .sizing
	                   CLAY_TEXT(CLAY_STRING("Sizing"), infoTitleConfig);
	                   CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                       CLAY_TEXT(CLAY_STRING("width: "), infoTextConfig);
	                       Clay__RenderDebugLayoutSizing(layoutConfig.Sizing.Width, infoTextConfig);
	                   }
	                   CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                       CLAY_TEXT(CLAY_STRING("height: "), infoTextConfig);
	                       Clay__RenderDebugLayoutSizing(layoutConfig.Sizing.Height, infoTextConfig);
	                   }
	                   // .padding
	                   CLAY_TEXT(CLAY_STRING("Padding"), infoTitleConfig);
	                   CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoPadding") }) {
	                       CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(layoutConfig.padding.Left), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(layoutConfig.padding.Right), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(layoutConfig.padding.Top), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig);
	                       CLAY_TEXT(Clay__IntToString(layoutConfig.padding.Bottom), infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                   }
	                   // .childGap
	                   CLAY_TEXT(CLAY_STRING("Child Gap"), infoTitleConfig);
	                   CLAY_TEXT(Clay__IntToString(layoutConfig.childGap), infoTextConfig);
	                   // .childAlignment
	                   CLAY_TEXT(CLAY_STRING("Child Alignment"), infoTitleConfig);
	                   CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                       CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
	                       string alignX = CLAY_STRING("LEFT");
	                       if (layoutConfig.childAlignment.x == CLAY_ALIGN_X_CENTER) {
	                           alignX = CLAY_STRING("CENTER");
	                       } else if (layoutConfig.childAlignment.x == CLAY_ALIGN_X_RIGHT) {
	                           alignX = CLAY_STRING("RIGHT");
	                       }
	                       CLAY_TEXT(alignX, infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
	                       string alignY = CLAY_STRING("TOP");
	                       if (layoutConfig.childAlignment.y == CLAY_ALIGN_Y_CENTER) {
	                           alignY = CLAY_STRING("CENTER");
	                       } else if (layoutConfig.childAlignment.y == CLAY_ALIGN_Y_BOTTOM) {
	                           alignY = CLAY_STRING("BOTTOM");
	                       }
	                       CLAY_TEXT(alignY, infoTextConfig);
	                       CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                   }
	               }
	               for (int32 elementConfigIndex = 0; elementConfigIndex < selectedItem.layoutElement.elementConfigs.length; ++elementConfigIndex) {
	                   Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&selectedItem.layoutElement.elementConfigs, elementConfigIndex);
	                   Clay__RenderDebugViewElementConfigHeader(selectedItem.elementId.stringId, elementConfig.type);
	                   switch (elementConfig.type) {
	                       case CLAY__ELEMENT_CONFIG_TYPE_SHARED: {
	                           Clay_SharedElementConfig *sharedConfig = elementConfig.config.sharedElementConfig;
	                           CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM }}) {
	                               // .backgroundColor
	                               CLAY_TEXT(CLAY_STRING("Background Color"), infoTitleConfig);
	                               Clay__RenderDebugViewColor(sharedConfig.backgroundColor, infoTextConfig);
	                               // .cornerRadius
	                               CLAY_TEXT(CLAY_STRING("Corner Radius"), infoTitleConfig);
	                               Clay__RenderDebugViewCornerRadius(sharedConfig.cornerRadius, infoTextConfig);
	                           }
	                           break;
	                       }
	                       case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
	                           Clay_TextElementConfig *textConfig = elementConfig.config.textElementConfig;
	                           CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                               // .fontSize
	                               CLAY_TEXT(CLAY_STRING("Font Size"), infoTitleConfig);
	                               CLAY_TEXT(Clay__IntToString(textConfig.fontSize), infoTextConfig);
	                               // .fontId
	                               CLAY_TEXT(CLAY_STRING("Font ID"), infoTitleConfig);
	                               CLAY_TEXT(Clay__IntToString(textConfig.fontId), infoTextConfig);
	                               // .lineHeight
	                               CLAY_TEXT(CLAY_STRING("Line Height"), infoTitleConfig);
	                               CLAY_TEXT(textConfig.lineHeight == 0 ? CLAY_STRING("auto") : Clay__IntToString(textConfig.lineHeight), infoTextConfig);
	                               // .letterSpacing
	                               CLAY_TEXT(CLAY_STRING("Letter Spacing"), infoTitleConfig);
	                               CLAY_TEXT(Clay__IntToString(textConfig.letterSpacing), infoTextConfig);
	                               // .wrapMode
	                               CLAY_TEXT(CLAY_STRING("Wrap Mode"), infoTitleConfig);
	                               string wrapMode = CLAY_STRING("WORDS");
	                               if (textConfig.wrapMode == CLAY_TEXT_WRAP_NONE) {
	                                   wrapMode = CLAY_STRING("NONE");
	                               } else if (textConfig.wrapMode == CLAY_TEXT_WRAP_NEWLINES) {
	                                   wrapMode = CLAY_STRING("NEWLINES");
	                               }
	                               CLAY_TEXT(wrapMode, infoTextConfig);
	                               // .textAlignment
	                               CLAY_TEXT(CLAY_STRING("Text Alignment"), infoTitleConfig);
	                               string textAlignment = CLAY_STRING("LEFT");
	                               if (textConfig.textAlignment == CLAY_TEXT_ALIGN_CENTER) {
	                                   textAlignment = CLAY_STRING("CENTER");
	                               } else if (textConfig.textAlignment == CLAY_TEXT_ALIGN_RIGHT) {
	                                   textAlignment = CLAY_STRING("RIGHT");
	                               }
	                               CLAY_TEXT(textAlignment, infoTextConfig);
	                               // .textColor
	                               CLAY_TEXT(CLAY_STRING("Text Color"), infoTitleConfig);
	                               Clay__RenderDebugViewColor(textConfig.textColor, infoTextConfig);
	                           }
	                           break;
	                       }
	                       case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
	                           Clay_ImageElementConfig *imageConfig = elementConfig.config.imageElementConfig;
	                           CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoImageBody"), .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                               // .sourceDimensions
	                               CLAY_TEXT(CLAY_STRING("Source Dimensions"), infoTitleConfig);
	                               CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoImageDimensions") }) {
	                                   CLAY_TEXT(CLAY_STRING("{ width: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(imageConfig.sourceDimensions.X), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(imageConfig.sourceDimensions.Y), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                               }
	                               // Image Preview
	                               CLAY_TEXT(CLAY_STRING("Preview"), infoTitleConfig);
	                               CLAY({ .layout = { .sizing = { .X = CLAY_SIZING_GROW(0, imageConfig.sourceDimensions.X) }}, .image = *imageConfig }) {}
	                           }
	                           break;
	                       }
	                       case CLAY__ELEMENT_CONFIG_TYPE_SCROLL: {
	                           Clay_ScrollElementConfig *scrollConfig = elementConfig.config.scrollElementConfig;
	                           CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                               // .vertical
	                               CLAY_TEXT(CLAY_STRING("Vertical"), infoTitleConfig);
	                               CLAY_TEXT(scrollConfig.vertical ? CLAY_STRING("true") : CLAY_STRING("false") , infoTextConfig);
	                               // .horizontal
	                               CLAY_TEXT(CLAY_STRING("Horizontal"), infoTitleConfig);
	                               CLAY_TEXT(scrollConfig.horizontal ? CLAY_STRING("true") : CLAY_STRING("false") , infoTextConfig);
	                           }
	                           break;
	                       }
	                       case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: {
	                           Clay_FloatingElementConfig *floatingConfig = elementConfig.config.floatingElementConfig;
	                           CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                               // .offset
	                               CLAY_TEXT(CLAY_STRING("Offset"), infoTitleConfig);
	                               CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                                   CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(floatingConfig.offset.x), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(floatingConfig.offset.y), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                               }
	                               // .expand
	                               CLAY_TEXT(CLAY_STRING("Expand"), infoTitleConfig);
	                               CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                                   CLAY_TEXT(CLAY_STRING("{ width: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(floatingConfig.expand.X), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(floatingConfig.expand.Y), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                               }
	                               // .zIndex
	                               CLAY_TEXT(CLAY_STRING("z-index"), infoTitleConfig);
	                               CLAY_TEXT(Clay__IntToString(floatingConfig.zIndex), infoTextConfig);
	                               // .parentId
	                               CLAY_TEXT(CLAY_STRING("Parent"), infoTitleConfig);
	                               Clay_LayoutElementHashMapItem *hashItem = getHashMapItem(floatingConfig.parentId);
	                               CLAY_TEXT(hashItem.elementId.stringId, infoTextConfig);
	                           }
	                           break;
	                       }
	                       case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
	                           Clay_BorderElementConfig *borderConfig = elementConfig.config.borderElementConfig;
	                           CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoBorderBody"), .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
	                               CLAY_TEXT(CLAY_STRING("Border Widths"), infoTitleConfig);
	                               CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                                   CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(borderConfig.X.left), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(borderConfig.X.right), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(borderConfig.X.top), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig);
	                                   CLAY_TEXT(Clay__IntToString(borderConfig.X.bottom), infoTextConfig);
	                                   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
	                               }
	                               // .textColor
	                               CLAY_TEXT(CLAY_STRING("Border Color"), infoTitleConfig);
	                               Clay__RenderDebugViewColor(borderConfig.color, infoTextConfig);
	                           }
	                           break;
	                       }
	                       case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM:
	                       default: break;
	                   }
	               }
	           }
	       } else {
	           CLAY({ .id = CLAY_ID("Clay__DebugViewWarningsScrollPane"), .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(300)}, .childGap = 6, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_2, .scroll = { .horizontal = true, .vertical = true } }) {
	               Clay_TextElementConfig *warningConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
	               CLAY({ .id = CLAY_ID("Clay__DebugViewWarningItemHeader"), .layout = { .sizing = {.Y = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childGap = 8, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
	                   CLAY_TEXT(CLAY_STRING("Warnings"), warningConfig);
	               }
	               CLAY({ .id = CLAY_ID("Clay__DebugViewWarningsTopBorder"), .layout = { .sizing = { .X = CLAY_SIZING_GROW(0), .Y = CLAY_SIZING_FIXED(1)} }, .backgroundColor = {200, 200, 200, 255} }) {}
	               int32 previousWarningsLength = c.warnings.length;
	               for (int32 i = 0; i < previousWarningsLength; i++) {
	                   Clay__Warning warning = c.warnings.internalArray[i];
	                   CLAY({ .id = CLAY_IDI("Clay__DebugViewWarningItem", i), .layout = { .sizing = {.Y = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childGap = 8, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
	                       CLAY_TEXT(warning.baseMessage, warningConfig);
	                       if (warning.dynamicMessage.length > 0) {
	                           CLAY_TEXT(warning.dynamicMessage, warningConfig);
	                       }
	                   }
	               }
	           }
	       }
	   }
	*/
}

var debugViewWidth uint32 = 400
var debugViewHighlightColor = colorex.RGBA{R: 168, G: 66, B: 28, A: 100}

/*
Clay__WarningArray Clay__WarningArray_Allocate_Arena(int32 capacity, Clay_Arena *arena) {
    size_t totalSizeBytes = capacity * sizeof(string);
    Clay__WarningArray array = {.capacity = capacity, .length = 0};
    uintptr_t nextAllocOffset = arena.nextAllocation + (64 - (arena.nextAllocation % 64));
    if (nextAllocOffset + totalSizeBytes <= arena.capacity) {
        array.internalArray = (Clay__Warning*)((uintptr_t)arena.memory + (uintptr_t)nextAllocOffset);
        arena.nextAllocation = nextAllocOffset + totalSizeBytes;
    }
    else {
        Clay__currentContext.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
            .errorType = ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
            .errorText = CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
            .userData = Clay__currentContext.errorHandler.userData });
    }
    return array;
}

Clay__Warning *Clay__WarningArray_Add(Clay__WarningArray *array, Clay__Warning item)
{
    if (array.length < array.capacity) {
        array.internalArray[array.length++] = item;
        return &array.internalArray[array.length - 1];
    }
    return &CLAY__WARNING_DEFAULT;
}

any Clay__Array_Allocate_Arena(int32 capacity, uint32 itemSize, Clay_Arena *arena)
{
    size_t totalSizeBytes = capacity * itemSize;
    uintptr_t nextAllocOffset = arena.nextAllocation + (64 - (arena.nextAllocation % 64));
    if (nextAllocOffset + totalSizeBytes <= arena.capacity) {
        arena.nextAllocation = nextAllocOffset + totalSizeBytes;
        return (any)((uintptr_t)arena.memory + (uintptr_t)nextAllocOffset);
    }
    else {
        Clay__currentContext.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
                .errorText = CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
                .userData = Clay__currentContext.errorHandler.userData });
    }
    return CLAY__NULL;
}

bool Clay__Array_RangeCheck(int32 index, int32 length)
{
    if (index < length && index >= 0) {
        return true;
    }
    context := GetCurrentContext();
    context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
            .errorType = ERROR_TYPE_INTERNAL_ERROR,
            .errorText = CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
            .userData = context.errorHandler.userData });
    return false;
}

bool Clay__Array_AddCapacityCheck(int32 length, int32 capacity)
{
    if (length < capacity) {
        return true;
    }
    context := GetCurrentContext();
    context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
        .errorType = ERROR_TYPE_INTERNAL_ERROR,
        .errorText = CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
        .userData = context.errorHandler.userData });
    return false;
}
*/
