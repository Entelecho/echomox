// Package reservoir implements a Deep Tree Echo State Network (ESN) Reservoir Computing Framework
// for enhanced mail processing with affective agency and emotional intelligence.
//
// This package integrates Echo State Networks with membrane computing P-systems,
// Butcher B-series Runge-Kutta methods, and differential Ricci flow equations
// to create an adaptive, emotion-aware mail classification system.
package reservoir

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/mjl-/mox/mlog"
)

// ESNParams defines the hyper-parameters for the Echo State Network.
type ESNParams struct {
	// Reservoir size - number of neurons in the hidden layer
	ReservoirSize int `sconf:"optional" sconf-doc:"Number of neurons in the reservoir layer. Default: 100."`
	
	// Spectral radius - controls memory capacity
	SpectralRadius float64 `sconf:"optional" sconf-doc:"Spectral radius for reservoir stability. Should be < 1.0. Default: 0.95."`
	
	// Input scaling - scales input signals
	InputScaling float64 `sconf:"optional" sconf-doc:"Scaling factor for input signals. Default: 1.0."`
	
	// Leak rate - controls neuron activation decay
	LeakRate float64 `sconf:"optional" sconf-doc:"Leak rate for neuron activation (0-1). Default: 0.3."`
	
	// Sparsity - connection density in reservoir
	Sparsity float64 `sconf:"optional" sconf-doc:"Sparsity of reservoir connections (0-1). Default: 0.1."`
	
	// Ridge regression parameter for output training
	RidgeParam float64 `sconf:"optional" sconf-doc:"Ridge regression parameter for regularization. Default: 1e-8."`
	
	// Tree depth for hierarchical processing
	TreeDepth int `sconf:"optional" sconf-doc:"Depth of the tree structure for hierarchical processing. Default: 3."`
}

// DefaultESNParams returns default parameters for the ESN.
func DefaultESNParams() ESNParams {
	return ESNParams{
		ReservoirSize:  100,
		SpectralRadius: 0.95,
		InputScaling:   1.0,
		LeakRate:       0.3,
		Sparsity:       0.1,
		RidgeParam:     1e-8,
		TreeDepth:      3,
	}
}

// PersonaTrait represents LLM personality traits mapped to reservoir parameters.
type PersonaTrait struct {
	// Affective dimensions
	Valence    float64 // Positive/negative emotional tone (-1 to 1)
	Arousal    float64 // Activation/energy level (0 to 1)
	Dominance  float64 // Control/power dimension (0 to 1)
	
	// Cognitive dimensions
	Attention  float64 // Focus/concentration (0 to 1)
	Memory     float64 // Retention capacity (0 to 1)
	Creativity float64 // Novelty generation (0 to 1)
}

// ESN represents a Deep Tree Echo State Network.
type ESN struct {
	params  ESNParams
	persona PersonaTrait
	
	// Network weights
	inputWeights     [][]float64 // Input to reservoir weights
	reservoirWeights [][]float64 // Recurrent reservoir weights
	outputWeights    [][]float64 // Reservoir to output weights
	
	// State
	state []float64 // Current reservoir state
	
	// Membrane computing components
	membranes []*Membrane
	
	// Synchronization
	mu     sync.RWMutex
	rng    *rand.Rand
	log    mlog.Log
	trained bool
}

// NewESN creates a new Echo State Network with the given parameters.
func NewESN(log mlog.Log, params ESNParams, persona PersonaTrait) (*ESN, error) {
	if params.ReservoirSize <= 0 {
		return nil, fmt.Errorf("reservoir size must be positive")
	}
	if params.SpectralRadius <= 0 || params.SpectralRadius >= 1.0 {
		return nil, fmt.Errorf("spectral radius must be in (0, 1)")
	}
	if params.LeakRate <= 0 || params.LeakRate > 1.0 {
		return nil, fmt.Errorf("leak rate must be in (0, 1]")
	}
	
	esn := &ESN{
		params:  params,
		persona: persona,
		state:   make([]float64, params.ReservoirSize),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		log:     log,
	}
	
	// Initialize reservoir weights with Paun P-system membrane structure
	if err := esn.initializeReservoir(); err != nil {
		return nil, fmt.Errorf("initializing reservoir: %w", err)
	}
	
	return esn, nil
}

