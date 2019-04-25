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
// - any supported date format;
// - 'y' or 'yesterday', case-insensitively;
// - any value that can be parse into an integer, as a relative date.
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
