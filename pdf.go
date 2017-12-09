package main

import (
	"time"
	"fmt"
	"strings"
	"os"
	"github.com/signintech/gopdf"
	"log"
	"bytes"
)

type CodeSheet struct {
	label int
	date  time.Time
}

func (sheet CodeSheet) Filename() string {
	return fmt.Sprintf("%s-%02d", sheet.date.Format("2006-01-02"), sheet.label) + ".pdf"
}

func (sheet CodeSheet) Pathname() string {
	return strings.Join([]string{PdfDirName, sheet.Filename()}, "/")
}

func (sheet CodeSheet) writeFile() {
	const TopMargin = 21.0
	const LeftMargin = 18.0
	const CellHeight = 200.0
	const CellWidth = 140.0

	const ImageTopMargin = 20.0
	const ImageLeftMargin = 20.0
	const ImageWidth = 100.0
	const ImageHeight = 100.0

	const CodeTextTopMargin = 120.0
	const CodeTextHeight = 26.0

	const TextTopMargin = 144.0
	const TextLeftMargin = 20.0
	const TextHeight = 25.0

	getPos := func(num int) (float64, float64) {
		w := (num - 1) % 4
		h := (num - 1) / 4
		return LeftMargin + float64(w)*CellWidth, TopMargin + float64(h)*CellHeight
	}

	os.Mkdir(PdfDirName, os.ModePerm)
	var err error

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595, H: 842}}) //595.28, 841.89 = A4

	fontData, err := Asset("data/msyh.ttf")
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(fontData)
	err = pdf.AddTTFFontByReader("Microsoft YaHei", r)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = pdf.SetFont("Microsoft YaHei", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	for page := 1; page <= PdfPage; page++ {
		pdf.AddPage()
		for num := 1; num <= 16; num++ {
			code := Code{Spec, sheet.label, sheet.date, page, num}
			baseX, baseY := getPos(num)

			// image
			err = pdf.Image(code.Pathname(), baseX+ImageLeftMargin, baseY+ImageTopMargin, &gopdf.Rect{ImageWidth, ImageHeight})
			if err != nil {
				log.Print(err.Error())
				return
			}

			// code
			pdf.SetX(baseX)
			pdf.SetY(baseY + CodeTextTopMargin)
			pdf.SetFont("Microsoft YaHei", "", 10)
			err = pdf.CellWithOption(&gopdf.Rect{CellWidth, CodeTextHeight}, code.String(), gopdf.CellOption{Align: gopdf.Middle | gopdf.Center})
			if err != nil {
				log.Print(err.Error())
				return
			}

			// text 1
			pdf.SetX(baseX + TextLeftMargin)
			pdf.SetY(baseY + TextTopMargin)
			pdf.SetFont("Microsoft YaHei", "", 12)
			err = pdf.CellWithOption(&gopdf.Rect{CellWidth, TextHeight}, "米数：____________", gopdf.CellOption{Align: gopdf.Middle | gopdf.Left})
			if err != nil {
				log.Print(err.Error())
				return
			}

			// text 1
			pdf.SetX(baseX + TextLeftMargin)
			pdf.SetY(baseY + TextTopMargin + TextHeight)
			pdf.SetFont("Microsoft YaHei", "", 12)
			err = pdf.CellWithOption(&gopdf.Rect{CellWidth, TextHeight}, "花号：____________", gopdf.CellOption{Align: gopdf.Middle | gopdf.Left})
			if err != nil {
				log.Print(err.Error())
				return
			}

		}
	}

	pdf.WritePdf(sheet.Pathname())

}

func genPDF(date time.Time, label int) {
	sheet := CodeSheet{label, date}
	fmt.Printf("正在生成PDF文档 %s ...\n", sheet.Filename())
	sheet.writeFile()
	fmt.Println("已完成PDF文档")
}

//qrcode.WriteFile("https://example.org", qrcode.Highest, 256, "qr.png")
//pdf := gofpdf.New("P", "mm", "A4", "")
//pdf.AddPage()
//pdf.SetFont("Arial", "B", 16)
//pdf.Cell(40, 10, "Hello, world")
//err := pdf.OutputFileAndClose("hello.pdf")
//
//fmt.Println(err);
