//go:build ext

package ext

import (
	_ "github.com/daichirata/hammer"
)

//go:generate go build -o $GO_ARTIFACT_PATH github.com/daichirata/hammer

