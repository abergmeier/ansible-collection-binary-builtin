package git

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/crypto/ssh"
)

var (
	suffix  *big.Int
	tempDir string
)

func init() {
	var err error
	suffix, err = rand.Int(rand.Reader, big.NewInt(999999999999))
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == "windows" {
		tempDir = os.Getenv("TEMP")
	} else {
		tempDir = "/tmp"
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	userDir := filepath.Join(tempDir, fmt.Sprintf("%s_%d", t.Name(), suffix.Int64()))
	defer os.RemoveAll(userDir)
	mockSSHConfig(t, userDir)
	r := runner{
		userDir: userDir,
	}

	dest := filepath.Join(userDir, "go-git-checkout")

	err := r.Run(&Parameters{
		Repo:    "https://github.com/go-git/go-git.git",
		Dest:    dest,
		Version: "d525a514057f97bc2b183e2c67f542dd6f0ac0aa",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSSHAuth(t *testing.T) {
	t.Parallel()

	userDir := filepath.Join(tempDir, fmt.Sprintf("%s_%d", t.Name(), suffix.Int64()))
	defer os.RemoveAll(userDir)
	mockSSHConfig(t, userDir)
	r := runner{
		userDir: userDir,
	}
	pk, err := r.getSSHPublicKeys("ssh://user@server/project.git")
	if err != nil {
		t.Fatal(err)
	}
	if pk.User != "user" {
		t.Fatalf("Unexpected user: %s", pk.User)
	}
}

func mockSSHConfig(t *testing.T, userDir string) {
	sshPath := filepath.Join(userDir, ".ssh")
	err := os.MkdirAll(sshPath, 0700)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}

	writePrivateKey(t, privateKey, sshPath)
	writePublicKey(t, &privateKey.PublicKey, sshPath)
}

func encodePrivateKey(t *testing.T, pk *rsa.PrivateKey) []byte {
	der := x509.MarshalPKCS1PrivateKey(pk)

	// Private key in PEM format
	return pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   der,
	})
}

func writePrivateKey(t *testing.T, privateKey *rsa.PrivateKey, sshPath string) {
	pem := encodePrivateKey(t, privateKey)
	privPath := filepath.Join(sshPath, "id_rsa")
	err := ioutil.WriteFile(privPath, pem, 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func writePublicKey(t *testing.T, privateKey *rsa.PublicKey, sshPath string) {
	publicKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	bytes := ssh.MarshalAuthorizedKey(publicKey)
	pubPath := filepath.Join(sshPath, "id_rsa.pub")
	err = ioutil.WriteFile(pubPath, bytes, 0600)
	if err != nil {
		t.Fatal(err)
	}
}
