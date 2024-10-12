package response

// Auth 请求/响应数据结构体
type Auth struct {
	ID     uint     `json:"id"`
	Mobile string   `json:"mobile" binding:"required"`
	Email  string   `json:"email" binding:"required"`
	Name   string   `json:"name" binding:"required"`
	Avatar string   `json:"avatar" binding:"required"`
	Sex    int      `json:"sex"`
	Status int      `json:"status,omitempty"`
	Roles  []string `json:"roles"`
}
