package btmetaspider

import (
    "log"
    "net"
    "time"
)

type Nspider struct {
    bootstrapAddr []string
    nodeChan chan *NodeInfo
    krpc *Krpc
}

func NewNspider(bootstrapAddr []string, nc chan *NodeInfo, krpc *Krpc) *Nspider {
    return &Nspider{
        bootstrapAddr: bootstrapAddr,
        nodeChan: nc,
        krpc: krpc,
    }
}

func (p *Nspider) Run() {
    for {
        select {
        case node := <-p.nodeChan:
            log.Println("node:", node)
            addr := &net.UDPAddr{IP: node.ip, Port: node.port}
            go p.krpc.QFindNode(addr, node.id)
        case <-time.After(time.Second*3):
            p.bootstrap()
        }
    }
}

func (p *Nspider) bootstrap() {
    for _, addr := range p.bootstrapAddr {
        addr, _ := net.ResolveUDPAddr("udp", addr)
        go p.krpc.QFindNode(addr, "")
    }
}
