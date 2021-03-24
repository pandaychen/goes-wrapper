package balancer

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

type RoundRobinWrapper interface {
	AddNode(key string, node interface{}, weight int) error
	RemoveAllNodes()
	ResetAllNodes()

	GetNextNode() interface{}

	//GetNext2() interface{}
}

// Nginx的weight round robin实现
type NginxWeightRoundrobinNode struct {
	NodeKey         string
	NodeMetadata    interface{} //save all..(like grpc balancer.SubConn)
	InitWeight      int         //初始化权重
	CurrentWeight   int
	EffectiveWeight int //每次pick之后的更新的权重值
}

type NginxWeightRoundrobin struct {
	lock   sync.Mutex
	Nodes  []*NginxWeightRoundrobinNode
	Count  int
	Logger *zap.Logger
}

func NewNginxWeightRoundrobin(logger *zap.Logger) *NginxWeightRoundrobin {
	return &NginxWeightRoundrobin{Logger: logger}
}

//增加权重节点
func (r *NginxWeightRoundrobin) AddNode(key string, server interface{}, weight int) error {
	for _, tnode := range r.Nodes {
		if tnode.NodeKey == key {
			return errors.New("node exists")
		}
	}

	node := &NginxWeightRoundrobinNode{
		NodeKey:         key,
		NodeMetadata:    server,
		InitWeight:      weight,
		EffectiveWeight: weight,
		CurrentWeight:   0}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Nodes = append(r.Nodes, node)
	r.Count++
}

func (r *NginxWeightRoundrobin) ResetAllNodes() {
	for _, node := range r.Nodes {
		node.EffectiveWeight = node.InitWeight
		node.CurrentWeight = 0
	}
}

func (r *NginxWeightRoundrobin) RemoveAllNodes() {
	r.Nodes = r.Nodes[:0]
	r.Count = 0
}

func (r *NginxWeightRoundrobin) GetNextNode() *NginxWeightRoundrobinNode {
	var chosen *NginxWeightRoundrobinNode

	if r.Count == 0 {
		return nil
	} else if r.Count == 1 {
		return r.Nodes[0]
	} else {
		total := 0
		//range all nodes, choose a probably node
		for i := 0; i < len(r.Nodes); i++ {
			cnode := r.Nodes[i]
			if cnode == nil {
				continue
			} else {
				cnode.CurrentWeight += cnode.EffectiveWeight
				total += cnode.EffectiveWeight
				if cnode.EffectiveWeight < cnode.InitWeight {
					cnode.EffectiveWeight++
				}

				if chosen == nil || cnode.CurrentWeight > chosen.CurrentWeight {
					// 更换节点
					chosen = cnode
				}
			}
		}

		if chosen == nil {
			return nil
		}

		//更新选定节点的CurrentWeight值
		chosen.CurrentWeight -= total
		return chosen
	}

	return nil
}
