// Package reservoir - Affective computing and emotional intelligence components
package reservoir

import (
	"context"
	"fmt"
	"math"
	"strings"
)

// EmotionDimension represents a dimension in the Differential Emotion Theory (DET) framework.
type EmotionDimension struct {
	Name      string  // Dimension name (e.g., "joy", "anger", "fear")
	Value     float64 // Current value (-1 to 1)
	Intensity float64 // Intensity/magnitude (0 to 1)
}

// AffectiveState represents the complete affective/emotional state.
type AffectiveState struct {
	// Primary emotions (Differential Emotion Theory)
	Joy      float64 // Positive emotion
	Sadness  float64 // Negative emotion
	Anger    float64 // Hostile emotion
	Fear     float64 // Anxious emotion
	Disgust  float64 // Aversive emotion
	Interest float64 // Engagement emotion
	Surprise float64 // Unexpected emotion
	
	// PAD model dimensions
	Valence   float64 // Positive/negative (-1 to 1)
	Arousal   float64 // Activation level (0 to 1)
	Dominance float64 // Control/power (0 to 1)
	
	// Cognitive dimensions
	Attention   float64 // Focus level (0 to 1)
	Complexity  float64 // Cognitive load (0 to 1)
	Uncertainty float64 // Ambiguity level (0 to 1)
}

// DefaultAffectiveState returns a neutral affective state.
func DefaultAffectiveState() AffectiveState {
	return AffectiveState{
		Joy:         0.5,
		Sadness:     0.0,
		Anger:       0.0,
		Fear:        0.0,
		Disgust:     0.0,
		Interest:    0.5,
		Surprise:    0.0,
		Valence:     0.0,
		Arousal:     0.5,
		Dominance:   0.5,
		Attention:   0.7,
		Complexity:  0.5,
		Uncertainty: 0.3,
	}
}

// AffectiveAgent represents an agent with emotional intelligence and personality.
type AffectiveAgent struct {
	Persona      PersonaTrait   // Base personality traits
	CurrentState AffectiveState // Current emotional state
	History      []AffectiveState // History of states
}

// NewAffectiveAgent creates a new affective agent with the given persona.
func NewAffectiveAgent(persona PersonaTrait) *AffectiveAgent {
	state := DefaultAffectiveState()
	
	// Initialize state based on persona
	state.Valence = persona.Valence
	state.Arousal = persona.Arousal
	state.Dominance = persona.Dominance
	state.Attention = persona.Attention
	
	return &AffectiveAgent{
		Persona:      persona,
		CurrentState: state,
		History:      make([]AffectiveState, 0),
	}
}

// ProcessMessage analyzes a message and updates affective state.
func (aa *AffectiveAgent) ProcessMessage(ctx context.Context, content string) AffectiveState {
	// Save current state to history
	aa.History = append(aa.History, aa.CurrentState)
	
	// Analyze emotional content
	emotions := aa.analyzeEmotionalContent(content)
	
	// Update state based on analysis
	aa.updateState(emotions)
	
	return aa.CurrentState
}

// analyzeEmotionalContent performs basic emotional content analysis.
func (aa *AffectiveAgent) analyzeEmotionalContent(content string) map[string]float64 {
	emotions := make(map[string]float64)
	
	// Convert to lowercase for matching
	lower := strings.ToLower(content)
	
	// Simple keyword-based emotion detection
	// In production, this would use more sophisticated NLP
	
	// Joy indicators
	joyWords := []string{"happy", "joy", "great", "excellent", "wonderful", "love", "pleased", "delighted"}
	emotions["joy"] = aa.countKeywords(lower, joyWords) * 0.1
	
	// Sadness indicators
	sadnessWords := []string{"sad", "unhappy", "disappointed", "unfortunate", "regret", "sorry"}
	emotions["sadness"] = aa.countKeywords(lower, sadnessWords) * 0.1
	
	// Anger indicators
	angerWords := []string{"angry", "furious", "outraged", "mad", "annoyed", "frustrated", "hate"}
	emotions["anger"] = aa.countKeywords(lower, angerWords) * 0.15
	
	// Fear indicators
	fearWords := []string{"afraid", "scared", "worried", "anxious", "nervous", "concerned", "fear"}
	emotions["fear"] = aa.countKeywords(lower, fearWords) * 0.1
	
	// Disgust indicators
	disgustWords := []string{"disgusting", "revolting", "nasty", "awful", "terrible", "horrible"}
	emotions["disgust"] = aa.countKeywords(lower, disgustWords) * 0.15
	
	// Interest indicators
	interestWords := []string{"interesting", "curious", "wonder", "question", "inquiry", "explore"}
	emotions["interest"] = aa.countKeywords(lower, interestWords) * 0.08
	
	// Surprise indicators
	surpriseWords := []string{"surprise", "unexpected", "amazing", "astonishing", "shocking", "wow"}
	emotions["surprise"] = aa.countKeywords(lower, surpriseWords) * 0.1
	
	// Spam indicators (treated as disgust/anger)
	spamWords := []string{"click here", "buy now", "free", "urgent", "limited time", "act now", "winner"}
	spamScore := aa.countKeywords(lower, spamWords) * 0.2
	emotions["disgust"] += spamScore
	emotions["anger"] += spamScore * 0.5
	
	return emotions
}

