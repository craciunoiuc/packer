// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build amd64 && (darwin || windows || linux)

package command

import (
	"log"
	"os"
	"runtime"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/cli"
	"golang.org/x/mod/sumdb/dirhash"
)

type testCasePluginsInstall struct {
	name                                     string
	Meta                                     Meta
	inPluginFolder                           map[string]string
	expectedPackerConfigDirHashBeforeInstall string
	packerConfigDir                          string
	pluginSourceArgs                         []string
	want                                     int
	dirFiles                                 []string
	expectedPackerConfigDirHashAfterInstall  string
}

func TestPluginsInstallCommand_Run(t *testing.T) {

	cfg := &configDirSingleton{map[string]string{}}

	tests := []testCasePluginsInstall{
		{
			name: "already-installed-no-op",
			Meta: TestMetaFile(t),
			inPluginFolder: map[string]string{
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_darwin_amd64":                "1",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_darwin_amd64_SHA256SUM":      "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_windows_amd64.exe":           "1.exe",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_windows_amd64.exe_SHA256SUM": "07d8453027192ee0c4120242e6e84e2ca2328b8e0f506e2f818a1a5b82790a0b",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_linux_amd64":                 "1.out",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_linux_amd64_SHA256SUM":       "59031c50e0dfeedfde2b4e9445754804dce3f29e4efa737eead0ca9b4f5b85a5",
			},
			expectedPackerConfigDirHashBeforeInstall: "h1:Q5qyAOdD43hL3CquQdVfaHpOYGf0UsZ/+wVA9Ry6cbA=",
			packerConfigDir:                          cfg.dir("1_pkr_plugins_config"),
			pluginSourceArgs:                         []string{"github.com/sylviamoss/comment", "v0.2.18"},
			want:                                     0,
			dirFiles:                                 nil,
			expectedPackerConfigDirHashAfterInstall:  "h1:Q5qyAOdD43hL3CquQdVfaHpOYGf0UsZ/+wVA9Ry6cbA=",
		},
		{
			name: "install-newer-version",
			Meta: TestMetaFile(t),
			inPluginFolder: map[string]string{
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_darwin_amd64":                "1",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_darwin_amd64_SHA256SUM":      "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_windows_amd64.exe":           "1.exe",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_windows_amd64.exe_SHA256SUM": "07d8453027192ee0c4120242e6e84e2ca2328b8e0f506e2f818a1a5b82790a0b",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_linux_amd64":                 "1.out",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_linux_amd64_SHA256SUM":       "59031c50e0dfeedfde2b4e9445754804dce3f29e4efa737eead0ca9b4f5b85a5",
			},
			expectedPackerConfigDirHashBeforeInstall: "h1:Q5qyAOdD43hL3CquQdVfaHpOYGf0UsZ/+wVA9Ry6cbA=",
			packerConfigDir:                          cfg.dir("2_pkr_plugins_config"),
			pluginSourceArgs:                         []string{"github.com/sylviamoss/comment", "v0.2.19"},
			want:                                     0,
			dirFiles: []string{
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_darwin_amd64",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_darwin_amd64_SHA256SUM",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_linux_amd64",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_linux_amd64_SHA256SUM",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_windows_amd64.exe",
				"github.com/sylviamoss/comment/packer-plugin-comment_v0.2.18_x5.0_windows_amd64.exe_SHA256SUM",
				map[string]string{
					"darwin":  "github.com/sylviamoss/comment/packer-plugin-comment_v0.2.19_x5.0_darwin_amd64_SHA256SUM",
					"linux":   "github.com/sylviamoss/comment/packer-plugin-comment_v0.2.19_x5.0_linux_amd64_SHA256SUM",
					"windows": "github.com/sylviamoss/comment/packer-plugin-comment_v0.2.19_x5.0_windows_amd64.exe_SHA256SUM",
				}[runtime.GOOS],
				map[string]string{
					"darwin":  "github.com/sylviamoss/comment/packer-plugin-comment_v0.2.19_x5.0_darwin_amd64",
					"linux":   "github.com/sylviamoss/comment/packer-plugin-comment_v0.2.19_x5.0_linux_amd64",
					"windows": "github.com/sylviamoss/comment/packer-plugin-comment_v0.2.19_x5.0_windows_amd64.exe",
				}[runtime.GOOS],
			},
			expectedPackerConfigDirHashAfterInstall: map[string]string{
				"darwin":  "h1:ORwcCYUx8z/5n/QvuTJo2vrgKpfJA4AxlNg1G9/BCDI=",
				"linux":   "h1:CGym0+Nd0LEANgzqL0wx/LDjRL8bYwlpZ0HajPJo/hs=",
				"windows": "h1:ag0/C1YjP7KoEDYOiJHE0K+lhFgs0tVgjriWCXVT1fg=",
			}[runtime.GOOS],
		},
		{
			name:                                     "unsupported-non-github-source-address",
			Meta:                                     TestMetaFile(t),
			inPluginFolder:                           nil,
			expectedPackerConfigDirHashBeforeInstall: "h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
			packerConfigDir:                          cfg.dir("3_pkr_plugins_config"),
			pluginSourceArgs:                         []string{"example.com/sylviamoss/comment", "v0.2.19"},
			want:                                     1,
			dirFiles:                                 nil,
			expectedPackerConfigDirHashAfterInstall:  "h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
		{
			name:                                     "multiple-source-addresses-provided",
			Meta:                                     TestMetaFile(t),
			inPluginFolder:                           nil,
			expectedPackerConfigDirHashBeforeInstall: "h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
			packerConfigDir:                          cfg.dir("4_pkr_plugins_config"),
			pluginSourceArgs:                         []string{"github.com/sylviamoss/comment", "v0.2.18", "github.com/sylviamoss/comment", "v0.2.19"},
			want:                                     cli.RunResultHelp,
			dirFiles:                                 nil,
			expectedPackerConfigDirHashAfterInstall:  "h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
		{
			name:                                     "no-source-address-provided",
			Meta:                                     TestMetaFile(t),
			inPluginFolder:                           nil,
			expectedPackerConfigDirHashBeforeInstall: "h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
			packerConfigDir:                          cfg.dir("5_pkr_plugins_config"),
			pluginSourceArgs:                         []string{},
			want:                                     cli.RunResultHelp,
			dirFiles:                                 nil,
			expectedPackerConfigDirHashAfterInstall:  "h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("starting %s", tt.name)
			log.Printf("%#v", tt)
			t.Cleanup(func() {
				_ = os.RemoveAll(tt.packerConfigDir)
			})
			t.Setenv("PACKER_CONFIG_DIR", tt.packerConfigDir)
			createFiles(tt.packerConfigDir, tt.inPluginFolder)

			hash, err := dirhash.HashDir(tt.packerConfigDir, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("HashDir: %v", err)
			}
			if diff := cmp.Diff(tt.expectedPackerConfigDirHashBeforeInstall, hash); diff != "" {
				t.Errorf("unexpected dir hash before plugins install: +found -expected %s", diff)
			}

			c := &PluginsInstallCommand{
				Meta: tt.Meta,
			}

			if err := c.CoreConfig.Components.PluginConfig.Discover(); err != nil {
				t.Fatalf("Failed to discover plugins: %s", err)
			}

			c.CoreConfig.Components.PluginConfig.KnownPluginFolders = []string{tt.packerConfigDir}
			if got := c.Run(tt.pluginSourceArgs); got != tt.want {
				t.Errorf("PluginsInstallCommand.Run() = %v, want %v", got, tt.want)
			}

			if tt.dirFiles != nil {
				dirFiles, err := dirhash.DirFiles(tt.packerConfigDir, "")
				if err != nil {
					t.Fatalf("DirFiles: %v", err)
				}
				sort.Strings(tt.dirFiles)
				sort.Strings(dirFiles)
				if diff := cmp.Diff(tt.dirFiles, dirFiles); diff != "" {
					t.Errorf("found files differ: %v", diff)
				}
			}

			hash, err = dirhash.HashDir(tt.packerConfigDir, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("HashDir: %v", err)
			}
			if diff := cmp.Diff(tt.expectedPackerConfigDirHashAfterInstall, hash); diff != "" {
				t.Errorf("unexpected dir hash after plugins install: %s", diff)
			}
		})
	}
}
