package main

import (
	"time"
	"os"
)

const QrcodeBaseDirName = "qrcodes"
const PdfDirName = "pdf"
const PdfPage = 500

const Spec = 1
const QrcodeGenThreads = 8

func initEnv() (time.Time, int) {
	date := time.Now()

	for i := 99; i >= 1; i-- {
		path := CodeSheet{i, date}.Pathname()
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return date, i + 1
		}
	}
	return date, 1
}

func main() {
	date, label := initEnv()
	genQrcodes(date, label)
	genPDF(date, label)
}
