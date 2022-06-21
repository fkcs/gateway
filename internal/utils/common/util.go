package common

import (
	"fmt"
)

const (
	UserPrefix  = "gateway:lock:user"
	CompPrefix  = "gateway:lock:comp"
	TokenPrefix = "gateway:token"
	SignPrefix  = "gateway:sign"
)

func GetRedisUserKey(userId string) string {
	return fmt.Sprintf("%s:%s", UserPrefix, userId)
}

func GetRedisCompKey(component string) string {
	return fmt.Sprintf("%s:%s", CompPrefix, component)
}

func GetRedisTokenKey(token string) string {
	return fmt.Sprintf("%s:%s", TokenPrefix, token)
}

func GetRedisSignKey(userId string) string {
	return fmt.Sprintf("%s:%s", SignPrefix, userId)
}
