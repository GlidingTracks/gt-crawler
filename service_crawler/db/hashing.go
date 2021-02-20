package db

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"io"
)




// HashFileWithSha256 hashes the content of the file with SHA256.
// Note: this will result in 64-character long string.
func HashFileWithSha256(buf []byte) string {
	sha := sha256.New()
	_, err := sha.Write(buf)
	if err != nil && err != io.EOF {
		logrus.Fatal(err)
	}

	sum := sha.Sum(nil)
	return hex.EncodeToString(sum)
}
