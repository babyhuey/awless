# AWS Services To Add

Ranked by usefulness for infrastructure management. Services marked with `*` use the existing EC2 API client and are especially easy to add.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). A service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, so that step is easy to miss.

CloudTrail and CloudWatch Logs previously headed this list and have since been added, along with EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM and EFS. All of those are writable rather than list-only.

| Rank | Service | Resources | Complexity |
|------|---------|-----------|------------|
| 1 | ElastiCache | cache clusters, replication groups, cache subnet groups | Low |
| 2 | EventBridge | event buses, rules, targets | Low |
| 3 | Step Functions | state machines, executions | Low |
| 4 | WAF v2 | web ACLs, IP sets, rule groups | Medium |
| 5 | AWS Config | config rules, compliance status | Medium |
| 6 | Kinesis | streams | Low |
| 7 | Redshift | clusters, cluster subnet groups | Low |
| 8 | CodePipeline | pipelines, pipeline executions | Low |
| 9 | CodeBuild | build projects, builds | Low |
| 10 | Elastic Beanstalk | applications, environments | Low |
| 11 | Transit Gateway* | transit gateways, attachments, route tables | Low |
| 12 | VPC Endpoints* | VPC endpoints, endpoint services | Low |
| 13 | CodeDeploy | applications, deployment groups, deployments | Medium |
| 14 | Glue | databases, tables, crawlers, jobs | Medium |
| 15 | SES v2 | email identities, configuration sets | Low |
| 16 | Cognito | user pools, identity pools | Low |
| 17 | MSK | clusters, configurations | Low |
| 18 | Amazon MQ | brokers, configurations | Low |
| 19 | FSx | file systems, volumes, backups | Low |
| 20 | Global Accelerator | accelerators, listeners, endpoint groups | Medium |
| 21 | VPC Peering* | peering connections | Low |
| 22 | Cloud Map | namespaces, services, instances | Low |
| 23 | AWS Backup | backup plans, vaults, recovery points | Low |
