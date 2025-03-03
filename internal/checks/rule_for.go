package checks

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/model"

	"github.com/cloudflare/pint/internal/discovery"
	"github.com/cloudflare/pint/internal/output"
	"github.com/cloudflare/pint/internal/parser"
)

const (
	RuleForCheckName = "rule/for"
)

func NewRuleForCheck(minFor, maxFor time.Duration, severity Severity) RuleForCheck {
	return RuleForCheck{
		minFor:   minFor,
		maxFor:   maxFor,
		severity: severity,
	}
}

type RuleForCheck struct {
	severity Severity
	minFor   time.Duration
	maxFor   time.Duration
}

func (c RuleForCheck) Meta() CheckMeta {
	return CheckMeta{IsOnline: true}
}

func (c RuleForCheck) String() string {
	return fmt.Sprintf("%s(%s:%s)", RuleForCheckName, output.HumanizeDuration(c.minFor), output.HumanizeDuration(c.maxFor))
}

func (c RuleForCheck) Reporter() string {
	return RuleForCheckName
}

func (c RuleForCheck) Check(ctx context.Context, path string, rule parser.Rule, entries []discovery.Entry) (problems []Problem) {
	if rule.AlertingRule == nil {
		return nil
	}

	var forDur model.Duration
	var fragment string
	var lines []int
	if rule.AlertingRule.For != nil {
		forDur, _ = model.ParseDuration(rule.AlertingRule.For.Value.Value)
		fragment = rule.AlertingRule.For.Value.Value
		lines = rule.AlertingRule.For.Lines()
	}
	if fragment == "" {
		fragment = rule.AlertingRule.Alert.Value.Value
		lines = rule.AlertingRule.Alert.Lines()
	}

	if time.Duration(forDur) < c.minFor {
		problems = append(problems, Problem{
			Fragment: fragment,
			Lines:    lines,
			Reporter: c.Reporter(),
			Text:     fmt.Sprintf("this alert rule must have a 'for' field with a minimum duration of %s", output.HumanizeDuration(c.minFor)),
			Severity: c.severity,
		})
	}

	if c.maxFor > 0 && time.Duration(forDur) > c.maxFor {
		problems = append(problems, Problem{
			Fragment: fragment,
			Lines:    lines,
			Reporter: c.Reporter(),
			Text:     fmt.Sprintf("this alert rule must have a 'for' field with a maximum duration of %s", output.HumanizeDuration(c.maxFor)),
			Severity: c.severity,
		})
	}

	return problems
}
