package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"gin-blog/config"
	"log/slog"
	"strings"

	"github.com/thanhpk/randstr"
)

type EmailData struct {
	URL      string // 验证链接
	UserName string // 用户名即邮箱地址
	Subject  string // 邮件主题
}

// Format 将邮箱地址转换成小写，并去除空格
func Format(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// GetCode 生成随机字符串
func GetCode() string {
	code := randstr.String(24)
	return code
}

// Encode 生成 base64 编码
func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return data
}

// Decode 解码 base64
func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", errors.New("emailVertify failed, decode error!!!")
	}
	return string(data), nil
}

// GenEmailVerificationInfo 生成加密后的 base64 字符串
func GenEmailVerificationInfo(email, password string) string {
	code := GetCode()
	info := Encode(email + "|" + password + "|" + code)
	return info
}

// ParseEmailVerificationInfo 返回解析 base64字符串后的 邮箱地址 和 code
func ParseEmailVerificationInfo(info string) (string, string, error) {
	data, err := Decode(info)
	if err != nil {
		return "", "", err
	}

	str := strings.Split(data, "|")
	if len(str) != 3 {
		return "", "", errors.New("错误的验证信息格式")
	}

	return str[0], str[1], nil
}

// GetEmailVerifyURL 生成验证链接
func GetEmailVerifyURL(info string) string {
	baseURL := config.Conf.Server.Port
	if baseURL[0] == ':' {
		baseURL = fmt.Sprintf("localhost%s", baseURL)
	}
	return fmt.Sprintf("%s/api/email/verify?info=%s", baseURL, info)
}

func GetEmailData(email, info string) *EmailData {
	return &EmailData{
		URL:      GetEmailVerifyURL(info),
		UserName: email,
		Subject:  "请完成账号注册",
	}
}

func SendEmail(email string, data *EmailData) error {
	// TODO：真实对接邮件 API
	slog.Debug("<发送邮件>")
	slog.Debug("To：" + email)
	slog.Debug("Subject：" + data.Subject)
	slog.Debug("Link：" + data.URL)
	return nil
}
