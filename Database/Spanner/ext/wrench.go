//go:build ext

package ext

import (
	_ "github.com/cloudspannerecosystem/wrench"
)

//go:generate go build -o $GO_ARTIFACT_PATH github.com/cloudspannerecosystem/wrench

