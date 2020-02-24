package domain

type ScreenID string

type Screen struct {
	ID ScreenID

	Width, Height int
	X, Y          int
}
