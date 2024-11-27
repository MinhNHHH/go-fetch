package fetch

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	asci "github.com/minhnh/fetch/internal/ascii"
)

type ClientDetail struct {
	SysInfor SystemInfor
	AsciiArt *AsciiArt
}

func DefaultArtSys() string {
	return runtime.GOOS
}

func HandleClient(cmd []string) {
	done := make(chan bool)
	var sysInfor SystemInfor

	go func() {
		sysInfor = NewSysInfor()
		done <- true
	}()

	if <-done {
		ascii := DefaultArt(DefaultArtSys())
		client := &ClientDetail{
			AsciiArt: ascii,
			SysInfor: sysInfor,
		}
		client.handleCommand(cmd)
	}
}

func (c *ClientDetail) handleCommand(command []string) {
	disable := []string{}
	seemore := []string{}
	switch command[0] {
	case "list":
	case "source":
		if len(command) >= 2 {
			ascii, err := NewAsciiArt(command[1])
			if err != nil {
				log.Fatalf(err.Error())
			}
			c.AsciiArt = ascii
		}
	case "disable":
		disable = command[1:]
	case "ascii_distro":
		if len(command) >= 2 {
			ascii := DefaultArt(command[1])
			c.AsciiArt = ascii
		}
	case "ascii_color":
		if len(command) >= 2 {
			for i, color := range command[1:] {
				codeColor := fmt.Sprintf("${c%d}", i+1)
				asci.PlaceHolder[codeColor] = asci.CodeColor[color]
			}
		}
	}
	c.PrintInfor(disable, seemore)
}

func (c *ClientDetail) CountPattern(input string) int {
	count := 0
	for label := range asci.PlaceHolder {
		if label == "${c0}" {
			count += strings.Count(input, label) * 4
		} else {
			count += strings.Count(input, label) * 7
		}
	}
	return count
}

func (c *ClientDetail) replacePlaceHolder(input string) string {
	for label, color := range asci.PlaceHolder {
		input = strings.ReplaceAll(input, label, color)
	}
	return input
}

func (c *ClientDetail) PrintInfor(disable, seemore []string) {
	listInfor := c.SysInfor.ListSysInfor(disable, seemore)
	maxLines := Max(len(c.AsciiArt.Lines), len(listInfor))
	asciiLine, sysInformLine := "", ""
	for i := 0; i < maxLines; i++ {
		pattern := c.CountPattern(c.AsciiArt.Lines[i])
		if i < len(c.AsciiArt.Lines) {
			c.AsciiArt.Lines[i] = c.replacePlaceHolder(c.AsciiArt.Lines[i])
			asciiLine = c.AsciiArt.Lines[i]
		}

		if i < len(listInfor) {
			sysInformLine = listInfor[i]
		} else {
			sysInformLine = ""
		}
		originalDistance := c.AsciiArt.MaxCleanLen + pattern
		padding := 5
		fmt.Printf("%-*s %s\n", originalDistance+padding, asciiLine, sysInformLine)
	}
}

// drawColorBoxesInLine builds and returns a string representing multiple colored boxes of given width, height, and background colors in one line
func DrawColorBoxesInLine(colorCodes []string, width int, height int) string {
	var sb strings.Builder
	reset := "\033[0m"

	for i := 0; i < height; i++ {
		for _, colorCode := range colorCodes {
			sb.WriteString(colorCode) // Set background color
			for j := 0; j < width; j++ {
				sb.WriteString(" ") // Print space to form the box
			}
			sb.WriteString(reset) // Reset color to separate boxes
		}
	}

	return sb.String()
}
