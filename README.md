# Prometheus Metrics Generator

**NOTE:** Work in Progress, unstable!!!

A CLI and a library to generate Prometheus metrics.

The main idea is to create a library that allows us to compose time series.
A user can then leverage the library to build whatever they want.
Accompaning the library, there is a CLI tool that can expose metrics configured in a config file.
It also has the option to enable an HTTP server on a configurable port with an API to modify the series being
generated (either create new ones, delete or modify existing ones). These changes are not persistent and will be lost
once the service starts.

There are constructs to allow the user to easily put the series together.
Once those series have been put together, one can either expose those time series using the usual HTTP endpoint or can
make use of the remote write protocol to inject the time series directly into prometheus.

## Use cases

* Use the tool to get a better understanding of how prometheus works. In particular, how calculations on timeseries are
performed.
* Aid in the development of Grafana dashboards. Generate data that can be used to be displayed in grafana. Grafana
allows for random data to be generated, but if one wants to see the information using prometheus as a source, one is out
of luck. Usually, we want to see how the dashboard we put together looks like when there is data in it, but without
swapping the promQL queries with the random data generator!
* Test alerting. One can put together fake time series to test whether alerts will go off in the expected way!

When putting grafana dashboards, if following GitOps, the reviewed is left with a JSON blob. Ofc this won't be reviewed
and tested. So instead, for each MR, the pipeline should generate a temporary dashboard with the contents and give it
a temporary name. It should also generate predefined data, so the dashboard can be properly reviewed when there is in
fact data! Finally, the dashboard should be deleted at some point. Either define an expiry time, or delete the
dashboard when the MR is merged. Or both, in case the MR doesn't get merged!

## Testing setup

To run prometheus and grafana, you should execute the following command:

```sh
docker compose -f test/docker-compose.yaml up
```

When configuring the prometheus datasource in Grafana, make sure to use the following URL: `http://prometheus:9090`.
This is because grafana can only access prometheus via it's docker name.

For the PromLens UI, you should use the following URL: `http://localhost:9090`.
This is because the PromLens webpage is the one making the requests, and therefore can only access prometheus via the
exposed port on the host system.

URLs:

* Prometheus UI: <http://localhost:9090>
* Grafana: <http://localhost:3000>
* PromLens: <http://localhost:8080>

## Prometheus Remote Write

Docs:

* <https://prometheus.io/docs/concepts/remote_write_spec/>
* <https://prometheus.io/docs/practices/remote_write/>
* <https://prometheus.io/docs/prometheus/latest/feature_flags/#remote-write-receiver>
* <https://prometheus.io/docs/prometheus/latest/storage/#remote-storage-integrations>
* <https://prometheus.io/docs/prometheus/latest/querying/api/#remote-write-receiver>

Golang libraries:

* <https://github.com/m3dbx/prometheus_remote_client_golang>
* <https://github.com/castai/promwrite>

There is barely any code to implement the prometheus remote write.
Instead of having these dependencies, let's write the code ourselves!

## Nomenclature

* `Sample` - A measure containing a pair of (timestamp in milliseconds, value as a float64).
* `Label` - A pair of (key, value).
* `Series` or `TimeSeries` - A list of samples, identified by a unique set of labels (including the `__name__` label).
* `Metric Family` or `Metric Name` - The metric name (without any labels). A metric family can be, and usually is, made
up of several timeseries.
* `Scrape` - A snapshot of samples at a given point in time for a particular system being observed.
* `Metric Type` - The type of the metric (e.g.: counter, gauge, histogram, summary).

## Components

This library is made up of several components.

### Scraper

A Scraper generates scrapes that are then passed into a DataIterator to generate the time series values.

The scraper is the building block that allows data iterators to generate samples.

### Data Iterator

A data iterator can either be discrete or continuous.

A discrete data iterator, generates a new data point for each scrape.

A continuous data iterator, calculates what the measure should be at the point of the scrape timestamp.

#### Discrete Data Iterator

There are various types of discrete data iterators.

The following is a list of the supported discrete data iterators:

* Linear Segment
* Custom Values
* Random
* Void
* Join
* Loop

## Ideas

* Add other discrete segments, like the ability to Add, Subtract, multiply and Divide 2 different segments.
* Create trignometric related segments (ex: sin, cos, tan, etc).
* Add support for Histogram!
* Add same functionality to continuous segments as the one we have in the discrete segments.
