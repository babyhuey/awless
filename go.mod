module github.com/wallix/awless

go 1.26.1

require (
	github.com/aws/aws-sdk-go-v2 v1.41.3
	github.com/aws/aws-sdk-go-v2/config v1.32.11
	github.com/aws/aws-sdk-go-v2/credentials v1.19.11
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.19
	github.com/aws/aws-sdk-go-v2/service/acm v1.37.21
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.33.7
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.41.12
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.64.2
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.71.7
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.60.2
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.55.7
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.55.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.64.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.56.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.294.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.55.4
	github.com/aws/aws-sdk-go-v2/service/ecs v1.73.1
	github.com/aws/aws-sdk-go-v2/service/efs v1.41.12
	github.com/aws/aws-sdk-go-v2/service/eks v1.80.2
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.33.21
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.54.8
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.5
	github.com/aws/aws-sdk-go-v2/service/kms v1.50.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.88.2
	github.com/aws/aws-sdk-go-v2/service/rds v1.116.2
	github.com/aws/aws-sdk-go-v2/service/route53 v1.62.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.96.4
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.41.3
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.13
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.23
	github.com/aws/aws-sdk-go-v2/service/ssm v1.68.2
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.8
	github.com/aws/smithy-go v1.24.2
	github.com/boltdb/bolt v1.3.1
	github.com/boombuler/barcode v1.1.0
	github.com/chzyer/readline v1.5.1
	github.com/fatih/color v1.18.0
	github.com/gorilla/mux v1.8.1
	github.com/mitchellh/ioprogress v0.0.0-20180201004757-6a23b12fa88e
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v1.1.3
	github.com/spf13/cobra v1.10.2
	github.com/wallix/awless-scheduler v0.0.6
	github.com/wallix/triplestore v0.0.0-20180213143850-4099dd913851
	golang.org/x/crypto v0.48.0
	golang.org/x/sync v0.22.0
	golang.org/x/term v0.40.0
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.16 // indirect
	github.com/clipperhouse/displaywidth v0.6.2 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20190725054713-01f96b0aa0cd // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.1.0 // indirect
	github.com/olekukonko/ll v0.1.4-0.20260115111900-9e59c2286df0 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/src-d/gcfg v1.4.0 // indirect
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)
