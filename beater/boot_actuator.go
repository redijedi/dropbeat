package beater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const metricsStatsPath = "/metrics"
const healthStatsPath = "/health"

// HealthStats contains the json health data
type HealthStats struct {
	Status    string `json:"status"`
	DiskSpace struct {
		Status    string `json:"status"`
		Total     uint64 `json:"total"`
		Free      uint64 `json:"free"`
		Threshold uint64 `json:"threshold"`
	} `json:"diskSpace"`
	DB struct {
		Status   string `json:"status"`
		Database string `json:"database"`
		Hello    uint64 `json:"hello"`
	} `json:"db"`
}

// MetricsStats contains the json metrics data
type MetricsStats struct {
	Mem struct {
		Total uint64 `json:"total"`
		Free  uint64 `json:"free"`
	} `json:"mem"`
	Processors  uint64  `json:"processors"`
	LoadAverage float64 `json:"load_average"`
	Uptime      struct {
		Total    uint64 `json:"total"`
		Instance uint64 `json:"instance"`
	} `json:"uptime"`
	Heap struct {
		Total     uint64 `json:"total"`
		Committed uint64 `json:"committed"`
		Init      uint64 `json:"init"`
		Used      uint64 `json:"used"`
	} `json:"heap"`
	NonHeap struct {
		Total     uint64 `json:"total"`
		Committed uint64 `json:"committed"`
		Init      uint64 `json:"init"`
		Used      uint64 `json:"used"`
	} `json:"non_heap"`
	Threads struct {
		Total        uint64 `json:"total"`
		TotalStarted uint64 `json:"started"`
		Peak         uint64 `json:"peak"`
		Daemon       uint64 `json:"daemon"`
	} `json:"non_heap"`
	Classes struct {
		Total    uint64 `json:"total"`
		Loaded   uint64 `json:"loaded"`
		Unloaded uint64 `json:"unloaded"`
	} `json:"classes"`
	GC struct {
		Scavenge struct {
			Count uint64 `json:"count"`
			Time  uint64 `json:"time"`
		} `json:"scavenge"`
		Marksweep struct {
			Count uint64 `json:"count"`
			Time  uint64 `json:"time"`
		} `json:"marksweep"`
	} `json:"gc"`
	HTTP struct {
		SessionsMax    int64  `json:"max_sessions"`
		SessionsActive uint64 `json:"active_sessions"`
	} `json:"http"`
	DataSource struct {
		PrimaryActive uint64  `json:"primary_active"`
		PrimaryUsage  float64 `json:"primary_usage"`
	} `json:"data_source"`
	GaugeResponse struct {
		Actuator    float64 `json:"actuator,omitempty"`
		Autoconfig  float64 `json:"autoconfig,omitempty"`
		Beans       float64 `json:"beans,omitempty"`
		Configprops float64 `json:"configprops,omitempty"`
		Dump        float64 `json:"dump,omitempty"`
		Env         float64 `json:"env,omitempty"`
		Health      float64 `json:"health,omitempty"`
		Info        float64 `json:"info,omitempty"`
		Root        float64 `json:"root,omitempty"`
		Trace       float64 `json:"trace,omitempty"`
		Unmapped    float64 `json:"unmapped,omitempty"`
	} `json:"gauge_response"`
	Status struct {
		TWO00 struct {
			Actuator    uint64 `json:"actuator,omitempty"`
			Autoconfig  uint64 `json:"autoconfig,omitempty"`
			Beans       uint64 `json:"beans,omitempty"`
			Configprops uint64 `json:"configprops,omitempty"`
			Dump        uint64 `json:"dump,omitempty"`
			Env         uint64 `json:"env,omitempty"`
			Health      uint64 `json:"health,omitempty"`
			Info        uint64 `json:"info,omitempty"`
			Root        uint64 `json:"root,omitempty"`
			Trace       uint64 `json:"trace,omitempty"`
		} `json:"200"`
	} `json:"status"`
}

