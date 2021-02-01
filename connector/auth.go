package connector

import "errors"

var authFunc = func(uid string, token string, sessionData map[string]string) error {
	if uid == "" || token == "" {
		return errors.New("认证失败")
	}
	return nil
}

// RegisteAuth Registe auth func
func RegisteAuth(auth func(uid string, token string, sessionData map[string]string) error) {
	authFunc = auth
}
