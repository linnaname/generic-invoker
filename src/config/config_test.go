package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigs_LoadAndGetConfigByServiceName(t *testing.T) {
	assert.Nil(t, configs)
	config, err := LoadAndGetConfigByServiceName("UserService")
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Param["id"], "uint64")
	assert.Equal(t, config.Result["name"], "string")

	assert.NotNil(t, configs)
	mconfig, err := LoadAndGetConfigByServiceName("MemberService")
	assert.Error(t, err)
	assert.Nil(t, mconfig)

}
