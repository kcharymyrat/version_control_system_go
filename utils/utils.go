package utils

import (
	"fmt"
	"log"
	"os"
	"version_control_go/common"
)

func Interaction() {
	command := CmdRead()
	if description, exists := common.Commands[command]; exists {
		fmt.Println(description)
	} else if command == "--help" {
		fmt.Print(common.HELP_MESSAGE)
	} else {
		fmt.Printf("'%v' is not a SVCS command.\n", command)
	}
}

func CmdRead() string {
	// config := flag.String("name", "", Commands["config"])
	// add := flag.String("add", "", Commands["add"])
	// log := flag.String("log", "", Commands["log"])
	// commit := flag.String("commit", "", Commands["commit"])
	// checkout := flag.String("checkout", "", Commands["checkout"])

	if len(os.Args) == 1 {
		return "--help"
	}

	if len(os.Args) > 2 {
		log.Fatal("Error! Expected 2 arguments only!")
	}
	return os.Args[1]

}
