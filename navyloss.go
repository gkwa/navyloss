package navyloss

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	LogFormat string `long:"log-format" choice:"text" choice:"json" default:"text" description:"Log format"`
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose debug information, each -v bumps log level"`
	Age       string `short:"a" long:"age" description:"Age parameter in the format 1y, 10M, 10m, 200s, 34d, 1y23d, 2d20s, etc." required:"true"`
	logLevel  slog.Level
}

func Execute() int {
	if err := parseFlags(); err != nil {
		return 1
	}

	if err := setLogLevel(); err != nil {
		return 1
	}

	if err := setupLogger(); err != nil {
		return 1
	}

	if err := run(); err != nil {
		slog.Error("run failed", "error", err)
		return 1
	}

	return 0
}

func parseFlags() error {
	_, err := flags.Parse(&opts)
	return err
}

func run() error {
	seconds, err := AgeToSeconds(opts.Age)
	if err != nil {
		fmt.Printf("Error parsing age parameter: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s = %f seconds\n", opts.Age, seconds)

	t := showDateGivenSecondsAgo(time.Now(), seconds)
	fmt.Printf("Date %d seconds ago: %s\n", int64(seconds), t.Format(time.RFC3339))

	return nil
}

func AgeToSeconds(age string) (float64, error) {
	re := regexp.MustCompile(`(\d+(\.\d+)?)([yMmdsw])`)
	matches := re.FindAllStringSubmatch(age, -1)

	var totalSeconds float64

	for _, match := range matches {
		value, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, err
		}

		unit := match[3]
		switch unit {
		case "y":
			totalSeconds += value * (60 * time.Minute * 24 * 365).Seconds() // seconds in a year
		case "M":
			totalSeconds += value * (60 * time.Minute * 24 * 7 * 30).Seconds() // seconds in a month
		case "w":
			totalSeconds += value * (60 * time.Minute * 24 * 7).Seconds() // seconds in a week
		case "d":
			totalSeconds += value * (60 * time.Minute * 24).Seconds() // seconds in a day
		case "m":
			totalSeconds += value * (1 * time.Minute * 60).Seconds() // seconds in a minute
		case "s":
			totalSeconds += value
		}
	}

	return totalSeconds, nil
}

func showDateGivenSecondsAgo(currentTime time.Time, seconds float64) time.Time {
	return currentTime.Add(-time.Second * time.Duration(seconds)).Truncate(time.Second)
}
