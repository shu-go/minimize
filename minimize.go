package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/shu-go/gli"
)

// Version is app version
var Version string

func init() {
	if Version == "" {
		Version = "dev-" + time.Now().Format("20060102")
	}
}

type globalCmd struct {
	Restore bool `cli:"restore, r"  help:"restore from minimized/maximized"`
}

func (g globalCmd) Run() error {
	ppid := os.Getppid()

	wins, err := listAllWindows()
	if err != nil {
		return err
	}

	for _, w := range wins {
		if w.PID == ppid {
			if g.Restore {
				showWindow.Call(uintptr(w.Handle), SW_RESTORE)
			} else {
				showWindow.Call(uintptr(w.Handle), SW_MINIMIZE)
			}
			break
		}
	}
	return nil
}

func main() {
	//defer rog.DoneDebugging()()

	app := gli.NewWith(&globalCmd{})
	app.Name = "minimize"
	app.Desc = "minimize/restore parent window"
	app.Version = Version
	app.Usage = `minimize [-r]`
	app.Copyright = "(C) 2017 Shuhei Kubota"

	app.Run(os.Args)
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
	SW_RESTORE  = 9
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
