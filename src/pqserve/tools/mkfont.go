package main

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pebbe/util"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"

	"fmt"
	"os"
)

const (
	basesize = 200
)

var (
	x = util.CheckErr
)

func main() {

	names := []string{"fontRegular", "fontBold"}
	fonts := make([]*truetype.Font, 2)

	var err error
	fonts[0], err = truetype.Parse(goregular.TTF)
	x(err)
	fonts[1], err = truetype.Parse(gobold.TTF)
	x(err)

	fmt.Printf("// GENERATED FILE. DO NOT EDIT.\n\npackage main\n\nvar fontBaseSize = %d\n\n", int(basesize))

	for i := 0; i < 2; i++ {

		// The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y. A
		// visual depiction of what these metrics are is at
		// https://developer.apple.com/library/mac/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png

		fc := truetype.NewFace(fonts[i], &truetype.Options{Size: basesize})
		b, _ := font.BoundString(fc, "X")
		fmt.Printf("var %sAscent = %d\n\n", names[i], -b.Min.Y.Round())
		b, _ = font.BoundString(fc, "g")
		fmt.Printf("var %sDescent = %d\n\n", names[i], b.Max.Y.Round())
		fmt.Printf("var %sSizes = []uint8{\n", names[i])
		ctx := freetype.NewContext()
		ctx.SetFont(fonts[i])
		ctx.SetFontSize(basesize)
		for c := rune(0); c < rune(256*256); c++ {
			s := string(c)
			p, err := ctx.DrawString(s, freetype.Pt(0, 0))
			x(err)
			i := p.X.Round()
			if i < 0 || i > 255 {
				fmt.Fprintln(os.Stderr, "rune %d dx = %d\n", int(c), i)
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
