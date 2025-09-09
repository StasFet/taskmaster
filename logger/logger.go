package logger

import (
	"log"
	"os"
)

var outputFile, _ = os.OpenFile("latest.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

var DB = log.New(outputFile, "[DB] ", log.Lshortfile|log.Ltime|log.LUTC)
var API = log.New(outputFile, "[API] ", log.Lshortfile|log.Ltime|log.LUTC)
var GEN = log.New(outputFile, "[GEN] ", log.Lshortfile|log.Ltime|log.LUTC)