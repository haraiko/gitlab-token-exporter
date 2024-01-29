GitLab Token Metrics Exporter
Overview

This Go project serves as a metrics exporter for GitLab token information. It fetches data from a GitLab instance and exposes relevant metrics via Prometheus.
Features

    Metric Collection: Gathers information about GitLab tokens, including creation date and expiration date.
    Prometheus Integration: Exposes metrics in a format compatible with Prometheus.
    Configuration via Environment Variables: Utilizes environment variables for configuration, making it flexible for various environments.

Prerequisites

    Go installed
    Prometheus for scraping metrics

Configuration

Configure the project by setting the following environment variables:

    GITLAB_API_URL: The GitLab API URL for fetching token information.
    PRIVATE_TOKEN: Your GitLab private access token.
    USER_ID: The user ID for whom the tokens will be retrieved.

Create a file named token.env and populate it with the above environment variables:

env

GITLAB_API_URL=https://gitlab.example.com/api/v4/personal_access_tokens?user_id=123
PRIVATE_TOKEN=your_private_token_here
USER_ID=123

Usage

    Clone the repository:

    bash

git clone https://github.com/yourusername/gitlab-token-metrics-exporter.git

Change into the project directory:

bash

cd gitlab-token-metrics-exporter

Build and run the application:

bash

    go build
    ./gitlab-token-metrics-exporter

    The metrics will be available at http://localhost:9960/metrics.

Prometheus Configuration

Add the following job to your Prometheus configuration to scrape metrics from this exporter:

yaml

- job_name: 'gitlab-token-metrics'
  static_configs:
    - targets: ['localhost:9960']

Metrics

    gitlab_access_tokens_request_duration_seconds: Duration of GitLab API requests.
    gitlab_access_tokens_api_errors_total: Total number of errors when making GitLab API requests.
    gitlab_access_tokens_creation_date_seconds: Creation date of GitLab token in seconds since epoch.
    gitlab_access_tokens_expiry_date_seconds: Expiry date of GitLab token in seconds since epoch.
