# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge, Step Functions and WAF v2 have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| - | WAF v2 (partial) | web ACL and rule group **creation** only; listing and IP sets are done | Medium |
| 1 | AWS Config | config rules, compliance status | Medium |
| 2 | Kinesis | streams | Low |
| 3 | Redshift | clusters, cluster subnet groups | Low |
| 4 | CodePipeline | pipelines, pipeline executions | Low |
| 5 | CodeBuild | build projects, builds | Low |
| 6 | Elastic Beanstalk | applications, environments | Low |
| 7 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 8 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 9 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 10 | Glue | databases, tables, crawlers, jobs | Medium |
| 11 | SES v2 | email identities, configuration sets | Low |
| 12 | Cognito | user pools, identity pools | Low |
| 13 | MSK | clusters, configurations | Low |
| 14 | Amazon MQ | brokers, configurations | Low |
| 15 | FSx | file systems, volumes, backups | Low |
| 16 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 17 | VPC Peering* | peering connections | Low |
| 18 | Cloud Map | namespaces, services, instances | Low |
| 19 | AWS Backup | backup plans, vaults, recovery points | Low |
