package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	syscall "syscall"

	"v2ray.com/core"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/common/platform"
	_ "v2ray.com/core/main/distro/all"
)

var (
	configFiles cmdarg.Arg // list of configuration files
	configDir   string
	version     = flag.Bool("version", false, "Show current version of V2Ray.")
	testConfig  = flag.Bool("test", false, "Test the configuration and exit without starting.")
	format      = flag.String("format", "json", "Format of input file, acceptable values are json, pb and toml.")
)

func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func dirExists(dir string) bool {
	if dir == "" {
		return false
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

func getConfigFilePath() cmdarg.Arg {
	if dirExists(configDir) {
		const defaultExt = ".json"
		entries, err := os.ReadDir(configDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read config dir: %v\n", err)
			return configFiles
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(entry.Name(), defaultExt) {
				configFiles = append(configFiles, filepath.Join(configDir, entry.Name()))
			}
		}
		return configFiles
	}

	if len(configFiles) > 0 {
		return configFiles
	}

	// try to find config in default locations
	defaultPaths := []string{
		platform.GetConfigurationPath(),
		filepath.Join(".", "config.json"),
	}
	for _, p := range defaultPaths {
		if fileExists(p) {
			return cmdarg.Arg{p}
		}
	}

	return configFiles
}

func startV2Ray() (core.Server, error) {
	configFiles := getConfigFilePath()

	config, err := core.LoadConfig(*format, configFiles[0], configFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	server, err := core.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return server, nil
}

func main() {
	// register -config and -c flags for specifying config files
	flag.Var(&configFiles, "config", "Config file(s) for V2Ray, multiple config files can be specified.")
	flag.Var(&configFiles, "c", "Short alias of -config.")
	flag.StringVar(&configDir, "confdir", "", "A directory with config files. All .json files in the directory will be loaded.")

	flag.Parse()

	printVersion()

	if *version {
		os.Exit(0)
	}

	server, err := startV2Ray()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// exit with a specific code to indicate config error
		os.Exit(23)
	}

	if *testConfig {
		fmt.Println("Configuration OK.")
		os.Exit(0)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start V2Ray:", err)
		os.Exit(1)
	}
	defer server.Close()

	// handle OS signals for graceful shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	<-osSignals
	fmt.Println("V2Ray shutting down...")
}

func printVersion() {
	version := core.VersionStatement()
	for _, s := range version {
		fmt.Println(s)
	}
	fmt.Printf("Go runtime: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
