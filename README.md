# GitLab Token Metrics Exporter

## Overview

This Go project fetches GitLab token metrics and exports them for Prometheus.

## Features

- **Metric Collection**: Gathers GitLab token info, including creation and expiration dates.
- **Prometheus Integration**: Exposes metrics compatible with Prometheus.
- **Environment Config**: Configured via environment variables.

## Prerequisites

- [Go](https://golang.org/) installed
- [Prometheus](https://prometheus.io/) for scraping metrics

## Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/gitlab-token-metrics-exporter.git

2. Change into the project directory:
   ```bash
   cd gitlab-token-metrics-exporter
3. Build and run
   ```bash
   go build && ./gitlab-token-metrics-exporter

Metrics available at http://localhost:9960/metrics.

## Prometheus Configuration
    
    - job_name: 'gitlab-token-metrics'
      static_configs:
        - targets: ['localhost:9960']

### Metrics

    Request Duration: Duration of GitLab API requests.
    API Errors: Total errors when making GitLab API requests.
    Creation Date: Creation date of GitLab token in seconds since epoch.
    Expiry Date: Expiry date of GitLab token in seconds since epoch.
    
