package define

import (
	"sync"

	"github.com/coderguang/GameEngine_go/sgwx/sgwxopenid"
)

type SecureWxOpenid struct {
	Data map[string]*sgwxopenid.SWxOpenid
	Lock sync.RWMutex
}


