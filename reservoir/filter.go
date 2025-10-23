// Package reservoir - Integration with mail filtering system
package reservoir

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/mjl-/mox/message"
	"github.com/mjl-/mox/mlog"
)

// FilterConfig contains configuration for the reservoir-enhanced filter.
type FilterConfig struct {
	ESNParams ESNParams    `sconf:"optional" sconf-doc:"Echo State Network parameters for reservoir computing."`
	Persona   PersonaTrait `sconf:"optional" sconf-doc:"Personality traits for affective computing."`
	
	// Integration parameters
	EnableReservoir bool    `sconf:"optional" sconf-doc:"Enable reservoir computing enhancement. Default: false."`
	EnableAffective bool    `sconf:"optional" sconf-doc:"Enable affective computing. Default: false."`
	ReservoirWeight float64 `sconf:"optional" sconf-doc:"Weight of reservoir prediction (0-1). Default: 0.3."`
	
	// Membrane computing
	MembraneDepth int `sconf:"optional" sconf-doc:"Depth of P-system membrane hierarchy. Default: 3."`
}

// DefaultFilterConfig returns default configuration.
func DefaultFilterConfig() FilterConfig {
	return FilterConfig{
		ESNParams:       DefaultESNParams(),
		Persona:         DefaultPersonaTrait(),
		EnableReservoir: false,
		EnableAffective: false,
		ReservoirWeight: 0.3,
		MembraneDepth:   3,
	}
}

// DefaultPersonaTrait returns a balanced default persona.
func DefaultPersonaTrait() PersonaTrait {
	return PersonaTrait{
		Valence:    0.2,  // Slightly positive
		Arousal:    0.6,  // Moderately alert
		Dominance:  0.5,  // Neutral control
		Attention:  0.8,  // High attention
		Memory:     0.7,  // Good memory
		Creativity: 0.5,  // Moderate creativity
	}
}

// ReservoirFilter combines traditional Bayesian filtering with reservoir computing.
type ReservoirFilter struct {
	config FilterConfig
	log    mlog.Log
	
	// Reservoir computing components
	esn             *ESN
	affectiveAgent  *AffectiveAgent
	membraneSystem  *MembraneSystem
	
	// Statistics
	messagesProcessed int
	reservoirEnabled  bool
}

// NewReservoirFilter creates a new reservoir-enhanced filter.
func NewReservoirFilter(log mlog.Log, config FilterConfig) (*ReservoirFilter, error) {
	rf := &ReservoirFilter{
		config:            config,
		log:               log,
		messagesProcessed: 0,
		reservoirEnabled:  config.EnableReservoir,
	}
	
	if config.EnableReservoir {
		// Initialize ESN
		esn, err := NewESN(log, config.ESNParams, config.Persona)
		if err != nil {
			return nil, fmt.Errorf("creating ESN: %w", err)
		}
		rf.esn = esn
		
		// Initialize membrane system
		rf.membraneSystem = NewMembraneSystem(config.MembraneDepth)
		
		log.Debug("reservoir computing initialized", 
			slog.Int("reservoir_size", config.ESNParams.ReservoirSize),
			slog.Int("membrane_depth", config.MembraneDepth))
	}
	
	if config.EnableAffective {
		// Initialize affective agent
		rf.affectiveAgent = NewAffectiveAgent(config.Persona)
		log.Debug("affective computing initialized")
	}
	
	return rf, nil
}

// ClassifyResult contains classification results from the reservoir filter.
type ClassifyResult struct {
	BayesianProb    float64 // Probability from Bayesian filter
	ReservoirProb   float64 // Probability from reservoir computing
	AffectiveProb   float64 // Probability from affective analysis
	CombinedProb    float64 // Combined probability
	AffectiveState  *AffectiveState // Emotional state (if affective enabled)
	MembraneObjects []Object // Objects from membrane processing
}

// ClassifyMessage classifies a message using reservoir computing enhancement.
func (rf *ReservoirFilter) ClassifyMessage(ctx context.Context, part *message.Part, bayesianProb float64) (*ClassifyResult, error) {
	result := &ClassifyResult{
		BayesianProb: bayesianProb,
		CombinedProb: bayesianProb, // Default to Bayesian if reservoir disabled
	}
	
	// Extract text content from message
	content := rf.extractTextContent(part)
	
	// Affective analysis
	if rf.config.EnableAffective && rf.affectiveAgent != nil {
		state := rf.affectiveAgent.ProcessMessage(ctx, content)
		result.AffectiveState = &state
		result.AffectiveProb = rf.affectiveAgent.GetSpamProbability()
		
		rf.log.Debug("affective analysis", 
			slog.Float64("spam_prob", result.AffectiveProb),
			slog.Float64("valence", state.Valence),
			slog.Float64("arousal", state.Arousal))
	}
	
	// Reservoir computing analysis
	if rf.config.EnableReservoir && rf.esn != nil {
		// Convert content to feature vector
		features := rf.extractFeatures(content)
		
		// Update ESN with features
		if err := rf.esn.Update(ctx, features); err != nil {
			return nil, fmt.Errorf("updating ESN: %w", err)
		}
		
		// Process through membrane system
		if rf.membraneSystem != nil {
			if err := rf.processMembraneSystem(content); err != nil {
				rf.log.Debug("membrane processing warning", slog.Any("err", err))
			}
			result.MembraneObjects = rf.membraneSystem.CollectResults()
		}
		
		// Get prediction (if trained)
		if rf.esn.trained {
			prediction, err := rf.esn.Predict(ctx)
			if err == nil && len(prediction) > 0 {
				// First output is spam probability
				result.ReservoirProb = sigmoid(prediction[0])
				
				rf.log.Debug("reservoir prediction", 
					slog.Float64("spam_prob", result.ReservoirProb))
			}
		} else {
			// Not trained, use a heuristic based on reservoir state
			state := rf.esn.GetState()
			activation := 0.0
			for _, s := range state {
				activation += s * s
			}
			activation = activation / float64(len(state))
			result.ReservoirProb = sigmoid(activation)
		}
	}
	
	// Combine predictions
	result.CombinedProb = rf.combinePredictions(result)
	
	rf.messagesProcessed++
	
	return result, nil
}

