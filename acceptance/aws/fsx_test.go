package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	fsxtypes "github.com/aws/aws-sdk-go-v2/service/fsx/types"
)

// FSx is four different file systems behind one API, each with its own configuration block.
// This checks the Lustre block decodes and that none of the other three is sent.
func TestCreateFsxfilesystemLustre(t *testing.T) {
	lustre := docFile(t, "lustre.json",
		`{"deploymentType": "PERSISTENT_2", "perUnitStorageThroughput": 125, "dataCompressionType": "LZ4"}`)

	mock := NewMock().On("CreateFileSystem", &fsx.CreateFileSystemOutput{
		FileSystem: &fsxtypes.FileSystem{FileSystemId: awssdk.String("fs-1234abcd")},
	})

	Template("create fsxfilesystem type=LUSTRE capacity=1200 subnets=subnet-1 " +
		"securitygroups=sg-1234 lustre-file=" + lustre).
		Mock(mock).
		ExpectCalls("CreateFileSystem").
		ExpectCommandResult("fs-1234abcd").
		ExpectRevert("delete fsxfilesystem id=fs-1234abcd").
		Run(t)

	in := mock.InputFor("CreateFileSystem").(*fsx.CreateFileSystemInput)
	if got := string(in.FileSystemType); got != "LUSTRE" {
		t.Errorf("FileSystemType: got %q, want LUSTRE", got)
	}
	if got := awssdk.ToInt32(in.StorageCapacity); got != 1200 {
		t.Errorf("StorageCapacity: got %d, want 1200", got)
	}
	if in.LustreConfiguration == nil {
		t.Fatal("LustreConfiguration was not decoded")
	}
	if got := string(in.LustreConfiguration.DeploymentType); got != "PERSISTENT_2" {
		t.Errorf("LustreConfiguration.DeploymentType: got %q", got)
	}
	if got := awssdk.ToInt32(in.LustreConfiguration.PerUnitStorageThroughput); got != 125 {
		t.Errorf("PerUnitStorageThroughput: got %d, want 125", got)
	}
	// Sending a configuration block for the wrong type is rejected by AWS, so the others
	// must stay nil.
	if in.WindowsConfiguration != nil || in.OntapConfiguration != nil || in.OpenZFSConfiguration != nil {
		t.Error("only the Lustre configuration should be set")
	}
}

// The four configuration blocks are alternatives; two at once is meaningless.
func TestCreateFsxfilesystemConfigurationIsExclusive(t *testing.T) {
	lustre := docFile(t, "lustre.json", `{"deploymentType": "SCRATCH_2"}`)
	windows := docFile(t, "windows.json", `{"throughputCapacity": 32}`)

	err := Template("create fsxfilesystem type=LUSTRE capacity=1200 subnets=subnet-1 " +
		"lustre-file=" + lustre + " windows-file=" + windows).
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected the per-type configuration files to be mutually exclusive")
	}
}

func TestDeleteFsxfilesystem(t *testing.T) {
	mock := NewMock().On("DeleteFileSystem", &fsx.DeleteFileSystemOutput{})

	Template("delete fsxfilesystem id=fs-1234abcd").
		Mock(mock).
		ExpectCalls("DeleteFileSystem").
		Run(t)

	in := mock.InputFor("DeleteFileSystem").(*fsx.DeleteFileSystemInput)
	if got := awssdk.ToString(in.FileSystemId); got != "fs-1234abcd" {
		t.Errorf("FileSystemId: got %q", got)
	}
}

// A backup is taken of a file system or of a volume, never both.
func TestCreateFsxbackup(t *testing.T) {
	mock := NewMock().On("CreateBackup", &fsx.CreateBackupOutput{
		Backup: &fsxtypes.Backup{BackupId: awssdk.String("backup-1234")},
	})

	Template("create fsxbackup filesystem=fs-1234abcd").
		Mock(mock).
		ExpectCalls("CreateBackup").
		ExpectCommandResult("backup-1234").
		ExpectRevert("delete fsxbackup id=backup-1234").
		Run(t)

	in := mock.InputFor("CreateBackup").(*fsx.CreateBackupInput)
	if got := awssdk.ToString(in.FileSystemId); got != "fs-1234abcd" {
		t.Errorf("FileSystemId: got %q", got)
	}
	if in.VolumeId != nil {
		t.Error("VolumeId should be unset when a file system was given")
	}
}

func TestCreateFsxbackupRejectsBoth(t *testing.T) {
	err := Template("create fsxbackup filesystem=fs-1234abcd volume=fsvol-5678").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected filesystem and volume to be mutually exclusive")
	}
}

func TestDeleteFsxbackup(t *testing.T) {
	mock := NewMock().On("DeleteBackup", &fsx.DeleteBackupOutput{})

	Template("delete fsxbackup id=backup-1234").
		Mock(mock).
		ExpectCalls("DeleteBackup").
		Run(t)

	in := mock.InputFor("DeleteBackup").(*fsx.DeleteBackupInput)
	if got := awssdk.ToString(in.BackupId); got != "backup-1234" {
		t.Errorf("BackupId: got %q", got)
	}
}
