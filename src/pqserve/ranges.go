package main

import (
	"fmt"
	"math"
	"time"
)

const (
	dr_hour = iota
	dr_day
	dr_month
	dr_year
	dr_dec
	dr_cent
)

var (
	dagen   = []string{"zo", "ma", "di", "wo", "do", "vr", "za"}
	maanden = []string{"", "jan", "feb", "maa", "apr", "mei", "jun", "jul", "aug", "sep", "okt", "nov", "dec"}
)

type drange struct {
	min, max time.Time
	r        int
	s        []string
}

type frange struct {
	min, step float64
	s         []string
}

type irange struct {
	min, step int
	s         []string
	indexed   bool
}

func newDrange(min, max time.Time, hasTime bool) *drange {
	dr := drange{
		s: make([]string, 0),
	}

	// tijdzone strippen
	min = time.Date(min.Year(), min.Month(), min.Day(), min.Hour(), min.Minute(), min.Second(), 0, time.UTC)
	max = time.Date(max.Year(), max.Month(), max.Day(), max.Hour(), max.Minute(), max.Second(), 0, time.UTC)

	if max.Year()-min.Year() > 300 /* 300 jaar */ {
		dr.r = dr_cent
	} else {
		dur := max.Sub(min)
		if hasTime && dur < time.Hour*24*4 /* 4 dagen */ {
			dr.r = dr_hour
		} else if dur < time.Hour*24*50 /* 50 dagen */ {
			dr.r = dr_day
		} else if dur < time.Hour*24*365*3 /* 3 jaar */ {
			dr.r = dr_month
		} else if dur < time.Hour*24*365*30 /* 30 jaar */ {
			dr.r = dr_year
		} else {
			dr.r = dr_dec
		}
	}
	switch dr.r {
	case dr_hour:
		dr.min = time.Date(min.Year(), min.Month(), min.Day(), 0, 0, 0, 0, time.UTC)
		dr.max = time.Date(max.Year(), max.Month(), max.Day(), 0, 0, 0, 0, time.UTC)
		dr.max = dr.max.AddDate(0, 0, 1)
		for d := dr.min; d.Before(dr.max); d = d.AddDate(0, 0, 1) {
			wd := dagen[d.Weekday()]
			dag := d.Day()
			maand := maanden[d.Month()]
			jaar := d.Year()
			for h := 0; h < 24; h++ {
				dr.s = append(dr.s, fmt.Sprintf("%s %2d %s %d %02d:00-%02d:59", wd, dag, maand, jaar, h, h))
			}
		}
	case dr_day:
		dr.min = time.Date(min.Year(), min.Month(), min.Day(), 0, 0, 0, 0, time.UTC)
		dr.max = time.Date(max.Year(), max.Month(), max.Day(), 0, 0, 0, 0, time.UTC)
		dr.max = dr.max.AddDate(0, 0, 1)
		for d := min; d.Before(dr.max); d = d.AddDate(0, 0, 1) {
			dr.s = append(dr.s, fmt.Sprintf("%s %2d %s %d", dagen[d.Weekday()], d.Day(), maanden[d.Month()], d.Year()))
		}
	case dr_month:
		dr.min = time.Date(min.Year(), min.Month(), 1, 0, 0, 0, 0, time.UTC)
		dr.max = time.Date(max.Year(), max.Month(), 1, 0, 0, 0, 0, time.UTC)
		dr.max = dr.max.AddDate(0, 1, 0)
		for d := min; d.Before(dr.max); d = d.AddDate(0, 1, 0) {
			dr.s = append(dr.s, fmt.Sprintf("%s %d", maanden[d.Month()], d.Year()))
		}
	case dr_year:
		dr.min = time.Date(min.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.max = time.Date(max.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.max = dr.max.AddDate(1, 0, 0)
		for y := min.Year(); y < dr.max.Year(); y++ {
			dr.s = append(dr.s, fmt.Sprint(y))
		}
	case dr_dec:
		dr.min = time.Date(10*(min.Year()/10), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.max = time.Date(10*(max.Year()/10), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.max = dr.max.AddDate(10, 0, 0)
		for y := dr.min.Year(); y < dr.max.Year(); y += 10 {
			dr.s = append(dr.s, fmt.Sprintf("%d–%d", y, y+9))
		}
	case dr_cent:
		dr.min = time.Date(100*(min.Year()/100), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.max = time.Date(100*(max.Year()/100), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.max = dr.max.AddDate(100, 0, 0)
		for y := dr.min.Year(); y < dr.max.Year(); y += 100 {
			dr.s = append(dr.s, fmt.Sprintf("%d–%d", y, y+99))
		}
	}

	return &dr
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
	if ir.step == 1 && max-min > 1000 {
		return &ir
	}
	ir.indexed = true
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

func (dr *drange) value(val time.Time) (string, int) {
	if val.Before(dr.min) || dr.max.Before(val) {
		return "UNDEF", -1
	}
	i := -1
	switch dr.r {
	case dr_hour:
		for d := dr.min; d.Before(dr.max); d = d.Add(time.Hour) {
			if val.Before(d) {
				break
			}
			i++
		}
	case dr_day:
		for d := dr.min; d.Before(dr.max); d = d.AddDate(0, 0, 1) {
			if val.Before(d) {
				break
			}
			i++
		}
	case dr_month:
		for d := dr.min; d.Before(dr.max); d = d.AddDate(0, 1, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	case dr_year:
		for d := dr.min; d.Before(dr.max); d = d.AddDate(1, 0, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	case dr_dec:
		for d := dr.min; d.Before(dr.max); d = d.AddDate(10, 0, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	case dr_cent:
		for d := dr.min; d.Before(dr.max); d = d.AddDate(100, 0, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	}
	if i < 0 || i >= len(dr.s) {
		return "UNDEF", -1
	}
	return dr.s[i], i
}

func (fr *frange) value(val float64) (string, int) {
	i := int((val - fr.min) / fr.step)
	if i < 0 || i >= len(fr.s) {
		return "UNDEF", -1
	}
	return fr.s[i], i
}

func (ir *irange) value(val int) (string, int) {
	if !ir.indexed {
		return fmt.Sprint(val), val
	}
	i := (val - ir.min) / ir.step
	if i < 0 || i >= len(ir.s) {
		return "UNDEF", -1
	}
	return ir.s[i], i
}

func printDate(t time.Time, hasTime bool) string {
	if hasTime {
		return fmt.Sprintf("%s %2d %s %d %02d:%02d", dagen[t.Weekday()], t.Day(), maanden[t.Month()], t.Year(), t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%s %2d %s %d", dagen[t.Weekday()], t.Day(), maanden[t.Month()], t.Year())
}
