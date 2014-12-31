package providers

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	types "github.com/mchmarny/thingz-commons"
)

const (
	// TODO: refactor for each provider to describe itself
	STRATEGY_CPU  = "cpu"
	STRATEGY_CPUS = "cpus"
	STRATEGY_MEM  = "mem"
	STRATEGY_SWAP = "swap"
	STRATEGY_LOAD = "load"
	STRATEGY_PROC = "proc"

	FORMAT_ERROR = "Invalid strategy format: "
)

// Provider describes the metric provider functionality
type Provider interface {

	// SetFrequency of execution
	SetFrequency(time.Duration)

	// Provide metric group
	Provide(out chan<- *types.MetricCollection) error
}

func GetProviders(plan string) (map[string]Provider, error) {

	if len(plan) < 1 {
		return nil, errors.New("Plan required")
	}

	provs := make(map[string]Provider)

	for _, p := range strings.Split(plan, ",") {

		// get strategy
		strategy := strings.Split(strings.Trim(p, " "), ":")
		if len(strategy) != 2 {
			log.Fatal(FORMAT_ERROR)
			return nil, errors.New(FORMAT_ERROR + p)
		}

		// frequency of execution for this strategy
		n, err := strconv.Atoi(strategy[1])
		if err != nil {
			log.Fatal(err)
			return nil, errors.New(FORMAT_ERROR + p)
		}

		freq := time.Duration(int32(n)) * time.Second
		group := strings.ToLower(strings.Trim(strategy[0], " "))

		// TODO: spool these into a map first
		switch group {
		case STRATEGY_CPU:
			provs[group] = CPUProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_CPUS:
			provs[group] = CPUSProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_MEM:
			provs[group] = MemoryProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_SWAP:
			provs[group] = SwapProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_LOAD:
			provs[group] = LoadProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_PROC:
			provs[group] = ProcProvider{
				Group:     group,
				Frequency: freq,
			}
		default:
			log.Fatal(FORMAT_ERROR)
			return nil, errors.New(FORMAT_ERROR + p)
		}

	}

	return provs, nil

}
