package util

import "os"

func Env(key string) string {
	if os.Getenv(key) == "" {
		panic("missing env : " + key)
	}
	return os.Getenv(key)
}
