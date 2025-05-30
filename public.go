package clay

import (
	"github.com/igadmg/gamemath/vector2"
)

/// Clay Git Reference - clay used to convert, keep it and update as you update goclay code.
// ad49977f1b37ccd7664333181f30f575d08d3838
///

// Returns the size, in bytes, of the minimum amount of memory Clay requires to operate at its current settings.
///func MinMemorySize() uint32 {}

// Creates an arena for clay to use for its internal allocations, given a certain capacity in bytes and a pointer to an allocation of at least that size.
// Intended to be used with clay.MinMemorySize in the following way:
// uint32_t minMemoryRequired = clay.MinMemorySize();
// clay.Arena clayMemory = clay.CreateArenaWithCapacityAndMemory(minMemoryRequired, malloc(minMemoryRequired));
///CLAY_DLL_EXPORT clay.Arena clay.CreateArenaWithCapacityAndMemory(size_t capacity, void *memory);

// Sets the state of the "pointer" (i.e. the mouse or touch) in Clay's internal data. Used for detecting and responding to mouse events in the debug view,
// as well as for clay.Hovered() and scroll element handling.
func (c *Context) SetPointerState(position Vector2, pointerDown bool) {
	if c.booleanWarnings.maxElementsExceeded {
		return
	}

	c.pointerInfo.Position = position
	c.pointerOverIds = c.pointerOverIds[:0]

	var dfsBuffer []int
	treeNodeVisited := make([]bool, len(c.layoutElements))
	for rootIndex := len(c.layoutElementTreeRoots) - 1; rootIndex >= 0; rootIndex-- {
		dfsBuffer = c.layoutElementChildrenBuffer[:0]
		root := c.layoutElementTreeRoots[rootIndex]
		dfsBuffer = append(dfsBuffer, root.layoutElementIndex)
		treeNodeVisited[0] = false
		found := false
		for len(dfsBuffer) > 0 {
			if treeNodeVisited[len(dfsBuffer)-1] {
				dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]
				continue
			}
			treeNodeVisited[len(dfsBuffer)-1] = true
			currentElement := c.layoutElements[dfsBuffer[len(dfsBuffer)-1]]
			mapItem := c.getHashMapItem(currentElement.id) // TODO think of a way around this, maybe the fact that it's essentially a binary tree limits the cost, but the worst case is not great
			clipElementId := uint32(0)                     //TODO: fix c.layoutElementClipElementIds[(int32)(currentElement-c.layoutElements.internalArray)]
			clipItem := c.getHashMapItem(clipElementId)
			if mapItem != nil {
				elementBox := mapItem.boundingBox.AddPositionXY(root.pointerOffset.X, -root.pointerOffset.Y)

				if elementBox.Contains(position) && (clipElementId == 0 || clipItem.boundingBox.Contains(position)) {
					if mapItem.onHoverFunction != nil {
						mapItem.onHoverFunction(mapItem.elementId, c.pointerInfo, mapItem.hoverFunctionUserData)
					}
					c.pointerOverIds = append(c.pointerOverIds, mapItem.elementId)
					found = true

					if mapItem.idAlias != 0 {
						c.pointerOverIds = append(c.pointerOverIds, ElementId{id: mapItem.idAlias})
					}
				}
				if elementHasConfig[*TextElementConfig](&currentElement) {
					dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]
					continue
				}
				for i := len(currentElement.children) - 1; i >= 0; i-- {
					dfsBuffer = append(dfsBuffer, currentElement.children[i])
					treeNodeVisited[len(dfsBuffer)-1] = false // TODO needs to be ranged checked
				}
			} else {
				dfsBuffer = dfsBuffer[:len(dfsBuffer)-1]
			}
		}

		if found {
			rootElement := c.layoutElements[root.layoutElementIndex]
			if config, ok := findElementConfigWithType[*FloatingElementConfig](&rootElement); ok {
				if config.PointerCaptureMode == POINTER_CAPTURE_MODE_CAPTURE {
					break
				}
			}
		}
	}

	if pointerDown {
		if c.pointerInfo.State == POINTER_DATA_PRESSED_THIS_FRAME {
			c.pointerInfo.State = POINTER_DATA_PRESSED
		} else if c.pointerInfo.State != POINTER_DATA_PRESSED {
			c.pointerInfo.State = POINTER_DATA_PRESSED_THIS_FRAME
		}
	} else {
		if c.pointerInfo.State == POINTER_DATA_RELEASED_THIS_FRAME {
			c.pointerInfo.State = POINTER_DATA_RELEASED
		} else if c.pointerInfo.State != POINTER_DATA_RELEASED {
			c.pointerInfo.State = POINTER_DATA_RELEASED_THIS_FRAME
		}
	}
}

