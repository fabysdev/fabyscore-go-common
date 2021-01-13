package secret

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fabysdev/fabyscore-go-common/env"
)

// secrets is the in-memory cache for the secrets.
var mu sync.RWMutex
var secrets = map[string]string{}
var replacer = strings.NewReplacer("-", "_", ".", "_")

// ErrNotFound is used if the secret is not found.
var ErrNotFound = errors.New("secret not found")

// Get returns the secret value from the in-memory cache or tries to load the secret.
// If no value is found a ErrNotFound error will be returned.
func Get(path string) (string, error) {
	// cache
	if s, ok := secrets[path]; ok {
		return s, nil
	}

	// file
	if f, err := os.Stat(path); !os.IsNotExist(err) && !f.IsDir() {
		value, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		if len(value) == 0 {
			return "", ErrNotFound
		}

		s := string(value)
		secrets[path] = s

		return s, nil
	}

	// env
	envKey := replacer.Replace(filepath.Base(path))

	value := env.StringDefault(strings.ToUpper(envKey), "")
	if value != "" {
		secrets[path] = value
		return value, nil
	}

	return "", ErrNotFound
}
