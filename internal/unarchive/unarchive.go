package unarchive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Parameters to feed to this module
type Parameters struct {
	Src  string
	Dest string
}

func Run(p *Parameters) error {

	r, err := zip.OpenReader(p.Src)
	if err != nil {
		return err
	}
	defer r.Close()

	return writeFilesToDir(r.File, p.Dest)
}

func writeFilesToDir(files []*zip.File, dest string) error {
	for _, f := range files {

		fp := filepath.Join(dest, f.Name)

		// Prevent ZipSlip
		if !strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("Illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}

		err := writeZippedFile(f, fp)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeZippedFile(f *zip.File, fp string) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	return copyFileContent(r, fp, f.Mode())
}

func copyFileContent(r io.Reader, fp string, mode os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	if err != nil {
		return err
	}
	outFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, r)
	if err != nil {
		return err
	}
	return outFile.Sync()
}
