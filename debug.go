package goclay

import "github.com/igadmg/goex/image/colorex"

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
	                                   CLAY_TEXT(Clay__IntToString(borderConfig.width.bottom), infoTextConfig);
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
