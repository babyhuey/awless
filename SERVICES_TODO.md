# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge, Step Functions, WAF v2, AWS Config, Kinesis, Redshift, CodePipeline, CodeBuild, Elastic Beanstalk, Transit Gateway, VPC Endpoints, CodeDeploy, Glue and SES v2 have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| 1 | Cognito | user pools, identity pools | Low |
| 2 | MSK | clusters, configurations | Low |
| 3 | Amazon MQ | brokers, configurations | Low |
| 4 | FSx | file systems, volumes, backups | Low |
| 5 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 6 | VPC Peering* | peering connections | Low |
| 7 | Cloud Map | namespaces, services, instances | Low |
| 8 | AWS Backup | backup plans, vaults, recovery points | Low |
