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
	Period    string `short:"p" long:"period" description:"Time period parameter in the format 1y, 10M, 10m, 200s, 34d, 1y23d, 2d20s, etc." required:"true"`
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
	duration, err := DurationFromString(opts.Period)
	if err != nil {
		fmt.Printf("Error parsing period parameter: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s = %f seconds, %v\n", opts.Period, duration.Seconds(), duration)

	t := showDateGivenSecondsAgo(time.Now(), duration.Seconds())
	fmt.Printf("Date %d seconds ago: %s\n", int64(duration), t.Format(time.RFC3339))

	return nil
}

func DurationFromString(period string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+(\.\d+)?)([yMwdhms])`)
	matches := re.FindAllStringSubmatch(period, -1)

	var totalDuration time.Duration

	for _, match := range matches {
		value, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, err
		}

		unit := match[3]
		switch unit {
		case "y":
			totalDuration += time.Duration(value) * 365 * 24 * time.Hour
		case "M":
			totalDuration += time.Duration(value) * 30 * 24 * time.Hour
		case "w":
			totalDuration += time.Duration(value) * 7 * 24 * time.Hour
		case "d":
			totalDuration += time.Duration(value) * 24 * time.Hour
		case "h":
			totalDuration += time.Duration(value) * time.Hour
		case "m":
			totalDuration += time.Duration(value) * time.Minute
		case "s":
			totalDuration += time.Duration(value) * time.Second
		default:
			return 0, fmt.Errorf("unknown duration unit: %s", unit)
		}

		slog.Debug("parsed duration",
			"value", value,
			"unit", unit,
			"totalDuration", totalDuration,
		)
	}

	slog.Debug("parsed duration",
		"totalDuration", totalDuration,
	)

	return totalDuration, nil
}

func showDateGivenSecondsAgo(currentTime time.Time, seconds float64) time.Time {
	return currentTime.Add(-time.Second * time.Duration(seconds)).Truncate(time.Second)
}