// RawMetricsStats contains the raw metics data
type RawMetricsStats struct {
	Mem                         uint64  `json:"mem"`
	MemFree                     uint64  `json:"mem.free"`
	Processors                  uint64  `json:"processors"`
	InstanceUptime              uint64  `json:"instance.uptime"`
	Uptime                      uint64  `json:"uptime"`
	SystemloadAverage           float64 `json:"systemload.average"`
	HeapCommitted               uint64  `json:"heap.committed"`
	HeapInit                    uint64  `json:"heap.init"`
	HeapUsed                    uint64  `json:"heap.used"`
	Heap                        uint64  `json:"heap"`
	NonheapCommitted            uint64  `json:"nonheap.committed"`
	NonheapInit                 uint64  `json:"nonheap.init"`
	NonheapUsed                 uint64  `json:"nonheap.used"`
	Nonheap                     uint64  `json:"nonheap"`
	ThreadsPeak                 uint64  `json:"threads.peak"`
	ThreadsDaemon               uint64  `json:"threads.daemon"`
	ThreadsTotalStarted         uint64  `json:"threads.totalStarted"`
	Threads                     uint64  `json:"threads"`
	Classes                     uint64  `json:"classes"`
	ClassesLoaded               uint64  `json:"classes.loaded"`
	ClassesUnloaded             uint64  `json:"classes.unloaded"`
	GCPsScavengeCount           uint64  `json:"gc.ps_scavenge.count"`
	GCPsScavengeTime            uint64  `json:"gc.ps_scavenge.time"`
	GCPsMarksweepCount          uint64  `json:"gc.ps_marksweep.count"`
	GCPsMarksweepTime           uint64  `json:"gc.ps_marksweep.time"`
	HTTPSessionsMax             int64   `json:"httpsessions.max"`
	HTTPSessionsActive          uint64  `json:"httpsessions.active"`
	DateSourcePrimaryActive     uint64  `json:"datasource.primary.active"`
	DateSourcePrimaryUsage      float64 `json:"datasource.primary.usage"`
	GaugeResponseActuator       float64 `json:"gauge.response.actuator"`
	GaugeResponseBeans          float64 `json:"gauge.response.beans"`
	GaugeResponseTrace          float64 `json:"gauge.response.trace"`
	GaugeResponseAutoconfig     float64 `json:"gauge.response.autoconfig"`
	GaugeResponseDump           float64 `json:"gauge.response.dump"`
	GaugeResponseHealth         float64 `json:"gauge.response.health"`
	GaugeResponseRoot           float64 `json:"gauge.response.root"`
	GaugeResponseUnmapped       float64 `json:"gauge.response.unmapped"`
	GaugeResponseInfo           float64 `json:"gauge.response.info"`
	GaugeResponseEnv            float64 `json:"gauge.response.env"`
	GaugeResponseConfigprops    float64 `json:"gauge.response.configprops"`
	CounterStatus200Actuator    uint64  `json:"counter.status.200.actuator"`
	CounterStatus200Autoconfig  uint64  `json:"counter.status.200.autoconfig"`
	CounterStatus200Beans       uint64  `json:"counter.status.200.beans"`
	CounterStatus200Configprops uint64  `json:"counter.status.200.configprops"`
	CounterStatus200Dump        uint64  `json:"counter.status.200.dump"`
	CounterStatus200Env         uint64  `json:"counter.status.200.env"`
	CounterStatus200Health      uint64  `json:"counter.status.200.health"`
	CounterStatus200Info        uint64  `json:"counter.status.200.info"`
	CounterStatus200Root        uint64  `json:"counter.status.200.root"`
	CounterStatus200Trace       uint64  `json:"counter.status.200.trace"`
}

