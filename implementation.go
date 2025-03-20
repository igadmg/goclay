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
	elementIndex        int32
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
	renderCommands              []any ///Clay_RenderCommandArray
	openLayoutElementStack      []int
	layoutElementChildren       []int
	layoutElementChildrenBuffer []int
	textElementData             []TextElementData
	imageElementPointers        []int32
	reusableElementIndexBuffer  []int32
	layoutElementClipElementIds []int32
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
	layoutElementIdStrings             []string
	wrappedTextLines                   []WrappedTextLine
	layoutElementTreeNodeArray         []LayoutElementTreeNode
	layoutElementTreeRoots             []LayoutElementTreeRoot
	layoutElementsHashMapInternal      []LayoutElementHashMapItem
	layoutElementsHashMap              map[uint32]*LayoutElementHashMapItem
	measureTextHashMapInternal         []MeasureTextCacheItem
	measureTextHashMapInternalFreeList []int32
	measureTextHashMap                 []int32
	measuredWords                      []MeasuredWord
	measuredWordsFreeList              []int32
	openClipElementStack               []int32
	pointerOverIds                     []ElementId
	scrollContainerDatas               []ScrollContainerDataInternal
	treeNodeVisited                    []bool
	dynamicStringData                  []byte
	debugElementData                   []DebugElementData
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

#ifdef CLAY_WASM
    __attribute__((import_module("clay"), import_name("measureTextFunction"))) vector2.Float32 Clay__MeasureText(Clay_StringSlice text, Clay_TextElementConfig *config, void *userData);
    __attribute__((import_module("clay"), import_name("queryScrollOffsetFunction"))) vector2.Float32 Clay__QueryScrollOffset(uint32 elementId, void *userData);
#else
    vector2.Float32 (*Clay__MeasureText)(Clay_StringSlice text, Clay_TextElementConfig *config, void *userData);
    vector2.Float32 (*Clay__QueryScrollOffset)(uint32 elementId, void *userData);
#endif
*/

func getOpenLayoutElement() *LayoutElement {
	context := GetCurrentContext()
	return &context.layoutElements[context.openLayoutElementStack[len(context.openLayoutElementStack)-1]]
}

func getParentElementId() uint32 {
	context := GetCurrentContext()
	return context.layoutElements[context.openLayoutElementStack[len(context.openLayoutElementStack)-2]].id
}

func storeLayoutConfig(config LayoutConfig) *LayoutConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_LayoutConfig
	}
	context.layoutConfigs = append(context.layoutConfigs, config)
	return &context.layoutConfigs[len(context.layoutConfigs)-1]
}

func storeTextElementConfig(config TextElementConfig) *TextElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_TextElementConfig
	}
	context.textElementConfigs = append(context.textElementConfigs, config)
	return &context.textElementConfigs[len(context.textElementConfigs)-1]
}

func storeImageElementConfig(config ImageElementConfig) *ImageElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_ImageElementConfig
	}
	context.imageElementConfigs = append(context.imageElementConfigs, config)
	return &context.imageElementConfigs[len(context.imageElementConfigs)-1]
}

func storeFloatingElementConfig(config FloatingElementConfig) *FloatingElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_FloatingElementConfig
	}
	context.floatingElementConfigs = append(context.floatingElementConfigs, config)
	return &context.floatingElementConfigs[len(context.floatingElementConfigs)-1]
}

func storeCustomElementConfig(config CustomElementConfig) *CustomElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_CustomElementConfig
	}
	context.customElementConfigs = append(context.customElementConfigs, config)
	return &context.customElementConfigs[len(context.customElementConfigs)-1]
}

func storeScrollElementConfig(config ScrollElementConfig) *ScrollElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_ScrollElementConfig
	}
	context.scrollElementConfigs = append(context.scrollElementConfigs, config)
	return &context.scrollElementConfigs[len(context.scrollElementConfigs)-1]
}

func storeBorderElementConfig(config BorderElementConfig) *BorderElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_BorderElementConfig
	}
	context.borderElementConfigs = append(context.borderElementConfigs, config)
	return &context.borderElementConfigs[len(context.borderElementConfigs)-1]
}

func storeSharedElementConfig(config SharedElementConfig) *SharedElementConfig {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return &default_SharedElementConfig
	}
	context.sharedElementConfigs = append(context.sharedElementConfigs, config)
	return &context.sharedElementConfigs[len(context.sharedElementConfigs)-1]
}

/*
	Clay_ElementConfig Clay__AttachElementConfig(Clay_ElementConfigUnion config, Clay__ElementConfigType type) {
	    context := GetCurrentContext();
	    if (context.booleanWarnings.maxElementsExceeded) {
	        return CLAY__INIT(Clay_ElementConfig) CLAY__DEFAULT_STRUCT;
	    }
	    openLayoutElement := getOpenLayoutElement();
	    openLayoutElement.elementConfigs.length++;
	    return *Clay__ElementConfigArray_Add(&context.elementConfigs, CLAY__INIT(Clay_ElementConfig) { .Type = type, .config = config });
	}

	Clay_ElementConfigUnion Clay__FindElementConfigWithType(Clay_LayoutElement *element, Clay__ElementConfigType type) {
	    for (int32 i = 0; i < element.elementConfigs.length; i++) {
	        Clay_ElementConfig *config = Clay__ElementConfigArraySlice_Get(&element.elementConfigs, i);
	        if (config.type == type) {
	            return config.config;
	        }
	    }
	    return CLAY__INIT(Clay_ElementConfigUnion) { NULL };
	}

	Clay_ElementId Clay__HashNumber(const uint32 offset, const uint32 seed) {
	    uint32 hash = seed;
	    hash += (offset + 48);
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += (hash << 3);
	    hash ^= (hash >> 11);
	    hash += (hash << 15);
	    return CLAY__INIT(Clay_ElementId) { .id = hash + 1, .offset = offset, .baseId = seed, .stringId = STRING_DEFAULT }; // Reserve the hash result of zero as "null id"
	}
*/

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

