package x11

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/patrislav/wmapi/domain"
)

type Client struct {
	xu *xgbutil.XUtil
}

func NewClient() (*Client, error) {
	xu, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	return &Client{xu: xu}, nil
}

func (c *Client) ListScreens() []domain.Screen {
	return []domain.Screen{}
}

func (c *Client) ListScreenWorkspaces(id domain.ScreenID) []domain.Desktop {
	return []domain.Desktop{}
}

func (c *Client) ListWorkspaceWindows(id domain.DesktopID) []domain.Window {
	return []domain.Window{}
}

func (c *Client) ListWindows() ([]domain.Window, error) {
	ids, err := ewmh.ClientListGet(c.xu)
	if err != nil {
		return nil, err
	}
	winch := make(chan domain.Window, len(ids))
	for _, xid := range ids {
		go func(xid xproto.Window) {
			winch <- windowData(c.xu, xid)
		}(xid)
	}
	wins := make([]domain.Window, 0, len(ids))
	for range ids {
		wins = append(wins, <-winch)
	}
	return wins, nil
}

func (c *Client) GetWindowByID(id domain.WindowID) *domain.Window {
	// Get a list of all managed window ID's, id should be one of them
	xid := xproto.Window(id)
	ids, err := ewmh.ClientListGet(c.xu)
	if err != nil {
		return nil
	}
	for _, clientID := range ids {
		if xid == clientID {
			win := windowData(c.xu, xid)
			return &win
		}
	}
	return nil
}

func getWindowName(xu *xgbutil.XUtil, id xproto.Window) string {
	name, err := ewmh.WmNameGet(xu, id)
	if err != nil {
		name, _ = icccm.WmNameGet(xu, id)
	}
	return name
}

func getWindowGeom(xu *xgbutil.XUtil, id xproto.Window) *domain.Geom {
	dgeom, err := xwindow.New(xu, id).DecorGeometry()
	if err != nil {
		return nil
	}
	x, y, w, h := dgeom.Pieces()
	return &domain.Geom{X: x, Y: y, Width: w, Height: h}
}

func getWindowStates(xu *xgbutil.XUtil, id xproto.Window) []string {
	states, err := ewmh.WmStateGet(xu, id)
	if err != nil {
		return []string{}
	}
	return states
}

func getWindowClassInstance(xu *xgbutil.XUtil, id xproto.Window) (class, instance string) {
	data, err := icccm.WmClassGet(xu, id)
	if err != nil {
		return
	}
	return data.Class, data.Instance
}

func getWindowProtocols(xu *xgbutil.XUtil, id xproto.Window) []string {
	protocols, err := icccm.WmProtocolsGet(xu, id)
	if err != nil {
		return []string{}
	}
	return protocols
}

func getWindowType(xu *xgbutil.XUtil, id xproto.Window) []string {
	typ, err := ewmh.WmWindowTypeGet(xu, id)
	if err != nil {
		return []string{}
	}
	return typ
}

func getWindowPID(xu *xgbutil.XUtil, id xproto.Window) uint {
	pid, err := ewmh.WmPidGet(xu, id)
	if err != nil {
		return 0
	}
	return pid
}

func windowData(xu *xgbutil.XUtil, id xproto.Window) domain.Window {
	name := getWindowName(xu, id)
	geom := getWindowGeom(xu, id)
	states := getWindowStates(xu, id)
	class, instance := getWindowClassInstance(xu, id)
	protocols := getWindowProtocols(xu, id)
	types := getWindowType(xu, id)
	pid := getWindowPID(xu, id)
	desktop := getWindowDesktop(xu, id)

	return domain.Window{
		ID:        domain.WindowID(id),
		Name:      name,
		Type:      types,
		Geom:      geom,
		States:    states,
		Class:     class,
		Instance:  instance,
		Protocols: protocols,
		PID:       pid,
		Desktop:   desktop,
	}
}

func getWindowDesktop(xu *xgbutil.XUtil, id xproto.Window) *domain.DesktopID {
	index, err := ewmh.WmDesktopGet(xu, id)
	if err != nil {
		return nil
	}
	desktop := domain.DesktopID(index)
	return &desktop
}
