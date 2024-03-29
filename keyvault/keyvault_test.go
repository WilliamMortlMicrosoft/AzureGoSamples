// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package keyvault

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/marstr/randname"

	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/config"
	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/util"
	"github.com/WilliamMortlMicrosoft/AzureGoSamples/resources"
)

var (
	kvName  = randname.GenerateWithPrefix("vault-sample-go-", 5)
	keyName = randname.GenerateWithPrefix("key-sample-go-", 5)
)

// TestMain sets up the environment and initiates tests.
func TestMain(m *testing.M) {
	var err error
	err = config.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err.Error())
	}

	err = config.AddFlags()
	if err != nil {
		log.Fatalf("failed to parse flags: %v\n", err.Error())
	}
	flag.Parse()

	code := m.Run()
	os.Exit(code)
}

func ExampleSetVaultPermissions() {
	var groupName = config.GenerateGroupName("KeyVault")
	config.SetGroupName(groupName)

	ctx := context.Background()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, config.GroupName())
	if err != nil {
		util.PrintAndLog(err.Error())
	}

	_, err = CreateVault(ctx, kvName)
	if err != nil {
		util.PrintAndLog(err.Error())
	}
	util.PrintAndLog("vault created")

	_, err = SetVaultPermissions(ctx, kvName)
	if err != nil {
		util.PrintAndLog(err.Error())
	}
	util.PrintAndLog("set vault permissions")

	_, err = CreateKey(ctx, kvName, keyName)
	if err != nil {
		util.PrintAndLog(err.Error())
	}
	util.PrintAndLog("created key")

	// Output:
	// vault created
	// set vault permissions
	// created key
}
