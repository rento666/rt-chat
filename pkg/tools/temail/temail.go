package temail

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"rt-chat/pkg/tools/ttrans"
	"strconv"
	"time"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	host := viper.GetString("email.host")
	port, _ := strconv.Atoi(viper.GetString("email.port"))
	username := viper.GetString("email.username")
	password := viper.GetString("email.password")

	m := gomail.NewMessage()
	m.SetHeader("From", username)   // 发送方的邮箱地址
	m.SetHeader("To", to)           // 接收方的邮箱地址
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(host, port, username, password) // SMTP服务器的地址、端口、发送方的邮箱地址、发送方的邮箱密码

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}
	return nil
}

// SendVerificationCode 利用邮箱发送验证码
func SendVerificationCode(email, code string, expire time.Duration) error {
	minute := strconv.FormatInt(int64(ttrans.DurationToMinute(expire)), 10)
	if minute == strconv.Itoa(0) {
		minute = "30"
	}
	return SendEmail(email, "Verification Code", "你的验证码是: "+code+"，有效期为30分钟，剩余时间"+minute+"分钟！")
}
