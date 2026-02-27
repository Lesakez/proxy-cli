//go:build windows

package logger

import (
	"os"
	"syscall"
	"unsafe"
)

const enableVirtualTerminalProcessing = 0x0004

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode = kernel32.NewProc("GetConsoleMode")
	setConsoleMode = kernel32.NewProc("SetConsoleMode")
)

// enableVT enables ANSI/VT100 escape codes on Windows 10+ consoles.
func enableVT() {
	handle := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	setConsoleMode.Call(uintptr(handle), uintptr(mode|enableVirtualTerminalProcessing))
}
