package frplib

import "sync"

type FrpLogCallback interface {
	OnLog(instanceID string, typ string, level string, message string)
}

var logBridge = struct {
	sync.RWMutex
	callback FrpLogCallback
}{}

func SetLogCallback(callback FrpLogCallback) {
	logBridge.Lock()
	logBridge.callback = callback
	logBridge.Unlock()
}

func emitLog(instanceID, typ, level, message string) {
	logBridge.RLock()
	callback := logBridge.callback
	logBridge.RUnlock()

	if callback == nil {
		return
	}

	go callback.OnLog(instanceID, typ, level, sanitizeLog(message))
}

func sanitizeLog(message string) string {
	return message
}
