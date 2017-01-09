package smtp

import (
	"net"

	"github.com/ian-kent/linkio"
)

// ChaosMonkey should be implemented by chaos monkeys!
type ChaosMonkey interface {
	// Accept is called for each incoming connection. Returning false closes the connection.
	Accept(ctx string, conn net.Conn) bool
	// LinkSpeed sets the maximum connection throughput (in one direction)
	LinkSpeed(ctx string) *linkio.Throughput

	// ValidRCPT is called for the RCPT command. Returning false signals an invalid recipient.
	ValidRCPT(ctx, rcpt string) bool
	// ValidMAIL is called for the MAIL command. Returning false signals an invalid sender.
	ValidMAIL(ctx, mail string) bool
	// ValidAUTH is called after authentication. Returning false signals invalid authentication.
	ValidAUTH(ctx, mechanism string, args ...string) bool

	// Disconnect is called after every read. Returning true will close the connection.
	Disconnect(ctx string) bool
}
