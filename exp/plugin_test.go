package exp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSnmpBasePlugin(t *testing.T) {
	// Configure the plugin to use the configuration defined in the test data.
	// The plugin constructor will look for the plugin config, so it needs
	// to be set to correctly initialize the base plugin instance.
	if err := os.Setenv("PLUGIN_CONFIG", "./testdata/config.yml"); err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("PLUGIN_CONFIG")

	// Create the base plugin
	plugin, err := NewSnmpBasePlugin(&PluginMetadata{
		Name:        "test",
		Maintainer:  "test",
		Description: "test",
		VCS:         "test",
	})

	assert.NoError(t, err)
	assert.NotNil(t, plugin)
}

func TestNewSnmpBasePlugin_NoConfig(t *testing.T) {
	plugin, err := NewSnmpBasePlugin(&PluginMetadata{
		Name:        "test",
		Maintainer:  "test",
		Description: "test",
		VCS:         "test",
	})

	assert.Error(t, err)
	assert.Nil(t, plugin)
}

func TestNewSnmpBasePlugin_NoName(t *testing.T) {
	plugin, err := NewSnmpBasePlugin(&PluginMetadata{
		Maintainer:  "test",
		Description: "test",
		VCS:         "test",
	})

	assert.Error(t, err)
	assert.Equal(t, ErrNoName, err)
	assert.Nil(t, plugin)
}

func TestNewSnmpBasePlugin_NoMaintainer(t *testing.T) {
	plugin, err := NewSnmpBasePlugin(&PluginMetadata{
		Name:        "test",
		Description: "test",
		VCS:         "test",
	})

	assert.Error(t, err)
	assert.Equal(t, ErrNoMaintainer, err)
	assert.Nil(t, plugin)
}
