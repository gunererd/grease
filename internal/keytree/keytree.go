package keytree

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/types"
)

// KeyAction represents what to do when a key sequence is matched
type KeyAction struct {
	Execute func(e types.Editor) (tea.Model, tea.Cmd)
}

// KeyNode represents a node in tree
type KeyNode struct {
	children map[string]*KeyNode
	action   *KeyAction
}

// KeyTree manages key sequences
type KeyTree struct {
	root    *KeyNode
	current *KeyNode
	timeout time.Duration
	lastKey time.Time
}

// NewKeyTree creates a new KeyTree with default timeout of 1 second
func NewKeyTree() *KeyTree {
	return &KeyTree{
		root:    &KeyNode{children: make(map[string]*KeyNode)},
		timeout: time.Second * 1,
	}
}

// Add registers a new key sequence with action
func (kt *KeyTree) Add(sequence []string, action KeyAction) {
	if len(sequence) == 0 {
		return
	}

	node := kt.root
	for _, key := range sequence {
		if node.children[key] == nil {
			node.children[key] = &KeyNode{children: make(map[string]*KeyNode)}
		}
		node = node.children[key]
	}
	node.action = &action
}

// handled: true if the key was consumed as part of a sequence (whether complete or partial)
// handled: false if the key doesn't match any sequence
func (kt *KeyTree) Handle(key string, e types.Editor) (handled bool, model tea.Model, cmd tea.Cmd) {
	now := time.Now()

	// Reset if timeout exceeded
	if kt.current != nil && now.Sub(kt.lastKey) > kt.timeout {
		kt.current = kt.root
	}

	// Start from root if no current node
	if kt.current == nil {
		kt.current = kt.root
	}

	// Traverse a branch
	next := kt.current.children[key]
	if next == nil {
		// Key doesn't match any sequence from branch
		kt.current = kt.root
		return false, e, nil
	}

	// Key is part of a sequence
	kt.current = next
	kt.lastKey = now

	// Check if we've reached an action
	if kt.current.action != nil {
		action := kt.current.action
		kt.current = kt.root
		model, cmd = action.Execute(e)
		return true, model, cmd
	}

	// Key is part of an incomplete sequence
	return true, e, nil
}

func (kt *KeyTree) Reset() {
	kt.current = kt.root
}
