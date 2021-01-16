package screen

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	libUser32               = syscall.NewLazyDLL("user32.dll")
	funcEnumWindows         = libUser32.NewProc("EnumWindows")
	funcGetWindowTextLength = libUser32.NewProc("GetWindowTextLengthW")
	funcGetWindowText       = libUser32.NewProc("GetWindowTextW")
	funcGetWindowRect       = libUser32.NewProc("GetWindowRect")
	// funcGetDesktopWindow, _    = syscall.GetProcAddress(syscall.Handle(libUser32), "GetDesktopWindow")
	// funcEnumDisplayMonitors, _ = syscall.GetProcAddress(syscall.Handle(libUser32), "EnumDisplayMonitors")
)

type searchContext struct {
	name   string
	bounds image.Rectangle
}

// SaveImage *image.RGBA to filePath with PNG format.
func SaveImage(img image.Image, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err, "unable to save image")
		return
	}
	defer file.Close()
	if file != nil && img != nil {
		png.Encode(file, img)
	}
}

// GetScreenCoords get the screen coords of a window
func GetScreenCoords(windowName string) image.Rectangle {
	var ctx searchContext
	ctx.name = windowName

	funcEnumWindows.Call(syscall.NewCallback(getWindowCallback), uintptr(unsafe.Pointer(&ctx)))
	//fmt.Println(ctx.bounds)
	return ctx.bounds
}

func getWindowCallback(hwnd win.HWND, user uintptr) uintptr {
	length := getWindowTextLength(hwnd)

	if length > 0 {
		var ctx *searchContext
		ctx = (*searchContext)(unsafe.Pointer(user))
		winName := getWindowText(hwnd, length)
		//fmt.Printf("Name -%s-\n", winName)

		if ctx.name == winName {
			var rect win.RECT

			getWindowRect(hwnd, &rect)

			if rect.Right-rect.Left > 600 {
				//fmt.Printf("Found %s : bot %d left %d top %d right %d\n", winName, rect.Bottom, rect.Left, rect.Top, rect.Right)
				ctx.bounds = image.Rect(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom))
				return uintptr(0)
			}
		}
	}
	//fmt.Println(rect)
	return uintptr(1)
}

func getWindowTextLength(hwnd win.HWND) int {
	//syscall.Syscall(funcGetWindowTextLength, 1, uintptr(hwnd), 0, 0)
	//name := make([]uint16, 0, 25)
	//funcGetWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&name)), 24)
	//fmt.Println(name)
	//syscall.Syscall(funcGetWindowTextLength, 1, uintptr(hwnd), 0, 0)
	ret, _, _ := funcGetWindowTextLength.Call(uintptr(hwnd))
	return int(ret) + 1
}

func getWindowText(hwnd win.HWND, len int) string {
	name := make([]uint16, len)
	funcGetWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&name[0])), uintptr(len))

	return syscall.UTF16ToString(name)
}

func getWindowRect(hwnd win.HWND, rect *win.RECT) {

	funcGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect)))
}
