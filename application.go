package tview

import (
	"sync"

	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/ubcell"
)

// Application represents the top node of an application.
//
// It is not strictly required to use this class as none of the other classes
// depend on it. However, it provides useful tools to set up an application and
// plays nicely with all widgets.
type Application struct {
	sync.RWMutex

	//The screen configuration
	cfg *Config

	// The application's screen.
	screen ubcell.Screen

	// The primitive which currently has the keyboard focus.
	focus Primitive

	// The root primitive to be seen on the screen.
	root Primitive

	// Whether or not the application resizes the root primitive.
	rootFullscreen bool

	// An optional capture function which receives a key event and returns the
	// event to be forwarded to the default input handler (nil if nothing should
	// be forwarded).
	keyCapture func(event pixelgl.Event) pixelgl.Event

	// An optional callback function which is invoked just before the root
	// primitive is drawn.
	beforeDraw func(screen ubcell.Screen) bool

	// An optional callback function which is invoked after the root primitive
	// was drawn.
	afterDraw func(screen ubcell.Screen)

	// If this value is true, the application has entered suspended mode.
	suspended bool
}

// NewApplication creates and returns a new application.
func NewApplication(cfg *Config) (*Application, error) {
	return &Application{cfg: cfg}, nil
}

func (a *Application) Screen() ubcell.Screen {
	return a.screen
}

// SetKeyCapture sets a function which captures all key events before they are
// forwarded to the key event handler of the primitive which currently has
// focus. This function can then choose to forward that key event (or a
// different one) by returning it or stop the key event processing by returning
// nil.
//
// Note that this also affects the default event handling of the application
// itself: Such a handler can intercept the Ctrl-C event which closes the
// applicatoon.
func (a *Application) SetKeyCapture(capture func(event pixelgl.Event) pixelgl.Event) *Application {

	a.keyCapture = capture
	return a
}

// GetKeyCapture returns the function installed with SetKeyCapture() or nil
// if no such function has been installed.
func (a *Application) GetKeyCapture() func(event pixelgl.Event) pixelgl.Event {

	return a.keyCapture
}

// Run starts the application and thus the event loop. This function returns
// when Stop() was called.
func (a *Application) Run() error {

	var err error
	a.Lock()

	// Make a screen.
	a.screen, err = ubcell.NewScreen(a.cfg)

	if err != nil {
		a.Unlock()
		return err
	}

	if err = a.screen.Init(); err != nil {
		a.Unlock()
		return err
	}

	// We catch panics to clean up because they mess up the terminal.
	defer func() {
		if p := recover(); p != nil {
			if a.screen != nil {
				a.screen.Fini()
			}
			panic(p)
		}
	}()

	// Draw the screen for the first time.
	a.Unlock()

	a.Draw()
	//post event

	// Start event loop.
	for {
		a.Lock()
		screen := a.screen
		if a.suspended {
			a.suspended = false // Clear previous suspended flag.
		}
		a.Unlock()
		if screen == nil {
			break
		}

		// Wait for next event - blocking...
		event := a.screen.PollEvent()
		if event == nil {
			a.Lock()
			if a.suspended {
				// This screen was renewed due to suspended mode.
				a.suspended = false
				a.Unlock()
				continue // Resume.
			}
			a.Unlock()

			// The screen was finalized. Exit the loop.
			break
		}

		switch event := event.(type) {

		case *pixelgl.CursorEvent, *pixelgl.ScrollEvent:
			a.RLock()
			p := a.root
			a.RUnlock()
			if handler := p.MouseHandler(); p != nil {
				handler(event, func(p Primitive) {
					a.SetFocus(p)
				})
			}

		//do nothing for now

		case *pixelgl.KeyEv, *pixelgl.ChaEv:

			a.RLock()
			p := a.focus
			a.RUnlock()

			//do key-specific stuff here
			if ev, ok := event.(*pixelgl.KeyEv); ok && ev.Act == pixelgl.RELEASE {
				// Ctrl-C closes the application.
				if *ev == pixelgl.KeyCtrlC {
					a.Stop()
				}
				break
			}

			// Intercept keys.
			//rename!
			if a.keyCapture != nil {
				event = a.keyCapture(event)
				if event == nil {
					break
				}
			}

			// Pass other key events to the currently focused primitive.
			if p != nil {
				if handler := p.KeyHandler(); handler != nil {
					handler(event, func(p Primitive) {
						a.SetFocus(p)
					})

					a.Draw()

				}
			}
		case *pixelgl.ResizeEvent:
			a.Lock()
			screen := a.screen
			a.Unlock()
			screen.Clear()
			a.Draw()

		}
	}

	return nil
}

// Stop stops the application, causing Run() to return.
func (a *Application) Stop() {

	if a.screen == nil {
		return
	}
	a.screen.Fini()
	a.screen = nil
}

