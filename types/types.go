package types

type Configurations struct {
	IBMQConfigDetails IBMQConfigDetails `validate:"required" json:"MQConfigDetails"`
}

type IBMQConfigDetails struct {
	MQIp           string `validate:"required" json:"mq_ip"`
	MQPort         int    `validate:"required" json:"mq_port"`
	MQConnection   string `validate:"required" json:"mq_connection_name"`
	MQManager      string `validate:"required" json:"mq_manager"`
	MQChannel      string `validate:"required" json:"mq_channel"`
	MQQueue        string `validate:"required" json:"mq_queue"`
	MQTLSCipher    string `validate:"required" json:"mq_tls_cipher"`
	MQRepoLoc      string `validate:"required" json:"mq_repo_loc"`
	SyncPointFLag  bool   `json:"sync_point_flag"`
	MQWaitInterval int32  `validate:"required" json:"mq_wait_interval"`
	MQBufferBytes  int    `validate:"required" json:"mq_buffer"`
}
