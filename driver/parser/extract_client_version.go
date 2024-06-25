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

package parser

import (
	"golang.org/x/exp/maps"
)

const DefaultClientVersion = "latest"

// ExtractClientVersion returns the list of all referenced client versions in a scenario.
func (s *Scenario) ExtractClientVersion() ([]string, error) {
	cvs := make(map[string]struct{}) //set of seen client versions

	for _, node := range s.Nodes {
		cv, err := node.ExtractClientVersion()
		if err != nil {
			return nil, err
		}

		if _, ok := cvs[cv]; !ok {
			cvs[cv] = struct{}{}
		}
	}

	return maps.Keys(cvs), nil
}

// ExtractClientVersion returns the node's client versions.
func (n *Node) ExtractClientVersion() (string, error) {
	if n.Client.Version == nil {
		return DefaultClientVersion, nil
	}

	return *n.Client.Version, nil
}
