# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2025-03-02

### Added
- Initial release of APT Exporter
- Core metrics:
  - `<prefix>_updates_available`: Number of available package updates
  - `<prefix>_security_updates_available`: Number of available security updates
  - `<prefix>_seconds_since_last_update`: Seconds since last successful apt update
  - `<prefix>_reboot_required`: 1 if a reboot is required, 0 otherwise
- Collector metrics:
  - `<prefix>_collector_success`: 1 if the last collection was successful, 0 otherwise
  - `<prefix>_collector_duration_seconds`: Duration of the last collection in seconds
  - `<prefix>_collector_last_timestamp`: Timestamp of the last collection
- Configuration options:
  - Check interval
  - Listen address
  - APT check path
  - Update stamp path
  - Reboot required file
  - Log level
  - Command timeout
  - Metrics endpoint
  - Metric prefix
- Command-line flags:
  - `-config`: Path to configuration file
  - `-version`: Show version information
  - `-skip-path-validation`: Skip validation of file paths
- Comprehensive test suite
- GitHub Actions workflows for CI and releases
- Makefile with targets for building, testing, and releasing
- Grafana dashboard for visualizing metrics

### Changed
- Improved error handling throughout the codebase
- Enhanced configuration validation
- Better path handling with absolute path resolution

[v0.1.0]: https://github.com/ncecere/apt-exporter/releases/tag/v0.1.0
