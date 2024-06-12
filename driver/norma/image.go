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
	"context"
	"os"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver/node"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma image`

var imageCommand = cli.Command{
	Name:  "image",
	Usage: "manages client docker images created by Norma.",
	Subcommands: []*cli.Command{
		{
			Name:   "ls",
			Usage:  "list client images",
			Action: ls,
		},
		{
			Name:   "build",
			Usage:  "build a client image",
			Action: notImplemented,
		},
		{
			Name:   "rm",
			Usage:  "remove one or more client images",
			Action: notImplemented,
		},
		{
			Name:   "purge",
			Usage:  "remove all client images",
			Action: notImplemented,
		},
	},
}

// ls list all docker images created by norma
func ls(ctx *cli.Context) (err error) {
	d, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer d.Close()

	filters := filters.NewArgs()
	filters.Add("reference", node.OperaDockerImageName)

	images, err := d.ImageList(context.Background(), types.ImageListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE",
	})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetColumnSeparator(" ")

	for _, image := range images {
		repository := "<none>"
		tag := "<none>"

		if len(image.RepoTags) > 0 {
			splitted := strings.Split(image.RepoTags[0], ":")
			repository = splitted[0]
			tag = splitted[1]
		} else if len(image.RepoDigests) > 0 {
			repository = strings.Split(image.RepoDigests[0], "@")[0]
		}

		duration := units.HumanDuration(
			time.Now().UTC().Sub(time.Unix(image.Created, 0)),
		) + " ago"
		size := units.HumanSizeWithPrecision(float64(image.Size), 3)

		table.Append([]string{
			repository, tag, image.ID[7:19], duration, size,
		})
	}

	table.Render()

	return nil
}

// notImplemented() is a placeholder func
func notImplemented(ctx *cli.Context) (err error) {
	return nil
}
