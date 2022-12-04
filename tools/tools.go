package tools

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/switchupcb/disgo"
	"golang.org/x/sync/errgroup"
)

var (
	// Signals represents common termination signals used to terminate the program.
	//
	// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
	Signals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	}
)

// InterceptSignal blocks until the provided signals are intercepted.
//
// Upon receiving a signal, the given sessions are gracefully disconnected.
func InterceptSignal(signals []os.Signal, sessions ...*disgo.Session) error {
	signalChannel := make(chan os.Signal, 1)

	// set the syscalls that signalChannel is sent.
	signal.Notify(signalChannel, signals...)

	// block the calling goroutine until a signal is received.
	<-signalChannel

	disgo.Logger.Info().Msg("Closing sessions due to signal...")

	eg := errgroup.Group{}
	for _, session := range sessions {
		s := session

		eg.Go(func() error {
			if err := s.Disconnect(); err != nil {
				err = fmt.Errorf("error closing connection to Discord Gateway: %w", err)
				disgo.LogSession(disgo.Logger.Error(), s.ID).Err(err).Msg("")

				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		disgo.Logger.Warn().Msg("Not all sessions were closed successfully.")

		return fmt.Errorf("error during signal intercept for termination: %w", err)
	}

	disgo.Logger.Info().Msg("Closed sessions successfully.")

	return nil
}
