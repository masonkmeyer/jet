package jet

import (
	"fmt"
	"time"
)

// Ago returns a human-readable string representing the time elapsed since the pastTime
func Ago(pastTime time.Time) string {
	elapsed := time.Since(pastTime)

	// Calculate elapsed time in different units
	seconds := int(elapsed.Seconds())
	minutes := int(elapsed.Minutes())
	hours := int(elapsed.Hours())
	days := int(elapsed.Hours() / 24)
	weeks := int(elapsed.Hours() / 168)
	months := int(elapsed.Hours() / 720)

	// Format the output
	if months > 0 {
		return fmt.Sprintf("%d months ago", months)
	} else if weeks > 0 {
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if days > 1 {
		return fmt.Sprintf("%d days ago", days)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours ago", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if seconds > 0 {
		return fmt.Sprintf("%d seconds ago", seconds)
	}

	return ""
}
