package parser

import (
	"strings"
	"testing"
)

func TestTimeRange_UnconstraintInputIsAccepted(t *testing.T) {
	if err := checkTimeInterval(nil, nil, 10); err != nil {
		t.Errorf("nil-time range should be accepted")
	}
}

func TestTimeRange_LeftSidedConstraintInputIsAccepted(t *testing.T) {
	start := new(float32)
	*start = 5
	if err := checkTimeInterval(start, nil, 10); err != nil {
		t.Errorf("nil-time range should be accepted")
	}
}

func TestTimeRange_RightSidedConstraintInputIsAccepted(t *testing.T) {
	end := new(float32)
	*end = 5
	if err := checkTimeInterval(nil, end, 10); err != nil {
		t.Errorf("nil-time range should be accepted")
	}
}

func TestTimeRange_NegativeStartTimeIsDetected(t *testing.T) {
	start := new(float32)
	*start = -5
	err := checkTimeInterval(start, nil, 10)
	if err == nil {
		t.Errorf("negative start time should not be allowed")
	}
	if !strings.Contains(err.Error(), "start time must be >= 0") {
		t.Errorf("incorrect issue reported: %v", err)
	}
}

func TestTimeRange_EndTimeBiggerThanDurationIsDetected(t *testing.T) {
	end := new(float32)
	*end = 15
	err := checkTimeInterval(nil, end, 10)
	if err == nil {
		t.Errorf("too large end time should not be allowed")
	}
	if !strings.Contains(err.Error(), "end time must be <= scenario duration") {
		t.Errorf("incorrect issue reported: %v", err)
	}
}

func TestTimeRange_StartTimeBiggerThanEndTimeIsDetected(t *testing.T) {
	start := new(float32)
	*start = 5
	end := new(float32)
	*end = 5
	err := checkTimeInterval(start, end, 10)
	if err != nil {
		t.Errorf("having the same start and end time should be allowed")
	}
	*end = 4
	err = checkTimeInterval(start, end, 10)
	if err == nil {
		t.Errorf("end time before start time should be detected")
	}
	if !strings.Contains(err.Error(), "end time must be >= start time") {
		t.Errorf("incorrect issue reported: %v", err)
	}
}

func TestWaveCheck_CorrectWaveDefinitionIsExcepted(t *testing.T) {
	wave := Wave{}
	wave.Max = 20
	wave.Period = 60
	if err := wave.Check(); err != nil {
		t.Errorf("issue reported for valid wave: %v", err)
	}
	wave.Min = new(float32)
	*wave.Min = 10
	if err := wave.Check(); err != nil {
		t.Errorf("issue reported for valid wave: %v", err)
	}
}

func TestWaveCheck_NegativeMinimumIsDetected(t *testing.T) {
	wave := Wave{Min: new(float32)}
	*wave.Min = -1
	if err := wave.Check(); err == nil {
		t.Errorf("negative minimum of wave should be detected")
	}
}

func TestWaveCheck_NegativeMaximumIsDetected(t *testing.T) {
	wave := Wave{Max: -1}
	if err := wave.Check(); err == nil {
		t.Errorf("negative maximum of wave should be detected")
	}
}

func TestWaveCheck_MinGreaterMaxIsDetected(t *testing.T) {
	wave := Wave{Min: new(float32), Max: 10}
	*wave.Min = 20
	if err := wave.Check(); err == nil {
		t.Errorf("minimum > maximium should be detected")
	}
}

func TestWaveCheck_NonPositivePeriodeIsDetected(t *testing.T) {
	wave := Wave{Period: 0}
	if err := wave.Check(); err == nil {
		t.Errorf("period length of 0 should not be allowed")
	}
	wave.Period = -1
	if err := wave.Check(); err == nil {
		t.Errorf("neagtive period length should not be allowed")
	}
}

