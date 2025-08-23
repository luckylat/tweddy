package style

import "image/color"

// ApplicationColors defines the overall application color scheme
type ApplicationColors struct {
	MainBackground   color.NRGBA
	EditorBackground color.NRGBA
	HeaderBackground color.NRGBA
	ButtonBackground color.NRGBA
	ButtonText       color.NRGBA
	ButtonHover      color.NRGBA
	BorderColor      color.NRGBA
}

// TabBarColors defines the color scheme for tab bars
type TabBarColors struct {
	Background       color.NRGBA
	ActiveTab        color.NRGBA
	InactiveTab      color.NRGBA
	ActiveText       color.NRGBA
	InactiveText     color.NRGBA
	CloseButton      color.NRGBA
	CloseButtonHover color.NRGBA
}

// Theme combines all color configurations for the application
type Theme struct {
	Application ApplicationColors
	TabBar      TabBarColors
}

// DefaultTheme returns the default theme for the application
func DefaultTheme() Theme {
	return Theme{
		Application: ApplicationColors{
			MainBackground:   color.NRGBA{R: 248, G: 249, B: 250, A: 255}, // Light gray background
			EditorBackground: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White editor background
			HeaderBackground: color.NRGBA{R: 230, G: 230, B: 230, A: 255}, // Medium gray header for better button contrast
			ButtonBackground: color.NRGBA{R: 52, G: 73, B: 94, A: 255},   // Dark blue-gray button background
			ButtonText:       color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White button text
			ButtonHover:      color.NRGBA{R: 72, G: 93, B: 114, A: 255},  // Lighter blue-gray button hover
			BorderColor:      color.NRGBA{R: 180, G: 180, B: 180, A: 255}, // Medium gray border
		},
		TabBar: TabBarColors{
			Background:       color.NRGBA{R: 220, G: 220, B: 220, A: 255}, // Light gray background for contrast
			ActiveTab:        color.NRGBA{R: 70, G: 130, B: 180, A: 255},  // Steel blue for active tab
			InactiveTab:      color.NRGBA{R: 190, G: 190, B: 190, A: 255}, // Medium gray for inactive tab
			ActiveText:       color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White text for active tab
			InactiveText:     color.NRGBA{R: 60, G: 60, B: 60, A: 255},    // Dark gray text for inactive tab
			CloseButton:      color.NRGBA{R: 120, G: 120, B: 120, A: 255}, // Medium gray close button
			CloseButtonHover: color.NRGBA{R: 220, G: 60, B: 60, A: 255},   // Bright red close button on hover
		},
	}
}

// DefaultTabBarColors returns the default color scheme for tab bars
// Deprecated: Use DefaultTheme().TabBar instead
func DefaultTabBarColors() TabBarColors {
	return DefaultTheme().TabBar
}