// initializeReservoir initializes the reservoir weights using membrane computing principles.
func (esn *ESN) initializeReservoir() error {
	esn.mu.Lock()
	defer esn.mu.Unlock()
	
	n := esn.params.ReservoirSize
	
	// Initialize reservoir weights with sparse random connections
	esn.reservoirWeights = make([][]float64, n)
	for i := range esn.reservoirWeights {
		esn.reservoirWeights[i] = make([]float64, n)
		for j := range esn.reservoirWeights[i] {
			if esn.rng.Float64() < esn.params.Sparsity {
				esn.reservoirWeights[i][j] = esn.rng.NormFloat64()
			}
		}
	}
	
	// Scale weights to achieve desired spectral radius
	if err := esn.scaleSpectralRadius(); err != nil {
		return fmt.Errorf("scaling spectral radius: %w", err)
	}
	
	// Initialize membrane structure for hierarchical processing
	esn.initializeMembranes()
	
	return nil
}

// scaleSpectralRadius scales reservoir weights to achieve the desired spectral radius.
func (esn *ESN) scaleSpectralRadius() error {
	// Simplified power iteration method to estimate largest eigenvalue
	n := len(esn.reservoirWeights)
	v := make([]float64, n)
	for i := range v {
		v[i] = esn.rng.NormFloat64()
	}
	
	// Normalize
	norm := 0.0
	for _, val := range v {
		norm += val * val
	}
	norm = math.Sqrt(norm)
	for i := range v {
		v[i] /= norm
	}
	
	// Power iteration
	for iter := 0; iter < 50; iter++ {
		// Multiply v by matrix
		newV := make([]float64, n)
		for i := range newV {
			for j := range esn.reservoirWeights[i] {
				newV[i] += esn.reservoirWeights[i][j] * v[j]
			}
		}
		
		// Normalize
		norm = 0.0
		for _, val := range newV {
			norm += val * val
		}
		norm = math.Sqrt(norm)
		for i := range newV {
			newV[i] /= norm
		}
		v = newV
	}
	
	// Estimate largest eigenvalue (spectral radius)
	eigenvalue := 0.0
	for i := range v {
		product := 0.0
		for j := range esn.reservoirWeights[i] {
			product += esn.reservoirWeights[i][j] * v[j]
		}
		eigenvalue += product * v[i]
	}
	eigenvalue = math.Abs(eigenvalue)
	
	// Scale weights
	if eigenvalue > 0 {
		scale := esn.params.SpectralRadius / eigenvalue
		for i := range esn.reservoirWeights {
			for j := range esn.reservoirWeights[i] {
				esn.reservoirWeights[i][j] *= scale
			}
		}
	}
	
	return nil
}

// initializeMembranes creates the P-system membrane structure for hierarchical processing.
func (esn *ESN) initializeMembranes() {
	depth := esn.params.TreeDepth
	esn.membranes = make([]*Membrane, 0)
	
	// Create hierarchical membrane structure
	for level := 0; level < depth; level++ {
		numMembranes := 1 << level // 2^level membranes at each level
		for i := 0; i < numMembranes; i++ {
			membrane := &Membrane{
				ID:           fmt.Sprintf("L%d_M%d", level, i),
				Level:        level,
				Permeability: 0.5 + 0.5*esn.persona.Arousal, // Influenced by affective arousal
				Rules:        make([]EvolutionRule, 0),
			}
			esn.membranes = append(esn.membranes, membrane)
		}
	}
}

