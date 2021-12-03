package utils

var ScreenBufferActive = false

// EnterScreenBuffer Switches terminal to alternate screen buffer to retain command history
//  of host process
func EnterScreenBuffer() {
	print("\033[?1049h")
	ScreenBufferActive = true
}

// ExitScreenBuffer Exits the alternate screen buffer and returns to that of host process
func ExitScreenBuffer() {
	print("\033[?1049l")
	ScreenBufferActive = false
}
