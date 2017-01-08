package btmetaspider

import (
    "net"
)

//type UDPPackage struct {
//    addr *net.UDPAddr
//    data []byte
//}

type Config struct {
    Addr string
    NodeBufferSize int
    PeerBufferSize int
}

type NodeInfo struct {
    id string
    ip net.IP
    port int
}

type PeerInfo struct {
    infohash string
    ip net.IP
    port int
}
