// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command testdata-refresh regenerates test_pb_data/data.db with all pending
// migrations applied, so that tests.NewTestApp does not re-run migrations on
// every bootstrap. Run it after a PocketBase upgrade (new core migrations) or
// after changing pb_migrations.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/tests"
)

func main() {
	target := "test_pb_data"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	app, err := tests.NewTestApp(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create test app:", err)
		os.Exit(1)
	}

	// Apply any pending migrations, then close cleanly so SQLite checkpoints
	// the WAL into data.db before we copy it back into the source directory.
	if err := app.RunAllMigrations(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to run migrations:", err)
		os.Exit(1)
	}
	if err := app.ResetBootstrapState(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to close app cleanly:", err)
		os.Exit(1)
	}

	in, err := os.ReadFile(filepath.Join(app.DataDir(), "data.db"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read migrated db:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(target, "data.db"), in, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write migrated db:", err)
		os.Exit(1)
	}
	for _, suffix := range []string{"-wal", "-shm"} {
		_ = os.Remove(filepath.Join(target, "data.db"+suffix))
	}
	app.Cleanup()

	fmt.Println("refreshed", filepath.Join(target, "data.db"))
}
