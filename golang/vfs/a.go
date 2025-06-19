package vfs

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"turtle/credentials"
	"turtle/lg"
)

func GetEnvOrDefault(key string, defaultValue string) string {
	var value = os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetExeFile() string {
	exePath, err := os.Executable()

	if err != nil {
		return "failedtogetworkdir"
	} else {
		return exePath
	}
}

func IsInDevelopment() bool {
	exePath, err := os.Executable()

	if IsLinux() {
		return false
	}

	if err != nil {
		lg.LogI("Error:", err)
		return false
	}
	exeDir := filepath.Dir(exePath)

	lg.LogW(exeDir)

	// Check if the binary is running from a temporary directory
	if strings.Contains(exeDir, "\\AppData\\") ||
		strings.Contains(exeDir, "src-go") ||
		strings.Contains(exeDir, "go-build") {
		return true
	} else {
		return false
	}
}

// FindAllFilesWithExtension - Finds all files with the given extension in a directory (optionally recursive)
func FindAllFilesWithExtension(folderPath string, extension string, recursive bool) ([]string, error) {
	var fileList []string
	if recursive {
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), "."+extension) {
				fileList = append(fileList, path)
			}
			return nil
		})
		return fileList, err
	}

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "."+extension) {
			fileList = append(fileList, filepath.Join(folderPath, file.Name()))
		}
	}
	return fileList, nil
}

// FindAllFoldersWithExtension - Finds all folders containing files with the given extension
func FindAllFoldersWithExtension(folderPath string, extension string) ([]string, error) {
	var folderList []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "."+extension) {
			folderList = append(folderList, filepath.Dir(path))
		}
		return nil
	})
	return folderList, err
}

// GetFileName - Returns file name from a path, optionally including the extension
func GetFileName(path string, includeExtension bool) string {
	if includeExtension {
		return filepath.Base(path)
	}
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// Exists - Checks if the path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileFolder - Returns the folder for a given file path
func GetFileFolder(filePath string) string {
	return filepath.Dir(filePath)
}

// FindPreviewFromPostfix - Finds preview image by appending a postfix
func FindPreviewFromPostfix(filePath string, postFix string) string {
	folder := GetFileFolder(filePath)
	name := GetFileName(filePath, false)
	searched := filepath.Join(folder, name+postFix)
	if Exists(searched) {
		return searched
	}
	return ""
}

// GetPreviewPathInFolder - Returns the path of a preview.png in a folder
func GetPreviewPathInFolder(folder string) string {
	path := filepath.Join(folder, "preview.png")
	if Exists(path) {
		return path
	}
	return ""
}

// GetFirstFileWithExtension - Returns the first file with the given extension in a folder
func GetFirstFileWithExtension(folderPath string, extension string) string {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), "."+extension) {
			return filepath.SkipDir
		}
		return nil
	})
	if err == nil {
		return ""
	}
	return err.Error()
}

// IsLinux - Check if the system is Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// GetWorkingDirectory - Returns the appropriate working directory based on OS
func GetWorkingDirectory() string {
	if IsLinux() {
		return credentials.LinuxWorkspace()
	}
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "TurtleEngine")
}

func GetDarwinWorkspace() string {

	fromEnv := GetEnvOrDefault("INFINITY_TWIN_DARWIN_WORKSPACE", "")

	if fromEnv == "" {
		usr, err := user.Current()
		appSupportPath := filepath.Join(usr.HomeDir, "Library", "Application Support")
		if err == nil {
			return appSupportPath + "/" + "TurtleEngine"
		} else {
			lg.LogE(err)
		}

	} else {
		return fromEnv + "/" + "TurtleEngine"
	}

	return "../infinity_twin_storage" + "/" + "TurtleEngine"

}

// GetFileFolderNew - Similar to GetFileFolder but normalizes path
func GetFileFolderNew(filePath string) string {
	return filepath.Dir(filepath.Clean(filePath))
}

// WriteFileToWD - Writes bytes to a file in the working directory
func WriteFileToWD(folder, filePath string, data []byte) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder, filePath)

	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, data, 0644)
}

// WriteFileToWDNew - Similar to WriteFileToWD but without a folder prefix
func WriteFileToWDNew(filePath string, data []byte) error {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, data, 0644)
}

// DeleteFolderNew - Deletes a folder in the working directory
func DeleteFolderNew(folderPath string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folderPath)
	return os.RemoveAll(finalPath)
}

// DeleteFolderNew - Deletes a folder in the working directory
func DeleteFolder(folderPath string, subfolder string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folderPath, subfolder)
	return os.RemoveAll(finalPath)
}

// DeleteFileNew - Deletes a file in the working directory
func DeleteFileNew(filePath string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	return os.Remove(finalPath)
}

// MakeDirs - Creates a directory (recursive)
func MakeDirs(folder string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder)
	return os.MkdirAll(finalPath, 0755)
}

// OpenWDFolder - Opens a folder in the working directory (on Windows)
func OpenWDFolder(folder string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folder)
	cmd := exec.Command("explorer", finalPath)
	return cmd.Start()
}

// WriteFileStringToWD - Writes a string to a file in the working directory
func WriteFileStringToWD(folder, filePath, content string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder, filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, []byte(content), 0644)
}

func WriteFileStringToWDNew(filePath, content string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, []byte(content), 0644)
}

func WriteFileStringToWDNewBase64(filePath, content string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, filePath)

	// Decode Base64 content
	decodedContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return fmt.Errorf("Failed to decode Base64 content: %w", err)
	}

	// Create necessary directories
	err = os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories for path %s: %w", finalPath, err)
	}

	// Write the decoded content to the file
	err = os.WriteFile(finalPath, decodedContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", finalPath, err)
	}

	return nil
}

// ListFiles - Returns a list of all files in a folder
func ListFiles(folder string) ([]string, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), folder)
	return filepath.Glob(filepath.Join(finalPath, "*"))
}

// GetFileBytesFromWDNew - Reads bytes from a file in the working directory
func GetFileBytesFromWDNew(filePath string) ([]byte, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	return os.ReadFile(finalPath)
}
func GetFileBytesFromWD(folder string, filePath string) ([]byte, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), folder, filePath)
	return os.ReadFile(finalPath)
}

// IsDesktop - Checks if the system is a Windows desktop
func IsDesktop() bool {
	return runtime.GOOS == "windows"
}

// GetFileStringFromWDNew - Reads a string from a file in the working directory
func GetFileStringFromWDNew(filePath string) (string, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	data, err := os.ReadFile(finalPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetFileStringFromWD(folder string, filePath string) (string, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), folder, filePath)
	data, err := os.ReadFile(finalPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetFilePathFromWD(container string, file string) string {
	finalPath := filepath.Join(GetWorkingDirectory(), file)
	return finalPath

}
