# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge and Step Functions have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| 1 | WAF v2 | web ACLs, IP sets, rule groups | Medium |
| 2 | AWS Config | config rules, compliance status | Medium |
| 3 | Kinesis | streams | Low |
| 4 | Redshift | clusters, cluster subnet groups | Low |
| 5 | CodePipeline | pipelines, pipeline executions | Low |
| 6 | CodeBuild | build projects, builds | Low |
| 7 | Elastic Beanstalk | applications, environments | Low |
| 8 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 9 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 10 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 11 | Glue | databases, tables, crawlers, jobs | Medium |
| 12 | SES v2 | email identities, configuration sets | Low |
| 13 | Cognito | user pools, identity pools | Low |
| 14 | MSK | clusters, configurations | Low |
| 15 | Amazon MQ | brokers, configurations | Low |
| 16 | FSx | file systems, volumes, backups | Low |
| 17 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 18 | VPC Peering* | peering connections | Low |
| 19 | Cloud Map | namespaces, services, instances | Low |
| 20 | AWS Backup | backup plans, vaults, recovery points | Low |