// countKeywords counts occurrences of keywords in text.
func (aa *AffectiveAgent) countKeywords(text string, keywords []string) float64 {
	count := 0.0
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			count++
		}
	}
	return count
}

// updateState updates the affective state based on emotional analysis.
func (aa *AffectiveAgent) updateState(emotions map[string]float64) {
	// Decay rate for temporal dynamics
	decayRate := 0.1
	
	// Update primary emotions with momentum
	aa.CurrentState.Joy = (1-decayRate)*aa.CurrentState.Joy + decayRate*emotions["joy"]
	aa.CurrentState.Sadness = (1-decayRate)*aa.CurrentState.Sadness + decayRate*emotions["sadness"]
	aa.CurrentState.Anger = (1-decayRate)*aa.CurrentState.Anger + decayRate*emotions["anger"]
	aa.CurrentState.Fear = (1-decayRate)*aa.CurrentState.Fear + decayRate*emotions["fear"]
	aa.CurrentState.Disgust = (1-decayRate)*aa.CurrentState.Disgust + decayRate*emotions["disgust"]
	aa.CurrentState.Interest = (1-decayRate)*aa.CurrentState.Interest + decayRate*emotions["interest"]
	aa.CurrentState.Surprise = (1-decayRate)*aa.CurrentState.Surprise + decayRate*emotions["surprise"]
	
	// Compute PAD dimensions from primary emotions
	aa.computePADDimensions()
	
	// Update cognitive dimensions
	aa.updateCognitiveDimensions()
	
	// Clamp values to valid ranges
	aa.clampState()
}

// computePADDimensions computes Valence-Arousal-Dominance from primary emotions.
func (aa *AffectiveAgent) computePADDimensions() {
	// Valence: positive vs negative
	positive := aa.CurrentState.Joy + aa.CurrentState.Interest
	negative := aa.CurrentState.Sadness + aa.CurrentState.Anger + aa.CurrentState.Fear + aa.CurrentState.Disgust
	aa.CurrentState.Valence = math.Tanh(positive - negative)
	
	// Arousal: activation level
	aa.CurrentState.Arousal = (aa.CurrentState.Anger + aa.CurrentState.Fear + 
		aa.CurrentState.Surprise + aa.CurrentState.Interest) / 4.0
	
	// Dominance: control/power
	aa.CurrentState.Dominance = (aa.CurrentState.Anger + aa.CurrentState.Joy - 
		aa.CurrentState.Fear - aa.CurrentState.Sadness) / 4.0
}

// updateCognitiveDimensions updates cognitive processing dimensions.
func (aa *AffectiveAgent) updateCognitiveDimensions() {
	// Attention is affected by arousal and interest
	aa.CurrentState.Attention = 0.7*aa.Persona.Attention + 0.3*(aa.CurrentState.Arousal+aa.CurrentState.Interest)/2.0
	
	// Uncertainty increases with surprise and fear
	aa.CurrentState.Uncertainty = (aa.CurrentState.Surprise + aa.CurrentState.Fear) / 2.0
}

// clampState ensures all values are in valid ranges.
func (aa *AffectiveAgent) clampState() {
	clamp := func(v, min, max float64) float64 {
		if v < min {
			return min
		}
		if v > max {
			return max
		}
		return v
	}
	
	aa.CurrentState.Joy = clamp(aa.CurrentState.Joy, 0, 1)
	aa.CurrentState.Sadness = clamp(aa.CurrentState.Sadness, 0, 1)
	aa.CurrentState.Anger = clamp(aa.CurrentState.Anger, 0, 1)
	aa.CurrentState.Fear = clamp(aa.CurrentState.Fear, 0, 1)
	aa.CurrentState.Disgust = clamp(aa.CurrentState.Disgust, 0, 1)
	aa.CurrentState.Interest = clamp(aa.CurrentState.Interest, 0, 1)
	aa.CurrentState.Surprise = clamp(aa.CurrentState.Surprise, 0, 1)
	
	aa.CurrentState.Valence = clamp(aa.CurrentState.Valence, -1, 1)
	aa.CurrentState.Arousal = clamp(aa.CurrentState.Arousal, 0, 1)
	aa.CurrentState.Dominance = clamp(aa.CurrentState.Dominance, 0, 1)
	
	aa.CurrentState.Attention = clamp(aa.CurrentState.Attention, 0, 1)
	aa.CurrentState.Complexity = clamp(aa.CurrentState.Complexity, 0, 1)
	aa.CurrentState.Uncertainty = clamp(aa.CurrentState.Uncertainty, 0, 1)
}

