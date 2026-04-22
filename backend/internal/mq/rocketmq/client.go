package rocketmq

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	appconfig "pilipili-go/backend/internal/config"
)

type Client struct {
	cfg       appconfig.RocketMQConfig
	endpoints []string
}

func New(cfg appconfig.RocketMQConfig) (*Client, error) {
	client := &Client{
		cfg:       cfg,
		endpoints: splitEndpoints(cfg.NameServerAddr),
	}

	if !cfg.Enabled {
		return client, nil
	}
	if len(client.endpoints) == 0 {
		return nil, fmt.Errorf("rocketmq nameserver endpoint is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DialTimeoutSecond)*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping rocketmq nameserver: %w", err)
	}
	return client, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.cfg.Enabled
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("rocketmq client is unavailable")
	}
	if !c.cfg.Enabled {
		return nil
	}
	if len(c.endpoints) == 0 {
		return fmt.Errorf("rocketmq nameserver endpoint is required")
	}

	dialer := &net.Dialer{
		Timeout: time.Duration(c.cfg.DialTimeoutSecond) * time.Second,
	}
	conn, err := dialer.DialContext(ctx, "tcp", c.endpoints[0])
	if err != nil {
		return err
	}
	return conn.Close()
}

func (c *Client) Topic(name string) string {
	if c == nil {
		return name
	}

	topic := strings.TrimSpace(name)
	prefix := strings.TrimSpace(c.cfg.TopicPrefix)
	if prefix == "" || topic == "" {
		return topic
	}
	return prefix + "." + topic
}

func (c *Client) ProducerGroup() string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.cfg.ProducerGroup)
}

func (c *Client) ConsumerGroup(name string) string {
	if c == nil {
		return strings.TrimSpace(name)
	}

	group := strings.TrimSpace(name)
	prefix := strings.TrimSpace(c.cfg.ConsumerGroupPrefix)
	if prefix == "" || group == "" {
		return group
	}
	return prefix + "." + group
}

func splitEndpoints(raw string) []string {
	parts := strings.Split(raw, ",")
	endpoints := make([]string, 0, len(parts))
	for _, part := range parts {
		endpoint := strings.TrimSpace(part)
		if endpoint == "" {
			continue
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}
