package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"pnetshieldq/config"
)

type OSSClient struct {
	Client     *oss.Client
	BucketName string
}

func NewOSSClient(cfg config.Config) (*OSSClient, error) {
	client, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccessKeyId, cfg.OSSAccessKeySecret)
	if err != nil {
		return nil, err
	}

	return &OSSClient{
		Client:     client,
		BucketName: cfg.OSSBucketName,
	}, nil
}

func (oc *OSSClient) UploadFile(objectName, filePath string) (string, error) {
	bucket, err := oc.Client.Bucket(oc.BucketName)
	if err != nil {
		return "", err
	}

	err = bucket.PutObjectFromFile(objectName, filePath)
	if err != nil {
		return "", err
	}

	fileURL := "https://grbqwe.oss-cn-hangzhou.aliyuncs.com/" + objectName // 根据实际情况调整URL拼接方式
	return fileURL, nil
}

func (oc *OSSClient) DeleteFile(objectName string) error {
	bucket, err := oc.Client.Bucket(oc.BucketName)
	if err != nil {
		return err
	}

	return bucket.DeleteObject(objectName)
}
func (oc *OSSClient) GetFile(objectName string) ([]byte, error) {
	bucket, err := oc.Client.Bucket(oc.BucketName)
	if err != nil {
		return nil, err
	}

	// 获取文件内容
	body, err := bucket.GetObject(objectName)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// 读取内容到字节数组
	fileContent, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}
