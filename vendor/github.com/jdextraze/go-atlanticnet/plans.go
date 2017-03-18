package atlanticnet

import "log"

type Plan struct {
	Bandwidth        int     `json:"bandwidth"`
	CentOsCapable    string  `json:"centos_capable"`
	CPanelCapable    string  `json:"cpanel_capable"`
	Disk             string  `json:"disk"`
	DisplayBandwidth string  `json:"display_bandwidth"`
	DisplayDisk      string  `json:"display_disk"`
	DisplayRam       string  `json:"display_ram"`
	FreeTransfer     string  `json:"free_transfer"`
	NumCpu           string  `json:"num_cpu"`
	OsType           string  `json:"ostype"`
	PlanName         string  `json:"plan_name"`
	Platform         string  `json:"platform"`
	Ram              string  `json:"ram"`
	RatePerHour      float64 `json:"rate_per_hr"`
}

type describePlanResponse struct {
	Plans     map[string]Plan
	RequestId string `json:"requestid"`
}

type describePlanDTO struct {
	BaseResponse
	Response describePlanResponse `json:"describe-planresponse"`
}

func (c *client) DescribePlan(name, platform string) ([]Plan, error) {
	in := map[string]string{}
	if name != "" {
		in["plan_name"] = name
	}
	if platform != "" {
		in["platform"] = platform
	}

	out := describePlanDTO{}

	err := c.doRequest("describe-plan", in, &out)
	if err != nil {
		return nil, err
	} else if out.Error != nil {
		return nil, out.Error
	} else if c.debug {
		log.Println(out)
	}

	plans := make([]Plan, len(out.Response.Plans))
	i := 0
	for _, plan := range out.Response.Plans {
		plans[i] = plan
		i += 1
	}
	return plans, nil
}
