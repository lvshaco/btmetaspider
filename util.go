package btmetaspider

import (
//    "math/rand"
    "crypto/rand"
    "github.com/lvshaco/bencode"
)

func assertbdecode(data []byte) interface{} {
    v, err := bencode.Decode(data)
    if err != nil {
        panic(err)
    }
    return v
}

func assertbencode(it interface{}) []byte {
    v, err := bencode.Encode(it)
    if err != nil {
        panic(err)
    }
    return v
}

func randomNodeId() string {
    b := make([]byte, 20)
    rand.Read(b)
    return string(b)
}

func getNeighbor(targetid string, myid string) string {
    return targetid[:10] + myid[10:]
}

func transactionId() string {
    b := make([]byte, 2)
    rand.Read(b)
    return string(b)
}
