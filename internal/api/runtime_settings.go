// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"context"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

func (r *Router) defaultLimits() store.LimitsSettings {
	return store.LimitsSettings{
		MinTextBytes:    r.config.Limits.MinTextBytes,
		MaxTextBytes:    r.config.Limits.MaxTextBytes,
		MinImageBytes:   r.config.Limits.MinImageBytes,
		MaxImageBytes:   r.config.Limits.MaxImageBytes,
		MinFileBytes:    r.config.Limits.MinFileBytes,
		MaxFileBytes:    r.config.Limits.MaxFileBytes,
		MinLinkBytes:    r.config.Limits.MinLinkBytes,
		MaxLinkBytes:    r.config.Limits.MaxLinkBytes,
		MaxRequestBytes: r.config.Limits.MaxRequestBytes,
	}
}

func (r *Router) currentLimits(ctx context.Context) store.LimitsSettings {
	limits, err := r.store.LoadLimitsSettings(ctx, r.defaultLimits())
	if err != nil {
		return r.defaultLimits()
	}
	return limits
}

func (r *Router) currentAdminToken(ctx context.Context) string {
	token, err := r.store.LoadAdminToken(ctx, r.config.Auth.Token)
	if err != nil {
		return r.config.Auth.Token
	}
	return token
}
