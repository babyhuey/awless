# AWS Services To Add

**Nothing outstanding.** Every service on this list has been added.

`AGENTS.md` has the step-by-step procedure under [Adding a New AWS Service](./AGENTS.md#adding-a-new-aws-service). It is kept current because it records the steps that fail *silently* — a service that is generated but never registered in `aws/services/init.go` compiles fine and is invisible at runtime, and a wrong `awsName` is dropped without a word.

The README's [supported services table](./README.md#supported-aws-services) is the current inventory.

## Deliberately not included

These came up while working through the list and were left out for a reason, rather than missed.

| Resource | Why |
|---|---|
| Step Functions executions, CodePipeline executions, CodeDeploy deployment history | Run history rather than infrastructure, and listing costs one call per parent |
| Global Accelerator endpoint groups | A third level down, behind a listener, with a document-shaped configuration |
| FSx volumes | Only exist for two of the four file system types, and list per file system |
| MSK and Amazon MQ configurations | Versioned blobs of engine properties rather than infrastructure |
| Cloud Map instances | Registered by whatever runs the workload, usually ECS, rather than by hand |
| AWS Backup recovery points | The backups themselves; a delete command here exists only to destroy data |
| VPC endpoint services | The provider side of PrivateLink, which is a different job from consuming an endpoint |

Adding any of them is a normal service addition — the procedure in `AGENTS.md` covers it.
