package btmetaspider

import (
    "log"
)

type Mspider struct {
    peerChan chan *PeerInfo
}

func NewMspider(pc chan *PeerInfo) *Mspider {
    return &Mspider{ peerChan: pc }
}

func (p *Mspider) Run() {
    for {
        for peer := range p.peerChan {
            log.Println("peer:", peer)
        }
    }
}
