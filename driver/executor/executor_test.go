package executor

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/golang/mock/gomock"
)

func TestExecutor_RunEmptyScenario(t *testing.T) {
	ctrl := gomock.NewController(t)
	clock := NewSimClock()
	net := driver.NewMockNetwork(ctrl)
	scenario := parser.Scenario{
		Name:     "Test",
		Duration: 10,
	}

	if err := Run(clock, net, &scenario); err != nil {
		t.Errorf("failed to run empty scenario: %v", err)
	}
	want := Seconds(10)
	if got := clock.Now(); got < want {
		t.Errorf("scenario execution did not complete all steps, expected end time %v, got %v", want, got)
	}
}

func TestExecutor_RunSingleNodeScenario(t *testing.T) {

	clock := NewSimClock()
	scenario := parser.Scenario{
		Name:     "Test",
		Duration: 10,
		Nodes: []parser.Node{{
			Name:  "A",
			Start: New[float32](3),
			End:   New[float32](7),
		}},
	}

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	node := driver.NewMockNode(ctrl)

	// In this scenario, a node is expected to be created and shut down.
	gomock.InOrder(
		net.EXPECT().CreateNode(gomock.Any()).Return(node, nil),
		node.EXPECT().Stop(),
		node.EXPECT().Cleanup(),
	)

	if err := Run(clock, net, &scenario); err != nil {
		t.Errorf("failed to run scenario: %v", err)
	}
	want := Seconds(10)
	if got := clock.Now(); got < want {
		t.Errorf("scenario execution did not complete all steps, expected end time %v, got %v", want, got)
	}
}

func TestExecutor_RunMultipleNodeScenario(t *testing.T) {

	clock := NewSimClock()
	scenario := parser.Scenario{
		Name:     "Test",
		Duration: 10,
		Nodes: []parser.Node{{
			Name:      "A",
			Instances: New(2),
			Start:     New[float32](3),
			End:       New[float32](7),
		}},
	}

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)

	// In this scenario, two nodes are created and stopped.
	gomock.InOrder(
		net.EXPECT().CreateNode(gomock.Any()).Return(node1, nil),
		node1.EXPECT().Stop(),
		node1.EXPECT().Cleanup(),
	)
	gomock.InOrder(
		net.EXPECT().CreateNode(gomock.Any()).Return(node2, nil),
		node2.EXPECT().Stop(),
		node2.EXPECT().Cleanup(),
	)

	if err := Run(clock, net, &scenario); err != nil {
		t.Errorf("failed to run scenario: %v", err)
	}
	want := Seconds(10)
	if got := clock.Now(); got < want {
		t.Errorf("scenario execution did not complete all steps, expected end time %v, got %v", want, got)
	}
}

func TestExecutor_RunSingleApplicationScenario(t *testing.T) {

	clock := NewSimClock()
	scenario := parser.Scenario{
		Name:     "Test",
		Duration: 10,
		Applications: []parser.Application{{
			Name:  "A",
			Start: New[float32](3),
			End:   New[float32](7),
			Rate:  parser.Rate{Constant: New[float32](10)},
		}},
	}

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	app := driver.NewMockApplication(ctrl)

	// In this scenario, an application is expected to be created and shut down.
	net.EXPECT().CreateApplication(gomock.Any()).Return(app, nil)
	app.EXPECT().Start()
	app.EXPECT().Stop()

	if err := Run(clock, net, &scenario); err != nil {
		t.Errorf("failed to run scenario: %v", err)
	}
	want := Seconds(10)
	if got := clock.Now(); got < want {
		t.Errorf("scenario execution did not complete all steps, expected end time %v, got %v", want, got)
	}
}

func TestExecutor_RunMultipleApplicationScenario(t *testing.T) {

	clock := NewSimClock()
	scenario := parser.Scenario{
		Name:     "Test",
		Duration: 10,
		Applications: []parser.Application{{
			Name:      "A",
			Instances: New(2),
			Start:     New[float32](3),
			End:       New[float32](7),
			Rate:      parser.Rate{Constant: New[float32](10)},
		}},
	}

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	app1 := driver.NewMockApplication(ctrl)
	app2 := driver.NewMockApplication(ctrl)

	// In this scenario, an application is expected to be created and shut down.
	net.EXPECT().CreateApplication(gomock.Any()).Return(app1, nil)
	net.EXPECT().CreateApplication(gomock.Any()).Return(app2, nil)
	app1.EXPECT().Start()
	app1.EXPECT().Stop()
	app2.EXPECT().Start()
	app2.EXPECT().Stop()

	if err := Run(clock, net, &scenario); err != nil {
		t.Errorf("failed to run scenario: %v", err)
	}
	want := Seconds(10)
	if got := clock.Now(); got < want {
		t.Errorf("scenario execution did not complete all steps, expected end time %v, got %v", want, got)
	}
}

func New[T any](value T) *T {
	res := new(T)
	*res = value
	return res
}
