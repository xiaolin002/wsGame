package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

// TestHashPassword 测试 HashPassword 函数
func TestHashPassword(t *testing.T) {
	password := "1"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword 失败: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword 返回空字符串")
	}
	log.Printf("hash: %s", hash)
}

// TestCheckPasswordHash 测试 CheckPasswordHash 函数
func TestCheckPasswordHash(t *testing.T) {
	password := "1"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword 失败: %v", err)
	}
	log.Printf("hash: %s", hash)

	// 测试正确密码
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash 验证正确密码失败")
	}

	// 测试错误密码
	wrongPassword := "wrongpassword"
	if CheckPasswordHash(wrongPassword, hash) {
		t.Error("CheckPasswordHash 验证错误密码通过")
	}
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
