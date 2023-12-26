package internal

import "errors"

var Users = map[string]int{
	"hash": 10000,
}

// 获取用户uid
func GetUid(token string) (int, error) {
	uid, ok := Users[token]
	if ok {
		return uid, nil
	} else {
		return 0, errors.New("error token")
	}
}
