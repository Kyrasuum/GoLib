package env

import (
	"fmt"
	"os"
	"strconv"
)

func CondAssignStrEnv(variable *string, env string) {
	def := *variable
	*variable = os.Getenv(env)
	if *variable == "" {
		fmt.Printf("Unset: %s, using: %s\n", env, def)
		*variable = def
	}
}

func CondAssignIntEnv(variable *int, env string) {
	def := *variable
	varstr := os.Getenv(env)
	if varstr == "" {
		fmt.Printf("Unset: %s, using: %d\n", env, def)
		*variable = def
		return
	}

	value, err := strconv.Atoi(varstr)
	if err != nil {
		fmt.Printf("Invalid integer value for %s, setting to %d\n", env, def)
		value = def
	}
	*variable = value
}