// Initialize Clay's internal arena and setup required data before layout can begin. Only needs to be called once.
// - arena can be created using clay.CreateArenaWithCapacityAndMemory()
// - layoutDimensions are the initial bounding dimensions of the layout (i.e. the screen width and height for a full screen layout)
// - errorHandler is used by Clay to inform you if something has gone wrong in configuration or layout.
func Initialize(layoutDimensions Dimensions, errorHandler ErrorHandler) *Context {
	// DEFAULTS
	if errorHandler.ErrorHandlerFunction == nil {
		errorHandler.ErrorHandlerFunction = errorHandlerFunctionDefault
	}

	context := &Context{
		maxElementCount:              defaultMaxElementCount,
		maxMeasureTextCacheWordCount: defaultMaxMeasureTextWordCacheCount,
		errorHandler:                 errorHandler,
		layoutDimensions:             layoutDimensions,
		//internalArena:                arena,
	}

	if oldContext := GetCurrentContext(); oldContext != nil {
		context.maxElementCount = oldContext.maxElementCount
		context.maxMeasureTextCacheWordCount = oldContext.maxMeasureTextCacheWordCount
	}

	SetCurrentContext(context)
	context.initializePersistentMemory()
	context.initializeEphemeralMemory()

	context.measureTextHashMapInternal = append(context.measureTextHashMapInternal, MeasureTextCacheItem{}) // Reserve the 0 value to mean "no next element"
	context.layoutDimensions = layoutDimensions

	return context
}

// Returns the Context that clay is currently using. Used when using multiple instances of clay simultaneously.
func GetCurrentContext() *Context {
	return currentContext
}

// Sets the context that clay will use to compute the layout.
// Used to restore a context saved fromGetCurrentContext when using multiple instances of clay simultaneously.
func SetCurrentContext(context *Context) {
	currentContext = context
}

