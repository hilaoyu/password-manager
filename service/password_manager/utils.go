package password_manager

import (
	"github.com/hilaoyu/go-utils/utilEnc"
	"strings"
)

func UtilPasswordToSecret(password string) string {

	b1 := utilEnc.Md5(password)

	b2 := utilEnc.Md5(password + b1)

	return strings.ToLower(b2)
}
