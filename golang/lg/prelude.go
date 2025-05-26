package lg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var isFileLogged = false

const RED = "\033[91m"
const GREEN = "\u001b[32m"

func InitLogDisc() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return "c"
	} else {
		// On Windows, the disk is typically the first character followed by ":"
		// On Unix-like systems, you might get a volume or mount point
		fmt.Println("Current working directory:", cwd)
		if len(cwd) >= 2 && cwd[1] == ':' {
			fmt.Println("Disk from CWD:", strings.ToUpper(cwd[:1]))
			return strings.ToUpper(cwd[:1])
		}
	}

	return "c"
}

var DISC = InitLogDisc()

func InitFileLogging() *os.File {
	os.MkdirAll("./logs", 0755)

	isFileLogged = true

	currentTime := time.Now()

	// Format the filename with the current date and time up to seconds
	fileName := currentTime.Format("2006-01-02_15-04-05") + ".log"

	file, err := os.OpenFile("logs/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogE("Failed to open file: ", err)
		return nil
	}

	// Save the original os.Stdout for restoring later
	_ = os.Stdout //Return original loggin

	// Redirect os.Stdout to the file
	os.Stdout = file

	return file
}

func LogI(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	print("\033[94m", "Info", file, line, logs...)
}

func LogW(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	print("\u001b[35m", "Info", file, line, logs...)

}

func LogOk(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	print(GREEN, "Ok", file, line, logs...)

}

func LogEson(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	for _, log := range logs {

		// Pretty-print JSON
		prettyJSON, err := json.MarshalIndent(log, "", "    ")

		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		print(RED, "Error", file, line, string(prettyJSON))

	}

}

func LogOkson(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	for _, log := range logs {

		// Pretty-print JSON
		prettyJSON, err := json.MarshalIndent(log, "", "    ")

		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		print(GREEN, "Ok", file, line, string(prettyJSON))

	}

}

func LogERaw(logs ...any) {
	printraw("\033[91m", logs...)
}

func LogE(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	print("\033[91m", "Error", file, line, logs...)

}

func Log(logs ...any) {
	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	printNoColor("Log", file, line, logs...)

}

func LogStackTraceErr(logs ...any) {

	// Get the caller's information
	LogERaw("\t\t\t\t\t\t\t\t\t--------------------[StackTrace]----------------------")
	for i, index := range []int{1, 2, 3, 4, 5, 6, 7, 8 /* 9, 10, 11, 12 */} {
		_, file, line, ok := runtime.Caller(index) // 1 means get the caller of this function
		if ok {
			if i == 0 {
				print("\033[91m", "StackTrace", file, line, logs...)
			} else {
				print("\033[91m", "StackTrace", file, line)
			}
		} else {
			break
		}
	}
	LogERaw("\t\t\t\t\t\t\t\t\t--------------------[StackTrace]----------------------")
}

func print(color string, level string, file string, line int, logs ...any) {

	if runtime.GOOS == "linux" {
		file = strings.Replace(file, fmt.Sprintf("/mnt/%s", DISC), fmt.Sprintf("%s:", DISC), 1)
	}

	if isFileLogged {
		fmt.Print(file, ":", line, " [", level, "]:")
		for _, log := range logs {
			fmt.Print(" ", log)
		}
		fmt.Print("\n")
	} else {

		fmt.Print(color, file, ":", line, " [", level, "]:")
		for _, log := range logs {
			if log == "" || log == " " {
				fmt.Print("█")
			} else if log == "\n" {
				fmt.Println("\\n")
			} else {
				fmt.Print(" ", log)
			}
		}
		fmt.Print("\033[0m \n")
	}

}

func printNoColor(level string, file string, line int, logs ...any) {

	if runtime.GOOS == "linux" {
		file = strings.Replace(file, fmt.Sprintf("/mnt/%s", DISC), fmt.Sprintf("%s:", DISC), 1)
	}

	if isFileLogged {
		fmt.Print(file, ":", line, " [", level, "]:")
		for _, log := range logs {
			fmt.Print(" ", log)
		}
		fmt.Print("\n")
	} else {

		fmt.Print(file, ":", line, " [", level, "]:")
		for _, log := range logs {
			if log == "" || log == " " {
				fmt.Print("█")
			} else if log == "\n" {
				fmt.Println("\\n")
			} else {
				fmt.Print(" ", log)
			}
		}
		fmt.Print("\n")
	}

}

func printraw(color string, logs ...any) {
	fmt.Print(color)
	for _, log := range logs {
		fmt.Print(" ", log)
	}
	fmt.Print("\033[0m \n")
}

func LogBreakPoint(logs ...any) {

	// Get the caller's information
	_, file, line, ok := runtime.Caller(1) // 1 means get the caller of this function
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}

	print("\u001b[35m", "Info", file, line, logs...)
	fmt.Println("Press enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	// Wait for the user to press Enter
	_, _ = reader.ReadString('\n')

}
