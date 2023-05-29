package monitoring

import "testing"

func TestSyncedSeries_CanAddAndRetrieveData(t *testing.T) {
	series := SyncedSeries[Time, int]{}

	if got, want := len(series.GetRange(Time(0), Time(10))), 0; got != want {
		t.Errorf("length of time series not as expected, wanted %d, got %d", got, want)
	}

	if err := series.Append(Time(0), 12); err != nil {
		t.Errorf("error appending data point: %v", err)
	}

	if got, want := len(series.GetRange(Time(0), Time(10))), 1; got != want {
		t.Errorf("length of time series not as expected, wanted %d, got %d", got, want)
	}

	if err := series.Append(Time(5), 8); err != nil {
		t.Errorf("error appending data point: %v", err)
	}

	if got, want := len(series.GetRange(Time(0), Time(10))), 2; got != want {
		t.Errorf("length of time series not as expected, wanted %d, got %d", got, want)
	}

	if err := series.Append(Time(10), 8); err != nil {
		t.Errorf("error appending data point: %v", err)
	}

	if got, want := len(series.GetRange(Time(0), Time(10))), 2; got != want {
		t.Errorf("length of time series not as expected, wanted %d, got %d", got, want)
	}
}

func TestSyncedSeries_AppendFailsIfTimeIsNotProgressing(t *testing.T) {
	series := SyncedSeries[Time, int]{}

	if err := series.Append(Time(10), 12); err != nil {
		t.Errorf("error appending data point: %v", err)
	}

	if err := series.Append(Time(9), 12); err == nil {
		t.Errorf("append of earlier data point should have failed")
	}

	if err := series.Append(Time(10), 12); err == nil {
		t.Errorf("append of existing data point should have failed")
	}

	if err := series.Append(Time(11), 12); err != nil {
		t.Errorf("error appending data point: %v", err)
	}
}

func TestSyncedSeries_GetLatestReturnsLastAppendedValue(t *testing.T) {
	series := SyncedSeries[Time, int]{}

	db := func(t, v int) DataPoint[Time, int] {
		return DataPoint[Time, int]{Time(t), v}
	}

	if got := series.GetLatest(); got != nil {
		t.Errorf("last element should initially be nil")
	}

	if err := series.Append(Time(0), 12); err != nil {
		t.Errorf("error appending data point: %v", err)
	}

	if got, want := series.GetLatest(), db(0, 12); got == nil || *got != want {
		t.Errorf("latest element not as expected, wanted %d, got %d", got, want)
	}

	if err := series.Append(Time(2), 8); err != nil {
		t.Errorf("error appending data point: %v", err)
	}

	if got, want := series.GetLatest(), db(2, 8); got == nil || *got != want {
		t.Errorf("latest element not as expected, wanted %d, got %d", got, want)
	}
}
