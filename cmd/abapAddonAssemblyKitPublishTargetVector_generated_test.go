package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbapAddonAssemblyKitPublishTargetVectorCommand(t *testing.T) {

	testCmd := AbapAddonAssemblyKitPublishTargetVectorCommand()

	// only high level testing performed - details are tested in step generation procudure
	assert.Equal(t, "abapAddonAssemblyKitPublishTargetVector", testCmd.Use, "command name incorrect")

}
