# awless Usage Guide

A comprehensive guide to all awless commands and supported resources.

## Table of Contents

- [Quick Start](#quick-start)
- [Core Commands](#core-commands)
  - [list / ls](#list--ls)
  - [show](#show)
  - [create](#create)
  - [delete](#delete)
  - [update](#update)
  - [attach / detach](#attach--detach)
  - [start / stop / restart](#start--stop--restart)
  - [check](#check)
  - [copy / import](#copy--import)
  - [ssh](#ssh)
  - [run](#run)
  - [sync](#sync)
  - [log](#log)
  - [revert](#revert)
  - [inspect](#inspect)
  - [switch](#switch)
  - [whoami](#whoami)
  - [config](#config)
  - [completion](#completion)
- [Service Reference](#service-reference)
- [Output & Filtering](#output--filtering)
- [Templates](#templates)

---

## Quick Start

```bash
# Build from source
git clone https://github.com/babyhuey/awless.git
cd awless && go build -o awless .

# Your ~/.aws/credentials and ~/.aws/config are used automatically
awless whoami

# List resources
awless list instances
awless list loggroups

# Show a specific resource
awless show my-instance

# Create resources
awless create instance type=t2.micro distro=debian
```

---

## Core Commands

### list / ls

List cloud resources with filtering, sorting, and formatting.

```bash
awless list <resource-type>     # List a specific resource type
awless list <service-name>      # List ALL resources in a service
awless ls <resource-type>       # 'ls' is an alias for 'list'
```

**Flags:**

| Flag | Description | Example |
|------|-------------|---------|
| `--format` | Output format: `table`, `csv`, `tsv`, `json` | `--format json` |
| `--filter` | Filter by property (case insensitive, exact match) | `--filter state=running` |
| `--tag` | Filter by EC2 tag key=value (case sensitive) | `--tag Env=Production` |
| `--tag-key` | Filter by tag key existence | `--tag-key Department` |
| `--tag-value` | Filter by tag value | `--tag-value Staging` |
| `--columns` | Select columns to display | `--columns id,name,state` |
| `--sort` | Sort by column | `--sort uptime` |
| `--reverse` | Reverse sort order | `--sort created --reverse` |
| `--ids` | Output only resource IDs | `--ids` |
| `--no-headers` | Suppress column headers | `--no-headers` |
| `--local` | Use locally cached data (offline) | `--local` |
| `-r` | Specify region | `-r us-west-2` |

**Examples:**

```bash
# EC2
awless list instances --sort uptime
awless list instances --filter state=running,type=t2.micro
awless list instances --tag Env=Production,Dept=Marketing
awless list vpcs --format json
awless list securitygroups --filter vpc=vpc-12345

# IAM
awless list users --format csv --columns name,created
awless list roles
awless list policies
awless list groups

# S3
awless list buckets
awless list s3objects --filter bucket=my-bucket

# CloudWatch Logs
awless list loggroups
awless list loggroups --format json
awless list loggroups --sort created --reverse
awless list loggroups --filter name=/aws/lambda

# EKS
awless list eksclusters
awless list eksnodegroups

# DynamoDB
awless list dynamodbtables

# Secrets Manager / KMS
awless list secrets
awless list keys

# API Gateway
awless list apigateways
awless list apigatewayroutes
awless list apigatewaystages

# SSM
awless list ssmparameters

# EFS
awless list filesystems
awless list mounttargets

# CloudTrail
awless list trails

# DNS
awless list zones
awless list records

# Monitoring
awless list alarms
awless list metrics

# Containers
awless list containercluster
awless list containertasks

# List entire service at once
awless list infra
awless list access
awless list storage
awless list dns
awless list monitoring
awless list cloudwatchlogs
```

---

### show

Display detailed properties and relations of a specific resource.

```bash
awless show <resource-name-or-id>
awless show <resource-name-or-id> --local   # offline mode
```

You can reference resources by name, ID, or ARN:

```bash
awless show my-instance
awless show i-0abc123def456
awless show /aws/lambda/my-function    # loggroup by name
awless show my-vpc --local             # use cached data
```

---

### create

Create AWS resources. Use `-h` on any create command to see required and optional parameters.

```bash
awless create <resource-type> param1=value1 param2=value2
awless create <resource-type> -h    # show help for this resource
```

**EC2 examples:**

```bash
# Instances
awless create instance type=t2.micro distro=debian
awless create instance type=t2.micro distro=amazonlinux:amzn2
awless create instance type=t2.micro distro=redhat::7.2 lock=true
awless create instance type=t2.micro image=ami-123456 subnet=sub-456 keypair=mykey
awless create instance type=t2.micro ebs-optimized=true userdata=/path/to/script.sh

# Networking
awless create vpc cidr=10.0.0.0/16 name=my-vpc
awless create subnet cidr=10.0.1.0/24 vpc=my-vpc availabilityzone=us-east-1a
awless create securitygroup vpc=my-vpc description="Web servers" name=web-sg
awless create internetgateway
awless create natgateway elasticip-id=eipalloc-123 subnet=my-subnet
awless create routetable vpc=my-vpc
awless create route table=rtb-123 cidr=0.0.0.0/0 gateway=igw-456
awless create elasticip domain=vpc

# Storage
awless create volume availabilityzone=us-east-1a size=100 type=gp2
awless create snapshot volume=vol-123 description="backup"
awless create keypair name=my-keypair
```

**IAM examples:**

```bash
awless create user name=jsmith
awless create group name=developers
awless create role name=lambda-role principal-service="lambda.amazonaws.com"
awless create policy name=my-policy effect=Allow action="s3:GetObject" resource="arn:aws:s3:::my-bucket/*"
awless create accesskey user=jsmith
awless create loginprofile username=jsmith password=TempPass123!
awless create instanceprofile name=my-profile
awless create mfadevice name=my-mfa user=jsmith
```

**S3 examples:**

```bash
awless create bucket name=my-bucket acl=private
awless create s3object bucket=my-bucket file=/path/to/file key=folder/file.txt
```

**RDS examples:**

```bash
awless create database engine=postgres id=mydb size=db.t2.micro dbname=appdb username=admin password=secret123
awless create dbsubnetgroup name=my-dbsubnet description="DB subnets" subnets=[sub-123,sub-456]
```

**Load Balancer examples:**

```bash
awless create loadbalancer name=my-alb subnets=[sub-123,sub-456] securitygroups=[sg-789]
awless create targetgroup name=my-tg port=80 protocol=HTTP vpc=vpc-123
awless create listener actiontype=forward loadbalancer=arn:... port=80 protocol=HTTP targetgroup=arn:...
```

**Lambda:**

```bash
awless create function handler=index.handler name=my-func role=arn:aws:iam::123:role/lambda-role runtime=python3.9
```

**DNS:**

```bash
awless create zone callerreference=ref1 name=example.com
awless create record zone=/hostedzone/Z123 name=www.example.com type=A value=1.2.3.4 ttl=300
```

**SNS / SQS:**

```bash
awless create topic name=my-topic
awless create subscription topic=arn:aws:sns:... protocol=email endpoint=user@example.com
awless create queue name=my-queue
```

**CloudFormation:**

```bash
awless create stack name=my-stack template-file=/path/to/template.yml
```

**Other services:**

```bash
awless create certificate domains=example.com            # ACM
awless create distribution origin-domain=mybucket.s3.amazonaws.com  # CloudFront
awless create repository name=my-ecr-repo                # ECR
awless create containercluster name=my-cluster            # ECS
awless create tag resource=i-123 key=Env value=Production # Tags
```

**Supported create resource types (45):**
accesskey, alarm, appscalingpolicy, appscalingtarget, bucket, certificate, classicloadbalancer, containercluster, database, dbsubnetgroup, distribution, elasticip, function, group, image, instance, instanceprofile, internetgateway, keypair, launchconfiguration, listener, loadbalancer, loginprofile, mfadevice, natgateway, networkinterface, policy, queue, record, repository, role, route, routetable, s3object, scalinggroup, scalingpolicy, securitygroup, snapshot, stack, subnet, subscription, tag, targetgroup, topic, user, volume, vpc, zone

---

### delete

Delete AWS resources.

```bash
awless delete <resource-type> <params>
awless delete <resource-type> -h    # show parameters
```

**Examples:**

```bash
awless delete instance id=i-0abc123
awless delete vpc id=vpc-123
awless delete securitygroup id=sg-456
awless delete volume id=vol-789
awless delete user name=jsmith
awless delete role name=old-role
awless delete bucket name=my-bucket
awless delete record zone=/hostedzone/Z123 name=www.example.com type=A value=1.2.3.4 ttl=300
awless delete function name=my-func      # Lambda
awless delete topic arn=arn:aws:sns:...
```

---

### update

Update existing resources.

```bash
awless update <resource-type> <params>
```

**Examples:**

```bash
awless update instance id=i-123 type=t2.large
awless update securitygroup id=sg-123 inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=443
awless update securitygroup id=sg-123 inbound=revoke protocol=tcp cidr=10.0.0.0/8 portrange=22
awless update subnet id=sub-123 public=true
awless update s3object bucket=my-bucket key=file.txt acl=public-read
awless update record zone=/hostedzone/Z123 name=www.example.com type=A value=2.3.4.5 ttl=60
awless update scalinggroup name=my-asg max-size=10 min-size=2
awless update image id=ami-123 operation=add description="Updated image"
awless update policy arn=arn:aws:iam::123:policy/my-policy
awless update bucket name=my-bucket acl=public-read
awless update distribution id=E123 enable=true
```

**Supported update resource types (15):**
bucket, classicloadbalancer, containertask, distribution, image, instance, loginprofile, policy, record, s3object, scalinggroup, securitygroup, stack, subnet, targetgroup

---

### attach / detach

Attach or detach resources to/from each other.

```bash
awless attach <resource-type> <params>
awless detach <resource-type> <params>
```

**Examples:**

```bash
# Attach policy to user/group/role
awless attach policy arn=arn:aws:iam::123:policy/ReadOnly user=jsmith
awless attach policy arn=arn:aws:iam::123:policy/ReadOnly group=developers
awless attach policy arn=arn:aws:iam::123:policy/ReadOnly role=my-role
awless detach policy arn=arn:aws:iam::123:policy/ReadOnly user=jsmith

# Attach user to group
awless attach user name=jsmith group=developers
awless detach user name=jsmith group=developers

# Attach volume to instance
awless attach volume id=vol-123 instance=i-456 device=/dev/sdf
awless detach volume id=vol-123 instance=i-456 device=/dev/sdf

# Attach internet gateway to VPC
awless attach internetgateway id=igw-123 vpc=vpc-456
awless detach internetgateway id=igw-123 vpc=vpc-456

# Attach elastic IP
awless attach elasticip id=eipalloc-123 instance=i-456

# Attach security group to instance
awless attach securitygroup id=sg-123 instance=i-456

# Attach route table to subnet
awless attach routetable id=rtb-123 subnet=sub-456

# Attach role to instance profile
awless attach role name=my-role instanceprofile=my-profile

# Attach instance to ELB
awless attach instance id=i-123 targetgroup=arn:...

# Attach MFA device
awless attach mfadevice id=arn:... user=jsmith mfa-code-1=123456 mfa-code-2=789012

# Attach network interface
awless attach networkinterface id=eni-123 instance=i-456
```

---

### start / stop / restart

Control the lifecycle of running resources.

```bash
awless start instance id=i-123
awless stop instance id=i-123
awless restart instance id=i-123

awless start database id=mydb
awless stop database id=mydb
awless restart database id=mydb

awless start alarm name=my-alarm    # enable
awless stop alarm name=my-alarm     # disable

awless start containertask ...
awless stop containertask ...
```

---

### check

Wait for a resource to reach a specific state.

```bash
awless check instance id=i-123 state=running timeout=180
awless check database id=mydb state=available timeout=600
awless check loadbalancer id=arn:... state=active timeout=300
awless check natgateway id=nat-123 state=available
awless check certificate arn=arn:... state=issued
awless check distribution id=E123 state=Deployed
awless check securitygroup id=sg-123 state=unused
awless check volume id=vol-123 state=available
awless check scalinggroup name=my-asg count=3 timeout=300
```

---

### copy / import

```bash
awless copy image name=my-copy source-id=ami-123 source-region=us-west-2
awless copy snapshot source-id=snap-123 source-region=eu-west-1
awless import image ...
```

---

### ssh

SSH into EC2 instances using name, ID, or IP.

```bash
awless ssh my-instance                               # by name
awless ssh i-0abc123def456                            # by instance ID
awless ssh 34.215.29.221                              # by IP
awless ssh my-instance --through jump-server          # via bastion/jump host
awless ssh db-private --private                       # connect to private IP
awless ssh 172.31.77.151 --port 2222 --through my-proxy --through-port 23
```

---

### run

Execute template files (`.aws` extension) for complex multi-step operations.

```bash
awless run /path/to/template.aws
awless run https://raw.githubusercontent.com/.../template.aws
```

Templates use the awless template language:

```
# my-infra.aws
vpc = create vpc cidr=10.0.0.0/16 name=my-vpc
subnet = create subnet cidr=10.0.1.0/24 vpc=$vpc name=public
igw = create internetgateway
attach internetgateway id=$igw vpc=$vpc
create instance subnet=$subnet type=t2.micro distro=debian name=web-server
```

See the [Templates wiki](https://github.com/wallix/awless/wiki/Templates) for more.

---

### sync

Manually sync cloud resources to local graph storage for offline use.

```bash
awless sync            # sync all services
```

After syncing, use `--local` with list/show for offline queries:

```bash
awless list instances --local
awless show my-vpc --local
```

Configure which resources to sync:

```bash
awless config set aws.cloudwatchlogs.sync false          # disable entire service
awless config set aws.cloudwatchlogs.loggroup.sync false  # disable specific resource
```

---

### log

View the history of all cloud actions performed through awless.

```bash
awless log           # show all actions
awless log --filter delete    # filter log entries
```

---

### revert

Rollback a previously executed template or action.

```bash
awless log                     # find the action ID
awless revert <revert-id>      # undo that action
```

---

### inspect

Run analysis inspectors on your cloud resources.

```bash
awless inspect -i bucket_sizer
awless inspect -h              # list available inspectors
```

---

### switch

Switch between AWS profiles and regions.

```bash
awless switch admin eu-west-2    # switch profile and region
awless switch us-west-1          # switch region only
awless switch mfa                # switch to MFA profile
```

---

### whoami

Display your current AWS identity.

```bash
awless whoami
```

---

### config

Manage awless configuration.

```bash
awless config list                          # show all config
awless config set aws.region us-west-2      # set a value
awless config get aws.region                # get a value
awless config unset aws.region              # remove a value
```

---

### completion

Set up shell autocompletion.

```bash
# Bash
source <(awless completion bash)

# Zsh
source <(awless completion zsh)
```

---

## Service Reference

All resource types, grouped by service name. Use the pluralized form with `awless list`.

| Service Name | Resource Types (singular) | List Command |
|-------------|--------------------------|--------------|
| **infra** | instance, subnet, vpc, keypair, securitygroup, volume, internetgateway, natgateway, routetable, availabilityzone, image, importimagetask, elasticip, snapshot, networkinterface, classicloadbalancer, loadbalancer, targetgroup, listener, database, dbsubnetgroup, launchconfiguration, scalinggroup, scalingpolicy, repository, containercluster, containertask, container, containerinstance, certificate | `awless ls instances`, `awless ls vpcs`, etc. |
| **access** | user, group, role, policy, accesskey, instanceprofile, mfadevice | `awless ls users`, `awless ls roles`, etc. |
| **storage** | bucket, s3object | `awless ls buckets`, `awless ls s3objects` |
| **messaging** | subscription, topic, queue | `awless ls topics`, `awless ls queues` |
| **dns** | zone, record | `awless ls zones`, `awless ls records` |
| **lambda** | function | `awless ls functions` |
| **monitoring** | metric, alarm | `awless ls metrics`, `awless ls alarms` |
| **cdn** | distribution | `awless ls distributions` |
| **cloudformation** | stack | `awless ls stacks` |
| **eks** | ekscluster, eksnodegroup | `awless ls eksclusters`, `awless ls eksnodegroups` |
| **dynamodb** | dynamodbtable | `awless ls dynamodbtables` |
| **secretsmanager** | secret, key | `awless ls secrets`, `awless ls keys` |
| **apigateway** | apigateway, apigatewayroute, apigatewaystage | `awless ls apigateways`, `awless ls apigatewayroutes` |
| **ssm** | ssmparameter | `awless ls ssmparameters` |
| **efs** | filesystem, mounttarget | `awless ls filesystems`, `awless ls mounttargets` |
| **cloudtrail** | trail | `awless ls trails` |
| **cloudwatchlogs** | loggroup | `awless ls loggroups` |

---

## Output & Filtering

### Format options

```bash
awless list instances --format table   # default, Markdown-compatible
awless list instances --format csv
awless list instances --format tsv
awless list instances --format json
```

### Filtering by properties

Filters use exact matching (case insensitive):

```bash
awless list instances --filter state=running
awless list instances --filter state=running,type=t2.micro    # AND logic
awless list instances --filter state=running --filter type=t2.micro  # same as above
awless list volumes --filter state=in-use --filter type=gp2
awless list loggroups --filter name=/aws/lambda
```

### Filtering by tags

Tags are case sensitive:

```bash
awless list instances --tag Env=Production
awless list instances --tag Env=Production,Dept=Marketing
awless list vpcs --tag-key Dept --tag-key Internal
awless list volumes --tag-value Purchased
```

### Selecting columns

```bash
awless list instances --columns id,name,state,type,publicip
awless list users --columns name,created
```

### Sorting

```bash
awless list instances --sort uptime
awless list instances --sort launched --reverse
awless list loggroups --sort created --reverse
```

### IDs only

Useful for scripting:

```bash
awless list instances --ids --filter state=running
```

---

## Templates

Templates are files (`.aws` extension) that describe multi-step infrastructure operations. They use a simple `VERB ENTITY param=value` syntax.

### Syntax

```
# Comments start with #
# Variables capture resource IDs
myvar = create <resource> param1=value1 param2=value2

# Reference variables with $
create <other-resource> ref=$myvar

# Supported verbs: create, delete, update, attach, detach, start, stop, check, copy
```

### Example: VPC with public subnet

```
# Create networking
myvpc = create vpc cidr=10.0.0.0/16 name=prod-vpc
mysubnet = create subnet cidr=10.0.1.0/24 vpc=$myvpc name=public-subnet
igw = create internetgateway
attach internetgateway id=$igw vpc=$myvpc

# Create route table with internet access
rt = create routetable vpc=$myvpc
create route table=$rt cidr=0.0.0.0/0 gateway=$igw
attach routetable id=$rt subnet=$mysubnet

# Create security group
sg = create securitygroup vpc=$myvpc description="Web traffic" name=web-sg
update securitygroup id=$sg inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=80
update securitygroup id=$sg inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=443

# Launch instance
create instance subnet=$mysubnet type=t2.micro keypair=my-key securitygroup=$sg name=web-1
```

### Running templates

```bash
awless run my-infra.aws
awless run my-infra.aws -e         # extra verbose
```

Templates are logged and can be reverted:

```bash
awless log
awless revert <id>
```

---

## Read-Only Resources

The following newer services are currently **read-only** (list/show only, no create/delete):

- **CloudWatch Logs**: `loggroup` — list and show log groups
- **CloudTrail**: `trail` — list and show trails
- **EFS**: `filesystem`, `mounttarget` — list and show
- **SSM**: `ssmparameter` — list and show
- **API Gateway v2**: `apigateway`, `apigatewayroute`, `apigatewaystage` — list and show
- **EKS**: `ekscluster`, `eksnodegroup` — list and show
- **DynamoDB**: `dynamodbtable` — list and show
- **Secrets Manager**: `secret` — list and show
- **KMS**: `key` — list and show

To manage these resources with CRUD operations, use the AWS CLI or Console directly.

---

## Global Flags

These flags work with most commands:

| Flag | Description |
|------|-------------|
| `--local` | Use locally synced data (offline) |
| `-r`, `--region` | Override AWS region |
| `-p`, `--profile` | Override AWS profile |
| `--force` | Skip confirmation prompts |
| `-v` | Verbose output |
| `-e` | Extra verbose output |
| `--no-sync` | Don't sync after command execution |
