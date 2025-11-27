//go:build ext

package ext

import (
	_ "go.mercari.io/yo"
)

//go:generate go build -o $GO_ARTIFACT_PATH go.mercari.io/yo

