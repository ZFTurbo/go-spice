package prettier

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Purple = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"

// Start your project and start time clock.
func Start(projectName string, projectVersion string, projectAuthor string) {
	fmt.Printf("%s\n==================== %s ====================\n", Green, projectName)
	fmt.Printf("Version: %s\n", projectVersion)
	fmt.Printf("Author: %s\n", projectAuthor)
	fmt.Printf("Licence: MIT %s\n\n", Reset)
}

// Show time from timer start
func ShowTime(timer ...time.Time) {
	if len(timer) > 0 {
		fmt.Printf("%s\nTime past from start: %f\n", Cyan, time.Now().Sub(timer[0]).Seconds())
	}
}

// Create default progress bar
func DefaultBar(steps int, barName string) *progressbar.ProgressBar {
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
func Info(messages map[string]interface{}) {
	fmt.Println(Yellow + "----Info----")

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

	fmt.Println("------------" + Reset)
	fmt.Println()
}

// Print error message.
func Error(massage string, err error) {
	fmt.Printf("\n%s\n%s", massage, Red)
	fmt.Println("\n", Reset)

	if err != nil {
		log.Fatal(err)
	}
}

// Print end of the programm and show exeqution time.
func End(timer ...time.Time) {
	if len(timer) > 0 {
		fmt.Printf("%s\n\nExeqution time in seconds: %f", Green, time.Now().Sub(timer[0]).Seconds())
	}
	fmt.Printf("%s\n====================END====================\n\n%s", Green, Reset)
}
