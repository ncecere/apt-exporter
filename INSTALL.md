# Installing APT Exporter

This document provides detailed instructions for installing the APT Exporter on your system.

## Prerequisites

- Linux (Debian/Ubuntu-based system with APT)
- `update-notifier-common` package (provides the `/usr/lib/update-notifier/apt-check` script)

### Installing Prerequisites

On Debian/Ubuntu systems, install the required package:

```bash
sudo apt-get update
sudo apt-get install update-notifier-common
```

This package provides the apt-check script that the exporter uses to collect information about available updates.

## Installation Methods

### Method 1: Using Pre-built Binaries from GitHub Releases

This is the recommended method for most users.

1. **Download the latest release**

   Visit the [GitHub Releases page](https://github.com/ncecere/apt-exporter/releases) and download the appropriate binary for your architecture. For example, for a 64-bit x86 system:

   ```bash
   # Create a directory for the exporter
   mkdir -p ~/apt-exporter
   cd ~/apt-exporter
   
   # Download the latest release (replace X.Y.Z with the actual version, e.g., 0.1.0)
   curl -L https://github.com/ncecere/apt-exporter/releases/download/vX.Y.Z/apt-exporter_X.Y.Z_linux_amd64.tar.gz -o apt-exporter.tar.gz
   
   # Extract the archive
   tar -xzf apt-exporter.tar.gz
   ```

2. **Make the binary executable**

   ```bash
   chmod +x apt-exporter
   ```

3. **Move the binary to a system location** (optional)

   ```bash
   sudo mv apt-exporter /usr/local/bin/
   ```

4. **Create a configuration directory and copy the config file**

   ```bash
   sudo mkdir -p /etc/apt-exporter
   sudo cp config.yml /etc/apt-exporter/
   ```

5. **Edit the configuration file as needed**

   ```bash
   sudo nano /etc/apt-exporter/config.yml
   ```

### Method 2: Building from Source

If you prefer to build from source or need to make custom modifications:

1. **Install Go**

   Make sure you have Go 1.18 or higher installed:

   ```bash
   # Check Go version
   go version
   
   # Install Go if needed (on Debian/Ubuntu)
   sudo apt-get update
   sudo apt-get install golang-go
   ```

2. **Clone the repository**

   ```bash
   git clone https://github.com/ncecere/apt-exporter.git
   cd apt-exporter
   ```

3. **Build the binary**

   ```bash
   make build
   ```

   Or manually:

   ```bash
   go build -o apt-exporter cmd/apt-exporter/main.go
   ```

4. **Install the binary**

   ```bash
   sudo cp apt-exporter /usr/local/bin/
   sudo mkdir -p /etc/apt-exporter
   sudo cp config.yml /etc/apt-exporter/
   ```

## Setting Up as a System Service

### Using Systemd

1. **Create a systemd service file**

   ```bash
   sudo nano /etc/systemd/system/apt-exporter.service
   ```

2. **Add the following content to the file**

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
   
   > **Note**: The exporter needs to run as root to access system files like apt-check. If you prefer not to run as root, you can create a dedicated user with appropriate permissions to access these files.

3. **Enable and start the service**

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable apt-exporter
   sudo systemctl start apt-exporter
   ```

4. **Check the service status**

   ```bash
   sudo systemctl status apt-exporter
   ```

## Verifying the Installation

1. **Check if the exporter is running**

   ```bash
   ps aux | grep apt-exporter
   ```

2. **Test the metrics endpoint**

   ```bash
   curl http://localhost:9100/metrics
   ```

   You should see Prometheus metrics output, including the APT metrics.

## Troubleshooting

### Common Issues

1. **Permission denied when accessing apt-check**

   The exporter needs permission to access the apt-check script. Make sure it's running as a user with appropriate permissions or modify the systemd service to run as root.

2. **Empty or unexpected output from apt-check**

   If you see errors like "unexpected output format from apt-check" in the logs, it could be due to:
   
   - The `update-notifier-common` package is not installed
   - The apt-check script is not working correctly on your system
   
   The exporter will handle empty or unexpected output gracefully by assuming zero updates, but you should check if the apt-check script works correctly:
   
   ```bash
   # Test the apt-check script directly
   /usr/lib/update-notifier/apt-check
   
   # Check the exit code
   echo $?
   
   # Check the stderr output (apt-check outputs to stderr)
   /usr/lib/update-notifier/apt-check 2>&1
   ```
   
   The output should be in the format "X;Y" where X is the number of updates and Y is the number of security updates.

3. **Metrics not updating**

   Check the log output for errors:

   ```bash
   sudo journalctl -u apt-exporter
   ```

4. **Configuration file not found**

   Make sure the path to the configuration file is correct in your command or systemd service file.

## Updating

To update to a newer version:

1. **Download the new version**

   Follow the same steps as in the installation section to download the latest release.

2. **Replace the binary**

   ```bash
   sudo systemctl stop apt-exporter
   sudo cp apt-exporter /usr/local/bin/
   sudo systemctl start apt-exporter
   ```

## Uninstalling

To remove the APT Exporter:

```bash
sudo systemctl stop apt-exporter
sudo systemctl disable apt-exporter
sudo rm /etc/systemd/system/apt-exporter.service
sudo rm /usr/local/bin/apt-exporter
sudo rm -rf /etc/apt-exporter
