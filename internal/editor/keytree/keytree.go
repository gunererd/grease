package keytree

import (
	"time"

	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

// KeyAction represents what to do when a key sequence is matched
type KeyAction struct {
	Before  func(e types.Editor) types.Editor
	Execute func(e types.Editor) types.Editor
	After   func(e types.Editor) types.Editor
}

// KeyNode represents a node in tree
type KeyNode struct {
	children map[string]*KeyNode
	action   *KeyAction
}

// KeyTree manages key sequences
type KeyTree struct {
	roots   map[state.Mode]*KeyNode
	current *KeyNode
	timeout time.Duration
	lastKey time.Time
	mode    state.Mode
}

// NewKeyTree creates a new KeyTree with default timeout of 1 second
func NewKeyTree() *KeyTree {
	return &KeyTree{
		roots:   make(map[state.Mode]*KeyNode),
		timeout: time.Second * 1,
	}
}

// Add registers a new key sequence with action for a specific mode
func (kt *KeyTree) Add(mode state.Mode, sequence []string, action KeyAction) {
	if len(sequence) == 0 {
		return
	}

	// Initialize root node for mode if it doesn't exist
	if kt.roots[mode] == nil {
		kt.roots[mode] = &KeyNode{children: make(map[string]*KeyNode)}
	}

	node := kt.roots[mode]
	for _, key := range sequence {
		if node.children[key] == nil {
			node.children[key] = &KeyNode{children: make(map[string]*KeyNode)}
		}
		node = node.children[key]
	}
	node.action = &action
}

// SetMode updates the current mode and resets the current node
func (kt *KeyTree) SetMode(mode state.Mode) {
	kt.mode = mode
	kt.current = nil
}

// Handle processes a key press for the current mode
func (kt *KeyTree) Handle(key string, e types.Editor) (handled bool, model types.Editor) {
	now := time.Now()

	// Reset if timeout exceeded
	if kt.current != nil && now.Sub(kt.lastKey) > kt.timeout {
		kt.current = nil
	}

	// Get root node for current mode
	root := kt.roots[kt.mode]
	if root == nil {
		return false, e
	}

	// Start from root if no current node
	if kt.current == nil {
		kt.current = root
	}

	// Traverse a branch
	next := kt.current.children[key]
	if next == nil {
		// Key doesn't match any sequence from branch
		kt.current = root
		return false, e
	}

	// Key is part of a sequence
	kt.current = next
	kt.lastKey = now

	// Check if we've reached an actio
	if kt.current.action != nil {
		action := kt.current.action
		kt.current = root
		if action.Before != nil {
			e = action.Before(e)
		}
		if action.Execute != nil {
			e = action.Execute(e)
		}
		if action.After != nil {
			e = action.After(e)
		}
		return true, e
	}

	// Key is part of an incomplete sequence
	return true, e
}

func (kt *KeyTree) Reset() {
	kt.current = kt.roots[kt.mode]
}
