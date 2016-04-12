package driver

import (
	"fmt"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
	"github.com/jdextraze/go-atlanticnet"
	"io"
	"os"
	"os/user"
	"strconv"
)

type Driver struct {
	*drivers.BaseDriver
	APIKey         string
	APISecret      string
	OrigSSHKeyPath string

	ImageID    string
	PlanName   string
	VmLocation string

	InstanceID string
	SSHKeyId   string

	client atlanticnet.Client
}

const (
	defaultImageId    = "ubuntu-14.04_64bit"
	defaultPlanName   = "XS"
	defaultVmLocation = "USWEST1"
	SSHUser           = "root"
	SSHPort           = 22
)

var (
	defaultSshKeyPath  string
	atlanticNetRegions = [...]string{
		"USEAST1",
		"USEAST2",
		"USCENTRAL1",
		"USWEST1",
		"CAEAST1",
		"EUWEST1",
	}
)

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	defaultSshKeyPath = user.HomeDir + "/.ssh/id_rsa"
}

func NewDriver(hostName, storePath string) *Driver {
	d := &Driver{
		ImageID:    defaultImageId,
		PlanName:   defaultPlanName,
		VmLocation: defaultVmLocation,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
			SSHUser:     SSHUser,
			SSHPort:     SSHPort,
		},
	}
	return d
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_API_KEY",
			Name:   "atlantic-net-api-key",
			Usage:  "Atlantic.Net API key",
		},
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_API_SECRET",
			Name:   "atlantic-net-api-secret",
			Usage:  "Atlantic.Net API secret",
		},
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_SSH_KEY_ID",
			Name:   "atlantic-net-ssh-key-id",
			Usage:  "Atlantic.Net SSH key id",
		},
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_SSH_KEY_PATH",
			Name:   "atlantic-net-ssh-key-path",
			Usage:  "Atlantic.Net SSH key path",
			Value:  defaultSshKeyPath,
		},
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_IMAGE_ID",
			Name:   "atlantic-net-image-id",
			Usage:  "Atlantic.Net image id",
			Value:  defaultImageId,
		},
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_PLAN_NAME",
			Name:   "atlantic-net-plan-name",
			Usage:  "Atlantic.Net plan name",
			Value:  defaultPlanName,
		},
		mcnflag.StringFlag{
			EnvVar: "ATLANTIC_NET_VM_LOCATION",
			Name:   "atlantic-net-vm-location",
			Usage:  "Atlantic.Net vm location",
			Value:  defaultVmLocation,
		},
	}
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) DriverName() string {
	return "atlanticnet"
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIKey = flags.String("atlantic-net-api-key")
	d.APISecret = flags.String("atlantic-net-api-secret")
	d.SSHKeyId = flags.String("atlantic-net-ssh-key-id")
	d.OrigSSHKeyPath = flags.String("atlantic-net-ssh-key-path")
	d.ImageID = flags.String("atlantic-net-image-id")
	d.VmLocation = flags.String("atlantic-net-vm-location")
	d.PlanName = flags.String("atlantic-net-plan-name")
	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")

	if d.APIKey == "" {
		return fmt.Errorf("Atlantic.Net driver requires the --atlantic-net-api-key option")
	}
	if d.APISecret == "" {
		return fmt.Errorf("Atlantic.Net driver requires the --atlantic-net-api-secret option")
	}
	if d.SSHKeyId == "" {
		return fmt.Errorf("Atlantic.Net driver requires the --atlantic-net-ssh-key-id option")
	}
	return nil
}

