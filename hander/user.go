package hander

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// 注册
func RegisterUser(c *gin.Context) {
	userName := c.PostForm("UserName")
	passWord := c.PostForm("PassWord")
	rePassWord := c.PostForm("RePassWord")
	email := c.PostForm("Email")
	phone := c.PostForm("Phone")
	err := validateInput(userName, passWord, email, phone)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": err,
		})
		return
	}
	if passWord != rePassWord {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "两次密码不一致",
		})
		return
	}
}

func validateInput(username, password, email, phone string) error {
	// 验证用户名不为空
	if username == "" {
		return errors.New("用户名不能为空")
	}

	// 验证密码不少于8位且包含字符和数字
	if len(password) < 8 || !containsDigitAndLetter(password) {
		return errors.New("密码必须至少包含8位字符和数字")
	}

	// 验证邮箱格式
	if !isValidEmail(email) {
		return errors.New("邮箱地址格式不正确")
	}

	// 验证手机号格式
	if !isValidPhone(phone) {
		return errors.New("手机号格式不正确")
	}

	return nil
}

func containsDigitAndLetter(s string) bool {
	hasDigit := false
	hasLetter := false
	for _, c := range s {
		if '0' <= c && c <= '9' {
			hasDigit = true
		}
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
			hasLetter = true
		}
	}
	return hasDigit && hasLetter
}

func isValidEmail(email string) bool {
	// 使用正则表达式验证邮箱格式
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

func isValidPhone(phone string) bool {
	// 使用正则表达式验证手机号格式（简化版）
	phonePattern := `^\d{11}$`
	return regexp.MustCompile(phonePattern).MatchString(phone)
}
