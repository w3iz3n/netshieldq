package utils

import (
	"encoding/base64"
	"io/ioutil"
)

func GetDefaultAvatar() ([]byte, error) {
	defaultAvatar, err := ioutil.ReadFile("E:\\netshieldq\\pnetshieldq\\assets\\default_avatar.png")
	if err != nil {
		return nil, err
	}
	return defaultAvatar, nil
}

func GetDefaultAvatarBase64() (string, error) {
	defaultAvatar, err := GetDefaultAvatar()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(defaultAvatar), nil
}
