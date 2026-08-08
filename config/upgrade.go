/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bootswithdefer/awless/database"
)

const (
	lastUpgradeCheckDBKey = "upgrade.lastcheck"
)

func VerifyNewVersionAvailable(ctx context.Context, url string, messaging io.Writer) error {
	return database.Execute(func(db *database.DB) error {
		last, err := db.GetTimeValue(lastUpgradeCheckDBKey)
		if err != nil {
			return err
		}

		upgradeFreq := getCheckUpgradeFrequency()
		if upgradeFreq < 0 {
			return nil
		}

		if time.Since(last) > upgradeFreq {
			// Advisory only: failing to reach the release endpoint must never
			// affect the command the user actually ran.
			_ = notifyIfUpgrade(ctx, url, messaging)
		}

		return db.SetTimeValue(lastUpgradeCheckDBKey, time.Now())
	})
}

func notifyIfUpgrade(ctx context.Context, url string, messaging io.Writer) error {
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "awless-client-"+Version)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	latest := struct {
		Version, URL string
	}{}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&latest); err == nil {
		if IsSemverUpgrade(Version, latest.Version) {
			var install string
			switch BuildFor {
			case "brew":
				install = "Run `brew upgrade awless`"
			case "zip", "targz":
				ext := "tar.gz"
				if runtime.GOOS == "windows" {
					ext = "zip"
				}
				// Matches GoReleaser's archive name template:
				// {{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}
				// The tag keeps its leading v; the version inside the file name
				// does not.
				install = fmt.Sprintf("Run `wget https://github.com/bootswithdefer/awless/releases/download/%s/awless_%s_%s_%s.%s`",
					latest.Version, strings.TrimPrefix(latest.Version, "v"), runtime.GOOS, runtime.GOARCH, ext)
			default:
				install = "Run `go install github.com/bootswithdefer/awless@latest`"
			}
			fmt.Fprintf(messaging, "New version %s available. Checkout the latest features at https://github.com/bootswithdefer/awless/blob/master/CHANGELOG.md\n%s\n", latest.Version, install)
		}
	}

	return nil
}

const semverLen = 3

type semver [semverLen]int

var ErrSemverInvalidFormat = errors.New("semver invalid format")

func IsSemverUpgrade(current, latest string) bool {
	i, err := CompareSemver(current, latest)
	if err != nil {
		return false
	}

	return i < 0
}

func CompareSemver(current, latest string) (int, error) {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	dot := func(r rune) bool {
		return r == '.'
	}
	cFields := strings.FieldsFunc(current, dot)
	lFields := strings.FieldsFunc(latest, dot)

	if len(cFields) != semverLen || len(lFields) != semverLen {
		return 0, ErrSemverInvalidFormat
	}

	currents := new(semver)
	for i, f := range cFields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return 0, ErrSemverInvalidFormat
		}
		currents[i] = num
	}

	latests := new(semver)
	for i, f := range lFields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return 0, ErrSemverInvalidFormat
		}
		latests[i] = num
	}

	for i := 0; i < semverLen; i++ {
		if latests[i] > currents[i] {
			return -1, nil
		} else if latests[i] == currents[i] {
			continue
		} else {
			return 1, nil
		}
	}

	return 0, nil
}
