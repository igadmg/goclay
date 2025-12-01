package clay

import (
	"math"
	"strconv"
)

var CLAY__DEBUGVIEW_COLOR_1 Color = Color{R: 58, G: 56, B: 52, A: 255}
var CLAY__DEBUGVIEW_COLOR_2 Color = Color{R: 62, G: 60, B: 58, A: 255}
var CLAY__DEBUGVIEW_COLOR_3 Color = Color{R: 141, G: 133, B: 135, A: 255}
var CLAY__DEBUGVIEW_COLOR_4 Color = Color{R: 238, G: 226, B: 231, A: 255}
var CLAY__DEBUGVIEW_COLOR_SELECTED_ROW Color = Color{R: 102, G: 80, B: 78, A: 255}
var CLAY__DEBUGVIEW_ROW_HEIGHT float32 = 30
var CLAY__DEBUGVIEW_OUTER_PADDING uint16 = 10
var CLAY__DEBUGVIEW_INDENT_WIDTH int32 = 16
var Clay__DebugView_TextNameConfig TextElementConfig = TextElementConfig{
	TextColor: Color{R: 238, G: 226, B: 231, A: 255},
	FontSize:  16,
	WrapMode:  TEXT_WRAP_NONE,
}
var Clay__DebugView_ScrollViewItemLayoutConfig LayoutConfig
var debugViewWidth uint32 = 400
var debugViewHighlightColor = Color{R: 168, G: 66, B: 28, A: 100}

/*
#pragma region DebugTools

	typedef struct {
	    string label;
	    Color color;
	} Clay__DebugElementConfigTypeLabelConfig;

	Clay__DebugElementConfigTypeLabelConfig Clay__DebugGetElementConfigTypeLabel(Clay__ElementConfigType type) {
	    switch (type) {
	        case CLAY__ELEMENT_CONFIG_TYPE_SHARED: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { "Shared", {243,134,48,255,
			} };
	        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { "Text", {105,210,231,255,
			} };
	        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { "Image", {121,189,154,255,
			} };
	        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { "Floating", {250,105,0,255,
			} };
	        case CLAY__ELEMENT_CONFIG_TYPE_CLIP: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) {"Scroll", {242, 196, 90, 255,
			} };
	        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) {"Border", {108, 91, 123, 255,
			} };
	        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { "Custom", {11,72,107,255,
			} };
	        default: break;
	    }
	    return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { "Error", {0,0,0,255,
		} };
	}
*/
type RenderDebugLayoutData struct {
	rowCount                int32
	selectedElementRowIndex int32
}

// Returns row count

