package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	dt "github.com/takanoriyanagitani/go-dostime2time"
	. "github.com/takanoriyanagitani/go-dostime2time/util"
)

type TimeWriter func(time.Time) IO[Void]

func (t TimeWriter) ToSimpleLocalTimeWriterLocal() SimpleLocalTimeWriter {
	return func(s dt.SimpleLocalTime) IO[Void] {
		var ltime time.Time = s.ToLocalTime()
		return t(ltime)
	}
}

func (t TimeWriter) ToSimpleLocalTimeWriterUTC() SimpleLocalTimeWriter {
	return func(s dt.SimpleLocalTime) IO[Void] {
		var ltime time.Time = s.ToUtcTime()
		return t(ltime)
	}
}

type SimpleLocalTimeWriter func(dt.SimpleLocalTime) IO[Void]

type DostimeWriter func(dt.Dostime) IO[Void]

func (s SimpleLocalTimeWriter) ToDostimeWriter() DostimeWriter {
	return func(d dt.Dostime) IO[Void] {
		var lt dt.SimpleLocalTime = d.ToSimpleLocalTimeUnchecked()
		return s(lt)
	}
}

type StringWriter func(string) IO[Void]

type Writer struct{ io.Writer }

func (w Writer) ToStringWriter() StringWriter {
	return func(s string) IO[Void] {
		return func(_ context.Context) (Void, error) {
			_, e := io.WriteString(w.Writer, s+"\n")
			return Empty, e
		}
	}
}

func (s StringWriter) ToTimeWriter(layout string) TimeWriter {
	return func(t time.Time) IO[Void] {
		var stime string = t.Format(layout)
		return s(stime)
	}
}

var string2stdout StringWriter = Writer{os.Stdout}.ToStringWriter()

var time2stdoutRfc3339 TimeWriter = string2stdout.ToTimeWriter(time.RFC3339)

var ltime2stdoutLocal SimpleLocalTimeWriter = time2stdoutRfc3339.
	ToSimpleLocalTimeWriterLocal()

var ltime2stdout SimpleLocalTimeWriter = ltime2stdoutLocal

var dtime2stdout DostimeWriter = ltime2stdout.ToDostimeWriter()

var args IO[[]string] = Of(os.Args[1:])

var rawDostimeFromArg IO[string] = Bind(
	args,
	Lift(func(s []string) (string, error) {
		var sz int = len(s)
		switch 0 < sz {
		case true:
			return s[0], nil
		default:
			return "", fmt.Errorf("unexpected arg count: %v", sz)
		}
	}),
)

var dostimeIntFromArg IO[int] = Bind(
	rawDostimeFromArg,
	Lift(strconv.Atoi),
)

var signedFromArg IO[dt.Signed] = Bind(
	dostimeIntFromArg,
	Lift(func(i int) (dt.Signed, error) { return dt.Signed(int32(i)), nil }),
)

var dostimeFromArg IO[dt.Dostime] = Bind(
	signedFromArg,
	Lift(func(s dt.Signed) (dt.Dostime, error) { return s.ToDostime(), nil }),
)

var arg2dostime2ltime2stdout IO[Void] = Bind(
	dostimeFromArg,
	dtime2stdout,
)

func main() {
	_, e := arg2dostime2ltime2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
