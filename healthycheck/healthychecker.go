package healthycheck

import (
	"fmt"
	"sync"
	"sync/atomic"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"          //httprouter
)

//CommonHealthyChecker是HealthyHandler的实例化
type CommonHealthyChecker struct {
	Lock  sync.RWMutex
	Count int64

	//SAVE ALL LIVENESS
	LivenessProbe map[string]HealthyCheckFunc

	//SAVE ALL READNESS
	ReadnessProbe map[string]HealthyCheckFunc
}

func NewCommonHealthyChecker() *CommonHealthyChecker {
	hl := &CommonHealthyChecker{
		Lock:          sync.RWMutex{},
		Count:         0,
		LivenessProbe: make(map[string]HealthyCheckFunc),
		ReadnessProbe: make(map[string]HealthyCheckFunc),
	}

	return hl
}

//添加readless探针
func (h *CommonHealthyChecker) SetReadnessProbe(name string, check HealthyCheckFunc) error {
	h.Lock.Lock()
	defer h.Lock.Unlock()

	if _, ok := h.ReadnessProbe[name]; ok {
		return fmt.Errorf("%s has exist in readness", name)
	}

	h.ReadnessProbe[name] = check

	return nil
}

func (h *CommonHealthyChecker) GetTotalCount() int64 {
	return atomic.LoadInt64(&h.Count)
}

//添加livenessless探针
func (h *CommonHealthyChecker) HandleSetLivenessProbe(name string, check HealthyCheckFunc) error {
	h.Lock.Lock()
	defer h.Lock.Unlock()

	if _, ok := h.LivenessProbe[name]; ok {
		return fmt.Errorf("%s has exist in liveness", name)
	}

	h.LivenessProbe[name] = check

	return nil
}

func (h *CommonHealthyChecker) HandleDoSingleCheck(name string) error {
	//TODO:pandaychen
	return nil
}

func (h *CommonHealthyChecker) HandleDoEndpointReadnessCheck() *HealthyCheckResult {
	//TODO: 接口频率压制
	rlist := NewHealthyCheckResult()
	h.Lock.RLock()
	defer h.Lock.RUnlock()

	for name, checkfunc := range h.ReadnessProbe {
		if err := checkfunc(); err != nil {
			rlist.PushErr(name, err)
			return rlist
		}

		rlist.PushSucc(name,"HTTP")
	}

	return rlist
}

func (h *CommonHealthyChecker) HandleDoEndpointLivenessCheck() *HealthyCheckResult {
	//TODO: 接口频率压制
	rlist := NewHealthyCheckResult()
	h.Lock.RLock()
	defer h.Lock.RUnlock()

	for name, checkfunc := range h.LivenessProbe {
		if err := checkfunc(); err != nil {
			rlist.PushErr(name, err)
			return rlist
		}

		rlist.PushSucc(name,"HTTP")
	}

	return rlist
}

//Send HTTP RESPONSE
func (h *CommonHealthyChecker) HandleGetEndpointReadnessCheckResult(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	atomic.AddInt64(&h.Count, 1)

	result := h.HandleDoEndpointReadnessCheck()
	if result.Code != 0 {
		ctx.Response.SetStatusCode(fasthttp.StatusServiceUnavailable)
	}

	ctx.Write(result.ToJson())

	return nil
}


//Send HTTP RESPONSE
func (h *CommonHealthyChecker) HandlePingCheckResult(ctx *routing.Context) error {
	atomic.AddInt64(&h.Count, 1)

	ctx.Write([]byte("200 OK"))

	return nil
}