func (c *Context) renderDebugLayoutElementsList(initialRootsLength int32, highlightedRowIndex int32) RenderDebugLayoutData {
	/*
		Clay__int32_tArray dfsBuffer = context.reusableElementIndexBuffer;
		Clay__DebugView_ScrollViewItemLayoutConfig = CLAY__INIT(Clay_LayoutConfig) {
		Sizing: Sizing{ Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT) },
		ChildGap: 6,
		ChildAlignment: ChildAlignment { Y: ALIGN_Y_CENTER }};
	*/
	var layoutData RenderDebugLayoutData
	/*
		uint32 highlightedElementId = 0;

		for (int32 rootIndex = 0; rootIndex < initialRootsLength; ++rootIndex) {
			dfsBuffer.length = 0;
			Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex);
			Clay__int32_tArray_Add(&dfsBuffer, (int32)root.layoutElementIndex);
			context.treeNodeVisited.internalArray[0] = false;
			if (rootIndex > 0) {
				CLAY({ .id = c.IDI("Clay__DebugView_EmptyRowOuter", rootIndex),
				Layout: LayoutConfig{
				Sizing: Sizing{Width: c.SIZING_GROW(0),
				},
					Padding: Padding{CLAY__DEBUGVIEW_INDENT_WIDTH / 2, 0, 0, 0,
					} } }) {
					CLAY({ .id = c.IDI("Clay__DebugView_EmptyRow", rootIndex),
					Layout: LayoutConfig{
					Sizing: Sizing{ Width: c.SIZING_GROW(0), Height: c.SIZING_FIXED((float32)CLAY__DEBUGVIEW_ROW_HEIGHT) },
					}, Border: BorderElementConfig { color: CLAY__DEBUGVIEW_COLOR_3,
					width: BorderWidth { .top = 1 } } }) {}
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
				CLAY({ .id = c.IDI("Clay__DebugView_ElementOuter", currentElement.id),
				Layout: Clay__DebugView_ScrollViewItemLayoutConfig }) {
					// Collapse icon / button
					if (!(elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement.childrenOrTextContent.children.length == 0)) {
						CLAY({
							.id = c.IDI("Clay__DebugView_CollapseElement", currentElement.id),
							Layout: LayoutConfig {
							Sizing: Sizing{
							Width: c.SIZING_FIXED(16),
							Height: c.SIZING_FIXED(16),
							},
							ChildAlignment: ChildAlignment { ALIGN_X_CENTER, ALIGN_Y_CENTER,
							} },
							CornerRadius: CORNER_RADIUS(4),
							Border: BorderElementConfig { color: CLAY__DEBUGVIEW_COLOR_3,
							width: BorderWidth {1, 1, 1, 1, 0,
							} },
						}) {
							c.TEXT((currentElementData && currentElementData.debugData.collapsed) ? "+" : "-", c.TEXT_CONFIG(TextElementConfig {
							textColor: CLAY__DEBUGVIEW_COLOR_4,
							fontSize: 16 }));
						}
					} else { // Square dot for empty containers
						CLAY({
						Layout: LayoutConfig{
						Sizing: Sizing{
						Width: c.SIZING_FIXED(16),
						Height: c.SIZING_FIXED(16),
						},
						ChildAlignment: ChildAlignment { ALIGN_X_CENTER, ALIGN_Y_CENTER } } }) {
							CLAY({
							Layout: LayoutConfig{
							Sizing: Sizing{
							Width: c.SIZING_FIXED(8),
							Height: c.SIZING_FIXED(8),
							} },
							BackgroundColor: CLAY__DEBUGVIEW_COLOR_3, CornerRadius: CORNER_RADIUS(2) }) {}
						}
					}
					// Collisions and offscreen info
					if (currentElementData) {
						if (currentElementData.debugData.collision) {
							CLAY({
							Layout: LayoutConfig{ Padding: { 8, 8, 2, 2 },
							}, Border: BorderElementConfig { color: Color{177, 147, 8, 255,
							},
							width: BorderWidth {1, 1, 1, 1, 0,
							} } }) {
								c.TEXT("Duplicate ID", c.TEXT_CONFIG(TextElementConfig {
								textColor: CLAY__DEBUGVIEW_COLOR_3,
								fontSize: 16 }));
							}
						}
						if (offscreen) {
							CLAY({
							Layout: LayoutConfig{ Padding: { 8, 8, 2, 2 } }, Border: BorderElementConfig {  .color = CLAY__DEBUGVIEW_COLOR_3,
							width: BorderWidth { 1, 1, 1, 1, 0,
							} } }) {
								c.TEXT("Offscreen", c.TEXT_CONFIG(TextElementConfig {
								textColor: CLAY__DEBUGVIEW_COLOR_3,
								fontSize: 16 }));
							}
						}
					}
					string idString = context.layoutElementIdStrings.internalArray[currentElementIndex];
					if (idString.length > 0) {
						c.TEXT(idString, offscreen ? c.TEXT_CONFIG(TextElementConfig {
						textColor: CLAY__DEBUGVIEW_COLOR_3,
						fontSize: 16 }) : &Clay__DebugView_TextNameConfig);
					}
					for (int32 elementConfigIndex = 0; elementConfigIndex < currentElement.elementConfigs.length; ++elementConfigIndex) {
						Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, elementConfigIndex);
						if (elementConfig.type == CLAY__ELEMENT_CONFIG_TYPE_SHARED) {
							Color labelColor = {243,134,48,90,
							};
							labelColor.a = 90;
							Color backgroundColor = elementConfig.config.sharedElementConfig.backgroundColor;
							Clay_CornerRadius radius = elementConfig.config.sharedElementConfig.cornerRadius;
							if (backgroundColor.a > 0) {
								CLAY({
								Layout: LayoutConfig{ Padding: { 8, 8, 2, 2 } },
								BackgroundColor: labelColor, CornerRadius: CORNER_RADIUS(4), Border: BorderElementConfig { color: labelColor,
								width: BorderWidth { 1, 1, 1, 1, 0,
								} } }) {
									c.TEXT("Color", c.TEXT_CONFIG(TextElementConfig {
									textColor: offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4,
									fontSize: 16 }));
								}
							}
							if (radius.bottomLeft > 0) {
								CLAY({
								Layout: LayoutConfig{ Padding: { 8, 8, 2, 2 } },
								BackgroundColor: labelColor, CornerRadius: CORNER_RADIUS(4), Border: BorderElementConfig { color: labelColor,
								width: BorderWidth { 1, 1, 1, 1, 0 } } }) {
									c.TEXT("Radius", c.TEXT_CONFIG(TextElementConfig {
									textColor: offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4,
									fontSize: 16 }));
								}
							}
							continue;
						}
						Clay__DebugElementConfigTypeLabelConfig config = Clay__DebugGetElementConfigTypeLabel(elementConfig.type);
						Color backgroundColor = config.color;
						backgroundColor.a = 90;
						CLAY({
						Layout: LayoutConfig{ Padding: { 8, 8, 2, 2 } },
						BackgroundColor: backgroundColor, CornerRadius: CORNER_RADIUS(4), Border: BorderElementConfig { color: config.color,
						width: BorderWidth { 1, 1, 1, 1, 0 } } }) {
							c.TEXT(config.label, c.TEXT_CONFIG(TextElementConfig {
							textColor: offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4,
							fontSize: 16 }));
						}
					}
				}

				// Render the text contents below the element as a non-interactive row
				if (elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
					layoutData.rowCount++;
					Clay__TextElementData *textElementData = currentElement.childrenOrTextContent.textElementData;
					Clay_TextElementConfig *rawTextConfig = offscreen ? c.TEXT_CONFIG(TextElementConfig {
					textColor: CLAY__DEBUGVIEW_COLOR_3,
					fontSize: 16 }) : &Clay__DebugView_TextNameConfig;
					CLAY({
					Layout: LayoutConfig{
					Sizing: Sizing{ Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT),
					},
					ChildAlignment: ChildAlignment { Y: ALIGN_Y_CENTER } } }) {
						CLAY({
						Layout: LayoutConfig{
						Sizing: Sizing{Width: c.SIZING_FIXED(CLAY__DEBUGVIEW_INDENT_WIDTH + 16) } } }) {}
						c.TEXT(CLAY_STRING("\""), rawTextConfig);
						c.TEXT(textElementData.text.length > 40 ? (CLAY__INIT(string) { .length = 40, .chars = textElementData.text.chars }) : textElementData.text, rawTextConfig);
						if (textElementData.text.length > 40) {
							c.TEXT("...", rawTextConfig);
						}
						c.TEXT(CLAY_STRING("\""), rawTextConfig);
					}
				} else if (currentElement.childrenOrTextContent.children.length > 0) {
					Clay__OpenElement();
					Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) {
					Layout: LayoutConfig{ Padding: { .left = 8 } } });
					Clay__OpenElement();
					Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) {
					Layout: LayoutConfig{ Padding: { .left = CLAY__DEBUGVIEW_INDENT_WIDTH },
					}, Border: BorderElementConfig { color: CLAY__DEBUGVIEW_COLOR_3,
					width: BorderWidth { .left = 1 } },
					});
					Clay__OpenElement();
					Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) {
					Layout: LayoutConfig{ LayoutDirection: TOP_TO_BOTTOM } });
				}

				layoutData.rowCount++;
				if (!(elementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (currentElementData && currentElementData.debugData.collapsed))) {
					for (int32 i = currentElement.childrenOrTextContent.children.length - 1; i >= 0; --i) {
						Clay__int32_tArray_Add(&dfsBuffer, currentElement.childrenOrTextContent.children.elements[i]);
						context.treeNodeVisited.internalArray[dfsBuffer.length - 1] = false; // TODO: needs to be ranged checked
					}
				}
			}
		}

		if (context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
			Clay_ElementId collapseButtonId = hashString("Clay__DebugView_CollapseElement");
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
			CLAY({ .id = c.ID("Clay__DebugView_ElementHighlight"),
			Layout: LayoutConfig{
			Sizing: Sizing{
			Width: c.SIZING_GROW(0),
			Height: c.SIZING_GROW(0),
			} }, Floating: FloatingElementConfig{ .parentId = highlightedElementId, .zIndex = 32767,
				pointerCaptureMode: POINTER_CAPTURE_MODE_PASSTHROUGH,
			attachTo: CLAY_ATTACH_TO_ELEMENT_WITH_ID } }) {
				CLAY({ .id = c.ID("Clay__DebugView_ElementHighlightRectangle"),
				Layout: LayoutConfig{
				Sizing: Sizing{
				Width: c.SIZING_GROW(0),
				Height: c.SIZING_GROW(0),
				} },
				BackgroundColor: debugViewHighlightColor }) {}
			}
		}
	*/
	return layoutData
}

