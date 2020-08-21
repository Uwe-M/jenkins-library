package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbapAddonAssemblyKitReserveNextPackagesCommand(t *testing.T) {

	testCmd := AbapAddonAssemblyKitReserveNextPackagesCommand()

	// only high level testing performed - details are tested in step generation procudure
	assert.Equal(t, "abapAddonAssemblyKitReserveNextPackages", testCmd.Use, "command name incorrect")

}