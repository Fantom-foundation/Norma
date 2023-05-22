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

func TestSyncedSeries_GetAt(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	series := SyncedSeries[BlockNumber, int]{}
	for _, item := range data {
		_ = series.Append(BlockNumber(item), item)
	}

	if series.Size() != len(data) {
		t.Errorf("sizes do not mathc: %v != %v", series.Size(), len(data))
	}

	for i := 0; i < series.Size(); i++ {
		point := series.GetAt(i)
		if point.Value != data[i] {
			t.Errorf("values do not match: %v != %v", point.Value, data[i])
		}
		if point.Position != BlockNumber(data[i]) {
			t.Errorf("values do not match: %v != %v", point.Position, data[i])
		}
	}
}

func TestSyncedSeries_Size(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	series := SyncedSeries[BlockNumber, int]{}
	for _, item := range data {
		_ = series.Append(BlockNumber(item), item)
	}

	if series.Size() != len(data) {
		t.Errorf("sizes do not mathc: %v != %v", series.Size(), len(data))
	}
}