// Updates the state of Clay's internal scroll data, updating scroll content positions if scrollDelta is non zero, and progressing momentum scrolling.
// - enableDragScrolling when set to true will enable mobile device like "touch drag" scroll of scroll containers, including momentum scrolling after the touch has ended.
// - scrollDelta is the amount to scroll this frame on each axis in pixels.
// - deltaTime is the time in seconds since the last "frame" (scroll update)
func (c *Context) UpdateScrollContainers(enableDragScrolling bool, scrollDelta Vector2, deltaTime float32) {
	/*
	   context := GetCurrentContext();
	       bool isPointerActive = enableDragScrolling && (context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED || context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME);
	       // Don't apply scroll events to ancestors of the inner element
	       int32_t highestPriorityElementIndex = -1;
	       Clay__ScrollContainerDataInternal *highestPriorityScrollData = CLAY__NULL;
	       for (int32_t i = 0; i < context.scrollContainerDatas.length; i++) {
	           Clay__ScrollContainerDataInternal *scrollData = Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i);
	           if (!scrollData.openThisFrame) {
	               Clay__ScrollContainerDataInternalArray_RemoveSwapback(&context.scrollContainerDatas, i);
	               continue;
	           }
	           scrollData.openThisFrame = false;
	           Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(scrollData.elementId);
	           // Element isn't rendered this frame but scroll offset has been retained
	           if (!hashMapItem) {
	               Clay__ScrollContainerDataInternalArray_RemoveSwapback(&context.scrollContainerDatas, i);
	               continue;
	           }

	           // Touch / click is released
	           if (!isPointerActive && scrollData.pointerScrollActive) {
	               float xDiff = scrollData.scrollPosition.x - scrollData.scrollOrigin.x;
	               if (xDiff < -10 || xDiff > 10) {
	                   scrollData.scrollMomentum.x = (scrollData.scrollPosition.x - scrollData.scrollOrigin.x) / (scrollData.momentumTime * 25);
	               }
	               float yDiff = scrollData.scrollPosition.y - scrollData.scrollOrigin.y;
	               if (yDiff < -10 || yDiff > 10) {
	                   scrollData.scrollMomentum.y = (scrollData.scrollPosition.y - scrollData.scrollOrigin.y) / (scrollData.momentumTime * 25);
	               }
	               scrollData.pointerScrollActive = false;

	               scrollData.pointerOrigin = CLAY__INIT(Clay_Vector2){0,0};
	               scrollData.scrollOrigin = CLAY__INIT(Clay_Vector2){0,0};
	               scrollData.momentumTime = 0;
	           }

	           // Apply existing momentum
	           scrollData.scrollPosition.x += scrollData.scrollMomentum.x;
	           scrollData.scrollMomentum.x *= 0.95f;
	           bool scrollOccurred = scrollDelta.x != 0 || scrollDelta.y != 0;
	           if ((scrollData.scrollMomentum.x > -0.1f && scrollData.scrollMomentum.x < 0.1f) || scrollOccurred) {
	               scrollData.scrollMomentum.x = 0;
	           }
	           scrollData.scrollPosition.x = CLAY__MIN(CLAY__MAX(scrollData.scrollPosition.x, -(CLAY__MAX(scrollData.contentSize.width - scrollData.layoutElement.dimensions.width, 0))), 0);

	           scrollData.scrollPosition.y += scrollData.scrollMomentum.y;
	           scrollData.scrollMomentum.y *= 0.95f;
	           if ((scrollData.scrollMomentum.y > -0.1f && scrollData.scrollMomentum.y < 0.1f) || scrollOccurred) {
	               scrollData.scrollMomentum.y = 0;
	           }
	           scrollData.scrollPosition.y = CLAY__MIN(CLAY__MAX(scrollData.scrollPosition.y, -(CLAY__MAX(scrollData.contentSize.height - scrollData.layoutElement.dimensions.height, 0))), 0);

	           for (int32_t j = 0; j < context.pointerOverIds.length; ++j) { // TODO n & m are small here but this being n*m gives me the creeps
	               if (scrollData.layoutElement.id == Clay__ElementIdArray_Get(&context.pointerOverIds, j).id) {
	                   highestPriorityElementIndex = j;
	                   highestPriorityScrollData = scrollData;
	               }
	           }
	       }

	       if (highestPriorityElementIndex > -1 && highestPriorityScrollData) {
	           Clay_LayoutElement *scrollElement = highestPriorityScrollData.layoutElement;
	           Clay_ScrollElementConfig *scrollConfig = findElementConfigWithType(scrollElement, CLAY__ELEMENT_CONFIG_TYPE_SCROLL).scrollElementConfig;
	           bool canScrollVertically = scrollConfig.vertical && highestPriorityScrollData.contentSize.height > scrollElement.dimensions.height;
	           bool canScrollHorizontally = scrollConfig.horizontal && highestPriorityScrollData.contentSize.width > scrollElement.dimensions.width;
	           // Handle wheel scroll
	           if (canScrollVertically) {
	               highestPriorityScrollData.scrollPosition.y = highestPriorityScrollData.scrollPosition.y + scrollDelta.y * 10;
	           }
	           if (canScrollHorizontally) {
	               highestPriorityScrollData.scrollPosition.x = highestPriorityScrollData.scrollPosition.x + scrollDelta.x * 10;
	           }
	           // Handle click / touch scroll
	           if (isPointerActive) {
	               highestPriorityScrollData.scrollMomentum = CLAY__INIT(Clay_Vector2)CLAY__DEFAULT_STRUCT;
	               if (!highestPriorityScrollData.pointerScrollActive) {
	                   highestPriorityScrollData.pointerOrigin = context.pointerInfo.position;
	                   highestPriorityScrollData.scrollOrigin = highestPriorityScrollData.scrollPosition;
	                   highestPriorityScrollData.pointerScrollActive = true;
	               } else {
	                   float scrollDeltaX = 0, scrollDeltaY = 0;
	                   if (canScrollHorizontally) {
	                       float oldXScrollPosition = highestPriorityScrollData.scrollPosition.x;
	                       highestPriorityScrollData.scrollPosition.x = highestPriorityScrollData.scrollOrigin.x + (context.pointerInfo.position.x - highestPriorityScrollData.pointerOrigin.x);
	                       highestPriorityScrollData.scrollPosition.x = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.x, 0), -(highestPriorityScrollData.contentSize.width - highestPriorityScrollData.boundingBox.width));
	                       scrollDeltaX = highestPriorityScrollData.scrollPosition.x - oldXScrollPosition;
	                   }
	                   if (canScrollVertically) {
	                       float oldYScrollPosition = highestPriorityScrollData.scrollPosition.y;
	                       highestPriorityScrollData.scrollPosition.y = highestPriorityScrollData.scrollOrigin.y + (context.pointerInfo.position.y - highestPriorityScrollData.pointerOrigin.y);
	                       highestPriorityScrollData.scrollPosition.y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.y, 0), -(highestPriorityScrollData.contentSize.height - highestPriorityScrollData.boundingBox.height));
	                       scrollDeltaY = highestPriorityScrollData.scrollPosition.y - oldYScrollPosition;
	                   }
	                   if (scrollDeltaX > -0.1f && scrollDeltaX < 0.1f && scrollDeltaY > -0.1f && scrollDeltaY < 0.1f && highestPriorityScrollData.momentumTime > 0.15f) {
	                       highestPriorityScrollData.momentumTime = 0;
	                       highestPriorityScrollData.pointerOrigin = context.pointerInfo.position;
	                       highestPriorityScrollData.scrollOrigin = highestPriorityScrollData.scrollPosition;
	                   } else {
	                        highestPriorityScrollData.momentumTime += deltaTime;
	                   }
	               }
	           }
	           // Clamp any changes to scroll position to the maximum size of the contents
	           if (canScrollVertically) {
	               highestPriorityScrollData.scrollPosition.y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.y, 0), -(highestPriorityScrollData.contentSize.height - scrollElement.dimensions.height));
	           }
	           if (canScrollHorizontally) {
	               highestPriorityScrollData.scrollPosition.x = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.x, 0), -(highestPriorityScrollData.contentSize.width - scrollElement.dimensions.width));
	           }
	       }
	*/
}

