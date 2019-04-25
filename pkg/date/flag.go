package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const longDateFormat = "2006-1-2"

type flag struct {
	*time.Time
}

// Flag will create a Flag for use as a pflag.Value
func Flag() *flag {
	return &flag{}
}

// String will return the date in yyyy-mm-dd format, or an empty string if one
// has not been set.
func (f flag) String() string {
	if f.Time == nil {
		return ""
	}
	return f.Time.Format(longDateFormat)
}

// Type returns the string that represents the type of flag.
func (flag) Type() string {
	return "date"
}

// Set parses the given string, attempting to create a logical date from its
// content. Set will match:
// - any value that can be parse into an integer, as a relative date.
// - 'y' or 'yesterday', case-insensitively;
// - any supported date format;
// - m{1,2}d{1,2} and use the current year for the year value
func (f *flag) Set(value string) error {
	val := strings.TrimSpace(value)
	if val == "" {
		return errors.New("no value given")
	}
	val = strings.ToLower(value)

	for _, parse := range []func(string) (time.Time, error){
		parseYesterday,
		parseRelative,
		format(longDateFormat).parse,
		monthDateParser{getYear: func() int { return time.Now().Year() }}.parse,
	} {
		d, err := parse(val)
		if err == nil {
			*f = flag{Time: &d}
			return nil
		}
	}

	// TODO: use multi error here? We don't want to only provide last error
	return fmt.Errorf("unsupported date value: %+v", value)
}

func parseYesterday(val string) (time.Time, error) {
	for _, valid := range []string{"yesterday", "y"} {
		if val == valid {
			return time.Now().Add(-time.Hour * 24), nil
		}
	}
	return time.Time{}, errors.New("unsupported value")
}

func parseRelative(val string) (time.Time, error) {
	i, err := strconv.Atoi(val)
	d := time.Now().Add(time.Hour * (24 * time.Duration(i)))
	return d, errors.Wrap(err, "converting ascii to integer")
}

type format string

func (fmt format) parse(val string) (time.Time, error) {
	d, err := time.Parse(string(fmt), val)
	return d, errors.Wrapf(err, "parsing into format:%s", fmt)
}

type monthDateParser struct {
	getYear func() int
}

func (mdp monthDateParser) parse(val string) (time.Time, error) {
	const shortDateFormat = "1-2"
	d, err := time.Parse(string(shortDateFormat), val)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "parsing into format:%s", shortDateFormat)
	}

	return time.Date(mdp.getYear(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC), nil

}
