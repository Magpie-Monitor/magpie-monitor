package envs

import (
	"fmt"
	"os"
)

func ValidateEnvs(message string, envKeys []string) {

	for _, env := range envKeys {
		_, isSet := os.LookupEnv(env)
		if !isSet {
			panic(fmt.Sprintf(message, env))
		}
	}

}