// GetSpamProbability computes spam probability based on affective state.
func (aa *AffectiveAgent) GetSpamProbability() float64 {
	// High disgust and anger indicate spam
	spamSignal := (aa.CurrentState.Disgust + aa.CurrentState.Anger) / 2.0
	
	// Low interest and high uncertainty also indicate spam
	spamSignal += (1.0 - aa.CurrentState.Interest) * 0.3
	spamSignal += aa.CurrentState.Uncertainty * 0.2
	
	// Negative valence increases spam probability
	if aa.CurrentState.Valence < 0 {
		spamSignal += math.Abs(aa.CurrentState.Valence) * 0.3
	}
	
	// Normalize to [0, 1]
	probability := math.Tanh(spamSignal)
	return math.Max(0, math.Min(1, probability))
}

// GetEngagementScore computes engagement score based on affective state.
func (aa *AffectiveAgent) GetEngagementScore() float64 {
	// High interest and attention indicate engagement
	engagement := (aa.CurrentState.Interest + aa.CurrentState.Attention) / 2.0
	
	// Positive valence increases engagement
	if aa.CurrentState.Valence > 0 {
		engagement += aa.CurrentState.Valence * 0.3
	}
	
	// Moderate arousal is best for engagement
	optimalArousal := 0.6
	arousalFactor := 1.0 - math.Abs(aa.CurrentState.Arousal-optimalArousal)
	engagement *= arousalFactor
	
	return math.Max(0, math.Min(1, engagement))
}

// GenerateReport generates a human-readable report of the affective state.
func (aa *AffectiveAgent) GenerateReport() string {
	state := aa.CurrentState
	
	report := fmt.Sprintf("Affective State Report:\n")
	report += fmt.Sprintf("  Primary Emotions:\n")
	report += fmt.Sprintf("    Joy:      %.2f\n", state.Joy)
	report += fmt.Sprintf("    Sadness:  %.2f\n", state.Sadness)
	report += fmt.Sprintf("    Anger:    %.2f\n", state.Anger)
	report += fmt.Sprintf("    Fear:     %.2f\n", state.Fear)
	report += fmt.Sprintf("    Disgust:  %.2f\n", state.Disgust)
	report += fmt.Sprintf("    Interest: %.2f\n", state.Interest)
	report += fmt.Sprintf("    Surprise: %.2f\n", state.Surprise)
	report += fmt.Sprintf("  PAD Dimensions:\n")
	report += fmt.Sprintf("    Valence:   %.2f\n", state.Valence)
	report += fmt.Sprintf("    Arousal:   %.2f\n", state.Arousal)
	report += fmt.Sprintf("    Dominance: %.2f\n", state.Dominance)
	report += fmt.Sprintf("  Cognitive:\n")
	report += fmt.Sprintf("    Attention:   %.2f\n", state.Attention)
	report += fmt.Sprintf("    Complexity:  %.2f\n", state.Complexity)
	report += fmt.Sprintf("    Uncertainty: %.2f\n", state.Uncertainty)
	report += fmt.Sprintf("  Derived Metrics:\n")
	report += fmt.Sprintf("    Spam Probability: %.2f\n", aa.GetSpamProbability())
	report += fmt.Sprintf("    Engagement Score: %.2f\n", aa.GetEngagementScore())
	
	return report
}

// ApplyRicciFlowToEmotion applies Ricci flow dynamics to emotional states.
// This provides geometric regularization of the emotional state space.
func (aa *AffectiveAgent) ApplyRicciFlowToEmotion(dt float64) {
	// Ricci flow equation: ∂g/∂t = -2Ric
	// Applied to emotional state manifold
	
	// Compute emotional curvature (simplified)
	// High curvature areas (extreme emotions) flow toward lower curvature (balance)
	
	emotions := []float64{
		aa.CurrentState.Joy,
		aa.CurrentState.Sadness,
		aa.CurrentState.Anger,
		aa.CurrentState.Fear,
		aa.CurrentState.Disgust,
		aa.CurrentState.Interest,
		aa.CurrentState.Surprise,
	}
	
	// Compute mean
	mean := 0.0
	for _, e := range emotions {
		mean += e
	}
	mean /= float64(len(emotions))
	
	// Apply flow toward mean (curvature correction)
	flowRate := dt * 0.1 // Small flow rate
	
	aa.CurrentState.Joy += flowRate * (mean - aa.CurrentState.Joy)
	aa.CurrentState.Sadness += flowRate * (mean - aa.CurrentState.Sadness)
	aa.CurrentState.Anger += flowRate * (mean - aa.CurrentState.Anger)
	aa.CurrentState.Fear += flowRate * (mean - aa.CurrentState.Fear)
	aa.CurrentState.Disgust += flowRate * (mean - aa.CurrentState.Disgust)
	aa.CurrentState.Interest += flowRate * (mean - aa.CurrentState.Interest)
	aa.CurrentState.Surprise += flowRate * (mean - aa.CurrentState.Surprise)
	
	aa.clampState()
}