/*
	uint32 Clay__HashTextWithConfig(string *text, Clay_TextElementConfig *config) {
	    uint32 hash = 0;
	    uintptr_t pointerAsNumber = (uintptr_t)text.chars;

	    if (config.hashStringContents) {
	        uint32 maxLengthToHash = min(text.length, 256);
	        for (uint32 i = 0; i < maxLengthToHash; i++) {
	            hash += text.chars[i];
	            hash += (hash << 10);
	            hash ^= (hash >> 6);
	        }
	    } else {
	        hash += pointerAsNumber;
	        hash += (hash << 10);
	        hash ^= (hash >> 6);
	    }

	    hash += text.length;
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += config.fontId;
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += config.fontSize;
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += config.lineHeight;
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += config.letterSpacing;
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += config.wrapMode;
	    hash += (hash << 10);
	    hash ^= (hash >> 6);

	    hash += (hash << 3);
	    hash ^= (hash >> 11);
	    hash += (hash << 15);
	    return hash + 1; // Reserve the hash result of zero as "null id"
	}

	Clay__MeasuredWord *Clay__AddMeasuredWord(Clay__MeasuredWord word, Clay__MeasuredWord *previousWord) {
	    context := GetCurrentContext();
	    if (context.measuredWordsFreeList.length > 0) {
	        uint32 newItemIndex = Clay__int32_tArray_GetValue(&context.measuredWordsFreeList, (int)context.measuredWordsFreeList.length - 1);
	        context.measuredWordsFreeList.length--;
	        Clay__MeasuredWordArray_Set(&context.measuredWords, (int)newItemIndex, word);
	        previousWord.next = (int32)newItemIndex;
	        return Clay__MeasuredWordArray_Get(&context.measuredWords, (int)newItemIndex);
	    } else {
	        previousWord.next = (int32)context.measuredWords.length;
	        return Clay__MeasuredWordArray_Add(&context.measuredWords, word);
	    }
	}

	Clay__MeasureTextCacheItem *Clay__MeasureTextCached(string *text, Clay_TextElementConfig *config) {
	    context := GetCurrentContext();
	    #ifndef CLAY_WASM
	    if (!Clay__MeasureText) {
	        if (!context.booleanWarnings.textMeasurementFunctionNotSet) {
	            context.booleanWarnings.textMeasurementFunctionNotSet = true;
	            context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	                    .errorType = CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED,
	                    .errorText = CLAY_STRING("Clay's internal MeasureText function is null. You may have forgotten to call Clay_SetMeasureTextFunction(), or passed a NULL function pointer by mistake."),
	                    .userData = context.errorHandler.userData });
	        }
	        return &Clay__MeasureTextCacheItem_DEFAULT;
	    }
	    #endif
	    uint32 id = Clay__HashTextWithConfig(text, config);
	    uint32 hashBucket = id % (context.maxMeasureTextCacheWordCount / 32);
	    int32 elementIndexPrevious = 0;
	    int32 elementIndex = context.measureTextHashMap.internalArray[hashBucket];
	    while (elementIndex != 0) {
	        Clay__MeasureTextCacheItem *hashEntry = Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, elementIndex);
	        if (hashEntry.id == id) {
	            hashEntry.generation = context.generation;
	            return hashEntry;
	        }
	        // This element hasn't been seen in a few frames, delete the hash map item
	        if (context.generation - hashEntry.generation > 2) {
	            // Add all the measured words that were included in this measurement to the freelist
	            int32 nextWordIndex = hashEntry.measuredWordsStartIndex;
	            while (nextWordIndex != -1) {
	                Clay__MeasuredWord *measuredWord = Clay__MeasuredWordArray_Get(&context.measuredWords, nextWordIndex);
	                Clay__int32_tArray_Add(&context.measuredWordsFreeList, nextWordIndex);
	                nextWordIndex = measuredWord.next;
	            }

	            int32 nextIndex = hashEntry.nextIndex;
	            Clay__MeasureTextCacheItemArray_Set(&context.measureTextHashMapInternal, elementIndex, CLAY__INIT(Clay__MeasureTextCacheItem) { .measuredWordsStartIndex = -1 });
	            Clay__int32_tArray_Add(&context.measureTextHashMapInternalFreeList, elementIndex);
	            if (elementIndexPrevious == 0) {
	                context.measureTextHashMap.internalArray[hashBucket] = nextIndex;
	            } else {
	                Clay__MeasureTextCacheItem *previousHashEntry = Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, elementIndexPrevious);
	                previousHashEntry.nextIndex = nextIndex;
	            }
	            elementIndex = nextIndex;
	        } else {
	            elementIndexPrevious = elementIndex;
	            elementIndex = hashEntry.nextIndex;
	        }
	    }

	    int32 newItemIndex = 0;
	    Clay__MeasureTextCacheItem newCacheItem = { .measuredWordsStartIndex = -1, .id = id, .generation = context.generation };
	    Clay__MeasureTextCacheItem *measured = NULL;
	    if (context.measureTextHashMapInternalFreeList.length > 0) {
	        newItemIndex = Clay__int32_tArray_GetValue(&context.measureTextHashMapInternalFreeList, context.measureTextHashMapInternalFreeList.length - 1);
	        context.measureTextHashMapInternalFreeList.length--;
	        Clay__MeasureTextCacheItemArray_Set(&context.measureTextHashMapInternal, newItemIndex, newCacheItem);
	        measured = Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, newItemIndex);
	    } else {
	        if (context.measureTextHashMapInternal.length == context.measureTextHashMapInternal.capacity - 1) {
	            if (!context.booleanWarnings.maxTextMeasureCacheExceeded) {
	                context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	                        .errorType = CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
	                        .errorText = CLAY_STRING("Clay ran out of capacity while attempting to measure text elements. Try using Clay_SetMaxElementCount() with a higher value."),
	                        .userData = context.errorHandler.userData });
	                context.booleanWarnings.maxTextMeasureCacheExceeded = true;
	            }
	            return &Clay__MeasureTextCacheItem_DEFAULT;
	        }
	        measured = Clay__MeasureTextCacheItemArray_Add(&context.measureTextHashMapInternal, newCacheItem);
	        newItemIndex = context.measureTextHashMapInternal.length - 1;
	    }

	    int32 start = 0;
	    int32 end = 0;
	    float lineWidth = 0;
	    float measuredWidth = 0;
	    float measuredHeight = 0;
	    float spaceWidth = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = 1, .chars = SPACECHAR.chars, .baseChars = SPACECHAR.chars }, config, context.measureTextUserData).X;
	    Clay__MeasuredWord tempWord = { .next = -1 };
	    Clay__MeasuredWord *previousWord = &tempWord;
	    while (end < text.length) {
	        if (context.measuredWords.length == context.measuredWords.capacity - 1) {
	            if (!context.booleanWarnings.maxTextMeasureCacheExceeded) {
	                context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	                    .errorType = CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED,
	                    .errorText = CLAY_STRING("Clay has run out of space in it's internal text measurement cache. Try using Clay_SetMaxMeasureTextCacheWordCount() (default 16384, with 1 unit storing 1 measured word)."),
	                    .userData = context.errorHandler.userData });
	                context.booleanWarnings.maxTextMeasureCacheExceeded = true;
	            }
	            return &Clay__MeasureTextCacheItem_DEFAULT;
	        }
	        char current = text.chars[end];
	        if (current == ' ' || current == '\n') {
	            int32 length = end - start;
	            vector2.Float32 dimensions = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = length, .chars = &text.chars[start], .baseChars = text.chars }, config, context.measureTextUserData);
	            measuredHeight = max(measuredHeight, dimensions.Y);
	            if (current == ' ') {
	                dimensions.X += spaceWidth;
	                previousWord = Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = start, .length = length + 1, .X = dimensions.X, .next = -1 }, previousWord);
	                lineWidth += dimensions.X;
	            }
	            if (current == '\n') {
	                if (length > 0) {
	                    previousWord = Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = start, .length = length, .X = dimensions.X, .next = -1 }, previousWord);
	                }
	                previousWord = Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = end + 1, .length = 0, .X = 0, .next = -1 }, previousWord);
	                lineWidth += dimensions.X;
	                measuredWidth = max(lineWidth, measuredWidth);
	                measured.containsNewlines = true;
	                lineWidth = 0;
	            }
	            start = end + 1;
	        }
	        end++;
	    }
	    if (end - start > 0) {
	        vector2.Float32 dimensions = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = end - start, .chars = &text.chars[start], .baseChars = text.chars }, config, context.measureTextUserData);
	        Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = start, .length = end - start, .X = dimensions.X, .next = -1 }, previousWord);
	        lineWidth += dimensions.X;
	        measuredHeight = max(measuredHeight, dimensions.Y);
	    }
	    measuredWidth = max(lineWidth, measuredWidth);

	    measured.measuredWordsStartIndex = tempWord.next;
	    measured.unwrappedDimensions.X = measuredWidth;
	    measured.unwrappedDimensions.Y = measuredHeight;

	    if (elementIndexPrevious != 0) {
	        Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, elementIndexPrevious).nextIndex = newItemIndex;
	    } else {
	        context.measureTextHashMap.internalArray[hashBucket] = newItemIndex;
	    }
	    return measured;
	}
*/

func addHashMapItem(elementId ElementId, layoutElement *LayoutElement, idAlias uint32) *LayoutElementHashMapItem {
	context := GetCurrentContext()
	if len(context.layoutElementsHashMapInternal) == cap(context.layoutElementsHashMapInternal)-1 {
		return nil
	}

	item := LayoutElementHashMapItem{
		elementId:     elementId,
		layoutElement: layoutElement,
		nextIndex:     -1,
		generation:    context.generation + 1,
		idAlias:       idAlias,
	}

	context.layoutElementsHashMapInternal = append(context.layoutElementsHashMapInternal, item)
	context.layoutElementsHashMap[elementId.id] = &context.layoutElementsHashMapInternal[len(context.layoutElementsHashMapInternal)-1]

	return context.layoutElementsHashMap[elementId.id]
}

func Clay__GetHashMapItem(id uint32) *LayoutElementHashMapItem {
	context := GetCurrentContext()
	r, ok := context.layoutElementsHashMap[id]
	if !ok {
		return &default_LayoutElementHashMapItem
	}

	return r
}

/*
	Clay_ElementId Clay__GenerateIdForAnonymousElement(Clay_LayoutElement *openLayoutElement) {
	    context := GetCurrentContext();
	    Clay_LayoutElement *parentElement = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length - 2));
	    Clay_ElementId elementId = Clay__HashNumber(parentElement.childrenOrTextContent.children.length, parentElement.id);
	    openLayoutElement.id = elementId.id;
	    addHashMapItem(elementId, openLayoutElement, 0);
	    Clay__StringArray_Add(&context.layoutElementIdStrings, elementId.stringId);
	    return elementId;
	}
*/

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

func closeElement() {
	context := GetCurrentContext()
	// TODO: implement
	//if (context.booleanWarnings.maxElementsExceeded) {
	//   return;
	//}

	openLayoutElement := getOpenLayoutElement()
	layoutConfig := openLayoutElement.layoutConfig
	elementHasScrollHorizontal := false
	elementHasScrollVertical := false

	for _, config := range openLayoutElement.elementConfigs {
		switch c := config.(type) {
		case *ScrollElementConfig:
			elementHasScrollHorizontal = c.horizontal
			elementHasScrollVertical = c.vertical
			context.openClipElementStack = context.openClipElementStack[:len(context.openClipElementStack)-1]
			break
		case *FloatingElementConfig:
			context.openClipElementStack = context.openClipElementStack[:len(context.openClipElementStack)-1]
		}
	}

	// Attach children to the current open element // TODO: have no idea
	openLayoutElement.children = context.layoutElementChildren[len(context.layoutElementChildren):len(context.layoutElementChildren)]
	if layoutConfig.LayoutDirection == LEFT_TO_RIGHT {
		openLayoutElement.dimensions.X = (float32)(layoutConfig.Padding.Left + layoutConfig.Padding.Right)
		for i := range openLayoutElement.children {
			childIndex := context.layoutElementChildrenBuffer[len(context.layoutElementChildrenBuffer)-(int)(len(openLayoutElement.children)+i)]
			child := context.layoutElements[childIndex]
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
			context.layoutElementChildren = append(context.layoutElementChildren, childIndex)
		}

		childGap := (float32)(max(len(openLayoutElement.children)-1, 0) * int(layoutConfig.ChildGap))
		openLayoutElement.dimensions.X += childGap // TODO this is technically a bug with childgap and scroll containers
		openLayoutElement.minDimensions.X += childGap
	} else if layoutConfig.LayoutDirection == TOP_TO_BOTTOM {
		openLayoutElement.dimensions.Y = (float32)(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)
		for i := range openLayoutElement.children {
			childIndex := context.layoutElementChildrenBuffer[len(context.layoutElementChildrenBuffer)-len(openLayoutElement.children)+i]
			child := context.layoutElements[childIndex]
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
			context.layoutElementChildren = append(context.layoutElementChildren, childIndex)
		}
		childGap := (float32)(max(len(openLayoutElement.children)-1, 0) * int(layoutConfig.ChildGap))
		openLayoutElement.dimensions.Y += childGap // TODO this is technically a bug with childgap and scroll containers
		openLayoutElement.minDimensions.Y += childGap
	}

	context.layoutElementChildrenBuffer = context.layoutElementChildrenBuffer[:len(context.layoutElementChildrenBuffer)-len(openLayoutElement.children)]

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
	context.openLayoutElementStack, closingElementIndex = slicesex.RemoveSwapback(context.openLayoutElementStack, len(context.openLayoutElementStack)-1)
	openLayoutElement = getOpenLayoutElement()

	if !elementIsFloating && len(context.openLayoutElementStack) > 1 {
		openLayoutElement.children = append(openLayoutElement.children, closingElementIndex)
		context.layoutElementChildrenBuffer = append(context.layoutElementChildrenBuffer, closingElementIndex)
	}
}

