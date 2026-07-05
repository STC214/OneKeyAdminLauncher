package main

import (
	"os"
	"runtime"

	"program-launch-manager/internal/process"
	"program-launch-manager/internal/winui"
)

func main() {
	runtime.LockOSThread()

	if !process.EnsureAdmin() {
		os.Exit(0)
	}
	os.Exit(winui.Run())
}
