package hander

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"regexp"
	"testProject/test/config"
	"testProject/test/dao"
	"testProject/test/middleware"
	"testProject/test/models"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	UserName   string `json:"userName" binding:"required"`
	PassWord   string `json:"passWord" binding:"required"`
	RePassWord string `json:"rePassWord" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
}
type RegisterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 注册
func RegisterUser(c *gin.Context) {
	var response RegisterResponse
	var request RegisterRequest
	user := &models.User{}
	c.ShouldBind(&request)
	err := validateInput(request.UserName, request.PassWord, request.RePassWord, request.Email, request.Phone)
	if err != nil {
		response.Code = 400
		response.Message = err.Error()
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	user1, _ := dao.GetUserByUsername(config.DB, request.UserName)
	if user1 != nil {
		response.Code = 400
		response.Message = "该用户已存在"
		c.JSON(http.StatusOK, response)
		return
	}
	user.UserName = request.UserName
	user.Email = request.Email
	user.Salt, _ = generateSalt()
	user.PassWord, _ = encryptPassword(request.PassWord, user.Salt)
	user.Phone = request.Phone
	err2 := dao.CreateUser(config.DB, user)
	if err2 != nil {
		response.Code = 400
		response.Message = err2.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	response.Code = 200
	response.Message = "注册成功"
	c.JSON(http.StatusOK, response)
	return
}

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	PassWord string `json:"passWord" binding:"required"`
}
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string
}

// 登录
func Login(c *gin.Context) {
	var response LoginResponse
	var request LoginRequest
	c.ShouldBind(&request)
	err := validateInput(request.UserName, request.PassWord, request.PassWord, "1111111@qq.com", "13914444444")
	if err != nil {
		response.Code = 400
		response.Message = err.Error()
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	user, _ := dao.GetUserByUsername(config.DB, request.UserName)
	if user == nil {
		response.Code = 400
		response.Message = "该用户不存在"
		c.JSON(http.StatusOK, response)
		return
	}
	if !verifyPassword(request.PassWord, user.Salt, user.PassWord) {
		response.Code = 400
		response.Message = "密码错误"
		c.JSON(http.StatusOK, response)
		return
	}

	token := middleware.GenerateToken(int(user.ID), user.UserName)
	response.Code = 200
	response.Message = "登录成功"
	response.Token = token
	c.JSON(http.StatusOK, response)
	return
}

type ResetRequest struct {
	UserName    string `json:"userName" binding:"required"`
	PassWord    string `json:"passWord" binding:"required"`
	NewPassWord string `json:"newPassWord" binding:"required"`
}
type ResetResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 修改密码
func ResetPassword(c *gin.Context) {
	var response ResetResponse
	var request ResetRequest
	c.ShouldBind(&request)
	err := validateInput(request.UserName, request.PassWord, request.PassWord, "1111111@qq.com", "13914444444")
	if request.PassWord == request.NewPassWord {
		response.Code = 400
		response.Message = "新密码不能和原密码一致"
		c.JSON(http.StatusOK, response)
		return
	}
	if len(request.NewPassWord) < 8 || !containsDigitAndLetter(request.NewPassWord) {
		response.Code = 400
		response.Message = "密码必须至少包含8位字符和数字"
		c.JSON(http.StatusOK, response)
		return
	}
	if err != nil {
		response.Code = 400
		response.Message = err.Error()
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	user, _ := dao.GetUserByUsername(config.DB, request.UserName)
	if user == nil {
		response.Code = 400
		response.Message = "该用户不存在"
		c.JSON(http.StatusOK, response)
		return
	}
	if !verifyPassword(request.PassWord, user.Salt, user.PassWord) {
		response.Code = 400
		response.Message = "密码错误"
		c.JSON(http.StatusOK, response)
		return
	}
	passwrod, _ := encryptPassword(request.NewPassWord, user.Salt)
	err2 := dao.UpdatePassword(config.DB, passwrod, int(user.ID))
	if err2 == nil {
		response.Code = 200
		response.Message = "修改成功"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Code = 400
	response.Message = err2.Error()
	c.JSON(http.StatusOK, response)
	return
}

// 验证输入参数合法性
func validateInput(username, password, repassword, email, phone string) error {
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
	if password != repassword {
		return errors.New("两次密码不一致")
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

// 生成随机盐值
func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// 加密密码
func encryptPassword(password string, salt string) (string, error) {
	// 将密码和盐值拼接在一起
	saltedPassword := password + salt

	// 使用SHA256哈希算法对拼接后的密码进行加密
	hash := sha256.Sum256([]byte(saltedPassword))

	// 将加密后的结果转换为十六进制字符串
	encryptedPassword := hex.EncodeToString(hash[:])

	return encryptedPassword, nil
}

// 验证密码是否正确
func verifyPassword(password string, salt string, hashedPassword string) bool {
	// 对输入的密码进行加密
	encryptedPassword, _ := encryptPassword(password, salt)

	// 比较加密后的密码与已存储的哈希密码是否一致
	return encryptedPassword == hashedPassword
}
