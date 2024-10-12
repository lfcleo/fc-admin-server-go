package request

type CreateDictType struct {
	Name   string `json:"name"`   //字典名称
	Type   string `json:"type"`   //字典类型
	Status int    `json:"status"` //字典状态,默认=1正常,=2停用
	Notes  string `json:"notes"`  //备注
}
