package domain

import (
	"strconv"
)

type WindowID int

func NewWindowID(s string) (WindowID, error) {
	id, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return WindowID(id), nil
}

type Window struct {
	ID      WindowID
	Desktop *DesktopID

	Name     string
	Type     []string
	Geom     *Geom
	Class    string
	Instance string
	PID      uint

	Protocols []string
	States    []string
}
