// Package registry is the central hub for Drivers and other types that
// follow the init-time registration.
package registry

import (
	"strings"
	"sync"

	"github.com/akutz/goof"

	"github.com/emccode/libstorage/api/types/drivers"
	"github.com/emccode/libstorage/api/types/http"
)

var (
	storExecsCtors    = map[string]drivers.NewStorageExecutor{}
	storExecsCtorsRWL = &sync.RWMutex{}

	lclStorDriverCtors    = map[string]drivers.NewLocalStorageDriver{}
	lclStorDriverCtorsRWL = &sync.RWMutex{}

	remStorDriverCtors    = map[string]drivers.NewRemoteStorageDriver{}
	remStorDriverCtorsRWL = &sync.RWMutex{}

	osDriverCtors    = map[string]drivers.NewOSDriver{}
	osDriverCtorsRWL = &sync.RWMutex{}

	intDriverCtors    = map[string]drivers.NewIntegrationDriver{}
	intDriverCtorsRWL = &sync.RWMutex{}

	routers    = []http.Router{}
	routersRWL = &sync.RWMutex{}
)

// RegisterRouter registers a Router.
func RegisterRouter(router http.Router) {
	routersRWL.Lock()
	defer routersRWL.Unlock()
	routers = append(routers, router)
}

// RegisterStorageExecutor registers a StorageExecutor.
func RegisterStorageExecutor(name string, ctor drivers.NewStorageExecutor) {
	storExecsCtorsRWL.Lock()
	defer storExecsCtorsRWL.Unlock()
	storExecsCtors[strings.ToLower(name)] = ctor
}

// RegisterLocalStorageDriver registers a LocalStorageDriver.
func RegisterLocalStorageDriver(
	name string, ctor drivers.NewLocalStorageDriver) {
	lclStorDriverCtorsRWL.Lock()
	defer lclStorDriverCtorsRWL.Unlock()
	lclStorDriverCtors[strings.ToLower(name)] = ctor
}

// RegisterRemoteStorageDriver registers a RemoteStorageDriver.
func RegisterRemoteStorageDriver(
	name string, ctor drivers.NewRemoteStorageDriver) {
	remStorDriverCtorsRWL.Lock()
	defer remStorDriverCtorsRWL.Unlock()
	remStorDriverCtors[strings.ToLower(name)] = ctor
}

// RegisterOSDriver registers a OSDriver.
func RegisterOSDriver(name string, ctor drivers.NewOSDriver) {
	osDriverCtorsRWL.Lock()
	defer osDriverCtorsRWL.Unlock()
	osDriverCtors[strings.ToLower(name)] = ctor
}

// RegisterIntegrationDriver registers a IntegrationDriver.
func RegisterIntegrationDriver(name string, ctor drivers.NewIntegrationDriver) {
	intDriverCtorsRWL.Lock()
	defer intDriverCtorsRWL.Unlock()
	intDriverCtors[strings.ToLower(name)] = ctor
}

// NewStorageExecutor returns a new instance of the executor specified by the
// executor name.
func NewStorageExecutor(name string) (drivers.StorageExecutor, error) {

	var ok bool
	var ctor drivers.NewStorageExecutor

	func() {
		storExecsCtorsRWL.RLock()
		defer storExecsCtorsRWL.RUnlock()
		ctor, ok = storExecsCtors[strings.ToLower(name)]
	}()

	if !ok {
		return nil, goof.WithField("executor", name, "invalid executor name")
	}

	return ctor(), nil
}

// NewLocalStorageDriver returns a new instance of the driver specified by the
// driver name.
func NewLocalStorageDriver(name string) (drivers.LocalStorageDriver, error) {

	var ok bool
	var ctor drivers.NewLocalStorageDriver

	func() {
		lclStorDriverCtorsRWL.RLock()
		defer lclStorDriverCtorsRWL.RUnlock()
		ctor, ok = lclStorDriverCtors[strings.ToLower(name)]
	}()

	if !ok {
		return nil, goof.WithField("driver", name, "invalid driver name")
	}

	return ctor(), nil
}