// GetHealthStats returns the health statistics
func (bt *Dropbeat) GetHealthStats(u url.URL) (*HealthStats, error) {
	res, err := http.Get(strings.TrimSuffix(u.String(), "/") + healthStatsPath)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	stats := &HealthStats{}
	err = json.Unmarshal([]byte(body), &stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// GetMetricsStats returns the metrics statistics
func (bt *Dropbeat) GetMetricsStats(u url.URL) (*MetricsStats, error) {
	res, err := http.Get(strings.TrimSuffix(u.String(), "/") + metricsStatsPath)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rawStats := &RawMetricsStats{}
	err = json.Unmarshal([]byte(body), &rawStats)
	if err != nil {
		return nil, err
	}

	// Transform into usable JSON format
	stats := &MetricsStats{}
	stats.Mem.Free = rawStats.MemFree
	stats.Mem.Total = rawStats.Mem
	stats.Processors = rawStats.Processors
	stats.LoadAverage = rawStats.SystemloadAverage
	stats.Uptime.Total = rawStats.Uptime
	stats.Uptime.Instance = rawStats.InstanceUptime
	stats.Heap.Total = rawStats.Heap
	stats.Heap.Init = rawStats.HeapInit
	stats.Heap.Committed = rawStats.HeapCommitted
	stats.Heap.Used = rawStats.HeapUsed
	stats.NonHeap.Total = rawStats.Nonheap
	stats.NonHeap.Init = rawStats.NonheapInit
	stats.NonHeap.Committed = rawStats.NonheapCommitted
	stats.NonHeap.Used = rawStats.NonheapUsed
	stats.Threads.Total = rawStats.Threads
	stats.Threads.TotalStarted = rawStats.ThreadsTotalStarted
	stats.Threads.Peak = rawStats.ThreadsPeak
	stats.Threads.Daemon = rawStats.ThreadsDaemon
	stats.Classes.Total = rawStats.Classes
	stats.Classes.Loaded = rawStats.ClassesLoaded
	stats.Classes.Unloaded = rawStats.ClassesUnloaded
	stats.GC.Scavenge.Count = rawStats.GCPsScavengeCount
	stats.GC.Scavenge.Time = rawStats.GCPsScavengeTime
	stats.GC.Marksweep.Count = rawStats.GCPsMarksweepCount
	stats.GC.Marksweep.Time = rawStats.GCPsMarksweepTime
	stats.HTTP.SessionsActive = rawStats.HTTPSessionsActive
	stats.HTTP.SessionsMax = rawStats.HTTPSessionsMax
	stats.DataSource.PrimaryActive = rawStats.DateSourcePrimaryActive
	stats.DataSource.PrimaryUsage = rawStats.DateSourcePrimaryUsage
	stats.GaugeResponse.Actuator = rawStats.GaugeResponseActuator
	stats.GaugeResponse.Autoconfig = rawStats.GaugeResponseAutoconfig
	stats.GaugeResponse.Beans = rawStats.GaugeResponseBeans
	stats.GaugeResponse.Configprops = rawStats.GaugeResponseConfigprops
	stats.GaugeResponse.Dump = rawStats.GaugeResponseDump
	stats.GaugeResponse.Env = rawStats.GaugeResponseEnv
	stats.GaugeResponse.Health = rawStats.GaugeResponseHealth
	stats.GaugeResponse.Info = rawStats.GaugeResponseInfo
	stats.GaugeResponse.Root = rawStats.GaugeResponseRoot
	stats.GaugeResponse.Trace = rawStats.GaugeResponseTrace
	stats.GaugeResponse.Unmapped = rawStats.GaugeResponseUnmapped
	stats.Status.TWO00.Actuator = rawStats.CounterStatus200Actuator
	stats.Status.TWO00.Autoconfig = rawStats.CounterStatus200Autoconfig
	stats.Status.TWO00.Beans = rawStats.CounterStatus200Beans
	stats.Status.TWO00.Configprops = rawStats.CounterStatus200Configprops
	stats.Status.TWO00.Dump = rawStats.CounterStatus200Dump
	stats.Status.TWO00.Env = rawStats.CounterStatus200Env
	stats.Status.TWO00.Health = rawStats.CounterStatus200Health
	stats.Status.TWO00.Info = rawStats.CounterStatus200Info
	stats.Status.TWO00.Root = rawStats.CounterStatus200Root
	stats.Status.TWO00.Trace = rawStats.CounterStatus200Trace

	return stats, nil
}
