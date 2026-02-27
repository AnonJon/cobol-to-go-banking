package model

type Control struct {
	Name     string `json:"name" db:"name"`
	ValueNum int    `json:"valueNum" db:"value_num"`
	ValueStr string `json:"valueStr" db:"value_str"`
}
