package quicbenchmark

import (
	"time"

	"github.com/quic-go/quic-go"
)

const Addr = "localhost:4242"

const Message = `明月当空照古今，你好啊,
                 明月当空照古今，你好啊,
				 明月当空照古今，你好啊,
				 明月当空照古今，你好啊,
				 明月当空照古今，你好啊,
				 明月当空照古今，你好啊
				 `

// 1.调大缓冲区

var QuicConf = &quic.Config{
	Allow0RTT:                  true,
	HandshakeIdleTimeout:       time.Second * 120,
	MaxIncomingStreams:         1 << 20, // 40万个incoming
	MaxIdleTimeout:             time.Minute * 10,
	MaxStreamReceiveWindow:     1 << 15, // 32M
	MaxConnectionReceiveWindow: 2 << 20, // 2G
	// Tracer: func(ctx context.Context, p logging.Perspective,
	// 	ci quic.ConnectionID) *logging.ConnectionTracer {
	// 	return qlog.NewConnectionTracer(NewBufferedWriteCloser(wc, f), p, ci)
	// },
}
