# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge, Step Functions, WAF v2, AWS Config and Kinesis have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| - | WAF v2 (partial) | web ACL and rule group **creation** only; listing and IP sets are done | Medium |
| 1 | Redshift | clusters, cluster subnet groups | Low |
| 2 | CodePipeline | pipelines, pipeline executions | Low |
| 3 | CodeBuild | build projects, builds | Low |
| 4 | Elastic Beanstalk | applications, environments | Low |
| 5 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 6 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 7 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 8 | Glue | databases, tables, crawlers, jobs | Medium |
| 9 | SES v2 | email identities, configuration sets | Low |
| 10 | Cognito | user pools, identity pools | Low |
| 11 | MSK | clusters, configurations | Low |
| 12 | Amazon MQ | brokers, configurations | Low |
| 13 | FSx | file systems, volumes, backups | Low |
| 14 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 15 | VPC Peering* | peering connections | Low |
| 16 | Cloud Map | namespaces, services, instances | Low |
| 17 | AWS Backup | backup plans, vaults, recovery points | Low |
