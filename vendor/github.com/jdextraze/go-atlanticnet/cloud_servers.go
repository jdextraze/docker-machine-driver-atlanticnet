package atlanticnet

import (
	"log"
	"strconv"
)

type Instance struct {
	Id                 string `json:"InstanceId"`
	CuId               string `json:"cu_id"`
	RatePerHour        string `json:"rate_per_hr"`
	VmBandwith         string `json:"vm_bandwith"`
	VmCpuReq           string `json:"vm_cpu_req"`
	VmCreatedDate      string `json:"vm_created_date"`
	VmDescription      string `json:"vm_description"`
	VmDiskReq          string `json:"vm_disk_req"`
	VmImage            string `json:"vm_image"`
	VmImageDisplayName string `json:"vm_image_display_name"`
	VmIpAddress        string `json:"vm_ip_address"`
	VmName             string `json:"vm_name"`
	VmNetworkReq       string `json:"vm_network_req"`
	VmOsArchitecture   string `json:"vm_os_architecture"`
	VmPlanName         string `json:"vm_plan_name"`
	VmRamReq           string `json:"vm_ram_req"`
	VmStatus           Status `json:"vm_status"`
}

type listInstancesResponse struct {
	InstancesSet map[string]Instance `json:"instancesSet"`
	RequestId    string              `json:"requestid"`
}

type listInstancesDTO struct {
	BaseResponse
	Response listInstancesResponse `json:"list-instancesresponse"`
}

func (c *client) ListInstances() ([]Instance, error) {
	out := listInstancesDTO{}
	err := c.doRequest("list-instances", nil, &out)
	if err != nil {
		return nil, err
	} else if out.Error != nil {
		return nil, out.Error
	} else if c.debug {
		log.Println(out)
	}

	instances := make([]Instance, len(out.Response.InstancesSet))
	i := 0
	for _, plan := range out.Response.InstancesSet {
		instances[i] = plan
		i += 1
	}
	return instances, nil
}

type RunInstance struct {
	Id        int    `json:"instanceid"`
	IpAddress string `json:"ip_address"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

type runInstanceResponse struct {
	InstancesSet map[string]RunInstance `json:"instancesSet"`
	RequestId    string                 `json:"requestid"`
}

type runInstanceDTO struct {
	BaseResponse
	RunInstanceResponse runInstanceResponse `json:"run-instanceresponse"`
}

type RunInstanceRequest struct {
	ServerName   string
	ImageId      string
	PlanName     string
	VMLocation   string
	EnableBackup bool
	CloneImage   string
	ServerQty    int
	KeyId        string
}

func (r RunInstanceRequest) ToMap() map[string]string {
	req := map[string]string{
		"servername":  r.ServerName,
		"imageid":     r.ImageId,
		"planname":    r.PlanName,
		"vm_location": r.VMLocation,
	}
	if r.EnableBackup {
		req["enablebackup"] = "Y"
	}
	if r.CloneImage != "" {
		req["cloneimage"] = r.CloneImage
	}
	if r.ServerQty > 0 {
		req["serverqty"] = strconv.Itoa(r.ServerQty)
	}
	if r.KeyId != "" {
		req["key_id"] = r.KeyId
	}
	return req
}

func (c *client) RunInstance(req RunInstanceRequest) ([]RunInstance, error) {
	out := runInstanceDTO{}
	err := c.doRequest("run-instance", req.ToMap(), &out)
	if err != nil {
		return nil, err
	} else if out.Error != nil {
		return nil, out.Error
	} else if c.debug {
		log.Println(out)
	}

	instances := make([]RunInstance, len(out.RunInstanceResponse.InstancesSet))
	i := 0
	for _, plan := range out.RunInstanceResponse.InstancesSet {
		instances[i] = plan
		i += 1
	}
	return instances, nil
}

type TerminateInstance struct {
	Id      string `json:"InstanceId"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type terminateInstanceResponse struct {
	InstancesSet map[string]TerminateInstance `json:"instancesSet"`
	RequestId    string                       `json:"requestid"`
}

type terminateInstanceDTO struct {
	BaseResponse
	Response terminateInstanceResponse `json:"terminate-instanceresponse"`
}

func (c *client) TerminateInstance(instanceId string) ([]TerminateInstance, error) {
	in := map[string]string{
		"instanceid": instanceId,
	}
	out := terminateInstanceDTO{}
	err := c.doRequest("terminate-instance", in, &out)
	if err != nil {
		return nil, err
	} else if out.Error != nil {
		return nil, out.Error
	} else if c.debug {
		log.Println(out)
	}

	instances := make([]TerminateInstance, len(out.Response.InstancesSet))
	i := 0
	for _, plan := range out.Response.InstancesSet {
		instances[i] = plan
		i += 1
	}
	return instances, nil
}

type InstanceDescription struct {
	Instance
	ClonedFrom                  string      `json:"cloned_from"`
	DisallowDeletion            string      `json:"disallow_deletion"`
	Removed                     string      `json:"removed"`
	ReprovisioningProcessedDate string      `json:"reprovisioning_processed_date"`
	VmId                        string      `json:"vm_id"`
	VmIpGateway                 string      `json:"vm_ip_gateway"`
	VmIpSubnet                  string      `json:"vm_ip_subnet"`
	RemovedDate                 string      `json:"vm_removed_date"`
	VmUsername                  string      `json:"vm_username"`
	VmVncPassword               string      `json:"vm_vnc_password"`
	VncPort                     interface{} `json:"vnc_port"` // TODO can be string or bool
}

type describeInstanceResponse struct {
	InstancesSet map[string]InstanceDescription `json:"instanceSet"`
	RequestId    string                         `json:"requestid"`
}

type describeInstanceDTO struct {
	BaseResponse
	Response describeInstanceResponse `json:"describe-instanceresponse"`
}

func (c *client) DescribeInstance(instanceId string) (*InstanceDescription, error) {
	in := map[string]string{
		"instanceid": instanceId,
	}
	out := describeInstanceDTO{}
	err := c.doRequest("describe-instance", in, &out)
	if err != nil {
		return nil, err
	} else if out.Error != nil {
		return nil, out.Error
	} else if c.debug {
		log.Println(out)
	}

	for _, instanceDescription := range out.Response.InstancesSet {
		return &instanceDescription, nil
	}

	return nil, nil
}

type RebootInstance struct {
	InstanceId string `json:"InstanceId"`
	Message    string `json:"message"`
	Value      string `json:"value"`
}

type rebootInstanceResponse struct {
	Return    map[string]RebootInstance `json:"instancesSet"`
	RequestId string                    `json:"requestid"`
}

type rebootInstanceDTO struct {
	BaseResponse
	Response rebootInstanceResponse `json:"reboot-instanceresponse"`
}

func (c *client) RebootInstance(instanceId string, rebootType RebootType) (*RebootInstance, error) {
	in := map[string]string{"instanceid": instanceId}
	if rebootType != "" {
		in["reboottype"] = string(rebootType)
	}

	out := rebootInstanceDTO{}
	err := c.doRequest("reboot-instance", in, &out)
	if err != nil {
		return nil, err
	} else if out.Error != nil {
		return nil, out.Error
	} else if c.debug {
		log.Println(out)
	}

	for _, rebootInstance := range out.Response.Return {
		return &rebootInstance, nil
	}

	return nil, nil
}
