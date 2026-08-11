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
	"github.com/aws/aws-sdk-go-v2/service/redshift"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateRedshiftcluster struct {
	_        string `action:"create" entity:"redshiftcluster" awsAPI:"redshift" awsCall:"CreateCluster" awsInput:"redshift.CreateClusterInput" awsOutput:"redshift.CreateClusterOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *redshift.Client
	ID       *string `awsName:"ClusterIdentifier" awsType:"awsstr" templateName:"id"`
	Username *string `awsName:"MasterUsername" awsType:"awsstr" templateName:"username"`
	Type     *string `awsName:"NodeType" awsType:"awsstr" templateName:"type"`
	Password *string `awsName:"MasterUserPassword" awsType:"awsstr" templateName:"password"`
	// Letting Redshift manage the credential in Secrets Manager avoids putting a
	// password on the command line, where it would also reach the template log.
	ManagePassword *bool     `awsName:"ManageMasterPassword" awsType:"awsbool" templateName:"manage-password"`
	Nodes          *int64    `awsName:"NumberOfNodes" awsType:"awsint64" templateName:"nodes"`
	ClusterType    *string   `awsName:"ClusterType" awsType:"awsstr" templateName:"cluster-type"`
	DBName         *string   `awsName:"DBName" awsType:"awsstr" templateName:"dbname"`
	SubnetGroup    *string   `awsName:"ClusterSubnetGroupName" awsType:"awsstr" templateName:"subnet-group"`
	Securitygroups []*string `awsName:"VpcSecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	Public         *bool     `awsName:"PubliclyAccessible" awsType:"awsbool" templateName:"public"`
	Encrypted      *bool     `awsName:"Encrypted" awsType:"awsbool" templateName:"encrypted"`
	Zone           *string   `awsName:"AvailabilityZone" awsType:"awsstr" templateName:"zone"`
}

// Exactly one of password and manage-password: AWS rejects both together, and with
// neither the call fails for want of a credential.
func (cmd *CreateRedshiftcluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"), params.Key("username"), params.Key("type"),
		params.OnlyOneOf(params.Key("password"), params.Key("manage-password")),
		params.Opt(
			params.Suggested("nodes", "encrypted"),
			"cluster-type", "dbname", "subnet-group", "securitygroups", "public", "zone",
		),
	))
}

func (cmd *CreateRedshiftcluster) ExtractResult(i any) string {
	out, ok := i.(*redshift.CreateClusterOutput)
	if !ok || out.Cluster == nil {
		return StringValue(cmd.ID)
	}
	return awssdk.ToString(out.Cluster.ClusterIdentifier)
}

type DeleteRedshiftcluster struct {
	_      string `action:"delete" entity:"redshiftcluster" awsAPI:"redshift" awsCall:"DeleteCluster" awsInput:"redshift.DeleteClusterInput" awsOutput:"redshift.DeleteClusterOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *redshift.Client
	ID     *string `awsName:"ClusterIdentifier" awsType:"awsstr" templateName:"id"`
	// Redshift refuses to delete without either a final snapshot name or an explicit
	// instruction to skip one, so one of the two is required rather than defaulted —
	// silently skipping a snapshot on a data warehouse is not a safe default.
	Snapshot     *string `awsName:"FinalClusterSnapshotIdentifier" awsType:"awsstr" templateName:"snapshot"`
	SkipSnapshot *bool   `awsName:"SkipFinalClusterSnapshot" awsType:"awsbool" templateName:"skip-snapshot"`
}

func (cmd *DeleteRedshiftcluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.OnlyOneOf(params.Key("snapshot"), params.Key("skip-snapshot")),
	))
}

type UpdateRedshiftcluster struct {
	_              string `action:"update" entity:"redshiftcluster" awsAPI:"redshift" awsCall:"ModifyCluster" awsInput:"redshift.ModifyClusterInput" awsOutput:"redshift.ModifyClusterOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *redshift.Client
	ID             *string   `awsName:"ClusterIdentifier" awsType:"awsstr" templateName:"id"`
	Type           *string   `awsName:"NodeType" awsType:"awsstr" templateName:"type"`
	Nodes          *int64    `awsName:"NumberOfNodes" awsType:"awsint64" templateName:"nodes"`
	ClusterType    *string   `awsName:"ClusterType" awsType:"awsstr" templateName:"cluster-type"`
	Securitygroups []*string `awsName:"VpcSecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	Public         *bool     `awsName:"PubliclyAccessible" awsType:"awsbool" templateName:"public"`
	Version        *string   `awsName:"ClusterVersion" awsType:"awsstr" templateName:"version"`
	Password       *string   `awsName:"MasterUserPassword" awsType:"awsstr" templateName:"password"`
}

func (cmd *UpdateRedshiftcluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.Opt("type", "nodes", "cluster-type", "securitygroups", "public", "version", "password"),
	))
}

func (cmd *UpdateRedshiftcluster) ExtractResult(i any) string {
	return StringValue(cmd.ID)
}

type CreateRedshiftsubnetgroup struct {
	_           string `action:"create" entity:"redshiftsubnetgroup" awsAPI:"redshift" awsCall:"CreateClusterSubnetGroup" awsInput:"redshift.CreateClusterSubnetGroupInput" awsOutput:"redshift.CreateClusterSubnetGroupOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *redshift.Client
	Name        *string   `awsName:"ClusterSubnetGroupName" awsType:"awsstr" templateName:"name"`
	Description *string   `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Subnets     []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
}

// Description is required by the API rather than optional, unlike the equivalent on RDS
// and ElastiCache.
func (cmd *CreateRedshiftsubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("description"), params.Key("subnets"),
	))
}

func (cmd *CreateRedshiftsubnetgroup) ExtractResult(i any) string {
	out, ok := i.(*redshift.CreateClusterSubnetGroupOutput)
	if !ok || out.ClusterSubnetGroup == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.ClusterSubnetGroup.ClusterSubnetGroupName)
}

type DeleteRedshiftsubnetgroup struct {
	_      string `action:"delete" entity:"redshiftsubnetgroup" awsAPI:"redshift" awsCall:"DeleteClusterSubnetGroup" awsInput:"redshift.DeleteClusterSubnetGroupInput" awsOutput:"redshift.DeleteClusterSubnetGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *redshift.Client
	Name   *string `awsName:"ClusterSubnetGroupName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteRedshiftsubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}
