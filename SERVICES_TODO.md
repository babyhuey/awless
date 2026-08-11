# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge, Step Functions, WAF v2, AWS Config, Kinesis, Redshift and CodePipeline have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| - | CodePipeline (partial) | pipeline **creation**; listing, deletion and running are done | Medium |
| - | WAF v2 (partial) | web ACL and rule group **creation** only; listing and IP sets are done | Medium |
| 1 | CodeBuild | build projects, builds | Low |
| 2 | Elastic Beanstalk | applications, environments | Low |
| 3 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 4 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 5 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 6 | Glue | databases, tables, crawlers, jobs | Medium |
| 7 | SES v2 | email identities, configuration sets | Low |
| 8 | Cognito | user pools, identity pools | Low |
| 9 | MSK | clusters, configurations | Low |
| 10 | Amazon MQ | brokers, configurations | Low |
| 11 | FSx | file systems, volumes, backups | Low |
| 12 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 13 | VPC Peering* | peering connections | Low |
| 14 | Cloud Map | namespaces, services, instances | Low |
| 15 | AWS Backup | backup plans, vaults, recovery points | Low |
