package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	messari "github.com/rajendrarusmana/coinspeed-worker/client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	gauges map[string]prometheus.Gauge
)

func main() {
	ExportMessariMetrics()
	http.Handle("/", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func ExportMessariMetrics() {
	client := messari.NewClient("")
	ctx := context.Background()
	viper.SetConfigName("application") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	gauges := map[string]prometheus.Gauge{}
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(func() {
		log.Info().Msg("Getting data from assets..")
		params := map[string]interface{}{
			"fields": "id,slug,name,symbol,metrics/market_data/price_usd,metrics/market_data/price_btc,metrics/market_data/ohlcv_last_1_hour,",
			"limit":  200,
		}
		assets, err := client.GetAllAssets(ctx, params)
		if err != nil {
			log.Error().Msg("Failed getting assets from messari")
		}

		for _, asset := range assets {

			if gauges[fmt.Sprintf("cryptos_price_btc_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_price_btc_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_price_btc"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_price_btc_%s", asset.Slug)].Set(asset.Metrics.MarketData.PriceBTC)

			if gauges[fmt.Sprintf("cryptos_price_usd_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_price_usd_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_price_usd"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_price_usd_%s", asset.Slug)].Set(asset.Metrics.MarketData.PriceUSD)

			if gauges[fmt.Sprintf("cryptos_open_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_open_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_open_1_hour"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_open_%s", asset.Slug)].Set(asset.Metrics.MarketData.OHLCVLast1Hour.Open)

			if gauges[fmt.Sprintf("cryptos_high_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_high_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_high_1_hour"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_high_%s", asset.Slug)].Set(asset.Metrics.MarketData.OHLCVLast1Hour.High)

			if gauges[fmt.Sprintf("cryptos_low_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_low_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_low_1_hour"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_low_%s", asset.Slug)].Set(asset.Metrics.MarketData.OHLCVLast1Hour.Low)

			if gauges[fmt.Sprintf("cryptos_close_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_close_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_close_1_hour"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_close_%s", asset.Slug)].Set(asset.Metrics.MarketData.OHLCVLast1Hour.Close)

			if gauges[fmt.Sprintf("cryptos_volume_%s", asset.Slug)] == nil {
				gauges[fmt.Sprintf("cryptos_volume_%s", asset.Slug)] = promauto.NewGauge(prometheus.GaugeOpts{
					Name:        fmt.Sprintf("cryptos_volume_1_hour"),
					ConstLabels: prometheus.Labels{"asset": asset.Slug, "symbol": asset.Symbol, "name": asset.Name},
				})
			}

			gauges[fmt.Sprintf("cryptos_volume_%s", asset.Slug)].Set(asset.Metrics.MarketData.OHLCVLast1Hour.Volume)
		}
	})

	s.StartAsync()
}
