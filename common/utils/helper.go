package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// CompareHashAndPassword
//
//	@Description: 密码比较
//	@param e
//	@param p
//	@return bool
//	@return error
func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Encrypt
//
//	@Description: 加密
//	@param password
//	@return err
func Encrypt(password string) (err error, pwd string) {
	if password == "" {
		return
	}
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		pwd = string(hash)
		return
	}
}