// Updates the layout dimensions in response to the window or outer container being resized.
func (c *Context) SetLayoutDimensions(dimensions Dimensions) {
	c.layoutDimensions = dimensions
}

// Called before starting any layout declarations.
func (c *Context) BeginLayout() {
	c.initializeEphemeralMemory()
	c.generation++
	c.dynamicElementIndex = 0
	// Set up the root container that covers the entire window
	rootDimensions := c.layoutDimensions
	if c.debugModeEnabled {
		rootDimensions.X -= (float32)(debugViewWidth)
	}
	c.booleanWarnings = BooleanWarnings{}
	c.openElement()
	c.configureOpenElement(&ElementDeclaration{
		Id: c.ID("Clay__RootContainer"),
		Layout: LayoutConfig{
			Sizing: Sizing{
				c.SIZING_FIXED(rootDimensions.X),
				c.SIZING_FIXED(rootDimensions.Y),
			},
		},
	})
	c.openLayoutElementStack = append(c.openLayoutElementStack, 0)
	c.layoutElementTreeRoots = append(c.layoutElementTreeRoots, LayoutElementTreeRoot{layoutElementIndex: 0})
}

// Called when all layout declarations are finished.
// Computes the layout and generates and returns the array of render commands to draw.
func (c *Context) EndLayout() []RenderCommand {
	c.closeElement()
	elementsExceededBeforeDebugView := c.booleanWarnings.maxElementsExceeded
	if c.debugModeEnabled && !elementsExceededBeforeDebugView {
		c.warningsEnabled = false
		c.Clay__RenderDebugView()
		c.warningsEnabled = true
	}
	if c.booleanWarnings.maxElementsExceeded {
		message := ""
		if !elementsExceededBeforeDebugView {
			message = "Clay Error: Layout elements exceeded Clay__maxElementCount after adding the debug-view to the layout."
		} else {
			message = "Clay Error: Layout elements exceeded Clay__maxElementCount"
		}
		c.addRenderCommand(RenderCommand{
			BoundingBox: MakeBoundingBox(c.layoutDimensions.ScaleF(0.5).AddX(-59*4), vector2.Zero[float32]()),
			RenderData: TextRenderData{
				StringContents: message,
				TextColor:      Color{R: 255, G: 0, B: 0, A: 255},
				FontSize:       16,
			},
		})
	} else {
		c.calculateFinalLayout()
	}

	return c.renderCommands
}

// Returns layout data such as the final calculated bounding box for an element with a given ID.
// The returned clay.ElementData contains a `found` bool that will be true if an element with the provided ID was found.
// This ID can be calculated either with CLAY_ID() for string literal IDs, or clay.GetElementId for dynamic strings.
///CLAY_DLL_EXPORT clay.ElementData clay.GetElementData(clay.ElementId id);

// Returns true if the pointer position provided by clay.SetPointerState is within the current element's bounding box.
// Works during element declaration, e.g. CLAY({ .backgroundColor = clay.Hovered() ? BLUE : RED });
///CLAY_DLL_EXPORT bool clay.Hovered(void);