// NewRemoteStorageDriver returns a new instance of the driver specified by the
// driver name.
func NewRemoteStorageDriver(name string) (drivers.RemoteStorageDriver, error) {

	var ok bool
	var ctor drivers.NewRemoteStorageDriver

	func() {
		remStorDriverCtorsRWL.RLock()
		defer remStorDriverCtorsRWL.RUnlock()
		ctor, ok = remStorDriverCtors[strings.ToLower(name)]
	}()

	if !ok {
		return nil, goof.WithField("driver", name, "invalid driver name")
	}

	return ctor(), nil
}

// NewOSDriver returns a new instance of the driver specified by the
// driver name.
func NewOSDriver(name string) (drivers.OSDriver, error) {

	var ok bool
	var ctor drivers.NewOSDriver

	func() {
		osDriverCtorsRWL.RLock()
		defer osDriverCtorsRWL.RUnlock()
		ctor, ok = osDriverCtors[strings.ToLower(name)]
	}()

	if !ok {
		return nil, goof.WithField("driver", name, "invalid driver name")
	}

	return ctor(), nil
}

// NewIntegrationDriver returns a new instance of the driver specified by the
// driver name.
func NewIntegrationDriver(name string) (drivers.IntegrationDriver, error) {

	var ok bool
	var ctor drivers.NewIntegrationDriver

	func() {
		intDriverCtorsRWL.RLock()
		defer intDriverCtorsRWL.RUnlock()
		ctor, ok = intDriverCtors[strings.ToLower(name)]
	}()

	if !ok {
		return nil, goof.WithField("driver", name, "invalid driver name")
	}

	return ctor(), nil
}

// StorageExecutors returns a channel on which new instances of all registered
// storage executors can be received.
func StorageExecutors() <-chan drivers.StorageExecutor {
	c := make(chan drivers.StorageExecutor)
	go func() {
		storExecsCtorsRWL.RLock()
		defer storExecsCtorsRWL.RUnlock()
		for _, ctor := range storExecsCtors {
			c <- ctor()
		}
		close(c)
	}()
	return c
}

// LocalStorageDrivers returns a channel on which new instances of all
// registered remote local drivers can be received.
func LocalStorageDrivers() <-chan drivers.LocalStorageDriver {
	c := make(chan drivers.LocalStorageDriver)
	go func() {
		lclStorDriverCtorsRWL.RLock()
		defer lclStorDriverCtorsRWL.RUnlock()
		for _, ctor := range lclStorDriverCtors {
			c <- ctor()
		}
		close(c)
	}()
	return c
}

// RemoteStorageDrivers returns a channel on which new instances of all
// registered remote storage drivers can be received.
func RemoteStorageDrivers() <-chan drivers.RemoteStorageDriver {
	c := make(chan drivers.RemoteStorageDriver)
	go func() {
		remStorDriverCtorsRWL.RLock()
		defer remStorDriverCtorsRWL.RUnlock()
		for _, ctor := range remStorDriverCtors {
			c <- ctor()
		}
		close(c)
	}()
	return c
}

// OSDrivers returns a channel on which new instances of all registered
// OS drivers can be received.
func OSDrivers() <-chan drivers.OSDriver {
	c := make(chan drivers.OSDriver)
	go func() {
		osDriverCtorsRWL.RLock()
		defer osDriverCtorsRWL.RUnlock()
		for _, ctor := range osDriverCtors {
			c <- ctor()
		}
		close(c)
	}()
	return c
}

// IntegrationDrivers returns a channel on which new instances of all registered
// integration drivers can be received.
func IntegrationDrivers() <-chan drivers.IntegrationDriver {
	c := make(chan drivers.IntegrationDriver)
	go func() {
		intDriverCtorsRWL.RLock()
		defer intDriverCtorsRWL.RUnlock()
		for _, ctor := range intDriverCtors {
			c <- ctor()
		}
		close(c)
	}()
	return c
}

// Routers returns a channel on which new instances of all registered routers
// can be received.
func Routers() <-chan http.Router {
	c := make(chan http.Router)
	go func() {
		routersRWL.RLock()
		defer routersRWL.RUnlock()
		for _, r := range routers {
			c <- r
		}
		close(c)
	}()
	return c
}
