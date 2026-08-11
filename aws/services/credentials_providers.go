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

package awsservices

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"

	"github.com/bootswithdefer/awless/logger"
)

// stsCacheDuration is how long cached STS credentials are trusted when the credential
// itself does not say. A credential that reports its own expiry is preferred; this is
// only the fallback, and matches the SDK's default assumed-role session of 15 minutes.
var stsCacheDuration = 15 * time.Minute

// expiryMargin keeps a credential from being served in the last moments of its life,
// where a long-running sync could start with a valid credential and finish without one.
const expiryMargin = 1 * time.Minute

// Serialized to the credentials cache file, so the field name is the on-disk key.
type cachedCredential struct {
	aws.Credentials
	Expiration time.Time `json:"Expiration"`
}

func (c *cachedCredential) isExpired() bool {
	return c.Expiration.Before(time.Now().UTC())
}

// cacheExpiry is when a retrieved credential should stop being served from cache.
//
// The credential's own expiry is authoritative: an assumed-role session honors
// duration_seconds, so a profile asking for an hour was previously re-prompting for MFA
// four times over that hour against a flat 15-minute assumption.
func cacheExpiry(c aws.Credentials) time.Time {
	if c.CanExpire && !c.Expires.IsZero() {
		return c.Expires.UTC().Add(-expiryMargin)
	}
	return time.Now().UTC().Add(stsCacheDuration)
}

type fileCacheProvider struct {
	creds   aws.CredentialsProvider
	mu      sync.Mutex
	curr    *cachedCredential
	profile string
	log     *logger.Logger
}

func (f *fileCacheProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	// Fetchers run concurrently, and the SDK asks for credentials once per request, so
	// this is a contended path.
	f.mu.Lock()
	defer f.mu.Unlock()

	// Serve from memory first. Without this the file below is stat'd and read on every
	// single AWS request.
	if f.curr != nil && !f.curr.isExpired() {
		return f.curr.Credentials, nil
	}

	awlessCache := os.Getenv("__AWLESS_CACHE")
	if awlessCache == "" {
		return f.creds.Retrieve(ctx)
	}
	credFolder := filepath.Join(awlessCache, "credentials")
	fold := &folder{credFolder}
	// An empty profile means the default one; without this the file is named
	// "aws-profile-.json".
	profile := f.profile
	if profile == "" {
		profile = "default"
	}
	credFile := fmt.Sprintf("aws-profile-%s.json", profile)

	if content, ok := fold.getFileContent(credFile); ok {
		var cached *cachedCredential
		if err := json.Unmarshal(content, &cached); err != nil {
			return aws.Credentials{}, err
		}
		f.log.ExtraVerbosef("loading credentials from '%s'", filepath.Join(credFolder, credFile))
		if !cached.isExpired() {
			f.curr = cached
			return cached.Credentials, nil
		}
	}
	credValue, err := f.creds.Retrieve(ctx)
	if err != nil {
		// Verbose rather than an error: this is returned to the caller, which reports
		// it. Logging at error level here printed every failure twice.
		f.log.ExtraVerbosef("retrieving credentials: %s", err)
		return credValue, err
	}

	switch credValue.Source {
	case stscreds.ProviderName:
		// Only assumed-role credentials are cached. Static keys are already on disk in
		// the shared credentials file, so copying them somewhere else would widen their
		// exposure for no benefit.
		cred := &cachedCredential{credValue, cacheExpiry(credValue)}
		f.curr = cred
		content, err := json.Marshal(cred)
		if err != nil {
			return credValue, err
		}
		if err = fold.putFileContent(credFile, content); err != nil {
			return credValue, fmt.Errorf("error writing cache file: %s", err.Error())
		}
		f.log.ExtraVerbosef("credentials cached in '%s'", filepath.Join(credFolder, credFile))
		return credValue, nil
	}
	return credValue, nil
}

type folder struct {
	path string
}

func (f *folder) getFileContent(filename string) (content []byte, ok bool) {
	if _, err := os.Stat(f.path); err != nil {
		return
	}
	credPath := filepath.Join(f.path, filename)

	if _, readerr := os.Stat(credPath); readerr != nil {
		return
	}
	var err error
	if content, err = os.ReadFile(credPath); err != nil {
		return
	}
	ok = true
	return
}

func (f *folder) putFileContent(filename string, content []byte) error {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		if err := os.MkdirAll(f.path, 0700); err != nil {
			return fmt.Errorf("creating credentials cache dir %s: %w", f.path, err)
		}
	}

	return os.WriteFile(filepath.Join(f.path, filename), content, 0600)
}
