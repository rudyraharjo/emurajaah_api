package models

type RequestLoginWithGoogle struct {
	Alias         string `json:"alias"`
	Password      string `json:"password"`
	TokenFirebase string `json:"token_firebse"`
}

type ResponseLoginFail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseLoginSuccess struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}