/*
bool Clay__MemCmp(const char *s1, const char *s2, int32 length);
#if !defined(CLAY_DISABLE_SIMD) && (defined(__x86_64__) || defined(_M_X64) || defined(_M_AMD64))
    bool Clay__MemCmp(const char *s1, const char *s2, int32 length) {
        while (length >= 16) {
            __m128i v1 = _mm_loadu_si128((const __m128i *)s1);
            __m128i v2 = _mm_loadu_si128((const __m128i *)s2);

            if (_mm_movemask_epi8(_mm_cmpeq_epi8(v1, v2)) != 0xFFFF) { // If any byte differs
                return false;
            }

            s1 += 16;
            s2 += 16;
            length -= 16;
        }

        // Handle remaining bytes
        while (length--) {
            if (*s1 != *s2) {
                return false;
            }
            s1++;
            s2++;
        }

        return true;
    }
#elif !defined(CLAY_DISABLE_SIMD) && defined(__aarch64__)
    bool Clay__MemCmp(const char *s1, const char *s2, int32 length) {
        while (length >= 16) {
            uint8x16_t v1 = vld1q_u8((const uint8_t *)s1);
            uint8x16_t v2 = vld1q_u8((const uint8_t *)s2);

            // Compare vectors
            if (vminvq_u32(vreinterpretq_u32_u8(vceqq_u8(v1, v2))) != 0xFFFFFFFF) { // If there's a difference
                return false;
            }

            s1 += 16;
            s2 += 16;
            length -= 16;
        }

        // Handle remaining bytes
        while (length--) {
            if (*s1 != *s2) {
                return false;
            }
            s1++;
            s2++;
        }

        return true;
    }
#else
    bool Clay__MemCmp(const char *s1, const char *s2, int32 length) {
        for (int32 i = 0; i < length; i++) {
            if (s1[i] != s2[i]) {
                return false;
            }
        }
        return true;
    }
#endif
*/

func openElement() {
	context := GetCurrentContext()
	// TODO: implement
	//if (context.layoutElements.length == context.layoutElements.capacity - 1 || context.booleanWarnings.maxElementsExceeded) {
	//    context.booleanWarnings.maxElementsExceeded = true;
	//    return;
	//}

	layoutElement := LayoutElement{}
	context.layoutElements = append(context.layoutElements, layoutElement)
	context.openLayoutElementStack = append(context.openLayoutElementStack, len(context.layoutElements)-1)
	if len(context.openClipElementStack) > 0 {
		context.layoutElementClipElementIds[len(context.layoutElements)-1] = context.openClipElementStack[len(context.openClipElementStack)-1]
	} else {
		context.layoutElementClipElementIds[len(context.layoutElements)-1] = 0
	}
}

/*
	void Clay__OpenTextElement(string text, Clay_TextElementConfig *textConfig) {
	    context := GetCurrentContext();
	    if (context.layoutElements.length == context.layoutElements.capacity - 1 || context.booleanWarnings.maxElementsExceeded) {
	        context.booleanWarnings.maxElementsExceeded = true;
	        return;
	    }
	    Clay_LayoutElement *parentElement = getOpenLayoutElement();

	    Clay_LayoutElement layoutElement = CLAY__DEFAULT_STRUCT;
	    Clay_LayoutElement *textElement = Clay_LayoutElementArray_Add(&context.layoutElements, layoutElement);
	    if (context.openClipElementStack.length > 0) {
	        Clay__int32_tArray_Set(&context.layoutElementClipElementIds, context.layoutElements.length - 1, Clay__int32_tArray_GetValue(&context.openClipElementStack, (int)context.openClipElementStack.length - 1));
	    } else {
	        Clay__int32_tArray_Set(&context.layoutElementClipElementIds, context.layoutElements.length - 1, 0);
	    }

	    Clay__int32_tArray_Add(&context.layoutElementChildrenBuffer, context.layoutElements.length - 1);
	    Clay__MeasureTextCacheItem *textMeasured = Clay__MeasureTextCached(&text, textConfig);
	    Clay_ElementId elementId = Clay__HashNumber(parentElement.childrenOrTextContent.children.length, parentElement.id);
	    textElement.id = elementId.id;
	    addHashMapItem(elementId, textElement, 0);
	    Clay__StringArray_Add(&context.layoutElementIdStrings, elementId.stringId);
	    vector2.Float32 textDimensions = { .X = textMeasured.unwrappedDimensions.X, .Y = textConfig.lineHeight > 0 ? (float32)textConfig.lineHeight : textMeasured.unwrappedDimensions.Y };
	    textElement.dimensions = textDimensions;
	    textElement.minDimensions = CLAY__INIT(vector2.Float32) { .X = textMeasured.unwrappedDimensions.Y, .Y = textDimensions.Y }; // TODO not sure this is the best way to decide min width for text
	    textElement.childrenOrTextContent.textElementData = Clay__TextElementDataArray_Add(&context.textElementData, CLAY__INIT(Clay__TextElementData) { .text = text, .preferredDimensions = textMeasured.unwrappedDimensions, .elementIndex = context.layoutElements.length - 1 });
	    textElement.elementConfigs = CLAY__INIT(Clay__ElementConfigArraySlice) {
	            .length = 1,
	            .internalArray = Clay__ElementConfigArray_Add(&context.elementConfigs, CLAY__INIT(Clay_ElementConfig) { .Type = CLAY__ELEMENT_CONFIG_TYPE_TEXT, .config = { .textElementConfig = textConfig }})
	    };
	    textElement.layoutConfig = &CLAY_LAYOUT_DEFAULT;
	    parentElement.childrenOrTextContent.children.length++;
	}
*/

func attachId(elementId ElementId) ElementId {
	context := GetCurrentContext()
	if context.booleanWarnings.maxElementsExceeded {
		return default_ElementId
	}
	openLayoutElement := getOpenLayoutElement()
	idAlias := openLayoutElement.id
	openLayoutElement.id = elementId.id
	addHashMapItem(elementId, openLayoutElement, idAlias)
	context.layoutElementIdStrings = append(context.layoutElementIdStrings, elementId.stringId)
	return elementId
}

func configureOpenElement(declaration *ElementDeclaration) {
	//context := GetCurrentContext()
	openLayoutElement := getOpenLayoutElement()
	openLayoutElement.layoutConfig = storeLayoutConfig(declaration.Layout)

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
		//	        context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
		//	                .errorType = CLAY_ERROR_TYPE_PERCENTAGE_OVER_1,
		//	                .errorText = CLAY_STRING("An element was configured with CLAY_SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2."),
		//	                .userData = context.errorHandler.userData });
	}

	//openLayoutElementId := declaration.Id

	//openLayoutElement.elementConfigs.internalArray = &context.elementConfigs.internalArray[context.elementConfigs.length]
	//sharedConfig := (*SharedElementConfig)(nil)
	//if declaration.backgroundColor.A > 0 {
	//	sharedConfig = storeSharedElementConfig(SharedElementConfig{backgroundColor: declaration.backgroundColor})
	//	//Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
	//}
	/*
	   if (!Clay__MemCmp((char *)(&declaration.cornerRadius), (char *)(&Clay__CornerRadius_DEFAULT), sizeof(Clay_CornerRadius))) {
	       if (sharedConfig) {
	           sharedConfig.cornerRadius = declaration.cornerRadius;
	       } else {
	           sharedConfig = storeSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .cornerRadius = declaration.cornerRadius });
	           Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
	       }
	   }
	   if (declaration.userData != 0) {
	       if (sharedConfig) {
	           sharedConfig.userData = declaration.userData;
	       } else {
	           sharedConfig = storeSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .userData = declaration.userData });
	           Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
	       }
	   }
	   if (declaration.image.imageData) {
	       Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .imageElementConfig = storeImageElementConfig(declaration.image) }, CLAY__ELEMENT_CONFIG_TYPE_IMAGE);
	       Clay__int32_tArray_Add(&context.imageElementPointers, context.layoutElements.length - 1);
	   }
	   if (declaration.floating.attachTo != CLAY_ATTACH_TO_NONE) {
	       Clay_FloatingElementConfig floatingConfig = declaration.floating;
	       // This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
	       Clay_LayoutElement *hierarchicalParent = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length - 2));
	       if (hierarchicalParent) {
	           uint32 clipElementId = 0;
	           if (declaration.floating.attachTo == CLAY_ATTACH_TO_PARENT) {
	               // Attach to the element's direct hierarchical parent
	               floatingConfig.parentId = hierarchicalParent.id;
	               if (context.openClipElementStack.length > 0) {
	                   clipElementId = Clay__int32_tArray_GetValue(&context.openClipElementStack, (int)context.openClipElementStack.length - 1);
	               }
	           } else if (declaration.floating.attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID) {
	               Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingConfig.parentId);
	               if (!parentItem) {
	                   context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	                           .errorType = CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
	                           .errorText = CLAY_STRING("A floating element was declared with a parentId, but no element with that ID was found."),
	                           .userData = context.errorHandler.userData });
	               } else {
	                   clipElementId = Clay__int32_tArray_GetValue(&context.layoutElementClipElementIds, (int32)(parentItem.layoutElement - context.layoutElements.internalArray));
	               }
	           } else if (declaration.floating.attachTo == CLAY_ATTACH_TO_ROOT) {
	               floatingConfig.parentId = hashString(CLAY_STRING("Clay__RootContainer"), 0, 0).id;
	           }
	           if (!openLayoutElementId.id) {
	               openLayoutElementId = hashString(CLAY_STRING("Clay__FloatingContainer"), context.layoutElementTreeRoots.length, 0);
	           }
	           int32 currentElementIndex = Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length - 1);
	           Clay__int32_tArray_Set(&context.layoutElementClipElementIds, currentElementIndex, clipElementId);
	           Clay__int32_tArray_Add(&context.openClipElementStack, clipElementId);
	           Clay__LayoutElementTreeRootArray_Add(&context.layoutElementTreeRoots, CLAY__INIT(Clay__LayoutElementTreeRoot) {
	                   .layoutElementIndex = Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length - 1),
	                   .parentId = floatingConfig.parentId,
	                   .clipElementId = clipElementId,
	                   .zIndex = floatingConfig.zIndex,
	           });
	           Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .floatingElementConfig = storeFloatingElementConfig(floatingConfig) }, CLAY__ELEMENT_CONFIG_TYPE_FLOATING);
	       }
	   }
	   if (declaration.custom.customData) {
	       Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .customElementConfig = storeCustomElementConfig(declaration.custom) }, CLAY__ELEMENT_CONFIG_TYPE_CUSTOM);
	   }

	   if (openLayoutElementId.id != 0) {
	       attachId(openLayoutElementId);
	   } else if (openLayoutElement.id == 0) {
	       openLayoutElementId = Clay__GenerateIdForAnonymousElement(openLayoutElement);
	   }

	   if (declaration.scroll.horizontal | declaration.scroll.vertical) {
	       Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .scrollElementConfig = storeScrollElementConfig(declaration.scroll) }, CLAY__ELEMENT_CONFIG_TYPE_SCROLL);
	       Clay__int32_tArray_Add(&context.openClipElementStack, (int)openLayoutElement.id);
	       // Retrieve or create cached data to track scroll position across frames
	       Clay__ScrollContainerDataInternal *scrollOffset = CLAY__NULL;
	       for (int32 i = 0; i < context.scrollContainerDatas.length; i++) {
	           Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i);
	           if (openLayoutElement.id == mapping.elementId) {
	               scrollOffset = mapping;
	               scrollOffset.layoutElement = openLayoutElement;
	               scrollOffset.openThisFrame = true;
	           }
	       }
	       if (!scrollOffset) {
	           scrollOffset = Clay__ScrollContainerDataInternalArray_Add(&context.scrollContainerDatas, CLAY__INIT(Clay__ScrollContainerDataInternal){.layoutElement = openLayoutElement, .scrollOrigin = {-1,-1}, .elementId = openLayoutElement.id, .openThisFrame = true});
	       }
	       if (context.externalScrollHandlingEnabled) {
	           scrollOffset.scrollPosition = Clay__QueryScrollOffset(scrollOffset.elementId, context.queryScrollOffsetUserData);
	       }
	   }
	   if (!Clay__MemCmp((char *)(&declaration.border.X), (char *)(&Clay__BorderWidth_DEFAULT), sizeof(Clay_BorderWidth))) {
	       Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .borderElementConfig = storeBorderElementConfig(declaration.border) }, CLAY__ELEMENT_CONFIG_TYPE_BORDER);
	   }
	*/
}