func (c *Context) renderDebugLayoutSizing(sizing AnySizingAxis, infoTextConfig *TextElementConfig) {
	c.CLAY_TEXT(SizingAxisTypeString(sizing), infoTextConfig)
	switch mm := sizing.(type) {
	case SizingAxisMinMax:
		c.CLAY_TEXT("(", infoTextConfig)
		if mm.GetMinMax().Min != 0 {
			c.CLAY_TEXT("min: ", infoTextConfig)
			c.CLAY_TEXT(strconv.Itoa(int(mm.GetMinMax().Min)), infoTextConfig)
			if mm.GetMinMax().Max != math.MaxFloat32 {
				c.CLAY_TEXT(", ", infoTextConfig)
			}
		}
		if mm.GetMinMax().Max != math.MaxFloat32 {
			c.CLAY_TEXT("max: ", infoTextConfig)
			c.CLAY_TEXT(strconv.Itoa(int(mm.GetMinMax().Max)), infoTextConfig)
		}
		c.CLAY_TEXT(")", infoTextConfig)
	}
}

func (c *Context) Clay__RenderDebugViewElementConfigHeader(elementId string, config AnyElementConfig) {
	/*
		Clay__DebugElementConfigTypeLabelConfig config = Clay__DebugGetElementConfigTypeLabel(type);
			    Color backgroundColor = config.color;
			    backgroundColor.a = 90;
			    CLAY({
				Layout: LayoutConfig{
				Sizing: Sizing{ Width: c.SIZING_GROW(0) },
				 Padding: PaddingCLAY_PADDING_ALL(CLAY__DEBUGVIEW_OUTER_PADDING),
				 ChildAlignment: ChildAlignment { Y: ALIGN_Y_CENTER } } }) {
			        CLAY({
					Layout: LayoutConfig{ Padding: { 8, 8, 2, 2 } },
					BackgroundColor: backgroundColor, CornerRadius: CORNER_RADIUS(4), Border: BorderElementConfig { color: config.color,
					width: BorderWidth { 1, 1, 1, 1, 0 } } }) {
			            c.TEXT(config.label, c.TEXT_CONFIG(TextElementConfig {
						textColor: CLAY__DEBUGVIEW_COLOR_4,
						fontSize: 16 }));
			        }
			        CLAY({
					Layout: LayoutConfig{
					Sizing: Sizing{ Width: c.SIZING_GROW(0) } } }) {}
			        c.TEXT(elementId, c.TEXT_CONFIG(TextElementConfig {
					textColor: CLAY__DEBUGVIEW_COLOR_3,
					fontSize: 16,
					wrapMode: TEXT_WRAP_NONE }));
			    }
	*/
}

