package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type metricsRequest struct {
	Hostname string `uri:"hostname" binding:"required"`
}

type Metricsbridge struct {
	redis     *redis.Client
	scrapeURL string
	hostname  string
}

func NewMetricsbridge(r *redis.Client, scrapeURL string, hostname string) *Metricsbridge {
	return &Metricsbridge{r, scrapeURL, hostname}
}

func (mb *Metricsbridge) GetMetricsHandler(c *gin.Context) {
	r := metricsRequest{}
	if err := c.ShouldBindUri(&r); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	res := mb.redis.Get(c, makeKey(r.Hostname))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			c.String(http.StatusNotFound, err.Error())
			return
		} else {
			log.Error().Err(err).Msgf("[METRICS] НЕ МОГУ ПОЛУЧИТЬ ДАННЫЕ ИЗ REDIS")
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	if data, err := res.Bytes(); err != nil {
		log.Error().Err(err).Msgf("[METRICS] НЕ МОГУ ПРОЧИТАТЬ ДАННЫЕ ИЗ REDIS")
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Data(http.StatusOK, "text/plain", data)
	}
}

// RunMetricsTransfer - локальные метрики текущего сервера перекладываются в удалённый редис
func (mb *Metricsbridge) RunMetricsTransfer(ctx context.Context) error {
	log.Info().Msgf("[METRICS] ЗАПУЩЕНА ПЕРЕДАЧА МЕТРИК КАК ХОСТА %s", mb.hostname)
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			t := time.Now()
			if err := mb.scrape(ctx, mb.hostname); err != nil {
				log.Error().Err(err).Msgf("[METRICS] НЕ МОГУ СТАЩИТЬ МЕТРИКИ")
			} else {
				log.Info().Err(err).Msgf("[METRICS] SCRAPE УДАЧНЫЙ. t:%s", time.Now().Sub(t))
			}
		}
	}
}

func (mb *Metricsbridge) scrape(ctx context.Context, hostname string) error {
	resp, err := http.Get(mb.scrapeURL)
	if err != nil {
		return fmt.Errorf("[METRICS] неудачный ответ от %s: %w", mb.scrapeURL, err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[METRICS] данные не читаются от %s: %w", mb.scrapeURL, err)
	}

	res := mb.redis.Set(ctx, makeKey(hostname), data, time.Minute)
	if err := res.Err(); err != nil {
		return fmt.Errorf("[METRICS] ОШИБКА ЗАПИСИ В REDIS: %w", err)
	}

	return nil
}

func makeKey(hostname string) string {
	return fmt.Sprintf("host:%s", hostname)
}
