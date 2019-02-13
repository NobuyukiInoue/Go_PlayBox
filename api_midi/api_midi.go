package api_midi

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// MyMIDI ... MyMIDI struct
type MyMIDI struct {
	initData   int
	MIDIMAPPER int
	h          *uint
	dll        *syscall.DLL
}

// PrintInfo ... Print struct MyMIDI members information.
func (pm *MyMIDI) PrintInfo() {
	fmt.Printf("addr = %x, initData = %x\n", &pm.initData, pm.initData)
	fmt.Printf("addr = %x, MIDIMAPPER = %x\n", &pm.MIDIMAPPER, pm.MIDIMAPPER)
	fmt.Printf("addr = %x, h = %x\n", &pm.h, pm.h)
}

// Init ... MIDI Init
func (pm *MyMIDI) Init(initData int) {
	dll, err := syscall.LoadDLL("winmm.dll")
	if err != nil {
		panic(err)
	}

	pm.dll = dll
	pm.MIDIMAPPER = -1
	pm.initData = initData

	proc, err := pm.dll.FindProc("midiOutOpen")
	if err != nil {
		panic(err)
	}

	proc.Call(uintptr(unsafe.Pointer(&pm.h)), uintptr(pm.MIDIMAPPER), uintptr(0), uintptr(0), uintptr(0))

	proc, err = pm.dll.FindProc("midiOutShortMsg")
	if err != nil {
		panic(err)
	}

	proc.Call(uintptr(unsafe.Pointer(pm.h)), uintptr(pm.initData))
}

// Out ... MIDI Output with time.Sleep
func (pm *MyMIDI) Out(outData int, length time.Duration) {
	proc, err := pm.dll.FindProc("midiOutShortMsg")
	if err != nil {
		panic(err)
	}

	proc.Call(uintptr(unsafe.Pointer(pm.h)), uintptr(outData))
	time.Sleep(length * time.Millisecond)
}

// OutOnly ... MIDI Output without timeSleep
func (pm *MyMIDI) OutOnly(outData int) {
	proc, err := pm.dll.FindProc("midiOutShortMsg")
	if err != nil {
		panic(err)
	}

	proc.Call(uintptr(unsafe.Pointer(pm.h)), uintptr(outData))
}

// Close ... MIDI Close
func (pm *MyMIDI) Close() {
	proc, err := pm.dll.FindProc("midiOutReset")
	if err != nil {
		panic(err)
	}

	proc.Call(uintptr(unsafe.Pointer(pm.h)))
}
