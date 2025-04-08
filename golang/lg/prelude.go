package lg

import (
	"encoding/json"
	"fmt"
	"runtime"
)

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

	print("\u001b[32m", "Ok", file, line, logs...)

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

		print("\033[91m", "Error", file, line, string(prettyJSON))

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

func LogStackTraceErr(logs ...any) {

	// Get the caller's information
	LogERaw("\t\t\t\t\t\t\t\t\t--------------------[Error]----------------------")
	for i := range []int{2, 3, 4, 5, 6, 7, 8} {
		_, file, line, ok := runtime.Caller(i) // 1 means get the caller of this function
		if ok {
			if i == 0 {
				print("\033[91m", "Error", file, line, logs...)
			} else {
				print("\033[91m", "Error", file, line)
			}
		} else {
			return
		}
	}
	LogERaw("\t\t\t\t\t\t\t\t\t--------------------[End]----------------------")
}

func print(color string, level string, file string, line int, logs ...any) {
	fmt.Print(color, file, ":", line, ": [", level, "]:")
	for _, log := range logs {
		fmt.Print(" ", log)
	}
	fmt.Print("\033[0m \n")
}

func printraw(color string, logs ...any) {
	fmt.Print(color)
	for _, log := range logs {
		fmt.Print(" ", log)
	}
	fmt.Print("\033[0m \n")
}