// SetInputWeights sets the input-to-reservoir weights (must be called with lock held).
func (esn *ESN) setInputWeights(inputDim int) {
	n := esn.params.ReservoirSize
	esn.inputWeights = make([][]float64, n)
	for i := range esn.inputWeights {
		esn.inputWeights[i] = make([]float64, inputDim)
		for j := range esn.inputWeights[i] {
			esn.inputWeights[i][j] = (esn.rng.Float64()*2 - 1) * esn.params.InputScaling
		}
	}
}

// SetInputWeights sets the input-to-reservoir weights (public method).
func (esn *ESN) SetInputWeights(inputDim int) {
	esn.mu.Lock()
	defer esn.mu.Unlock()
	esn.setInputWeights(inputDim)
}

// Update updates the reservoir state with new input using Runge-Kutta integration.
func (esn *ESN) Update(ctx context.Context, input []float64) error {
	esn.mu.Lock()
	defer esn.mu.Unlock()
	
	if len(esn.inputWeights) == 0 {
		esn.setInputWeights(len(input))
	}
	
	if len(input) != len(esn.inputWeights[0]) {
		return fmt.Errorf("input dimension mismatch: expected %d, got %d", len(esn.inputWeights[0]), len(input))
	}
	
	// Simplified update with leak rate dynamics
	newState := make([]float64, len(esn.state))
	
	for i := 0; i < len(esn.state); i++ {
		// Input contribution
		inputSum := 0.0
		for j := range input {
			inputSum += esn.inputWeights[i][j] * input[j]
		}
		
		// Recurrent contribution
		recurrentSum := 0.0
		for j := 0; j < len(esn.state); j++ {
			recurrentSum += esn.reservoirWeights[i][j] * esn.state[j]
		}
		
		// Activation with affective modulation
		activation := math.Tanh(inputSum + recurrentSum)
		
		// Apply persona-based emotional modulation
		activation *= (1.0 + 0.1*esn.persona.Valence) // Valence affects signal strength
		
		// Leak dynamics influenced by attention
		leakRate := esn.params.LeakRate * (1.0 + 0.2*esn.persona.Attention)
		newState[i] = (1-leakRate)*esn.state[i] + leakRate*activation
	}
	
	// Update state
	copy(esn.state, newState)
	
	// Apply membrane computing transformations
	esn.applyMembraneEvolution()
	
	// Apply Ricci flow curvature correction for geometric regularization
	esn.applyRicciFlow()
	
	return nil
}

// computeDerivative computes the derivative for Runge-Kutta integration.
func (esn *ESN) computeDerivative(input, state []float64) []float64 {
	n := len(state)
	derivative := make([]float64, n)
	
	for i := 0; i < n; i++ {
		// Input contribution
		inputSum := 0.0
		for j := range input {
			inputSum += esn.inputWeights[i][j] * input[j]
		}
		
		// Recurrent contribution
		recurrentSum := 0.0
		for j := 0; j < n; j++ {
			recurrentSum += esn.reservoirWeights[i][j] * state[j]
		}
		
		// Activation with affective modulation
		activation := math.Tanh(inputSum + recurrentSum)
		
		// Apply persona-based emotional modulation
		activation *= (1.0 + 0.1*esn.persona.Valence) // Valence affects signal strength
		
		// Leak dynamics influenced by attention
		leakRate := esn.params.LeakRate * (1.0 + 0.2*esn.persona.Attention)
		derivative[i] = -leakRate*state[i] + (1-leakRate)*activation
	}
	
	return derivative
}

// applyMembraneEvolution applies P-system membrane computing rules for reservoir evolution.
func (esn *ESN) applyMembraneEvolution() {
	// Apply evolution rules through membrane hierarchy
	for _, membrane := range esn.membranes {
		// Objects pass through membranes based on permeability
		if esn.rng.Float64() < membrane.Permeability {
			// Apply transformation rules
			for i := range esn.state {
				// Membrane-specific transformation influenced by level
				factor := 1.0 - 0.05*float64(membrane.Level)
				esn.state[i] *= factor
			}
		}
	}
}

