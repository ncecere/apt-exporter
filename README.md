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

| Metric Name | Description | Type |
|-------------|-------------|------|
| `<prefix>_updates_available` | Number of available package updates | Gauge |
| `<prefix>_security_updates_available` | Number of available security updates | Gauge |
| `<prefix>_seconds_since_last_update` | Seconds since last successful apt update | Gauge |
| `<prefix>_reboot_required` | 1 if a reboot is required, 0 otherwise | Gauge |

The prefix is configurable in the configuration file (default: `ubuntu`).

## Requirements

- Go 1.18 or higher
- Debian/Ubuntu-based system with APT
- `/usr/lib/update-notifier/apt-check` script (provided by the `update-notifier-common` package)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ncecere/apt-exporter.git
cd apt-exporter

# Build the binary
go build -o apt-exporter cmd/apt-exporter/main.go

# Copy the binary to a location in your PATH
sudo cp apt-exporter /usr/local/bin/

# Copy and adjust the configuration file
sudo mkdir -p /etc/apt-exporter
sudo cp config.yml /etc/apt-exporter/
```

### Using Go Install

```bash
go install github.com/ncecere/apt-exporter/cmd/apt-exporter@latest
```

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
```

## Running as a Service

### Systemd

Create a systemd service file at `/etc/systemd/system/apt-exporter.service`:

```ini
[Unit]
Description=APT Exporter
After=network.target

[Service]
User=nobody
Group=nogroup
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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
