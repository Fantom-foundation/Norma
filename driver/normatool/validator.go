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

package main

import (
	"crypto/ecdsa"
	"fmt"
	"path"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
	"github.com/Fantom-foundation/go-opera/valkeystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

// fakeKey returns a hex privatekey from a list of prepared fake privatekey
func fakeKey(n uint32) (k *ecdsa.PrivateKey, retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("failed to get key #%d; %s", n, err)
		}
	}()
	k = evmcore.FakeKey(n)
	return k, nil
}

// Run with `go run ./driver/normatool validator`
var validatorCommand = cli.Command{
	Name:  "validator",
	Usage: "mimics sonictool validator",
	Subcommands: []*cli.Command{
		{
			Name:   "from",
			Usage:  "create new validator key from validator-id or private key",
			Action: generateValidatorFrom,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "validator-id",
					Usage:   "validator id from which to generate norma-privatenet validator keys",
					Aliases: []string{"id"},
				},
				&cli.StringFlag{
					Name:    "validator-private-key",
					Usage:   "validator private key from which to generate norma-privatenet validator keys",
					Aliases: []string{"k"},
				},
				&cli.PathFlag{
					Name:    "data-directory",
					Usage:   "data directory for sonic database",
					Aliases: []string{"d", "dd", "dir", "datadir"},
				},
			},
		},
	},
}

// generateValidatorFrom takes a validator-id or validator-privkey and generates corresponding validator-pubkey and pubkeyfile
func generateValidatorFrom(ctx *cli.Context) (err error) {
	datadir := ctx.String("data-directory")
	if datadir == "" {
		fmt.Println("--data-directory unset; skipping secret file generation.")
	}

	id := ctx.Int("validator-id")
	privkey := ctx.String("validator-private-key")

	if id == 0 && privkey == "" {
		return fmt.Errorf("At least one target (--validator-id or --validator-private-key) required for <normatool validator from>")
	}

	if id != 0 && privkey != "" {
		return fmt.Errorf("Target ambiguous (--validator-id and --validator-private-key provided for <normatool validator from>)")
	}

	var privateKeyECDSA *ecdsa.PrivateKey
	if id != 0 {
		privateKeyECDSA, err = fakeKey(uint32(id))
		if err != nil {
			return fmt.Errorf("failed to fake privatekey id %d; %s", id, err)
		}
	} else {
		privateKeyECDSA, err = crypto.ToECDSA(hexutil.MustDecode(privkey))
		if err != nil {
			return fmt.Errorf("failed to decode provided privkey %s; %s", privkey, err)
		}
	}

	// mimicing sonictool
	privateKey := crypto.FromECDSA(privateKeyECDSA)
	publicKey := validatorpk.PubKey{
		Raw:  crypto.FromECDSAPub(&privateKeyECDSA.PublicKey),
		Type: validatorpk.Types.Secp256k1,
	}

	var pathSecretFile string = ""
	if datadir != "" {
		valKeystore := valkeystore.NewDefaultFileRawKeystore(path.Join(datadir, "keystore", "validator"))
		err = valKeystore.Add(publicKey, privateKey, "password")
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		// Sanity check
		_, err = valKeystore.Get(publicKey, "password")
		if err != nil {
			return fmt.Errorf("failed to decrypt the account: %w", err)
		}

		pathSecretFile = valKeystore.PathOf(publicKey)
	}

	// Print it out for user
	// <pubkey> <pub address> <path-to-secret>
	fmt.Printf("%s %s %s",
		publicKey.String(),
		crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		pathSecretFile,
	)

	return nil
}
