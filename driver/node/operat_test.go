package node

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
)

func TestImplements(t *testing.T) {
	var inst OperaNode
	var _ driver.Node = &inst

}
