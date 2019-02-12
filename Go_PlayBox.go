package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

/* Go text processing suport(SJIS)
https://github.com/golang/text

** Download/Install
The easiest way to install is to run go get -u golang.org/x/text.
You can also manually git clone the repository to $GOPATH/src/golang.org/x/text.
*/

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

// ScaleDefs ... ScaleDefs struct
type ScaleDefs struct {
	scale string
	note  string
}

// PlayData ... struct PlayData
type PlayData struct {
	scale  string
	note   string
	length int
}

// Load the scale definition file.
func loadDefFile(filename string) []ScaleDefs {
	fp, err := os.Open(filename)
	if os.IsNotExist(err) {
		fmt.Printf("%s not found.", filename)
		return nil
	}

	// File Open
	var defs []ScaleDefs
	scanner := bufio.NewScanner(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
	for scanner.Scan() {
		tempStr := scanner.Text()

		pos := strings.Index(tempStr, "//")
		if pos >= 0 {
			tempStr = tempStr[:pos]
		}

		tempStr = strings.Replace(tempStr, " ", "", -1)
		tempStr = strings.Replace(tempStr, "\t", "", -1)
		flds := strings.Split(tempStr, "=")

		if tempStr != "" && len(flds) >= 2 {
			currentDef := new(ScaleDefs)
			currentDef.scale, currentDef.note = flds[0], flds[1]
			defs = append(defs, *currentDef)
		}
	}

	return defs
}

// Load the score file.
func loadPlayFile(filename string) []PlayData {
	fp, err := os.Open(filename)
	if os.IsNotExist(err) {
		fmt.Printf("%s not found.", filename)
		return nil
	}

	// File Open
	var pData []PlayData
	scanner := bufio.NewScanner(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
	for scanner.Scan() {
		tempStr := scanner.Text()

		pos := strings.Index(tempStr, "//")
		if pos >= 0 {
			tempStr = tempStr[:pos]
		}

		tempStr = strings.Replace(tempStr, " ", "", -1)
		tempStr = strings.Replace(tempStr, "\t", "", -1)
		flds := strings.Split(tempStr, "=")

		if tempStr != "" && len(flds) >= 2 {
			currentData := new(PlayData)
			currentData.scale = flds[0]
			currentData.note = ""
			currentData.length, err = strconv.Atoi(flds[1])

			if err != nil {
				fmt.Printf("%s Atoi() error!!", flds[1])
			}

			pData = append(pData, *currentData)
		}
	}

	return pData
}

// Search the musical scale character string and set the note number.
func replaceScale2Freq(defs *[]ScaleDefs, pData *[]PlayData) {
	for i := 0; i < len(*pData); i++ {
		scale := strings.Split((*pData)[i].scale, ",")
		for _, temp := range scale {
			for _, currentLen := range *defs {
				if temp == currentLen.scale {
					if (*pData)[i].note == "" {
						(*pData)[i].note = currentLen.note
					} else {
						(*pData)[i].note += "," + currentLen.note
					}
					break
				}
			}
		}
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage)\n" +
			"go run " + os.Args[0] + " musicDataFile <timbre>\n")
		return
	}

	noteNumberFile := "./note-number.dat"
	if fileExists(noteNumberFile) == false {
		fmt.Printf("%s not found.", noteNumberFile)
		return
	}

	if fileExists(os.Args[1]) == false {
		fmt.Printf("%s not found.", os.Args[1])
		return
	}

	var timbre int
	if len(os.Args) >= 3 {
		readval, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("%s ... timbre is a not integer!!", os.Args[2])
			return
		}
		timbre = readval
	} else {
		timbre = 1
	}

	// Load the scale definition file.
	defs := loadDefFile(noteNumberFile)

	// Load PlayData file.
	pData := loadPlayFile(os.Args[1])

	// Set note number.
	replaceScale2Freq(&defs, &pData)

	initData := timbre*256 + 0xc0

	// Set Integer Size.
	const intSize = 32 << (^uint(0) >> 63)

	fmt.Printf("intSize = %d, initData = 0x%04x\n", intSize, initData)

	// Initialize the MyMIDI struct and functions
	pm := new(MyMIDI)
	pm.Init(initData)

	fmt.Printf("Load Done. Play start!!\n")

	for i, currentpData := range pData {
		if currentpData.note != "" {
			fmt.Printf("[%d] = %s( %s ), %d [ms]\n", i, currentpData.scale, currentpData.note, currentpData.length)
			cnote := strings.Split(currentpData.note, ",")

			for _, data := range cnote {
				// Press the keyboad.
				playOn := "0x7f" + data + "90"
				playData, err := strconv.ParseInt(playOn, 0, intSize)
				if err != nil {

				}

				pm.OutOnly(int(playData))
			}

			// Keep pressed.
			time.Sleep(time.Duration(currentpData.length) * time.Millisecond)

			for _, data := range cnote {
				// Release the keyboad.
				playOff := "0x7f" + data + "80"
				playData, err := strconv.ParseInt(playOff, 0, intSize)
				if err != nil {

				}

				pm.OutOnly(int(playData))
			}
		} else {
			// Rest.
			fmt.Printf("[%d] = rest, %d [ms]\n", i, currentpData.length)
			time.Sleep(time.Duration(currentpData.length) * time.Millisecond)
		}
	}

	pm.Close()
	fmt.Println()
}
