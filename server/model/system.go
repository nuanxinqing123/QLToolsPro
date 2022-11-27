package model

// RemoteVersion 系统远程信息
type RemoteVersion struct {
	LocVersion string
	Version    string
	Notice     string
	Update     string
}

// SoftWareGOOS 更新版本架构
type SoftWareGOOS struct {
	Framework string `json:"framework" binding:"required"`
}
