package btmetaspider

import (
    "net"
    "log"
)

type Pspider struct {
    listenConn *net.UDPConn
    krpc *Krpc
    nodeChan chan *NodeInfo
    peerChan chan *PeerInfo
}

func NewPspider(cfg *Config) *Pspider {
    log.Println("Listen on udp: "+cfg.Addr)
    addr, err := net.ResolveUDPAddr("udp", cfg.Addr)
    if err != nil {
        log.Println(err)
        return nil
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Println(err)
        return nil
    }
    nc := make(chan *NodeInfo, cfg.NodeBufferSize)
    pc := make(chan *PeerInfo, cfg.PeerBufferSize)
    krpc := NewKrpc(conn,
        func(ni *NodeInfo) {
            nc <- ni
        },
        func(pi *PeerInfo) {
            pc <- pi
        },
    )
    return &Pspider{
        listenConn: conn,
        krpc: krpc,
        nodeChan: nc,
        peerChan: pc,
    }
}

func (p *Pspider) Run() {
    listenConn := p.listenConn
    for {
        data := make([]byte, 1024)
        nread, addr, err := listenConn.ReadFromUDP(data)
        if err != nil {
            log.Fatal(err)
        }
        go p.handle(addr, data[:nread])
    }
}

func (p *Pspider) handle(addr *net.UDPAddr, pkg []byte) {
    defer func() {
        if e := recover(); e != nil {
            log.Println("Error Pspider handle:", e.(error))
        }
    }()
    p.krpc.Handle(addr, pkg)
}

func (p *Pspider) NodeChan() chan *NodeInfo {
    return p.nodeChan
}

func (p *Pspider) PeerChan() chan *PeerInfo {
    return p.peerChan
}

func (p *Pspider) Krpc() *Krpc {
    return p.krpc
}
