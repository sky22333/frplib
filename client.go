package frplib

import (
	"context"
	"time"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/config/source"
)

func StartClientWithID(id, configToml string) string {
	return errorString(startClientWithID(id, configToml))
}

func StopClientWithID(id string) string {
	return errorString(stopClientWithID(id))
}

func ReloadClientWithID(id, configToml string) string {
	return errorString(reloadClientWithID(id, configToml))
}

func IsClientRunningWithID(id string) bool {
	return globalManager.isRunning(instanceTypeClient, id)
}

func startClientWithID(id, configToml string) error {
	return globalManager.start(instanceTypeClient, id, configToml, newClientService)
}

func stopClientWithID(id string) error {
	return globalManager.stop(instanceTypeClient, id)
}

func reloadClientWithID(id, configToml string) error {
	return globalManager.reload(instanceTypeClient, id, configToml, validateClientConfig, newClientService)
}

type clientService struct {
	svc *client.Service
}

func (s *clientService) Close() error {
	s.svc.Close()
	return nil
}

func (s *clientService) GracefulClose(d time.Duration) error {
	s.svc.GracefulClose(d)
	return nil
}

func validateClientConfig(configPath string) error {
	common, proxies, visitors, isLegacyFormat, err := config.LoadClientConfig(configPath, true)
	if err != nil {
		return newError(ErrInvalidToml, "parse frpc TOML failed: %v", err)
	}
	if isLegacyFormat {
		return newError(ErrInvalidToml, "legacy frpc config format is not supported")
	}
	if common == nil || proxies == nil || visitors == nil {
		return nil
	}
	return nil
}

func newClientService(ctx context.Context, configPath string) (closeableService, error) {
	common, proxies, visitors, isLegacyFormat, err := config.LoadClientConfig(configPath, true)
	if err != nil {
		return nil, newError(ErrInvalidToml, "parse frpc TOML failed: %v", err)
	}
	if isLegacyFormat {
		return nil, newError(ErrInvalidToml, "legacy frpc config format is not supported")
	}

	configSource := source.NewConfigSource()
	if err := configSource.ReplaceAll(proxies, visitors); err != nil {
		return nil, newError(ErrStartFailed, "load frpc config source failed: %v", err)
	}
	aggregator := source.NewAggregator(configSource)

	svc, err := client.NewService(client.ServiceOptions{
		Common:                 common,
		ConfigFilePath:         configPath,
		ConfigSourceAggregator: aggregator,
	})
	if err != nil {
		return nil, newError(ErrStartFailed, "create frpc service failed: %v", err)
	}

	go func() {
		if err := svc.Run(ctx); err != nil {
			emitLog("", instanceTypeClient, "error", err.Error())
		}
	}()

	return &clientService{svc: svc}, nil
}
