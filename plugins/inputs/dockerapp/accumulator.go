package dockerapp

import (
	"time"

	"github.com/influxdata/telegraf"
)

type Metric struct {
	containerid string
	valueType   telegraf.ValueType
	measurement string
	fields      map[string]interface{}
	tags        map[string]string
	t           time.Time
}

type accumulator struct {
	containerid string
	metrics     chan Metric
	errors      chan error
	precision   time.Duration
}

func NewAccumulator(containerid string, metrics chan Metric, errors chan error) *accumulator {
	acc := accumulator{
		containerid: containerid,
		metrics:     metrics,
		errors:      errors,
		precision:   time.Nanosecond,
	}
	return &acc
}

func (ac *accumulator) MakeMetric(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	valueType telegraf.ValueType,
	t time.Time,
) Metric {
	return Metric{
		containerid: ac.containerid,
		valueType:   valueType,
		measurement: measurement,
		fields:      fields,
		tags:        tags,
		t:           t}
}

func (ac *accumulator) AddFields(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	m := ac.MakeMetric(measurement, fields, tags, telegraf.Untyped, ac.getTime(t))
	ac.metrics <- m

}

func (ac *accumulator) AddGauge(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	m := ac.MakeMetric(measurement, fields, tags, telegraf.Gauge, ac.getTime(t))
	ac.metrics <- m
}

func (ac *accumulator) AddCounter(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	m := ac.MakeMetric(measurement, fields, tags, telegraf.Counter, ac.getTime(t))
	ac.metrics <- m
}

func (ac *accumulator) AddSummary(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	m := ac.MakeMetric(measurement, fields, tags, telegraf.Summary, ac.getTime(t))
	ac.metrics <- m
}

func (ac *accumulator) AddHistogram(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	m := ac.MakeMetric(measurement, fields, tags, telegraf.Histogram, ac.getTime(t))
	ac.metrics <- m
}

func (ac *accumulator) AddError(err error) {
	if err != nil {
		ac.errors <- err
	}
}

func (ac *accumulator) SetPrecision(precision, interval time.Duration) {
	if precision > 0 {
		ac.precision = precision
		return
	}
	switch {
	case interval >= time.Second:
		ac.precision = time.Second
	case interval >= time.Millisecond:
		ac.precision = time.Millisecond
	case interval >= time.Microsecond:
		ac.precision = time.Microsecond
	default:
		ac.precision = time.Nanosecond
	}
}

func (ac accumulator) getTime(t []time.Time) time.Time {
	var timestamp time.Time
	if len(t) > 0 {
		timestamp = t[0]
	} else {
		timestamp = time.Now()
	}
	return timestamp.Round(ac.precision)
}
