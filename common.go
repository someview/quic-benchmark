package quicbenchmark

import (
	"time"

	"github.com/quic-go/quic-go"
)

const Addr = "localhost:4242"

const Message = "foobar"

var QuicConf = &quic.Config{
	Allow0RTT:            true,
	HandshakeIdleTimeout: time.Second * 120,
	MaxIncomingStreams:   1 << 20, // 40万个incoming
	MaxIdleTimeout:       time.Minute * 10,
	// Tracer: func(ctx context.Context, p logging.Perspective,
	// 	ci quic.ConnectionID) *logging.ConnectionTracer {
	// 	return qlog.NewConnectionTracer(NewBufferedWriteCloser(wc, f), p, ci)
	// },
}
