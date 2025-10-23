// Package reservoir - Membrane computing components for P-system reservoir evolution
package reservoir

import (
	"fmt"
	"math"
)

// Membrane represents a P-system membrane in the hierarchical computing structure.
// Membranes contain objects and evolution rules, and interact through permeability.
type Membrane struct {
	ID           string         // Unique identifier
	Level        int            // Depth level in hierarchy (0 = outermost)
	Permeability float64        // How easily objects pass through (0-1)
	Objects      []Object       // Objects contained in this membrane
	Rules        []EvolutionRule // Evolution rules for this membrane
	Parent       *Membrane      // Parent membrane (nil for root)
	Children     []*Membrane    // Child membranes
}

// Object represents a computational object in a membrane.
type Object struct {
	Type     string  // Object type identifier
	Value    float64 // Object value/charge
	Charge   int     // Positive/negative/neutral charge
	Mobility float64 // How easily the object moves (0-1)
}

// EvolutionRule represents a P-system evolution rule.
// Rules transform objects within or across membranes.
type EvolutionRule struct {
	Name        string   // Rule identifier
	Priority    int      // Execution priority (higher = first)
	InputTypes  []string // Required input object types
	OutputTypes []string // Produced output object types
	Conditions  []string // Conditions for rule activation
	Transform   func([]Object) []Object // Transformation function
}

// NewMembrane creates a new membrane with the specified properties.
func NewMembrane(id string, level int, permeability float64) *Membrane {
	return &Membrane{
		ID:           id,
		Level:        level,
		Permeability: permeability,
		Objects:      make([]Object, 0),
		Rules:        make([]EvolutionRule, 0),
		Children:     make([]*Membrane, 0),
	}
}

// AddChild adds a child membrane to this membrane.
func (m *Membrane) AddChild(child *Membrane) {
	child.Parent = m
	m.Children = append(m.Children, child)
}

// AddObject adds an object to the membrane.
func (m *Membrane) AddObject(obj Object) {
	m.Objects = append(m.Objects, obj)
}

// AddRule adds an evolution rule to the membrane.
func (m *Membrane) AddRule(rule EvolutionRule) {
	m.Rules = append(m.Rules, rule)
}

// Evolve applies evolution rules to the objects in the membrane.
// This implements the P-system computation step.
func (m *Membrane) Evolve() error {
	// Sort rules by priority
	sortedRules := make([]EvolutionRule, len(m.Rules))
	copy(sortedRules, m.Rules)
	
	// Simple bubble sort by priority
	for i := 0; i < len(sortedRules); i++ {
		for j := i + 1; j < len(sortedRules); j++ {
			if sortedRules[j].Priority > sortedRules[i].Priority {
				sortedRules[i], sortedRules[j] = sortedRules[j], sortedRules[i]
			}
		}
	}
	
	// Apply rules
	newObjects := make([]Object, 0)
	usedObjects := make(map[int]bool)
	
	for _, rule := range sortedRules {
		// Find matching objects for this rule
		matches := m.findMatches(rule, usedObjects)
		for _, match := range matches {
			// Mark objects as used
			for _, idx := range match {
				usedObjects[idx] = true
			}
			
			// Get input objects
			inputObjs := make([]Object, len(match))
			for i, idx := range match {
				inputObjs[i] = m.Objects[idx]
			}
			
			// Apply transformation
			if rule.Transform != nil {
				outputObjs := rule.Transform(inputObjs)
				newObjects = append(newObjects, outputObjs...)
			}
		}
	}
	
	// Keep unused objects and add new ones
	remainingObjects := make([]Object, 0)
	for i, obj := range m.Objects {
		if !usedObjects[i] {
			remainingObjects = append(remainingObjects, obj)
		}
	}
	m.Objects = append(remainingObjects, newObjects...)
	
	return nil
}

// findMatches finds sets of objects that match the rule's input requirements.
func (m *Membrane) findMatches(rule EvolutionRule, usedObjects map[int]bool) [][]int {
	matches := make([][]int, 0)
	
	if len(rule.InputTypes) == 0 {
		return matches
	}
	
	// Simple implementation: find first complete match
	match := make([]int, 0)
	typeNeeded := make(map[string]bool)
	for _, t := range rule.InputTypes {
		typeNeeded[t] = true
	}
	
	for i, obj := range m.Objects {
		if usedObjects[i] {
			continue
		}
		if typeNeeded[obj.Type] {
			match = append(match, i)
			delete(typeNeeded, obj.Type)
			if len(typeNeeded) == 0 {
				matches = append(matches, match)
				break
			}
		}
	}
	
	return matches
}

// PassObjects moves objects between membranes based on permeability and mobility.
func (m *Membrane) PassObjects(target *Membrane) error {
	if target == nil {
		return fmt.Errorf("target membrane is nil")
	}
	
	remaining := make([]Object, 0)
	for _, obj := range m.Objects {
		// Probability of passing through membrane
		passProb := m.Permeability * obj.Mobility
		if math.Abs(float64(obj.Charge)) > 0 {
			// Charged objects are more likely to move
			passProb *= 1.2
		}
		
		if passProb > 0.5 { // Simplified threshold
			target.AddObject(obj)
		} else {
			remaining = append(remaining, obj)
		}
	}
	m.Objects = remaining
	
	return nil
}

