package b5

import "errors"

type memoryManager struct {
	currentStack     int
	valueDefinitions []map[string]interface{}
}

func (m memoryManager) newValue(key string, value interface{}) {
	m.valueDefinitions[m.currentStack][key] = value
}

func (m memoryManager) getValue(key string) (interface{}, error) {
	var v interface{}

	s := m.currentStack
	for {
		if s >= len(m.valueDefinitions) || s < 0 {
			return nil, errors.New(`could not find identifier "` + key + `"`)
		}

		if value, ok := m.valueDefinitions[s][key]; ok {
			v = value
			break
		}

		s--
	}

	return v, nil
}

func (m memoryManager) pushStack() {
	m.currentStack++
	m.valueDefinitions = append(m.valueDefinitions, make(map[string]interface{}))
}

func (m memoryManager) popStack() {
	m.valueDefinitions = m.valueDefinitions[0 : len(m.valueDefinitions)-1]
	m.currentStack--
}

