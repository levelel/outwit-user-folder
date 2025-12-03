//go:build windows
// +build windows

package rime

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// 获取 macOS/Windows Rime 配置目录
func getRimeDir() string {
	var dir string
	switch runtime.GOOS {
	case "darwin": // macOS
		u, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		dir = filepath.Join(u.HomeDir, "Library/witty")
	case "windows": // Windows
		dir = getOutwitDir()
	default:
		log.Fatalf("Unsupported OS: %s so far", runtime.GOOS)
	}

	return dir
}

func getOutwitDir() string {
	keyPath := `Software\witty\outwit`
	valueName := "RimeUserDir"

	// Get from Windows registry
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err != nil {
		log.Printf("Failed to open registry key: %v\n", err)
		// Fallback to default dir
		return getDefaultOutwitDir()
	}
	defer k.Close()

	rimeUserDir, _, err := k.GetStringValue(valueName)
	if err != nil {
		log.Printf("Failed to read registry value: %v\n", err)
		// Fallback to default dir
		return getDefaultOutwitDir()
	}

	return rimeUserDir
}

func getDefaultOutwitDir() string {
	appData := os.Getenv("APPDATA") // AppData\Roaming
	if appData == "" {
		log.Fatalln("APPDATA environment variable is not set.")
	}
	return filepath.Join(appData, "witty")
}
