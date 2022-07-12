package prettier

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Prettier formate for console output
type Prettier struct {
	Reset  string
	Red    string
	Green  string
	Yellow string
	Blue   string
	Purple string
	Cyan   string
	Gray   string
	White  string
	timer  interface{}
}

// Create new instance of Prettier.
func NewPrettier() *Prettier {
	prettier := &Prettier{}
	prettier.timer = nil

	if runtime.GOOS == "windows" {
		prettier.Reset = ""
		prettier.Red = ""
		prettier.Green = ""
		prettier.Yellow = ""
		prettier.Blue = ""
		prettier.Purple = ""
		prettier.Cyan = ""
		prettier.Gray = ""
		prettier.White = ""
	} else {
		prettier.Reset = "\033[0m"
		prettier.Red = "\033[31m"
		prettier.Green = "\033[32m"
		prettier.Yellow = "\033[33m"
		prettier.Blue = "\033[34m"
		prettier.Purple = "\033[35m"
		prettier.Cyan = "\033[36m"
		prettier.Gray = "\033[37m"
		prettier.White = "\033[97m"
	}

	return prettier
}

// Start your project and start time clock.
func (prettier *Prettier) Start(projectName string, projectVersion string, projectAuthor string) {
	fmt.Printf("%s\n==================== %s ====================\n", prettier.Green, projectName)
	fmt.Printf("Version: %s\n", projectVersion)
	fmt.Printf("Author: %s\n", projectAuthor)
	fmt.Printf("Licence: MIT %s\n\n", prettier.Reset)
}

// Set timer for your project
func (prettier *Prettier) SetTimer() {
	prettier.timer = time.Now()
}

// Show time from timer start
func (prettier *Prettier) ShowTime() {
	if entryTimer, ok := prettier.timer.(time.Time); ok {
		fmt.Printf("%s\nTime past from start: %f\n", prettier.Cyan, time.Now().Sub(entryTimer).Seconds())
	} else {
		fmt.Printf("\n%s\n%s", "Timer is not set!\nUse prettier.SetTimer().", prettier.Red)
		log.Fatal()
	}
}

// Create default progress bar
func (prettier *Prettier) DefaultBar(steps int, barName string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(steps,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[cyan]"+barName+"[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[red]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	return bar
}

// Print info block.
func (prettier *Prettier) Info(messages map[string]interface{}) {
	fmt.Println(prettier.Cyan + "----Info----")

	for key, value := range messages {
		if entryMessage, ok := value.(string); ok {
			fmt.Println(key + entryMessage)
		}
		if entryMessage, ok := value.(float64); ok {
			fmt.Println(key + strconv.FormatFloat(entryMessage, 'e', 8, 64))
		}
		if entryMessage, ok := value.(int); ok {
			fmt.Println(key + strconv.FormatInt(int64(entryMessage), 10))
		}
	}

	fmt.Println("------------" + prettier.Reset)
	fmt.Println()
}

// Print error message.
func (prettier *Prettier) Error(massage string, err error) {
	fmt.Printf("\n%s\n%s", massage, prettier.Red)
	fmt.Println("\n", prettier.Reset)
	log.Fatal(err)
}

// Print end of the programm and show exeqution time.
func (prettier *Prettier) End() {
	if entryTimer, ok := prettier.timer.(time.Time); ok {
		fmt.Printf("%s\nExeqution time in seconds: %f\n", prettier.Cyan, time.Now().Sub(entryTimer).Seconds())
	}
	fmt.Printf("%s\n====================END====================\n\n%s", prettier.Green, prettier.Reset)
}
