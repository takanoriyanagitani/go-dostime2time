package dostime2time

import (
	"time"
)

type DosTime16 uint16

type DosDate16 uint16

type Dostime struct {
	DosDate16
	DosTime16
}

func (d DosDate16) ToDostime(t DosTime16) Dostime {
	return Dostime{
		DosDate16: d,
		DosTime16: t,
	}
}

func (d DosDate16) YearRaw() uint8 {
	var u uint16 = uint16(d)
	return uint8(u >> 9)
}

func (d DosDate16) Year() uint16 { return 1980 + uint16(d.YearRaw()) }

func (d DosDate16) MonthRaw() uint8 {
	var u uint16 = uint16(d)
	var m uint16 = u >> 5
	return uint8(m & 0xf)
}

func (d DosDate16) MonthUnchecked() time.Month {
	var u uint8 = d.MonthRaw()
	var i int = int(u)
	return time.Month(i)
}

func (d DosDate16) Day() uint8 {
	var u uint16 = uint16(d)
	return uint8(u & 0x1f)
}

func (t DosTime16) Hour() uint8 {
	var u uint16 = uint16(t)
	return uint8(u >> 11)
}

func (t DosTime16) Minute() uint8 {
	var u uint16 = uint16(t)
	var m uint16 = u >> 5
	return uint8(m & 0x3f)
}

func (t DosTime16) Second() uint8 {
	var u uint16 = uint16(t)
	var s uint16 = u & 0x1f
	return uint8(s << 1)
}

type SimpleLocalTime struct {
	Year  uint16     `json:"year"`
	Month time.Month `json:"month"`
	Day   uint8      `json:"day"`

	Hour   uint8 `json:"hour"`
	Minute uint8 `json:"minute"`
	Second uint8 `json:"second"`
}

func (s SimpleLocalTime) ToLocalTime() time.Time {
	return time.Date(
		int(s.Year),
		s.Month,
		int(s.Day),
		int(s.Hour),
		int(s.Minute),
		int(s.Second),
		0,
		time.Local,
	)
}

func (s SimpleLocalTime) ToUtcTime() time.Time {
	return time.Date(
		int(s.Year),
		s.Month,
		int(s.Day),
		int(s.Hour),
		int(s.Minute),
		int(s.Second),
		0,
		time.UTC,
	)
}

func (s Dostime) ToSimpleLocalTimeUnchecked() SimpleLocalTime {
	var d DosDate16 = s.DosDate16
	var t DosTime16 = s.DosTime16
	return SimpleLocalTime{
		Year:  d.Year(),
		Month: d.MonthUnchecked(),
		Day:   d.Day(),

		Hour:   t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}
}

type JsonNumber float64

type Signed int32

type Unsigned uint32

func (u Unsigned) ToDostime() Dostime {
	var h uint32 = uint32(u) >> 16
	var l uint32 = uint32(u) & 0xffff
	var d DosDate16 = DosDate16(uint16(h))
	var t DosTime16 = DosTime16(uint16(l))
	return d.ToDostime(t)
}

func (s Signed) ToDostime() Dostime {
	var i int32 = int32(s)
	var u Unsigned = Unsigned(uint32(i))
	return u.ToDostime()
}

func (j JsonNumber) ToDostime() Dostime {
	var f float64 = float64(j)
	var negative bool = f < 0
	switch negative {
	case true:
		var s Signed = Signed(int32(f))
		return s.ToDostime()
	default:
		var u Unsigned = Unsigned(uint32(f))
		return u.ToDostime()
	}
}
