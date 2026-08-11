package awsdoc

import (
	"bytes"
	"fmt"
)

func AwlessCommandDefinitionsDoc(action, entity, fallbackDef string) string {
	if v, ok := CommandDefinitionsDoc[action+"."+entity]; ok {
		return v
	}
	return fallbackDef
}

var CommandDefinitionsDoc = map[string]string{
	"copy.image":                 "Copy an EC2 image from given source region to current awless region",
	"create.classicloadbalancer": "Create a ELB Classic Loadbalancer (recommended only for EC2 Classic instances).\n\nYou should favor newer AWS load balancers. See `awless create loadbalancer -h`.",
}

func AwlessExamplesDoc(action, entity string) string {
	return exampleDoc(action + "." + entity)
}

func exampleDoc(key string) string {
	examples, ok := cliExamplesDoc[key]
	if ok {
		var buf bytes.Buffer
		for i, ex := range examples {
			fmt.Fprintf(&buf, "  %s", ex)
			if i != len(examples)-1 {
				buf.WriteByte('\n')
			}
		}
		return buf.String()
	}
	return ""
}

var cliExamplesDoc = map[string][]string{
	"attach.elasticip": {
		"awless attach elasticip id=eipalloc-1c517b26 instance=@redis",
	},
	"attach.instanceprofile": {
		"awless attach instanceprofile instance=@redis name=MyProfile replace=true",
	},
	"attach.internetgateway": {
		"awless attach internetgateway id=igw-636c0504 vpc=vpc-1aba387c",
	},
	"attach.listener": {
		"awless attach listener certificate=@www.mysite.com id=arn:aws:elasticloadbalancing:.../00683da53db92e54",
		"awless attach listener certificate=arn:aws:acm:...a7b691c218 id=arn:aws:elasticloadbalancing:.../00683da53db92e54",
	},
	"attach.policy": {
		"awless attach policy role=MyNewRole service=ec2 access=readonly",
		"awless attach policy user=jsmith service=s3 access=readonly",
	},
	"attach.role": {
		"awless attach role instanceprofile=MyProfile name=MyRole",
	},
	"attach.routetable": {
		"awless attach routetable id=rtb-306da254 subnet=@my-subnet",
	},
	"attach.securitygroup": {
		"awless attach securitygroup id=sg-0714247d instance=@redis",
	},
	"attach.user": {
		"awless attach user name=jsmith group=AdminGroup",
	},
	"attach.volume": {
		"awless attach volume id=vol-123oefwejf device=/dev/sdh instance=@redis",
	},
	"authenticate.registry": {
		"awless authenticate registry",
	},
	"check.database": {
		"awless check database id=@mydb state=available timeout=180",
	},
	"check.distribution": {
		"awless check distribution id=@mydistr state=Deployed timeout=180",
	},
	"check.instance": {
		"awless check instance id=@redis state=running timeout=180",
	},
	"check.loadbalancer": {
		"awless check loadbalancer id=@myloadb state=active timeout=180",
	},
	"check.natgateway": {
		"awless check natgateway id=@mynat state=active timeout=180",
	},
	"check.scalinggroup": {
		"awless check scalinggroup name=MyAutoScalingGroup count=3 timeout=180",
	},
	"check.securitygroup": {
		"awless check securitygroup id=@mysshsecgroup state=unused timeout=180",
	},
	"check.volume": {
		"awless check volume id=vol-12r1o3rp state=available timeout=180",
	},
	"copy.image": {
		"awless copy image name=my-ami-name source-id=ami-23or2or source-region=us-west-2",
	},
	"copy.snapshot": {
		"awless copy snapshot source-id=efwqwdr2or source-region=us-west-2",
	},
	"create.accesskey": {
		"awless create accesskey user=jsmith no-prompt=true",
	},
	"create.alarm": {
		" awless create alarm namespace=AWS/EC2 dimensions=AutoScalingGroupName:instancesScalingGroup evaluation-periods=2 metric=CPUUtilization name=scaleinAlarm operator=GreaterThanOrEqualToThreshold period=300 statistic-function=Average threshold=75",
	},
	"create.appscalingpolicy": {
		" awless create appscalingpolicy dimension=ecs:service:DesiredCount name=ScaleOutPolicy resource=service/my-ecs-cluster/my-service-deployment-name service-namespace=ecs stepscaling-adjustment-type=ChangeInCapacity stepscaling-adjustments=0::+1 type=StepScaling stepscaling-aggregation-type=Average stepscaling-cooldown=60",
	},
	"create.appscalingtarget": {
		"awless create appscalingtarget dimension=ecs:service:DesiredCount min-capacity=2 max-capacity=10 resource=service/my-ecs-cluster/my-service-deployment-nameource role=arn:aws:iam::519101889238:role/ecsAutoscaleRole service-namespace=ecs",
	},
	"create.secret": {
		"awless create secret name=db-password secret=s3cr3t",
		"awless create secret name=prod/api-key secret=abc123 description=\"API key for prod\"",
	},
	"update.secret": {
		"awless update secret id=db-password secret=rotated-value",
	},
	"delete.secret": {
		"awless delete secret id=db-password",
		"awless delete secret id=db-password recovery-window=7",
		"awless delete secret id=db-password force=true",
	},
	"create.ssmparameter": {
		"awless create ssmparameter name=/app/db/host value=db.internal",
		"awless create ssmparameter name=/app/db/password value=s3cr3t type=SecureString",
	},
	"update.ssmparameter": {
		"awless update ssmparameter name=/app/db/host value=db2.internal",
	},
	"delete.ssmparameter": {
		"awless delete ssmparameter name=/app/db/host",
	},
	"create.dynamodbtable": {
		"awless create dynamodbtable name=users partition-key=id",
		"awless create dynamodbtable name=events partition-key=id partition-type=S sort-key=ts sort-type=N",
		"awless create dynamodbtable name=sessions partition-key=id billing-mode=PROVISIONED read-capacity=10 write-capacity=5",
	},
	"delete.dynamodbtable": {
		"awless delete dynamodbtable name=users",
	},
	"create.ekscluster": {
		"awless create ekscluster name=prod role=arn:aws:iam::123456789012:role/eksClusterRole subnets=@subnet-a,@subnet-b",
	},
	"delete.ekscluster": {
		"awless delete ekscluster name=prod",
	},
	"create.eksnodegroup": {
		"awless create eksnodegroup name=workers cluster=prod role=arn:aws:iam::123456789012:role/eksNodeRole subnets=@subnet-a,@subnet-b instance-type=t3.medium",
		"awless create eksnodegroup name=workers cluster=prod role=arn:aws:iam::123456789012:role/eksNodeRole subnets=@subnet-a min-size=1 max-size=5 desired-size=3",
	},
	"delete.eksnodegroup": {
		"awless delete eksnodegroup name=workers cluster=prod",
	},
	"create.filesystem": {
		"awless create filesystem token=my-filesystem encrypted=true",
		"awless create filesystem token=shared-data performance-mode=maxIO throughput-mode=elastic",
	},
	"delete.filesystem": {
		"awless delete filesystem id=fs-0123456789abcdef",
	},
	"create.trail": {
		"awless create trail name=org-audit bucket=my-cloudtrail-bucket multiregion=true log-validation=true",
	},
	"delete.trail": {"awless delete trail name=org-audit"},
	"start.trail":  {"awless start trail name=org-audit"},
	"stop.trail":   {"awless stop trail name=org-audit"},
	"create.loggroup": {
		"awless create loggroup name=/aws/lambda/my-function",
	},
	"delete.loggroup": {"awless delete loggroup name=/aws/lambda/my-function"},
	"update.loggroup": {
		"awless update loggroup name=/aws/lambda/my-function retention=30",
	},
	"create.cachecluster": {
		"awless create cachecluster id=sessions engine=redis type=cache.t3.micro nodes=1",
		"awless create cachecluster id=sessions engine=redis type=cache.t3.micro nodes=1 subnet-group=cache-private port=6379",
		"awless create cachecluster id=sessions-replica-1 replication-group=sessions-group",
	},
	"delete.cachecluster": {
		"awless delete cachecluster id=sessions",
		"awless delete cachecluster id=sessions snapshot=sessions-final",
	},
	"update.cachecluster": {
		"awless update cachecluster id=sessions nodes=3 apply-immediately=true",
		"awless update cachecluster id=sessions type=cache.t3.small",
	},
	"create.cachesubnetgroup": {
		"awless create cachesubnetgroup name=cache-private subnets=subnet-1234,subnet-5678 description=\"Private cache subnets\"",
	},
	"delete.cachesubnetgroup": {"awless delete cachesubnetgroup name=cache-private"},
	"update.cachesubnetgroup": {
		"awless update cachesubnetgroup name=cache-private subnets=subnet-1234,subnet-5678,subnet-9abc",
	},
	"create.replicationgroup": {
		"awless create replicationgroup id=sessions-group description=\"Session cache\" engine=redis type=cache.t3.micro clusters=2",
		"awless create replicationgroup id=sessions-group description=\"Session cache\" engine=redis type=cache.t3.micro clusters=3 automatic-failover=true multi-az=true",
	},
	"delete.replicationgroup": {
		"awless delete replicationgroup id=sessions-group",
		"awless delete replicationgroup id=sessions-group retain-primary=true",
	},
	"create.eventbus": {
		"awless create eventbus name=orders description=\"Order domain events\"",
	},
	"delete.eventbus": {"awless delete eventbus name=orders"},
	"create.eventrule": {
		"awless create eventrule name=nightly-report schedule=\"cron(0 6 * * ? *)\"",
		"awless create eventrule name=on-instance-stop pattern='{\"source\":[\"aws.ec2\"]}' eventbus=orders",
	},
	"update.eventrule": {
		"awless update eventrule name=nightly-report schedule=\"rate(1 hour)\" state=DISABLED",
	},
	"delete.eventrule": {
		"awless delete eventrule name=nightly-report",
		"awless delete eventrule name=nightly-report force=true",
	},
	"start.eventrule": {"awless start eventrule name=nightly-report"},
	"stop.eventrule":  {"awless stop eventrule name=nightly-report"},
	"attach.eventtarget": {
		"awless attach eventtarget rule=nightly-report id=report-lambda arn=arn:aws:lambda:us-west-2:123456789012:function:report",
	},
	"detach.eventtarget": {
		"awless detach eventtarget rule=nightly-report id=report-lambda",
	},
	"create.statemachine": {
		"awless create statemachine name=order-flow role=arn:aws:iam::123456789012:role/StepFunctions definition-file=/home/jsmith/order-flow.json",
		"awless create statemachine name=fast-flow role=arn:aws:iam::123456789012:role/StepFunctions definition-file=/home/jsmith/flow.json type=EXPRESS",
	},
	"delete.statemachine": {
		"awless delete statemachine arn=arn:aws:states:us-west-2:123456789012:stateMachine:order-flow",
	},
	"update.statemachine": {
		"awless update statemachine arn=arn:aws:states:us-west-2:123456789012:stateMachine:order-flow definition-file=/home/jsmith/order-flow-v2.json",
	},
	"start.execution": {
		"awless start execution statemachine=arn:aws:states:us-west-2:123456789012:stateMachine:order-flow",
		"awless start execution statemachine=arn:aws:states:us-west-2:123456789012:stateMachine:order-flow name=order-4711 input='{\"orderId\":4711}'",
	},
	"stop.execution": {
		"awless stop execution arn=arn:aws:states:us-west-2:123456789012:execution:order-flow:order-4711 cause=Superseded",
	},
	"create.accelerator": {"awless create accelerator name=global-api enabled=true"},
	"delete.accelerator": {"awless delete accelerator arn=arn:aws:globalaccelerator::123456789012:accelerator/abcd"},
	"update.accelerator": {
		"awless update accelerator arn=arn:aws:globalaccelerator::123456789012:accelerator/abcd enabled=false",
	},
	"create.acceleratorlistener": {
		"awless create acceleratorlistener accelerator=arn:aws:globalaccelerator::123456789012:accelerator/abcd protocol=TCP ports-file=/home/jsmith/ports.json",
	},
	"delete.acceleratorlistener": {
		"awless delete acceleratorlistener arn=arn:aws:globalaccelerator::123456789012:accelerator/abcd/listener/1234",
	},
	"create.fsxfilesystem": {
		"awless create fsxfilesystem type=LUSTRE capacity=1200 subnets=subnet-1 lustre-file=/home/jsmith/lustre.json",
		"awless create fsxfilesystem type=WINDOWS capacity=32 subnets=subnet-1,subnet-2 windows-file=/home/jsmith/windows.json",
	},
	"delete.fsxfilesystem": {"awless delete fsxfilesystem id=fs-1234abcd"},
	"create.fsxbackup": {
		"awless create fsxbackup filesystem=fs-1234abcd",
	},
	"delete.fsxbackup": {"awless delete fsxbackup id=backup-1234abcd"},
	"create.broker": {
		"awless create broker name=orders engine=RABBITMQ type=mq.m5.large mode=SINGLE_INSTANCE public=false users-file=/home/jsmith/users.json",
	},
	"delete.broker": {"awless delete broker id=b-1234abcd"},
	"create.kafkacluster": {
		"awless create kafkacluster name=events version=3.6.0 brokers=3 subnets=subnet-1,subnet-2,subnet-3 type=kafka.m5.large storage=100",
	},
	"delete.kafkacluster": {
		"awless delete kafkacluster arn=arn:aws:kafka:us-west-2:123456789012:cluster/events/abcd",
	},
	"create.userpool": {
		"awless create userpool name=customers deletion-protection=ACTIVE auto-verified=email",
	},
	"delete.userpool": {"awless delete userpool id=us-west-2_abcDEF123"},
	"create.identitypool": {
		"awless create identitypool name=customers-identities allow-unauthenticated=false",
	},
	"delete.identitypool": {"awless delete identitypool id='us-west-2:11111111-2222-3333-4444-555555555555'"},
	"create.emailidentity": {
		"awless create emailidentity name=mail.example.com",
		"awless create emailidentity name=noreply@example.com configuration-set=transactional",
	},
	"delete.emailidentity":    {"awless delete emailidentity name=mail.example.com"},
	"create.configurationset": {"awless create configurationset name=transactional"},
	"delete.configurationset": {"awless delete configurationset name=transactional"},
	"create.gluedatabase":     {"awless create gluedatabase name=analytics description=\"Analytics catalog\""},
	"delete.gluedatabase":     {"awless delete gluedatabase name=analytics"},
	"create.crawler": {
		"awless create crawler name=events-crawler role=AWSGlueServiceRole database=analytics targets-file=/home/jsmith/targets.json",
	},
	"delete.crawler": {"awless delete crawler name=events-crawler"},
	"start.crawler":  {"awless start crawler name=events-crawler"},
	"stop.crawler":   {"awless stop crawler name=events-crawler"},
	"create.job": {
		"awless create job name=etl-events role=AWSGlueServiceRole command=glueetl script=s3://scripts/etl.py glue-version=4.0 worker-type=G.1X workers=2",
	},
	"delete.job": {"awless delete job name=etl-events"},
	"start.job": {
		"awless start job name=etl-events",
		"awless start job name=etl-events workers=10 worker-type=G.2X",
	},
	"create.deployapplication": {
		"awless create deployapplication name=web-api platform=Server",
	},
	"delete.deployapplication": {"awless delete deployapplication name=web-api"},
	"create.deploymentgroup": {
		"awless create deploymentgroup name=prod application=web-api role=arn:aws:iam::123456789012:role/CodeDeploy config=CodeDeployDefault.OneAtATime",
		"awless create deploymentgroup name=prod application=web-api role=arn:aws:iam::123456789012:role/CodeDeploy scalinggroups=web-asg",
	},
	"delete.deploymentgroup": {"awless delete deploymentgroup name=prod application=web-api"},
	"create.deployment": {
		"awless create deployment application=web-api group=prod revision-file=/home/jsmith/revision.json",
	},
	"stop.deployment": {
		"awless stop deployment id=d-ABCDEF123 rollback=true",
	},
	"create.transitgateway": {
		"awless create transitgateway description=\"Shared services hub\"",
		"awless create transitgateway description=hub auto-accept=disable default-association=disable",
	},
	"delete.transitgateway": {"awless delete transitgateway id=tgw-1234"},
	"create.transitgatewayattachment": {
		"awless create transitgatewayattachment transitgateway=tgw-1234 vpc=vpc-5678 subnets=subnet-1,subnet-2",
	},
	"delete.transitgatewayattachment": {"awless delete transitgatewayattachment id=tgw-attach-1234"},
	"create.transitgatewayroutetable": {
		"awless create transitgatewayroutetable transitgateway=tgw-1234",
	},
	"delete.transitgatewayroutetable": {"awless delete transitgatewayroutetable id=tgw-rtb-1234"},
	"create.vpcendpoint": {
		"awless create vpcendpoint vpc=vpc-1234 service=com.amazonaws.us-west-2.s3 type=Gateway routetables=rtb-1234",
		"awless create vpcendpoint vpc=vpc-1234 service=com.amazonaws.us-west-2.secretsmanager type=Interface subnets=subnet-1,subnet-2 securitygroups=sg-1234 private-dns=true",
	},
	"delete.vpcendpoint": {"awless delete vpcendpoint id=vpce-1234"},
	"create.pipeline": {
		"awless create pipeline definition-file=/home/jsmith/build-and-deploy.json",
	},
	"create.webacl": {
		"awless create webacl name=public-api default-action-file=/home/jsmith/allow.json visibility-file=/home/jsmith/visibility.json rules-file=/home/jsmith/rules.json",
	},
	"delete.webacl": {
		"awless delete webacl name=public-api",
		"awless delete webacl name=edge-acl scope=CLOUDFRONT",
	},
	"create.rulegroup": {
		"awless create rulegroup name=rate-limits capacity=100 visibility-file=/home/jsmith/visibility.json rules-file=/home/jsmith/rules.json",
	},
	"delete.rulegroup": {"awless delete rulegroup name=rate-limits"},
	"create.ipset": {
		"awless create ipset name=blocklist addresses=203.0.113.0/24,198.51.100.7/32",
		"awless create ipset name=cdn-blocklist addresses=203.0.113.0/24 scope=CLOUDFRONT description=\"Blocked at the edge\"",
	},
	"delete.ipset": {
		"awless delete ipset name=blocklist",
		"awless delete ipset name=cdn-blocklist scope=CLOUDFRONT",
	},
	"update.ipset": {
		"awless update ipset name=blocklist addresses=203.0.113.0/24,198.51.100.7/32,192.0.2.0/24",
	},
	"create.configrule": {
		"awless create configrule name=s3-versioning source=S3_BUCKET_VERSIONING_ENABLED",
		"awless create configrule name=ebs-encrypted source=ENCRYPTED_VOLUMES resource-types=AWS::EC2::Volume description=\"Volumes must be encrypted\"",
	},
	"update.configrule": {
		"awless update configrule name=s3-versioning source=S3_BUCKET_VERSIONING_ENABLED frequency=TwentyFour_Hours",
	},
	"delete.configrule": {"awless delete configrule name=s3-versioning"},
	"create.stream": {
		"awless create stream name=clickstream mode=ON_DEMAND",
		"awless create stream name=clickstream mode=PROVISIONED shards=4",
	},
	"delete.stream": {
		"awless delete stream name=clickstream",
		"awless delete stream name=clickstream force=true",
	},
	"update.stream": {
		"awless update stream name=clickstream shards=8",
	},
	"create.redshiftcluster": {
		"awless create redshiftcluster id=analytics username=admin type=ra3.xlplus manage-password=true cluster-type=single-node",
		"awless create redshiftcluster id=analytics username=admin type=ra3.xlplus manage-password=true nodes=4 encrypted=true subnet-group=warehouse-private",
	},
	"delete.redshiftcluster": {
		"awless delete redshiftcluster id=analytics snapshot=analytics-final",
		"awless delete redshiftcluster id=analytics skip-snapshot=true",
	},
	"update.redshiftcluster": {
		"awless update redshiftcluster id=analytics nodes=8",
	},
	"create.redshiftsubnetgroup": {
		"awless create redshiftsubnetgroup name=warehouse-private description=\"Private warehouse subnets\" subnets=subnet-1234,subnet-5678",
	},
	"delete.redshiftsubnetgroup": {"awless delete redshiftsubnetgroup name=warehouse-private"},
	"delete.pipeline":            {"awless delete pipeline name=build-and-deploy"},
	"create.buildproject": {
		"awless create buildproject name=api-build role=arn:aws:iam::123456789012:role/CodeBuild source-type=GITHUB source-location=https://github.com/acme/api env-type=LINUX_CONTAINER image=aws/codebuild/standard:7.0 compute-type=BUILD_GENERAL1_SMALL artifact-type=NO_ARTIFACTS",
	},
	"update.buildproject": {
		"awless update buildproject name=api-build compute-type=BUILD_GENERAL1_MEDIUM timeout=30",
	},
	"delete.buildproject": {"awless delete buildproject name=api-build"},
	"start.buildproject": {
		"awless start buildproject name=api-build",
		"awless start buildproject name=api-build source-version=release-2.0",
	},
	"stop.buildproject": {
		"awless stop buildproject build='api-build:11111111-2222-3333-4444-555555555555'",
	},
	"create.application": {"awless create application name=web-api description=\"Public API\""},
	"delete.application": {
		"awless delete application name=web-api",
		"awless delete application name=web-api force=true",
	},
	"update.application": {"awless update application name=web-api description=\"Public API v2\""},
	"create.environment": {
		"awless create environment name=web-api-prod application=web-api solution-stack=\"64bit Amazon Linux 2023 v4.0.0 running Go 1\"",
		"awless create environment name=web-api-staging application=web-api config-template=saved-config version=build-42",
	},
	"delete.environment": {
		"awless delete environment name=web-api-prod",
		"awless delete environment name=web-api-prod force=true",
	},
	"update.environment": {"awless update environment name=web-api-prod version=build-43"},
	"start.pipeline":     {"awless start pipeline name=build-and-deploy"},
	"stop.pipeline": {
		"awless stop pipeline name=build-and-deploy execution='12345678-1234-1234-1234-123456789012' reason=\"Superseded by a newer commit\"",
	},
	"create.apigateway": {
		"awless create apigateway name=my-api protocol=HTTP",
		"awless create apigateway name=my-api protocol=HTTP target=arn:aws:lambda:us-west-2:123456789012:function:handler",
	},
	"delete.apigateway": {"awless delete apigateway id=abc123"},
	"create.apigatewayroute": {
		"awless create apigatewayroute api=abc123 route-key=\"GET /items\" target=integrations/xyz789",
	},
	"delete.apigatewayroute": {"awless delete apigatewayroute api=abc123 id=xyz789"},
	"create.apigatewaystage": {
		"awless create apigatewaystage api=abc123 name=prod autodeploy=true",
	},
	"delete.apigatewaystage": {"awless delete apigatewaystage api=abc123 name=prod"},
	"attach.alarm": {
		"awless attach alarm action-arn=arn:aws:sns:us-west-2:123456789012:alerts name=my-alarm",
	},
	"attach.classicloadbalancer": {
		"awless attach classicloadbalancer name=my-classicloadbalancer instance=i-0123456789abcdef0",
	},
	"attach.containertask": {
		"awless attach containertask container-name=web image=nginx:latest memory-hard-limit=512 name=my-containertask",
	},
	"attach.instance": {
		"awless attach instance id=i-0123456789abcdef0 targetgroup=arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd1234",
	},
	"attach.mfadevice": {
		"awless attach mfadevice id=arn:aws:iam::123456789012:mfa/jsmith mfa-code-1=123456 mfa-code-2=654321 user=jsmith",
	},
	"attach.networkinterface": {
		"awless attach networkinterface device-index=1 id=eni-0123456789abcdef0 instance=i-0123456789abcdef0",
	},
	"check.certificate": {
		"awless check certificate arn=arn:aws:iam::123456789012:policy/my-policy state=available timeout=30",
	},
	"check.networkinterface": {
		"awless check networkinterface id=eni-0123456789abcdef0 state=available timeout=30",
	},
	"create.certificate": {
		"awless create certificate domains=example.com",
	},
	"create.function": {
		"awless create function handler=index.handler name=my-function role=arn:aws:iam::123456789012:role/my-role runtime=python3.12",
	},
	"create.instanceprofile": {
		"awless create instanceprofile name=my-instanceprofile",
	},
	"create.internetgateway": {
		"awless create internetgateway",
	},
	"create.keypair": {
		"awless create keypair name=my-keypair",
	},
	"create.launchconfiguration": {
		"awless create launchconfiguration distro=amazonlinux name=my-launchconfiguration type=t3.micro",
	},
	"create.listener": {
		"awless create listener actiontype=forward loadbalancer=my-loadbalancer port=80 protocol=HTTP targetgroup=arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd1234",
	},
	"create.loadbalancer": {
		"awless create loadbalancer name=my-loadbalancer subnets=subnet-0123456789abcdef0,subnet-0fedcba9876543210",
	},
	"create.loginprofile": {
		"awless create loginprofile password='S3cure!Pass' username=jsmith",
	},
	"create.mfadevice": {
		"awless create mfadevice name=my-mfadevice",
	},
	"create.natgateway": {
		"awless create natgateway elasticip-id=eipalloc-0123456789abcdef0 subnet=subnet-0123456789abcdef0",
	},
	"create.networkinterface": {
		"awless create networkinterface subnet=subnet-0123456789abcdef0",
	},
	"create.policy": {
		"awless create policy action=s3:GetObject effect=Allow name=my-policy resource=arn:aws:s3:::my-bucket/*",
	},
	"create.queue": {
		"awless create queue name=my-queue",
	},
	"create.record": {
		"awless create record name=my-record ttl=300 type=A values=10.0.0.1,10.0.0.2 zone=/hostedzone/Z3P5QSUBK4POTI",
	},
	"create.repository": {
		"awless create repository name=my-repository",
	},
	"create.role": {
		"awless create role name=my-role",
	},
	"create.route": {
		"awless create route cidr=10.0.0.0/24 gateway=igw-0123456789abcdef0 table=rtb-0123456789abcdef0",
	},
	"create.routetable": {
		"awless create routetable vpc=vpc-0123456789abcdef0",
	},
	"create.s3object": {
		"awless create s3object bucket=my-bucket file=/tmp/report.pdf",
	},
	"create.scalinggroup": {
		"awless create scalinggroup launchconfiguration=my-launchconfig max-size=4 min-size=1 name=my-scalinggroup subnets=subnet-0123456789abcdef0,subnet-0fedcba9876543210",
	},
	"create.scalingpolicy": {
		"awless create scalingpolicy adjustment-scaling=1 adjustment-type=ChangeInCapacity name=my-scalingpolicy scalinggroup=my-scalinggroup",
	},
	"create.snapshot": {
		"awless create snapshot volume=vol-0123456789abcdef0",
	},
	"create.stack": {
		"awless create stack name=my-stack template-file=/tmp/template.yml",
	},
	"create.subnet": {
		"awless create subnet cidr=10.0.0.0/24 vpc=vpc-0123456789abcdef0",
	},
	"create.subscription": {
		"awless create subscription endpoint=ops@example.com protocol=HTTP topic=arn:aws:sns:us-west-2:123456789012:alerts",
	},
	"create.tag": {
		"awless create tag key=Env resource=arn:aws:s3:::my-bucket/* value=production",
	},
	"create.targetgroup": {
		"awless create targetgroup name=my-targetgroup port=80 protocol=HTTP vpc=vpc-0123456789abcdef0",
	},
	"create.topic": {
		"awless create topic name=my-topic",
	},
	"create.user": {
		"awless create user name=my-user",
	},
	"create.volume": {
		"awless create volume availabilityzone=us-west-2a size=20",
	},
	"create.vpc": {
		"awless create vpc cidr=10.0.0.0/24",
	},
	"create.zone": {
		"awless create zone callerreference=my-zone-2024-01-01 name=my-zone",
	},
	"delete.accesskey": {
		"awless delete accesskey id=AKIAIOSFODNN7EXAMPLE user=my-accesskey",
	},
	"delete.alarm": {
		"awless delete alarm name=my-alarm",
	},
	"delete.appscalingpolicy": {
		"awless delete appscalingpolicy dimension=ecs:service:DesiredCount name=my-appscalingpolicy resource=arn:aws:s3:::my-bucket/* service-namespace=ecs",
	},
	"delete.appscalingtarget": {
		"awless delete appscalingtarget dimension=ecs:service:DesiredCount resource=arn:aws:s3:::my-bucket/* service-namespace=ecs",
	},
	"delete.bucket": {
		"awless delete bucket name=my-bucket",
	},
	"delete.certificate": {
		"awless delete certificate arn=arn:aws:iam::123456789012:policy/my-policy",
	},
	"delete.classicloadbalancer": {
		"awless delete classicloadbalancer name=my-classicloadbalancer",
	},
	"delete.containercluster": {
		"awless delete containercluster id=my-cluster",
	},
	"delete.containertask": {
		"awless delete containertask name=my-containertask",
	},
	"delete.database": {
		"awless delete database id=my-database",
	},
	"delete.dbsubnetgroup": {
		"awless delete dbsubnetgroup name=my-dbsubnetgroup",
	},
	"delete.distribution": {
		"awless delete distribution id=E1PA6795SAMPLE",
	},
	"delete.elasticip": {
		"awless delete elasticip id=eipalloc-0123456789abcdef0",
	},
	"delete.function": {
		"awless delete function id=my-function",
	},
	"delete.group": {
		"awless delete group name=my-group",
	},
	"delete.image": {
		"awless delete image id=ami-0123456789abcdef0",
	},
	"delete.instance": {
		"awless delete instance ids=i-0123456789abcdef0,i-0fedcba9876543210",
	},
	"delete.instanceprofile": {
		"awless delete instanceprofile name=my-instanceprofile",
	},
	"delete.internetgateway": {
		"awless delete internetgateway id=igw-0123456789abcdef0",
	},
	"delete.keypair": {
		"awless delete keypair name=my-keypair",
	},
	"delete.launchconfiguration": {
		"awless delete launchconfiguration name=my-launchconfiguration",
	},
	"delete.listener": {
		"awless delete listener id=arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-lb/abcd1234/ef567890",
	},
	"delete.loadbalancer": {
		"awless delete loadbalancer id=arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-lb/abcd1234",
	},
	"delete.loginprofile": {
		"awless delete loginprofile username=jsmith",
	},
	"delete.mfadevice": {
		"awless delete mfadevice id=arn:aws:iam::123456789012:mfa/jsmith",
	},
	"delete.natgateway": {
		"awless delete natgateway id=nat-0123456789abcdef0",
	},
	"delete.networkinterface": {
		"awless delete networkinterface id=eni-0123456789abcdef0",
	},
	"delete.policy": {
		"awless delete policy arn=arn:aws:iam::123456789012:policy/my-policy",
	},
	"delete.queue": {
		"awless delete queue url=https://example.com/image.vmdk",
	},
	"delete.record": {
		"awless delete record name=my-record ttl=300 type=A values=10.0.0.1,10.0.0.2 zone=/hostedzone/Z3P5QSUBK4POTI",
	},
	"delete.repository": {
		"awless delete repository name=my-repository",
	},
	"delete.role": {
		"awless delete role name=my-role",
	},
	"delete.route": {
		"awless delete route cidr=10.0.0.0/24 table=rtb-0123456789abcdef0",
	},
	"delete.routetable": {
		"awless delete routetable id=rtb-0123456789abcdef0",
	},
	"delete.s3object": {
		"awless delete s3object bucket=my-bucket name=my-s3object",
	},
	"delete.scalinggroup": {
		"awless delete scalinggroup name=my-scalinggroup",
	},
	"delete.scalingpolicy": {
		"awless delete scalingpolicy id=my-scalingpolicy",
	},
	"delete.securitygroup": {
		"awless delete securitygroup id=sg-0123456789abcdef0",
	},
	"delete.snapshot": {
		"awless delete snapshot id=snap-0123456789abcdef0",
	},
	"delete.stack": {
		"awless delete stack name=my-stack",
	},
	"delete.subnet": {
		"awless delete subnet id=subnet-0123456789abcdef0",
	},
	"delete.subscription": {
		"awless delete subscription id=arn:aws:sns:us-west-2:123456789012:alerts:abcd1234",
	},
	"delete.tag": {
		"awless delete tag key=Env resource=arn:aws:s3:::my-bucket/*",
	},
	"delete.targetgroup": {
		"awless delete targetgroup id=arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd1234",
	},
	"delete.topic": {
		"awless delete topic id=arn:aws:sns:us-west-2:123456789012:alerts",
	},
	"delete.volume": {
		"awless delete volume id=vol-0123456789abcdef0",
	},
	"delete.vpc": {
		"awless delete vpc id=vpc-0123456789abcdef0",
	},
	"delete.zone": {
		"awless delete zone id=/hostedzone/Z3P5QSUBK4POTI",
	},
	"detach.alarm": {
		"awless detach alarm action-arn=arn:aws:sns:us-west-2:123456789012:alerts name=my-alarm",
	},
	"detach.classicloadbalancer": {
		"awless detach classicloadbalancer name=my-classicloadbalancer instance=i-0123456789abcdef0",
	},
	"detach.containertask": {
		"awless detach containertask container-name=web name=my-containertask",
	},
	"detach.elasticip": {
		"awless detach elasticip association=eipassoc-0123456789abcdef0",
	},
	"detach.instance": {
		"awless detach instance id=i-0123456789abcdef0 targetgroup=arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd1234",
	},
	"detach.instanceprofile": {
		"awless detach instanceprofile instance=i-0123456789abcdef0 name=my-instanceprofile",
	},
	"detach.internetgateway": {
		"awless detach internetgateway id=igw-0123456789abcdef0 vpc=vpc-0123456789abcdef0",
	},
	"detach.mfadevice": {
		"awless detach mfadevice id=arn:aws:iam::123456789012:mfa/jsmith user=jsmith",
	},
	"detach.networkinterface": {
		"awless detach networkinterface instance=i-0123456789abcdef0 id=eni-0123456789abcdef0",
	},
	"detach.policy": {
		"awless detach policy user=jsmith arn=arn:aws:iam::123456789012:policy/my-policy",
	},
	"detach.role": {
		"awless detach role instanceprofile=my-instance-profile name=my-role",
	},
	"detach.routetable": {
		"awless detach routetable association=eipassoc-0123456789abcdef0",
	},
	"detach.securitygroup": {
		"awless detach securitygroup id=sg-0123456789abcdef0 instance=i-0123456789abcdef0",
	},
	"detach.user": {
		"awless detach user group=developers name=my-user",
	},
	"detach.volume": {
		"awless detach volume device=/dev/sdh id=vol-0123456789abcdef0 instance=i-0123456789abcdef0",
	},
	"import.image": {
		"awless import image snapshot=snap-0123456789abcdef0",
	},
	"restart.database": {
		"awless restart database id=my-database",
	},
	"restart.instance": {
		"awless restart instance ids=i-0123456789abcdef0,i-0fedcba9876543210",
	},
	"start.alarm": {
		"awless start alarm names=my-alarm",
	},
	"start.containertask": {
		"awless start containertask cluster=my-cluster desired-count=2 name=my-containertask type=service",
	},
	"start.database": {
		"awless start database id=my-database",
	},
	"start.instance": {
		"awless start instance ids=i-0123456789abcdef0,i-0fedcba9876543210",
	},
	"stop.alarm": {
		"awless stop alarm names=my-alarm",
	},
	"stop.containertask": {
		"awless stop containertask cluster=my-cluster type=service",
	},
	"stop.database": {
		"awless stop database id=my-database",
	},
	"stop.instance": {
		"awless stop instance ids=i-0123456789abcdef0,i-0fedcba9876543210",
	},
	"update.bucket": {
		"awless update bucket name=my-bucket",
	},
	"update.containertask": {
		"awless update containertask cluster=my-cluster deployment-name=web-service",
	},
	"update.distribution": {
		"awless update distribution id=E1PA6795SAMPLE",
	},
	"update.instance": {
		"awless update instance id=i-0123456789abcdef0",
	},
	"update.loginprofile": {
		"awless update loginprofile password='S3cure!Pass' username=jsmith",
	},
	"update.policy": {
		"awless update policy action=s3:GetObject arn=arn:aws:iam::123456789012:policy/my-policy effect=Allow resource=arn:aws:s3:::my-bucket/*",
	},
	"update.record": {
		"awless update record name=my-record ttl=300 type=A values=10.0.0.1,10.0.0.2 zone=/hostedzone/Z3P5QSUBK4POTI",
	},
	"update.s3object": {
		"awless update s3object acl=public-read bucket=my-bucket name=my-s3object",
	},
	"update.scalinggroup": {
		"awless update scalinggroup name=my-scalinggroup",
	},
	"update.stack": {
		"awless update stack name=my-stack",
	},
	"update.subnet": {
		"awless update subnet id=subnet-0123456789abcdef0",
	},
	"update.targetgroup": {
		"awless update targetgroup id=arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd1234",
	},
	"create.bucket": {
		"awless create bucket name=my-bucket-name acl=public-read",
	},
	"create.containercluster": {
		"awless create containercluster name=mycluster",
	},
	"create.classicloadbalancer": {
		"awless create classicloadbalancer name=my-loadb subnets=[sub-123,sub-456] listeners=HTTPS:443:HTTP:80 securitygroups=sg-54321",
		"awless create classicloadbalancer name=my-loadb healthcheck-path=/health/ping listeners=TCP:80:TCP:8080 tags=Env:Test,Created:Awless",
		"awless create classicloadbalancer name=my-loadb listeners=[TCP:5000:TCP:5000,HTTPS:443:HTTP:80]",
	},
	"create.database": {
		"awless create database engine=postgres id=mystartup-prod-db subnetgroup=@my-dbsubnetgroup password=notsafe dbname=mydb size=5 type=db.t2.small username=admin vpcsecuritygroups=@postgres_sg",
	},
	"create.dbsubnetgroup": {
		"awless create dbsubnetgroup name=mydbsubnetgroup description=\"subnets for peps db\" subnets=[@my-firstsubnet, @my-secondsubnet]",
	},
	"create.distribution": {
		"awless create distribution origin-domain=mybucket.s3.amazonaws.com",
	},
	"create.elasticip": {
		"awless create elasticip domain=vpc",
	},
	"create.group": {
		"awless create group name=admins",
	},
	"create.image": {
		"awless create image instance=@my-instance-name name=redis-image description='redis prod image'",
		"awless create image instance=i-0ee436a45561c04df name=redis-image reboot=true",
		"awless create image instance=@redis-prod name=redis-prod-image",
	},
	"create.instance": {
		"awless create instance image=ami-123456 # Start to create instance from specific image",
		"awless create instance keypair=jsmith type=t2.micro subnet=@my-subnet",
		"awless create instance image=ami-123456 keypair=jsmith",
		"awless create instance name=redis type=t2.nano keypair=jsmith userdata=/home/jsmith/data.sh",
		"", // create empty line for clarity
		"awless create instance distro=redhat type=t2.micro",
		"awless create instance distro=coreos name=redis-prod",
		"awless create instance distro=redhat::7.2 type=t2.micro",
		"awless create instance distro=canonical:ubuntu role=MyInfraReadOnlyRole",
		"awless create instance distro=debian:debian:jessie lock=true",
		"awless create instance distro=amazonlinux securitygroup=@my-ssh-secgroup",
		"awless create instance distro=amazonlinux:::::instance-store",
		"awless create instance distro=amazonlinux:amzn2",
	},
	"create.securitygroup": {
		"awless create securitygroup vpc=@myvpc name=ssh-only description=ssh-access",
	},
	"delete.user": {
		"awless delete user name=john",
	},
	"update.classicloadbalancer": {
		"awless update classicloadbalancer name=my-loadb health-target=HTTP:80/health health-interval=30 health-timeout=5 healthy-threshold=10 unhealthy-threshold=2",
	},
	"update.image": {
		"awless update image id=@my-image description=new-description",
		"awless update image id=ami-bd6bb2c5 groups=all operation=add # Make an AMI public",
		"awless update image id=ami-bd6bb2c5 groups=all operation=remove # Make an AMI private",
		"awless update image id=@my-image accounts=3456728198326 operation=add # Grants launch permission to an AWS account",
		"awless update image id=@my-image accounts=[3456728198326,546371829387] operation=remove  # Remove launch permission to multiple AWS accounts",
	},
	"update.securitygroup": {
		"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=26257",
		"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp securitygroup=sg-123457 portrange=8080",
	},
}
