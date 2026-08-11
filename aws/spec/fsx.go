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
	"github.com/aws/aws-sdk-go-v2/service/fsx"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// The entity is fsxfilesystem: EFS already owns "filesystem", and the two are different
// products that would be genuinely confusing to conflate.
//
// FSx is really four file systems behind one API — Windows, Lustre, ONTAP and OpenZFS — each
// with its own configuration block and none of them interchangeable. The per-type settings
// therefore come from a file rather than being flattened into flags that only apply to one
// type.
type CreateFsxfilesystem struct {
	_              string `action:"create" entity:"fsxfilesystem" awsAPI:"fsx" awsCall:"CreateFileSystem" awsInput:"fsx.CreateFileSystemInput" awsOutput:"fsx.CreateFileSystemOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *fsx.Client
	Type           *string   `awsName:"FileSystemType" awsType:"awsstr" templateName:"type"`
	Capacity       *int64    `awsName:"StorageCapacity" awsType:"awsint64" templateName:"capacity"`
	Subnets        []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
	Securitygroups []*string `awsName:"SecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	StorageType    *string   `awsName:"StorageType" awsType:"awsstr" templateName:"storage-type"`
	KmsKey         *string   `awsName:"KmsKeyId" awsType:"awsstr" templateName:"kms-key"`
	// Exactly one of these applies, chosen by the type above.
	WindowsFile *string `awsName:"WindowsConfiguration" awsType:"awsfiletostruct" templateName:"windows-file"`
	LustreFile  *string `awsName:"LustreConfiguration" awsType:"awsfiletostruct" templateName:"lustre-file"`
	OntapFile   *string `awsName:"OntapConfiguration" awsType:"awsfiletostruct" templateName:"ontap-file"`
	OpenZFSFile *string `awsName:"OpenZFSConfiguration" awsType:"awsfiletostruct" templateName:"openzfs-file"`
}

// The four configuration files are mutually exclusive: passing the Lustre block for a Windows
// file system is rejected, and passing two is meaningless.
func (cmd *CreateFsxfilesystem) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("type"), params.Key("capacity"), params.Key("subnets"),
		params.Opt(
			params.Suggested("securitygroups", "storage-type"),
			"kms-key",
		),
		params.OnlyOneOf(
			params.Key("windows-file"),
			params.Key("lustre-file"),
			params.Key("ontap-file"),
			params.Key("openzfs-file"),
		),
	))
}

func (cmd *CreateFsxfilesystem) ExtractResult(i any) string {
	out, ok := i.(*fsx.CreateFileSystemOutput)
	if !ok || out.FileSystem == nil {
		return ""
	}
	return awssdk.ToString(out.FileSystem.FileSystemId)
}

type DeleteFsxfilesystem struct {
	_      string `action:"delete" entity:"fsxfilesystem" awsAPI:"fsx" awsCall:"DeleteFileSystem" awsInput:"fsx.DeleteFileSystemInput" awsOutput:"fsx.DeleteFileSystemOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *fsx.Client
	ID     *string `awsName:"FileSystemId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteFsxfilesystem) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CreateFsxbackup struct {
	_          string `action:"create" entity:"fsxbackup" awsAPI:"fsx" awsCall:"CreateBackup" awsInput:"fsx.CreateBackupInput" awsOutput:"fsx.CreateBackupOutput"`
	logger     *logger.Logger
	graph      cloud.GraphAPI
	api        *fsx.Client
	FileSystem *string `awsName:"FileSystemId" awsType:"awsstr" templateName:"filesystem"`
	// An ONTAP or OpenZFS backup is taken of a volume rather than the whole file system,
	// so one or the other is given.
	Volume *string `awsName:"VolumeId" awsType:"awsstr" templateName:"volume"`
}

func (cmd *CreateFsxbackup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.OnlyOneOf(params.Key("filesystem"), params.Key("volume")),
	))
}

func (cmd *CreateFsxbackup) ExtractResult(i any) string {
	out, ok := i.(*fsx.CreateBackupOutput)
	if !ok || out.Backup == nil {
		return ""
	}
	return awssdk.ToString(out.Backup.BackupId)
}

type DeleteFsxbackup struct {
	_      string `action:"delete" entity:"fsxbackup" awsAPI:"fsx" awsCall:"DeleteBackup" awsInput:"fsx.DeleteBackupInput" awsOutput:"fsx.DeleteBackupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *fsx.Client
	ID     *string `awsName:"BackupId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteFsxbackup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}
