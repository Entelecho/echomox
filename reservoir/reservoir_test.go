package reservoir

import (
	"context"
	"testing"

	"github.com/mjl-/mox/mlog"
)

func TestNewESN(t *testing.T) {
	log := mlog.New("test", nil)
	params := DefaultESNParams()
	persona := DefaultPersonaTrait()
	
	esn, err := NewESN(log, params, persona)
	if err != nil {
		t.Fatalf("failed to create ESN: %v", err)
	}
	
	if esn == nil {
		t.Fatal("ESN is nil")
	}
	
	if len(esn.state) != params.ReservoirSize {
		t.Errorf("expected state size %d, got %d", params.ReservoirSize, len(esn.state))
	}
	
	if len(esn.membranes) == 0 {
		t.Error("expected membranes to be initialized")
	}
}

func TestESNUpdate(t *testing.T) {
	log := mlog.New("test", nil)
	params := DefaultESNParams()
	params.ReservoirSize = 50 // Smaller for testing
	persona := DefaultPersonaTrait()
	
	esn, err := NewESN(log, params, persona)
	if err != nil {
		t.Fatalf("failed to create ESN: %v", err)
	}
	
	// Create input
	input := []float64{0.5, 0.3, 0.8, 0.1, 0.9}
	
	// Update state
	err = esn.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to update ESN: %v", err)
	}
	
	// Check state changed
	state := esn.GetState()
	hasNonZero := false
	for _, s := range state {
		if s != 0 {
			hasNonZero = true
			break
		}
	}
	
	if !hasNonZero {
		t.Error("expected some non-zero state values after update")
	}
}

func TestESNReset(t *testing.T) {
	log := mlog.New("test", nil)
	params := DefaultESNParams()
	params.ReservoirSize = 30
	persona := DefaultPersonaTrait()
	
	esn, err := NewESN(log, params, persona)
	if err != nil {
		t.Fatalf("failed to create ESN: %v", err)
	}
	
	// Update with some input
	input := []float64{0.5, 0.3, 0.8}
	err = esn.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to update ESN: %v", err)
	}
	
	// Reset
	esn.Reset()
	
	// Check all state is zero
	state := esn.GetState()
	for i, s := range state {
		if s != 0 {
			t.Errorf("expected state[%d] to be 0 after reset, got %f", i, s)
		}
	}
}

func TestESNTrainOutput(t *testing.T) {
	log := mlog.New("test", nil)
	params := DefaultESNParams()
	params.ReservoirSize = 20
	persona := DefaultPersonaTrait()
	
	esn, err := NewESN(log, params, persona)
	if err != nil {
		t.Fatalf("failed to create ESN: %v", err)
	}
	
	// Create training data
	states := [][]float64{
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{0.5, 0.4, 0.3, 0.2, 0.1},
		{0.3, 0.3, 0.3, 0.3, 0.3},
	}
	
	targets := [][]float64{
		{0.0},
		{1.0},
		{0.5},
	}
	
	// Train
	err = esn.TrainOutput(context.Background(), states, targets)
	if err != nil {
		t.Fatalf("failed to train output: %v", err)
	}
	
	if !esn.trained {
		t.Error("expected ESN to be marked as trained")
	}
}

func TestMembraneSystem(t *testing.T) {
	ms := NewMembraneSystem(3)
	
	if ms.Root == nil {
		t.Fatal("root membrane is nil")
	}
	
	if len(ms.All) == 0 {
		t.Fatal("expected some membranes in system")
	}
	
	// Test object injection
	obj := Object{
		Type:     "test",
		Value:    1.0,
		Charge:   0,
		Mobility: 0.5,
	}
	
	err := ms.InjectObject("root", obj)
	if err != nil {
		t.Fatalf("failed to inject object: %v", err)
	}
	
	if len(ms.Root.Objects) != 1 {
		t.Errorf("expected 1 object in root, got %d", len(ms.Root.Objects))
	}
}

