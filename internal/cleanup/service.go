// SPDX-License-Identifier: GPL-3.0-only

package cleanup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
)

// Service owns the retention policy, background worker loop, and the latest
// cleanup execution summary.
type Service struct {
	store    *store.SQLiteStore
	defaults store.CleanupSettings

	stopOnce sync.Once
	stopCh   chan struct{}
	doneCh   chan struct{}
}

// NewService creates the retention service and seeds the persistent cleanup
// policy from config defaults when the database has not stored overrides yet.
func NewService(dbStore *store.SQLiteStore, cfg *config.Config) (*Service, error) {
	service := &Service{
		store: dbStore,
		defaults: store.CleanupSettings{
			TTLHours:        cfg.Storage.TTLHours,
			MaxItems:        cfg.Storage.MaxItems,
			MaxTotalSizeMB:  cfg.Storage.MaxTotalSizeMB,
			IntervalMinutes: cfg.Cleaner.IntervalMinutes,
			Enabled:         cfg.Cleaner.Enabled,
		},
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}

	if err := service.defaults.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	settings, err := service.store.LoadCleanupSettings(ctx, service.defaults)
	if err != nil {
		return nil, err
	}

	if err := service.store.RefreshExpirations(ctx, settings.TTLHours); err != nil {
		return nil, err
	}

	go service.loop()
	return service, nil
}

// Close stops the background worker.
func (s *Service) Close() {
	if s == nil {
		return
	}

	s.stopOnce.Do(func() {
		close(s.stopCh)
		<-s.doneCh
	})
}

// GetSettings returns the current persisted cleanup policy.
func (s *Service) GetSettings(ctx context.Context) (store.CleanupSettings, error) {
	return s.store.LoadCleanupSettings(ctx, s.defaults)
}

// UpdateSettings persists a new cleanup policy and refreshes item expirations
// immediately so the worker and manual cleanup behave consistently.
func (s *Service) UpdateSettings(ctx context.Context, settings store.CleanupSettings) (store.CleanupSettings, error) {
	if err := settings.Validate(); err != nil {
		return store.CleanupSettings{}, err
	}

	if err := s.store.SaveCleanupSettings(ctx, settings); err != nil {
		return store.CleanupSettings{}, err
	}

	if err := s.store.RefreshExpirations(ctx, settings.TTLHours); err != nil {
		return store.CleanupSettings{}, err
	}

	return settings, nil
}

// GetStatus returns the latest cleanup run summary.
func (s *Service) GetStatus(ctx context.Context) (store.CleanupStatus, error) {
	status, err := s.store.LoadCleanupStatus(ctx)
	if err != nil {
		return store.CleanupStatus{}, err
	}

	storageStatus, err := s.store.GetStorageStatus(ctx)
	if err != nil {
		return store.CleanupStatus{}, err
	}

	status.HistoryCount = storageStatus.HistoryCount
	status.FavoriteCount = storageStatus.FavoriteCount
	status.TotalBytes = storageStatus.TotalBytes
	status.FileBytes = storageStatus.FileBytes
	return status, nil
}

// GetStorageStatus returns a point-in-time storage summary.
func (s *Service) GetStorageStatus(ctx context.Context) (store.StorageStatus, error) {
	return s.store.GetStorageStatus(ctx)
}

// RunNow executes one cleanup pass on demand.
func (s *Service) RunNow(ctx context.Context, reason string) (store.CleanupStatus, error) {
	return s.runCleanup(ctx, reason)
}

func (s *Service) loop() {
	defer close(s.doneCh)

	for {
		settings, err := s.GetSettings(context.Background())
		if err != nil {
			s.recordFailure("background-loop", err)
			select {
			case <-time.After(30 * time.Second):
				continue
			case <-s.stopCh:
				return
			}
		}

		waitDuration := time.Duration(settings.IntervalMinutes) * time.Minute
		if waitDuration <= 0 {
			waitDuration = 30 * time.Minute
		}

		select {
		case <-time.After(waitDuration):
			if !settings.Enabled {
				continue
			}

			runCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_, err := s.runCleanup(runCtx, "scheduled")
			cancel()
			if err != nil {
				s.recordFailure("scheduled", err)
			}
		case <-s.stopCh:
			return
		}
	}
}

