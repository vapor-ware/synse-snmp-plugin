package mibs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

func tagEquals(t *testing.T, t1, t2 *sdk.Tag) {
	assert.Equal(t, t1.Namespace, t2.Namespace)
	assert.Equal(t, t1.Annotation, t2.Annotation)
	assert.Equal(t, t1.Label, t2.Label)
}

func TestSnmpDevice_String(t *testing.T) {
	d := SnmpDevice{
		OID:     "1.2.3.4",
		Info:    "info",
		Type:    "type",
		Handler: "handler",
		Output:  "output",
	}

	str := d.String()
	assert.Equal(t, "[SnmpDevice 1.2.3.4: info]", str)
}

func TestSnmpDevice_String_Empty(t *testing.T) {
	d := SnmpDevice{}

	str := d.String()
	assert.Equal(t, "[SnmpDevice : ]", str)
}

func TestSnmpDevice_ToDevice(t *testing.T) {
	d := SnmpDevice{
		OID:     "1.2.3.4",
		Info:    "testDevice",
		Type:    "temperature",
		Handler: "temperature",
		Output:  "temperature",
		Tags: []*sdk.Tag{
			core.TagOrPanic("vaporio/test:device"),
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
		Context: map[string]string{
			"abc": "123",
		},
		Alias:        "test-device",
		WriteTimeout: 1 * time.Second,
	}

	dev, err := d.ToDevice()
	assert.NoError(t, err)
	assert.NotNil(t, dev)

	assert.Equal(t, "temperature", dev.Type)
	assert.Equal(t, "testDevice", dev.Info)
	assert.Equal(t, "temperature", dev.Handler)
	assert.Equal(t, "test-device", dev.Alias)
	assert.Equal(t, "temperature", dev.Output)
	assert.Equal(t, 1*time.Second, dev.WriteTimeout)
	assert.Equal(t, []sdk.Transformer(nil), dev.Transforms)
	assert.Equal(t, map[string]interface{}{
		"foo": "bar",
		"oid": "1.2.3.4",
	}, dev.Data)
	assert.Equal(t, map[string]string{
		"abc": "123",
		"oid": "1.2.3.4",
	}, dev.Context)

	assert.Len(t, dev.Tags, 4)
	tagEquals(t, dev.Tags[0], &sdk.Tag{Namespace: "protocol", Annotation: "", Label: "snmp"})
	tagEquals(t, dev.Tags[1], &sdk.Tag{Namespace: "snmp", Annotation: "oid", Label: "1.2.3.4"})
	tagEquals(t, dev.Tags[2], &sdk.Tag{Namespace: "snmp", Annotation: "name", Label: "testDevice"})
	tagEquals(t, dev.Tags[3], &sdk.Tag{Namespace: "vaporio", Annotation: "test", Label: "device"})
}

func TestSnmpDevice_ToDevice2(t *testing.T) {
	// Same test, different configuration.
	d := SnmpDevice{
		OID:     "1.2.3.4",
		Info:    "info with spaces",
		Type:    "state",
		Handler: "state",
		Output:  "state",
	}

	dev, err := d.ToDevice()
	assert.NoError(t, err)
	assert.NotNil(t, dev)

	assert.Equal(t, "state", dev.Type)
	assert.Equal(t, "info with spaces", dev.Info)
	assert.Equal(t, "state", dev.Handler)
	assert.Equal(t, "", dev.Alias)
	assert.Equal(t, "state", dev.Output)
	assert.Equal(t, 0*time.Second, dev.WriteTimeout)
	assert.Equal(t, []sdk.Transformer(nil), dev.Transforms)
	assert.Equal(t, map[string]interface{}{
		"oid": "1.2.3.4",
	}, dev.Data)
	assert.Equal(t, map[string]string{
		"oid": "1.2.3.4",
	}, dev.Context)

	assert.Len(t, dev.Tags, 3)
	tagEquals(t, dev.Tags[0], &sdk.Tag{Namespace: "protocol", Annotation: "", Label: "snmp"})
	tagEquals(t, dev.Tags[1], &sdk.Tag{Namespace: "snmp", Annotation: "oid", Label: "1.2.3.4"})
	tagEquals(t, dev.Tags[2], &sdk.Tag{Namespace: "snmp", Annotation: "name", Label: "infoWithSpaces"})
}

func TestSnmpDevice_ToDevice_BadOutput(t *testing.T) {
	d := SnmpDevice{
		OID:     "1.2.3.4",
		Info:    "testDevice",
		Type:    "temperature",
		Handler: "temperature",
		Output:  "unknown-output",
	}

	dev, err := d.ToDevice()
	assert.Error(t, err)
	assert.Nil(t, dev)
}
