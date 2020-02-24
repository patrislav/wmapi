package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"image/png"
	"net/http"

	"github.com/patrislav/wmapi/domain"
)

type geom struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type windowPayload struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Desktop   *uint    `json:"desktop"`
	Type      []string `json:"type,omitempty"`
	Geom      *geom    `json:"geometry,omitempty"`
	Class     string   `json:"class,omitempty"`
	Instance  string   `json:"instance,omitempty"`
	States    []string `json:"states,omitempty"`
	Protocols []string `json:"protocols,omitempty"`
	PID       uint     `json:"pid,omitempty"`
}

func (windowPayload) Render(w http.ResponseWriter, r *http.Request) error { return nil }

type windowListPayload []windowPayload

func (windowListPayload) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (s *Server) getWindow(w http.ResponseWriter, r *http.Request) {
	urlWinID := chi.URLParam(r, "windowID")
	winID, err := domain.NewWindowID(urlWinID)
	if err != nil {
		render.Render(w, r, errPayload{statusCode: 404, Message: fmt.Sprintf("Invalid window ID: %s", urlWinID)})
		return
	}
	window := s.wm.GetWindowByID(winID)
	if window == nil {
		render.Render(w, r, errPayload{statusCode: 404, Message: fmt.Sprintf("Window not found: %s", urlWinID)})
		return
	}
	render.Render(w, r, getWindowPayload(*window))
}

func (s *Server) listAllWindows(w http.ResponseWriter, r *http.Request) {
	windows, err := s.wm.ListWindows()
	if err != nil {
		render.Render(w, r, errPayload{statusCode: 500, Message: fmt.Sprintf("Could not retrieve windows: %v", err)})
		return
	}
	payload := make(windowListPayload, len(windows))
	for i, win := range windows {
		payload[i] = getWindowPayload(win)
	}
	render.Render(w, r, payload)
}

func (s *Server) getWindowIcon(w http.ResponseWriter, r *http.Request) {
	urlWinID := chi.URLParam(r, "windowID")
	winID, err := domain.NewWindowID(urlWinID)
	if err != nil {
		render.Render(w, r, errPayload{statusCode: 404, Message: fmt.Sprintf("Invalid window ID: %s", urlWinID)})
		return
	}
	icon, err := s.wm.GetWindowIcon(winID)
	if err != nil {
		render.Render(w, r, errPayload{statusCode: 404, Message: fmt.Sprintf("Window has no icon: %s", urlWinID)})
		return
	}
	if err := png.Encode(w, icon); err != nil {
		render.Render(w, r, errPayload{statusCode: 500, Message: fmt.Sprintf("Failed to render icon: %s", urlWinID)})
		return
	}
}

func parseGeom(g *domain.Geom) *geom {
	if g == nil {
		return nil
	}
	return &geom{
		X:      g.X,
		Y:      g.Y,
		Width:  g.Width,
		Height: g.Height,
	}
}

func getWindowPayload(win domain.Window) windowPayload {
	var desktop *uint
	if win.Desktop != nil {
		index := uint(*win.Desktop)
		desktop = &index
	}
	return windowPayload{
		ID:        int(win.ID),
		Name:      win.Name,
		Type:      win.Type,
		Geom:      parseGeom(win.Geom),
		States:    win.States,
		Class:     win.Class,
		Instance:  win.Instance,
		Protocols: win.Protocols,
		PID:       win.PID,
		Desktop:   desktop,
	}
}
