package main

import (
	"fmt"
	"math"
)

type frange struct {
	min, step float64
	s         []string
}

type irange struct {
	min, step int
	s         []string
}

func newFrange(min, max float64) *frange {
	fr := frange{
		step: math.Pow(10, math.Floor(math.Log10(float64(max-min))-.5)) / 5,
		s:    make([]string, 0),
	}
	fr.min = fr.step * math.Floor(min/fr.step)
	for i, f := 0, fr.min; f <= max; i++ {
		f = fr.min + float64(i)*fr.step
		fr.s = append(fr.s, fmt.Sprintf("%g – %g", float32(f), float32(f+fr.step)))
	}
	fr.s = append(fr.s, fr.s[len(fr.s)-1])
	return &fr
}

func newIrange(min, max, count int) *irange {
	ir := irange{
		min:  min,
		step: int(math.Pow(10, math.Floor(math.Log10(float64(max-min))-.5))),
		s:    make([]string, 0),
	}
	if count <= 20 || ir.step < 1 {
		ir.step = 1
	}
	if ir.step >= 100 {
		ir.step /= 5
	} else if ir.step == 10 {
		ir.step /= 2
	}
	if min >= 0 {
		ir.min = ir.step * (min / ir.step)
	} else {
		ir.min = ir.step * (min/ir.step - 1)
	}
	ln := len(fmt.Sprint(ir.min))
	if l := len(fmt.Sprint(max)); l > ln {
		ln = l
	}
	f := fmt.Sprintf("%%%dd", ln)
	if ir.step > 1 {
		f = f + " – " + f
	}
	for i := ir.min; i <= max; i += ir.step {
		if ir.step == 1 {
			ir.s = append(ir.s, fmt.Sprintf(f, i))
		} else {
			ir.s = append(ir.s, fmt.Sprintf(f, i, i+ir.step-1))
		}
	}
	return &ir
}

func (fr *frange) value(val float64) (string, int) {
	i := int((val - fr.min) / fr.step)
	if i < 0 || i >= len(fr.s) {
		return "UNDEF", -1
	}
	return fr.s[i], i
}

func (ir *irange) value(val int) (string, int) {
	i := (val - ir.min) / ir.step
	if i < 0 || i >= len(ir.s) {
		return "UNDEF", -1
	}
	return ir.s[i], i
}
