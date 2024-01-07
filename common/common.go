package common

// Directory and File Configuration
const (
	VCS_DIR          = "./vcs"
	CONFIG_FILE_NAME = "config.txt"
	INDEX_FILE_NAME  = "index.txt"
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
