package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BrokerIP             string `json:"broker_ip"`
	BrokerPort           string `json:"broker_port"`
	DBUser               string `json:"db_user"`
	DBPassword           string `json:"db_password"`
	MQTTUserTable        string `json:"mqtt_user_table"`
	MQTTACLTable         string `json:"mqtt_acl_table"`
	DBHost               string `json:"db_host"`
	DBPort               string `json:"db_port"`
	DBName               string `json:"db_name"`
	EmqxIp               string `json:"emqx_ip"`
	EmqxManagementPort   string `json:"emqx_management_port"`
	EmqxApiPort          string `json:"emqx_api_port"`
	EmqxManagementApiKey string `json:"emqx_management_api_key"`
	ACLName              string `json:"acl_name"`
	ACLPassword          string `json:"acl_password"`
	OSSAccessKeyId       string `json:"oss_access_key_id"`
	OSSAccessKeySecret   string `json:"oss_access_key_secret"`
	OSSBucketName        string `json:"oss_bucket_name"`
	OSSEndpoint          string `json:"oss_endpoint"`
	OSSRegion            string `json:"oss_region"`
}

type OSSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
}

func LoadConfig(file string) (Config, error) {
	var config Config

	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return config, err // 确保解析错误能被返回
	}

	return config, nil
}

func LoadOSSConfig() (*OSSConfig, error) {
	file, err := os.Open("./config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg OSSConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
