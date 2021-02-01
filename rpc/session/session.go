package session

// Get a value from session
func (session *Session) Get(key string) interface{} {
	return session.Data[key]
}

// Set a value to session
func (session *Session) Set(key string, v string) {
	if session.Data == nil {
		session.Data = make(map[string]string)
	}
	session.Data[key] = v
}
