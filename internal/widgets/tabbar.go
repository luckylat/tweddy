package widgets

import (
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"simple-editor/internal/style"
)

// TabBar represents the visual tab bar component
type TabBar struct {
	theme        *material.Theme
	tabManager   *TabManager
	tabButtons   map[int]*widget.Clickable
	closeButtons map[int]*widget.Clickable
	appTheme     style.Theme
}

// NewTabBar creates a new tab bar
func NewTabBar(theme *material.Theme, tabManager *TabManager) *TabBar {
	return &TabBar{
		theme:        theme,
		tabManager:   tabManager,
		tabButtons:   make(map[int]*widget.Clickable),
		closeButtons: make(map[int]*widget.Clickable),
		appTheme:     style.DefaultTheme(),
	}
}

// drawBottomBorder draws a bottom border line for the tab bar
func (tb *TabBar) drawBottomBorder(gtx layout.Context) {
	borderHeight := unit.Dp(1)
	
	defer op.Offset(image.Point{X: 0, Y: gtx.Dp(borderHeight)}).Push(gtx.Ops).Pop()
	
	rect := clip.Rect{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Dp(borderHeight),
		},
	}
	
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: tb.appTheme.Application.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// Layout renders the tab bar
func (tb *TabBar) Layout(gtx layout.Context) layout.Dimensions {
	tabs := tb.tabManager.GetTabs()
	activeIndex := tb.tabManager.GetActiveTabIndex()

	// Ensure we have buttons for all tabs
	for _, tab := range tabs {
		if _, exists := tb.tabButtons[tab.ID]; !exists {
			tb.tabButtons[tab.ID] = &widget.Clickable{}
			tb.closeButtons[tab.ID] = &widget.Clickable{}
		}
	}

	// Create tab buttons
	children := make([]layout.FlexChild, 0, len(tabs))

	for i, tab := range tabs {
		tab := tab // capture loop variable
		i := i     // capture loop variable

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return tb.layoutTab(gtx, tab, i == activeIndex)
		}))

		// Add spacer between tabs
		if i < len(tabs)-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Width: unit.Dp(2)}.Layout(gtx)
			}))
		}
	}

	// Add gray background to entire tab bar with border
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			// Gray background for entire tab bar
			paint.ColorOp{Color: tb.appTheme.TabBar.Background}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			// Add padding around the entire tab bar
			inset := layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
				Left:   unit.Dp(4),
				Right:  unit.Dp(4),
			}
			
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx, children...)
				})
			},
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			tb.drawBottomBorder(gtx)
			return layout.Dimensions{}
		}),
	)
}

// layoutTab renders a single tab
func (tb *TabBar) layoutTab(gtx layout.Context, tab *Tab, isActive bool) layout.Dimensions {
	tabButton := tb.tabButtons[tab.ID]
	closeButton := tb.closeButtons[tab.ID]

	// Handle tab click
	if tabButton.Clicked(gtx) {
		tb.tabManager.SetActiveTabByID(tab.ID)
	}

	// Handle close button click
	if closeButton.Clicked(gtx) {
		tb.tabManager.CloseTabByID(tab.ID)
	}

	// Apply inset for better padding
	inset := layout.Inset{
		Top:    unit.Dp(4),
		Bottom: unit.Dp(4),
		Left:   unit.Dp(8),
		Right:  unit.Dp(4),
	}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Use Background widget for better visual separation
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				// Different background color for active tab
				var bgColor color.NRGBA
				if isActive {
					bgColor = tb.appTheme.TabBar.ActiveTab
				} else {
					bgColor = tb.appTheme.TabBar.InactiveTab
				}
				paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					// Tab content using Clickable directly with material styling
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return material.Clickable(gtx, tabButton, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(tb.theme, tab.GetDisplayName())
							if isActive {
								label.Color = tb.appTheme.TabBar.ActiveText
							} else {
								label.Color = tb.appTheme.TabBar.InactiveText
							}
							return label.Layout(gtx)
						})
					}),
					// Close button with better styling
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if len(tb.tabManager.GetTabs()) > 1 {
							return material.Clickable(gtx, closeButton, func(gtx layout.Context) layout.Dimensions {
								inset := layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}
								return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									closeLabel := material.Caption(tb.theme, "×")
									closeLabel.Color = tb.appTheme.TabBar.CloseButton
									if closeButton.Hovered() {
										closeLabel.Color = tb.appTheme.TabBar.CloseButtonHover
									}
									return closeLabel.Layout(gtx)
								})
							})
						}
						return layout.Dimensions{}
					}),
				)
			},
		)
	})
}