// Bind a callback that will be called when the pointer position provided by clay.SetPointerState is within the current element's bounding box.
// - onHoverFunction is a function pointer to a user defined function.
// - userData is a pointer that will be transparently passed through when the onHoverFunction is called.
///CLAY_DLL_EXPORT void clay.OnHover(void (*onHoverFunction)(clay.ElementId elementId, clay.PointerData pointerData, intptr_t userData), intptr_t userData);

// An imperative function that returns true if the pointer position provided by clay.SetPointerState is within the element with the provided ID's bounding box.
// This ID can be calculated either with CLAY_ID() for string literal IDs, or clay.GetElementId for dynamic strings.
///CLAY_DLL_EXPORT bool clay.PointerOver(clay.ElementId elementId);

// Returns data representing the state of the scrolling element with the provided ID.
// The returned clay.ScrollContainerData contains a `found` bool that will be true if a scroll element was found with the provided ID.
// An imperative function that returns true if the pointer position provided by clay.SetPointerState is within the element with the provided ID's bounding box.
// This ID can be calculated either with CLAY_ID() for string literal IDs, or clay.GetElementId for dynamic strings.
///CLAY_DLL_EXPORT clay.ScrollContainerData clay.GetScrollContainerData(clay.ElementId id);

// Binds a callback function that Clay will call to determine the dimensions of a given string slice.
// - measureTextFunction is a user provided function that adheres to the interface clay.Dimensions (clay.StringSlice text, clay.TextElementConfig *config, void *userData);
// - userData is a pointer that will be transparently passed through when the measureTextFunction is called.
func (c *Context) SetMeasureTextFunction(fn MeasureTextFn, userData any) {
	measureText = fn
	c.measureTextUserData = userData
}

// Experimental - Used in cases where Clay needs to integrate with a system that manages its own scrolling containers externally.
// Please reach out if you plan to use this function, as it may be subject to change.
func (c *Context) SetQueryScrollOffsetFunction(fn QueryScrollOffsetFn, userData any) {
	queryScrollOffset = fn
	c.queryScrollOffsetUserData = userData
}

// Enables and disables Clay's internal debug tools.
// This state is retained and does not need to be set each frame.
///CLAY_DLL_EXPORT void SetDebugModeEnabled(bool enabled);

// Returns true if Clay's internal debug tools are currently enabled.
///CLAY_DLL_EXPORT bool IsDebugModeEnabled(void);

// Enables and disables visibility culling. By default, Clay will not generate render commands for elements whose bounding box is entirely outside the screen.
///CLAY_DLL_EXPORT void SetCullingEnabled(bool enabled);

func (c *Context) SetRenderTranclucentEnabled(v bool) {
	c.renderTranslucent = v
}

// Returns the maximum number of UI elements supported by Clay's current configuration.
///CLAY_DLL_EXPORT int32_t GetMaxElementCount(void);

// Modifies the maximum number of UI elements supported by Clay's current configuration.
// This may require reallocating additional memory, and re-calling clay.Initialize();
///CLAY_DLL_EXPORT void SetMaxElementCount(int32_t maxElementCount);

// Returns the maximum number of measured "words" (whitespace seperated runs of characters) that Clay can store in its internal text measurement cache.
///CLAY_DLL_EXPORT int32_t GetMaxMeasureTextCacheWordCount(void);

// Modifies the maximum number of measured "words" (whitespace seperated runs of characters) that Clay can store in its internal text measurement cache.
// This may require reallocating additional memory, and re-calling clay.Initialize();
///CLAY_DLL_EXPORT void SetMaxMeasureTextCacheWordCount(int32_t maxMeasureTextCacheWordCount);

// Resets Clay's internal text measurement cache, useful if memory to represent strings is being re-used.
// Similar behaviour can be achieved on an individual text element level by using TextElementConfig.hashStringContents
///CLAY_DLL_EXPORT void ResetMeasureTextCache(void);

func (c *Context) CLAY(e ElementDeclaration, fns ...func()) {
	if !c.openElement() {
		return
	}

	c.configureOpenElement(&e)
	for _, fn := range fns {
		fn()
	}
	c.closeElement()
}

func (c *Context) TEXT_CONFIG(config TextElementConfig) *TextElementConfig {
	return c.storeTextElementConfig(config)
}

func (c *Context) TEXT(text string, config *TextElementConfig) {
	c.openTextElement(text, config)
}
