// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package resources

import (
	"context"
	"log"

	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/config"
)

// Cleanup deletes the rescource group created for the sample
func Cleanup(ctx context.Context) error {
	if config.KeepResources() {
		log.Println("Hybrid resources cleanup: keeping resources")
		return nil
	}
	log.Println("Hybrid resources cleanup: deleting resources")
	_, err := DeleteGroup(ctx)
	return err
}
