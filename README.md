# APT Exporter

A Prometheus exporter for APT package management information on Debian/Ubuntu systems.

## Overview

APT Exporter collects and exposes metrics related to APT package management, including:

- Number of available package updates
- Number of available security updates
- Time since the last APT update
- Whether a system reboot is required

These metrics are exposed in Prometheus format, allowing you to monitor your Debian/Ubuntu systems' update status and configure alerts for outdated packages or required reboots.

## Metrics

### Core Metrics

| Metric Name | Description | Type |
|-------------|-------------|------|
| `<prefix>_updates_available` | Number of available package updates | Gauge |
| `<prefix>_security_updates_available` | Number of available security updates | Gauge |
| `<prefix>_seconds_since_last_update` | Seconds since last successful apt update | Gauge |
| `<prefix>_reboot_required` | 1 if a reboot is required, 0 otherwise | Gauge |

### Collector Metrics

| Metric Name | Description | Type |
|-------------|-------------|------|
| `<prefix>_collector_success` | 1 if the last collection was successful, 0 otherwise | Gauge |
| `<prefix>_collector_duration_seconds` | Duration of the last collection in seconds | Gauge |
| `<prefix>_collector_last_timestamp` | Timestamp of the last collection | Gauge |

The prefix is configurable in the configuration file (default: `ubuntu`).

### Go Runtime Metrics

By default, the exporter does not expose Go runtime metrics (like memory usage, goroutines, GC stats, etc.). This keeps the metrics output clean and focused on APT-related information. If you need these metrics for debugging or monitoring the exporter itself, you can modify the code in `cmd/apt-exporter/main.go` to use the default Prometheus registry instead of a custom one.

## Requirements

- Linux (Debian/Ubuntu-based system with APT)
- Go 1.18 or higher (for building from source)
- `/usr/lib/update-notifier/apt-check` script (provided by the `update-notifier-common` package)

## Installation

For detailed installation instructions, please see [INSTALL.md](INSTALL.md).

### Quick Start

```bash
# Download the latest release for your architecture (example for amd64)
curl -L https://github.com/ncecere/apt-exporter/releases/download/v0.1.0/apt-exporter_0.1.0_linux_amd64.tar.gz -o apt-exporter.tar.gz

# Extract the archive
tar -xzf apt-exporter.tar.gz

# Make the binary executable and move it to a location in your PATH
chmod +x apt-exporter
sudo mv apt-exporter /usr/local/bin/

# Copy and adjust the configuration file
sudo mkdir -p /etc/apt-exporter
sudo cp config.yml /etc/apt-exporter/
```

See [INSTALL.md](INSTALL.md) for more detailed instructions, including:
- Building from source
- Setting up as a system service
- Troubleshooting common issues
- Updating and uninstalling

## Configuration

APT Exporter uses a YAML configuration file. By default, it looks for `config.yml` in the current directory, but you can specify a different path using the `-config` flag.

Example configuration:

```yaml
check_interval_seconds: 300
listen_address: ":9100"
apt_check_path: "/usr/lib/update-notifier/apt-check"
update_stamp_path: "/var/lib/apt/periodic/update-success-stamp"
reboot_required_file: "/var/run/reboot-required"
log_level: "info"
command_timeout_seconds: 10
metrics_endpoint: "/metrics"
metric_prefix: "ubuntu"
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `check_interval_seconds` | How often to check for updates (in seconds) | 300 |
| `listen_address` | IP:port where the HTTP server listens | ":9100" |
| `apt_check_path` | Path to the apt-check script | "/usr/lib/update-notifier/apt-check" |
| `update_stamp_path` | Path to the update success stamp file | "/var/lib/apt/periodic/update-success-stamp" |
| `reboot_required_file` | Path to the reboot-required file | "/var/run/reboot-required" |
| `log_level` | Logging level (debug, info, warn, error) | "info" |
| `command_timeout_seconds` | Timeout for external commands (in seconds) | 10 |
| `metrics_endpoint` | URL path for exposing metrics | "/metrics" |
| `metric_prefix` | Prefix added to all metric names | "ubuntu" |

## Usage

```bash
# Run with default configuration file (config.yml in current directory)
apt-exporter

# Run with a specific configuration file
apt-exporter -config /etc/apt-exporter/config.yml

# Show version information
apt-exporter -version

# Skip validation of file paths (useful for testing)
apt-exporter -skip-path-validation
```

## Running as a Service

### Systemd

Create a systemd service file at `/etc/systemd/system/apt-exporter.service`:

```ini
[Unit]
Description=APT Exporter
After=network.target

[Service]
# Run as root to ensure access to apt-check and other system files
User=root
Group=root
ExecStart=/usr/local/bin/apt-exporter -config /etc/apt-exporter/config.yml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable apt-exporter
sudo systemctl start apt-exporter
```

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'apt'
    static_configs:
      - targets: ['localhost:9100']
```

## Grafana Dashboard

A sample Grafana dashboard is available in the `dashboards` directory. Import this dashboard into your Grafana instance to visualize the APT metrics.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development

The project follows standard Go project layout:

- `cmd/apt-exporter/`: Main application entry point
- `internal/`: Internal packages
  - `config/`: Configuration handling
  - `collector/`: Metrics collection logic
  - `metrics/`: Prometheus metrics definitions

### Testing

Run the tests with:

```bash
go test ./...
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed list of changes in each release.

## Releasing

To create a new release:

1. Update the [CHANGELOG.md](CHANGELOG.md) file with the changes in the new version.

2. Create a new tag with the version number (following [Semantic Versioning](https://semver.org/)):

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
```

3. Push the tag to GitHub:

```bash
git push origin v1.0.0
```

4. GitHub Actions will automatically:
   - Build the binaries for different Linux architectures
   - Create a GitHub release
   - Upload the binaries to the release

## License

This project is licensed under the MIT License - see the LICENSE file for details.
