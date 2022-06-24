package commons

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func HashLocalFileMD5(sourcePath string) (string, error) {
	hashAlg := md5.New()
	return HashLocalFile(sourcePath, hashAlg)
}

func HashLocalFile(sourcePath string, hashAlg hash.Hash) (string, error) {
	f, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(hashAlg, f)
	if err != nil {
		return "", err
	}

	sumBytes := hashAlg.Sum(nil)
	sumString := hex.EncodeToString(sumBytes)

	return sumString, nil
}
