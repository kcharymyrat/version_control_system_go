// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// )

// const (
// 	VCS_DIR          = "./vcs"
// 	CONFIG_FILE_NAME = "config.txt"
// 	INDEX_FILE_NAME  = "index.txt"
// )

// const CONFIG = "config"
// const ADD = "add"
// const LOG = "log"
// const COMMIT = "commit"
// const CHECKOUT = "checkout"
// const HELP = "--help"

// const CommandsText = "These are SVCS commands:"
// const IS_NOT_COMMAND = "'%s' is not a SVCS command.\n"

// const HELP_MESSAGE = `
// These are SVCS commands:
// config     Get and set a username.
// add        Add a file to the index.
// log        Show commit logs.
// commit     Save changes.
// checkout   Restore a file.
// `

// var Commands = map[string]string{
// 	CONFIG:   "Get and set a username.",
// 	ADD:      "Add a file to the index.",
// 	LOG:      "Show commit logs.",
// 	COMMIT:   "Save changes.",
// 	CHECKOUT: "Restore a file.",
// 	HELP:     HELP_MESSAGE,
// }

// const whoAreYou = "Please, tell me who you are."
// const usernameIs = "The username is %s.\n"

// const addFileToIndex = "Add a file to the index."
// const trackedFiles = "Tracked files:"
// const fileIsTracked = "The file '%s' is tracked.\n"
// const canNotFindFile = "Can't find '%s'.\n"

// func CreateVcsDir() {
// 	err := os.MkdirAll(VCS_DIR, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func getConsoleInput() []string {
// 	// config := flag.String("name", "", Commands["config"])
// 	// add := flag.String("add", "", Commands["add"])
// 	// log := flag.String("log", "", Commands["log"])
// 	// commit := flag.String("commit", "", Commands["commit"])
// 	// checkout := flag.String("checkout", "", Commands["checkout"])

// 	if len(os.Args) == 1 {
// 		return []string{os.Args[0], HELP}
// 	}
// 	return os.Args
// }

// func Interaction() {
// 	consoleArgs := getConsoleInput()
// 	command := consoleArgs[1]
// 	if description, exists := Commands[command]; exists {
// 		CommandSwitchCases(command, description, consoleArgs)
// 	} else {
// 		fmt.Printf(IS_NOT_COMMAND, command)
// 	}
// }

// func CommandSwitchCases(command string, description string, consoleArgs []string) {
// 	switch command {
// 	case CONFIG:
// 		configCase(consoleArgs)
// 	case ADD:
// 		addCase(consoleArgs)
// 	case LOG:
// 		fmt.Println(description)
// 	case COMMIT:
// 		fmt.Println(description)
// 	case CHECKOUT:
// 		fmt.Println(description)
// 	default:
// 		fmt.Println(description)
// 	}
// }

// func configCase(consoleArgs []string) {
// 	configFilePath := VCS_DIR + "/" + CONFIG_FILE_NAME
// 	if len(consoleArgs) < 3 {
// 		file, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)
// 		if err != nil {
// 			fmt.Println(whoAreYou)
// 		}
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			fmt.Printf(usernameIs, scanner.Text())
// 			return
// 		}
// 		fmt.Println(whoAreYou)
// 	} else {
// 		// set their name or output an already existing name
// 		configFile := getOpenFileToWriteOnlyOrCreate(configFilePath)
// 		defer configFile.Close()
// 		if err := os.WriteFile(configFilePath, []byte(consoleArgs[2]), 0644); err != nil {
// 			log.Fatal(err) // exit the program if we have an unexpected error
// 		}
// 		fmt.Printf(usernameIs, consoleArgs[2])
// 	}
// }

// func addCase(consoleArgs []string) {
// 	indexFilePath := VCS_DIR + "/" + INDEX_FILE_NAME
// 	if len(consoleArgs) < 3 {
// 		file, err := os.OpenFile(indexFilePath, os.O_RDONLY|os.O_CREATE, 0644)
// 		if err != nil {
// 			fmt.Println(addFileToIndex)
// 		}
// 		defer file.Close()

// 		numFilesTracked := 0

// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			if numFilesTracked == 0 {
// 				fmt.Println(trackedFiles)
// 			}
// 			fmt.Println(scanner.Text())
// 			numFilesTracked++
// 		}
// 		if numFilesTracked < 1 {
// 			fmt.Println(addFileToIndex)
// 		}
// 	} else {
// 		// Make sure args file exists
// 		_, e := os.Open(consoleArgs[2])
// 		if e != nil {
// 			fmt.Printf(canNotFindFile, consoleArgs[2])
// 			return
// 		}

// 		// Append, Write or Create index.html
// 		indexFile, err := os.OpenFile(indexFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer indexFile.Close()

// 		b, err := fmt.Fprintln(indexFile, consoleArgs[2])
// 		if err != nil {
// 			log.Fatal(err, b)
// 		}
// 		fmt.Printf(fileIsTracked, consoleArgs[2])
// 	}
// }

// func getOpenFileToWriteOnlyOrCreate(path string) *os.File {
// 	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return file
// }

// func main() {
// 	CreateVcsDir()
// 	Interaction()
// }
