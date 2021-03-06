package rke

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"sort"

	gover "github.com/hashicorp/go-version"
)

func base64Encode(s string) string {
	if len(s) == 0 {
		return ""
	}
	data := []byte(s)

	return base64.StdEncoding.EncodeToString(data)
}

func base64Decode(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(s)

	return string(data), err
}

func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func toArrayString(in []interface{}) []string {
	out := make([]string, len(in))
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toArrayInterface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func toMapString(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toMapInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for i, v := range in {
		out[i] = v
	}
	return out
}

func fileExist(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}

func sortVersions(list map[string]string) ([]*gover.Version, error) {
	var versions []*gover.Version
	for key := range list {
		v, err := gover.NewVersion(key)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	sort.Sort(gover.Collection(versions))
	return versions, nil
}

func getLatestVersion(list map[string]string) (string, error) {
	sorted, err := sortVersions(list)
	if err != nil {
		return "", err
	}

	return sorted[len(sorted)-1].String(), nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func privateKeyToPEM(key *rsa.PrivateKey) string {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	return string(pemdata)
}

func certificateToPEM(cert *x509.Certificate) string {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return string(pemdata)
}
