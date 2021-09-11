// author: s0nnet
// time: 2020-09-01
// desc:

package util

import (
	"adp_backend/common"
	"adp_backend/util/crypto"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

// GenerateToken generates a jwt access token
func GenerateToken(username, role string, pri int32, exp int64) (string, error) {
	claim := make(jwt.MapClaims)
	claim["user"] = username
	claim["role"] = role
	claim["pri"] = pri
	claim["exp"] = exp
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtToken, err := token.SignedString([]byte(common.JWT_SECRET))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

// CheckPassStrength check a password strength
func CheckPassStrength(password string) string {
	high := `^(\w|\W)+{12,}$`
	middle := `^(((\d|[a-zA-Z])+)|(\d|[~!@#$%^&*_])+|([a-zA-Z]|[~!@#$%^&*_])+)$`
	low := `^(?:\d+|[a-zA-Z]+|[~!@#$%^&*_]+)$`

	if ok, err := regexp.MatchString(low, password); err == nil && ok {
		return "low"
	}
	if ok, err := regexp.MatchString(middle, password); err == nil && ok {
		return "middle"
	}
	if ok, err := regexp.MatchString(high, password); err == nil && ok {
		return "high"
	}

	return "low"
}

//handle ldap sting
//example:
//ldapAddr: ldap://DC01.domain02.com
//
//domain: domain02.com
//domainName: domain02
//dcHostName: DC01
//dn: DC=domain02,DC=com
func LDAPParse(ldapAddr string) (domain, domainName, dcHostName, dn string, err error) {
	ldap, err := url.Parse(ldapAddr)
	if err != nil {
		return "", "", "", "", err
	}
	FQDN := ldap.Host
	parts := strings.Split(FQDN, ".")

	//A.B.C.D A:域控制器 B:dcName B.C.D:域名
	dcHostName = parts[0]
	domainName = parts[1]
	domain = strings.Join(parts[1:], ".")
	dn = "DC=" + strings.Join(parts[1:], ",DC=")

	return
}

func Tar(src, dst, password string) (err error) {
	// 如果存在特殊字符，抛出异常，防止系统命令执行
	for _, ch := range []string{" ", "|", "&"} {
		if strings.Contains(src, ch) {
			return fmt.Errorf("illegal char in src")
		}
		if strings.Contains(dst, ch) {
			return fmt.Errorf("illegal char in dst")
		}
	}

	var tarCmd string
	if password == "" {
		tarCmd = fmt.Sprintf("tar -czvf %s %s", dst, src)
	} else {
		tarCmd = fmt.Sprintf("tar -czvf - %s | openssl des3 -salt -k %s -out %s", src, password, dst)
	}
	c := exec.Command("bash", "-c", tarCmd)
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

// PasswordEncrypt ldap password encrypt
func PasswordEncrypt(password string) (encrypt string, err error) {
	aesUtil := crypto.NewAes([]byte(common.RDX_CRYPT_SECRET))
	aesEncrypt, err := aesUtil.Encrypt(password)
	if err != nil {
		return "", err
	}
	sEnc := base64.StdEncoding.EncodeToString(aesEncrypt)
	return sEnc, nil
}

// PasswordDecode ldap password decode
func PasswordDecode(encrypt string) (password string, err error) {
	aesInstance := crypto.NewAes([]byte(common.RDX_CRYPT_SECRET))
	encByte, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	cfgStr, err := aesInstance.Decrypt(encByte)
	if err != nil {
		return "", err
	}

	return string(cfgStr), nil
}
