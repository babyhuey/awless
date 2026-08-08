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
	"attach.alarm":         {},
	"attach.containertask": {},
	"attach.elasticip": {
		"awless attach elasticip id=eipalloc-1c517b26 instance=@redis",
	},
	"attach.instance": {},
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
	"create.bucket": {
		"awless create bucket name=my-bucket-name acl=public-read",
	},
	"create.containercluster": {
		"awless create containercluster name=mycluster",
	},
	"create.classicloadbalancer": {
		"create classicloadbalancer name=my-loadb subnets=[sub-123,sub-456] listeners=HTTPS:443:HTTP:80 securitygroups=sg-54321",
		"create classicloadbalancer healthcheck-path=/health/ping listeners=TCP:80:TCP:8080 tags=Env:Test,Created:Awless",
		"create classicloadbalancer listeners=[TCP:5000:TCP:5000,HTTPS:443:HTTP:80]",
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
	"create.function": {},
	"create.group": {
		"awless create name=admins",
	},
	"create.image": {
		"awless create image instance=@my-instance-name name=redis-image description='redis prod image'",
		"awless create image instance=i-0ee436a45561c04df name=redis-image reboot=true",
		"awless create image instance=@redis-prod name=redis-prod-image",
	},
	"create.instance": {
		"awless create image=ami-123456 # Start to create instance from specific image",
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
	"create.instanceprofile":     {},
	"create.internetgateway":     {},
	"create.keypair":             {},
	"create.launchconfiguration": {},
	"create.listener":            {},
	"create.loadbalancer":        {},
	"create.loginprofile":        {},
	"create.natgateway":          {},
	"create.policy":              {},
	"create.queue":               {},
	"create.record":              {},
	"create.repository":          {},
	"create.role":                {},
	"create.route":               {},
	"create.routetable":          {},
	"create.s3object":            {},
	"create.scalinggroup":        {},
	"create.scalingpolicy":       {},
	"create.securitygroup": {
		"awless create securitygroup vpc=@myvpc name=ssh-only description=ssh-access",
		"(... see more params at `awless update securitygroup -h`)",
	},
	"create.snapshot":            {},
	"create.stack":               {},
	"create.subnet":              {},
	"create.subscription":        {},
	"create.tag":                 {},
	"create.targetgroup":         {},
	"create.topic":               {},
	"create.user":                {},
	"create.volume":              {},
	"create.vpc":                 {},
	"create.zone":                {},
	"delete.accesskey":           {},
	"delete.alarm":               {},
	"delete.appscalingpolicy":    {},
	"delete.appscalingtarget":    {},
	"delete.bucket":              {},
	"delete.containercluster":    {},
	"delete.containertask":       {},
	"delete.database":            {},
	"delete.dbsubnetgroup":       {},
	"delete.distribution":        {},
	"delete.elasticip":           {},
	"delete.function":            {},
	"delete.group":               {},
	"delete.image":               {},
	"delete.instance":            {},
	"delete.instanceprofile":     {},
	"delete.internetgateway":     {},
	"delete.keypair":             {},
	"delete.launchconfiguration": {},
	"delete.listener":            {},
	"delete.loadbalancer":        {},
	"delete.loginprofile":        {},
	"delete.natgateway":          {},
	"delete.policy":              {},
	"delete.queue":               {},
	"delete.record":              {},
	"delete.repository":          {},
	"delete.role":                {},
	"delete.route":               {},
	"delete.routetable":          {},
	"delete.s3object":            {},
	"delete.scalinggroup":        {},
	"delete.scalingpolicy":       {},
	"delete.securitygroup":       {},
	"delete.snapshot":            {},
	"delete.stack":               {},
	"delete.subnet":              {},
	"delete.subscription":        {},
	"delete.tag":                 {},
	"delete.targetgroup":         {},
	"delete.topic":               {},
	"delete.user": {
		"awless delete user name=john",
	},
	"delete.volume":          {},
	"delete.vpc":             {},
	"delete.zone":            {},
	"detach.alarm":           {},
	"detach.containertask":   {},
	"detach.elasticip":       {},
	"detach.instance":        {},
	"detach.instanceprofile": {},
	"detach.internetgateway": {},
	"detach.policy":          {},
	"detach.role":            {},
	"detach.routetable":      {},
	"detach.securitygroup":   {},
	"detach.user":            {},
	"detach.volume":          {},
	"import.image":           {},
	"start.alarm":            {},
	"start.containertask":    {},
	"start.instance":         {},
	"stop.alarm":             {},
	"stop.containertask":     {},
	"stop.instance":          {},
	"update.bucket":          {},
	"update.classicloadbalancer": {
		"awless update classicloadbalancer name=my-loadb health-target=HTTP:80/health health-interval=30 health-timeout=5 healthy-threshold=10 unhealthy-threshold=2",
	},
	"update.containertask": {},
	"update.distribution":  {},
	"update.instance":      {},
	"update.image": {
		"awless update image id=@my-image description=new-description",
		"awless update image id=ami-bd6bb2c5 groups=all operation=add # Make an AMI public",
		"awless update image id=ami-bd6bb2c5 groups=all operation=remove # Make an AMI private",
		"awless update image id=@my-image accounts=3456728198326 operation=add # Grants launch permission to an AWS account",
		"awless update image id=@my-image accounts=[3456728198326,546371829387] operation=remove  # Remove launch permission to multiple AWS accounts",
	},
	"update.loginprofile": {},
	"update.policy":       {},
	"update.record":       {},
	"update.s3object":     {},
	"update.scalinggroup": {},
	"update.securitygroup": {
		"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=26257",
		"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp securitygroup=sg-123457 portrange=8080",
	},
	"update.stack":       {},
	"update.subnet":      {},
	"update.targetgroup": {},
}