func (s *Service) runCleanup(ctx context.Context, reason string) (store.CleanupStatus, error) {
	settings, err := s.GetSettings(ctx)
	if err != nil {
		return store.CleanupStatus{}, err
	}

	if err := s.store.RefreshExpirations(ctx, settings.TTLHours); err != nil {
		return store.CleanupStatus{}, err
	}

	candidates, err := s.store.ListCleanupCandidates(ctx)
	if err != nil {
		return store.CleanupStatus{}, err
	}

	now := time.Now().UTC()
	deletedAt := now.Format(time.RFC3339)
	deletedIDs := make([]int64, 0)
	status := store.CleanupStatus{
		LastRunAt:     deletedAt,
		LastRunReason: reason,
	}

	activeCandidates := make([]store.CleanupCandidate, 0, len(candidates))

	for _, candidate := range candidates {
		if candidate.IsFavorite {
			activeCandidates = append(activeCandidates, candidate)
			continue
		}

		if candidate.ExpiresAt != "" {
			expiresAt, parseErr := time.Parse(time.RFC3339, candidate.ExpiresAt)
			if parseErr == nil && !now.Before(expiresAt) {
				if err := s.deleteCandidate(ctx, candidate, deletedAt); err != nil {
					return store.CleanupStatus{}, err
				}
				deletedIDs = append(deletedIDs, candidate.ID)
				status.DeletedExpired++
				if candidate.ItemType == "file" {
					status.DeletedFiles++
				}
				continue
			}
		}

		activeCandidates = append(activeCandidates, candidate)
	}

	nonFavorite := filterNonFavorite(activeCandidates)
	for len(activeCandidates) > settings.MaxItems && len(nonFavorite) > 0 {
		candidate := nonFavorite[0]
		nonFavorite = nonFavorite[1:]

		if err := s.deleteCandidate(ctx, candidate, deletedAt); err != nil {
			return store.CleanupStatus{}, err
		}
		deletedIDs = append(deletedIDs, candidate.ID)
		status.DeletedMaxItems++
		if candidate.ItemType == "file" {
			status.DeletedFiles++
		}
		activeCandidates = removeCandidateByID(activeCandidates, candidate.ID)
	}

	limitBytes := int64(settings.MaxTotalSizeMB) * 1024 * 1024
	totalBytes := sumCandidateBytes(activeCandidates)
	if totalBytes > limitBytes {
		fileCandidates := filterItemType(filterNonFavorite(activeCandidates), "file")
		otherCandidates := filterNotItemType(filterNonFavorite(activeCandidates), "file")
		deleteQueue := append(fileCandidates, otherCandidates...)

		for _, candidate := range deleteQueue {
			if totalBytes <= limitBytes {
				break
			}

			if err := s.deleteCandidate(ctx, candidate, deletedAt); err != nil {
				return store.CleanupStatus{}, err
			}
			deletedIDs = append(deletedIDs, candidate.ID)
			status.DeletedStorage++
			if candidate.ItemType == "file" {
				status.DeletedFiles++
			}
			totalBytes -= candidate.SizeBytes
			activeCandidates = removeCandidateByID(activeCandidates, candidate.ID)
		}
	}

	storageStatus, err := s.store.GetStorageStatus(ctx)
	if err != nil {
		return store.CleanupStatus{}, err
	}

	status.HistoryCount = storageStatus.HistoryCount
	status.FavoriteCount = storageStatus.FavoriteCount
	status.TotalBytes = storageStatus.TotalBytes
	status.FileBytes = storageStatus.FileBytes
	status.NonFavoriteFileSize = sumNonFavoriteFileBytes(activeCandidates)

	if err := s.store.SaveCleanupStatus(ctx, status); err != nil {
		return store.CleanupStatus{}, err
	}

	_ = deletedIDs
	return status, nil
}

func (s *Service) deleteCandidate(ctx context.Context, candidate store.CleanupCandidate, deletedAt string) error {
	if candidate.ItemType == "file" && candidate.Metadata.LocalPath != "" {
		if err := os.Remove(candidate.Metadata.LocalPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove local file %q: %w", candidate.Metadata.LocalPath, err)
		}
	}

	return s.store.MarkItemDeletedForCleanup(ctx, candidate.ID, deletedAt)
}

func (s *Service) recordFailure(reason string, runErr error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, _ := s.store.LoadCleanupStatus(ctx)
	status.LastRunAt = time.Now().UTC().Format(time.RFC3339)
	status.LastRunReason = reason
	status.LastError = runErr.Error()
	_ = s.store.SaveCleanupStatus(ctx, status)
}

func filterNonFavorite(candidates []store.CleanupCandidate) []store.CleanupCandidate {
	filtered := make([]store.CleanupCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if !candidate.IsFavorite {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func filterItemType(candidates []store.CleanupCandidate, itemType string) []store.CleanupCandidate {
	filtered := make([]store.CleanupCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.ItemType == itemType {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func filterNotItemType(candidates []store.CleanupCandidate, itemType string) []store.CleanupCandidate {
	filtered := make([]store.CleanupCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.ItemType != itemType {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func removeCandidateByID(candidates []store.CleanupCandidate, id int64) []store.CleanupCandidate {
	filtered := make([]store.CleanupCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.ID != id {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func sumCandidateBytes(candidates []store.CleanupCandidate) int64 {
	var total int64
	for _, candidate := range candidates {
		total += candidate.SizeBytes
	}
	return total
}

func sumNonFavoriteFileBytes(candidates []store.CleanupCandidate) int64 {
	var total int64
	for _, candidate := range candidates {
		if !candidate.IsFavorite && candidate.ItemType == "file" {
			total += candidate.SizeBytes
		}
	}
	return total
}
