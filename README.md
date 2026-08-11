<p align="center">
  <img src="docs/logo.svg" alt="awless" width="400">
</p>

<p align="center">
  <a href="https://github.com/bootswithdefer/awless/actions/workflows/ci.yml"><img src="https://github.com/bootswithdefer/awless/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://pkg.go.dev/github.com/bootswithdefer/awless"><img src="https://pkg.go.dev/badge/github.com/bootswithdefer/awless.svg" alt="Go Reference"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/bootswithdefer/awless" alt="License"></a>
  <img src="https://img.shields.io/github/go-mod/go-version/bootswithdefer/awless" alt="Go version">
  <img src="https://img.shields.io/badge/coverage-58%25-yellow" alt="Coverage">
  <img src="https://img.shields.io/badge/AWS%20SDK-v2-orange" alt="AWS SDK v2">
</p>

`awless` is a powerful, innovative and small surface command line interface (CLI) to manage Amazon Web Services.

> **This is a modernized fork of [wallix/awless](https://github.com/wallix/awless)**, which was last released in 2018. It migrates to AWS SDK v2, replaces `dep`/vendor with Go modules, adds 8 AWS services, and fixes a number of bugs that the SDK migration left behind. See [Changes in this fork](#changes-in-this-fork).

[Upstream Wiki](https://github.com/wallix/awless/wiki) | [Changelog](./CHANGELOG.md)

The upstream wiki is still the best guide to concepts — the template language, the sync
model and the getting-started tour are unchanged. Where this fork differs, it is noted
below.

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

### Homebrew

```sh
brew install --cask bootswithdefer/tap/awless
```

A cask rather than a formula because it ships the pre-built binary; Homebrew reserves
formulae for software built from source. Upgrade with `brew upgrade --cask awless`.

### From source

Requires Go 1.26+ (the version is pinned in `go.mod`):

```sh
go install github.com/bootswithdefer/awless@latest
```

Or clone and build:

```sh
git clone https://github.com/bootswithdefer/awless.git
cd awless
make build
```

### Pre-built binaries

Releases are built by GoReleaser for Linux, macOS and Windows on amd64 and arm64, with
checksums — see [Releases](https://github.com/bootswithdefer/awless/releases).

### Development setup

```sh
git config core.hooksPath .githooks   # gofmt + lint before each commit
make tools                            # install pinned dev tools
make verify                           # the full gate: fmt, vet, lint, race, vuln
```

`make help` lists every target. Contributors should read [AGENTS.md](./AGENTS.md), which
documents the code generation pipeline and the reflective command-spec system — the two
things most likely to surprise you.

### Configuration

If you have previously used the AWS CLI or aws-shell, you don't need to configure anything! Your config will be automatically loaded (i.e. `~/.aws/{credentials,config}`) and `awless` will prompt for any missing info (more at the [getting started wiki](https://github.com/wallix/awless/wiki/Getting-Started)).

# Supported AWS services

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
| **ElastiCache** | cache clusters, replication groups, cache subnet groups |
| **EventBridge** | event buses, rules, rule targets |
| **Step Functions** | state machines, executions |
| **WAF v2** | web ACLs, IP sets, rule groups |
| **AWS Config** | config rules, with compliance status |
| **Kinesis** | data streams |
| **Redshift** | clusters, cluster subnet groups |
| **CodePipeline** | pipelines |
| **CodeBuild** | build projects |
| **Elastic Beanstalk** | applications, environments |
| **Transit Gateway** | transit gateways, VPC attachments, route tables |
| **CodeDeploy** | applications, deployment groups, deployments |
| **Glue** | catalog databases, tables, crawlers, jobs |
| **SES v2** | email identities, configuration sets |
| **Cognito** | user pools, identity pools |
| **MSK** | Kafka clusters |
| **Amazon MQ** | brokers |
| **FSx** | file systems, backups |
| **Global Accelerator** | accelerators, listeners |
| **VPC Peering** | peering connections |
| **Cloud Map** | namespaces, services |
| **AWS Backup** | backup plans, vaults |
| **VPC Endpoints** | gateway and interface endpoints |

Services in **bold** are new in this fork. All of them support create/update/delete, not
just listing — see [Changes in this fork](#changes-in-this-fork).

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
      $ awless ls cacheclusters
      $ awless ls replicationgroups
      $ awless ls eventbuses
      $ awless ls eventrules
      $ awless ls statemachines
      $ awless ls ipsets
      $ awless ls webacls
      $ awless ls configrules
      $ awless ls streams
      $ awless ls redshiftclusters
      $ awless ls pipelines
      $ awless ls buildprojects
      $ awless ls applications
      $ awless ls environments
      $ awless ls transitgateways
      $ awless ls vpcendpoints
      $ awless ls deployapplications
      $ awless ls deploymentgroups
      $ awless ls gluedatabases
      $ awless ls crawlers
      $ awless ls jobs
      $ awless ls emailidentities
      $ awless ls userpools
      $ awless ls kafkaclusters
      $ awless ls brokers
      $ awless ls fsxfilesystems
      $ awless ls accelerators
      $ awless ls vpcpeerings
      $ awless ls namespaces
      $ awless ls backupplans
      $ ...
      (see awless ls -h)

- `awless run` : Create, update and delete complex infrastructures with smart defaults and sound auto-complete through awless templates.

      $ awless run ~/templates/my-infra.aws
      etc.

- **194 CRUD one-liners** integrated in the awless templating engine, each with `-h`
  documentation and worked examples:

      $ awless create instance -h
      $ awless create vpc -h
      $ awless attach policy -h
      $ awless create secret name=db-password secret=s3cr3t
      $ awless create ssmparameter name=/app/db/host value=db.internal type=SecureString
      $ awless create dynamodbtable name=users partition-key=id
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

### Platform

- **AWS SDK v2** throughout, replacing the v1 SDK
- **Go modules**, replacing `dep` and a vendor directory
- **Go 1.26**, with the toolchain pinned in `go.mod`
- **`github.com/bootswithdefer/triplestore` v1.0.0** for graph storage, the first tagged
  release that library has ever had
- **bbolt** replaces the archived `boltdb/bolt`, which crashes under the race detector

### New services

EKS, DynamoDB, Secrets Manager, KMS, API Gateway v2, SSM, EFS, CloudTrail and CloudWatch
Logs. All are writable, not list-only: 27 create/update/delete commands were added across
them.

### Bugs fixed

Nine of these came from the SDK v1 → v2 migration and were invisible to the compiler,
because `awless` maps template parameters onto AWS request structs by reflection:

- `create role`, `create/update distribution`, `attach/detach securitygroup` and four
  `containertask` commands all failed — the reflective call passed the request without the
  context that SDK v2 requires
- 35 parameters that take a list of strings panicked, because v1 modeled them as
  `[]*string` and v2 uses `[]string`
- 13 parameters that take a list of structs panicked, for the same reason
- `create instance count` reached neither `MinCount` nor `MaxCount`, from a one-character
  typo in a struct tag
- `update classicloadbalancer` silently sent an empty request: its field paths were
  spelled `Healthcheck` where the SDK field is `HealthCheck`
- `create distribution` silently dropped its origin domain, because indexed field paths
  such as `Origins.Items[0].DomainName` were not resolved
- `create/start/stop instance`, `create listener`, `create loadbalancer` and
  `create targetgroup` panicked on an empty AWS response
- `create vpc` without `name=` crashed while trying to tag the new VPC
- `awless log` crashed on a comment-only or empty template

And, found separately:

- **Secrets were written to the template log in cleartext** — command lines, fillers and
  messages all persisted `password=` values. They are now redacted at the persistence
  boundary.
- **The web UI listened on `0.0.0.0` with no authentication or timeouts.** It now defaults
  to loopback and requires an explicit opt-in to expose.
- **A deadlock in the IAM users fetcher**, plus an unsynchronized error flag.
- **Ten fan-out sites leaked goroutines** on the first error and ran unbounded, which
  throttled large accounts. All now use a bounded `errgroup`.
- **`awless` exited 0 on failure**, which broke shell scripting.
- **A denial of service in `triplestore`'s decoder**: an 11-byte input could request 3.8 GB
  and kill the process. Found with fuzzing.

### Upstream issues closed

- [#296](https://github.com/wallix/awless/issues/296): `--filter` now uses exact matching instead of substring matching
- [#281](https://github.com/wallix/awless/issues/281): Added `ebs-optimized` flag to `create instance`
- [#289](https://github.com/wallix/awless/issues/289): RDS endpoint now shown in `list databases`

### Removed

- **The scheduler is gone.** `--run-in`, `--revert-in` and the hidden `awless scheduler`
  command required a daemon that had been unbuildable since 2018 and exposed an
  unauthenticated endpoint that executed templates against the host's AWS credentials.
  Use EventBridge Scheduler, cron, or a CI schedule instead.

### Breaking changes

See [CHANGELOG.md](./CHANGELOG.md). The one to know about: `https-behaviour` is now
`https-behavior` on `create/update distribution`, for consistency with the American
spelling used everywhere else. Existing templates must be updated, and reverting a
previously logged distribution creation will fail because the stored command line carries
the old spelling.

### Engineering

- **200 acceptance tests covering all 194 commands**, running with no network access
- 51 linters enabled, including `errcheck`, `errorlint`, `noctx`, `contextcheck`,
  `musttag` and `bidichk`
- Native Go fuzzing, `govulncheck` in CI, and CI actions pinned to commit SHAs
- Ctrl-C cancels in-flight AWS calls, via a signal-cancelled root context threaded to
  every request

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

**Is this fork maintained?**

It is modernized rather than actively developed. CI is green, every command has an
acceptance test, and releases are published for Linux, macOS and Windows — install with
Homebrew or grab a binary from [Releases](https://github.com/bootswithdefer/awless/releases).

There is no open backlog. Decisions that look like gaps but are deliberate — `gosec` not
being a build gate, Classic ELB support being retained, five functions left untested — are
recorded under [Deliberate Omissions](./AGENTS.md#deliberate-omissions) in `AGENTS.md`.

# About

`awless` was originally created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and Francois-Xavier Aguessy at [WALLIX](https://github.com/wallix). This fork is maintained by [bootswithdefer](https://github.com/bootswithdefer).

`awless` is released under the Apache License.

    Disclaimer: Awless allows for easy resource creation with your cloud provider;
    we will not be responsible for any cloud costs incurred (even if you create a
    million instances using awless templates).

Contributors are welcome! Note that `awless` uses [triplestore](https://github.com/bootswithdefer/triplestore), another project originally developed at WALLIX and also forked here.
