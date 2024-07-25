package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func ClearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Print("\033[H\033[2J")
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported platform")
	}
}

func ShowMenu() {
	fmt.Println("======================= WELCOME TO LOG TOOL MENU =====================")
	fmt.Println("Select Menu")
	fmt.Println("1. List Log File")
	fmt.Println("2. Search Request and Response Log")
	fmt.Println("3. Download Log")
	fmt.Println("4. Print Latest Log")
	fmt.Println("5. Check ENV")
	fmt.Println("0. Exit")
	fmt.Println("To show help, type 'help'")
}

func ShowHelp() {
	fmt.Println("================================ HELP ================================")
	fmt.Println("List Command")
	fmt.Println("[up] : to exit from current menu")
	fmt.Println("[clear] : to clear screen")
	fmt.Println("[show] : to show menu")
}

func PrintSeparator(maxLength int) {
	fmt.Println(strings.Repeat("-", maxLength))
}
