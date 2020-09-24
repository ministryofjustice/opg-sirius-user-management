package server

import (
	"io"
)

type mockTemplate struct {
	count    int
	lastName string
	lastVars interface{}
}

func (m *mockTemplate) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars
	return nil
}
