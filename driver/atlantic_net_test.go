package driver

import (
	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetConfigFromFlags(t *testing.T) {
	driver := NewDriver("default", "path")

	flags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"atlantic-net-api-key":    "API_KEY",
			"atlantic-net-api-secret": "SECRET",
			"atlantic-net-ssh-key-id": "KEY_ID",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(flags)

	assert.NoError(t, err)
	assert.Empty(t, flags.InvalidFlags)
	assert.Equal(t, driver.ResolveStorePath("id_rsa"), driver.GetSSHKeyPath())
}

func TestDefaultSSHUserAndPort(t *testing.T) {
	driver := NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"atlantic-net-api-key":    "API_KEY",
			"atlantic-net-api-secret": "SECRET",
			"atlantic-net-ssh-key-id": "KEY_ID",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)

	sshPort, err := driver.GetSSHPort()
	assert.Equal(t, "root", driver.GetSSHUsername())
	assert.Equal(t, 22, sshPort)
	assert.NoError(t, err)
}
