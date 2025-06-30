# Zettelkasten Exporter Helm Chart

This Helm chart deploys the [zettelkasten-exporter](https://github.com/luissimas/zettelkasten-exporter) to Kubernetes.

It includes options to deploy VictoriaMetrics and Grafana as dependencies for a complete monitoring solution.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2+

## Installation

### Installing from OCI Registry (Recommended)

1.  **Log in to the GitHub Container Registry:**

    ```bash
    export CR_PAT=YOUR_GITHUB_PAT
    echo $CR_PAT | helm registry login ghcr.io --username YOUR_USERNAME --password-stdin
    ```

2.  **Install the Chart:**

    ```bash
    helm install my-release oci://ghcr.io/luissimas/zettelkasten-exporter --version <CHART_VERSION> -f values.yaml
    ```

    Replace `<CHART_VERSION>` with the version you want to deploy. You will still need a `values.yaml` file to configure the Zettelkasten repository URL.

### Installing from Local Files

1.  **Update Chart Dependencies:**

    ```bash
    helm dependency update ./charts/zettelkasten-exporter
    ```

2.  **Install the Chart:**

    To install the chart, you must provide the URL of your Zettelkasten Git repository. Create a `values.yaml` file to specify your configuration:

    ```yaml
    # values.yaml
    zettelkasten:
      git:
        url: "https://github.com/user/my-zettelkasten.git"

      # Provide the token directly (will be stored in a new secret).
      # The value should be base64 encoded.
      # githubToken: "YOUR_BASE64_ENCODED_TOKEN"
      #
      # Or, use an existing secret containing the token.
      existingSecret:
        name: "my-github-secret"
        key: "token"
    ```

    Install the chart with your custom values:

    ```bash
    helm install my-release -f values.yaml ./charts/zettelkasten-exporter
    ```

## Configuration

The following table lists the most common configurable parameters of the chart.

| Parameter                             | Description                                                               | Default                                                        |
| ------------------------------------- | ------------------------------------------------------------------------- | -------------------------------------------------------------- |
| `image.repository`                    | Image repository                                                          | `ghcr.io/luissimas/zettelkasten-exporter`                        |
| `image.tag`                           | Image tag                                                                 | `latest`                                                       |
| `zettelkasten.directory`              | The directory where your Zettelkasten is located.                         | `""`                                                           |
| `zettelkasten.git.url`                | The URL of your Zettelkasten Git repository.                              | `""`                                                           |
| `zettelkasten.git.branch`             | The branch to checkout.                                                   | `main`                                                         |
| `zettelkasten.githubToken`            | A base64 encoded GitHub token for private repositories.                   | `""`                                                           |
| `zettelkasten.existingSecret.name`    | The name of an existing secret containing the GitHub token.               | `""`                                                           |
| `zettelkasten.existingSecret.key`     | The key within the existing secret that holds the token.                  | `""`                                                           |
| `zettelkasten.ignoreFiles`            | A list of files to ignore when collecting metrics.                        | `[".git", ".obsidian", ".trash", "README.md"]`                  |
| `collectionInterval`                  | The interval at which to collect metrics.                                 | `5m`                                                           |
| `collectHistoricalMetrics`            | Collect historical metrics from the git history.                          | `true`                                                         |
| `ingress.enabled`                     | Enable or disable the Ingress resource.                                   | `false`                                                        |
| `victoriaMetrics.enabled`             | Enable or disable the VictoriaMetrics dependency.                         | `true`                                                         |
| `grafana.enabled`                     | Enable or disable the Grafana dependency.                                 | `true`                                                         |

## Accessing Grafana

Once the chart is deployed, follow the instructions printed in the console to get the Grafana URL and admin credentials.
