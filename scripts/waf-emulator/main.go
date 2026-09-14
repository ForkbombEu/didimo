// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command waf-emulator emulates a Cloudflare-proxied Credimi deployment.
//
// It reproduces the failure mode that breaks Temporal callbacks: browser-like
// requests are proxied to the origin, while non-browser clients (Go's
// net/http, curl, server SDKs) receive a managed-challenge interstitial
// instead of reaching the API.
//
// Use it to prove that worker-to-API traffic must not traverse the public,
// WAF-fronted URL. See docs/src/content/docs/software-architecture/developer-setup.md.
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	defaultListen   = ":8091"
	defaultUpstream = "http://localhost:8090"

	listenEnv   = "WAF_EMULATOR_LISTEN"
	upstreamEnv = "WAF_EMULATOR_UPSTREAM"
)

// nonBrowserUserAgent matches clients Cloudflare would classify as automated.
var nonBrowserUserAgent = regexp.MustCompile(
	`(?i)(curl|wget|go-http-client|python-requests|python-urllib|node-fetch|axios|okhttp|java/|libwww|httpclient)`,
)

const challengePage = `<!DOCTYPE html>
<html>
<head><title>Just a moment...</title></head>
<body>
<h1>Just a moment...</h1>
<p>Verifying you are human. This may take a few seconds.</p>
<p>Enable JavaScript and cookies to continue.</p>
<p>Attention Required! | Cloudflare</p>
</body>
</html>
`

func main() {
	listen := envOrDefault(listenEnv, defaultListen)
	upstreamRaw := envOrDefault(upstreamEnv, defaultUpstream)

	upstream, err := url.Parse(upstreamRaw)
	if err != nil {
		log.Fatalf("invalid %s %q: %v", upstreamEnv, upstreamRaw, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(upstream)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isBrowserLike(r) {
			log.Printf("proxied -> %s %s", r.Method, r.URL.Path)
			proxy.ServeHTTP(w, r)
			return
		}

		log.Printf("challenged <- %s %s ua=%q", r.Method, r.URL.Path, r.UserAgent())
		w.Header().Set("Server", "cloudflare")
		w.Header().Set("cf-mitigated", "challenge")
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(challengePage))
	})

	log.Printf(
		"waf-emulator listening on %s, proxying browser-like traffic to %s",
		listen,
		upstream,
	)
	log.Fatal(http.ListenAndServe(listen, handler))
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// isBrowserLike reports whether the request looks like a document navigation.
// Cloudflare's managed challenge passes those and challenges everyday HTTP
// clients, which is exactly the split that breaks server-to-server callbacks.
func isBrowserLike(r *http.Request) bool {
	if !strings.Contains(r.Header.Get("Accept"), "text/html") {
		return false
	}

	agent := strings.TrimSpace(r.UserAgent())
	if agent == "" {
		return false
	}

	return !nonBrowserUserAgent.MatchString(agent)
}
