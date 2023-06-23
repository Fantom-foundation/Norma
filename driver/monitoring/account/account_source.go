package accountmon

import (
	"log"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
)

// SensorFactory is a factory for sensors targeting selected applications.
type SensorFactory[T any] interface {
	CreateSensor(driver.Application, int) (utils.Sensor[T], error)
}

// periodicAccountDataSource is a generic data source periodically querying
// account-associated sensors for data.
type periodicAccountDataSource[T any] struct {
	*utils.PeriodicDataSource[mon.Account, T]
	factory SensorFactory[T]
}

// NewPeriodicAccountDataSource creates a new data source managing per-account
// sensor instances for a given metric and periodically collecting data from those.
func NewPeriodicAccountDataSource[T any](
	metric mon.Metric[mon.Account, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	factory SensorFactory[T],
) mon.Source[mon.Account, mon.Series[mon.Time, T]] {
	return newPeriodicAccountDataSource(metric, monitor, time.Second, factory)
}

// newPeriodicAccountDataSource is the same as NewPeriodicAccountDataSource but with
// a customizable sampling periode.
func newPeriodicAccountDataSource[T any](
	metric mon.Metric[mon.Account, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	period time.Duration,
	factory SensorFactory[T],
) mon.Source[mon.Account, mon.Series[mon.Time, T]] {
	res := &periodicAccountDataSource[T]{
		PeriodicDataSource: utils.NewPeriodicDataSourceWithPeriod(metric, monitor, period),
		factory:            factory,
	}

	monitor.Network().RegisterListener(res)
	for _, app := range monitor.Network().GetActiveApplications() {
		res.AfterApplicationCreation(app)
	}

	return res
}

func (s *periodicAccountDataSource[T]) AfterNodeCreation(driver.Node) {
	// ignored
}

func (s *periodicAccountDataSource[T]) AfterApplicationCreation(app driver.Application) {
	label := mon.App(app.Config().Name)
	for i := 0; i < app.Config().Accounts; i++ {
		sensor, err := s.factory.CreateSensor(app, i)
		if err != nil {
			log.Printf("failed to create sensor for metric %v / app %s / account %d: %v", s.GetMetric().Name, label, i, err)
			return
		}
		s.AddSubject(mon.Account{
			App: label,
			Id:  i,
		}, sensor)
	}
}