// ComputeDissolution computes membrane dissolution based on object concentration.
// Returns true if membrane should dissolve.
func (m *Membrane) ComputeDissolution() bool {
	// Membrane dissolves if it contains too many charged objects
	chargedCount := 0
	for _, obj := range m.Objects {
		if obj.Charge != 0 {
			chargedCount++
		}
	}
	
	threshold := 10 // Arbitrary threshold
	return chargedCount > threshold
}

// CreateDefaultRules creates default evolution rules for email processing.
func CreateDefaultRules() []EvolutionRule {
	rules := make([]EvolutionRule, 0)
	
	// Rule 1: Spam detection - transform high-value negative objects
	spamRule := EvolutionRule{
		Name:        "spam_detection",
		Priority:    100,
		InputTypes:  []string{"token", "negative_signal"},
		OutputTypes: []string{"spam_score"},
		Transform: func(objs []Object) []Object {
			// Combine negative signals
			totalValue := 0.0
			for _, obj := range objs {
				if obj.Charge < 0 {
					totalValue += obj.Value
				}
			}
			return []Object{{
				Type:     "spam_score",
				Value:    totalValue,
				Charge:   -1,
				Mobility: 0.8,
			}}
		},
	}
	rules = append(rules, spamRule)
	
	// Rule 2: Ham detection - transform positive signals
	hamRule := EvolutionRule{
		Name:        "ham_detection",
		Priority:    100,
		InputTypes:  []string{"token", "positive_signal"},
		OutputTypes: []string{"ham_score"},
		Transform: func(objs []Object) []Object {
			totalValue := 0.0
			for _, obj := range objs {
				if obj.Charge > 0 {
					totalValue += obj.Value
				}
			}
			return []Object{{
				Type:     "ham_score",
				Value:    totalValue,
				Charge:   1,
				Mobility: 0.8,
			}}
		},
	}
	rules = append(rules, hamRule)
	
	// Rule 3: Affective modulation - adjust scores based on emotional context
	affectiveRule := EvolutionRule{
		Name:        "affective_modulation",
		Priority:    50,
		InputTypes:  []string{"emotion_signal"},
		OutputTypes: []string{"modulated_signal"},
		Transform: func(objs []Object) []Object {
			result := make([]Object, 0)
			for _, obj := range objs {
				// Apply emotional modulation
				modulated := obj
				modulated.Value *= 1.1 // Boost emotional signals
				modulated.Type = "modulated_signal"
				result = append(result, modulated)
			}
			return result
		},
	}
	rules = append(rules, affectiveRule)
	
	return rules
}

// MembraneSystem represents a complete P-system with hierarchical membranes.
type MembraneSystem struct {
	Root     *Membrane   // Root membrane
	All      []*Membrane // All membranes in system
	StepCount int        // Number of evolution steps performed
}

// NewMembraneSystem creates a new membrane system with hierarchical structure.
func NewMembraneSystem(depth int) *MembraneSystem {
	root := NewMembrane("root", 0, 1.0)
	system := &MembraneSystem{
		Root:      root,
		All:       []*Membrane{root},
		StepCount: 0,
	}
	
	// Build hierarchical structure
	system.buildHierarchy(root, depth, 1)
	
	return system
}

// buildHierarchy recursively builds the membrane hierarchy.
func (ms *MembraneSystem) buildHierarchy(parent *Membrane, maxDepth, currentDepth int) {
	if currentDepth >= maxDepth {
		return
	}
	
	// Create 2 children at each level (binary tree)
	for i := 0; i < 2; i++ {
		id := fmt.Sprintf("%s_%d", parent.ID, i)
		permeability := 0.5 + 0.1*float64(currentDepth) // Deeper = more permeable
		child := NewMembrane(id, currentDepth, permeability)
		
		parent.AddChild(child)
		ms.All = append(ms.All, child)
		
		// Add default rules
		for _, rule := range CreateDefaultRules() {
			child.AddRule(rule)
		}
		
		// Recurse
		ms.buildHierarchy(child, maxDepth, currentDepth+1)
	}
}

// Step performs one evolution step on all membranes.
func (ms *MembraneSystem) Step() error {
	// Evolve all membranes
	for _, membrane := range ms.All {
		if err := membrane.Evolve(); err != nil {
			return fmt.Errorf("evolving membrane %s: %w", membrane.ID, err)
		}
	}
	
	// Pass objects between membranes
	for _, membrane := range ms.All {
		if membrane.Parent != nil {
			// Objects can pass to parent
			if err := membrane.PassObjects(membrane.Parent); err != nil {
				return fmt.Errorf("passing objects from %s to parent: %w", membrane.ID, err)
			}
		}
		
		// Objects can pass to children
		for _, child := range membrane.Children {
			if err := membrane.PassObjects(child); err != nil {
				return fmt.Errorf("passing objects from %s to child: %w", membrane.ID, err)
			}
		}
	}
	
	ms.StepCount++
	return nil
}

// InjectObject injects an object into a specific membrane.
func (ms *MembraneSystem) InjectObject(membraneID string, obj Object) error {
	for _, membrane := range ms.All {
		if membrane.ID == membraneID {
			membrane.AddObject(obj)
			return nil
		}
	}
	return fmt.Errorf("membrane %s not found", membraneID)
}

// CollectResults collects all objects from leaf membranes.
func (ms *MembraneSystem) CollectResults() []Object {
	results := make([]Object, 0)
	for _, membrane := range ms.All {
		if len(membrane.Children) == 0 { // Leaf membrane
			results = append(results, membrane.Objects...)
		}
	}
	return results
}
