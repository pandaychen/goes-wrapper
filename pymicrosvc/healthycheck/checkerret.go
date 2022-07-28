package healthycheck

//use fasthttp to build healthy-check cgi

import (
	"fmt"
	"encoding/json"
)

//single unit
type HealthyCheckRet struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	hctype	string `json:"type"`
}

type HealthyCheckResult struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data []HealthyCheckRet `json:"data"`
}

//new obj
func NewHealthyCheckResult() *HealthyCheckResult {
	return &HealthyCheckResult{
		Data: make([]HealthyCheckRet, 0),
	}
}

func (r *HealthyCheckResult) PushSucc(name,hctype string) {
	if hctype=="HTTP"{
		r.Data = append(r.Data, HealthyCheckRet{name, "200 OK",hctype})
	}else{
		r.Data = append(r.Data, HealthyCheckRet{name, "succ",hctype})
	}
}

func (r *HealthyCheckResult) PushErr(name string, e error) {
	r.Code += 1
	r.Msg += e.Error() + "\n"
}

func (r *HealthyCheckResult) ToJson() []byte {
	b, err := json.Marshal(r)
	if err!=nil{
		fmt.Println(err)
		return nil
	}
	return b
}
