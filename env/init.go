package env

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	// resolve env
	env := strings.ToLower(os.Getenv("ENV"))
	if env == "" {
		env = "dev"
	}

	// resolve files to load (.env.dev.local > .env.dev > .env; an environment variable will not be overridden if it is already defined)
	files := []string{
		".env." + env + ".local",
		".env." + env,
		".env",
	}

	// find directory with env files (the current working directory could be cmd/ and the env files are in the service root directory for example)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Str("code", "COMMON-ENV-WD").Err(err).Msg("error getting current working directory")
	}

	dir, err := findDirectory(wd, files)
	if err != nil {
		log.Fatal().Str("code", "COMMON-ENV-DIR").Err(err).Msg("error finding root directory")
	}

	if dir == "" {
		dir = wd
	}

	// load env files
	for _, f := range files {
		filePath := path.Join(dir, f)

		err := godotenv.Load(filePath)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal().Str("code", "COMMON-ENV-LOADFILE").Str("file", filePath).Err(err).Msg("error loading env file")
		}
	}
}

func findDirectory(dir string, files []string) (string, error) {
	for _, f := range files {
		exists, err := fileExists(path.Join(dir, f))
		if err != nil {
			return "", err
		}

		// we take the directory if we found atleast one of the files
		if exists {
			return dir, nil
		}
	}

	// try to find in parent
	if dir == "/" {
		return "", nil
	}

	return findDirectory(filepath.Dir(dir), files)
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
