package rest

import (
	"image"

	"github.com/patrislav/wmapi/domain"
)

type wm interface {
	ListScreens() []domain.Screen
	ListScreenWorkspaces(id domain.ScreenID) []domain.Desktop
	ListWorkspaceWindows(id domain.DesktopID) []domain.Window

	ListWindows() ([]domain.Window, error)
	GetWindowByID(id domain.WindowID) *domain.Window
	GetWindowIcon(id domain.WindowID) (image.Image, error)
}
