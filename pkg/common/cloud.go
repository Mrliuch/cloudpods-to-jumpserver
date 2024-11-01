package common

type Notice struct {
	Action          string          `json:"action"`        // create/delete
	ResourceType    string          `json:"resource_type"` // server
	Result          string          `json:"result"`        // succeed
	ResourceDetails ResourceDetails `json:"resource_details"`
}

type ResourceDetails struct {
	Host         string `json:"host"` // liuchen-hp-192-168-10-4
	HostID       string `json:"host_id"`
	Hypervisor   string `json:"hypervisor"`    // kvm
	IPS          string `json:"ips"`           // 192.168.10.x
	LoginAccount string `json:"login_account"` // root
	Password     string `json:"password"`      // passwd
	Name         string `json:"name"`          // master2-239
	OSType       string `json:"os_type"`       // Linux
	JMPID        string `json:"jmpid"`         // JumpServer返回ID
}
