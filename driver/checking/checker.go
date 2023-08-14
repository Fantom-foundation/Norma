package checking

import (
	"errors"
	"github.com/Fantom-foundation/Norma/driver"
)

// Checker do the network consistency check at the end of the scenario.
type Checker interface {
	Check(net driver.Network) error
}

func CheckNetworkConsistency(net driver.Network) error {
	checkers := []Checker{
		new(BlockHeightChecker),
		new(BlocksHashesChecker),
	}
	errs := make([]error, len(checkers))
	for i, checker := range checkers {
		errs[i] = checker.Check(net)
	}
	return errors.Join(errs...)
}
