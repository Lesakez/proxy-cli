//go:build !windows

package logger

// enableVT is a no-op on non-Windows platforms.
// ANSI codes are supported natively on Linux/macOS terminals.
func enableVT() {}
