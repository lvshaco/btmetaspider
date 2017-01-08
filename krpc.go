package btmetaspider

import (
    "log"
    "net"
    "encoding/json"
)

type Krpc struct {
    conn *net.UDPConn
    nid string
    cbGetNode func(ni *NodeInfo)
    cbGetPeer func(pi *PeerInfo)
}

func NewKrpc(conn *net.UDPConn,
        cbGetNode func(ni *NodeInfo),
        cbGetPeer func(pi *PeerInfo)) *Krpc {
    return &Krpc{
        conn: conn,
        nid: randomNodeId(),
        cbGetNode: cbGetNode,
        cbGetPeer: cbGetPeer}
}

func (p *Krpc) Handle(addr *net.UDPAddr, pkg []byte) {
    vs := assertbdecode(pkg)
    o, err := json.MarshalIndent(vs, "", " ")
    log.Println("pkg: ", len(pkg))
    if err == nil {
        log.Println("json v: ", o)
    } else {
        log.Println("json error: ", err)
    }
    v := vs.(map[string]interface{})
    t := v["t"].(string)
    y := v["y"].(string)
    switch y {
    case "r":
        r := v["r"].(map[string]interface{})
        id := r["id"].(string)
        nodes := r["nodes"].(string)
        p.HFindNode(addr, id, nodes)
    case "q":
        q := v["q"].(string)
        a := v["a"].(map[string]interface{})
        switch q {
        case "get_peers":
            p.HGetPeers(addr, t, a)
        case "announce_peer":
            p.HAnnouncePeer(addr, t, a)
        }
    }
}

func (p *Krpc) HFindNode(addr *net.UDPAddr, id string, nodes string) {
    if len(nodes) % 26 != 0 {
        panic("Invalid compact node info")
    }
    for i:=0; i<len(nodes); i=i+26 {
        node := nodes[i:i+26]
        id := nodes[:20]
        ip := net.IPv4(node[20], node[21], node[22], node[23])
        port := int((uint16(node[24])<<8) + uint16(node[25]))
        p.cbGetNode(&NodeInfo{
            id: id,
            ip: ip,
            port: port,
        })
    }
}

func (p *Krpc) HGetPeers(addr *net.UDPAddr, t string, a map[string]interface{}) {
    id := a["id"].(string)
    infohash := a["info_hash"].(string)
    r := map[string]interface{} {
        "id": getNeighbor(id, p.nid),
        "token": infohash[:2],
        "nodes": "",
    }
    p.response(addr, t, r)
}

func (p *Krpc) HAnnouncePeer(addr *net.UDPAddr, t string, a map[string]interface{}) {
    id := a["id"].(string)
    infohash := a["info_hash"].(string)
    token := a["token"].(string)
    if (len(token) != 2) || (token != infohash[:2]) {
        panic("Invalid token")
    }
    implied_port, ok := a["implied_port"]
    var port int
    if ok && (implied_port != 0) {
        port = addr.Port
    } else {
        port = a["port"].(int)
    }
    r := map[string]interface{} {
        "id": getNeighbor(id, p.nid),
    }
    p.response(addr, t, r)
    p.cbGetPeer(&PeerInfo{
        ip: addr.IP,
        port: port,
        infohash: infohash,
    })
}

func (p *Krpc) QFindNode(addr *net.UDPAddr, nid string) {
    if nid == "" {
        nid = p.nid
    } else {
        nid = getNeighbor(nid, p.nid)
    }
    t := transactionId()
    a := map[string]interface{} {
        "id": nid,
        "target": randomNodeId(),
    }
    p.query(addr, t, "find_node", a)
}

func (p *Krpc) query(addr *net.UDPAddr, t string, qt string, a map[string]interface{}) {
    v := map[string]interface{} {
        "t": t,
        "y": "q",
        "q": qt,
        "a": a,
    }
    p.send(addr, v)
}

func (p *Krpc) response(addr *net.UDPAddr, t string, r map[string]interface{}) {
    v := map[string]interface{} {
        "t": t,
        "y": "r",
        "r": r,
    }
    p.send(addr, v)
}

func (p *Krpc) send(addr *net.UDPAddr, v map[string]interface{}) {
    msg := assertbencode(v)
    p.conn.WriteToUDP(msg, addr)
}
