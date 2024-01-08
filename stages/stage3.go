package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const (
	VCS_DIR          = "./vcs"
	COMMITS_DIR_NAME = "commits"
	COMMITS_DIR      = VCS_DIR + "/" + COMMITS_DIR_NAME
	CONFIG_FILE_NAME = "config.txt"
	INDEX_FILE_NAME  = "index.txt"
	LOG_FILE_NAME    = "log.txt"
)

const CONFIG = "config"
const ADD = "add"
const LOG = "log"
const COMMIT = "commit"
const CHECKOUT = "checkout"
const HELP = "--help"

const CommandsText = "These are SVCS commands:"
const IS_NOT_COMMAND = "'%s' is not a SVCS command.\n"

const HELP_MESSAGE = `
These are SVCS commands:
config     Get and set a username.
add        Add a file to the index.
log        Show commit logs.
commit     Save changes.
checkout   Restore a file.
`

var Commands = map[string]string{
	CONFIG:   "Get and set a username.",
	ADD:      "Add a file to the index.",
	LOG:      "Show commit logs.",
	COMMIT:   "Save changes.",
	CHECKOUT: "Restore a file.",
	HELP:     HELP_MESSAGE,
}

const configFilePath = VCS_DIR + "/" + CONFIG_FILE_NAME
const indexFilePath = VCS_DIR + "/" + INDEX_FILE_NAME
const logFilePath = VCS_DIR + "/" + LOG_FILE_NAME

const whoAreYou = "Please, tell me who you are."
const usernameIs = "The username is %s."

const addFileToIndex = "Add a file to the index."
const trackedFiles = "Tracked files:"
const fileIsTracked = "The file '%s' is tracked.\n"
const canNotFindFile = "Can't find '%s'.\n"

const noCommitsYet = "No commits yet."
const changesCommited = "Changes are committed."
const messageWasNotPassed = "Message was not passed."
const logMessage = "commit %s\nAuthor: %s\n%s\n\n"
const nothingToCommit = "Nothing to commit."

func createDir(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func getConsoleInput() []string {
	if len(os.Args) == 1 {
		return []string{os.Args[0], HELP}
	}
	return os.Args
}

func Interaction() {
	createDir(VCS_DIR)
	consoleArgs := getConsoleInput()
	command := consoleArgs[1]
	if description, exists := Commands[command]; exists {
		CommandSwitchCases(command, description, consoleArgs)
	} else {
		fmt.Printf(IS_NOT_COMMAND, command)
	}
}

func CommandSwitchCases(command string, description string, consoleArgs []string) {
	switch command {
	case CONFIG:
		configCase(consoleArgs)
	case ADD:
		addCase(consoleArgs)
	case LOG:
		logCase(consoleArgs)
	case COMMIT:
		commitCase(consoleArgs)
	case CHECKOUT:
		fmt.Println(description)
	default:
		fmt.Println(description)
	}
}

func checkFileAndGetSliceOfLines(filePath string) ([]string, error) {
	// Check if the file exists
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err // File does not exist or cannot be opened
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	// Read file line by line
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err // Handle potential scanner errors
	}

	// Check if file is empty
	if len(lines) == 0 {
		return []string{"Empty file"}, nil
	}

	return lines, nil
}

func commitCase(consoleArgs []string) {
	rootPath := "."

	// Check if there is argument passed alongsing the "commit", if not print required message and return
	if len(consoleArgs) < 3 {
		fmt.Println(messageWasNotPassed)
		return
	}

	// Create "commit" directory
	createDir(COMMITS_DIR)

	// Get currently tracked files from index.txt initially - and get trackedFileNames
	trackedFileNames, err := checkFileAndGetSliceOfLines(indexFilePath)
	if err != nil {
		fmt.Println("No Files are tracked now! Use 'add'")
		return
	}

	// Get combined hash of the trackedFiles in current state in the root directory
	hashesStr := getHashesStrOfFilesInDirPath(rootPath, trackedFileNames)

	// Generate a currentCommitIdHash
	currentCommitIdHashStr := createHashedCommitId(hashesStr)

	// Get latest commit Id
	latestCommitId, err := getLatestCommitsDirectoryName(COMMITS_DIR)
	if err != nil {
		fmt.Println("No commits directory yet")
		return
	}

	// If it the same as currentCommit then nothing has changed
	if currentCommitIdHashStr == latestCommitId {
		fmt.Println(nothingToCommit)
		return
	}

	// Else
	// Create new folder inside commits file with currentCommitId
	newPath := COMMITS_DIR + "/" + currentCommitIdHashStr
	createDir(newPath)

	srcDir := rootPath
	dstDir := newPath

	err = copyFiles(srcDir, dstDir, trackedFileNames)
	if err != nil {
		fmt.Println("From src or dst")
	}

	// Create a log for the commit done
	logCommit(currentCommitIdHashStr, consoleArgs)

	fmt.Println(changesCommited)
}

