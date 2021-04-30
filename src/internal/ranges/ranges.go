package ranges

import (
	"fmt"
	"math"
	"time"
)

const (
	// DEZE WAARDES NIET MEER VERANDEREN
	Dr_hour = iota
	Dr_day
	Dr_month
	Dr_year
	Dr_dec
	Dr_cent
)

var (
	dagen   = []string{"zo", "ma", "di", "wo", "do", "vr", "za"}
	Maanden = []string{"", "jan", "feb", "maa", "apr", "mei", "jun", "jul", "aug", "sep", "okt", "nov", "dec"}
)

type Drange struct {
	Min, Max time.Time
	R        int
	S        []string
	Indexed  bool
}

type Frange struct {
	Min, Step float64
	S         []string
	Indexed   bool
}

type Irange struct {
	Min, Step int
	S         []string
	Indexed   bool
}

func NewDrange(min, max time.Time, count int, hasTime bool) *Drange {
	// tijdzone strippen
	min = time.Date(min.Year(), min.Month(), min.Day(), min.Hour(), min.Minute(), min.Second(), 0, time.UTC)
	max = time.Date(max.Year(), max.Month(), max.Day(), max.Hour(), max.Minute(), max.Second(), 0, time.UTC)

	/*
		if count < 21 && !hasTime {
			return OldDrange(min, max, size, dr_day, false)
		}
	*/

	var r int
	if max.Year()-min.Year() > 300 /* 300 jaar */ {
		r = Dr_cent
	} else {
		dur := max.Sub(min)
		if hasTime && dur < time.Hour*24*4 /* 4 dagen */ {
			r = Dr_hour
		} else if dur < time.Hour*24*50 /* 50 dagen */ {
			r = Dr_day
		} else if dur < time.Hour*24*365*3 /* 3 jaar */ {
			r = Dr_month
		} else if dur < time.Hour*24*365*30 /* 30 jaar */ {
			r = Dr_year
		} else {
			r = Dr_dec
		}
	}
	return OldDrange(min, max, r, true)
}

