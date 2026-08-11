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
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

	awsconfig "github.com/bootswithdefer/awless/aws/config"
	"github.com/bootswithdefer/awless/logger"
)

func ResolveRegionFromEnv() (region string) {
	cfg, err := newConfigResolver().resolve()
	if err == nil {
		region = cfg.Region
	}

	if awsconfig.IsValidRegion(region) {
		fmt.Fprintf(os.Stderr, "Found existing AWS region '%s'. Setting it as your default region.\n", region)
	} else {
		client := imds.NewFromConfig(cfg)
		output, imdsErr := client.GetRegion(context.Background(), &imds.GetRegionInput{})
		if imdsErr == nil {
			fmt.Fprintf(os.Stderr, "Found AWS region '%s' from local EC2 instance metadata. Setting it as your default region.\n", output.Region)
			region = output.Region
		}
	}

	if !awsconfig.IsValidRegion(region) {
		// Resolution is best-effort and the signature has no error to return, so an
		// aborted prompt leaves region empty for the caller to reject.
		selected, err := awsconfig.StdinRegionSelector()
		if err != nil {
			fmt.Fprintln(os.Stderr)
			return ""
		}
		region = selected
		fmt.Println()
	}

	return
}

type configResolver struct {
	region, profile           string
	profileSetterCallback     func(val string) error
	httpClient                *http.Client
	credentialHTTPClient      *http.Client
	logger                    *logger.Logger
	enableCredentialResolvers bool
}

func newConfigResolver() *configResolver {
	return &configResolver{
		credentialHTTPClient:  &http.Client{Timeout: 1 * time.Second},
		httpClient:            http.DefaultClient,
		profileSetterCallback: func(val string) error { return nil },
		logger:                logger.DiscardLogger,
	}
}

func (s *configResolver) withRegion(region string) *configResolver {
	s.region = region
	return s
}

func (s *configResolver) withProfile(profile string) *configResolver {
	s.profile = profile
	return s
}

func (s *configResolver) withCredentialResolvers() *configResolver {
	s.enableCredentialResolvers = true
	return s
}

func (s *configResolver) withProfileSetter(f func(val string) error) *configResolver {
	s.profileSetterCallback = f
	return s
}

func (s *configResolver) withLogger(l *logger.Logger) *configResolver {
	s.logger = l
	return s
}

func (s *configResolver) withNetworkMonitor(enableNetworkMonitor bool) *configResolver {
	// Network monitor request handlers are not supported in SDK v2.
	// This method is retained for API compatibility.
	return s
}

func (s *configResolver) resolve() (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	if s.region != "" {
		opts = append(opts, config.WithRegion(s.region))
	}
	if s.profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(s.profile))
	}

	opts = append(opts, config.WithHTTPClient(s.httpClient))

	if s.enableCredentialResolvers {
		// Without a token provider, a profile carrying mfa_serial does not prompt: the
		// SDK fails the whole load with "assume role with MFA enabled, but
		// AssumeRoleTokenProvider session option not set". This has to be passed to
		// LoadDefaultConfig rather than set afterwards, because the provider chain is
		// built during the load.
		opts = append(opts, config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
			o.TokenProvider = awsconfig.StdinMFATokenProvider
		}))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return aws.Config{}, err
	}

	if s.enableCredentialResolvers {
		// Cache assumed-role credentials on disk, keyed by profile, so MFA is entered
		// once per session rather than once per awless invocation.
		//
		// Not wrapped in aws.NewCredentialsCache: cfg.Credentials is already one, and
		// nesting two made every credential error report "failed to refresh cached
		// credentials" twice. fileCacheProvider holds its own in-memory copy instead.
		cfg.Credentials = &fileCacheProvider{
			creds:   cfg.Credentials,
			profile: s.profile,
			log:     s.logger,
		}
		if _, err = cfg.Credentials.Retrieve(context.Background()); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}
