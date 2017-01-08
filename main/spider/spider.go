package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/lvshaco/btmetaspider"
)

func main() {
    bootstrapAddr := []string {
        "router.bittorrent.com:6881",
        "dht.transmissionbt.com:6881",
        "router.utorrent.com:6881",
    }
    cfg := &btmetaspider.Config {
        Addr: "0.0.0.0:6882",
        NodeBufferSize: 256,
        PeerBufferSize: 256,
    }
    ps := btmetaspider.NewPspider(cfg)
    go ps.Run()

    ns := btmetaspider.NewNspider(bootstrapAddr, ps.NodeChan(), ps.Krpc())
    go ns.Run()

    ms := btmetaspider.NewMspider(ps.PeerChan())
    go ms.Run()

    sigChan := make(chan os.Signal, 1)

    signal.Notify(sigChan, syscall.SIGINT)
    sig := <-sigChan
    log.Println("Receive signal to exit", sig)
}
