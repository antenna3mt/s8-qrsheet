package main

import (
	"time"
	"fmt"
	"strings"
	"os"
	"sync"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"image"
	"image/color"
	"image/png"
)

type Code struct {
	spec  int
	label int
	date  time.Time
	page  int
	num   int
}

func (code Code) String() string {
	return fmt.Sprintf("%03d-%02d-%s-%05d", code.spec, code.label, code.date.Format("060102"), code.page*100+code.num)
}

func (code Code) Dirname() string {
	return QrcodeDirname(code.date, code.label)
}

func (code Code) Filename() string {
	return code.String() + ".png"
}

func (code Code) Pathname() string {
	return strings.Join([]string{code.Dirname(), code.Filename()}, "/")
}

func (code Code) writeFile() {
	qrCode, _ := qr.Encode(code.String(), qr.H, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)

	file, _ := os.Create(code.Pathname())
	defer file.Close()

	img := image.NewRGBA(qrCode.Bounds())
	for i := 0; i < qrCode.Bounds().Max.X; i++ {
		for j := 0; j < qrCode.Bounds().Max.Y; j++ {
			r, g, b, a := qrCode.At(i, j).RGBA()
			img.SetRGBA(i, j, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	png.Encode(file, img)
}

func QrcodeDirname(date time.Time, label int) string {
	return strings.Join([]string{QrcodeBaseDirName, fmt.Sprintf("%s-%02d", date.Format("2006-01-02"), label)}, "/")
}

func CodeGenerator(date time.Time, label int) chan Code {
	ch := make(chan Code, 100)
	go func() {
		for page := 1; page <= PdfPage; page++ {
			for num := 1; num <= 16; num++ {
				ch <- Code{Spec, label, date, page, num}
			}
		}
		close(ch)
	}()
	return ch
}

func genQrcodes(date time.Time, label int) {
	ch := CodeGenerator(date, label)
	os.Mkdir(QrcodeBaseDirName, os.ModePerm)
	os.Mkdir(QrcodeDirname(date, label), os.ModePerm)

	count := SafeCounter{n: 0}
	total := PdfPage * 16

	fmt.Println("正在生成二维码...")

	var wg sync.WaitGroup
	wg.Add(QrcodeGenThreads + 1)
	for i := 0; i < QrcodeGenThreads; i++ {
		go func() {
			defer wg.Done()
			for code := range ch {
				code.writeFile()
				go count.Inc()
			}
		}()
	}

	go func() {
		defer wg.Done()
		for (count.Value() < total) {
			fmt.Printf("已完成 %d%% (%d/%d)\n", 100*count.Value()/total, count.Value(), total)
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
	fmt.Println("已生成所有二维码")
}

//func genQrcode(date time.Time, label int, code string) {
//	path := qrcodePath(date, label, code)
//	fmt.Println(path)
//	//qrcode.WriteFile(code, qrcode.Highest, 256, path)
//}
//
//func genQecodes(date time.Time, label int, ch chan string, thread int) {
//	os.Mkdir(qrcodesDir(date, label), os.ModePerm)
//	var wg sync.WaitGroup
//	wg.Add(thread)
//	for i := 0; i < thread; i++ {
//		go func() {
//			defer wg.Done()
//			for code := range ch {
//				genQrcode(date, label, code)
//			}
//		}()
//	}
//	wg.Wait()
//}
