package awsfetch

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/bootswithdefer/awless/logger"
)

type AWSAPI struct {
	IAM                    *iam.Client
	Ec2                    *ec2.Client
	Elbv2                  *elasticloadbalancingv2.Client
	Elb                    *elasticloadbalancing.Client
	RDS                    *rds.Client
	Autoscaling            *autoscaling.Client
	ECR                    *ecr.Client
	ECS                    *ecs.Client
	Applicationautoscaling *applicationautoscaling.Client
	STS                    *sts.Client
	S3                     *s3.Client
	SNS                    *sns.Client
	SQS                    *sqs.Client
	Route53                *route53.Client
	Lambda                 *lambda.Client
	Cloudwatch             *cloudwatch.Client
	Cloudfront             *cloudfront.Client
	Cloudformation         *cloudformation.Client
	ACM                    *acm.Client
	EKS                    *eks.Client
	Dynamodb               *dynamodb.Client
	Secretsmanager         *secretsmanager.Client
	KMS                    *kms.Client
	Apigatewayv2           *apigatewayv2.Client
	SSM                    *ssm.Client
	EFS                    *efs.Client
	Cloudtrail             *cloudtrail.Client
	Cloudwatchlogs         *cloudwatchlogs.Client
	Elasticache            *elasticache.Client
}

type Config struct {
	Log   *logger.Logger
	Extra map[string]any
	APIs  *AWSAPI
}

func NewConfig(apis ...any) *Config {
	c := &Config{
		Extra: make(map[string]any),
		Log:   logger.DiscardLogger,
	}
	assignAPIs(c, apis...)
	return c
}

func (c *Config) getBoolDefaultTrue(key string) bool {
	if c.Extra == nil {
		return true
	}

	if b, ok := c.Extra[key].(bool); ok {
		return b
	}

	return true
}

func assignAPIs(c *Config, apis ...any) {
	c.APIs = new(AWSAPI)
	val := reflect.ValueOf(c.APIs).Elem()
	stru := val.Type()

	for _, api := range apis {
		if !reflect.ValueOf(api).IsValid() {
			continue
		}

		apiType := reflect.TypeOf(api)
		for i := range stru.NumField() {
			fieldType := stru.Field(i).Type
			if apiType.AssignableTo(fieldType) {
				val.Field(i).Set(reflect.ValueOf(api))
				break
			}
		}
	}
}
