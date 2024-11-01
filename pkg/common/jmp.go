package common

type Asset struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Nodes     []string   `json:"nodes"`
	Protocols []Protocol `json:"protocols"`
	IsActive  bool       `json:"is_active"`
	Platform  string     `json:"platform"`
	Accounts  []Accounts `json:"accounts"`
}

type Protocol struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type Accounts struct {
	Name       string `json:"name" default:"root"`
	Username   string `json:"username" default:"root"`
	Secret     string `json:"secret"`
	SecretType string `json:"secret_type" default:"password"`
	IsActive   bool   `json:"is_active" default:"true"`
	Privileged bool   `json:"privileged" default:"true"`
}

type Platform struct {
	PK int `json:"pk"`
}

// // 创建资产jumpserver返回数据
type RespJMP struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
