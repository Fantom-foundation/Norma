// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"
	"os"
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

func TestAutoCheck_DefaultValueIsValid(t *testing.T) {
	auto := Auto{}
	if err := auto.Check(); err != nil {
		t.Errorf("issue reported for valid auto-shape: %v", err)
	}
}

func TestAutoCheck_NegativeIncreaseIsDetected(t *testing.T) {
	auto := Auto{Increase: new(float32)}
	if err := auto.Check(); err == nil {
		t.Errorf("zero increase rate should be detected")
	}
	*auto.Increase = -10
	if err := auto.Check(); err == nil {
		t.Errorf("negative increase rate should be detected")
	}
}

func TestAutoCheck_InvalidDecreaseRateIsDetected(t *testing.T) {
	auto := Auto{Decrease: new(float32)}
	*auto.Decrease = 0
	if err := auto.Check(); err != nil {
		t.Errorf("zero decrease ratio should be fine")
	}
	*auto.Decrease = 1
	if err := auto.Check(); err != nil {
		t.Errorf("100%% decrease ratio should be fine")
	}
	*auto.Decrease = -0.1
	if err := auto.Check(); err == nil {
		t.Errorf("negative decrease rate should be detected")
	}
	*auto.Decrease = 1.1
	if err := auto.Check(); err == nil {
		t.Errorf(">100%% decrease rate should be detected")
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
		t.Errorf("negative period length should not be allowed")
	}
}

func TestSlopeCheck_NegativeStartRateIsDetected(t *testing.T) {
	slope := Slope{Start: -1}
	if err := slope.Check(); err == nil {
		t.Errorf("negative slope start rate should not be allowed")
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
	rate.Slope = new(Slope)
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

func TestRateCheck_InvalidSlopeRateIsDetected(t *testing.T) {
	scenario := Scenario{}
	rate := Rate{}
	rate.Slope = new(Slope)
	if err := rate.Check(&scenario); err != nil {
		t.Errorf("valid slope of %v should be fine, but received the error %v", *rate.Slope, err)
	}
	rate.Slope.Start = -10
	if err := rate.Check(&scenario); err == nil {
		t.Errorf("invalid slope specification should be detected")
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

func TestApplication_InvalidApplicationTypeIsDetected(t *testing.T) {
	scenario := Scenario{}
	app := Application{}
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "application type must be specified") {
		t.Errorf("missing type was not detected")
	}
	app.Type = "something_that_will_hopefully_never_exist"
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "unknown application type") {
		t.Errorf("invalid type was not detected")
	}
}

func TestApplication_NegativeInstanceCounterIsNotAllowed(t *testing.T) {
	scenario := Scenario{}
	app := Application{Name: "test", Type: "counter", Instances: new(int), Rate: Rate{Constant: new(float32)}}
	if err := app.Check(&scenario); err != nil {
		t.Errorf("default instance value should be valid, but got error: %v", err)
	}
	*app.Instances = -1
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "number of instances must be >= 0") {
		t.Errorf("negative instance counter was not detected")
	}
}

func TestApplication_NegativeUserCounterIsNotAllowed(t *testing.T) {
	scenario := Scenario{}
	users := 5
	app := Application{Name: "test", Type: "counter", Users: &users, Rate: Rate{Constant: new(float32)}}
	if err := app.Check(&scenario); err != nil {
		t.Errorf("default instance value should be valid, but got error: %v", err)
	}
	*app.Users = -1
	if err := app.Check(&scenario); err == nil || !strings.Contains(err.Error(), "number of users") {
		t.Errorf("negative user counter was not detected")
	}
}

func TestApplication_DetectsTimingIssue(t *testing.T) {
	scenario := Scenario{}
	app := Application{
		Name:  "test",
		Type:  "counter",
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
		Type: "counter",
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
		t.Errorf("negative duration was not detected")
	}
}

func TestScenario_NegativeNumberOfValidatorsIsDetected(t *testing.T) {
	scenario := Scenario{Name: "Test"}
	scenario.NumValidators = new(int)
	*scenario.NumValidators = -5
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "invalid number of validators: -5 <= 0") {
		t.Errorf("negative number of validators was not detected")
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

func TestScenario_CheatIssuesAreDetected(t *testing.T) {
	start := new(float32)
	*start = 70

	scenario := Scenario{
		Name:     "Test",
		Duration: 60,
		Cheats: []Cheat{
			{Name: "Test", Start: start},
		},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "start time must be <= scenario duration") {
		fmt.Println(err)
		t.Errorf("cheat issue was not detected")
	}
}

func TestScenario_NodeGenesisImportIssuesAreDetected(t *testing.T) {
	scenario := Scenario{
		Name:     "Test",
		Duration: 60,
		Nodes: []Node{
			{Genesis: Genesis{Import: "/does/not/exist.g"}},
		},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "provided genesis file does not exist") {
		t.Errorf("genesis does not exist but issue was not detected")
	}

	scenario = Scenario{
		Name:     "Test",
		Duration: 60,
		Nodes: []Node{
			{Genesis: Genesis{Import: "/does/exist.notg"}},
		},
	}
	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "provided path is not a genesis file") {
		t.Errorf("targeted file is not a genesis but issue was not detected")
	}
}

func TestScenario_NodeGenesisExportIssuesAreDetected(t *testing.T) {
	f, _ := os.Create("file_exists.g")

	scenario := Scenario{
		Name:     "Test",
		Duration: 60,
		Nodes: []Node{
			{Genesis: Genesis{Export: "file_exists.g"}},
		},
	}

	if err := scenario.Check(); err == nil || !strings.Contains(err.Error(), "provided genesis file already exists") {
		t.Errorf("genesis exists but issue was not detected")
	}

	f.Close()
	_ = os.Remove("file_exists.g")
}
