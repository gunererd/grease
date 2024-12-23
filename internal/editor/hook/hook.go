package hook

import (
	"github.com/gunererd/grease/internal/editor/types"
)

// Manager manages command hooks
type Manager struct {
	hooks []types.Hook
}

func NewManager() types.HookManager {
	return &Manager{}
}

func (m *Manager) AddHook(h types.Hook) {
	m.hooks = append(m.hooks, h)
}

func (m *Manager) RemoveHook(h types.Hook) {
	for i, hook := range m.hooks {
		if hook == h {
			m.hooks = append(m.hooks[:i], m.hooks[i+1:]...)
			break
		}
	}
}

func (m *Manager) GetHooks() []types.Hook {
	return m.hooks
}

func (m *Manager) ExecuteBeforeHooks(cmd types.Command, e types.Editor) {
	for _, h := range m.hooks {
		h.OnBeforeCommand(cmd, e)
	}
}

func (m *Manager) ExecuteAfterHooks(cmd types.Command, e types.Editor) {
	for _, h := range m.hooks {
		h.OnAfterCommand(cmd, e)
	}
}
