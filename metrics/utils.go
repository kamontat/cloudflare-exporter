package metrics

import "github.com/cloudflare/cloudflare-go"

func convZoneID(z cloudflare.Zone) string {
	return z.ID
}

func filterFreeZone(z cloudflare.Zone) bool {
	return z.Plan.ID == FREE_PLAN_ID
}

func filterNonFreeZone(z cloudflare.Zone) bool {
	return !filterFreeZone(z)
}

func filterEnterpriseZone(z cloudflare.Zone) bool {
	return z.Plan.ID == ENTERPRISE_ID
}

func filterNonEnterpriseZone(z cloudflare.Zone) bool {
	return !filterEnterpriseZone(z)
}

func toArr[T any, P any](m map[string]T, convert func(T) P, filter func(T) bool) (output []P) {
	output = make([]P, 0)
	for _, value := range m {
		if filter(value) {
			output = append(output, convert(value))
		}
	}
	return
}

func sliceChunk[T any](arr []T, size int) (output [][]T) {
	output = make([][]T, 0)
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}

		output = append(output, arr[i:end])
	}
	return
}
