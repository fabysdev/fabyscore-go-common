package env

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// String returns an environment variable of type string.
// If the variable is not set or is empty the process will exit fatally.
func String(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal().Str("code", "COMMON-ENV-STRING-VALUEMISSING").Str("variable", key).Msg("missing value in environment-variable")
	}

	return v
}

// StringDefault returns an environment variable of type string or the default value.
func StringDefault(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	return v
}

// Bool returns an environment variable of type bool.
// If the variable is not set or not parseable by strconv.ParseBool the process will exit fatally.
func Bool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal().Str("code", "COMMON-ENV-BOOL-VALUEMISSING").Str("variable", key).Msg("missing value in environment-variable")
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatal().Str("code", "COMMON-ENV-BOOL-VALUEINVALID").Str("variable", key).Str("value", v).Msg("invalid value in environment-variable")
	}

	return b
}

// BoolDefault returns an environment variable of type bool or the default value.
// If the variable is not parseable by strconv.ParseBool a WARN log will be generated and the default value will be returned.
func BoolDefault(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Warn().Str("code", "COMMON-ENV-BOOLDEFAULT-VALUEINVALID").Str("variable", key).Str("value", v).Msg("invalid value in environment-variable")
		return defaultValue
	}

	return b
}

// Int returns an environment variable of type int.
// If the variable is not set or is not parseable by strconv.Atoi the process will exit fatally.
func Int(key string) int {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal().Str("code", "COMMON-ENV-INT-VALUEMISSING").Str("variable", key).Msg("missing value in environment-variable")
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatal().Str("code", "COMMON-ENV-INT-VALUEINVALID").Str("variable", key).Str("value", v).Msg("invalid value in environment-variable")
	}

	return i
}

// IntDefault returns an environment variable of type int  or the default value.
// If the variable is not parseable by strconv.Atoi a WARN log will be generated and the default value will be returned.
func IntDefault(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		log.Warn().Str("code", "COMMON-ENV-INTDEFAULT-VALUEINVALID").Str("variable", key).Str("value", v).Msg("invalid value in environment-variable")
		return defaultValue
	}

	return i
}
