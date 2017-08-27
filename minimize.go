package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	ppid := os.Getppid()

	wins, _ := listAllWindows()
	for _, w := range wins {
		if w.PID == ppid {
			showWindow.Call(uintptr(w.Handle), SW_MINIMIZE)
			break
		}
	}
}

type (
	Window struct {
		Title  string
		Handle syscall.Handle
		PID    int
	}
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	isWindow                 = user32.NewProc("IsWindow")
	enumWindows              = user32.NewProc("EnumWindows")
	getWindowText            = user32.NewProc("GetWindowTextW")
	getWindowTextLength      = user32.NewProc("GetWindowTextLengthW")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	showWindow               = user32.NewProc("ShowWindow")
)

const (
	SW_MINIMIZE = 6
)

func listAllWindows() (wins []*Window, err error) {
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		b, _, _ := isWindow.Call(uintptr(hwnd))
		if b == 0 {
			return 1
		}

		title := ""
		/*
			tlen, _, _ := getWindowTextLength.Call(uintptr(hwnd))
			if tlen != 0 {
				tlen++
				buff := make([]uint16, tlen)
				getWindowText.Call(
					uintptr(hwnd),
					uintptr(unsafe.Pointer(&buff[0])),
					uintptr(tlen),
				)
				title = syscall.UTF16ToString(buff)
			}
		*/

		var processID uintptr
		getWindowThreadProcessId.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&processID)),
		)

		win := &Window{
			Title:  title,
			Handle: hwnd,
			PID:    int(processID),
		}
		wins = append(wins, win)

		return 1
	})

	a, _, _ := enumWindows.Call(cb, 0)
	if a == 0 {
		return nil, fmt.Errorf("USER32.EnumWindows returned FALSE")
	}

	return wins, nil
}
