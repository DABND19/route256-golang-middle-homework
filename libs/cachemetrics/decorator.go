package cachemetrics

import "time"

type Cache[KeyT comparable, ValueT any] interface {
	Get(key KeyT) (ValueT, bool)
	Set(key KeyT, value ValueT)
}

type CacheMetrics[KeyT comparable, ValueT any] struct {
	cacheTargetName string
	Cache[KeyT, ValueT]
}

func New[KeyT comparable, ValueT any](
	cache Cache[KeyT, ValueT],
	cacheTargetName string,
) *CacheMetrics[KeyT, ValueT] {
	return &CacheMetrics[KeyT, ValueT]{
		cacheTargetName,
		cache,
	}
}

func (d *CacheMetrics[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	timeBegin := time.Now()
	res, ok := d.Cache.Get(key)
	duration := time.Since(timeBegin)
	HistogramReadTime.WithLabelValues(d.cacheTargetName).Observe(duration.Seconds())
	if ok {
		HitsCounter.WithLabelValues(d.cacheTargetName).Inc()
	} else {
		MissesCounter.WithLabelValues(d.cacheTargetName).Inc()
	}
	return res, ok
}

func (d *CacheMetrics[KeyT, ValueT]) Set(key KeyT, value ValueT) {
	timeBegin := time.Now()
	d.Cache.Set(key, value)
	duration := time.Since(timeBegin)
	HistogramWriteTime.WithLabelValues(d.cacheTargetName).Observe(duration.Seconds())
}
