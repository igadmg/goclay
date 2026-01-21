package clay

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
	renderTranslucent             bool

	// Layout Elements / Render Commands
	layoutElements              []LayoutElement
	renderCommands              []RenderCommand
	openLayoutElementStack      []int
	layoutElementChildren       []int
	layoutElementChildrenBuffer []int
	textElementData             []TextElementData
	aspectRatioElementIndexes   []int
	reusableElementIndexBuffer  []int32
	layoutElementClipElementIds []int

	// Configs
	layoutConfigs             []LayoutConfig
	elementConfigs            []AnyElementConfig
	textElementConfigs        []TextElementConfig
	aspectRatioElementConfigs []AspectRatioElementConfig
	imageElementConfigs       []ImageElementConfig
	floatingElementConfigs    []FloatingElementConfig
	clipElementConfigs        []ClipElementConfig
	customElementConfigs      []CustomElementConfig
	borderElementConfigs      []BorderElementConfig
	sharedElementConfigs      []SharedElementConfig

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

func (c *Context) Finalize() []RenderCommand {
	rc := c.renderCommands
	c.initializeEphemeralMemory()
	return rc
}

// Ephemeral Memory - reset every frame
func (c *Context) initializeEphemeralMemory() {
	clear(c.layoutElementChildrenBuffer)
	c.layoutElementChildrenBuffer = c.layoutElementChildrenBuffer[:0]
	clear(c.layoutElements)
	c.layoutElements = c.layoutElements[:0]
	clear(c.warnings)
	c.warnings = c.warnings[:0]

	clear(c.layoutConfigs)
	c.layoutConfigs = c.layoutConfigs[:0]
	clear(c.elementConfigs)
	c.elementConfigs = c.elementConfigs[:0]
	clear(c.textElementConfigs)
	c.textElementConfigs = c.textElementConfigs[:0]
	clear(c.aspectRatioElementConfigs)
	c.aspectRatioElementConfigs = c.aspectRatioElementConfigs[:0]
	clear(c.imageElementConfigs)
	c.imageElementConfigs = c.imageElementConfigs[:0]
	clear(c.floatingElementConfigs)
	c.floatingElementConfigs = c.floatingElementConfigs[:0]
	clear(c.clipElementConfigs)
	c.clipElementConfigs = c.clipElementConfigs[:0]
	clear(c.customElementConfigs)
	c.customElementConfigs = c.customElementConfigs[:0]
	clear(c.borderElementConfigs)
	c.borderElementConfigs = c.borderElementConfigs[:0]
	clear(c.sharedElementConfigs)
	c.sharedElementConfigs = c.sharedElementConfigs[:0]

	clear(c.layoutElementIdStrings)
	c.layoutElementIdStrings = c.layoutElementIdStrings[:0]
	clear(c.wrappedTextLines)
	c.wrappedTextLines = c.wrappedTextLines[:0]
	clear(c.layoutElementTreeNodeArray1)
	c.layoutElementTreeNodeArray1 = c.layoutElementTreeNodeArray1[:0]
	clear(c.layoutElementTreeRoots)
	c.layoutElementTreeRoots = c.layoutElementTreeRoots[:0]
	clear(c.layoutElementChildren)
	c.layoutElementChildren = c.layoutElementChildren[:0]
	clear(c.openLayoutElementStack)
	c.openLayoutElementStack = c.openLayoutElementStack[:0]
	clear(c.textElementData)
	c.textElementData = c.textElementData[:0]
	clear(c.aspectRatioElementIndexes)
	c.aspectRatioElementIndexes = c.aspectRatioElementIndexes[:0]
	clear(c.openClipElementStack)
	c.openClipElementStack = c.openClipElementStack[:0]
	clear(c.reusableElementIndexBuffer)
	c.reusableElementIndexBuffer = c.reusableElementIndexBuffer[:0]
	clear(c.layoutElementClipElementIds)
	c.layoutElementClipElementIds = c.layoutElementClipElementIds[:0]
	clear(c.dynamicStringData)
	c.dynamicStringData = c.dynamicStringData[:0]

	c.renderCommands = c.renderCommands[:0]
}

