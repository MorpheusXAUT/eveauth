package web

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AssetChecksums struct {
	Checksums map[string]string
}

func SetupAssetChecksums() (*AssetChecksums, error) {
	checksums := &AssetChecksums{
		Checksums: make(map[string]string),
	}

	err := filepath.Walk("web/assets/css", checksums.CalculateFileChecksum)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk("web/assets/js", checksums.CalculateFileChecksum)
	if err != nil {
		return nil, err
	}

	return checksums, nil
}

func (checksums *AssetChecksums) CalculateFileChecksum(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	hash := md5.New()

	hash.Write(fileContent)

	checksums.Checksums[info.Name()] = fmt.Sprintf("%x", hash.Sum(nil))

	return nil
}
