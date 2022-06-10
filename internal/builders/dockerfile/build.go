// Copyright 2022 SLSA Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/slsa-framework/slsa-github-generator/internal/builders/dockerfile/docker"
)

// buildCmd returns the 'build' command.
func buildCmd() *cobra.Command {
	var filePath string
	var context string
	var tags string

	c := &cobra.Command{
		Use:   "build",
		Short: "Build and push a Docker image.",
		Long: `Build a Docker image from a Dockerfile and push it to an image repository.
This command assumes that it is being run in the context of a Github Actions
workflow.`,

		Run: func(cmd *cobra.Command, args []string) {
			tagList := strings.Split(tags, ",")
			if len(tagList) == 0 {
				check(errors.New("at least one tag name must be specified"))
			}

			if err := docker.Build(docker.BuildOpts{
				ContextDir: context,
				File:       filePath,
				Tags:       tagList,
			}); err != nil {
				check(fmt.Errorf("failed to build Docker image: %w", err))
			}

			for _, tag := range tagList {
				if err := docker.Push(docker.PushOpts{
					Tag: tag,
				}); err != nil {
					check(fmt.Errorf("failed to push to image repository: %w", err))
				}
			}
		},
	}

	// TODO(github.com/slsa-framework/slsa-github-generator/issues/57): flags
	c.Flags().StringVarP(&filePath, "file", "f", "./Dockerfile", "Path to the Dockerfile.")
	c.Flags().StringVarP(&context, "context", "c", ".", "A path to the build context.")
	c.Flags().StringVarP(&tags, "tags", "t", "", "A CSV list of name:tag.")

	return c
}
