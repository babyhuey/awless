
[![CI](https://github.com/bootswithdefer/awless/actions/workflows/ci.yml/badge.svg)](https://github.com/bootswithdefer/awless/actions/workflows/ci.yml)
[![Go Coverage](https://github.com/bootswithdefer/awless/wiki/coverage.svg)](https://raw.githack.com/wiki/bootswithdefer/awless/coverage.html)

`awless` is a powerful, innovative and small surface command line interface (CLI) to manage Amazon Web Services.

> This is a modernized fork of [wallix/awless](https://github.com/wallix/awless). It migrates from AWS SDK v1 to v2, uses Go modules, and adds support for newer AWS services.

[Upstream Wiki](https://github.com/wallix/awless/wiki) | [Changelog](https://github.com/wallix/awless/blob/master/CHANGELOG.md#readme)

# Why awless

`awless` stands out by having the following characteristics:

- small and hierarchical set of commands
- a simple/powerful text [templating language](https://github.com/wallix/awless/wiki/Templates) to create and **revert** fully-fledged infrastructures
- wrapping/composing AWS API calls when necessary to enrich behavior. Ex: ensure smart defaults, security best practices, etc.
- local log of all your cloud modifications done through `awless` to list/revert past actions
- sync to a local graph storage of your cloud representation
- exploration of your cloud infrastructure and resources interrelations, **even offline** using the local graph storage
- clearer and flexible terminal output's with: numerous formats (machine/human friendly), enriched resources's properties/relations when feasible
- connect easily using awless' **smart SSH** to your private & public instances

For more read our [FAQ](#faq) below (how `awless` compares to other tools, etc.)

# Install

### From source (recommended)

Requires Go 1.21+:

```sh
go install github.com/bootswithdefer/awless@latest
```

Or clone and build:

```sh
git clone https://github.com/bootswithdefer/awless.git
cd awless
go build -o awless .
```

### Development setup

After cloning, enable the pre-commit hooks to catch formatting and lint errors before pushing:

```sh
git config core.hooksPath .githooks
```

### Pre-built binaries

Download the latest `awless` binaries (Windows/Linux/macOS) [from Github releases](https://github.com/bootswithdefer/awless/releases/latest).

### Configuration

If you have previously used the AWS CLI or aws-shell, you don't need to configure anything! Your config will be automatically loaded (i.e. `~/.aws/{credentials,config}`) and `awless` will prompt for any missing info (more at the [getting started wiki](https://github.com/wallix/awless/wiki/Getting-Started)).

# Supported AWS services

`awless` supports the following AWS services:

| Service | Resources |
|---------|-----------|
| EC2 | instances, vpcs, subnets, security groups, keypairs, volumes, snapshots, images, elastic IPs, network interfaces, NAT gateways, internet gateways, route tables |
| IAM | users, groups, roles, policies, access keys, instance profiles, login profiles, MFA devices |
| S3 | buckets, s3objects |
| RDS | databases, DB subnet groups |
| ELBv2 | load balancers, target groups, listeners |
| Auto Scaling | launch configurations, scaling groups, scaling policies |
| Lambda | functions |
| SNS | topics, subscriptions |
| SQS | queues |
| Route53 | zones, records |
| CloudWatch | alarms, metrics |
| CloudFront | distributions |
| ECR | registries |
| ECS | container clusters, container tasks, containers |
| ACM | certificates |
| CloudFormation | stacks |
| Application Auto Scaling | app scaling targets, app scaling policies |
| **EKS** | clusters, node groups |
| **DynamoDB** | tables |
| **Secrets Manager** | secrets |
| **KMS** | keys |
| **API Gateway v2** | APIs, routes, stages |
| **Systems Manager (SSM)** | parameters |
| **EFS** | file systems, mount targets |
| **CloudTrail** | trails |
| **CloudWatch Logs** | log groups |

Services in **bold** are new in this fork.

# Main features

- **Aliasing of resources through their natural name** so you don't have to always use cryptic ids that are impossible to remember
- `awless show` : Explore the properties, relations, dependencies of a specific resource (even offline thanks to the sync) given only a *name* (or id/arn).

      $ awless show jsmith --local

- `awless list` : Clear and easy listing of multi-region cloud resources with filters via *resources properties* or *resources tags*.

      $ awless list instances --sort uptime --local
      $ awless list users --format csv --columns name,created
      $ awless list volumes --filter state=in-use --filter type=gp2
      $ awless list volumes --tag-value Purchased
      $ awless ls vpcs --tag-key Dept --tag-key Internal --format tsv
      $ awless ls instances --tag Env=Production,Dept=Marketing
      $ awless ls instances --filter state=running,type=t2.micro --format json
      $ awless ls s3objects --filter bucket=pdf-bucket -r us-west-2
      $ awless ls eksclusters
      $ awless ls dynamodbtables
      $ awless ls secrets
      $ awless ls ssmparameters
      $ awless ls filesystems
      $ ...
      (see awless ls -h)

- `awless run` : Create, update and delete complex infrastructures with smart defaults and sound auto-complete through awless templates.

      $ awless run ~/templates/my-infra.aws
      etc.

- **Hundreds of powerful CRUD CLI one-liners** integrated in the awless templating engine:

      $ awless create instance -h
      $ awless create vpc -h
      $ awless attach policy -h
      $ ...
      (see awless -h)

- `awless log` : Detailed and easy reporting of all the CLI template executions
- `awless revert` : Revert of executed templates and resources creation
- Create instances straight from a distro name. No need to know the region or AMI ;) (_free tier community bare distro only_, see `awless create instance -h`)

      $ awless create instance distro=debian
      $ awless create instance distro=coreos
      $ awless create instance distro=redhat::7.2 type=t2.micro
      $ awless create instance distro=debian:debian:jessie lock=true
      $ awless create instance distro=amazonlinux:amzn2
      $ awless create instance type=t2.micro ebs-optimized=true
      etc.

- Leveraging AWS `userdata` to provision instance on creation from remote (i.e http) or local scripts: `awless create instance ... userdata=/home/john/...`
- `awless ssh` : Clean and simple SSH to public & private instances using only a name

      $ awless ssh my-production-instance
      $ awless ssh redis-prod --through jump-server
      $ awless ssh 34.215.29.221
      $ awless ssh db-private --private
      $ awless ssh 172.31.77.151 --port 2222 --through my-proxy --through-port 23
      $ ...
      (see awless ssh -h)

- `awless switch` : Switch easily between AWS accounts (i.e. profile) and regions

      $ awless switch admin eu-west-2
      $ awless switch us-west-1
      $ awless switch mfa
      etc.

- `awless` transparently syncs cloud resources locally to a graph representation in order for the CLI to leverage data and their relations in other awless commands and in an offline manner ([more on the sync](https://github.com/wallix/awless/wiki/Getting-Started#sync))
- `awless sync` : Explicit and manual command to fetch & store resources locally. Then query & inspect your cloud offline
- Output listing formats either human (**default display is Markdown-compatible tables**) or machine readable (csv, tsv, json, ...): `--format`
- `awless inspect` : Leverage **experimental** and community inspectors which are interface implementation utilities to run analysis on your cloud resources graphs

      $ awless inspect -i bucket_sizer
      (see awless inspect -h)

- `awless completion` : CLI autocompletion for Unix/Linux's bash and zsh

# Getting started

Take the tour at [Getting Started (wiki)](https://github.com/wallix/awless/wiki/Getting-Started).

# Changes in this fork

- **AWS SDK v2**: Migrated from AWS SDK for Go v1 to v2
- **Go modules**: Replaced `dep`/vendor with Go modules
- **6 new AWS services**: EKS, DynamoDB, Secrets Manager/KMS, API Gateway v2, SSM, EFS
- **Bug fixes from upstream issues**:
  - [#296](https://github.com/wallix/awless/issues/296): `--filter` now uses exact matching instead of substring matching
  - [#281](https://github.com/wallix/awless/issues/281): Added `ebs-optimized` flag to `create instance`
  - [#289](https://github.com/wallix/awless/issues/289): RDS endpoint now shown in `list databases`

# FAQ

**There are already some AWS CLIs. What is `awless` unique approach?**

Three things that differentiate `awless` from other AWS CLIs:

* It has its own **compiled and very simple templating language** to build AWS infrastructures.
* Commands are made of _VERB + ENTITY [+ param=value]_ and are actually valid lines of the template language.
* It transparently syncs to a local graph a representation of the cloud resources and their relations.

**How do you create infrastructure with `awless`?**

You build infrastructure using `template files` or `command one-liners` that get compiled and run through `awless` builtin engine. Learn [more about the way templates work](https://github.com/wallix/awless/wiki/Templates).

Note that all your actions against the cloud are logged. Templates are revertible/rollbackable.

**How does `awless` compare to Terraform?**

Terraform is much broader in scope. `awless` takes a different approach:

- Favors simplicity with a straight forward, compiled and simple deployment language
- Employs an all-or-nothing deployment: does not keep state
- Provides rollback on any ran template
- Logs all actions against the cloud with rich, revertable logs

# About

`awless` was originally created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and Francois-Xavier Aguessy at [WALLIX](https://github.com/wallix). This fork is maintained by [bootswithdefer](https://github.com/bootswithdefer).

`awless` is released under the Apache License.

    Disclaimer: Awless allows for easy resource creation with your cloud provider;
    we will not be responsible for any cloud costs incurred (even if you create a
    million instances using awless templates).

Contributors are welcome! Note that `awless` uses [triplestore](https://github.com/bootswithdefer/triplestore) another project developed at WALLIX.
