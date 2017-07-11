// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLumiToTerraformName(t *testing.T) {
	assert.Equal(t, "", LumiToTerraformName(""))
	assert.Equal(t, "test", LumiToTerraformName("test"))
	assert.Equal(t, "test_name", LumiToTerraformName("testName"))
	assert.Equal(t, "test_name_pascal", LumiToTerraformName("TestNamePascal"))
	assert.Equal(t, "test_name", LumiToTerraformName("test_name"))
	assert.Equal(t, "test_name_", LumiToTerraformName("testName_"))
	assert.Equal(t, "t_e_s_t_n_a_m_e", LumiToTerraformName("TESTNAME"))
}

func TestTerraformToLumiName(t *testing.T) {
	assert.Equal(t, "", TerraformToLumiName("", false))
	assert.Equal(t, "test", TerraformToLumiName("test", false))
	assert.Equal(t, "testName", TerraformToLumiName("test_name", false))
	assert.Equal(t, "testName_", TerraformToLumiName("testName_", false))
	assert.Equal(t, "tESTNAME", TerraformToLumiName("t_e_s_t_n_a_m_e", false))
	assert.Equal(t, "", TerraformToLumiName("", true))
	assert.Equal(t, "Test", TerraformToLumiName("test", true))
	assert.Equal(t, "TestName", TerraformToLumiName("test_name", true))
	assert.Equal(t, "TestName_", TerraformToLumiName("testName_", true))
	assert.Equal(t, "TESTNAME", TerraformToLumiName("t_e_s_t_n_a_m_e", true))
}
