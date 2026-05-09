package labeler

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RuleFile is the structure of a YAML file containing labeling rules.
type RuleFile struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads a YAML file and returns the parsed rules.
// The file is expected to have the top-level key "rules".
func LoadRules(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("labeler: read rules file: %w", err)
	}
	var rf RuleFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("labeler: parse rules file: %w", err)
	}
	for i, r := range rf.Rules {
		if r.Label == "" {
			return nil, fmt.Errorf("labeler: rule %d has empty label", i)
		}
	}
	return rf.Rules, nil
}

// SaveRules writes rules to a YAML file, overwriting any existing content.
func SaveRules(path string, rules []Rule) error {
	rf := RuleFile{Rules: rules}
	data, err := yaml.Marshal(&rf)
	if err != nil {
		return fmt.Errorf("labeler: marshal rules: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("labeler: write rules file: %w", err)
	}
	return nil
}
