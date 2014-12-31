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

// ProviderConfig is the provider for CPU information
type ProviderConfig struct {
	Source    string
	Dimension string
	Frequency time.Duration
}

// Provider describes the metric provider functionality
type Provider interface {

	// Provide metric group
	Provide(out chan<- *types.MetricCollection) error
}

func getProviderConfig(src, dim string, freq time.Duration) *ProviderConfig {
	return &ProviderConfig{
		Source:    src,
		Dimension: dim,
		Frequency: freq,
	}
}

func GetProviders(src, plan string) (map[string]Provider, error) {

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
		dim := strings.ToLower(strings.Trim(strategy[0], " "))
		conf := getProviderConfig(src, dim, freq)

		// TODO: spool these into a map first
		switch dim {
		case STRATEGY_CPU:
			provs[dim] = CPUProvider{conf}
		case STRATEGY_CPUS:
			provs[dim] = CPUSProvider{conf}
		case STRATEGY_MEM:
			provs[dim] = MemoryProvider{conf}
		case STRATEGY_SWAP:
			provs[dim] = SwapProvider{conf}
		case STRATEGY_LOAD:
			provs[dim] = LoadProvider{conf}
		case STRATEGY_PROC:
			provs[dim] = ProcProvider{conf}
		default:
			log.Fatal(FORMAT_ERROR)
			return nil, errors.New(FORMAT_ERROR + p)
		}

	}

	return provs, nil

}
