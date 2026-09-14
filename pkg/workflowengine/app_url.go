// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflowengine

import (
	"os"
	"strings"
)

// AppURLConfigKey holds the public, user-facing base URL of the Credimi
// deployment (PocketBase Settings → App URL). It is persisted in workflow
// inputs, queue tickets, schedule configs, and memos, and is used to build
// links shown to users (test runs, emails, PR comments, screenshots).
const AppURLConfigKey = "app_url"

// InternalAppURLConfigKey optionally holds a deployment-local base URL for
// server-to-server HTTP calls from Temporal workflows and activities back to
// the Credimi API. When the app URL is proxied by a Cloudflare WAF that
// challenges non-browser traffic, set CREDIMI_INTERNAL_APP_URL so callbacks
// bypass the proxy. Handlers inject it into workflow config at start time.
const InternalAppURLConfigKey = "internal_app_url"

// InternalAppURLConfigKeyEnv is the environment variable handlers read to
// populate InternalAppURLConfigKey in workflow input configs.
const InternalAppURLConfigKeyEnv = "CREDIMI_INTERNAL_APP_URL"

// InternalAppURLOverride returns the deployment-local base URL from the
// environment, trimmed and empty when unset.
func InternalAppURLOverride() string {
	return strings.TrimSpace(os.Getenv(InternalAppURLConfigKeyEnv))
}

// InternalAppURLFromConfig returns the base URL that workflows and activities
// must use for HTTP callbacks to the Credimi API. It prefers
// InternalAppURLConfigKey and falls back to AppURLConfigKey, so behavior is
// unchanged in deployments that do not configure an internal override.
func InternalAppURLFromConfig(config map[string]any) string {
	if config != nil {
		if v, ok := config[InternalAppURLConfigKey].(string); ok {
			if trimmed := strings.TrimSpace(v); trimmed != "" {
				return trimmed
			}
		}
	}
	v, _ := config[AppURLConfigKey].(string)
	return strings.TrimSpace(v)
}

// WithInternalAppURL returns config unchanged when no deployment-local
// override is configured, and otherwise injects InternalAppURLConfigKey describing
// it. Handlers that build workflow input config use this so worker-side HTTP
// callbacks can bypass a Cloudflare WAF while user-facing links keep the
// public app_url.
func WithInternalAppURL(config map[string]any) map[string]any {
	if internalURL := InternalAppURLOverride(); internalURL != "" {
		config[InternalAppURLConfigKey] = internalURL
	}
	return config
}
