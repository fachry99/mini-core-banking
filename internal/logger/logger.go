package logger

import (
	"log"
	"os"
)

var Log = log.New(
	os.Stdout,
	"[MINI-CORE-BANKING] ",
	log.LstdFlags|log.Lshortfile,
)