func logCommit(currentCommitIdHashStr string, consoleArgs []string) {
	// get config.txt
	file, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(whoAreYou)
	}
	defer file.Close()

	// get current username from config.txt
	var username string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username = scanner.Text()
	}

	// get log.txt to append or create to write
	logFile := getOpenFileToAppendOrCreateWriteOnly(logFilePath)
	defer logFile.Close()

	// Append to log.txt
	stringToPass := currentCommitIdHashStr + " " + username + " " + consoleArgs[2] + "\n"
	_, err = logFile.WriteString(stringToPass)
	if err != nil {
		log.Fatal(err)
	}

}

func copyFiles(srcDir, dstDir string, trackedFileNames []string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			for _, fileName := range trackedFileNames {
				filePath := srcDir + "/" + entry.Name()
				if entry.Name() == fileName {
					err := copySingleFile(filePath, dstDir+"/"+fileName)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
		}
	}

	return nil
}

func copySingleFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file for writing
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure that any writes to destFile are committed to stable storage
	return destFile.Sync()
}

func getLatestCommitsDirectoryName(path string) (string, error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	var latestDir fs.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return "", err
			}
			if latestDir == nil || info.ModTime().After(latestDir.ModTime()) {
				latestDir = info
			}
		}
	}

	if latestDir != nil {
		return latestDir.Name(), nil
	}
	return "", nil
}

func getHashesStrOfFilesInDirPath(dirPath string, trackedFileNames []string) string {
	var hashesStr string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			for _, fileName := range trackedFileNames {
				filePath := dirPath + "/" + entry.Name()
				if entry.Name() == fileName {
					hashesStr += getMD5HashStrForFile(filePath)
				}
			}
		}
	}

	return hashesStr
}

func getMD5HashStrForFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	md5Hash := md5.New()
	if _, e := io.Copy(md5Hash, file); err != nil {
		log.Fatal(e)
	}

	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func createHashedCommitId(hashesStr string) string {
	// Hash the commit
	commitIdHash := sha256.New()
	commitIdHash.Write([]byte(hashesStr))

	return fmt.Sprintf("%x", commitIdHash.Sum(nil))
}

func logCase(consoleArgs []string) {
	const commit = "commit"
	const author = "author"
	const message = "message"

	// Open the log.txt file
	logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(noCommitsYet)
		return
	}

	// Generate slice of maps from log.txt lines
	var sliceOfMaps []map[string]string
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		lineStr := scanner.Text()
		lineSlice := strings.Split(strings.TrimSpace(lineStr), " ")
		lineMap := map[string]string{
			commit:  lineSlice[0],
			author:  lineSlice[1],
			message: strings.Join(lineSlice[2:], " "),
		}
		sliceOfMaps = append(sliceOfMaps, lineMap)
	}

	// If log.txt is empty
	if len(sliceOfMaps) < 1 {
		fmt.Println(noCommitsYet)
		return
	}

	// Loop from the end and print to the user
	for i := len(sliceOfMaps) - 1; i >= 0; i-- {
		fmt.Printf(logMessage, sliceOfMaps[i][commit], sliceOfMaps[i][author], sliceOfMaps[i][message])
	}
}

func addCase(consoleArgs []string) {
	if len(consoleArgs) < 3 {
		file, err := os.OpenFile(indexFilePath, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(addFileToIndex)
		}
		defer file.Close()

		numFilesTracked := 0

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if numFilesTracked == 0 {
				fmt.Println(trackedFiles)
			}
			fmt.Println(scanner.Text())
			numFilesTracked++
		}
		if numFilesTracked < 1 {
			fmt.Println(addFileToIndex)
		}
	} else {
		// Make sure args file exist
		_, e := os.Open(consoleArgs[2])
		if e != nil {
			fmt.Printf(canNotFindFile, consoleArgs[2])
			return
		}

		// Append, Write or Create index.html
		indexFile, err := os.OpenFile(indexFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer indexFile.Close()

		b, error := fmt.Fprintln(indexFile, consoleArgs[2])
		if error != nil {
			log.Fatal(error, b)
		}
		fmt.Printf(fileIsTracked, consoleArgs[2])
	}
}

func configCase(consoleArgs []string) {

	if len(consoleArgs) < 3 {
		file, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(whoAreYou)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Printf(usernameIs, scanner.Text())
			return
		}
		fmt.Println(whoAreYou)
	} else {
		// set their name or output an already existing name
		configFile := getOpenFileToWriteOnlyOrCreate(configFilePath)
		defer configFile.Close()
		if err := os.WriteFile(configFilePath, []byte(consoleArgs[2]), 0644); err != nil {
			log.Fatal(err) // exit the program if we have an unexpected error
		}
		fmt.Printf(usernameIs, consoleArgs[2])
	}
}

func getOpenFileToWriteOnlyOrCreate(path string) *os.File {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func getOpenFileToAppendOrCreateWriteOnly(path string) *os.File {
	// os.O_APPEND | os.O_CREATE | os.O_WRONLY opens the file in append mode,
	// creates it if it doesn't exist, and opens it in write-only mode.
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func main() {
	Interaction()
}
