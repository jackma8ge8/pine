package channelservice

import (
	"sync"
)

var mutex sync.Mutex
var channelStore = make(map[string]*Channel)

// CreateChannel 创建一个channel
func CreateChannel(channelID string) *Channel {

	mutex.Lock()
	defer mutex.Unlock()

	channelInstance, ok := channelStore[channelID]
	if ok {
		return channelInstance
	}
	channelIns := &Channel{}
	channelStore[channelID] = channelIns

	return channelIns
}
