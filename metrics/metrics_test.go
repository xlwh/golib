/* metrics_test.go - unit test for metrics */
/*
modification history
--------------------
2016/12/19, by Sijie Yang, create
*/
/*
DESCRIPTION
*/
package metrics

import (
	"reflect"
	"testing"
)

import (
	"www.baidu.com/golang-lib/module_state2"
)

type MockState struct {
	ModuleCounter0 *Counter
	ModuleCounter1 *Counter
	ModuleCounter2 *Counter
	ModuleCounter3 *Counter
	ModuleCounter4 *Counter
	ModuleCounter5 *Counter
	ModuleCounter6 *Counter
	ModuleCounter7 *Counter
	ModuleCounter8 *Counter
	ModuleCounter9 *Counter
}

func prepareMetricsState() (*Metrics, *MockState) {
	state := new(MockState)
	metrics := new(Metrics)
	metrics.Init(state, "METRICS", 20)
	return metrics, state
}

func prepareModuleState() module_state2.State {
	counters := []string{
		"MODULE_COUNTER0",
		"MODULE_COUNTER1",
		"MODULE_COUNTER2",
		"MODULE_COUNTER3",
		"MODULE_COUNTER4",
		"MODULE_COUNTER5",
		"MODULE_COUNTER6",
		"MODULE_COUNTER7",
		"MODULE_COUNTER8",
		"MODULE_COUNTER9",
	}
	var s module_state2.State
	s.Init()
	s.CountersInit(counters)
	return s
}

func TestMetricsGetAll(t *testing.T) {
	m, s := prepareMetricsState()
	s.ModuleCounter0.Inc(1)
	s.ModuleCounter2.Inc(1)
	s.ModuleCounter2.Inc(1)
	s.ModuleCounter4.Inc(4)
	s.ModuleCounter8.Inc(-2)

	d := m.GetAll()

	r := NewMetricsData("METRICS", KindTotal)
	r.Data["MODULE_COUNTER0"] = int64(1)
	r.Data["MODULE_COUNTER1"] = int64(0)
	r.Data["MODULE_COUNTER2"] = int64(2)
	r.Data["MODULE_COUNTER3"] = int64(0)
	r.Data["MODULE_COUNTER4"] = int64(4)
	r.Data["MODULE_COUNTER5"] = int64(0)
	r.Data["MODULE_COUNTER6"] = int64(0)
	r.Data["MODULE_COUNTER7"] = int64(0)
	r.Data["MODULE_COUNTER8"] = int64(-2)
	r.Data["MODULE_COUNTER9"] = int64(0)

	if !reflect.DeepEqual(d, r) {
		t.Errorf("GetAll(): expect %v, actual %v", r, d)
	}
}

func TestMetricsGetDiff(t *testing.T) {
	m, s := prepareMetricsState()

	// case 1
	d := m.GetDiff()
	r := NewMetricsData("METRICS", KindDelta)
	r.Data["MODULE_COUNTER0"] = int64(0)
	r.Data["MODULE_COUNTER1"] = int64(0)
	r.Data["MODULE_COUNTER2"] = int64(0)
	r.Data["MODULE_COUNTER3"] = int64(0)
	r.Data["MODULE_COUNTER4"] = int64(0)
	r.Data["MODULE_COUNTER5"] = int64(0)
	r.Data["MODULE_COUNTER6"] = int64(0)
	r.Data["MODULE_COUNTER7"] = int64(0)
	r.Data["MODULE_COUNTER8"] = int64(0)
	r.Data["MODULE_COUNTER9"] = int64(0)

	if !reflect.DeepEqual(d, r) {
		t.Errorf("GetAll(): expect %v, actual %v", r, d)
	}

	// case 2
	s.ModuleCounter0.Inc(1)
	s.ModuleCounter4.Inc(4)
	s.ModuleCounter8.Inc(-2)
	m.updateDiff()
	d = m.GetDiff()
	r = NewMetricsData("METRICS", KindDelta)
	r.Data["MODULE_COUNTER0"] = int64(1)
	r.Data["MODULE_COUNTER1"] = int64(0)
	r.Data["MODULE_COUNTER2"] = int64(0)
	r.Data["MODULE_COUNTER3"] = int64(0)
	r.Data["MODULE_COUNTER4"] = int64(4)
	r.Data["MODULE_COUNTER5"] = int64(0)
	r.Data["MODULE_COUNTER6"] = int64(0)
	r.Data["MODULE_COUNTER7"] = int64(0)
	r.Data["MODULE_COUNTER8"] = int64(-2)
	r.Data["MODULE_COUNTER9"] = int64(0)

	if !reflect.DeepEqual(d, r) {
		t.Errorf("GetAll(): expect %v, actual %v", r, d)
	}
}

type CaseStructA struct {
	c *Counter
}

type CaseStructB struct {
	c Counter
}

type CaseStructC struct {
	c *Counter
	i int64
}

func TestInvalidCounter(t *testing.T) {
	var m Metrics

	// case 1
	var s1 CaseStructA
	if err := m.Init(s1, "METRICS", 20); err == nil {
		t.Errorf("expect error: %s", ErrStructPtrType)
	}

	// case 2
	var s2 CaseStructB
	if err := m.Init(&s2, "METRICS", 20); err == nil {
		t.Errorf("expect error: %s", ErrStructFieldType)
	}

	// case 3
	var s3 CaseStructC
	if err := m.Init(&s3, "METRICS", 20); err == nil {
		t.Errorf("expect error: %s", ErrStructFieldType)
	}
}

func TestMetricsData(t *testing.T) {
	d1 := NewMetricsData("METRIX", KindTotal)
	d2 := NewMetricsData("METRIX", KindTotal)
	de := NewMetricsData("METRIX", KindDelta)

	d1.Data["MODULE_COUNTER1"] = 10
	d1.Data["MODULE_COUNTER2"] = 20
	d2.Data["MODULE_COUNTER2"] = 30
	d2.Data["MODULE_COUNTER3"] = 40
	de.Data["MODULE_COUNTER2"] = 10
	de.Data["MODULE_COUNTER3"] = 40

	d := d2.Diff(d1)
	if !reflect.DeepEqual(de, d) {
		t.Errorf("expect %v, actual %v", de, d)
	}
}

func BenchmarkMetricIncSingle(b *testing.B) {
	_, s := prepareMetricsState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ModuleCounter4.Inc(1)
	}
}

func BenchmarkStateIncSingle(b *testing.B) {
	s := prepareModuleState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Inc("MODULE_COUNTER4", 1)
	}
}

func BenchmarkMetricIncMulti(b *testing.B) {
	_, s := prepareMetricsState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ModuleCounter4.Inc(1)
		s.ModuleCounter9.Inc(1)
		s.ModuleCounter0.Inc(1)
		s.ModuleCounter6.Inc(1)
	}
}

func BenchmarkStateIncMutli(b *testing.B) {
	s := prepareModuleState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Inc("MODULE_COUNTER4", 1)
		s.Inc("MODULE_COUNTER9", 1)
		s.Inc("MODULE_COUNTER0", 1)
		s.Inc("MODULE_COUNTER6", 1)
	}
}

func BenchmarkMetricGet(b *testing.B) {
	m, _ := prepareMetricsState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.GetAll()
	}
}

func BenchmarkStateGet(b *testing.B) {
	s := prepareModuleState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetAll()
	}
}
