package main

var Commands = map[string]string{
	"config":   "Get and set a username.",
	"add":      "Add a file to the index.",
	"log":      "Show commit logs.",
	"commit":   "Save changes.",
	"checkout": "Restore a file.",
}

const HELP_MESSAGE = `
These are SVCS commands:
config     Get and set a username.
add        Add a file to the index.
log        Show commit logs.
commit     Save changes.
checkout   Restore a file.
`

func main() {
	Stage1()
}
