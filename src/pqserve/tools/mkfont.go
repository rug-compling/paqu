package main

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pebbe/util"
	"golang.org/x/image/font"

	"fmt"
	"io/ioutil"
	"os"
)

const (
	basesize = 134
)

var (
	x = util.CheckErr
)

func main() {

	fonts := [][2]string{
		[2]string{"fontRegular", "FreeSans.ttf"},
		[2]string{"fontBold", "FreeSansBold.ttf"},
	}

	fmt.Printf("// GENERATED FILE. DO NOT EDIT.\n\npackage main\n\nvar fontBaseSize = %d\n\n", int(basesize))

	for i := 0; i < 2; i++ {

		b, err := ioutil.ReadFile(fonts[i][1])
		x(err)
		fnt, err := truetype.Parse(b)
		x(err)

		name := fonts[i][0]

		// The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y. A
		// visual depiction of what these metrics are is at
		// https://developer.apple.com/library/mac/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png

		fc := truetype.NewFace(fnt, &truetype.Options{Size: basesize})

		bs, _ := font.BoundString(fc, "X")
		fmt.Printf("var %sAscent = %d\n\n", fonts[i][0], -bs.Min.Y.Round())

		bs, _ = font.BoundString(fc, "g")
		fmt.Printf("var %sDescent = %d\n\n", name, bs.Max.Y.Round())

		fmt.Printf("var %sSizes = []uint8{\n", name)
		ctx := freetype.NewContext()
		ctx.SetFont(fnt)
		ctx.SetFontSize(basesize)
		for c := rune(0); c < rune(256*256); c++ {
			s := string(c)
			p, err := ctx.DrawString(s, freetype.Pt(0, 0))
			x(err)
			i := p.X.Round()
			if i < 0 || i > 255 {
				fmt.Fprintf(os.Stderr, "rune %d dx = %d\n", int(c), i)
				if i < 0 {
					i = 0
				} else {
					i = 255
				}
			}
			fmt.Printf("%d,", i)
			if c%32 == 31 {
				fmt.Println()
			}
		}
		fmt.Println("}\n")
	}
}
