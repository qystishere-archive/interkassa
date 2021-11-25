package interkassa

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type signAlgorithm string

const (
	SignAlgorithmDefault signAlgorithm = SignAlgorithmSha256
	SignAlgorithmMD5     signAlgorithm = "md5"
	SignAlgorithmSha256  signAlgorithm = "sha256"
)

// Interkassa касса.
type Interkassa struct {
	config *Config
}

// Config конфигурация кассы.
type Config struct {
	ID            string
	SignAlgorithm signAlgorithm
	SignKey       string
	SignTestKey   string
}

// New создаёт новый экземпляр кассы.
func New(config Config) *Interkassa {
	return &Interkassa{
		config: &config,
	}
}

func (ik *Interkassa) sign(data string) string {
	var hasher hash.Hash
	switch ik.config.SignAlgorithm {
	case SignAlgorithmMD5:
		hasher = md5.New()
	case SignAlgorithmSha256:
		fallthrough
	default:
		hasher = sha256.New()
	}
	hasher.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
