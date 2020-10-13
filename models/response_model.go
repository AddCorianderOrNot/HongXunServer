package models

type Response struct {
	ErrCode int         `json:"errorCode"`
	ErrMsg  string      `json:"errorMsg"`
	Data    interface{} `json:"data"`
}
