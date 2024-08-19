package metrics

import (
	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func (f *fetcher) ZoneRequest(zoneRequestTotal *prometheus.CounterVec) {
	for _, zones := range sliceChunk(toArr(f.client.Zones, convZoneID, filterNonFreeZone), ZONES_LIMIT) {
		go func() {
			defer f.wg.Done()
			f.wg.Add(1)

			data, err := cloudflare.ZonesRequestsTotal(
				f.context, f.client.GQL,
				zones, f.last, f.now,
			)
			if err != nil {
				f.logger.Warn("Cannot fetching zone requests", zap.Error(err))
			}

			for _, z := range data.Viewer.Zones {
				zone := f.client.Zones[z.ZoneTag]
				if len(z.Requests) == 1 {
					requests := z.Requests[0]
					zoneRequestTotal.WithLabelValues(zone.Account.Name, zone.Name).Add(float64(requests.Sum.Visits))
				} else {
					f.logger.Warn("Zone doesn't have requests metric",
						zap.String("zone", zone.Name),
						zap.Int("metric-size", len(z.Requests)),
					)
				}
			}
		}()
	}
}