func TestMembraneEvolution(t *testing.T) {
	membrane := NewMembrane("test", 0, 0.5)
	
	// Add objects
	membrane.AddObject(Object{Type: "token", Value: 1.0, Charge: 0, Mobility: 0.5})
	membrane.AddObject(Object{Type: "positive_signal", Value: 0.8, Charge: 1, Mobility: 0.7})
	
	// Add a simple rule
	rule := EvolutionRule{
		Name:        "test_rule",
		Priority:    10,
		InputTypes:  []string{"token"},
		OutputTypes: []string{"processed"},
		Transform: func(objs []Object) []Object {
			return []Object{{
				Type:     "processed",
				Value:    2.0,
				Charge:   0,
				Mobility: 0.5,
			}}
		},
	}
	membrane.AddRule(rule)
	
	// Evolve
	err := membrane.Evolve()
	if err != nil {
		t.Fatalf("failed to evolve: %v", err)
	}
	
	// Check objects changed
	hasProcessed := false
	for _, obj := range membrane.Objects {
		if obj.Type == "processed" {
			hasProcessed = true
			break
		}
	}
	
	if !hasProcessed {
		t.Error("expected processed object after evolution")
	}
}

func TestAffectiveAgent(t *testing.T) {
	persona := DefaultPersonaTrait()
	agent := NewAffectiveAgent(persona)
	
	if agent == nil {
		t.Fatal("agent is nil")
	}
	
	// Test with positive message
	positiveMsg := "Thank you for your wonderful help! I'm very happy with the results."
	state := agent.ProcessMessage(context.Background(), positiveMsg)
	
	if state.Joy <= 0 {
		t.Error("expected some joy in positive message")
	}
	
	if state.Valence <= 0 {
		t.Error("expected positive valence for positive message")
	}
}

func TestAffectiveAgentSpamDetection(t *testing.T) {
	persona := DefaultPersonaTrait()
	agent := NewAffectiveAgent(persona)
	
	// Test with spam message
	spamMsg := "Click here to buy now! Free money! Limited time offer! Act now!"
	agent.ProcessMessage(context.Background(), spamMsg)
	
	spamProb := agent.GetSpamProbability()
	
	if spamProb < 0.2 {
		t.Errorf("expected high spam probability for spam message, got %f", spamProb)
	}
}

func TestReservoirFilter(t *testing.T) {
	log := mlog.New("test", nil)
	config := DefaultFilterConfig()
	config.EnableReservoir = true
	config.EnableAffective = true
	config.ESNParams.ReservoirSize = 30 // Smaller for testing
	
	filter, err := NewReservoirFilter(log, config)
	if err != nil {
		t.Fatalf("failed to create reservoir filter: %v", err)
	}
	
	if filter.esn == nil {
		t.Error("expected ESN to be initialized")
	}
	
	if filter.affectiveAgent == nil {
		t.Error("expected affective agent to be initialized")
	}
}

func TestFilterConfig(t *testing.T) {
	config := DefaultFilterConfig()
	
	if config.ESNParams.ReservoirSize <= 0 {
		t.Error("invalid reservoir size")
	}
	
	if config.ESNParams.SpectralRadius <= 0 || config.ESNParams.SpectralRadius >= 1 {
		t.Error("invalid spectral radius")
	}
	
	if config.ReservoirWeight < 0 || config.ReservoirWeight > 1 {
		t.Error("invalid reservoir weight")
	}
}

func TestPersonaTrait(t *testing.T) {
	persona := DefaultPersonaTrait()
	
	// Check all values are in valid ranges
	if persona.Valence < -1 || persona.Valence > 1 {
		t.Errorf("valence out of range: %f", persona.Valence)
	}
	
	if persona.Arousal < 0 || persona.Arousal > 1 {
		t.Errorf("arousal out of range: %f", persona.Arousal)
	}
	
	if persona.Dominance < 0 || persona.Dominance > 1 {
		t.Errorf("dominance out of range: %f", persona.Dominance)
	}
	
	if persona.Attention < 0 || persona.Attention > 1 {
		t.Errorf("attention out of range: %f", persona.Attention)
	}
	
	if persona.Memory < 0 || persona.Memory > 1 {
		t.Errorf("memory out of range: %f", persona.Memory)
	}
	
	if persona.Creativity < 0 || persona.Creativity > 1 {
		t.Errorf("creativity out of range: %f", persona.Creativity)
	}
}

func TestAffectiveStateReport(t *testing.T) {
	persona := DefaultPersonaTrait()
	agent := NewAffectiveAgent(persona)
	
	// Process a message
	agent.ProcessMessage(context.Background(), "This is a test message.")
	
	// Generate report
	report := agent.GenerateReport()
	
	if len(report) == 0 {
		t.Error("expected non-empty report")
	}
	
	// Check report contains key sections
	expectedSections := []string{
		"Affective State Report",
		"Primary Emotions",
		"PAD Dimensions",
		"Cognitive",
		"Spam Probability",
	}
	
	for _, section := range expectedSections {
		if !contains(report, section) {
			t.Errorf("report missing section: %s", section)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
