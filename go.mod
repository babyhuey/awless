module github.com/bootswithdefer/awless

go 1.26.1

toolchain go1.26.5

require (
	github.com/aws/aws-sdk-go-v2 v1.43.6
	github.com/aws/aws-sdk-go-v2/config v1.32.37
	github.com/aws/aws-sdk-go-v2/credentials v1.19.36
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.37
	github.com/aws/aws-sdk-go-v2/service/acm v1.44.1
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.37.6
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.45.6
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.72.1
	github.com/aws/aws-sdk-go-v2/service/backup v1.60.2
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.76.3
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.67.6
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.58.6
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.66.5
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.82.2
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.72.6
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.38.6
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.49.6
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.36.6
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.67.6
	github.com/aws/aws-sdk-go-v2/service/configservice v1.68.6
	github.com/aws/aws-sdk-go-v2/service/directconnect v1.44.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.63.3
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.321.2
	github.com/aws/aws-sdk-go-v2/service/ecr v1.60.6
	github.com/aws/aws-sdk-go-v2/service/ecs v1.90.2
	github.com/aws/aws-sdk-go-v2/service/efs v1.44.6
	github.com/aws/aws-sdk-go-v2/service/eks v1.91.1
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.56.6
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.37.6
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.36.6
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.58.7
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.48.6
	github.com/aws/aws-sdk-go-v2/service/fsx v1.68.6
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.38.6
	github.com/aws/aws-sdk-go-v2/service/glue v1.153.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.59.1
	github.com/aws/aws-sdk-go-v2/service/kafka v1.58.2
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.46.6
	github.com/aws/aws-sdk-go-v2/service/kms v1.55.6
	github.com/aws/aws-sdk-go-v2/service/lambda v1.101.4
	github.com/aws/aws-sdk-go-v2/service/mq v1.39.6
	github.com/aws/aws-sdk-go-v2/service/networkmanager v1.44.6
	github.com/aws/aws-sdk-go-v2/service/rds v1.124.3
	github.com/aws/aws-sdk-go-v2/service/redshift v1.65.6
	github.com/aws/aws-sdk-go-v2/service/route53 v1.65.8
	github.com/aws/aws-sdk-go-v2/service/s3 v1.107.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.44.6
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.43.6
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.66.6
	github.com/aws/aws-sdk-go-v2/service/sfn v1.45.6
	github.com/aws/aws-sdk-go-v2/service/sns v1.42.6
	github.com/aws/aws-sdk-go-v2/service/sqs v1.46.6
	github.com/aws/aws-sdk-go-v2/service/ssm v1.73.6
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.6
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.77.5
	github.com/aws/smithy-go v1.27.8
	github.com/boombuler/barcode v1.1.0
	github.com/bootswithdefer/triplestore v1.0.0
	github.com/chzyer/readline v1.5.1
	github.com/fatih/color v1.19.0
	github.com/go-git/go-git/v5 v5.19.2
	github.com/gorilla/mux v1.8.1
	github.com/oklog/ulid/v2 v2.1.2
	github.com/olekukonko/tablewriter v1.1.4
	github.com/spf13/cobra v1.10.2
	go.etcd.io/bbolt v1.5.0
	go.yaml.in/yaml/v3 v3.0.5
	golang.org/x/crypto v0.55.0
	golang.org/x/sync v0.22.0
	golang.org/x/term v0.45.0
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.30 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.6 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/cloudflare/circl v1.6.5 // indirect
	github.com/cyphar/filepath-securejoin v0.7.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.9.1 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.27 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.3.0 // indirect
	github.com/olekukonko/ll v0.1.8 // indirect
	github.com/pjbgf/sha1cd v0.6.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/skeema/knownhosts v1.3.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)