func (c *Context) Clay__RenderDebugViewColor(color Color, textConfig *TextElementConfig) {
	/*
		    CLAY({
			Layout: LayoutConfig{ .childAlignment = {Y: ALIGN_Y_CENTER,
			} } }) {
		        c.TEXT("{ r: ", textConfig);
		        c.TEXT(strconv.Itoa(color.r), textConfig);
		        c.TEXT(", g: ", textConfig);
		        c.TEXT(strconv.Itoa(color.g), textConfig);
		        c.TEXT(", b: ", textConfig);
		        c.TEXT(strconv.Itoa(color.b), textConfig);
		        c.TEXT(", a: ", textConfig);
		        c.TEXT(strconv.Itoa(color.a), textConfig);
		        c.TEXT(" }", textConfig);
		        CLAY({
				Layout: LayoutConfig{
				Sizing: Sizing{ Width: c.SIZING_FIXED(10) } } }) {}
		        CLAY({
				Layout: LayoutConfig{
				Sizing: Sizing{ c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 8),
				Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 8),
				} },
				BackgroundColor: color, CornerRadius: CORNER_RADIUS(4), Border: BorderElementConfig { color: CLAY__DEBUGVIEW_COLOR_4,
				width: BorderWidth { 1, 1, 1, 1, 0 } } }) {}
		    }
	*/
}

func (c *Context) renderDebugViewCornerRadius(cornerRadius CornerRadius, textConfig *TextElementConfig) {
	/*
		    CLAY({
			Layout: LayoutConfig{ .childAlignment = {Y: ALIGN_Y_CENTER,
			} } }) {
		        c.TEXT("{ topLeft: ", textConfig);
		        c.TEXT(strconv.Itoa(cornerRadius.topLeft), textConfig);
		        c.TEXT(", topRight: ", textConfig);
		        c.TEXT(strconv.Itoa(cornerRadius.topRight), textConfig);
		        c.TEXT(", bottomLeft: ", textConfig);
		        c.TEXT(strconv.Itoa(cornerRadius.bottomLeft), textConfig);
		        c.TEXT(", bottomRight: ", textConfig);
		        c.TEXT(strconv.Itoa(cornerRadius.bottomRight), textConfig);
		        c.TEXT(" }", textConfig);
		    }
	*/
}

