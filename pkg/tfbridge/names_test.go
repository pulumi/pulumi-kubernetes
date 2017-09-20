// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPulumiToTerraformName(t *testing.T) {
	assert.Equal(t, "", PulumiToTerraformName(""))
	assert.Equal(t, "test", PulumiToTerraformName("test"))
	assert.Equal(t, "test_name", PulumiToTerraformName("testName"))
	assert.Equal(t, "test_name_pascal", PulumiToTerraformName("TestNamePascal"))
	assert.Equal(t, "test_name", PulumiToTerraformName("test_name"))
	assert.Equal(t, "test_name_", PulumiToTerraformName("testName_"))
	assert.Equal(t, "t_e_s_t_n_a_m_e", PulumiToTerraformName("TESTNAME"))
}

func TestTerraformToPulumiName(t *testing.T) {
	assert.Equal(t, "", TerraformToPulumiName("", false))
	assert.Equal(t, "test", TerraformToPulumiName("test", false))
	assert.Equal(t, "testName", TerraformToPulumiName("test_name", false))
	assert.Equal(t, "testName_", TerraformToPulumiName("testName_", false))
	assert.Equal(t, "tESTNAME", TerraformToPulumiName("t_e_s_t_n_a_m_e", false))
	assert.Equal(t, "", TerraformToPulumiName("", true))
	assert.Equal(t, "Test", TerraformToPulumiName("test", true))
	assert.Equal(t, "TestName", TerraformToPulumiName("test_name", true))
	assert.Equal(t, "TestName_", TerraformToPulumiName("testName_", true))
	assert.Equal(t, "TESTNAME", TerraformToPulumiName("t_e_s_t_n_a_m_e", true))
}
