package wrapper

// coroner investigates when a Session's goroutines are shutdown.
func (s *Session) coroner() {
	// wait until all the manager goroutines is closed.
	err := s.manager.coroner.Wait()

	s.Lock()

	// report the disconnection error
	s.manager.actionError <- err
	close(s.manager.actionError)

	// remove the session from the client.
	s.client_manager.RemoveGatewaySession(s.ID)

	s.logClose("coroner")
	s.Unlock()
}

// Wait blocks until the calling Session is inactive (due to a final disconnect),
// then returns the Session's state and the disconnection error (if it exists).
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
