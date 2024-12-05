package envs

import (
	"fmt"
	"os"
	"strconv"
)

func ValidateEnvs(message string, envKeys []string) {

	for _, env := range envKeys {
		_, isSet := os.LookupEnv(env)
		if !isSet {
			panic(fmt.Sprintf(message, env))
		}
	}

}

// Validates if env exists and converts it to int type, panics on error
func ConvertToInt(env string) int {
	ValidateEnvs("%s env variable not set", []string{env})

	envInt, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		panic(fmt.Sprintf("%s has to be numeric", env))
	}

	return envInt
}
