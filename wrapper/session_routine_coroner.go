package wrapper

// coroner investigates when a Session's goroutines are shutdown.
func (s *Session) coroner() {
	// wait until the manager goroutine is closed.
	if err := s.manager.coroner.Wait(); err != nil {
		LogSession(Logger.Error(), s.ID).Err(err).Msg("coroner manager routine error")
	}

	// Reset the session.
	putSession(s)

	Logger.Info().Msg("closed coroner routine")
}

// Wait blocks until the calling Session is inactive (due to a final disconnect),
// then returns the Session's state and the disconnection error (when it exists).
//
// If Wait() is called on a Session that isn't connected, it will return immediately
// with code SessionStateNew.
//
// A disconnected session is reset and placed into a memory pool,
// so do NOT modify a Session after it disconnects.
func (s *Session) Wait() (string, error) {
	if s.State() == SessionStateNew {
		return SessionStateNew, nil
	}

	err := s.manager.coroner.Wait()

	return s.State(), err
}
