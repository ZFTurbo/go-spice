package utils

import (
	"os"
	"pgsolver/pkg/prettier"
	"strconv"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Opacity(str string, tar string) float64 {
	if str == tar {
		return 0.35
	}
	return 1
}

func ParseInt(str string) int {
	if entry, err := strconv.ParseInt(str, 10, 0); err == nil {
		return int(entry)
	} else {
		prettier.Error("Can't parse int from passed string.", err)
		return 0
	}
}

func ParseFloat(str string) float64 {
	if entry, err := strconv.ParseFloat(str, 64); err == nil {
		return entry
	} else {
		prettier.Error("Can't parse float from passed string.", err)
		return 0
	}
}

func OpenFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		prettier.Error("Can't open file "+fileName+".", err)
		return nil
	}
	return file
}

func CreateFolder(folderPath string) {
	if fail, err := exists(folderPath); err == nil {
		if !fail {
			err = os.Mkdir(folderPath, os.ModePerm)
			if err != nil {
				prettier.Error("Can't create folder in derectory "+folderPath+".", err)
			}
		}
	} else {
		prettier.Error("Can't verfy directory existance.", err)
	}
}

func CreateFile(filePath string) *os.File {
	if file, err := os.Create(filePath); err == nil {
		return file
	} else {
		prettier.Error("Can't create a file at path "+filePath+".", err)
		return nil
	}
}
