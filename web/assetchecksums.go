package web

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// AssetChecksums stores MD5 checksums for all JS and CSS assets the web app uses
type AssetChecksums struct {
	// Checksums represents a map of checksums, indexed by the file's name
	Checksums map[string]string
}

// SetupAssetChecksums prepares the checksums of all files by walking the appropriate directories and calculating MD5 values
func SetupAssetChecksums() (*AssetChecksums, error) {
	checksums := &AssetChecksums{
		Checksums: make(map[string]string),
	}

	err := filepath.Walk("app/assets/css", checksums.CalculateFileChecksum)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk("app/assets/js", checksums.CalculateFileChecksum)
	if err != nil {
		return nil, err
	}

	return checksums, nil
}

// CalculateFileChecksum reads the content of a file and calculates its MD5 checksum
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
