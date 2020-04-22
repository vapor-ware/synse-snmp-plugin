package mibs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetGlobalMibs() {
	pluginMibs = map[string]*MIB{}
}

func TestRegister(t *testing.T) {
	defer resetGlobalMibs()
	assert.Empty(t, pluginMibs)

	mib1 := &MIB{Name: "mib-1"}
	mib2 := &MIB{Name: "mib-2"}

	err := Register(mib1, mib2)
	assert.NoError(t, err)

	assert.Len(t, pluginMibs, 2)
	assert.NotNil(t, pluginMibs["mib-1"])
	assert.NotNil(t, pluginMibs["mib-2"])
}

func TestRegister_NilMib(t *testing.T) {
	defer resetGlobalMibs()
	assert.Empty(t, pluginMibs)

	mib1 := &MIB{Name: "mib-1"}

	err := Register(mib1, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrNilMib, err)

	assert.Len(t, pluginMibs, 1) // the first non-nil mib was added
	assert.NotNil(t, pluginMibs["mib-1"])
}

func TestRegister_MibExists(t *testing.T) {
	defer resetGlobalMibs()
	assert.Empty(t, pluginMibs)

	mib1 := &MIB{Name: "mib-1"}
	mib2 := &MIB{Name: "mib-1"} // same name

	err := Register(mib1, mib2)
	assert.Error(t, err)
	assert.Equal(t, ErrMibExists, err)

	assert.Len(t, pluginMibs, 1) // the first non-duplicate mib was added
	assert.NotNil(t, pluginMibs["mib-1"])
}

func TestGet_Exists(t *testing.T) {
	defer resetGlobalMibs()

	pluginMibs["test-mib"] = &MIB{Name: "test-mib"}

	mib := Get("test-mib")
	assert.NotNil(t, mib)
	assert.Equal(t, pluginMibs["test-mib"], mib)
}

func TestGet_NotExists(t *testing.T) {
	defer resetGlobalMibs()

	mib := Get("test-mib")
	assert.Nil(t, mib)
}

func TestGetAll(t *testing.T) {
	defer resetGlobalMibs()

	mib1 := &MIB{Name: "mib-1"}
	mib2 := &MIB{Name: "mib-2"}
	mib3 := &MIB{Name: "mib-3"}

	pluginMibs["mib-1"] = mib1
	pluginMibs["mib-2"] = mib2
	pluginMibs["mib-3"] = mib3

	mibs := GetAll()
	assert.Len(t, mibs, 3)

	expected := []string{"mib-1", "mib-2", "mib-3"}
	for _, m := range mibs {
		assert.Contains(t, expected, m.Name)
	}
}

func TestGetAll_Empty(t *testing.T) {
	defer resetGlobalMibs()

	mibs := GetAll()
	assert.Empty(t, mibs)
}
