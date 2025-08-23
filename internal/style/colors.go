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


// Theme combines all color configurations for the application
type Theme struct {
	Application ApplicationColors
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
	}
}