func (d *Driver) PreCreateCheck() error {
	log.Info("Validating Atlantic.Net VPS parameters...")

	if err := d.validateSshKey(); err != nil {
		return err
	}

	if err := d.validateVmLocation(); err != nil {
		return err
	}

	if err := d.validatePlan(); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Create() error {
	log.Info("Creating Atlantic.Net VPS...")

	if err := d.copySSHKey(); err != nil {
		return err
	}

	instance, err := d.getClient().RunInstance(atlanticnet.RunInstanceRequest{
		ServerName: d.MachineName,
		ImageId:    d.ImageID,
		PlanName:   d.PlanName,
		VMLocation: d.VmLocation,
		KeyId:      d.SSHKeyId,
	})
	if err != nil {
		return err
	}
	d.InstanceID = strconv.Itoa(instance[0].Id)
	d.IPAddress = instance[0].IpAddress

	log.Infof("Created Atlantic.Net VPS ID: %s, Public IP: %s",
		d.InstanceID,
		d.IPAddress,
	)

	return nil
}

func (d *Driver) GetURL() (string, error) {
	s, err := d.GetState()
	if err != nil {
		return "", err
	}

	if s != state.Running {
		return "", drivers.ErrHostIsNotRunning
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" || d.IPAddress == "0" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}

func (d *Driver) GetState() (state.State, error) {
	machine, err := d.getClient().DescribeInstance(d.InstanceID)
	if err != nil {
		return state.Error, err
	}
	switch machine.VmStatus {
	case atlanticnet.StatusAwaitingCreation, atlanticnet.StatusCreating, atlanticnet.StatusRestarting:
		return state.Starting, nil
	case atlanticnet.StatusStopped:
		return state.Stopped, nil
	case atlanticnet.StatusRunning:
		return state.Running, nil
	}
	return state.Error, nil
}

func (d *Driver) Start() error {
	return fmt.Errorf("Atlantic.Net doesn`t support this. Please restart the machine instead.")
}

func (d *Driver) Stop() error {
	return fmt.Errorf("Atlantic.Net doesn`t support this. Please restart the machine instead.")
}

func (d *Driver) Remove() error {
	client := d.getClient()
	log.Debugf("removing %s", d.MachineName)
	terminatedInstances, err := client.TerminateInstance(d.InstanceID)
	if err != nil {
		return err
	}
	for _, v := range terminatedInstances {
		if v.Id == d.InstanceID && v.Result == "true" {
			return nil
		}
	}
	return fmt.Errorf("Error removing instance %s", d.InstanceID)
}

func (d *Driver) Restart() error {
	if vmState, err := d.GetState(); err != nil {
		return err
	} else if vmState == state.Starting {
		log.Info("Host is already starting")
		return nil
	}
	log.Debugf("restarting %s", d.MachineName)
	rebootedInstance, err := d.getClient().RebootInstance(d.InstanceID, atlanticnet.RebootTypeSoft)
	if err != nil {
		return err
	}
	if rebootedInstance.Value == "true" {
		return nil
	}
	return fmt.Errorf("Error rebooting instance %s", d.InstanceID)
}

func (d *Driver) Kill() error {
	return fmt.Errorf("Atlantic.Net doesn`t support this. Please restart the machine instead.")
}

func (d *Driver) getClient() atlanticnet.Client {
	log.Debug("getting client")
	if d.client == nil {
		d.client = atlanticnet.NewClient(d.APIKey, d.APISecret, false)
	}
	return d.client
}

func (d *Driver) validateSshKey() error {
	sshKeys, err := d.getClient().ListSshKeys()
	if err != nil {
		return err
	}
	for _, sshKey := range sshKeys {
		if sshKey.Id == d.SSHKeyId {
			return nil
		}
	}
	return fmt.Errorf("Ssh Key Id %s is invalid", d.SSHKeyId)
}

func (d *Driver) validateVmLocation() error {
	for _, region := range atlanticNetRegions {
		if region == d.VmLocation {
			return nil
		}
	}
	return fmt.Errorf("VM location %s is invalid", d.VmLocation)
}

func (d *Driver) validatePlan() error {
	plans, err := d.getClient().DescribePlan("", "linux")
	if err != nil {
		return err
	}
	for _, plan := range plans {
		if plan.PlanName == d.PlanName {
			return nil
		}
	}
	return fmt.Errorf("Plan name %s is invalid", d.PlanName)
}

func (d *Driver) copySSHKey() (err error) {
	in, err := os.Open(d.OrigSSHKeyPath)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(d.GetSSHKeyPath())
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
