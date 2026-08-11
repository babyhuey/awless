/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsspec

import (
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// Integer params are declared as *int64 with awsType "awsint64" even though the AWS
// fields are *int32. The reflective setter converts between widths, and this is the
// established convention here — see CreateDatabase's port.

type CreateCachecluster struct {
	_                 string `action:"create" entity:"cachecluster" awsAPI:"elasticache" awsCall:"CreateCacheCluster" awsInput:"elasticache.CreateCacheClusterInput" awsOutput:"elasticache.CreateCacheClusterOutput"`
	logger            *logger.Logger
	graph             cloud.GraphAPI
	api               *elasticache.Client
	ID                *string   `awsName:"CacheClusterId" awsType:"awsstr" templateName:"id"`
	Engine            *string   `awsName:"Engine" awsType:"awsstr" templateName:"engine"`
	EngineVersion     *string   `awsName:"EngineVersion" awsType:"awsstr" templateName:"engine-version"`
	Type              *string   `awsName:"CacheNodeType" awsType:"awsstr" templateName:"type"`
	Nodes             *int64    `awsName:"NumCacheNodes" awsType:"awsint64" templateName:"nodes"`
	SubnetGroup       *string   `awsName:"CacheSubnetGroupName" awsType:"awsstr" templateName:"subnet-group"`
	Securitygroups    []*string `awsName:"SecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	Port              *int64    `awsName:"Port" awsType:"awsint64" templateName:"port"`
	Zone              *string   `awsName:"PreferredAvailabilityZone" awsType:"awsstr" templateName:"zone"`
	ReplicationGroup  *string   `awsName:"ReplicationGroupId" awsType:"awsstr" templateName:"replication-group"`
	ParameterGroup    *string   `awsName:"CacheParameterGroupName" awsType:"awsstr" templateName:"parameter-group"`
	SnapshotRetention *int64    `awsName:"SnapshotRetentionLimit" awsType:"awsint64" templateName:"snapshot-retention"`
}

// A cache cluster is created either standalone, which needs an engine and a node type,
// or as a read replica joining an existing replication group, which inherits both. Those
// two forms are mutually exclusive; AWS rejects engine and node-type overrides on a
// replication-group member.
//
// The shared params sit above the OnlyOneOf rather than inside each branch. The param
// collector walks every branch, so anything repeated is listed twice in `-h`. The cost is
// that a param only meaningful to one form — `nodes` on a replica, say — is accepted here
// and rejected by AWS instead.
func (cmd *CreateCachecluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.OnlyOneOf(
			params.AllOf(params.Key("engine"), params.Key("type")),
			params.Key("replication-group"),
		),
		params.Opt(
			params.Suggested("nodes"),
			"engine-version", "subnet-group", "securitygroups", "port", "zone",
			"parameter-group", "snapshot-retention",
		),
	))
}

func (cmd *CreateCachecluster) ExtractResult(i any) string {
	out, ok := i.(*elasticache.CreateCacheClusterOutput)
	if !ok || out.CacheCluster == nil {
		// An empty response is not fatal here: the caller supplied the id, so the
		// graph still has something to key on. Six commands panicked on exactly this
		// shape before, so the guard is deliberate.
		return StringValue(cmd.ID)
	}
	return awssdk.ToString(out.CacheCluster.CacheClusterId)
}

type DeleteCachecluster struct {
	_        string `action:"delete" entity:"cachecluster" awsAPI:"elasticache" awsCall:"DeleteCacheCluster" awsInput:"elasticache.DeleteCacheClusterInput" awsOutput:"elasticache.DeleteCacheClusterOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *elasticache.Client
	ID       *string `awsName:"CacheClusterId" awsType:"awsstr" templateName:"id"`
	Snapshot *string `awsName:"FinalSnapshotIdentifier" awsType:"awsstr" templateName:"snapshot"`
}

func (cmd *DeleteCachecluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.Opt("snapshot"),
	))
}

type UpdateCachecluster struct {
	_                 string `action:"update" entity:"cachecluster" awsAPI:"elasticache" awsCall:"ModifyCacheCluster" awsInput:"elasticache.ModifyCacheClusterInput" awsOutput:"elasticache.ModifyCacheClusterOutput"`
	logger            *logger.Logger
	graph             cloud.GraphAPI
	api               *elasticache.Client
	ID                *string   `awsName:"CacheClusterId" awsType:"awsstr" templateName:"id"`
	Nodes             *int64    `awsName:"NumCacheNodes" awsType:"awsint64" templateName:"nodes"`
	Type              *string   `awsName:"CacheNodeType" awsType:"awsstr" templateName:"type"`
	EngineVersion     *string   `awsName:"EngineVersion" awsType:"awsstr" templateName:"engine-version"`
	Securitygroups    []*string `awsName:"SecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	MaintenanceWindow *string   `awsName:"PreferredMaintenanceWindow" awsType:"awsstr" templateName:"maintenance-window"`
	ApplyImmediately  *bool     `awsName:"ApplyImmediately" awsType:"awsbool" templateName:"apply-immediately"`
}

func (cmd *UpdateCachecluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.Opt(
			params.Suggested("apply-immediately"),
			"nodes", "type", "engine-version", "securitygroups", "maintenance-window",
		),
	))
}

type CreateCachesubnetgroup struct {
	_           string `action:"create" entity:"cachesubnetgroup" awsAPI:"elasticache" awsCall:"CreateCacheSubnetGroup" awsInput:"elasticache.CreateCacheSubnetGroupInput" awsOutput:"elasticache.CreateCacheSubnetGroupOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *elasticache.Client
	Name        *string   `awsName:"CacheSubnetGroupName" awsType:"awsstr" templateName:"name"`
	Description *string   `awsName:"CacheSubnetGroupDescription" awsType:"awsstr" templateName:"description"`
	Subnets     []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
}