// Persistent memory - initialized once and not reset
func (c *Context) initializePersistentMemory() {
	maxElementCount := c.maxElementCount
	maxMeasureTextCacheWordCount := c.maxMeasureTextCacheWordCount

	c.scrollContainerDatas = make([]ScrollContainerDataInternal, 0, 100)
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
	c.clipElementConfigs = make([]ClipElementConfig, 0, maxElementCount)
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
	c.aspectRatioElementIndexes = make([]int, 0, maxElementCount)
	c.renderCommands = make([]RenderCommand, 0, maxElementCount)
	c.openClipElementStack = make([]int, 0, maxElementCount)
	c.reusableElementIndexBuffer = make([]int32, 0, maxElementCount)
	c.layoutElementClipElementIds = make([]int, 0, maxElementCount)
	c.dynamicStringData = make([]byte, 0, maxElementCount)
}

func (c *Context) getOpenLayoutElement() *LayoutElement {
	return &c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-1]]
}

func (c *Context) getParentElementId() uint32 {
	return c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-2]].id
}

func (c *Context) getHashMapItem(id uint32) (e *LayoutElementHashMapItem, ok bool) {
	e, ok = c.layoutElementsHashMap[id]
	return
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

func (c *Context) storeAspectRatioElementConfig(config AspectRatioElementConfig) *AspectRatioElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_AspectRatioElementConfig
	}
	c.aspectRatioElementConfigs = append(c.aspectRatioElementConfigs, config)
	return &c.aspectRatioElementConfigs[len(c.aspectRatioElementConfigs)-1]
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

func (c *Context) storeClipElementConfig(config ClipElementConfig) *ClipElementConfig {
	if c.booleanWarnings.maxElementsExceeded {
		return &default_ClipElementConfig
	}
	c.clipElementConfigs = append(c.clipElementConfigs, config)
	return &c.clipElementConfigs[len(c.clipElementConfigs)-1]
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

func (c *Context) addMeasuredWord(word MeasuredWord, previousWord *MeasuredWord) *MeasuredWord {
	if len(c.measuredWordsFreeList) > 0 {
		newItemIndex := c.measuredWordsFreeList[len(c.measuredWordsFreeList)-1]
		c.measuredWordsFreeList = c.measuredWordsFreeList[:len(c.measuredWordsFreeList)-1]
		c.measuredWords = slicesex_Set(c.measuredWords, int(newItemIndex), word)
		previousWord.next = newItemIndex
		return &c.measuredWords[newItemIndex]
	} else {
		previousWord.next = int32(len(c.measuredWords))
		c.measuredWords = append(c.measuredWords, word)
		return &c.measuredWords[len(c.measuredWords)-1]
	}
}

func (c *Context) addHashMapItem(elementId ElementId, layoutElement *LayoutElement) *LayoutElementHashMapItem {
	if len(c.layoutElementsHashMapInternal) == cap(c.layoutElementsHashMapInternal)-1 {
		return nil
	}

	item := LayoutElementHashMapItem{
		elementId:     elementId,
		layoutElement: layoutElement,
		nextIndex:     -1,
		generation:    c.generation + 1,
	}

	c.layoutElementsHashMapInternal = append(c.layoutElementsHashMapInternal, item)
	c.layoutElementsHashMap[elementId.id] = &c.layoutElementsHashMapInternal[len(c.layoutElementsHashMapInternal)-1]

	return c.layoutElementsHashMap[elementId.id]
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

func (c *Context) generateIdForAnonymousElement(openLayoutElement *LayoutElement) ElementId {
	parentElement := c.layoutElements[c.openLayoutElementStack[len(c.openLayoutElementStack)-2]]
	elementId := hashNumber(uint32(len(parentElement.children)), parentElement.id)
	openLayoutElement.id = elementId.id
	c.addHashMapItem(elementId, openLayoutElement)
	c.layoutElementIdStrings = append(c.layoutElementIdStrings, elementId.stringId)
	return elementId
}
