package models

type Response struct {
	Success bool `json:"success"`
	Code int `json:"code"`
	Reason string `json:"reason,omitempty"`
	Result interface{} `json:"result,omitempty"`
}
