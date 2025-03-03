package utils

import (
	"math/rand"
	"time"
)

// 随机生成一个字符串
func GenerateRandomNickname() string {
	// 创建一个新的随机数生成器实例
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成一个随机字符串作为昵称后缀
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 3 // 你可以根据需要调整长度
	var nickname []byte
	for i := 0; i < length; i++ {
		nickname = append(nickname, charset[rng.Intn(len(charset))])
	}

	return "Client_" + string(nickname)
}
