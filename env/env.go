package env

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
 This package encompasses the methods to cleanly pull information of the correct type out of
 ENV variables that would normally be set by a deploy script on the container or box
*/

func GetEnv(key string) (string, bool) {
	str := os.Getenv(toEnvKey(key))
	if str == "" {
		return "", false
	}
	return str, true
}

func toEnvKey(name string) string {
	key := strings.ToUpper(strings.Replace(name, "-", "_", -1))
	return key
}

func FlagOrEnvDuration(p *time.Duration, name string, value time.Duration, usage string) {
	v, ok := GetEnv(name)
	if ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			value = time.Duration(i)
		}
	}
	flag.DurationVar(p, name, value, usage)
}
