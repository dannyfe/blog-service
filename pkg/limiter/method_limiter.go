package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"strings"
)

type MethodLimiter struct {
	*Limite
}

func NewMethodLimiter() LimiteIface {
	return MethodLimiter{
		Limite: &Limite{limiteBuckets: make(map[string]*ratelimit.Bucket)},
	}
}

func (m MethodLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	index := strings.Index(uri, "?")
	if index == -1 {
		return uri
	}

	return uri[:index]
}

func (m MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := m.limiteBuckets[key]
	return bucket, ok
}

func (m MethodLimiter) AddBuckets(rules ...LimiteBucketRule) LimiteIface {
	for _, rule := range rules {
		if _, ok := m.limiteBuckets[rule.Key]; ok {
			m.limiteBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}

	return m
}





