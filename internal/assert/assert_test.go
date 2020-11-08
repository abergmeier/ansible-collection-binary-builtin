package assert

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"go.starlark.net/starlark"
)

func TestAssertVersion(t *testing.T) {
	thread := &starlark.Thread{
		Name: t.Name(),
	}
	truth, err := evalExprTruth(thread, 0, `ansible_version.full is version('2.1', '>=')`)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !truth {
		t.Fatal("Unexpected false result")
	}
}

func TestAssertDirectory(t *testing.T) {
	thread := &starlark.Thread{
		Name: t.Name(),
	}

	dirName, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(dirName)

	truth, err := evalExprTruth(thread, 0, fmt.Sprintf(`"%s" is directory`, dirName))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !truth {
		t.Fatal("Unexpected false result")
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	defer f.Close()

	truth, err = evalExprTruth(thread, 0, fmt.Sprintf(`"%s" is directory`, f.Name()))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if truth {
		t.Fatal("Unexpected true result")
	}
}

func TestAssertExists(t *testing.T) {
	thread := &starlark.Thread{
		Name: t.Name(),
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	defer f.Close()

	truth, err := evalExprTruth(thread, 0, fmt.Sprintf(`"%s" is exists`, f.Name()))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !truth {
		t.Fatal("Unexpected false result")
	}
}

func TestAssertFile(t *testing.T) {
	thread := &starlark.Thread{
		Name: t.Name(),
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	defer f.Close()

	truth, err := evalExprTruth(thread, 0, fmt.Sprintf(`"%s" is file`, f.Name()))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !truth {
		t.Fatal("Unexpected false result")
	}

	dirName, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(dirName)

	truth, err = evalExprTruth(thread, 0, fmt.Sprintf(`"%s" is file`, dirName))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if truth {
		t.Fatal("Unexpected true result")
	}
}
