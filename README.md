# Changedetection.io Prometheus Exporter
This application exports the latest prices of watches configured inside a [changedetection.io](https://changedetection.io) instance as prometheus-compatible metrics. Additionally, the exporter exposes metrics regarding scraping statistics and certain system information in order to enable system monitoring as well.

- Exposes price metrics for eligible watches (must be of type @Offer).
- Exposes scrape metrics to monitor performance of your watches.
- Monitor instance metrics like queue size, uptime or overdue watches.
- Visualize watches with the same name (but different sources) to compare price developments.

Right now, there are no plans for further development (maybe an example Grafana dashboard at some point in time). If you feel something is missing, feel free to open up an issue or (even better) a pull request!

## Installation
The recommended way to run the exporter is as a part of your docker-compose stack or as a pod on your k8s cluster. Of course, you can also run it directly as a binary on your host, you have to manually compile this one yourself though.

### Docker (with docker-compose)
```yaml
services:
  changedection-exporter:
    image: ghcr.io/schaermu/changedetection.io-exporter:latest
    container_name: changedetection-exporter
    restart: unless-stopped
    environment:
      - CDIO_API_BASE_URL=http://changedetection
      - CDIO_API_KEY=...
    depends_on:
      changedetection:
        condition: service_started
```
This example assumes that you have started an instance of changedetection.io within the same docker network reachable via hostname `changedetection`.

### Kubernetes (tbd)
```
```

The exporter is configured using the following environment variables:
|Environment Variable|Default value|Mandatory?|
|---|---|---|
|`CDIO_API_BASE_URL`|-|yes|
|`CDIO_API_KEY`|-|yes|
|`PORT`|`9123`|no|
|`LOG_LEVEL`|`info`|no|

For all scenarios, setting both the `CDIO_API_BASE_URL` and a `CDIO_API_KEY` environment variable is mandatory, and the exporter will panic on startup if any of those is missing.

## Usage
Metrics can be access by requesting the path `/metrics` using the exporter's hostname and its configured port (or the default one of 9123).

If you want to read those metrics into Prometheus or VictoriaMetrics, you have to configure a scraper like this:

```yml
scrape_configs:
  - job_name: "changedetection"
    static_configs:
      - targets: ["changedetectionio-exporter:9123"]
```
If you haven't got any watches registered on your changedetection.io instance, you will simply get the system metrics read in:
|Metric name|Labels|Type|
|---|---|---|
|`changedetectionio_system_uptime`|`version`|Gauge|
|`changedetectionio_system_watch_count`|`version`|Gauge|
|`changedetectionio_system_overdue_watch_count`|`version`|Gauge|
|`changedetectionio_system_queue_size`|`version`|Gauge|

The label `version` contains the running version of changedetection.io (i.e. 0.45.6).

If you have any watches registered, you will additionally get the following metrics for each of those:
|Metric name|Labels|Type|
|---|---|---|
|`changedetectionio_watch_check_count`|`title`,`source`|Counter|
|`changedetectionio_watch_fetch_time`|`title`,`source`|Gauge|
|`changedetectionio_watch_notification_alert_count`|`title`,`source`|Counter|
|`changedetectionio_watch_last_check_status`|`title`,`source`|Gauge|
|`changedetectionio_watch_price`|`title`,`source`|Gauge|

**IMPORTANT**: the metric `changedetectionio_watch_price` will ONLY be exposed for watches that return price information in the shape of a JSON object (with the attribute `@type` set to `Offer`).

The label `title` should be pretty self-explanatory, it simply contains the title from changedetection.io. In order to make sure all those metrics are unique, an additional label `source` is being exported. It contains the **host-part** of the monitored URL (i.e. www.foobar.org, so including the subdomain).

## Contributing
There are two ways you can build and run the exporter locally: using the binary build or a docker image. For both options, there are `Makefile` targets:
```bash
# clean artifacts, run tests and build binary
$ make

# clean artifacts, compile and run binary
$ make run

# create a docker image (see note below)
$ make docker
```

**IMPORTANT**: before running the docker build, you have to bootstrap the multi-platform build system in docker by running `docker buildx create --name multi-builder --bootstrap --use`.

To run the tests, you can leverage `Makefile` targets as well:
```bash
# run tests
$ make tests

# run tests with code coverage output (coverage.html)
$ make cover

# run tests in watch mode to re-run on all code-changes
$ make watch
```
