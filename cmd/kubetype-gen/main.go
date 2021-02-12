// Copyright 2019 Istio Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"os"

	"github.com/golang/glog"
	"k8s.io/gengo/args"

	"istio.io/tools/cmd/kubetype-gen/generators"
	"istio.io/tools/cmd/kubetype-gen/scanner"
)

func main() {
	arguments := args.Default()

	arguments.GeneratedByCommentTemplate = "// Code generated by kubetype-gen. DO NOT EDIT."

	// Don't default the file header
	// arguments.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), "istio.io/tools/cmd/kubetype-gen/boilerplate.go.txt")

	scanner := scanner.Scanner{}

	if err := arguments.Execute(
		generators.NameSystems("", nil),
		generators.DefaultNameSystem(),
		scanner.Scan,
	); err != nil {
		glog.Errorf("Error: %v", err)
		os.Exit(1)
	}
	glog.V(2).Info("Completed successfully.")
}
