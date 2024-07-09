package tests

import (
	"fmt"
	"os"
)

func getContainerImage(name string) string {
	val, ok := os.LookupEnv("TESTCONTAINERS_HUB_IMAGE_NAME_PREFIX")
	if ok {
		return fmt.Sprintf("%s/%s", val, name)
	}
	return name
}