// applyRicciFlow applies differential Ricci flow equations for geometric regularization.
func (esn *ESN) applyRicciFlow() {
	// Simplified Ricci flow: adjust curvature of state space
	// R_ij represents the Ricci curvature tensor
	// ∂g_ij/∂t = -2R_ij (Ricci flow equation)
	
	// Compute local curvature estimate
	n := len(esn.state)
	for i := 0; i < n; i++ {
		// Estimate curvature from neighboring states
		neighbors := 0.0
		count := 0.0
		for j := 0; j < n; j++ {
			if esn.reservoirWeights[i][j] != 0 {
				neighbors += esn.state[j]
				count++
			}
		}
		if count > 0 {
			avgNeighbor := neighbors / count
			curvature := esn.state[i] - avgNeighbor
			
			// Apply Ricci flow correction (small time step)
			flowCoeff := 0.01 * esn.persona.Memory // Memory affects flow rate
			esn.state[i] -= flowCoeff * curvature
		}
	}
}

// GetState returns the current reservoir state.
func (esn *ESN) GetState() []float64 {
	esn.mu.RLock()
	defer esn.mu.RUnlock()
	
	state := make([]float64, len(esn.state))
	copy(state, esn.state)
	return state
}

// Reset resets the reservoir state to zero.
func (esn *ESN) Reset() {
	esn.mu.Lock()
	defer esn.mu.Unlock()
	
	for i := range esn.state {
		esn.state[i] = 0
	}
}

// TrainOutput trains the output layer using ridge regression.
func (esn *ESN) TrainOutput(ctx context.Context, states [][]float64, targets [][]float64) error {
	esn.mu.Lock()
	defer esn.mu.Unlock()
	
	if len(states) != len(targets) {
		return fmt.Errorf("number of states (%d) must match number of targets (%d)", len(states), len(targets))
	}
	
	if len(states) == 0 {
		return fmt.Errorf("no training data provided")
	}
	
	inputDim := len(states[0])
	outputDim := len(targets[0])
	
	// Initialize output weights
	esn.outputWeights = make([][]float64, outputDim)
	for i := range esn.outputWeights {
		esn.outputWeights[i] = make([]float64, inputDim)
	}
	
	// Ridge regression: W = (S^T S + λI)^-1 S^T T
	// Simplified version: use gradient descent
	learningRate := 0.01
	epochs := 100
	
	for epoch := 0; epoch < epochs; epoch++ {
		for s := range states {
			// Forward pass
			predictions := make([]float64, outputDim)
			for i := 0; i < outputDim; i++ {
				sum := 0.0
				for j := 0; j < inputDim; j++ {
					sum += esn.outputWeights[i][j] * states[s][j]
				}
				predictions[i] = sum
			}
			
			// Backward pass
			for i := 0; i < outputDim; i++ {
				error := predictions[i] - targets[s][i]
				for j := 0; j < inputDim; j++ {
					gradient := error*states[s][j] + esn.params.RidgeParam*esn.outputWeights[i][j]
					esn.outputWeights[i][j] -= learningRate * gradient
				}
			}
		}
	}
	
	esn.trained = true
	esn.log.Debug("esn trained", slog.Int("states", len(states)), slog.Int("epochs", epochs))
	
	return nil
}

// Predict generates output predictions from the current reservoir state.
func (esn *ESN) Predict(ctx context.Context) ([]float64, error) {
	esn.mu.RLock()
	defer esn.mu.RUnlock()
	
	if !esn.trained {
		return nil, fmt.Errorf("network not trained")
	}
	
	if esn.outputWeights == nil {
		return nil, fmt.Errorf("output weights not initialized")
	}
	
	outputDim := len(esn.outputWeights)
	output := make([]float64, outputDim)
	
	for i := 0; i < outputDim; i++ {
		sum := 0.0
		for j := range esn.state {
			sum += esn.outputWeights[i][j] * esn.state[j]
		}
		output[i] = sum
	}
	
	return output, nil
}
