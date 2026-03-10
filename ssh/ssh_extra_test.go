package ssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wallix/awless/logger"
)

func TestSSHConfigStringWithProxy(t *testing.T) {
	proxy := &Client{
		Port:    22,
		IP:      "52.26.181.76",
		User:    "ec2-user",
		Keypath: "/path/to/bastion-key",
	}
	c := &Client{
		Port:                  22,
		IP:                    "172.31.78.138",
		User:                  "ubuntu",
		Keypath:               "/path/to/key",
		StrictHostKeyChecking: true,
		Proxy:                 proxy,
	}

	config := c.SSHConfigString("MyHost")
	if !strings.Contains(config, "Hostname 172.31.78.138") {
		t.Fatalf("expected Hostname in config, got: %s", config)
	}
	if !strings.Contains(config, "User ubuntu") {
		t.Fatalf("expected User in config, got: %s", config)
	}
	if !strings.Contains(config, "IdentityFile /path/to/key") {
		t.Fatalf("expected IdentityFile in config, got: %s", config)
	}
	if !strings.Contains(config, "ProxyCommand") {
		t.Fatalf("expected ProxyCommand in config, got: %s", config)
	}
	if !strings.Contains(config, "ec2-user@52.26.181.76") {
		t.Fatalf("expected proxy user@ip in ProxyCommand, got: %s", config)
	}
}

func TestConnectStringWithProxy(t *testing.T) {
	proxy := &Client{
		Port:    22,
		IP:      "52.26.181.76",
		User:    "ec2-user",
		Keypath: "/path/to/bastion-key",
	}
	c := &Client{
		Port:                  22,
		IP:                    "172.31.78.138",
		User:                  "ubuntu",
		Keypath:               "/path/to/key",
		StrictHostKeyChecking: true,
		Proxy:                 proxy,
	}

	cs := c.ConnectString()
	if !strings.Contains(cs, "ubuntu@172.31.78.138") {
		t.Fatalf("expected user@ip in connect string, got: %s", cs)
	}
	if !strings.Contains(cs, "ProxyCommand") {
		t.Fatalf("expected ProxyCommand in connect string, got: %s", cs)
	}
	if !strings.Contains(cs, "-i /path/to/key") {
		t.Fatalf("expected key in connect string, got: %s", cs)
	}
}

func TestConnectStringWithProxyNoKey(t *testing.T) {
	proxy := &Client{
		Port: 22,
		IP:   "52.26.181.76",
		User: "ec2-user",
	}
	c := &Client{
		Port:                  22,
		IP:                    "172.31.78.138",
		User:                  "ubuntu",
		StrictHostKeyChecking: true,
		Proxy:                 proxy,
	}

	cs := c.ConnectString()
	if !strings.Contains(cs, "ProxyCommand") {
		t.Fatalf("expected ProxyCommand, got: %s", cs)
	}
}

func TestSSHConfigStringWithNonDefaultPort(t *testing.T) {
	c := &Client{
		Port:                  2222,
		IP:                    "1.2.3.4",
		User:                  "admin",
		StrictHostKeyChecking: true,
	}
	config := c.SSHConfigString("MyServer")
	if !strings.Contains(config, "Port 2222") {
		t.Fatalf("expected Port in config, got: %s", config)
	}
}

func TestSSHConfigStringWithStrictHostKeyCheckingOff(t *testing.T) {
	c := &Client{
		Port:                  22,
		IP:                    "1.2.3.4",
		User:                  "admin",
		StrictHostKeyChecking: false,
	}
	config := c.SSHConfigString("MyServer")
	if !strings.Contains(config, "StrictHostKeychecking no") {
		t.Fatalf("expected StrictHostKeychecking no in config, got: %s", config)
	}
}

func TestCloseAll(t *testing.T) {
	// nil client
	var c *Client
	if err := c.CloseAll(); err != nil {
		t.Fatalf("nil client CloseAll: %s", err)
	}

	// Client with no inner client or proxy
	c2 := &Client{}
	if err := c2.CloseAll(); err != nil {
		t.Fatalf("empty client CloseAll: %s", err)
	}
}

