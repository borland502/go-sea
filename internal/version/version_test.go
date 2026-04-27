package version

import (
	"runtime/debug"
	"testing"
)

func TestResolveMetadataUsesBuildInfoFallbacks(t *testing.T) {
	buildInfo := &debug.BuildInfo{
		Main: debug.Module{Version: "v0.1.1"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "abc123def456"},
			{Key: "vcs.time", Value: "2026-04-27T16:41:47Z"},
		},
	}

	versionValue := resolveVersion(buildInfo)
	commitValue, dateValue := resolveBuildSettings(buildInfo, "unknown", "unknown")

	if versionValue != "v0.1.1" {
		t.Fatalf("resolveVersion() = %q, want %q", versionValue, "v0.1.1")
	}

	if commitValue != "abc123def456" {
		t.Fatalf("resolveBuildSettings() commit = %q, want %q", commitValue, "abc123def456")
	}

	if dateValue != "2026-04-27T16:41:47Z" {
		t.Fatalf("resolveBuildSettings() date = %q, want %q", dateValue, "2026-04-27T16:41:47Z")
	}
}

func TestResolveVersionLeavesDevelBuildsAsDev(t *testing.T) {
	buildInfo := &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}

	if got := resolveVersion(buildInfo); got != "dev" {
		t.Fatalf("resolveVersion() = %q, want %q", got, "dev")
	}
}

func TestResolveBuildSettingsKeepsExistingValues(t *testing.T) {
	buildInfo := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "new-commit"},
			{Key: "vcs.time", Value: "2026-04-27T00:00:00Z"},
		},
	}

	commitValue, dateValue := resolveBuildSettings(buildInfo, "set-commit", "set-date")

	if commitValue != "set-commit" {
		t.Fatalf("resolveBuildSettings() commit = %q, want %q", commitValue, "set-commit")
	}

	if dateValue != "set-date" {
		t.Fatalf("resolveBuildSettings() date = %q, want %q", dateValue, "set-date")
	}
}
