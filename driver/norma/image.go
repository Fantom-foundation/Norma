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
	"io"
	"fmt"
	"errors"
	"path/filepath"

	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/driver/parser"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/archive"
	"github.com/urfave/cli/v2"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
)

// Run with `go run ./driver/norma image`

var imageCommand = cli.Command{
	Name:   "image",
	Usage:  "manages client docker images created by Norma.",
	Subcommands: []*cli.Command{
		{
			Name: "ls",
			Usage: "list client images",
			Action: imageLs,
		},
		{
			Name: "build",
			Usage: "build a client image",
			Action: imageBuild,
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name: "dockerfile",
					Usage: "Dockerfile used to build client image",
					Aliases: []string{"d"},
					EnvVars: []string{"CLIENT_DOCKERFILE"},
					Value:   "./Dockerfile",
				},
				&cli.PathFlag{
					Name: "scenario-file",
					Usage: "target scenario file to extract client versions from",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name: "client-version",
					Usage: "target client versions",
					Aliases: []string{"c"},
				},
			},
		},
		{
			Name: "rm",
			Usage: "remove one or more client images",
			Action: notImplemented,
		},
		{
			Name: "purge",
			Usage: "remove all client images",
			Action: imagePurge,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "force", 
					Aliases: []string{"f"},
					Usage: "force stop container before purging",
				},
			},
		},
	},
}

// ls list all docker images created by norma
func imageLs(ctx *cli.Context) (err error) {
	d, err := newDockerClient()
	if err != nil {
		return err
	}
	defer d.Close()

	filters := filters.NewArgs()
	filters.Add("reference", node.OperaDockerImageName)

	images, err := d.ImageList(context.Background(), types.ImageListOptions{
		All: true,
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

// imageBuild builds an image using 1. dockerfile 2. client version or scenarios
// Example1: norma image build -d /path/to/norma/Dockerfile -c latest
// Example2: norma image build -d /path/to/norma/Dockerfile -s scenarios/small.yml
//    in Example2, all referenced client versions are extracted and built
// Note: we also can set norma's directory for Dockerfile as an env NORMA_DOCKERFILE
//    so we can reduce Example1 to: norma image build -c latest
func imageBuild(ctx *cli.Context) (err error) {
	dockerfile := ctx.String("dockerfile")
	if dockerfile == "" {
		return fmt.Errorf("norma image build cannot proceed without dockerfile")
	}
	if _, err := os.Stat(dockerfile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("targeted dockefile does not exist; %s", err)
	}

	cv := ctx.String("client-version")
	sf := ctx.String("scenario-file")
	if cv == "" && sf == "" {
		return fmt.Errorf("Please provide target client version or scenario fil")
	}
	if cv != "" && sf != "" {
		return fmt.Errorf("Conflict: both target client version and scenario file are provided. Please provide one target.")
	}

	if cv != "" {
		return imageBuildFromClientVersion(dockerfile, cv)
	}

	scenario, err := parser.ParseFile(sf)
	if err != nil {
		return err
	}

	if err := scenario.Check(); err != nil {
		return err
	}

	cvs, err := scenario.ExtractClientVersion()
	if err != nil {
		return err
	}

	errs := []error{}
	for _, cv := range cvs {
		err := imageBuildFromClientVersion(dockerfile, cv)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// imageBuildFromClientVersion creates an image provided 1. dockerfile 2. client version
func imageBuildFromClientVersion(dockerfile string, version string) (err error) {
	fmt.Println(dockerfile, version)
	d, err := newDockerClient()
	if err != nil {
		return err
	}

	buildCtx, _ := archive.TarWithOptions(
		filepath.Dir(dockerfile),
		&archive.TarOptions{},
	)
	buildOpts := types.ImageBuildOptions {
		Dockerfile: filepath.Base(dockerfile),
		Tags:       []string{fmt.Sprintf("%s:%s", node.OperaDockerImageName, version)},
	}

	buildResp, err := d.ImageBuild(context.Background(), buildCtx, buildOpts)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, buildResp.Body)
	return nil
}


// purge removes all images, --force to also include currently running container
func imagePurge(ctx *cli.Context) (err error) {
	var force = ctx.Bool("force")

	d, err := newDockerClient()
	if err != nil {
		return err
	}
	
	filters := filters.NewArgs()
	filters.Add("reference", node.OperaDockerImageName)

	images, err := d.ImageList(context.Background(), types.ImageListOptions{
		Filters: filters,
	})
	for _, image := range images {
		d.ImageRemove(
			context.Background(), 
			image.ID[7:19], 
			types.ImageRemoveOptions{Force: force},
		)
	}

	return nil
}

// newDockerClient creates a docker cli client
func newDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(
		client.FromEnv, 
		client.WithAPIVersionNegotiation(),
	)
}

//notImplemented() is a placeholder func
func notImplemented(ctx *cli.Context) (err error) {
	return nil
}

