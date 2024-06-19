// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

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
