// SPDX-License-Identifier: GPL-3.0-only

package store

import "fmt"

// Validate keeps malformed cleanup policies from reaching the database or the
// background worker.
func (s CleanupSettings) Validate() error {
	if s.TTLHours <= 0 {
		return fmt.Errorf("ttl_hours must be greater than 0")
	}
	if s.MaxItems <= 0 {
		return fmt.Errorf("max_items must be greater than 0")
	}
	if s.MaxTotalSizeMB <= 0 {
		return fmt.Errorf("max_total_size_mb must be greater than 0")
	}
	if s.IntervalMinutes <= 0 {
		return fmt.Errorf("interval_minutes must be greater than 0")
	}
	return nil
}
