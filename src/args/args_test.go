package args

import "testing"

func TestParse_ConfigFlag_Short(t *testing.T) {
	p, err := Parse([]string{"-c", "/path/to/conf.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ConfigPath != "/path/to/conf.yaml" {
		t.Errorf("want ConfigPath %q, got %q", "/path/to/conf.yaml", p.ConfigPath)
	}
}

func TestParse_ConfigFlag_Long(t *testing.T) {
	p, err := Parse([]string{"--config", "/path/to/conf.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ConfigPath != "/path/to/conf.yaml" {
		t.Errorf("want ConfigPath %q, got %q", "/path/to/conf.yaml", p.ConfigPath)
	}
}

func TestParse_ConfigFlag_MissingValue(t *testing.T) {
	_, err := Parse([]string{"-c"})
	if err == nil {
		t.Fatal("expected error for missing -c value, got nil")
	}
}

func TestParse_ConfigFlag_NotSet(t *testing.T) {
	p, err := Parse([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ConfigPath != "" {
		t.Errorf("want empty ConfigPath, got %q", p.ConfigPath)
	}
}

func TestParse_ConfigFlag_WithOtherFlags(t *testing.T) {
	p, err := Parse([]string{"-c", "custom.yaml", "--mode", "lite"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ConfigPath != "custom.yaml" {
		t.Errorf("want ConfigPath %q, got %q", "custom.yaml", p.ConfigPath)
	}
	if p.Mode != ModeLite {
		t.Errorf("want mode lite, got %q", p.Mode)
	}
}