/*
	void HandleDebugViewCloseButtonInteraction(Clay_ElementId elementId, Clay_PointerData pointerInfo, intptr_t userData) {
	    context := GetCurrentContext();
	    (void) elementId; (void) pointerInfo; (void) userData;
	    if (pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
	        context.debugModeEnabled = false;
	    }
	}
*/
func (c *Context) Clay__RenderDebugView() {
	closeButtonId := hashString("Clay__DebugViewTopHeaderCloseButtonOuter")
	if c.pointerInfo.State == POINTER_DATA_PRESSED_THIS_FRAME {
		for _, elementId := range c.pointerOverIds {
			if elementId.id == closeButtonId.id {
				c.debugModeEnabled = false
				return
			}
		}
	}

	initialRootsLength := int32(len(c.layoutElementTreeRoots))
	initialElementsLength := int32(len(c.layoutElements))
	infoTextConfig := c.TEXT_CONFIG(TextElementConfig{
		TextColor: CLAY__DEBUGVIEW_COLOR_4,
		FontSize:  16,
		WrapMode:  TEXT_WRAP_NONE,
	})
	infoTitleConfig := c.TEXT_CONFIG(TextElementConfig{
		TextColor: CLAY__DEBUGVIEW_COLOR_3,
		FontSize:  16,
		WrapMode:  TEXT_WRAP_NONE,
	})
	scrollId := hashString("Clay__DebugViewOuterScrollPane")
	scrollYOffset := float32(0)
	pointerInDebugView := c.pointerInfo.Position.Y < c.layoutDimensions.Y-300
	for _, scrollContainerData := range c.scrollContainerDatas {
		if scrollContainerData.elementId == scrollId.id {
			if !c.externalScrollHandlingEnabled {
				scrollYOffset = scrollContainerData.scrollPosition.Y
			} else {
				pointerInDebugView = c.pointerInfo.Position.Y+scrollContainerData.scrollPosition.Y < c.layoutDimensions.Y-300
			}
			break
		}
	}
	highlightedRow := int32(-1)
	if pointerInDebugView {
		highlightedRow = (int32)((c.pointerInfo.Position.Y-scrollYOffset)/float32(CLAY__DEBUGVIEW_ROW_HEIGHT)) - 1
	}
	if c.pointerInfo.Position.X < c.layoutDimensions.X-float32(debugViewWidth) {
		highlightedRow = -1
	}
	var layoutData RenderDebugLayoutData
	c.CLAY_ID(c.ID("Clay__DebugView"), ElementDeclaration{
		Layout: LayoutConfig{
			Sizing: Sizing{
				Width:  c.SIZING_FIXED(float32(debugViewWidth)),
				Height: c.SIZING_FIXED(c.layoutDimensions.Y),
			},
			LayoutDirection: TOP_TO_BOTTOM,
		},
		Floating: FloatingElementConfig{
			ZIndex: 32765,
			AttachPoints: FloatingAttachPoints{
				Element: ATTACH_POINT_LEFT_CENTER,
				Parent:  ATTACH_POINT_RIGHT_CENTER,
			},
			AttachTo: ATTACH_TO_ROOT,
		},
		Border: BorderElementConfig{
			Color: CLAY__DEBUGVIEW_COLOR_3,
			Width: BorderWidth{Bottom: 1},
		},
	}, func() {
		c.CLAY(ElementDeclaration{
			Layout: LayoutConfig{
				Sizing: Sizing{
					Width:  c.SIZING_GROW(0),
					Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT),
				},
				Padding:        Padding{CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0},
				ChildAlignment: ChildAlignment{Y: ALIGN_Y_CENTER},
			},
			BackgroundColor: CLAY__DEBUGVIEW_COLOR_2,
		}, func() {
			c.CLAY_TEXT("Clay Debug Tools", infoTextConfig)
			c.CLAY(ElementDeclaration{
				Layout: LayoutConfig{
					Sizing: Sizing{Width: c.SIZING_GROW(0)},
				},
			})
			// Close button
			c.CLAY(ElementDeclaration{
				Layout: LayoutConfig{
					Sizing: Sizing{
						Width:  c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10),
						Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10),
					},
					ChildAlignment: ChildAlignment{ALIGN_X_CENTER, ALIGN_Y_CENTER},
				},
				BackgroundColor: Color{R: 217, G: 91, B: 67, A: 80},
				CornerRadius:    CORNER_RADIUS(4),
				Border: BorderElementConfig{
					Color: Color{R: 217, G: 91, B: 67, A: 255},
					Width: BorderWidth{1, 1, 1, 1, 0},
				},
			}, func() {
				//Clay_OnHover(HandleDebugViewCloseButtonInteraction, 0);
				c.CLAY_TEXT("x", c.TEXT_CONFIG(TextElementConfig{
					TextColor: CLAY__DEBUGVIEW_COLOR_4,
					FontSize:  16,
				}))
			})
		})
		c.CLAY(ElementDeclaration{
			Layout: LayoutConfig{
				Sizing: Sizing{
					Width:  c.SIZING_GROW(0),
					Height: c.SIZING_FIXED(1),
				},
			},
			BackgroundColor: CLAY__DEBUGVIEW_COLOR_3,
		})
		c.CLAY_ID(scrollId, ElementDeclaration{
			Layout: LayoutConfig{
				Sizing: Sizing{
					Width:  c.SIZING_GROW(0),
					Height: c.SIZING_GROW(0),
				},
			},
			Clip: ClipElementConfig{
				Horizontal:  true,
				Vertical:    true,
				ChildOffset: c.GetScrollOffset(),
			},
		}, func() {
			bgColor := CLAY__DEBUGVIEW_COLOR_2
			if ((initialElementsLength + initialRootsLength) & 1) != 0 {
				bgColor = CLAY__DEBUGVIEW_COLOR_1
			}
			c.CLAY(ElementDeclaration{
				Layout: LayoutConfig{
					Sizing: Sizing{
						Width:  c.SIZING_GROW(0),
						Height: c.SIZING_GROW(0),
					},
					LayoutDirection: TOP_TO_BOTTOM,
				},
				BackgroundColor: bgColor,
			}, func() {
				panelContentsId := hashString("Clay__DebugViewPaneOuter")
				// Element list
				c.CLAY_ID(panelContentsId, ElementDeclaration{
					Layout: LayoutConfig{
						Sizing: Sizing{
							Width:  c.SIZING_GROW(0),
							Height: c.SIZING_GROW(0),
						},
					},
					Floating: FloatingElementConfig{
						ZIndex:             32766,
						PointerCaptureMode: POINTER_CAPTURE_MODE_PASSTHROUGH,
						AttachTo:           ATTACH_TO_PARENT,
					},
				}, func() {
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{
							Sizing: Sizing{
								Width:  c.SIZING_GROW(0),
								Height: c.SIZING_GROW(0),
							},
							Padding:         Padding{CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0},
							LayoutDirection: TOP_TO_BOTTOM,
						},
					}, func() {
						layoutData = c.renderDebugLayoutElementsList(initialRootsLength, highlightedRow)
					})
				})
				panelContents, _ := c.getHashMapItem(panelContentsId.id)
				contentWidth := panelContents.layoutElement.dimensions.X
				c.CLAY(ElementDeclaration{
					Layout: LayoutConfig{
						Sizing:          Sizing{Width: c.SIZING_FIXED(contentWidth)},
						LayoutDirection: TOP_TO_BOTTOM,
					},
				})
				for i := int32(0); i < layoutData.rowCount; i++ {
					rowColor := CLAY__DEBUGVIEW_COLOR_1
					if (i & 1) == 0 {
						rowColor = CLAY__DEBUGVIEW_COLOR_2
					}
					if i == layoutData.selectedElementRowIndex {
						rowColor = CLAY__DEBUGVIEW_COLOR_SELECTED_ROW
					}
					if i == highlightedRow {
						rowColor.R *= 5
						rowColor.R /= 4
						rowColor.G *= 5
						rowColor.G /= 4
						rowColor.B *= 5
						rowColor.B /= 4
					}
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{
							Sizing: Sizing{
								Width:  c.SIZING_GROW(0),
								Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT),
							},
							LayoutDirection: TOP_TO_BOTTOM,
						},
						BackgroundColor: rowColor,
					})
				}
			})
		})

		c.CLAY(ElementDeclaration{
			Layout: LayoutConfig{
				Sizing: Sizing{
					Width:  c.SIZING_GROW(0),
					Height: c.SIZING_FIXED(1),
				},
			},
			BackgroundColor: CLAY__DEBUGVIEW_COLOR_3,
		})
		if c.debugSelectedElementId != 0 {
			selectedItem, _ := c.getHashMapItem(c.debugSelectedElementId)
			c.CLAY(ElementDeclaration{
				Layout: LayoutConfig{
					Sizing: Sizing{
						Width:  c.SIZING_GROW(0),
						Height: c.SIZING_FIXED(300),
					},
					LayoutDirection: TOP_TO_BOTTOM,
				},
				BackgroundColor: CLAY__DEBUGVIEW_COLOR_2,
				Clip:            ClipElementConfig{Vertical: true},
				Border: BorderElementConfig{
					Color: CLAY__DEBUGVIEW_COLOR_3,
					Width: BorderWidth{BetweenChildren: 1},
				},
			}, func() {
				c.CLAY(ElementDeclaration{
					Layout: LayoutConfig{
						Sizing: Sizing{
							Width:  c.SIZING_GROW(0),
							Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT + 8),
						},
						Padding:        Padding{CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0},
						ChildAlignment: ChildAlignment{Y: ALIGN_Y_CENTER},
					},
				}, func() {
					c.CLAY_TEXT("Layout Config", infoTextConfig)
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{
							Sizing: Sizing{Width: c.SIZING_GROW(0)},
						},
					})
					if selectedItem.elementId.stringId != "" {
						c.CLAY_TEXT(selectedItem.elementId.stringId, infoTitleConfig)
						/*
							if selectedItem.elementId.offset != 0 {
								c.TEXT(" (", infoTitleConfig)
								c.TEXT(strconv.Itoa(int(selectedItem.elementId.offset)), infoTitleConfig)
								c.TEXT(")", infoTitleConfig)
							}
						*/
					}
				})
				attributeConfigPadding := Padding{CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 8, 8}
				// Clay_LayoutConfig debug info
				c.CLAY(ElementDeclaration{
					Layout: LayoutConfig{
						Padding:         attributeConfigPadding,
						ChildGap:        8,
						LayoutDirection: TOP_TO_BOTTOM},
				}, func() {
					// .boundingBox
					c.CLAY_TEXT("Bounding Box", infoTitleConfig)
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
					}, func() {
						c.CLAY_TEXT("{ x: ", infoTextConfig)
						c.CLAY_TEXT(strconv.Itoa(int(selectedItem.boundingBox.X())), infoTextConfig)
						c.CLAY_TEXT(", y: ", infoTextConfig)
						c.CLAY_TEXT(strconv.Itoa(int(selectedItem.boundingBox.Y())), infoTextConfig)
						c.CLAY_TEXT(", width: ", infoTextConfig)
						c.CLAY_TEXT(strconv.Itoa(int(selectedItem.boundingBox.Width())), infoTextConfig)
						c.CLAY_TEXT(", height: ", infoTextConfig)
						c.CLAY_TEXT(strconv.Itoa(int(selectedItem.boundingBox.Height())), infoTextConfig)
						c.CLAY_TEXT(" }", infoTextConfig)
					})
					// .layoutDirection
					c.CLAY_TEXT("Layout Direction", infoTitleConfig)
					layoutConfig := selectedItem.layoutElement.layoutConfig
					c.CLAY_TEXT(layoutConfig.LayoutDirection.String(), infoTextConfig)
					// .sizing
					c.CLAY_TEXT("Sizing", infoTitleConfig)
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
					}, func() {
						c.CLAY_TEXT("width: ", infoTextConfig)
						c.renderDebugLayoutSizing(layoutConfig.Sizing.Width, infoTextConfig)
					})
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
					}, func() {
						c.CLAY_TEXT("height: ", infoTextConfig)
						c.renderDebugLayoutSizing(layoutConfig.Sizing.Height, infoTextConfig)
					})
					// .padding
					c.CLAY_TEXT("Padding", infoTitleConfig)
					c.CLAY_ID(c.ID("Clay__DebugViewElementInfoPadding"),
						ElementDeclaration{}, func() {
							c.CLAY_TEXT("{ left: ", infoTextConfig)
							c.CLAY_TEXT(strconv.Itoa(int(layoutConfig.Padding.Left)), infoTextConfig)
							c.CLAY_TEXT(", right: ", infoTextConfig)
							c.CLAY_TEXT(strconv.Itoa(int(layoutConfig.Padding.Right)), infoTextConfig)
							c.CLAY_TEXT(", top: ", infoTextConfig)
							c.CLAY_TEXT(strconv.Itoa(int(layoutConfig.Padding.Top)), infoTextConfig)
							c.CLAY_TEXT(", bottom: ", infoTextConfig)
							c.CLAY_TEXT(strconv.Itoa(int(layoutConfig.Padding.Bottom)), infoTextConfig)
							c.CLAY_TEXT(" }", infoTextConfig)
						})
					// .childGap
					c.CLAY_TEXT("Child Gap", infoTitleConfig)
					c.CLAY_TEXT(strconv.Itoa(int(layoutConfig.ChildGap)), infoTextConfig)
					// .childAlignment
					c.CLAY_TEXT("Child Alignment", infoTitleConfig)
					c.CLAY(ElementDeclaration{
						Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
					}, func() {
						c.CLAY_TEXT("{ x: ", infoTextConfig)
						alignX := "LEFT"
						switch layoutConfig.ChildAlignment.X {
						case ALIGN_X_CENTER:
							alignX = "CENTER"
						case ALIGN_X_RIGHT:
							alignX = "RIGHT"
						}
						c.CLAY_TEXT(alignX, infoTextConfig)
						c.CLAY_TEXT(", y: ", infoTextConfig)
						alignY := "TOP"
						switch layoutConfig.ChildAlignment.Y {
						case ALIGN_Y_CENTER:
							alignY = "CENTER"
						case ALIGN_Y_BOTTOM:
							alignY = "BOTTOM"
						}
						c.CLAY_TEXT(alignY, infoTextConfig)
						c.CLAY_TEXT(" }", infoTextConfig)
					})
				})
				for _, elementConfig := range selectedItem.layoutElement.elementConfigs {
					c.Clay__RenderDebugViewElementConfigHeader(selectedItem.elementId.stringId, elementConfig)
					switch cfg := elementConfig.(type) {
					case *SharedElementConfig:
						c.CLAY(ElementDeclaration{
							Layout: LayoutConfig{
								Padding:         attributeConfigPadding,
								ChildGap:        8,
								LayoutDirection: TOP_TO_BOTTOM,
							},
						}, func() {
							// .backgroundColor
							c.CLAY_TEXT("Background Color", infoTitleConfig)
							c.Clay__RenderDebugViewColor(cfg.backgroundColor, infoTextConfig)
							// .cornerRadius
							c.CLAY_TEXT("Corner Radius", infoTitleConfig)
							c.renderDebugViewCornerRadius(cfg.cornerRadius, infoTextConfig)
						})
					case *TextElementConfig:
						c.CLAY(ElementDeclaration{
							Layout: LayoutConfig{
								Padding:         attributeConfigPadding,
								ChildGap:        8,
								LayoutDirection: TOP_TO_BOTTOM,
							},
						}, func() {
							// .fontSize
							c.CLAY_TEXT("Font Size", infoTitleConfig)
							c.CLAY_TEXT(strconv.Itoa(int(cfg.FontSize)), infoTextConfig)
							// .fontId
							c.CLAY_TEXT("Font ID", infoTitleConfig)
							c.CLAY_TEXT(strconv.Itoa(int(cfg.FontId)), infoTextConfig)
							// .lineHeight
							c.CLAY_TEXT("Line Height", infoTitleConfig)
							lineHeight := "auto"
							if cfg.LineHeight != 0 {
								lineHeight = strconv.Itoa(int(cfg.LineHeight))
							}
							c.CLAY_TEXT(lineHeight, infoTextConfig)
							// .letterSpacing
							c.CLAY_TEXT("Letter Spacing", infoTitleConfig)
							c.CLAY_TEXT(strconv.Itoa(int(cfg.LetterSpacing)), infoTextConfig)
							// .wrapMode
							c.CLAY_TEXT("Wrap Mode", infoTitleConfig)
							c.CLAY_TEXT(cfg.WrapMode.String(), infoTextConfig)
							// .textAlignment
							c.CLAY_TEXT("Text Alignment", infoTitleConfig)
							c.CLAY_TEXT(cfg.TextAlignment.String(), infoTextConfig)
							// .textColor
							c.CLAY_TEXT("Text Color", infoTitleConfig)
							c.Clay__RenderDebugViewColor(cfg.TextColor, infoTextConfig)
						})
					case *ImageElementConfig:
						c.CLAY_ID(c.ID("Clay__DebugViewElementInfoImageBody"),
							ElementDeclaration{
								Layout: LayoutConfig{
									Padding:         attributeConfigPadding,
									ChildGap:        8,
									LayoutDirection: TOP_TO_BOTTOM,
								},
							}, func() {
								// .sourceDimensions
								c.CLAY_TEXT("Source Dimensions", infoTitleConfig)
								c.CLAY_ID(c.ID("Clay__DebugViewElementInfoImageDimensions"),
									ElementDeclaration{}, func() {
										c.CLAY_TEXT("{}", infoTextConfig)
									})
								// Image Preview
								c.CLAY_TEXT("Preview", infoTitleConfig)
								c.CLAY(ElementDeclaration{
									Layout: LayoutConfig{
										Sizing: Sizing{Width: c.SIZING_GROW()},
									},
									Image: *cfg,
								})
							})
					case *ClipElementConfig:
						c.CLAY(ElementDeclaration{
							Layout: LayoutConfig{
								Padding:         attributeConfigPadding,
								ChildGap:        8,
								LayoutDirection: TOP_TO_BOTTOM,
							},
						}, func() {
							// .vertical
							c.CLAY_TEXT("Vertical", infoTitleConfig)
							// TODO: fix
							//c.TEXT(scrollConfig.vertical ? "true" : "false" , infoTextConfig);
							// .horizontal
							c.CLAY_TEXT("Horizontal", infoTitleConfig)
							//c.TEXT(scrollConfig.horizontal ? "true" : "false" , infoTextConfig);
						})
					case *FloatingElementConfig:
						c.CLAY(ElementDeclaration{
							Layout: LayoutConfig{
								Padding:         attributeConfigPadding,
								ChildGap:        8,
								LayoutDirection: TOP_TO_BOTTOM,
							},
						}, func() {
							// .offset
							c.CLAY_TEXT("Offset", infoTitleConfig)
							c.CLAY(ElementDeclaration{
								Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
							}, func() {
								c.CLAY_TEXT("{ x: ", infoTextConfig)
								c.CLAY_TEXT(strconv.Itoa(int(cfg.Offset.X)), infoTextConfig)
								c.CLAY_TEXT(", y: ", infoTextConfig)
								c.CLAY_TEXT(strconv.Itoa(int(cfg.Offset.Y)), infoTextConfig)
								c.CLAY_TEXT(" }", infoTextConfig)
							})
							// .expand
							c.CLAY_TEXT("Expand", infoTitleConfig)
							c.CLAY(ElementDeclaration{
								Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
							}, func() {
								c.CLAY_TEXT("{ width: ", infoTextConfig)
								c.CLAY_TEXT(strconv.Itoa(int(cfg.Expand.X)), infoTextConfig)
								c.CLAY_TEXT(", height: ", infoTextConfig)
								c.CLAY_TEXT(strconv.Itoa(int(cfg.Expand.Y)), infoTextConfig)
								c.CLAY_TEXT(" }", infoTextConfig)
							})
							// .zIndex
							c.CLAY_TEXT("z-index", infoTitleConfig)
							c.CLAY_TEXT(strconv.Itoa(int(cfg.ZIndex)), infoTextConfig)
							// .parentId
							c.CLAY_TEXT("Parent", infoTitleConfig)
							hashItem, _ := c.getHashMapItem(cfg.ParentId)
							c.CLAY_TEXT(hashItem.elementId.stringId, infoTextConfig)
						})
					case *BorderElementConfig:
						c.CLAY_ID(c.ID("Clay__DebugViewElementInfoBorderBody"),
							ElementDeclaration{
								Layout: LayoutConfig{
									Padding:         attributeConfigPadding,
									ChildGap:        8,
									LayoutDirection: TOP_TO_BOTTOM,
								},
							}, func() {
								c.CLAY_TEXT("Border Widths", infoTitleConfig)
								c.CLAY(ElementDeclaration{
									Layout: LayoutConfig{LayoutDirection: LEFT_TO_RIGHT},
								}, func() {
									c.CLAY_TEXT("{ left: ", infoTextConfig)
									c.CLAY_TEXT(strconv.Itoa(int(cfg.Width.Left)), infoTextConfig)
									c.CLAY_TEXT(", right: ", infoTextConfig)
									c.CLAY_TEXT(strconv.Itoa(int(cfg.Width.Right)), infoTextConfig)
									c.CLAY_TEXT(", top: ", infoTextConfig)
									c.CLAY_TEXT(strconv.Itoa(int(cfg.Width.Top)), infoTextConfig)
									c.CLAY_TEXT(", bottom: ", infoTextConfig)
									c.CLAY_TEXT(strconv.Itoa(int(cfg.Width.Bottom)), infoTextConfig)
									c.CLAY_TEXT(" }", infoTextConfig)
								})
								// .textColor
								c.CLAY_TEXT("Border Color", infoTitleConfig)
								c.Clay__RenderDebugViewColor(cfg.Color, infoTextConfig)
							})
						break
					case *CustomElementConfig:
					default:
						break
					}
				}
			})
		} else {
			c.CLAY_ID(c.ID("Clay__DebugViewWarning nsScrollPane"),
				ElementDeclaration{
					Layout: LayoutConfig{
						Sizing: Sizing{
							Width:  c.SIZING_GROW(0),
							Height: c.SIZING_FIXED(300),
						},
						ChildGap:        6,
						LayoutDirection: TOP_TO_BOTTOM,
					},
					BackgroundColor: CLAY__DEBUGVIEW_COLOR_2,
					Clip: ClipElementConfig{
						Horizontal: true,
						Vertical:   true,
					},
				}, func() {
					warningConfig := c.TEXT_CONFIG(TextElementConfig{
						TextColor: CLAY__DEBUGVIEW_COLOR_4,
						FontSize:  16,
						WrapMode:  TEXT_WRAP_NONE,
					})
					c.CLAY_ID(c.ID("Clay__DebugViewWarningItemHeader"),
						ElementDeclaration{
							Layout: LayoutConfig{
								Sizing:         Sizing{Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)},
								Padding:        Padding{CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0},
								ChildGap:       8,
								ChildAlignment: ChildAlignment{Y: ALIGN_Y_CENTER},
							},
						}, func() {
							c.CLAY_TEXT("Warnings", warningConfig)
						})
					c.CLAY_ID(c.ID("Clay__DebugViewWarningsTopBorder"),
						ElementDeclaration{
							Layout: LayoutConfig{
								Sizing: Sizing{Width: c.SIZING_GROW(0), Height: c.SIZING_FIXED(1)},
							},
							BackgroundColor: Color{R: 200, G: 200, B: 200, A: 255},
						})
					previousWarningsLength := len(c.warnings)
					for i := 0; i < previousWarningsLength; i++ {
						warning := c.warnings[i]
						c.CLAY_ID(c.IDI("Clay__DebugViewWarningItem", uint32(i)),
							ElementDeclaration{
								Layout: LayoutConfig{
									Sizing:         Sizing{Height: c.SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)},
									Padding:        Padding{CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0},
									ChildGap:       8,
									ChildAlignment: ChildAlignment{Y: ALIGN_Y_CENTER},
								},
							}, func() {
								c.CLAY_TEXT(warning.baseMessage, warningConfig)
								if warning.dynamicMessage != "" {
									c.CLAY_TEXT(warning.dynamicMessage, warningConfig)
								}
							})
					}
				})
		}
	})
}
