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
	Loglevel string `json:"loglevel,omitempty"`
}

// InboundConfig defines an inbound proxy connection handler.
type InboundConfig struct {
	// Tag is a unique identifier for this inbound.
	Tag string `json:"tag,omitempty"`
	// Port is the port number to listen on.
	Port int `json:"port"`
	// Listen is the IP address to listen on.
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
type RoutingConfig struct {
	// DomainStrategy controls how domain names are resolved for routing.
	DomainStrategy string `json:"domainStrategy,omitempty"`
	// Rules is a list of routing rules.
	Rules []json.RawMessage `json:"rules,omitempty"`
}

// SniffingConfig controls traffic sniffing behavior on an inbound.
type SniffingConfig struct {
	// Enabled turns sniffing on or off.
	Enabled bool `json:"enabled"`
	// DestOverride lists protocols whose destination should be overridden.
	DestOverride []string `json:"destOverride,omitempty"`
}

// DetectConfigFormat infers the configuration format from the file extension.
func DetectConfigFormat(path string) (ConfigFormat, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".jsonc":
		return ConfigFormatJSON, nil
	case ".toml":
		return ConfigFormatTOML, nil
	case ".yaml", ".yml":
		return ConfigFormatYAML, nil
	default:
		return ConfigFormatJSON, fmt.Errorf("unknown config file extension: %s, defaulting to JSON", ext)
	}
}

// LoadConfigFromFile reads and parses a configuration file into a Config struct.
// Currently supports JSON format; TOML and YAML support may be added in future.
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	format, err := DetectConfigFormat(path)
	if err != nil {
		// Log warning but continue with JSON default
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	switch format {
	case ConfigFormatJSON:
		return parseJSONConfig(data)
	default:
		return nil, fmt.Errorf("config format not yet supported, please use JSON")
	}
}

// parseJSONConfig unmarshals raw JSON bytes into a Config struct.
func parseJSONConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %w", err)
	}
	return &cfg, nil
}
