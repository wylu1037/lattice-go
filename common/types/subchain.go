package types

const (
	// SubchainStatusRUNNING 子链（通道）运行中
	SubchainStatusRUNNING = "running"
	// SubchainStatusSTOP 子链（通道）已停止运行
	SubchainStatusSTOP = "stop"
)

// SubchainRunningStatus 子链的运行状态
//   - Status SubchainStatusRUNNING or SubchainStatusSTOP
//   - Running 是否在运行中，false:已停止 true:运行中
type SubchainRunningStatus struct {
	Status  string `json:"status"`
	Running bool   `json:"running"`
}