func TestSetLogger(t *testing.T) {
	c := &Client{logger: logger.DiscardLogger}
	newLogger := logger.New("test", 0)
	c.SetLogger(newLogger)
	if c.logger != newLogger {
		t.Fatal("SetLogger did not update logger")
	}
}

func TestSetStrictHostKeyChecking(t *testing.T) {
	c := &Client{StrictHostKeyChecking: true}
	c.SetStrictHostKeyChecking(false)
	if c.StrictHostKeyChecking != false {
		t.Fatal("expected StrictHostKeyChecking to be false")
	}
	c.SetStrictHostKeyChecking(true)
	if c.StrictHostKeyChecking != true {
		t.Fatal("expected StrictHostKeyChecking to be true")
	}
}

func TestFindPrivateKeyFromNameEmpty(t *testing.T) {
	_, ok := findPrivateKeyFromName("")
	if ok {
		t.Fatal("expected false for empty keyname")
	}
}

func TestFindPrivateKeyFromNameAbsolutePath(t *testing.T) {
	// Create temp file
	dir, err := os.MkdirTemp("", "ssh-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	keypath := filepath.Join(dir, "mykey.pem")
	if err := os.WriteFile(keypath, []byte("fake-key"), 0600); err != nil {
		t.Fatal(err)
	}

	// Absolute path should be found directly, ignoring keyFolders
	priv, ok := findPrivateKeyFromName(keypath, "/nonexistent")
	if !ok {
		t.Fatal("expected to find key at absolute path")
	}
	if priv.path != keypath {
		t.Fatalf("got %q, want %q", priv.path, keypath)
	}
}

func TestFindPrivateKeyFromNameNonexistent(t *testing.T) {
	_, ok := findPrivateKeyFromName("nonexistent-key-12345", "/tmp")
	if ok {
		t.Fatal("expected false for nonexistent key")
	}
}

func TestFindPrivateKeyFromNameInFolder(t *testing.T) {
	dir, err := os.MkdirTemp("", "ssh-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	keypath := filepath.Join(dir, "test-key.pem")
	if err := os.WriteFile(keypath, []byte("fake-key-data"), 0600); err != nil {
		t.Fatal(err)
	}

	// Search by name without extension
	priv, ok := findPrivateKeyFromName("test-key", dir)
	if !ok {
		t.Fatal("expected to find key in folder")
	}
	if priv.path != keypath {
		t.Fatalf("got path %q, want %q", priv.path, keypath)
	}
}

func TestFindPrivateKeyFromNameWithPemSuffix(t *testing.T) {
	dir, err := os.MkdirTemp("", "ssh-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	keypath := filepath.Join(dir, "mykey.pem")
	if err := os.WriteFile(keypath, []byte("fake-key-data"), 0600); err != nil {
		t.Fatal(err)
	}

	// Search with .pem extension
	priv, ok := findPrivateKeyFromName("mykey.pem", dir)
	if !ok {
		t.Fatal("expected to find key with .pem suffix")
	}
	if priv.path != keypath {
		t.Fatalf("got path %q, want %q", priv.path, keypath)
	}
}

func TestConnectStringNoKeyNoCustomPort(t *testing.T) {
	c := &Client{
		Port:                  22,
		IP:                    "10.0.0.1",
		User:                  "root",
		StrictHostKeyChecking: true,
	}

	cs := c.ConnectString()
	if strings.Contains(cs, "-i") {
		t.Fatalf("should not contain -i flag without keypath, got: %s", cs)
	}
	if strings.Contains(cs, "-p") {
		t.Fatalf("should not contain -p flag for default port, got: %s", cs)
	}
	if !strings.Contains(cs, "root@10.0.0.1") {
		t.Fatalf("expected user@ip, got: %s", cs)
	}
}