func TestRateCheck_NoOptionIsDetected(t *testing.T) {
	scenario := Scenario{}
	rate := Rate{}
	if err := rate.Check(&scenario); err == nil {
		t.Errorf("missing rate specification should be detected")
	}
}

func TestRateCheck_MultipleOptionsIsDetected(t *testing.T) {
	scenario := Scenario{}
	rate := Rate{}
	rate.Constant = new(float32)
	*rate.Constant = 10
	rate.Slope = new(float32)
	*rate.Slope = 15
	if err := rate.Check(&scenario); err == nil {
		t.Errorf("multiple rate specifications should be detected")
	}
}

func TestRateCheck_NegativeConstantRateIsDetected(t *testing.T) {
	scenario := Scenario{}
	rate := Rate{}
	rate.Constant = new(float32)
	if err := rate.Check(&scenario); err != nil {
		t.Errorf("vailid constant rate of %v should be fine, but received the error %v", *rate.Constant, err)
	}
	*rate.Constant = -10
	if err := rate.Check(&scenario); err == nil {
		t.Errorf("negative constant rate specification should be detected")
	}
}

func TestRateCheck_NegativeSlopeRateIsDetected(t *testing.T) {
	scenario := Scenario{}
	rate := Rate{}
	rate.Slope = new(float32)
	if err := rate.Check(&scenario); err != nil {
		t.Errorf("vailid constant rate of %v should be fine, but received the error %v", *rate.Slope, err)
	}
	*rate.Slope = -10
	if err := rate.Check(&scenario); err == nil {
		t.Errorf("negative slope rate specification should be detected")
	}
}

func TestRateCheck_InvalidWaveIsDetected(t *testing.T) {
	scenario := Scenario{}
	rate := Rate{}
	rate.Wave = new(Wave)
	if err := rate.Check(&scenario); err == nil {
		t.Errorf("invalid wave specification should be detected")
	}
}

func TestApplication_InvalidNameIsDetected(t *testing.T) {
	scenario := Scenario{}
	app := Application{}
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "application name must match") {
		t.Errorf("missing name was not detected")
	}
	app.Name = "  "
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "application name must match") {
		t.Errorf("missing name was not detected")
	}
	app.Name = "_something_with_underscores_"
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "application name must match") {
		t.Errorf("invalid name was not detected")
	}
}

func TestApplication_NegativeInstanceCounterIsNotAllowed(t *testing.T) {
	scenario := Scenario{}
	app := Application{Name: "test", Instances: new(int), Rate: Rate{Constant: new(float32)}}
	if err := app.Check(&scenario); err != nil {
		t.Errorf("default instance value should be valid, but got error: %v", err)
	}
	*app.Instances = -1
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "number of instances must be >= 0") {
		t.Errorf("negative instance counter was not detected")
	}
}

func TestApplication_NegativeAccountCounterIsNotAllowed(t *testing.T) {
	scenario := Scenario{}
	accounts := 5
	app := Application{Name: "test", Accounts: &accounts, Rate: Rate{Constant: new(float32)}}
	if err := app.Check(&scenario); err != nil {
		t.Errorf("default instance value should be valid, but got error: %v", err)
	}
	*app.Accounts = -1
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "number of accounts") {
		t.Errorf("negative account counter was not detected")
	}
}

func TestApplication_DetectsTimingIssue(t *testing.T) {
	scenario := Scenario{}
	app := Application{
		Name:  "test",
		Rate:  Rate{Constant: new(float32)},
		Start: new(float32),
	}
	if err := app.Check(&scenario); err != nil {
		t.Errorf("default start value should be valid, but got error: %v", err)
	}
	*app.Start = 10
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "end time must be >= start time") {
		t.Errorf("invalid start time was not detected")
	}
}

func TestApplication_DetectsShapeIssue(t *testing.T) {
	scenario := Scenario{}
	app := Application{
		Name: "test",
		Rate: Rate{Constant: new(float32)},
	}
	if err := app.Check(&scenario); err != nil {
		t.Errorf("default start value should be valid, but got error: %v", err)
	}
	*app.Rate.Constant = -10
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "transaction rate must be >= 0") {
		t.Errorf("invalid rate was not detected")
	}
}

