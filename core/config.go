package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigFormat represents the format of a configuration file.
type ConfigFormat int

const (
	// ConfigFormatJSON represents JSON configuration format.
	ConfigFormatJSON ConfigFormat = iota
	// ConfigFormatTOML represents TOML configuration format.
	ConfigFormatTOML
	// ConfigFormatYAML represents YAML configuration format.
	ConfigFormatYAML
)

// Config holds the top-level v2ray configuration.
type Config struct {
	// Log defines logging behavior.
	Log *LogConfig `json:"log,omitempty"`
	// Inbounds defines inbound proxy configurations.
	Inbounds []InboundConfig `json:"inbounds,omitempty"`
	// Outbounds defines outbound proxy configurations.
	Outbounds []OutboundConfig `json:"outbounds,omitempty"`
	// DNS defines DNS server configuration.
	DNS *DNSConfig `json:"dns,omitempty"`
	// Routing defines traffic routing rules.
	Routing *RoutingConfig `json:"routing,omitempty"`
}

// LogConfig defines the logging configuration.
type LogConfig struct {
	// Access is the path to the access log file.
	Access string `json:"access,omitempty"`
	// Error is the path to the error log file.
	Error string `json:"error,omitempty"`
	// Loglevel sets the logging verbosity: debug, info, warning, error, none.
	// Personal note: I prefer "warning" as default to reduce noise in logs.
	Loglevel string `json:"loglevel,omitempty"`
}

// InboundConfig defines an inbound proxy connection handler.
type InboundConfig struct {
	// Tag is a unique identifier for this inbound.
	Tag string `json:"tag,omitempty"`
	// Port is the port number to listen on.
	Port int `json:"port"`
	// Listen is the IP address to listen on. Defaults to 127.0.0.1 for security.
	// Personal note: keeping this as 127.0.0.1 only; never expose to 0.0.0.0 on untrusted networks.
	Listen string `json:"listen,omitempty"`
	// Protocol is the inbound proxy protocol (e.g. vmess, socks, http).
	Protocol string `json:"protocol"`
	// Settings contains protocol-specific settings as raw JSON.
	Settings json.RawMessage `json:"settings,omitempty"`
	// StreamSettings defines network transport settings.
	StreamSettings json.RawMessage `json:"streamSettings,omitempty"`
	// Sniffing enables traffic sniffing for routing purposes.
	Sniffing *SniffingConfig `json:"sniffing,omitempty"`
}

// OutboundConfig defines an outbound proxy connection handler.
type OutboundConfig struct {
	// Tag is a unique identifier for this outbound.
	Tag string `json:"tag,omitempty"`
	// Protocol is the outbound proxy protocol (e.g. vmess, freedom, blackhole).
	Protocol string `json:"protocol"`
	// Settings contains protocol-specific settings as raw JSON.
	Settings json.RawMessage `json:"settings,omitempty"`
	// StreamSettings defines network transport settings.
	StreamSettings json.RawMessage `json:"streamSettings,omitempty"`
}

// DNSConfig defines DNS server settings.
type DNSConfig struct {
	// Servers is a list of DNS server addresses.
	Servers []interface{} `json:"servers,omitempty"`
}

// RoutingConfig defines traffic routing rules.
type RoutingConfig struc