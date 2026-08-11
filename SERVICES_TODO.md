# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

ElastiCache, EventBridge, Step Functions, WAF v2, AWS Config, Kinesis, Redshift, CodePipeline, CodeBuild and Elastic Beanstalk have since been added and removed from the list, as were CloudTrail and CloudWatch Logs before them, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| - | CodePipeline (partial) | pipeline **creation**; listing, deletion and running are done | Medium |
| - | WAF v2 (partial) | web ACL and rule group **creation** only; listing and IP sets are done | Medium |
| 1 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 2 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 3 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 4 | Glue | databases, tables, crawlers, jobs | Medium |
| 5 | SES v2 | email identities, configuration sets | Low |
| 6 | Cognito | user pools, identity pools | Low |
| 7 | MSK | clusters, configurations | Low |
| 8 | Amazon MQ | brokers, configurations | Low |
| 9 | FSx | file systems, volumes, backups | Low |
| 10 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 11 | VPC Peering* | peering connections | Low |
| 12 | Cloud Map | namespaces, services, instances | Low |
| 13 | AWS Backup | backup plans, vaults, recovery points | Low |
