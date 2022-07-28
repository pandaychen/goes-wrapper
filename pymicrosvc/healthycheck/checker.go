package healthycheck

//use fasthttp to build healthy-check cgi

import (
	"fmt"
	"runtime"
	"github.com/qiangxue/fasthttp-routing"
)

type HealthyCheckFunc = func() error

type HealthyHandler interface {
	HandleAddLivenessProbe(name string, check HealthyCheckFunc) error
	HandleAddReadnessProbe(name string, check HealthyCheckFunc) error
	HandleDoEndpointCheck(ctx *routing.Context) error
	HandleDoSingleCheck(ctx *routing.Context,name string )error

	//获取所有的检查结果
	HandleGetEndpointReadnessCheckResult(ctx *routing.Context) (error)
}

////HealthyCheckFunc实例化

//构造Server检查
func CheckGoroutineCount(max int) HealthyCheckFunc {
	return func() error {
		cnt := runtime.NumGoroutine()
		if cnt > max {
			return fmt.Errorf("goroutines overflow, count: %d", cnt)
		}
		return nil
	}
}