func TestNode_InvalidNameIsDetected(t *testing.T) {
	scenario := Scenario{}
	node := Node{}
	if err := node.Check(&scenario); err == nil || !strings.Contains(err.Error(), "node name must match") {
		t.Errorf("missing name was not detected")
	}
	node.Name = "   "
	if err := node.Check(&scenario); err == nil || !strings.Contains(err.Error(), "node name must match") {
		t.Errorf("missing name was not detected")
	}
	node.Name = "_something_with_underscores_"
	if err := node.Check(&scenario); err == nil || !strings.Contains(err.Error(), "node name must match") {
		t.Errorf("missing name was not detected")
	}
}

func TestNode_NegativeInstanceCounterIsNotAllowed(t *testing.T) {
	scenario := Scenario{}
	node := Node{Name: "test", Instances: new(int)}
	if err := node.Check(&scenario); err != nil {
		t.Errorf("default instance value should be valid, but got error: %v", err)
	}
	*node.Instances = -1
	if err := node.Check(&scenario); err == nil || !strings.Contains(err.Error(), "number of instances must be >= 0") {
		t.Errorf("negative instance counter was not detected")
	}
}

func TestNode_DetectsTimingIssue(t *testing.T) {
	scenario := Scenario{}
	node := Node{
		Name:  "test",
		Start: new(float32),
	}
	if err := node.Check(&scenario); err != nil {
		t.Errorf("default start value should be valid, but got error: %v", err)
	}
	*node.Start = 10
	if err := node.Check(&scenario); err == nil || !strings.Contains(err.Error(), "end time must be >= start time") {
		t.Errorf("invalid start time was not detected")
	}
}

func TestScenario_MissingNameIsDetected(t *testing.T) {
	scenario := Scenario{}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "scenario name must not be empty") {
		t.Errorf("missing name was not detected")
	}
	scenario.Name = "  "
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "scenario name must not be empty") {
		t.Errorf("missing name was not detected")
	}
}

func TestScenario_NegativeDurationIsDetected(t *testing.T) {
	scenario := Scenario{Name: "Test"}
	scenario.Duration = -10
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "scenario duration must be > 0") {
		t.Errorf("neagative duration was not detected")
	}
}

func TestScenario_NegativeNumberOfValidatorsIsDetected(t *testing.T) {
	scenario := Scenario{Name: "Test"}
	scenario.NumValidators = new(int)
	*scenario.NumValidators = -5
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "invalid number of validators: -5 <= 0") {
		t.Errorf("neagative number of validators was not detected")
	}
}

func TestScenario_NodeNameCollisionIsDetected(t *testing.T) {
	scenario := Scenario{
		Name:     "Test",
		Duration: 60,
		Nodes: []Node{
			{Name: "A"},
			{Name: "B"},
			{Name: "A"},
		},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "node names must be unique") {
		t.Errorf("node name collision was not detected")
	}
}

func TestScenario_ApplicationNameCollisionIsDetected(t *testing.T) {
	scenario := Scenario{
		Name:     "Test",
		Duration: 60,
		Applications: []Application{
			{Name: "A"},
			{Name: "B"},
			{Name: "A"},
		},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "application names must be unique") {
		t.Errorf("application name collision was not detected")
	}
}

func TestScenario_NodeIssuesAreDetected(t *testing.T) {
	scenario := Scenario{
		Name:     "Test",
		Duration: 60,
		Nodes:    []Node{{}},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "node name must match") {
		t.Errorf("node issue was not detected")
	}
}

func TestScenario_ApplicationIssuesAreDetected(t *testing.T) {
	scenario := Scenario{
		Name:         "Test",
		Duration:     60,
		Applications: []Application{{}},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "application name must match") {
		t.Errorf("application issue was not detected")
	}
}
