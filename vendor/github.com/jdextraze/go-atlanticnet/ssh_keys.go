package atlanticnet

type SshKey struct {
	Id        string `json:"key_id"`
	Name      string `json:"key_name"`
	PublicKey string `json:"public_key"`
}

type listSshKeysResponse struct {
	KeysSet   map[string]SshKey `json:"KeysSet"`
	RequestId string            `json:"requestid"`
}

type listSshKeysDTO struct {
	BaseResponse
	Response listSshKeysResponse `json:"list-sshkeysresponse"`
}

func (c *client) ListSshKeys() ([]SshKey, error) {
	out := listSshKeysDTO{}
	err := c.doRequest("list-sshkeys", nil, &out)
	if err != nil {
		return []SshKey{}, err
	} else if out.Error != nil {
		return []SshKey{}, out.Error
	}

	sshKeys := make([]SshKey, len(out.Response.KeysSet))
	i := 0
	for _, sshKey := range out.Response.KeysSet {
		sshKeys[i] = sshKey
		i += 1
	}

	return sshKeys, nil
}
