# Deep Tree Echo State Network Reservoir Computing Framework

This package implements a Deep Tree Echo State Network (ESN) Reservoir Computing Framework integrated with:
- Paun P-System Membrane Computing for reservoir evolution
- Butcher B-Series Rooted Forest Runge-Kutta Ridge Regression
- Julia J-Surface Elementary Differential Ricci Flow Equations
- Differential Emotion Theory Framework for Affective Agency

## Overview

The reservoir computing framework enhances the mail server's spam filtering capabilities by adding:

1. **Echo State Networks (ESN)**: Recurrent neural networks with fixed random connections that act as a reservoir of computational dynamics
2. **Membrane Computing**: P-system hierarchical membranes that process and transform computational objects
3. **Affective Computing**: Emotional intelligence based on Differential Emotion Theory
4. **Geometric Flow**: Ricci flow equations for state space regularization

## Components

### Echo State Network (esn.go)

The ESN is a type of reservoir computing architecture with:
- Fixed random recurrent connections (the "reservoir")
- Trainable output layer
- Leak rate for temporal dynamics
- Spectral radius control for stability

Key parameters:
- `ReservoirSize`: Number of neurons (default: 100)
- `SpectralRadius`: Largest eigenvalue of reservoir matrix (default: 0.95)
- `LeakRate`: State decay rate (default: 0.3)
- `Sparsity`: Connection density (default: 0.1)

### Membrane Computing (membrane.go)

Implements Paun P-Systems with:
- Hierarchical membrane structure
- Objects with type, value, charge, and mobility
- Evolution rules for object transformation
- Object passage between membranes based on permeability

### Affective Computing (affective.go)

Provides emotional intelligence with:
- **Differential Emotion Theory**: Joy, Sadness, Anger, Fear, Disgust, Interest, Surprise
- **PAD Model**: Valence, Arousal, Dominance dimensions
- **Cognitive Dimensions**: Attention, Complexity, Uncertainty
- Spam detection based on emotional signals

### Integration (filter.go)

Combines traditional Bayesian filtering with reservoir computing:
- Weighted combination of Bayesian and reservoir predictions
- Feature extraction from email content
- Membrane-based object processing
- Affective state analysis

## Usage

### Basic Example

```go
import (
    "context"
    "github.com/mjl-/mox/reservoir"
    "github.com/mjl-/mox/mlog"
)

// Create configuration
config := reservoir.DefaultFilterConfig()
config.EnableReservoir = true
config.EnableAffective = true

// Create filter
log := mlog.New("reservoir", nil)
filter, err := reservoir.NewReservoirFilter(log, config)
if err != nil {
    // Handle error
}

// Classify a message
result, err := filter.ClassifyMessage(ctx, messagePart, bayesianProb)
if err != nil {
    // Handle error
}

// Use combined probability
spamProbability := result.CombinedProb
```

### Configuration

The framework can be configured through the `FilterConfig` structure:

```go
config := reservoir.FilterConfig{
    ESNParams: reservoir.ESNParams{
        ReservoirSize:  100,
        SpectralRadius: 0.95,
        LeakRate:       0.3,
        TreeDepth:      3,
    },
    Persona: reservoir.PersonaTrait{
        Valence:    0.2,
        Arousal:    0.6,
        Dominance:  0.5,
        Attention:  0.8,
        Memory:     0.7,
        Creativity: 0.5,
    },
    EnableReservoir: true,
    EnableAffective: true,
    ReservoirWeight: 0.3,
    MembraneDepth:   3,
}
```

### Persona Traits

The `PersonaTrait` structure maps LLM personality characteristics to reservoir hyper-parameters:

- **Valence**: Emotional tone (-1 to 1) - affects signal strength
- **Arousal**: Energy level (0 to 1) - affects membrane permeability
- **Dominance**: Control level (0 to 1) - affects prediction confidence
- **Attention**: Focus (0 to 1) - modulates leak rate
- **Memory**: Retention (0 to 1) - affects Ricci flow rate
- **Creativity**: Novelty (0 to 1) - affects reservoir diversity

## Mathematical Foundations

### Runge-Kutta Integration

State updates use 4th-order Runge-Kutta method for numerical stability:

```
k1 = f(x, t)
k2 = f(x + h*k1/2, t + h/2)
k3 = f(x + h*k2/2, t + h/2)
k4 = f(x + h*k3, t + h)
x(t+h) = x(t) + (h/6)*(k1 + 2*k2 + 2*k3 + k4)
```

### Ricci Flow

Geometric regularization using Ricci flow equation:

```
∂g/∂t = -2*Ric
```

This smooths the state manifold curvature, preventing extreme activation states.

### Membrane Evolution

P-system evolution follows:
1. Apply evolution rules (in priority order)
2. Transform matching objects
3. Pass objects between membranes
4. Check dissolution conditions

## Performance

The reservoir computing framework is designed to:
- Add minimal overhead to existing filtering
- Process messages in < 10ms on modern hardware
- Scale to 1000s of messages per second
- Provide interpretable results through affective states

## Future Enhancements

Potential improvements:
- Online learning for ESN output weights
- Adaptive reservoir topology
- Multi-modal input (text, images, attachments)
- Distributed membrane computing
- Integration with external ML models

## References

- Echo State Networks: Jaeger, H. (2001). "The echo state approach to analysing and training recurrent neural networks"
- Membrane Computing: Păun, G. (2000). "Computing with membranes"
- Differential Emotion Theory: Izard, C. E. (1977). "Human Emotions"
- Ricci Flow: Hamilton, R. S. (1982). "Three-manifolds with positive Ricci curvature"
- Butcher Series: Butcher, J. C. (1963). "Coefficients for the study of Runge-Kutta integration processes"

## License

This package is part of mox and is licensed under the same terms (MIT License).