// extractTextContent extracts text content from a message part.
func (rf *ReservoirFilter) extractTextContent(part *message.Part) string {
	var content strings.Builder
	
	// Extract subject
	if part.Envelope != nil && part.Envelope.Subject != "" {
		content.WriteString(part.Envelope.Subject)
		content.WriteString(" ")
	}
	
	// Extract body (simplified - reads up to 1MB)
	// For production, would need full MIME handling and charset decoding
	reader := part.Reader()
	if reader != nil {
		buf := make([]byte, 1024*1024) // 1MB limit
		n, _ := reader.Read(buf)
		if n > 0 {
			content.Write(buf[:n])
		}
	}
	
	return content.String()
}

// extractFeatures extracts feature vector from text content.
func (rf *ReservoirFilter) extractFeatures(content string) []float64 {
	features := make([]float64, 10) // Fixed-size feature vector
	
	lower := strings.ToLower(content)
	
	// Feature 0: Length (normalized)
	features[0] = math.Min(float64(len(content))/1000.0, 1.0)
	
	// Feature 1: Uppercase ratio
	upperCount := 0
	for _, r := range content {
		if r >= 'A' && r <= 'Z' {
			upperCount++
		}
	}
	if len(content) > 0 {
		features[1] = float64(upperCount) / float64(len(content))
	}
	
	// Feature 2: Digit ratio
	digitCount := 0
	for _, r := range content {
		if r >= '0' && r <= '9' {
			digitCount++
		}
	}
	if len(content) > 0 {
		features[2] = float64(digitCount) / float64(len(content))
	}
	
	// Feature 3: Special character ratio
	specialCount := strings.Count(content, "!") + strings.Count(content, "$") + strings.Count(content, "%")
	if len(content) > 0 {
		features[3] = float64(specialCount) / float64(len(content)) * 10.0
	}
	
	// Feature 4-9: Spam keyword indicators
	spamKeywords := [][]string{
		{"free", "buy", "click"},
		{"urgent", "act now", "limited"},
		{"winner", "prize", "congratulations"},
		{"viagra", "pharmacy", "pills"},
		{"loan", "credit", "debt"},
		{"make money", "work from home", "earn"},
	}
	
	for i, keywords := range spamKeywords {
		count := 0.0
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				count++
			}
		}
		features[4+i] = count
	}
	
	return features
}

// processMembraneSystem processes content through the membrane system.
func (rf *ReservoirFilter) processMembraneSystem(content string) error {
	// Inject objects into root membrane based on content analysis
	lower := strings.ToLower(content)
	
	// Positive signals
	positiveKeywords := []string{"thank", "please", "regards", "sincerely"}
	for _, kw := range positiveKeywords {
		if strings.Contains(lower, kw) {
			obj := Object{
				Type:     "positive_signal",
				Value:    1.0,
				Charge:   1,
				Mobility: 0.7,
			}
			rf.membraneSystem.InjectObject("root", obj)
		}
	}
	
	// Negative signals (spam indicators)
	negativeKeywords := []string{"click here", "buy now", "free money", "act now"}
	for _, kw := range negativeKeywords {
		if strings.Contains(lower, kw) {
			obj := Object{
				Type:     "negative_signal",
				Value:    1.5,
				Charge:   -1,
				Mobility: 0.9,
			}
			rf.membraneSystem.InjectObject("root", obj)
		}
	}
	
	// Perform evolution steps
	for i := 0; i < 3; i++ {
		if err := rf.membraneSystem.Step(); err != nil {
			return err
		}
	}
	
	return nil
}

// combinePredictions combines predictions from different sources.
func (rf *ReservoirFilter) combinePredictions(result *ClassifyResult) float64 {
	// Start with Bayesian
	combined := result.BayesianProb
	
	// Add reservoir if enabled
	if rf.config.EnableReservoir && result.ReservoirProb > 0 {
		// Weighted combination
		w := rf.config.ReservoirWeight
		combined = (1-w)*result.BayesianProb + w*result.ReservoirProb
	}
	
	// Add affective if enabled
	if rf.config.EnableAffective && result.AffectiveProb > 0 {
		// Affective gets small weight
		combined = 0.8*combined + 0.2*result.AffectiveProb
	}
	
	return combined
}

// sigmoid applies sigmoid function.
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// GetStats returns statistics about the filter.
func (rf *ReservoirFilter) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"messages_processed": rf.messagesProcessed,
		"reservoir_enabled":  rf.reservoirEnabled,
	}
	
	if rf.esn != nil {
		stats["esn_trained"] = rf.esn.trained
		stats["reservoir_size"] = rf.config.ESNParams.ReservoirSize
	}
	
	if rf.membraneSystem != nil {
		stats["membrane_depth"] = rf.config.MembraneDepth
		stats["membrane_steps"] = rf.membraneSystem.StepCount
	}
	
	return stats
}