// HandleEvents processes tab bar events and returns if tabs changed
func (tb *TabBar) HandleEvents(gtx layout.Context) bool {
	tabsChanged := false

	// Handle keyboard shortcuts for tab navigation
	for {
		event, ok := gtx.Event(
			key.Filter{Name: key.NameTab, Optional: key.ModCtrl},
			key.Filter{Name: key.NameTab, Optional: key.ModCtrl | key.ModShift},
			key.Filter{Name: "W", Optional: key.ModCtrl},
			key.Filter{Name: "T", Optional: key.ModCtrl},
			key.Filter{Name: key.NamePageUp, Optional: key.ModCtrl},
			key.Filter{Name: key.NamePageDown, Optional: key.ModCtrl},
			key.Filter{Name: "1", Optional: key.ModCtrl},
			key.Filter{Name: "2", Optional: key.ModCtrl},
			key.Filter{Name: "3", Optional: key.ModCtrl},
			key.Filter{Name: "4", Optional: key.ModCtrl},
			key.Filter{Name: "5", Optional: key.ModCtrl},
			key.Filter{Name: "6", Optional: key.ModCtrl},
			key.Filter{Name: "7", Optional: key.ModCtrl},
			key.Filter{Name: "8", Optional: key.ModCtrl},
			key.Filter{Name: "9", Optional: key.ModCtrl},
		)
		if !ok {
			break
		}

		if keyEvent, ok := event.(key.Event); ok && keyEvent.State == key.Press {
			tabs := tb.tabManager.GetTabs()
			currentIndex := tb.tabManager.GetActiveTabIndex()

			switch {
			case keyEvent.Name == key.NameTab && keyEvent.Modifiers.Contain(key.ModCtrl) && !keyEvent.Modifiers.Contain(key.ModShift):
				// Ctrl+Tab - next tab
				nextIndex := (currentIndex + 1) % len(tabs)
				tb.tabManager.SetActiveTab(nextIndex)
				tabsChanged = true
			case keyEvent.Name == key.NameTab && keyEvent.Modifiers.Contain(key.ModCtrl) && keyEvent.Modifiers.Contain(key.ModShift):
				// Ctrl+Shift+Tab - previous tab
				prevIndex := currentIndex - 1
				if prevIndex < 0 {
					prevIndex = len(tabs) - 1
				}
				tb.tabManager.SetActiveTab(prevIndex)
				tabsChanged = true
			case keyEvent.Name == key.NamePageUp && keyEvent.Modifiers.Contain(key.ModCtrl):
				// Ctrl+PageUp - previous tab (alternative)
				prevIndex := currentIndex - 1
				if prevIndex < 0 {
					prevIndex = len(tabs) - 1
				}
				tb.tabManager.SetActiveTab(prevIndex)
				tabsChanged = true
			case keyEvent.Name == key.NamePageDown && keyEvent.Modifiers.Contain(key.ModCtrl):
				// Ctrl+PageDown - next tab (alternative)
				nextIndex := (currentIndex + 1) % len(tabs)
				tb.tabManager.SetActiveTab(nextIndex)
				tabsChanged = true
			case keyEvent.Name == "W" && keyEvent.Modifiers.Contain(key.ModCtrl):
				// Ctrl+W - close current tab
				if len(tabs) > 1 {
					tb.tabManager.CloseTab(currentIndex)
					tabsChanged = true
				}
			case keyEvent.Name == "T" && keyEvent.Modifiers.Contain(key.ModCtrl):
				// Ctrl+T - new tab (browser-like)
				tb.tabManager.NewTab("")
				tabsChanged = true
			case keyEvent.Modifiers.Contain(key.ModCtrl) && len(keyEvent.Name) == 1:
				// Ctrl+1-9 - switch to specific tab number
				if keyEvent.Name >= "1" && keyEvent.Name <= "9" {
					tabNum := int(keyEvent.Name[0] - '1') // Convert '1'-'9' to 0-8
					if tabNum < len(tabs) {
						tb.tabManager.SetActiveTab(tabNum)
						tabsChanged = true
					}
				}
			}
		}
	}

	return tabsChanged
}
