package prettier

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	Reset  = ""
	Red    = ""
	Green  = ""
	Yellow = ""
	Blue   = ""
	Purple = ""
	Cyan   = ""
	Gray   = ""
	White  = ""
)

// Start your project and start time clock.
func Start(projectName string, projectVersion string, projectAuthor string, licence string) {
	fmt.Printf("%s\n==================== %s ====================\n", Green, projectName)
	if len(projectVersion) > 0 {
		fmt.Printf("Version: %s\n", projectVersion)
	}
	if len(projectAuthor) > 0 {
		fmt.Printf("Author: %s\n", projectAuthor)
	}
	if len(licence) > 0 {
		fmt.Printf("Licence: %s%s\n\n", licence, Reset)
	}
}

// Show time from timer start
func ShowTime(timer ...time.Time) {
	if len(timer) > 0 {
		fmt.Printf("%s\nTime past from start: %f\n", Cyan, time.Since(timer[0]).Seconds())
	}
}

// Create default progress bar
func DefaultBar(steps int, barName string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(steps,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(barName),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	return bar
}

// Print info block.
func Info(messages map[string]interface{}) {
	var keys []string

	for key := range messages {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	fmt.Println(Yellow + "----Info----")

	for _, key := range keys {
		if entryMessage, ok := messages[key].(string); ok {
			fmt.Println(key + entryMessage)
		}
		if entryMessage, ok := messages[key].(float64); ok {
			if entryMessage > 1e5 || entryMessage < 1e-3 {
				fmt.Println(key + strconv.FormatFloat(entryMessage, 'e', 8, 64))
			} else {
				fmt.Println(key + strconv.FormatFloat(entryMessage, 'f', 8, 64))
			}
		}
		if entryMessage, ok := messages[key].(int); ok {
			fmt.Println(key + strconv.FormatInt(int64(entryMessage), 10))
		}
	}

	fmt.Println("------------" + Reset)
	fmt.Println()
}

func Warrning(message string) {
	fmt.Printf("\n%s\n%s", Yellow, message)
	fmt.Println(Reset)
}

// Print error message and exit if error not nil.
func Error(message string, err error) {
	fmt.Printf("\n%s\n%s", Red, message)
	fmt.Println(Reset)
	if err != nil {
		log.Fatal(err)
	}
}

// Print end of the programm and show exeqution time.
func End(timer ...time.Time) {
	if len(timer) > 0 {
		fmt.Printf("%s\n\nExeqution time in seconds: %f", Green, time.Since(timer[0]).Seconds())
		fmt.Printf("%s\n====================END====================\n\n%s", Green, Reset)
	} else {
		fmt.Printf("%s====================END====================\n\n%s", Green, Reset)
	}
}
