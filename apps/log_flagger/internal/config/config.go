package config

import "os"

var (
	PATH_OPEN_BOT_DATA  = ""
	PATH_OUTPUT         = ""
	PATH_INPUT          = ""
	CSV_FIELD_CLIENT_IP = 0
	CSV_FIELD_FP_JA4    = 1
	CSV_FIELD_UA        = 2
	DEBUG               = false
	MODE_TEST           = os.Getenv("MODE_TEST") == "1"
	DEBUG_UA            = "___"
)

const VERSION = "1.0"
