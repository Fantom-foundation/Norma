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

package app

import (
	"fmt"
	"strings"
)

type appFactoryFunc func(context AppContext, feederId, appId uint32) (Application, error)

func NewApplication(appType string, context AppContext, feederId, appId uint32) (Application, error) {
	if factory := getFactory(appType); factory != nil {
		return factory(context, feederId, appId)
	}
	return nil, fmt.Errorf("unknown application type '%s'", appType)
}

func IsSupportedApplicationType(appType string) bool {
	return getFactory(appType) != nil
}

func getFactory(appType string) appFactoryFunc {
	switch strings.ToLower(appType) {
	case "erc20":
		return NewERC20Application
	case "counter", "":
		return NewCounterApplication
	case "store":
		return NewStoreApplication
	case "uniswap":
		return NewUniswapApplication
	}
	return nil
}
