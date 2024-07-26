package kernel

import (
	"Tools/config"
	"Tools/controller"
	"Tools/service"
	"Tools/util"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func StartUp() {
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logService := service.NewLogService(cfg)
	logging := controller.NewLogController(cfg, logService)

	util.ShowMenu()
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("$ ")

		if !scanner.Scan() {
			fmt.Println("Failed to read input")
			break
		}

		option := scanner.Text()
		if option == "clear" {
			util.ClearScreen()
			continue
		}
		if option == "show" {
			util.ShowMenu()
			continue
		}
		if option == "help" {
			util.ShowHelp()
			continue
		}

		optInt, err := strconv.Atoi(option)
		if err != nil {
			fmt.Println("Input is not valid, try again")
			continue
		}

		switch optInt {
		case 1:
			fmt.Println("=========================== LIST LOG FILE ============================")
			for true {
				fmt.Print("Input Parent Directory File (enter to default folder) : ")
				dirScan := bufio.NewScanner(os.Stdin)
				if dirScan.Scan() {
					dir := dirScan.Text()
					if dir == "" {
						dir = cfg.DefaultFolder
						fmt.Println("Directory : " + dir)
					}
					if dir == "up" {
						break
					}
					if dir == "clear" {
						util.ClearScreen()
						continue
					}
					if dir == "show" {
						util.ShowMenu()
						continue
					}
					if dir == "help" {
						util.ShowHelp()
						continue
					}
					err := logging.ListLogFile(dir)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("Failed to read directory input:", dirScan.Err())
				}
			}
		case 2:
			fmt.Println("=========================== SEARCH LOG ===============================")
			for true {
				fmt.Print("Input Directory File (enter to default folder): ")
				dirScan := bufio.NewScanner(os.Stdin)
				if !dirScan.Scan() {
					fmt.Println("Failed to read directory input:", dirScan.Err())
					continue
				}
				dir := dirScan.Text()
				if dir == "" {
					dir = cfg.DefaultFolder
					fmt.Println("Directory : " + dir)
				}
				if dir == "up" {
					break
				}
				if dir == "clear" {
					util.ClearScreen()
					continue
				}
				if dir == "show" {
					util.ShowMenu()
					continue
				}
				if dir == "help" {
					util.ShowHelp()
					continue
				}

				fmt.Print("Input Keyword: ")
				keywordScan := bufio.NewScanner(os.Stdin)
				if !keywordScan.Scan() {
					fmt.Println("Failed to read keyword input:", keywordScan.Err())
					continue
				}
				keyword := keywordScan.Text()
				if keyword == "up" {
					break
				}
				if keyword == "clear" {
					util.ClearScreen()
					continue
				}
				if keyword == "show" {
					util.ShowMenu()
					continue
				}
				if keyword == "help" {
					util.ShowHelp()
					continue
				}

				err := logging.SearchLog(dir, keyword)
				if err != nil {
					fmt.Println(err)
				}
			}
		case 3:
			fmt.Println("=========================== DOWNLOAD LOG =============================")
			for true {
				fmt.Print("Input Directory File (enter to default folder): ")
				dirScan := bufio.NewScanner(os.Stdin)
				if !dirScan.Scan() {
					fmt.Println("Failed to read directory input:", dirScan.Err())
					continue
				}
				dir := dirScan.Text()
				if dir == "" {
					dir = cfg.DefaultFolder
					fmt.Println("Directory : " + dir)
				}
				if dir == "up" {
					break
				}
				if dir == "clear" {
					util.ClearScreen()
					continue
				}
				if dir == "show" {
					util.ShowMenu()
					continue
				}
				if dir == "help" {
					util.ShowHelp()
					continue
				}
				err := logging.DownloadLog(dir)
				if err != nil {
					fmt.Println(err)
				}
			}
		case 4:
			fmt.Println("========================== PRINT LATEST LOG ==========================")
			for true {
				fmt.Print("Input Directory File (enter to default folder): ")
				dirScan := bufio.NewScanner(os.Stdin)
				if !dirScan.Scan() {
					fmt.Println("Failed to read directory input:", dirScan.Err())
					continue
				}
				dir := dirScan.Text()
				if dir == "" {
					dir = cfg.DefaultFolder
					fmt.Println("Directory : " + dir)
				}
				if dir == "up" {
					break
				}
				if dir == "clear" {
					util.ClearScreen()
					continue
				}
				if dir == "show" {
					util.ShowMenu()
					continue
				}
				if dir == "help" {
					util.ShowHelp()
					continue
				}

				fmt.Print("Input buffer size MB (Default 1MB): ")
				bufferSize := bufio.NewScanner(os.Stdin)
				if !bufferSize.Scan() {
					fmt.Println("Failed to read keyword input:", bufferSize.Err())
					continue
				}
				size := bufferSize.Text()
				if size == "" {
					size = "1"
				}
				if size == "up" {
					break
				}
				if size == "clear" {
					util.ClearScreen()
					continue
				}
				if size == "show" {
					util.ShowMenu()
					continue
				}
				if size == "help" {
					util.ShowHelp()
					continue
				}

				sizeInt, err := strconv.Atoi(size)
				if err != nil {
					fmt.Println("Input must be integer")
					continue
				}

				err = logging.PrintLatestLog(dir, sizeInt)
				if err != nil {
					fmt.Println(err)
				}
			}
		case 5:
			fmt.Println("============================= CHECK ENV ==============================")
			logging.CheckEnv()
		case 0:
			fmt.Println("Exited!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option, please try again")
		}
	}

}