func (cmd *CreateCachesubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("subnets"),
		params.Opt(params.Suggested("description")),
	))
}

func (cmd *CreateCachesubnetgroup) ExtractResult(i any) string {
	out, ok := i.(*elasticache.CreateCacheSubnetGroupOutput)
	if !ok || out.CacheSubnetGroup == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.CacheSubnetGroup.CacheSubnetGroupName)
}

type DeleteCachesubnetgroup struct {
	_      string `action:"delete" entity:"cachesubnetgroup" awsAPI:"elasticache" awsCall:"DeleteCacheSubnetGroup" awsInput:"elasticache.DeleteCacheSubnetGroupInput" awsOutput:"elasticache.DeleteCacheSubnetGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *elasticache.Client
	Name   *string `awsName:"CacheSubnetGroupName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteCachesubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type UpdateCachesubnetgroup struct {
	_           string `action:"update" entity:"cachesubnetgroup" awsAPI:"elasticache" awsCall:"ModifyCacheSubnetGroup" awsInput:"elasticache.ModifyCacheSubnetGroupInput" awsOutput:"elasticache.ModifyCacheSubnetGroupOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *elasticache.Client
	Name        *string   `awsName:"CacheSubnetGroupName" awsType:"awsstr" templateName:"name"`
	Description *string   `awsName:"CacheSubnetGroupDescription" awsType:"awsstr" templateName:"description"`
	Subnets     []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
}

func (cmd *UpdateCachesubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("description", "subnets"),
	))
}

type CreateReplicationgroup struct {
	_                 string `action:"create" entity:"replicationgroup" awsAPI:"elasticache" awsCall:"CreateReplicationGroup" awsInput:"elasticache.CreateReplicationGroupInput" awsOutput:"elasticache.CreateReplicationGroupOutput"`
	logger            *logger.Logger
	graph             cloud.GraphAPI
	api               *elasticache.Client
	ID                *string   `awsName:"ReplicationGroupId" awsType:"awsstr" templateName:"id"`
	Description       *string   `awsName:"ReplicationGroupDescription" awsType:"awsstr" templateName:"description"`
	Engine            *string   `awsName:"Engine" awsType:"awsstr" templateName:"engine"`
	EngineVersion     *string   `awsName:"EngineVersion" awsType:"awsstr" templateName:"engine-version"`
	Type              *string   `awsName:"CacheNodeType" awsType:"awsstr" templateName:"type"`
	Replicas          *int64    `awsName:"NumCacheClusters" awsType:"awsint64" templateName:"clusters"`
	SubnetGroup       *string   `awsName:"CacheSubnetGroupName" awsType:"awsstr" templateName:"subnet-group"`
	Securitygroups    []*string `awsName:"SecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	Port              *int64    `awsName:"Port" awsType:"awsint64" templateName:"port"`
	MultiAZ           *bool     `awsName:"MultiAZEnabled" awsType:"awsbool" templateName:"multi-az"`
	AutomaticFailover *bool     `awsName:"AutomaticFailoverEnabled" awsType:"awsbool" templateName:"automatic-failover"`
	AtRestEncryption  *bool     `awsName:"AtRestEncryptionEnabled" awsType:"awsbool" templateName:"at-rest-encryption"`
	TransitEncryption *bool     `awsName:"TransitEncryptionEnabled" awsType:"awsbool" templateName:"transit-encryption"`
	SnapshotRetention *int64    `awsName:"SnapshotRetentionLimit" awsType:"awsint64" templateName:"snapshot-retention"`
}

func (cmd *CreateReplicationgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"), params.Key("description"),
		params.Opt(
			params.Suggested("engine", "type", "clusters"),
			"engine-version", "subnet-group", "securitygroups", "port", "multi-az",
			"automatic-failover", "at-rest-encryption", "transit-encryption",
			"snapshot-retention",
		),
	))
}

func (cmd *CreateReplicationgroup) ExtractResult(i any) string {
	out, ok := i.(*elasticache.CreateReplicationGroupOutput)
	if !ok || out.ReplicationGroup == nil {
		return StringValue(cmd.ID)
	}
	return awssdk.ToString(out.ReplicationGroup.ReplicationGroupId)
}

type DeleteReplicationgroup struct {
	_        string `action:"delete" entity:"replicationgroup" awsAPI:"elasticache" awsCall:"DeleteReplicationGroup" awsInput:"elasticache.DeleteReplicationGroupInput" awsOutput:"elasticache.DeleteReplicationGroupOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *elasticache.Client
	ID       *string `awsName:"ReplicationGroupId" awsType:"awsstr" templateName:"id"`
	Snapshot *string `awsName:"FinalSnapshotIdentifier" awsType:"awsstr" templateName:"snapshot"`
	// RetainPrimaryCluster keeps the primary node running and removes only the
	// replicas, which is the difference between shrinking a group and destroying it.
	RetainPrimary *bool `awsName:"RetainPrimaryCluster" awsType:"awsbool" templateName:"retain-primary"`
}

func (cmd *DeleteReplicationgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.Opt("snapshot", "retain-primary"),
	))
}
