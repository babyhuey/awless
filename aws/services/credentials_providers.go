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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"

	"github.com/bootswithdefer/awless/logger"
)

// stsCacheDuration is the duration for which STS credentials are cached.
// In AWS SDK v2, the default assumed role session duration is 15 minutes.
var stsCacheDuration = 15 * time.Minute

type cachedCredential struct {
	aws.Credentials
	Expiration time.Time
}

func (c *cachedCredential) isExpired() bool {
	return c.Expiration.Before(time.Now().UTC())
}

type fileCacheProvider struct {
	creds   aws.CredentialsProvider
	curr    *cachedCredential
	profile string
	log     *logger.Logger
}

func (f *fileCacheProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	awlessCache := os.Getenv("__AWLESS_CACHE")
	if awlessCache == "" {
		return f.creds.Retrieve(ctx)
	}
	credFolder := filepath.Join(awlessCache, "credentials")
	fold := &folder{credFolder}
	credFile := fmt.Sprintf("aws-profile-%s.json", f.profile)

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
		f.log.Errorf("%s\n", err)
		return credValue, err
	}

	switch credValue.Source {
	case stscreds.ProviderName:
		cred := &cachedCredential{credValue, time.Now().UTC().Add(stsCacheDuration)}
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
