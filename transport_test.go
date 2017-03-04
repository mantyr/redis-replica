package redis_replica_test

import (
	"testing"
	replica "github.com/mantyr/redis-replica"
)

// TestNetTransportInterface проверяет что NetTransport реализует интерфейс Transport
func TestNetTranspoerInterface(t *testing.T) {
    var _ replica.Transport = (*replica.NetTransport)(nil)
}

// TestReaderTransportInterface проверяет что ReaderTransport реализует интерфейс Transport
func TestReaderTranspoerInterface(t *testing.T) {
    var _ replica.Transport = (*replica.ReaderTransport)(nil)
}

