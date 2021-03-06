package main

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
)

func defaultConfig() *config {
	return &config{
		Root:  "/var/lib/containerd",
		State: "/run/containerd",
		GRPC: grpcConfig{
			Socket: "/run/containerd/containerd.sock",
		},
		Debug: debug{
			Level:  "info",
			Socket: "/run/containerd/debug.sock",
		},
		Snapshotter: "overlay",
	}
}

// loadConfig loads the config from the provided path
func loadConfig(path string) error {
	md, err := toml.DecodeFile(path, conf)
	if err != nil {
		return err
	}
	conf.md = md
	return nil
}

// config specifies the containerd configuration file in the TOML format.
// It contains fields to configure various subsystems and containerd as a whole.
type config struct {
	// State is the path to a directory where containerd will store runtime state
	State string `toml:"state"`
	// Root is the path to a directory where containerd will store persistent data
	Root string `toml:"root"`
	// GRPC configuration settings
	GRPC grpcConfig `toml:"grpc"`
	// Debug and profiling settings
	Debug debug `toml:"debug"`
	// Metrics and monitoring settings
	Metrics metricsConfig `toml:"metrics"`
	// Snapshotter specifies which snapshot driver to use
	Snapshotter string `toml:"snapshotter"`
	// Plugins provides plugin specific configuration for the initialization of a plugin
	Plugins map[string]toml.Primitive `toml:"plugins"`
	// Enable containerd as a subreaper
	Subreaper bool `toml:"subreaper"`

	md toml.MetaData
}

func (c *config) decodePlugin(name string, v interface{}) error {
	p, ok := c.Plugins[name]
	if !ok {
		return nil
	}
	return c.md.PrimitiveDecode(p, v)
}

func (c *config) WriteTo(w io.Writer) (int64, error) {
	buf := bytes.NewBuffer(nil)
	e := toml.NewEncoder(buf)
	if err := e.Encode(c); err != nil {
		return 0, err
	}
	return io.Copy(w, buf)
}

type grpcConfig struct {
	Socket string `toml:"socket"`
	Uid    int    `toml:"uid"`
	Gid    int    `toml:"gid"`
}

type debug struct {
	Socket string `toml:"socket"`
	Level  string `toml:"level"`
}

type metricsConfig struct {
	Address string `toml:"address"`
}
