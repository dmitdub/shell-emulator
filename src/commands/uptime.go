package commands

import (
	"fmt"
	"time"
)

// Time constants used to format durations.
const (
	secondsPerMinute = 60
	secondsPerHour   = 3600
)

// Uptime implements the uptime command.
//
// It reports how long the emulator has been running since it
// was started. The output has the form "up 2 minutes, 15 seconds".
func (c *Context) Uptime(args []string) {
	elapsed := time.Since(c.Started)
	fmt.Fprintf(c.Out, "up %s\n", formatDuration(elapsed))
}

// formatDuration renders a duration in a human-readable form.
//
// Only hours, minutes, and seconds are included; zero-value
// components are omitted.
func formatDuration(d time.Duration) string {
	total := int(d.Seconds())

	hours := total / secondsPerHour
	minutes := (total % secondsPerHour) / secondsPerMinute
	seconds := total % secondsPerMinute

	parts := make([]string, 0, 3)
	if hours > 0 {
		parts = append(parts, plural(hours, "hour"))
	}
	if minutes > 0 {
		parts = append(parts, plural(minutes, "minute"))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, plural(seconds, "second"))
	}

	return joinParts(parts)
}

// plural formats a number with the given singular noun,
// appending an "s" for values other than one.
func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

// joinParts joins parts with commas and the word "and" before
// the last component.
func joinParts(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		return parts[0] + ", " + joinParts(parts[1:])
	}
}
