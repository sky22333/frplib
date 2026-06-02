package frplib

import "runtime/debug"

var globalManager = newManager()

// Version returns the upstream frp version embedded by this core.
func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/fatedier/frp" {
				if dep.Replace != nil {
					return dep.Replace.Version
				}
				return dep.Version
			}
		}
	}
	return "unknown"
}

func StartClient(configToml string) string {
	return errorString(startClientWithID(defaultInstanceID, configToml))
}

func StopClient() string {
	return errorString(stopClientWithID(defaultInstanceID))
}

func ReloadClient(configToml string) string {
	return errorString(reloadClientWithID(defaultInstanceID, configToml))
}

func IsClientRunning() bool {
	return IsClientRunningWithID(defaultInstanceID)
}

func StartServer(configToml string) string {
	return errorString(startServerWithID(defaultInstanceID, configToml))
}

func StopServer() string {
	return errorString(stopServerWithID(defaultInstanceID))
}

func ReloadServer(configToml string) string {
	return errorString(reloadServerWithID(defaultInstanceID, configToml))
}

func IsServerRunning() bool {
	return IsServerRunningWithID(defaultInstanceID)
}

func StopAll() string {
	return errorString(globalManager.stopAll())
}

func ListInstances() string {
	return globalManager.listInstances()
}