// Ephemeral Memory - reset every frame
func initializeEphemeralMemory(context *Context) {
	context.layoutElementChildrenBuffer = context.layoutElementChildrenBuffer[:0]
	context.layoutElements = context.layoutElements[:0]
	context.warnings = context.warnings[:0]

	context.layoutConfigs = context.layoutConfigs[:0]
	context.elementConfigs = context.elementConfigs[:0]
	context.textElementConfigs = context.textElementConfigs[:0]
	context.imageElementConfigs = context.imageElementConfigs[:0]
	context.floatingElementConfigs = context.floatingElementConfigs[:0]
	context.scrollElementConfigs = context.scrollElementConfigs[:0]
	context.customElementConfigs = context.customElementConfigs[:0]
	context.borderElementConfigs = context.borderElementConfigs[:0]
	context.sharedElementConfigs = context.sharedElementConfigs[:0]

	context.layoutElementIdStrings = context.layoutElementIdStrings[:0]
	context.wrappedTextLines = context.wrappedTextLines[:0]
	context.layoutElementTreeNodeArray = context.layoutElementTreeNodeArray[:0]
	context.layoutElementTreeRoots = context.layoutElementTreeRoots[:0]
	context.layoutElementChildren = context.layoutElementChildren[:0]
	context.openLayoutElementStack = context.openLayoutElementStack[:0]
	context.textElementData = context.textElementData[:0]
	context.imageElementPointers = context.imageElementPointers[:0]
	context.renderCommands = context.renderCommands[:0]
	context.treeNodeVisited = context.treeNodeVisited[:0]
	context.openClipElementStack = context.openClipElementStack[:0]
	context.reusableElementIndexBuffer = context.reusableElementIndexBuffer[:0]
	context.layoutElementClipElementIds = context.layoutElementClipElementIds[:0]
	context.dynamicStringData = context.dynamicStringData[:0]
}

// Persistent memory - initialized once and not reset
func initializePersistentMemory(context *Context) {
	maxElementCount := context.maxElementCount
	maxMeasureTextCacheWordCount := context.maxMeasureTextCacheWordCount

	context.scrollContainerDatas = make([]ScrollContainerDataInternal, 0, 10)
	context.layoutElementsHashMapInternal = make([]LayoutElementHashMapItem, 0, maxElementCount)
	context.layoutElementsHashMap = map[uint32]*LayoutElementHashMapItem{}
	context.measureTextHashMapInternal = make([]MeasureTextCacheItem, 0, maxElementCount)
	context.measureTextHashMapInternalFreeList = make([]int32, 0, maxElementCount)
	context.measuredWordsFreeList = make([]int32, 0, maxMeasureTextCacheWordCount)
	context.measureTextHashMap = make([]int32, 0, maxElementCount)
	context.measuredWords = make([]MeasuredWord, 0, maxMeasureTextCacheWordCount)
	context.pointerOverIds = make([]ElementId, 0, maxElementCount)
	context.debugElementData = make([]DebugElementData, 0, maxElementCount)

	context.layoutElementChildrenBuffer = make([]int, 0, maxElementCount)
	context.layoutElements = make([]LayoutElement, 0, maxElementCount)
	context.warnings = make([]Warning, 0, 100)

	context.layoutConfigs = make([]LayoutConfig, 0, maxElementCount)
	context.elementConfigs = make([]AnyElementConfig, 0, maxElementCount)
	context.textElementConfigs = make([]TextElementConfig, 0, maxElementCount)
	context.imageElementConfigs = make([]ImageElementConfig, 0, maxElementCount)
	context.floatingElementConfigs = make([]FloatingElementConfig, 0, maxElementCount)
	context.scrollElementConfigs = make([]ScrollElementConfig, 0, maxElementCount)
	context.customElementConfigs = make([]CustomElementConfig, 0, maxElementCount)
	context.borderElementConfigs = make([]BorderElementConfig, 0, maxElementCount)
	context.sharedElementConfigs = make([]SharedElementConfig, 0, maxElementCount)

	context.layoutElementIdStrings = make([]string, 0, maxElementCount)
	context.wrappedTextLines = make([]WrappedTextLine, 0, maxElementCount)
	context.layoutElementTreeNodeArray = make([]LayoutElementTreeNode, 0, maxElementCount)
	context.layoutElementTreeRoots = make([]LayoutElementTreeRoot, 0, maxElementCount)
	context.layoutElementChildren = make([]int, 0, maxElementCount)
	context.openLayoutElementStack = make([]int, 0, maxElementCount)
	context.textElementData = make([]TextElementData, 0, maxElementCount)
	context.imageElementPointers = make([]int32, 0, maxElementCount)
	context.renderCommands = make([]any, 0, maxElementCount)
	context.treeNodeVisited = make([]bool, maxElementCount)
	context.openClipElementStack = make([]int32, 0, maxElementCount)
	context.reusableElementIndexBuffer = make([]int32, 0, maxElementCount)
	context.layoutElementClipElementIds = make([]int32, 0, maxElementCount)
	context.dynamicStringData = make([]byte, 0, maxElementCount)
}

