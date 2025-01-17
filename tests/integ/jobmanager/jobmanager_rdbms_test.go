// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// +build integration_storage

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/linuxboot/contest/pkg/storage"
	"github.com/linuxboot/contest/plugins/storage/rdbms"
	"github.com/linuxboot/contest/tests/integ/common"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestJobManagerSuiteRdbmsStorage(t *testing.T) {
	testSuite := TestJobManagerSuite{}

	opts := []rdbms.Opt{
		rdbms.TestEventsFlushSize(1),
		rdbms.TestEventsFlushInterval(10 * time.Second),
	}
	storageLayer, err := common.NewStorage(opts...)
	if err != nil {
		panic(fmt.Sprintf("could not initialize rdbms storage layer: %v", err))
	}
	err = storage.SetStorage(storageLayer)
	require.NoError(t, err)
	testSuite.storage = storageLayer

	suite.Run(t, &testSuite)
}