// Suspend temporarily suspends the application by exiting terminal UI mode and
// invoking the provided function "f". When "f" returns, terminal UI mode is
// entered again and the application resumes.
//
// A return value of true indicates that the application was suspended and "f"
// was called. If false is returned, the application was already suspended,
// terminal UI mode was not exited, and "f" was not called.
//doesn't work currently...
//func (a *Application) Suspend(f func()) bool {

//
//	if a.suspended || a.screen == nil {
//		// Application is already suspended.
//		return false
//	}
//
//	// Enter suspended mode.
//	a.suspended = true
//	a.Stop()
//
//	// Deal with panics during suspended mode. Exit the program.
//	defer func() {
//		if p := recover(); p != nil {
//			fmt.Println(p)
//			os.Exit(1)
//		}
//	}()
//
//	// Wait for "f" to return.
//	f()
//
//	// Make a new screen and redraw.
//	var err error
//	a.screen, err = ubcell.NewScreen()
//	if err != nil {
//		panic(err)
//	}
//	if err = a.screen.Init(); err != nil {
//		panic(err)
//	}
//	a.Draw()
//
//	// Continue application loop.
//	return true
//}

// Draw refreshes the screen. It calls the Draw() function of the application's
// root primitive and then syncs the screen buffer.
func (a *Application) Draw() *Application {

	a.RLock()
	screen := a.screen
	root := a.root
	fullscreen := a.rootFullscreen
	before := a.beforeDraw
	after := a.afterDraw
	a.RUnlock()

	// Maybe we're not ready yet or not anymore.
	if screen == nil || root == nil {
		return a
	}

	// Resize if requested.
	if fullscreen && root != nil {
		width, height := screen.Size()
		root.SetRect(0, 0, width, height)
	}

	// Call before handler if there is one.
	if before != nil {
		if before(screen) {
			screen.Show()
			return a
		}

	}

	// Draw all primitives.

	root.Draw(screen)

	// Call after handler if there is one.
	if after != nil {
		after(screen)
	}

	// Sync screen.
	a.screen.Show()

	return a
}

// SetBeforeDrawFunc installs a callback function which is invoked just before
// the root primitive is drawn during screen updates. If the function returns
// true, drawing will not continue, i.e. the root primitive will not be drawn
// (and an after-draw-handler will not be called).
//
// Note that the screen is not cleared by the application. To clear the screen,
// you may call screen.Clear().
//
// Provide nil to uninstall the callback function.
func (a *Application) SetBeforeDrawFunc(handler func(screen ubcell.Screen) bool) *Application {

	a.beforeDraw = handler
	return a
}

// GetBeforeDrawFunc returns the callback function installed with
// SetBeforeDrawFunc() or nil if none has been installed.
func (a *Application) GetBeforeDrawFunc() func(screen ubcell.Screen) bool {

	return a.beforeDraw
}

// SetAfterDrawFunc installs a callback function which is invoked after the root
// primitive was drawn during screen updates.
//
// Provide nil to uninstall the callback function.
func (a *Application) SetAfterDrawFunc(handler func(screen ubcell.Screen)) *Application {

	a.afterDraw = handler
	return a
}

// GetAfterDrawFunc returns the callback function installed with
// SetAfterDrawFunc() or nil if none has been installed.
func (a *Application) GetAfterDrawFunc() func(screen ubcell.Screen) {

	return a.afterDraw
}

// SetRoot sets the root primitive for this application. If "fullscreen" is set
// to true, the root primitive's position will be changed to fill the screen.
//
// This function must be called at least once or nothing will be displayed when
// the application starts.
//
// It also calls SetFocus() on the primitive.
func (a *Application) SetRoot(root Primitive, fullscreen bool) *Application {

	a.Lock()
	a.root = root
	a.rootFullscreen = fullscreen
	if a.screen != nil {
		a.screen.Clear()
	}
	a.Unlock()

	a.SetFocus(root)

	return a
}

// ResizeToFullScreen resizes the given primitive such that it fills the entire
// screen.
func (a *Application) ResizeToFullScreen(p Primitive) *Application {

	a.RLock()
	width, height := a.screen.Size()
	a.RUnlock()
	p.SetRect(0, 0, width, height)
	return a
}

// SetFocus sets the focus on a new primitive. All key events will be redirected
// to that primitive. Callers must ensure that the primitive will handle key
// events.
//
// Blur() will be called on the previously focused primitive. Focus() will be
// called on the new primitive.
func (a *Application) SetFocus(p Primitive) *Application {

	a.Lock()
	if a.focus != nil {
		a.focus.Blur()
	}
	a.focus = p
	if a.screen != nil {
		a.screen.HideCursor()
	}

	a.Unlock()

	if p != nil {
		p.Focus(func(p Primitive) {
			a.SetFocus(p)
		})
	}

	return a
}

// GetFocus returns the primitive which has the current focus. If none has it,
// nil is returned.
func (a *Application) GetFocus() Primitive {

	a.RLock()
	defer a.RUnlock()
	return a.focus
}
