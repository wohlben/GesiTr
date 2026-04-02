package models

import (
	"testing"
)

func TestComputeLevel(t *testing.T) {
	tests := []struct {
		xp   float64
		want int
	}{
		{0, 0},
		{50, 0},
		{100, 1},
		{250, 2},
		{999, 9},
		{1000, 10},
		{10000, 100},
		{99999, 100}, // capped
	}
	for _, tt := range tests {
		if got := ComputeLevel(tt.xp); got != tt.want {
			t.Errorf("ComputeLevel(%v) = %v, want %v", tt.xp, got, tt.want)
		}
	}
}

func TestComputeTier(t *testing.T) {
	tests := []struct {
		level int
		want  MasteryTier
	}{
		{0, TierNovice},
		{10, TierNovice},
		{11, TierJourneyman},
		{30, TierJourneyman},
		{31, TierAdept},
		{50, TierAdept},
		{51, TierMaster},
		{99, TierMaster},
		{100, TierMastered},
	}
	for _, tt := range tests {
		if got := ComputeTier(tt.level); got != tt.want {
			t.Errorf("ComputeTier(%v) = %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestComputeRecencyMultiplier(t *testing.T) {
	tests := []struct {
		nDays     int
		baseLevel int
		want      float64
	}{
		{0, 0, 0},        // no days, floor is max(1, 0) = 1, min(0, 1) = 0
		{1, 0, 0.5},      // 0.5*1=0.5, cap=max(1,0)=1, min(0.5,1)=0.5
		{2, 0, 1.0},      // 0.5*2=1.0, cap=1, min(1,1)=1.0
		{10, 0, 1.0},     // 0.5*10=5, cap=1, min(5,1)=1.0
		{1, 2, 0.5},      // 0.5*1=0.5, cap=max(1,1)=1, min(0.5,1)=0.5
		{4, 2, 1.0},      // 0.5*4=2, cap=max(1,1)=1, min(2,1)=1.0
		{10, 10, 5.0},    // 0.5*10=5, cap=max(1,5)=5, min(5,5)=5.0
		{100, 10, 5.0},   // 0.5*100=50, cap=5, min(50,5)=5.0
		{100, 30, 15.0},  // 0.5*100=50, cap=15, min(50,15)=15.0
		{100, 100, 50.0}, // 0.5*100=50, cap=50, min(50,50)=50.0
	}
	for _, tt := range tests {
		if got := ComputeRecencyMultiplier(tt.nDays, tt.baseLevel); got != tt.want {
			t.Errorf("ComputeRecencyMultiplier(%v, %v) = %v, want %v", tt.nDays, tt.baseLevel, got, tt.want)
		}
	}
}

func TestComputeProgress(t *testing.T) {
	tests := []struct {
		xp    float64
		level int
		want  float64
	}{
		{0, 0, 0.0},
		{50, 0, 0.5},
		{150, 1, 0.5},
		{10000, 100, 1.0}, // maxed out
	}
	for _, tt := range tests {
		if got := ComputeProgress(tt.xp, tt.level); got != tt.want {
			t.Errorf("ComputeProgress(%v, %v) = %v, want %v", tt.xp, tt.level, got, tt.want)
		}
	}
}

func TestRelationshipTypeBonus(t *testing.T) {
	// equivalent
	bonus, ok := RelationshipTypeBonus("equivalent")
	if !ok || bonus != 0.5 {
		t.Errorf("equivalent: got (%v, %v), want (0.5, true)", bonus, ok)
	}

	// skill transfer
	bonus, ok = RelationshipTypeBonus("alternative")
	if !ok || bonus != 0.25 {
		t.Errorf("alternative: got (%v, %v), want (0.25, true)", bonus, ok)
	}

	// no transfer
	_, ok = RelationshipTypeBonus("antagonist")
	if ok {
		t.Error("antagonist should not contribute")
	}

	_, ok = RelationshipTypeBonus("accessory")
	if ok {
		t.Error("accessory should not contribute")
	}
}

func TestComputeContributionMultiplier(t *testing.T) {
	tests := []struct {
		strength float64
		relType  string
		want     float64
		wantOK   bool
	}{
		{1.0, "equivalent", 1.0, true},           // 0.5 + 0.5
		{1.0, "equipment_variation", 0.75, true}, // 0.5 + 0.25
		{0.8, "variant", 0.65, true},             // 0.4 + 0.25
		{0.5, "alternative", 0.5, true},          // 0.25 + 0.25
		{1.0, "antagonist", 0, false},            // no transfer
		{1.0, "accessory", 0, false},             // no transfer
	}
	for _, tt := range tests {
		got, ok := ComputeContributionMultiplier(tt.strength, tt.relType)
		if ok != tt.wantOK || got != tt.want {
			t.Errorf("ComputeContributionMultiplier(%v, %q) = (%v, %v), want (%v, %v)", tt.strength, tt.relType, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestEquipmentRelationshipTypeBonus(t *testing.T) {
	// equivalent contributes
	bonus, ok := EquipmentRelationshipTypeBonus("equivalent")
	if !ok || bonus != 0.5 {
		t.Errorf("equivalent: got (%v, %v), want (0.5, true)", bonus, ok)
	}

	// forked does not contribute
	_, ok = EquipmentRelationshipTypeBonus("forked")
	if ok {
		t.Error("forked should not contribute")
	}

	// unknown type does not contribute
	_, ok = EquipmentRelationshipTypeBonus("unknown")
	if ok {
		t.Error("unknown should not contribute")
	}
}

func TestComputeEquipmentContributionMultiplier(t *testing.T) {
	tests := []struct {
		strength float64
		relType  string
		want     float64
		wantOK   bool
	}{
		{1.0, "equivalent", 1.0, true},  // 0.5 + 0.5
		{0.8, "equivalent", 0.9, true},  // 0.4 + 0.5
		{0.5, "equivalent", 0.75, true}, // 0.25 + 0.5
		{1.0, "forked", 0, false},       // no transfer
	}
	for _, tt := range tests {
		got, ok := ComputeEquipmentContributionMultiplier(tt.strength, tt.relType)
		if ok != tt.wantOK || got != tt.want {
			t.Errorf("ComputeEquipmentContributionMultiplier(%v, %q) = (%v, %v), want (%v, %v)", tt.strength, tt.relType, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestFulfillmentContributionMultiplier(t *testing.T) {
	got := FulfillmentContributionMultiplier()
	if got != 0.75 {
		t.Errorf("FulfillmentContributionMultiplier() = %v, want 0.75", got)
	}
}