/*
const float CLAY__EPSILON = 0.01;

	bool Clay__FloatEqual(float left, float right) {
	    float subtracted = left - right;
	    return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON;
	}

	void Clay__SizeContainersAlongAxis(bool xAxis) {
	    context := GetCurrentContext();
	    Clay__int32_tArray bfsBuffer = context.layoutElementChildrenBuffer;
	    Clay__int32_tArray resizableContainerBuffer = context.openLayoutElementStack;
	    for (int32 rootIndex = 0; rootIndex < context.layoutElementTreeRoots.length; ++rootIndex) {
	        bfsBuffer.length = 0;
	        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex);
	        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context.layoutElements, (int)root.layoutElementIndex);
	        Clay__int32_tArray_Add(&bfsBuffer, (int32)root.layoutElementIndex);

	        // Size floating containers to their parents
	        if (elementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
	            Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
	            Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingElementConfig.parentId);
	            if (parentItem && parentItem != &Clay_LayoutElementHashMapItem_DEFAULT) {
	                Clay_LayoutElement *parentLayoutElement = parentItem.layoutElement;
	                if (rootElement.layoutConfig.sizing.Width.Type == CLAY__SIZING_TYPE_GROW) {
	                    rootElement.dimensions.X = parentLayoutElement.dimensions.X;
	                }
	                if (rootElement.layoutConfig.sizing.Height.Type == CLAY__SIZING_TYPE_GROW) {
	                    rootElement.dimensions.Y = parentLayoutElement.dimensions.Y;
	                }
	            }
	        }

	        rootElement.dimensions.X = min(max(rootElement.dimensions.X, rootElement.layoutConfig.sizing.Width.size.MinMax.Min), rootElement.layoutConfig.sizing.Width.size.MinMax.Max);
	        rootElement.dimensions.Y = min(max(rootElement.dimensions.Y, rootElement.layoutConfig.sizing.Height.size.MinMax.Min), rootElement.layoutConfig.sizing.Height.size.MinMax.Max);

	        for (int32 i = 0; i < bfsBuffer.length; ++i) {
	            int32 parentIndex = Clay__int32_tArray_GetValue(&bfsBuffer, i);
	            Clay_LayoutElement *parent = Clay_LayoutElementArray_Get(&context.layoutElements, parentIndex);
	            Clay_LayoutConfig *parentStyleConfig = parent.layoutConfig;
	            int32 growContainerCount = 0;
	            float parentSize = xAxis ? parent.dimensions.X : parent.dimensions.Y;
	            float parentPadding = (float32)(xAxis ? (parent.layoutConfig.padding.Left + parent.layoutConfig.padding.Right) : (parent.layoutConfig.padding.Top + parent.layoutConfig.padding.Bottom));
	            float innerContentSize = 0, totalPaddingAndChildGaps = parentPadding;
	            bool sizingAlongAxis = (xAxis && parentStyleConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) || (!xAxis && parentStyleConfig.layoutDirection == CLAY_TOP_TO_BOTTOM);
	            resizableContainerBuffer.length = 0;
	            float parentChildGap = parentStyleConfig.childGap;

	            for (int32 childOffset = 0; childOffset < parent.childrenOrTextContent.children.length; childOffset++) {
	                int32 childElementIndex = parent.childrenOrTextContent.children.elements[childOffset];
	                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, childElementIndex);
	                Clay_SizingAxis childSizing = xAxis ? childElement.layoutConfig.sizing.Width : childElement.layoutConfig.sizing.Height;
	                float childSize = xAxis ? childElement.dimensions.X : childElement.dimensions.Y;

	                if (!elementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && childElement.childrenOrTextContent.children.length > 0) {
	                    Clay__int32_tArray_Add(&bfsBuffer, childElementIndex);
	                }

	                if (childSizing.Type != CLAY__SIZING_TYPE_PERCENT
	                    && childSizing.Type != CLAY__SIZING_TYPE_FIXED
	                    && (!elementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (Clay__FindElementConfigWithType(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig.wrapMode == CLAY_TEXT_WRAP_WORDS)) // todo too many loops
	                    && (xAxis || !elementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_IMAGE))
	                ) {
	                    Clay__int32_tArray_Add(&resizableContainerBuffer, childElementIndex);
	                }

	                if (sizingAlongAxis) {
	                    innerContentSize += (childSizing.Type == CLAY__SIZING_TYPE_PERCENT ? 0 : childSize);
	                    if (childSizing.Type == CLAY__SIZING_TYPE_GROW) {
	                        growContainerCount++;
	                    }
	                    if (childOffset > 0) {
	                        innerContentSize += parentChildGap; // For children after index 0, the childAxisOffset is the gap from the previous child
	                        totalPaddingAndChildGaps += parentChildGap;
	                    }
	                } else {
	                    innerContentSize = max(childSize, innerContentSize);
	                }
	            }

	            // Expand percentage containers to size
	            for (int32 childOffset = 0; childOffset < parent.childrenOrTextContent.children.length; childOffset++) {
	                int32 childElementIndex = parent.childrenOrTextContent.children.elements[childOffset];
	                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, childElementIndex);
	                Clay_SizingAxis childSizing = xAxis ? childElement.layoutConfig.sizing.Width : childElement.layoutConfig.sizing.Height;
	                float *childSize = xAxis ? &childElement.dimensions.X : &childElement.dimensions.Y;
	                if (childSizing.Type == CLAY__SIZING_TYPE_PERCENT) {
	                    *childSize = (parentSize - totalPaddingAndChildGaps) * childSizing.size.percent;
	                    if (sizingAlongAxis) {
	                        innerContentSize += *childSize;
	                    }
	                    Clay__UpdateAspectRatioBox(childElement);
	                }
	            }

	            if (sizingAlongAxis) {
	                float sizeToDistribute = parentSize - parentPadding - innerContentSize;
	                // The content is too large, compress the children as much as possible
	                if (sizeToDistribute < 0) {
	                    // If the parent can scroll in the axis direction in this direction, don't compress children, just leave them alone
	                    Clay_ScrollElementConfig *scrollElementConfig = Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
	                    if (scrollElementConfig) {
	                        if (((xAxis && scrollElementConfig.horizontal) || (!xAxis && scrollElementConfig.vertical))) {
	                            continue;
	                        }
	                    }
	                    // Scrolling containers preferentially compress before others
	                    while (sizeToDistribute < -CLAY__EPSILON && resizableContainerBuffer.length > 0) {
	                        float largest = 0;
	                        float secondLargest = 0;
	                        float widthToAdd = sizeToDistribute;
	                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
	                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
	                            float childSize = xAxis ? child.dimensions.X : child.dimensions.Y;
	                            if (Clay__FloatEqual(childSize, largest)) { continue; }
	                            if (childSize > largest) {
	                                secondLargest = largest;
	                                largest = childSize;
	                            }
	                            if (childSize < largest) {
	                                secondLargest = max(secondLargest, childSize);
	                                widthToAdd = secondLargest - largest;
	                            }
	                        }

	                        widthToAdd = max(widthToAdd, sizeToDistribute / resizableContainerBuffer.length);

	                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
	                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
	                            float *childSize = xAxis ? &child.dimensions.X : &child.dimensions.Y;
	                            float minSize = xAxis ? child.minDimensions.X : child.minDimensions.Y;
	                            float previousWidth = *childSize;
	                            if (Clay__FloatEqual(*childSize, largest)) {
	                                *childSize += widthToAdd;
	                                if (*childSize <= minSize) {
	                                    *childSize = minSize;
	                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
	                                }
	                                sizeToDistribute -= (*childSize - previousWidth);
	                            }
	                        }
	                    }
	                // The content is too small, allow SIZING_GROW containers to expand
	                } else if (sizeToDistribute > 0 && growContainerCount > 0) {
	                    for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
	                        Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
	                        Clay__SizingType childSizing = xAxis ? child.layoutConfig.sizing.Width.Type : child.layoutConfig.sizing.Height.Type;
	                        if (childSizing != CLAY__SIZING_TYPE_GROW) {
	                            Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
	                        }
	                    }
	                    while (sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.length > 0) {
	                        float smallest = math.MaxFloat32;
	                        float secondSmallest = math.MaxFloat32;
	                        float widthToAdd = sizeToDistribute;
	                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
	                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
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
	                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
	                            float *childSize = xAxis ? &child.dimensions.X : &child.dimensions.Y;
	                            float maxSize = xAxis ? child.layoutConfig.sizing.Width.size.MinMax.Max : child.layoutConfig.sizing.Height.size.MinMax.Max;
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
	                }
	            // Sizing along the non layout axis ("off axis")
	            } else {
	                for (int32 childOffset = 0; childOffset < resizableContainerBuffer.length; childOffset++) {
	                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childOffset));
	                    Clay_SizingAxis childSizing = xAxis ? childElement.layoutConfig.sizing.Width : childElement.layoutConfig.sizing.Height;
	                    float *childSize = xAxis ? &childElement.dimensions.X : &childElement.dimensions.Y;

	                    if (!xAxis && elementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_IMAGE)) {
	                        continue; // Currently we don't support resizing aspect ratio images on the Y axis because it would break the ratio
	                    }

	                    // If we're laying out the children of a scroll panel, grow containers expand to the height of the inner content, not the outer container
	                    float maxSize = parentSize - parentPadding;
	                    if (elementHasConfig(parent, CLAY__ELEMENT_CONFIG_TYPE_SCROLL)) {
	                        Clay_ScrollElementConfig *scrollElementConfig = Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
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
	            }
	        }
	    }
	}

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

	void Clay__AddRenderCommand(Clay_RenderCommand renderCommand) {
	    context := GetCurrentContext();
	    if (context.renderCommands.length < context.renderCommands.capacity - 1) {
	        Clay_RenderCommandArray_Add(&context.renderCommands, renderCommand);
	    } else {
	        if (!context.booleanWarnings.maxRenderCommandsExceeded) {
	            context.booleanWarnings.maxRenderCommandsExceeded = true;
	            context.errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	                .errorType = CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
	                .errorText = CLAY_STRING("Clay ran out of capacity while attempting to create render commands. This is usually caused by a large amount of wrapping text elements while close to the max element capacity. Try using Clay_SetMaxElementCount() with a higher value."),
	                .userData = context.errorHandler.userData });
	        }
	    }
	}

	bool Clay__ElementIsOffscreen(Clay_BoundingBox *boundingBox) {
	    context := GetCurrentContext();
	    if (context.disableCulling) {
	        return false;
	    }

	    return (boundingBox.x > (float32)context.layoutDimensions.X) ||
	           (boundingBox.y > (float32)context.layoutDimensions.Y) ||
	           (boundingBox.x + boundingBox.X < 0) ||
	           (boundingBox.y + boundingBox.Y < 0);
	}

	void Clay__CalculateFinalLayout(void) {
	    context := GetCurrentContext();
	    // Calculate sizing along the X axis
	    Clay__SizeContainersAlongAxis(true);

	    // Wrap text
	    for (int32 textElementIndex = 0; textElementIndex < context.textElementData.length; ++textElementIndex) {
	        Clay__TextElementData *textElementData = Clay__TextElementDataArray_Get(&context.textElementData, textElementIndex);
	        textElementData.wrappedLines = CLAY__INIT(Clay__WrappedTextLineArraySlice) { .length = 0, .internalArray = &context.wrappedTextLines.internalArray[context.wrappedTextLines.length] };
	        Clay_LayoutElement *containerElement = Clay_LayoutElementArray_Get(&context.layoutElements, (int)textElementData.elementIndex);
	        Clay_TextElementConfig *textConfig = Clay__FindElementConfigWithType(containerElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig;
	        Clay__MeasureTextCacheItem *measureTextCacheItem = Clay__MeasureTextCached(&textElementData.text, textConfig);
	        float lineWidth = 0;
	        float lineHeight = textConfig.lineHeight > 0 ? (float32)textConfig.lineHeight : textElementData.preferredDimensions.Y;
	        int32 lineLengthChars = 0;
	        int32 lineStartOffset = 0;
	        if (!measureTextCacheItem.containsNewlines && textElementData.preferredDimensions.X <= containerElement.dimensions.X) {
	            Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { containerElement.dimensions,  textElementData.text });
	            textElementData.wrappedLines.length++;
	            continue;
	        }
	        float spaceWidth = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = 1, .chars = SPACECHAR.chars, .baseChars = SPACECHAR.chars }, textConfig, context.measureTextUserData).X;
	        int32 wordIndex = measureTextCacheItem.measuredWordsStartIndex;
	        while (wordIndex != -1) {
	            if (context.wrappedTextLines.length > context.wrappedTextLines.capacity - 1) {
	                break;
	            }
	            Clay__MeasuredWord *measuredWord = Clay__MeasuredWordArray_Get(&context.measuredWords, wordIndex);
	            // Only word on the line is too large, just render it anyway
	            if (lineLengthChars == 0 && lineWidth + measuredWord.X > containerElement.dimensions.X) {
	                Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { measuredWord.X, lineHeight }, { .length = measuredWord.length, .chars = &textElementData.text.chars[measuredWord.startOffset] } });
	                textElementData.wrappedLines.length++;
	                wordIndex = measuredWord.next;
	                lineStartOffset = measuredWord.startOffset + measuredWord.length;
	            }
	            // measuredWord.length == 0 means a newline character
	            else if (measuredWord.length == 0 || lineWidth + measuredWord.X > containerElement.dimensions.X) {
	                // Wrapped text lines list has overflowed, just render out the line
	                bool finalCharIsSpace = textElementData.text.chars[lineStartOffset + lineLengthChars - 1] == ' ';
	                Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { lineWidth + (finalCharIsSpace ? -spaceWidth : 0), lineHeight }, { .length = lineLengthChars + (finalCharIsSpace ? -1 : 0), .chars = &textElementData.text.chars[lineStartOffset] } });
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
	            Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { lineWidth, lineHeight }, {.length = lineLengthChars, .chars = &textElementData.text.chars[lineStartOffset] } });
	            textElementData.wrappedLines.length++;
	        }
	        containerElement.dimensions.Y = lineHeight * (float32)textElementData.wrappedLines.length;
	    }

	    // Scale vertical image heights according to aspect ratio
	    for (int32 i = 0; i < context.imageElementPointers.length; ++i) {
	        Clay_LayoutElement* imageElement = Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.imageElementPointers, i));
	        Clay_ImageElementConfig *config = Clay__FindElementConfigWithType(imageElement, CLAY__ELEMENT_CONFIG_TYPE_IMAGE).imageElementConfig;
	        imageElement.dimensions.Y = (config.sourceDimensions.Y / max(config.sourceDimensions.X, 1)) * imageElement.dimensions.X;
	    }

	    // Propagate effect of text wrapping, image aspect scaling etc. on height of parents
	    Clay__LayoutElementTreeNodeArray dfsBuffer = context.layoutElementTreeNodeArray1;
	    dfsBuffer.length = 0;
	    for (int32 i = 0; i < context.layoutElementTreeRoots.length; ++i) {
	        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, i);
	        context.treeNodeVisited.internalArray[dfsBuffer.length] = false;
	        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context.layoutElements, (int)root.layoutElementIndex) });
	    }
	    while (dfsBuffer.length > 0) {
	        Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
	        Clay_LayoutElement *currentElement = currentElementTreeNode.layoutElement;
	        if (!context.treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
	            context.treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
	            // If the element has no children or is the container for a text element, don't bother inspecting it
	            if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement.childrenOrTextContent.children.length == 0) {
	                dfsBuffer.length--;
	                continue;
	            }
	            // Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
	            for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; i++) {
	                context.treeNodeVisited.internalArray[dfsBuffer.length] = false;
	                Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[i]) });
	            }
	            continue;
	        }
	        dfsBuffer.length--;

	        // DFS node has been visited, this is on the way back up to the root
	        Clay_LayoutConfig *layoutConfig = currentElement.layoutConfig;
	        if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
	            // Resize any parent containers that have grown in height along their non layout axis
	            for (int32 j = 0; j < currentElement.childrenOrTextContent.children.length; ++j) {
	                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[j]);
	                float childHeightWithPadding = max(childElement.dimensions.Y + layoutConfig.padding.Top + layoutConfig.padding.Bottom, currentElement.dimensions.Y);
	                currentElement.dimensions.Y = min(max(childHeightWithPadding, layoutConfig.sizing.Height.size.MinMax.Min), layoutConfig.sizing.Height.size.MinMax.Max);
	            }
	        } else if (layoutConfig.layoutDirection == CLAY_TOP_TO_BOTTOM) {
	            // Resizing along the layout axis
	            float contentHeight = (float32)(layoutConfig.padding.Top + layoutConfig.padding.Bottom);
	            for (int32 j = 0; j < currentElement.childrenOrTextContent.children.length; ++j) {
	                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[j]);
	                contentHeight += childElement.dimensions.Y;
	            }
	            contentHeight += (float32)(max(currentElement.childrenOrTextContent.children.length - 1, 0) * layoutConfig.childGap);
	            currentElement.dimensions.Y = min(max(contentHeight, layoutConfig.sizing.Height.size.MinMax.Min), layoutConfig.sizing.Height.size.MinMax.Max);
	        }
	    }

	    // Calculate sizing along the Y axis
	    Clay__SizeContainersAlongAxis(false);

	    // Sort tree roots by z-index
	    int32 sortMax = context.layoutElementTreeRoots.length - 1;
	    while (sortMax > 0) { // todo dumb bubble sort
	        for (int32 i = 0; i < sortMax; ++i) {
	            Clay__LayoutElementTreeRoot current = *Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, i);
	            Clay__LayoutElementTreeRoot next = *Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, i + 1);
	            if (next.zIndex < current.zIndex) {
	                Clay__LayoutElementTreeRootArray_Set(&context.layoutElementTreeRoots, i, next);
	                Clay__LayoutElementTreeRootArray_Set(&context.layoutElementTreeRoots, i + 1, current);
	            }
	        }
	        sortMax--;
	    }

	    // Calculate final positions and generate render commands
	    context.renderCommands.length = 0;
	    dfsBuffer.length = 0;
	    for (int32 rootIndex = 0; rootIndex < context.layoutElementTreeRoots.length; ++rootIndex) {
	        dfsBuffer.length = 0;
	        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex);
	        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context.layoutElements, (int)root.layoutElementIndex);
	        vector2.Float32 rootPosition = CLAY__DEFAULT_STRUCT;
	        Clay_LayoutElementHashMapItem *parentHashMapItem = Clay__GetHashMapItem(root.parentId);
	        // Position root floating containers
	        if (elementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) && parentHashMapItem) {
	            Clay_FloatingElementConfig *config = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
	            vector2.Float32 rootDimensions = rootElement.dimensions;
	            Clay_BoundingBox parentBoundingBox = parentHashMapItem.boundingBox;
	            // Set X position
	            vector2.Float32 targetAttachPosition = CLAY__DEFAULT_STRUCT;
	            switch (config.attachPoints.parent) {
	                case CLAY_ATTACH_POINT_LEFT_TOP:
	                case CLAY_ATTACH_POINT_LEFT_CENTER:
	                case CLAY_ATTACH_POINT_LEFT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x; break;
	                case CLAY_ATTACH_POINT_CENTER_TOP:
	                case CLAY_ATTACH_POINT_CENTER_CENTER:
	                case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + (parentBoundingBox.X / 2); break;
	                case CLAY_ATTACH_POINT_RIGHT_TOP:
	                case CLAY_ATTACH_POINT_RIGHT_CENTER:
	                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + parentBoundingBox.X; break;
	            }
	            switch (config.attachPoints.element) {
	                case CLAY_ATTACH_POINT_LEFT_TOP:
	                case CLAY_ATTACH_POINT_LEFT_CENTER:
	                case CLAY_ATTACH_POINT_LEFT_BOTTOM: break;
	                case CLAY_ATTACH_POINT_CENTER_TOP:
	                case CLAY_ATTACH_POINT_CENTER_CENTER:
	                case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x -= (rootDimensions.X / 2); break;
	                case CLAY_ATTACH_POINT_RIGHT_TOP:
	                case CLAY_ATTACH_POINT_RIGHT_CENTER:
	                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x -= rootDimensions.X; break;
	            }
	            switch (config.attachPoints.parent) { // I know I could merge the x and y switch statements, but this is easier to read
	                case CLAY_ATTACH_POINT_LEFT_TOP:
	                case CLAY_ATTACH_POINT_RIGHT_TOP:
	                case CLAY_ATTACH_POINT_CENTER_TOP: targetAttachPosition.y = parentBoundingBox.y; break;
	                case CLAY_ATTACH_POINT_LEFT_CENTER:
	                case CLAY_ATTACH_POINT_CENTER_CENTER:
	                case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y = parentBoundingBox.y + (parentBoundingBox.Y / 2); break;
	                case CLAY_ATTACH_POINT_LEFT_BOTTOM:
	                case CLAY_ATTACH_POINT_CENTER_BOTTOM:
	                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y = parentBoundingBox.y + parentBoundingBox.Y; break;
	            }
	            switch (config.attachPoints.element) {
	                case CLAY_ATTACH_POINT_LEFT_TOP:
	                case CLAY_ATTACH_POINT_RIGHT_TOP:
	                case CLAY_ATTACH_POINT_CENTER_TOP: break;
	                case CLAY_ATTACH_POINT_LEFT_CENTER:
	                case CLAY_ATTACH_POINT_CENTER_CENTER:
	                case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y -= (rootDimensions.Y / 2); break;
	                case CLAY_ATTACH_POINT_LEFT_BOTTOM:
	                case CLAY_ATTACH_POINT_CENTER_BOTTOM:
	                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y -= rootDimensions.Y; break;
	            }
	            targetAttachPosition.x += config.offset.x;
	            targetAttachPosition.y += config.offset.y;
	            rootPosition = targetAttachPosition;
	        }
	        if (root.clipElementId) {
	            Clay_LayoutElementHashMapItem *clipHashMapItem = Clay__GetHashMapItem(root.clipElementId);
	            if (clipHashMapItem) {
	                // Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
	                if (context.externalScrollHandlingEnabled) {
	                    Clay_ScrollElementConfig *scrollConfig = Clay__FindElementConfigWithType(clipHashMapItem.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
	                    for (int32 i = 0; i < context.scrollContainerDatas.length; i++) {
	                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i);
	                        if (mapping.layoutElement == clipHashMapItem.layoutElement) {
	                            root.pointerOffset = mapping.scrollPosition;
	                            if (scrollConfig.horizontal) {
	                                rootPosition.x += mapping.scrollPosition.x;
	                            }
	                            if (scrollConfig.vertical) {
	                                rootPosition.y += mapping.scrollPosition.y;
	                            }
	                            break;
	                        }
	                    }
	                }
	                Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	                    .boundingBox = clipHashMapItem.boundingBox,
	                    .userData = 0,
	                    .id = Clay__HashNumber(rootElement.id, rootElement.childrenOrTextContent.children.length + 10).id, // TODO need a better strategy for managing derived ids
	                    .zIndex = root.zIndex,
	                    .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
	                });
	            }
	        }
	        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = rootElement, .position = rootPosition, .nextChildOffset = { .x = (float32)rootElement.layoutConfig.padding.Left, .y = (float32)rootElement.layoutConfig.padding.Top } });

	        context.treeNodeVisited.internalArray[0] = false;
	        while (dfsBuffer.length > 0) {
	            Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
	            Clay_LayoutElement *currentElement = currentElementTreeNode.layoutElement;
	            Clay_LayoutConfig *layoutConfig = currentElement.layoutConfig;
	            vector2.Float32 scrollOffset = CLAY__DEFAULT_STRUCT;

	            // This will only be run a single time for each element in downwards DFS order
	            if (!context.treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
	                context.treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;

	                Clay_BoundingBox currentElementBoundingBox = { currentElementTreeNode.position.x, currentElementTreeNode.position.y, currentElement.dimensions.X, currentElement.dimensions.Y };
	                if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
	                    Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
	                    vector2.Float32 expand = floatingElementConfig.expand;
	                    currentElementBoundingBox.x -= expand.X;
	                    currentElementBoundingBox.X += expand.X * 2;
	                    currentElementBoundingBox.y -= expand.Y;
	                    currentElementBoundingBox.Y += expand.Y * 2;
	                }

	                Clay__ScrollContainerDataInternal *scrollContainerData = CLAY__NULL;
	                // Apply scroll offsets to container
	                if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SCROLL)) {
	                    Clay_ScrollElementConfig *scrollConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;

	                    // This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
	                    for (int32 i = 0; i < context.scrollContainerDatas.length; i++) {
	                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i);
	                        if (mapping.layoutElement == currentElement) {
	                            scrollContainerData = mapping;
	                            mapping.boundingBox = currentElementBoundingBox;
	                            if (scrollConfig.horizontal) {
	                                scrollOffset.x = mapping.scrollPosition.x;
	                            }
	                            if (scrollConfig.vertical) {
	                                scrollOffset.y = mapping.scrollPosition.y;
	                            }
	                            if (context.externalScrollHandlingEnabled) {
	                                scrollOffset = CLAY__INIT(vector2.Float32) CLAY__DEFAULT_STRUCT;
	                            }
	                            break;
	                        }
	                    }
	                }

	                Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(currentElement.id);
	                if (hashMapItem) {
	                    hashMapItem.boundingBox = currentElementBoundingBox;
	                    if (hashMapItem.idAlias) {
	                        Clay_LayoutElementHashMapItem *hashMapItemAlias = Clay__GetHashMapItem(hashMapItem.idAlias);
	                        if (hashMapItemAlias) {
	                            hashMapItemAlias.boundingBox = currentElementBoundingBox;
	                        }
	                    }
	                }

	                int32 sortedConfigIndexes[20];
	                for (int32 elementConfigIndex = 0; elementConfigIndex < currentElement.elementConfigs.length; ++elementConfigIndex) {
	                    sortedConfigIndexes[elementConfigIndex] = elementConfigIndex;
	                }
	                sortMax = currentElement.elementConfigs.length - 1;
	                while (sortMax > 0) { // todo dumb bubble sort
	                    for (int32 i = 0; i < sortMax; ++i) {
	                        int32 current = sortedConfigIndexes[i];
	                        int32 next = sortedConfigIndexes[i + 1];
	                        Clay__ElementConfigType currentType = Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, current).type;
	                        Clay__ElementConfigType nextType = Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, next).type;
	                        if (nextType == CLAY__ELEMENT_CONFIG_TYPE_SCROLL || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
	                            sortedConfigIndexes[i] = next;
	                            sortedConfigIndexes[i + 1] = current;
	                        }
	                    }
	                    sortMax--;
	                }

	                bool emitRectangle = false;
	                // Create the render commands for this element
	                Clay_SharedElementConfig *sharedConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig;
	                if (sharedConfig && sharedConfig.backgroundColor.a > 0) {
	                   emitRectangle = true;
	                }
	                else if (!sharedConfig) {
	                    emitRectangle = false;
	                    sharedConfig = &Clay_SharedElementConfig_DEFAULT;
	                }
	                for (int32 elementConfigIndex = 0; elementConfigIndex < currentElement.elementConfigs.length; ++elementConfigIndex) {
	                    Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, sortedConfigIndexes[elementConfigIndex]);
	                    Clay_RenderCommand renderCommand = {
	                        .boundingBox = currentElementBoundingBox,
	                        .userData = sharedConfig.userData,
	                        .id = currentElement.id,
	                    };

	                    bool offscreen = Clay__ElementIsOffscreen(&currentElementBoundingBox);
	                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
	                    bool shouldRender = !offscreen;
	                    switch (elementConfig.type) {
	                        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
	                        case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
	                        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
	                            shouldRender = false;
	                            break;
	                        }
	                        case CLAY__ELEMENT_CONFIG_TYPE_SCROLL: {
	                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START;
	                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
	                                .scroll = {
	                                    .horizontal = elementConfig.config.scrollElementConfig.horizontal,
	                                    .vertical = elementConfig.config.scrollElementConfig.vertical,
	                                }
	                            };
	                            break;
	                        }
	                        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
	                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_IMAGE;
	                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
	                                .image = {
	                                    .backgroundColor = sharedConfig.backgroundColor,
	                                    .cornerRadius = sharedConfig.cornerRadius,
	                                    .sourceDimensions = elementConfig.config.imageElementConfig.sourceDimensions,
	                                    .imageData = elementConfig.config.imageElementConfig.imageData,
	                               }
	                            };
	                            emitRectangle = false;
	                            break;
	                        }
	                        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
	                            if (!shouldRender) {
	                                break;
	                            }
	                            shouldRender = false;
	                            Clay_ElementConfigUnion configUnion = elementConfig.config;
	                            Clay_TextElementConfig *textElementConfig = configUnion.textElementConfig;
	                            float naturalLineHeight = currentElement.childrenOrTextContent.textElementData.preferredDimensions.Y;
	                            float finalLineHeight = textElementConfig.lineHeight > 0 ? (float32)textElementConfig.lineHeight : naturalLineHeight;
	                            float lineHeightOffset = (finalLineHeight - naturalLineHeight) / 2;
	                            float yPosition = lineHeightOffset;
	                            for (int32 lineIndex = 0; lineIndex < currentElement.childrenOrTextContent.textElementData.wrappedLines.length; ++lineIndex) {
	                                Clay__WrappedTextLine *wrappedLine = Clay__WrappedTextLineArraySlice_Get(&currentElement.childrenOrTextContent.textElementData.wrappedLines, lineIndex);
	                                if (wrappedLine.line.length == 0) {
	                                    yPosition += finalLineHeight;
	                                    continue;
	                                }
	                                float offset = (currentElementBoundingBox.X - wrappedLine.dimensions.X);
	                                if (textElementConfig.textAlignment == CLAY_TEXT_ALIGN_LEFT) {
	                                    offset = 0;
	                                }
	                                if (textElementConfig.textAlignment == CLAY_TEXT_ALIGN_CENTER) {
	                                    offset /= 2;
	                                }
	                                Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	                                    .boundingBox = { currentElementBoundingBox.x + offset, currentElementBoundingBox.y + yPosition, wrappedLine.dimensions.X, wrappedLine.dimensions.Y },
	                                    .renderData = { .text = {
	                                        .stringContents = CLAY__INIT(Clay_StringSlice) { .length = wrappedLine.line.length, .chars = wrappedLine.line.chars, .baseChars = currentElement.childrenOrTextContent.textElementData.text.chars },
	                                        .textColor = textElementConfig.textColor,
	                                        .fontId = textElementConfig.fontId,
	                                        .fontSize = textElementConfig.fontSize,
	                                        .letterSpacing = textElementConfig.letterSpacing,
	                                        .lineHeight = textElementConfig.lineHeight,
	                                    }},
	                                    .userData = textElementConfig.userData,
	                                    .id = Clay__HashNumber(lineIndex, currentElement.id).id,
	                                    .zIndex = root.zIndex,
	                                    .commandType = CLAY_RENDER_COMMAND_TYPE_TEXT,
	                                });
	                                yPosition += finalLineHeight;

	                                if (!context.disableCulling && (currentElementBoundingBox.y + yPosition > context.layoutDimensions.Y)) {
	                                    break;
	                                }
	                            }
	                            break;
	                        }
	                        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: {
	                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_CUSTOM;
	                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
	                                .custom = {
	                                    .backgroundColor = sharedConfig.backgroundColor,
	                                    .cornerRadius = sharedConfig.cornerRadius,
	                                    .customData = elementConfig.config.customElementConfig.customData,
	                                }
	                            };
	                            emitRectangle = false;
	                            break;
	                        }
	                        default: break;
	                    }
	                    if (shouldRender) {
	                        Clay__AddRenderCommand(renderCommand);
	                    }
	                    if (offscreen) {
	                        // NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
	                        // Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
	                        // be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
	                    }
	                }

	                if (emitRectangle) {
	                    Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	                        .boundingBox = currentElementBoundingBox,
	                        .renderData = { .rectangle = {
	                                .backgroundColor = sharedConfig.backgroundColor,
	                                .cornerRadius = sharedConfig.cornerRadius,
	                        }},
	                        .userData = sharedConfig.userData,
	                        .id = currentElement.id,
	                        .zIndex = root.zIndex,
	                        .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
	                    });
	                }

	                // Setup initial on-axis alignment
	                if (!elementHasConfig(currentElementTreeNode.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
	                    vector2.Float32 contentSize = {0,0};
	                    if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
	                        for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
	                            Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
	                            contentSize.X += childElement.dimensions.X;
	                            contentSize.Y = max(contentSize.Y, childElement.dimensions.Y);
	                        }
	                        contentSize.X += (float32)(max(currentElement.childrenOrTextContent.children.length - 1, 0) * layoutConfig.childGap);
	                        float extraSpace = currentElement.dimensions.X - (float32)(layoutConfig.padding.Left + layoutConfig.padding.Right) - contentSize.X;
	                        switch (layoutConfig.childAlignment.x) {
	                            case CLAY_ALIGN_X_LEFT: extraSpace = 0; break;
	                            case CLAY_ALIGN_X_CENTER: extraSpace /= 2; break;
	                            default: break;
	                        }
	                        currentElementTreeNode.nextChildOffset.x += extraSpace;
	                    } else {
	                        for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
	                            Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
	                            contentSize.X = max(contentSize.X, childElement.dimensions.X);
	                            contentSize.Y += childElement.dimensions.Y;
	                        }
	                        contentSize.Y += (float32)(max(currentElement.childrenOrTextContent.children.length - 1, 0) * layoutConfig.childGap);
	                        float extraSpace = currentElement.dimensions.Y - (float32)(layoutConfig.padding.Top + layoutConfig.padding.Bottom) - contentSize.Y;
	                        switch (layoutConfig.childAlignment.y) {
	                            case CLAY_ALIGN_Y_TOP: extraSpace = 0; break;
	                            case CLAY_ALIGN_Y_CENTER: extraSpace /= 2; break;
	                            default: break;
	                        }
	                        currentElementTreeNode.nextChildOffset.y += extraSpace;
	                    }

	                    if (scrollContainerData) {
	                        scrollContainerData.contentSize = CLAY__INIT(vector2.Float32) { contentSize.X + (float32)(layoutConfig.padding.Left + layoutConfig.padding.Right), contentSize.Y + (float32)(layoutConfig.padding.Top + layoutConfig.padding.Bottom) };
	                    }
	                }
	            }
	            else {
	                // DFS is returning upwards backwards
	                bool closeScrollElement = false;
	                Clay_ScrollElementConfig *scrollConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
	                if (scrollConfig) {
	                    closeScrollElement = true;
	                    for (int32 i = 0; i < context.scrollContainerDatas.length; i++) {
	                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i);
	                        if (mapping.layoutElement == currentElement) {
	                            if (scrollConfig.horizontal) { scrollOffset.x = mapping.scrollPosition.x; }
	                            if (scrollConfig.vertical) { scrollOffset.y = mapping.scrollPosition.y; }
	                            if (context.externalScrollHandlingEnabled) {
	                                scrollOffset = CLAY__INIT(vector2.Float32) CLAY__DEFAULT_STRUCT;
	                            }
	                            break;
	                        }
	                    }
	                }

	                if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER)) {
	                    Clay_LayoutElementHashMapItem *currentElementData = Clay__GetHashMapItem(currentElement.id);
	                    Clay_BoundingBox currentElementBoundingBox = currentElementData.boundingBox;

	                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
	                    if (!Clay__ElementIsOffscreen(&currentElementBoundingBox)) {
	                        Clay_SharedElementConfig *sharedConfig = elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) ? Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig : &Clay_SharedElementConfig_DEFAULT;
	                        Clay_BorderElementConfig *borderConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).borderElementConfig;
	                        Clay_RenderCommand renderCommand = {
	                                .boundingBox = currentElementBoundingBox,
	                                .renderData = { .border = {
	                                    .color = borderConfig.color,
	                                    .cornerRadius = sharedConfig.cornerRadius,
	                                    .X = borderConfig.X
	                                }},
	                                .userData = sharedConfig.userData,
	                                .id = Clay__HashNumber(currentElement.id, currentElement.childrenOrTextContent.children.length).id,
	                                .commandType = CLAY_RENDER_COMMAND_TYPE_BORDER,
	                        };
	                        Clay__AddRenderCommand(renderCommand);
	                        if (borderConfig.X.betweenChildren > 0 && borderConfig.color.a > 0) {
	                            float halfGap = layoutConfig.childGap / 2;
	                            vector2.Float32 borderOffset = { (float32)layoutConfig.padding.Left - halfGap, (float32)layoutConfig.padding.Top - halfGap };
	                            if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
	                                for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
	                                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
	                                    if (i > 0) {
	                                        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	                                            .boundingBox = { currentElementBoundingBox.x + borderOffset.x + scrollOffset.x, currentElementBoundingBox.y + scrollOffset.y, (float32)borderConfig.X.betweenChildren, currentElement.dimensions.Y },
	                                            .renderData = { .rectangle = {
	                                                .backgroundColor = borderConfig.color,
	                                            } },
	                                            .userData = sharedConfig.userData,
	                                            .id = Clay__HashNumber(currentElement.id, currentElement.childrenOrTextContent.children.length + 1 + i).id,
	                                            .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
	                                        });
	                                    }
	                                    borderOffset.x += (childElement.dimensions.X + (float32)layoutConfig.childGap);
	                                }
	                            } else {
	                                for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
	                                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
	                                    if (i > 0) {
	                                        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	                                            .boundingBox = { currentElementBoundingBox.x + scrollOffset.x, currentElementBoundingBox.y + borderOffset.y + scrollOffset.y, currentElement.dimensions.X, (float32)borderConfig.X.betweenChildren },
	                                            .renderData = { .rectangle = {
	                                                    .backgroundColor = borderConfig.color,
	                                            } },
	                                            .userData = sharedConfig.userData,
	                                            .id = Clay__HashNumber(currentElement.id, currentElement.childrenOrTextContent.children.length + 1 + i).id,
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
	                    Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	                        .id = Clay__HashNumber(currentElement.id, rootElement.childrenOrTextContent.children.length + 11).id,
	                        .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
	                    });
	                }

	                dfsBuffer.length--;
	                continue;
	            }

	            // Add children to the DFS buffer
	            if (!elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
	                dfsBuffer.length += currentElement.childrenOrTextContent.children.length;
	                for (int32 i = 0; i < currentElement.childrenOrTextContent.children.length; ++i) {
	                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context.layoutElements, currentElement.childrenOrTextContent.children.elements[i]);
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
	                    context.treeNodeVisited.internalArray[newNodeIndex] = false;

	                    // Update parent offsets
	                    if (layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) {
	                        currentElementTreeNode.nextChildOffset.x += childElement.dimensions.X + (float32)layoutConfig.childGap;
	                    } else {
	                        currentElementTreeNode.nextChildOffset.y += childElement.dimensions.Y + (float32)layoutConfig.childGap;
	                    }
	                }
	            }
	        }

	        if (root.clipElementId) {
	            Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) { .id = Clay__HashNumber(rootElement.id, rootElement.childrenOrTextContent.children.length + 11).id, .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END });
	        }
	    }
	}

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
	            Clay_LayoutElementHashMapItem *currentElementData = Clay__GetHashMapItem(currentElement.id);
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
	                Clay_LayoutElementHashMapItem *highlightedItem = Clay__GetHashMapItem(elementId.offset);
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

	void Clay__RenderDebugViewElementConfigHeader(string elementId, Clay__ElementConfigType type) {
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

	void Clay__RenderDebugView(void) {
	    context := GetCurrentContext();
	    Clay_ElementId closeButtonId = hashString(CLAY_STRING("Clay__DebugViewTopHeaderCloseButtonOuter"), 0, 0);
	    if (context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
	        for (int32 i = 0; i < context.pointerOverIds.length; ++i) {
	            Clay_ElementId *elementId = Clay__ElementIdArray_Get(&context.pointerOverIds, i);
	            if (elementId.id == closeButtonId.id) {
	                context.debugModeEnabled = false;
	                return;
	            }
	        }
	    }

	    uint32 initialRootsLength = context.layoutElementTreeRoots.length;
	    uint32 initialElementsLength = context.layoutElements.length;
	    Clay_TextElementConfig *infoTextConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
	    Clay_TextElementConfig *infoTitleConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
	    Clay_ElementId scrollId = hashString(CLAY_STRING("Clay__DebugViewOuterScrollPane"), 0, 0);
	    float scrollYOffset = 0;
	    bool pointerInDebugView = context.pointerInfo.position.y < context.layoutDimensions.Y - 300;
	    for (int32 i = 0; i < context.scrollContainerDatas.length; ++i) {
	        Clay__ScrollContainerDataInternal *scrollContainerData = Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i);
	        if (scrollContainerData.elementId == scrollId.id) {
	            if (!context.externalScrollHandlingEnabled) {
	                scrollYOffset = scrollContainerData.scrollPosition.y;
	            } else {
	                pointerInDebugView = context.pointerInfo.position.y + scrollContainerData.scrollPosition.y < context.layoutDimensions.Y - 300;
	            }
	            break;
	        }
	    }
	    int32 highlightedRow = pointerInDebugView
	            ? (int32)((context.pointerInfo.position.y - scrollYOffset) / (float32)CLAY__DEBUGVIEW_ROW_HEIGHT) - 1
	            : -1;
	    if (context.pointerInfo.position.x < context.layoutDimensions.X - (float32)debugViewWidth) {
	        highlightedRow = -1;
	    }
	    Clay__RenderDebugLayoutData layoutData = CLAY__DEFAULT_STRUCT;
	    CLAY({ .id = CLAY_ID("Clay__DebugView"),
	         .layout = { .sizing = { CLAY_SIZING_FIXED((float32)debugViewWidth) , CLAY_SIZING_FIXED(context.layoutDimensions.Y) }, .layoutDirection = CLAY_TOP_TO_BOTTOM },
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
	                float contentWidth = Clay__GetHashMapItem(panelContentsId.id).layoutElement.dimensions.X;
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
	        if (context.debugSelectedElementId != 0) {
	            Clay_LayoutElementHashMapItem *selectedItem = Clay__GetHashMapItem(context.debugSelectedElementId);
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
	                        Clay__RenderDebugLayoutSizing(layoutConfig.sizing.Width, infoTextConfig);
	                    }
	                    CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
	                        CLAY_TEXT(CLAY_STRING("height: "), infoTextConfig);
	                        Clay__RenderDebugLayoutSizing(layoutConfig.sizing.Height, infoTextConfig);
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
	                                Clay_LayoutElementHashMapItem *hashItem = Clay__GetHashMapItem(floatingConfig.parentId);
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
	                int32 previousWarningsLength = context.warnings.length;
	                for (int32 i = 0; i < previousWarningsLength; i++) {
	                    Clay__Warning warning = context.warnings.internalArray[i];
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
	}

#pragma endregion
*/
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
            .errorType = CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
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
                .errorType = CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
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
            .errorType = CLAY_ERROR_TYPE_INTERNAL_ERROR,
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
        .errorType = CLAY_ERROR_TYPE_INTERNAL_ERROR,
        .errorText = CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
        .userData = context.errorHandler.userData });
    return false;
}
*/