func OldDrange(min, max time.Time, dtype int, indexed bool) *Drange {
	// tijdzone strippen
	min = time.Date(min.Year(), min.Month(), min.Day(), min.Hour(), min.Minute(), min.Second(), 0, time.UTC)
	max = time.Date(max.Year(), max.Month(), max.Day(), max.Hour(), max.Minute(), max.Second(), 0, time.UTC)

	dr := Drange{
		Min:     min,
		Max:     max,
		R:       dtype,
		S:       make([]string, 0),
		Indexed: indexed,
	}
	if !indexed {
		return &dr
	}

	switch dr.R {
	case Dr_hour:
		dr.Min = time.Date(min.Year(), min.Month(), min.Day(), 0, 0, 0, 0, time.UTC)
		dr.Max = time.Date(max.Year(), max.Month(), max.Day(), 0, 0, 0, 0, time.UTC)
		dr.Max = dr.Max.AddDate(0, 0, 1)
		for d := dr.Min; d.Before(dr.Max); d = d.AddDate(0, 0, 1) {
			wd := dagen[d.Weekday()]
			dag := d.Day()
			maand := Maanden[d.Month()]
			jaar := d.Year()
			for h := 0; h < 24; h++ {
				dr.S = append(dr.S, fmt.Sprintf("%s %2d %s %d %02d:00-%02d:59", wd, dag, maand, jaar, h, h))
			}
		}
	case Dr_day:
		dr.Min = time.Date(min.Year(), min.Month(), min.Day(), 0, 0, 0, 0, time.UTC)
		dr.Max = time.Date(max.Year(), max.Month(), max.Day(), 0, 0, 0, 0, time.UTC)
		dr.Max = dr.Max.AddDate(0, 0, 1)
		for d := min; d.Before(dr.Max); d = d.AddDate(0, 0, 1) {
			dr.S = append(dr.S, fmt.Sprintf("%s %2d %s %d", dagen[d.Weekday()], d.Day(), Maanden[d.Month()], d.Year()))
		}
	case Dr_month:
		dr.Min = time.Date(min.Year(), min.Month(), 1, 0, 0, 0, 0, time.UTC)
		dr.Max = time.Date(max.Year(), max.Month(), 1, 0, 0, 0, 0, time.UTC)
		dr.Max = dr.Max.AddDate(0, 1, 0)
		for d := min; d.Before(dr.Max); d = d.AddDate(0, 1, 0) {
			dr.S = append(dr.S, fmt.Sprintf("%s %d", Maanden[d.Month()], d.Year()))
		}
	case Dr_year:
		dr.Min = time.Date(min.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.Max = time.Date(max.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.Max = dr.Max.AddDate(1, 0, 0)
		for y := min.Year(); y < dr.Max.Year(); y++ {
			dr.S = append(dr.S, fmt.Sprint(y))
		}
	case Dr_dec:
		dr.Min = time.Date(10*(min.Year()/10), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.Max = time.Date(10*(max.Year()/10), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.Max = dr.Max.AddDate(10, 0, 0)
		for y := dr.Min.Year(); y < dr.Max.Year(); y += 10 {
			dr.S = append(dr.S, fmt.Sprintf("%d–%d", y, y+9))
		}
	case Dr_cent:
		dr.Min = time.Date(100*(min.Year()/100), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.Max = time.Date(100*(max.Year()/100), 1, 1, 0, 0, 0, 0, time.UTC)
		dr.Max = dr.Max.AddDate(100, 0, 0)
		for y := dr.Min.Year(); y < dr.Max.Year(); y += 100 {
			dr.S = append(dr.S, fmt.Sprintf("%d–%d", y, y+99))
		}
	}

	return &dr
}

func NewFrange(min, max float64) *Frange {
	if min == max {
		return OldFrange(min, 0, 1)
	}
	step := math.Pow(10, math.Floor(math.Log10(float64(max-min))-.5)) / 5
	min = step * math.Floor(min/step)
	size := 0
	for i, f := 0, min; f <= max; i++ {
		f = min + float64(i)*step
		size++
	}
	return OldFrange(min, step, size+1)
}

func OldFrange(fmin, fstep float64, size int) *Frange {
	fr := Frange{
		S:       make([]string, 0),
		Indexed: false,
	}
	if size == 1 {
		fr.S = append(fr.S, fmt.Sprintf("%g", float32(fmin)))
		return &fr
	}

	fr.Min = fmin
	fr.Step = fstep
	fr.Indexed = true

	for i := 0; i < size; i++ {
		f := fr.Min + float64(i)*fr.Step
		fr.S = append(fr.S, fmt.Sprintf("%g – %g", float32(f), float32(f+fr.Step)))
	}
	return &fr
}

func NewIrange(min, max, count int) *Irange {
	step := int(math.Pow(10, math.Floor(math.Log10(float64(max-min))-.5)))
	if count <= 20 || step < 1 {
		step = 1
	}
	if step == 1 && max-min > 1000 {
		return OldIrange(min, step, 0, false)
	}

	if step >= 100 {
		step /= 5
	} else if step == 10 {
		step /= 2
	}
	if min >= 0 {
		min = step * (min / step)
	} else {
		min = step * (min/step - 1)
	}
	size := 0
	for i := min; i <= max; i += step {
		size++
	}
	return OldIrange(min, step, size, true)
}

func OldIrange(imin, istep, size int, indexed bool) *Irange {
	ir := Irange{
		Min:     imin,
		Step:    istep,
		Indexed: indexed,
		S:       make([]string, 0),
	}
	if !indexed {
		return &ir
	}

	ln := len(fmt.Sprint(ir.Min))
	if l := len(fmt.Sprint(ir.Min + istep*(size-1))); l > ln {
		ln = l
	}
	f := fmt.Sprintf("%%%dd", ln)
	if ir.Step > 1 {
		f = f + " – " + f
	}
	iv := ir.Min
	for i := 0; i < size; i++ {
		if ir.Step == 1 {
			ir.S = append(ir.S, fmt.Sprintf(f, iv))
		} else {
			ir.S = append(ir.S, fmt.Sprintf(f, iv, iv+ir.Step-1))
		}
		iv += ir.Step
	}

	return &ir
}

func (dr *Drange) Value(val time.Time) (string, int) {
	if !dr.Indexed {
		return fmt.Sprintf("%s %2d %s %d", dagen[val.Weekday()], val.Day(), Maanden[val.Month()], val.Year()),
			val.Day() + 31*int(val.Month()) + 31*12*val.Year()
	}

	// tijdzone strippen
	val = time.Date(val.Year(), val.Month(), val.Day(), val.Hour(), val.Minute(), val.Second(), 0, time.UTC)

	if val.Before(dr.Min) || dr.Max.Before(val) {
		return "UNDEF", -1
	}
	i := -1
	switch dr.R {
	case Dr_hour:
		for d := dr.Min; d.Before(dr.Max); d = d.Add(time.Hour) {
			if val.Before(d) {
				break
			}
			i++
		}
	case Dr_day:
		for d := dr.Min; d.Before(dr.Max); d = d.AddDate(0, 0, 1) {
			if val.Before(d) {
				break
			}
			i++
		}
	case Dr_month:
		for d := dr.Min; d.Before(dr.Max); d = d.AddDate(0, 1, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	case Dr_year:
		for d := dr.Min; d.Before(dr.Max); d = d.AddDate(1, 0, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	case Dr_dec:
		for d := dr.Min; d.Before(dr.Max); d = d.AddDate(10, 0, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	case Dr_cent:
		for d := dr.Min; d.Before(dr.Max); d = d.AddDate(100, 0, 0) {
			if val.Before(d) {
				break
			}
			i++
		}
	}
	if i < 0 || i >= len(dr.S) {
		return "UNDEF", -1
	}
	return dr.S[i], i
}

func (fr *Frange) Value(val float64) (string, int) {
	if !fr.Indexed {
		return fmt.Sprintf("%g", float32(val)), 0
	}

	i := int((val - fr.Min) / fr.Step)
	if i < 0 || i >= len(fr.S) {
		return "UNDEF", 2147483647
	}
	return fr.S[i], i
}

func (ir *Irange) Value(val int) (string, int) {
	if !ir.Indexed {
		return fmt.Sprint(val), val
	}
	i := (val - ir.Min) / ir.Step
	if i < 0 || i >= len(ir.S) {
		return "UNDEF", -1
	}
	return ir.S[i], i
}

func (dr *Drange) sql(table string) string {
	if table != "" {
		table = "`" + table + "`."
	}
	val := table + "`dval`"
	if !dr.Indexed {
		val = "DATE(" + table + "`dval`)"
	} else {
		switch dr.R {
		case Dr_hour:
			val = "STR_TO_DATE(CONCAT(DATE(" + table + "`dval`), \",\", HOUR(" + table + "`dval`)), \"%Y-%m-%d,%H\")"
		case Dr_day:
			val = "DATE(" + table + "`dval`)"
		case Dr_month:
			val = "STR_TO_DATE(CONCAT(YEAR(" + table + "`dval`), \"-\", MONTH(" + table + "`dval`), \"-01\"), \"%Y-%m-%d\")"
		case Dr_year:
			val = "STR_TO_DATE(CONCAT(YEAR(" + table + "`dval`), \"-01-01\"), \"%Y-%m-%d\")"
		case Dr_dec:
			val = "STR_TO_DATE(CONCAT(10*FLOOR(YEAR(" + table + "`dval`)/10), \"-01-01\"), \"%Y-%m-%d\")"
		case Dr_cent:
			val = "STR_TO_DATE(CONCAT(100*FLOOR(YEAR(" + table + "`dval`)/100), \"-01-01\"), \"%Y-%m-%d\")"
		}
	}
	return val
}

func (fr *Frange) sql(table string) string {
	if table != "" {
		table = "`" + table + "`."
	}
	return fmt.Sprintf("%g * FLOOR(%s`fval`/%g)", fr.Step, table, fr.Step)
}

func (ir *Irange) sql(table string) string {
	if table != "" {
		table = "`" + table + "`."
	}
	return fmt.Sprintf("%d * FLOOR(%s`ival`/%d)", ir.Step, table, ir.Step)
}

func PrintDate(t time.Time, hasTime bool) string {
	if hasTime {
		return fmt.Sprintf("%s %2d %s %d %02d:%02d", dagen[t.Weekday()], t.Day(), Maanden[t.Month()], t.Year(), t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%s %2d %s %d", dagen[t.Weekday()], t.Day(), Maanden[t.Month()], t.Year())
}
