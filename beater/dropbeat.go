package beater

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/redijedi/dropbeat/config"
)

const selector = "dropbeat"

// Dropbeat is the container for the data that will be pushed out
type Dropbeat struct {
	period time.Duration
	urls   []*url.URL

	beatConfig *config.Config

	done   chan struct{}
	client publisher.Client

	metricsStats bool
	healthStats  bool
}

// New Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Dropbeat{
		done:       make(chan struct{}),
		beatConfig: &config,
	}
	return bt, nil
}

/// *** Beater interface methods ***///

// Config configures this beat
func (bt *Dropbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := b.RawConfig.Unpack(&bt.beatConfig)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	//define default URL if none provided
	var urlConfig []string
	if bt.beatConfig.Dropbeat.URLs != nil {
		urlConfig = bt.beatConfig.Dropbeat.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1"}
	}

	bt.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid Dropwizard Metrics URL: %v", err)
			return err
		}
		bt.urls[i] = u
	}

	if bt.beatConfig.Dropbeat.Stats.Metrics != nil {
		bt.metricsStats = *bt.beatConfig.Dropbeat.Stats.Metrics
	} else {
		bt.metricsStats = true
	}

	if bt.beatConfig.Dropbeat.Stats.Health != nil {
		bt.healthStats = *bt.beatConfig.Dropbeat.Stats.Health
	} else {
		bt.healthStats = true
	}

	if !bt.metricsStats && !bt.healthStats {
		return errors.New("Invalid statistics configuration")
	}

	return nil
}

// Setup sets up the beat
func (bt *Dropbeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Dropbeat.Period == "" {
		bt.beatConfig.Dropbeat.Period = "10s"
	}

	bt.client = b.Publisher.Connect()

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Dropbeat.Period)
	if err != nil {
		return err
	}

	logp.Debug(selector, "Init dropbeat")
	logp.Debug(selector, "Period %v\n", bt.period)
	logp.Debug(selector, "Watch %v", bt.urls)
	logp.Debug(selector, "Metrics statistics %t\n", bt.metricsStats)
	logp.Debug(selector, "Health statistics %t\n", bt.healthStats)

	return nil
}

// Run runs the beat
func (bt *Dropbeat) Run(b *beat.Beat) error {
	logp.Info("dropbeat is running! Hit CTRL-C to stop it.")

	for _, u := range bt.urls {
		go func(u *url.URL) {

			ticker := time.NewTicker(bt.period)
			counter := 1
			for {
				select {
				case <-bt.done:
					goto GotoFinish
				case <-ticker.C:
				}

				timerStart := time.Now()

				if bt.metricsStats {
					logp.Debug(selector, "Metrics stats for url: %v", u)
					metricsStats, err := bt.GetMetricsStats(*u)

					if err != nil {
						logp.Err("Error reading Metrics stats: %v", err)
					} else {
						logp.Debug(selector, "Metrics stats detail: %+v", metricsStats)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "metrics",
							"counter":    counter,
							"metrics":    metricsStats,
						}

						bt.client.PublishEvent(event)
						logp.Info("Dropwizard /metrics stats sent")
						counter++
					}
				}

				if bt.healthStats {
					logp.Debug(selector, "Health stats for url: %v", u)
					healthStats, err := bt.GetHealthStats(*u)

					if err != nil {
						logp.Err("Error reading Health stats: %v", err)
					} else {
						logp.Debug(selector, "Health stats detail: %+v", healthStats)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "health",
							"counter":    counter,
							"health":     healthStats,
						}

						bt.client.PublishEvent(event)
						logp.Info("Dropwizard /health stats sent")
						counter++
					}
				}

				timerEnd := time.Now()
				duration := timerEnd.Sub(timerStart)
				if duration.Nanoseconds() > bt.period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}(u)
	}

	<-bt.done
	return nil
}

// Cleanup cleans up the beat
func (bt *Dropbeat) Cleanup(b *beat.Beat) error {
	return nil
}

// Stop stops the beat
func (bt *Dropbeat) Stop() {
	logp.Debug(selector, "Stop dropbeat")
	close(bt.done)
}
