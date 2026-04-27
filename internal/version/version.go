// Package version holds build-time injected version metadata.
// These variables are set by the linker (typically via GoReleaser) at build time.
// When ldflags are absent, the package falls back to Go build metadata.
//
// Usage in ldflags:
//
//	-ldflags "-X github.com/borland502/go-sea/internal/version.Version=v1.2.3 ..."
package version

import "runtime/debug"

// Version is the semantic version of the build (e.g., "v1.2.3", or "dev" for debug builds).
var Version = "dev"

// Commit is the git commit hash included in this build.
var Commit = "unknown"

// Date is the build timestamp in RFC3339 format (e.g., "2026-04-24T22:15:30Z").
var Date = "unknown"

func init() {
	Version, Commit, Date = resolveMetadata(Version, Commit, Date)
}

func resolveMetadata(versionValue, commitValue, dateValue string) (string, string, string) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return versionValue, commitValue, dateValue
	}

	if versionValue == "dev" {
		versionValue = resolveVersion(buildInfo)
	}

	if commitValue == "unknown" || dateValue == "unknown" {
		commitValue, dateValue = resolveBuildSettings(buildInfo, commitValue, dateValue)
	}

	return versionValue, commitValue, dateValue
}

func resolveVersion(buildInfo *debug.BuildInfo) string {
	if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}

	return "dev"
}

func resolveBuildSettings(buildInfo *debug.BuildInfo, commitValue, dateValue string) (string, string) {
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			if commitValue == "unknown" && setting.Value != "" {
				commitValue = setting.Value
			}
		case "vcs.time":
			if dateValue == "unknown" && setting.Value != "" {
				dateValue = setting.Value
			}
		}
	}

	return commitValue, dateValue
}
