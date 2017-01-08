package btmetaspider

import (
    "log"
    "fmt"
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
            log.Println(fmt.Sprintf("peer %02x on %s:%d",
                peer.infohash, peer.ip.String(), peer.port))
        }
    }
}
