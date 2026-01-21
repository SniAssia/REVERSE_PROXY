package loadbalancer

import "errors"
var ErrNoAvailableBackends = errors.New("no available backends")
