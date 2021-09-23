package context

func (m *Manager) List() (map[string]Summary, error) {
	m.readProjectSpec()
	m.readConfig()
	m.initContexts()
	m.getLocalContexts()
	return m.contexts, m.err
}

func (m *Manager) initContexts() {
	if m.err != nil {
		return
	}
	m.contexts = make(map[string]Summary)
}

func (m *Manager) getLocalContexts() {
	if m.err != nil {
		return
	}
	for contextName := range m.projectSpec.Contexts {
		m.contexts[contextName] = Summary{Name: contextName}
	}
}
