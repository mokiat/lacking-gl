package app

import (
	"fmt"
	"image"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/log"
)

var (
	appLogger = log.Path("/lacking-gl/app")
	glLogger  = appLogger.Path("/opengl")
)

// Run starts a new application and opens a single window.
//
// The specified configuration is used to determine how the
// window is initialized.
//
// The specified controller will be used to send notifications
// on window state changes.
func Run(cfg *Config, controller app.Controller) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %w", err)
	}
	defer glfw.Terminate()

	var (
		windowWidth  = cfg.width
		windowHeight = cfg.height
		monitor      *glfw.Monitor
	)
	if cfg.fullscreen {
		monitor = glfw.GetPrimaryMonitor()
		videoMode := monitor.GetVideoMode()
		windowWidth = videoMode.Width
		windowHeight = videoMode.Height
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.SRGBCapable, glfw.True)
	if cfg.maximized {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	}

	window, err := glfw.CreateWindow(windowWidth, windowHeight, cfg.title, monitor, nil)
	if err != nil {
		return fmt.Errorf("failed to create glfw window: %w", err)
	}
	defer window.Destroy()

	if cfg.minWidth != nil || cfg.maxWidth != nil || cfg.minHeight != nil || cfg.maxHeight != nil {
		minWidth := glfw.DontCare
		if cfg.minWidth != nil {
			minWidth = *cfg.minWidth
		}
		minHeight := glfw.DontCare
		if cfg.minHeight != nil {
			minHeight = *cfg.minHeight
		}
		maxWidth := glfw.DontCare
		if cfg.maxWidth != nil {
			maxWidth = *cfg.maxWidth
		}
		maxHeight := glfw.DontCare
		if cfg.maxHeight != nil {
			maxHeight = *cfg.maxHeight
		}
		window.SetSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
	}

	if cfg.icon != "" {
		img, err := openImage(cfg.locator, cfg.icon)
		if err != nil {
			return fmt.Errorf("failed to open icon %q: %w", cfg.icon, err)
		}
		window.SetIcon([]image.Image{img})
	}

	window.MakeContextCurrent()
	defer glfw.DetachCurrentContext()
	glfw.SwapInterval(cfg.swapInterval)

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize opengl: %w", err)
	}

	if glLogger.DebugEnabled() {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(func(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
			switch severity {
			case gl.DEBUG_SEVERITY_LOW:
				glLogger.Debug(message)
			case gl.DEBUG_SEVERITY_MEDIUM:
				glLogger.Warn(message)
			case gl.DEBUG_SEVERITY_HIGH:
				glLogger.Error(message)
			default:
				glLogger.Debug(message)
			}
		}, gl.PtrOffset(0))
	}

	l := newLoop(cfg.locator, cfg.title, window, controller)

	if cfg.cursor != nil {
		cursor := l.CreateCursor(*cfg.cursor)
		defer cursor.Destroy()
		l.UseCursor(cursor)
		defer l.UseCursor(nil)
	}

	if !cfg.cursorVisible {
		l.SetCursorVisible(false)
	}

	return l.Run()
}
