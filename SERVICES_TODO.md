# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge, Step Functions, WAF v2 and AWS Config have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| - | WAF v2 (partial) | web ACL and rule group **creation** only; listing and IP sets are done | Medium |
| 1 | Kinesis | streams | Low |
| 2 | Redshift | clusters, cluster subnet groups | Low |
| 3 | CodePipeline | pipelines, pipeline executions | Low |
| 4 | CodeBuild | build projects, builds | Low |
| 5 | Elastic Beanstalk | applications, environments | Low |
| 6 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 7 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 8 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 9 | Glue | databases, tables, crawlers, jobs | Medium |
| 10 | SES v2 | email identities, configuration sets | Low |
| 11 | Cognito | user pools, identity pools | Low |
| 12 | MSK | clusters, configurations | Low |
| 13 | Amazon MQ | brokers, configurations | Low |
| 14 | FSx | file systems, volumes, backups | Low |
| 15 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 16 | VPC Peering* | peering connections | Low |
| 17 | Cloud Map | namespaces, services, instances | Low |
| 18 | AWS Backup | backup plans, vaults, recovery points | Low |
