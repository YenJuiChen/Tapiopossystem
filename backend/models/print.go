package models

type PrintRequest struct {
	SerialNo    string `json:"serial_no"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Amount      string `json:"amount"`
	Phoen       string `json:"phoen"`
	Date        string `json:"date"`
	Address     string `json:"address"`
	ServiceName string `json:"service_name"` // ✅ 新增這行對應 template
}
