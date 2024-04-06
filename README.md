# Prometheus Metrics Generator

**NOTE:** Unfinished work!

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

## Components

### Scraper

A Scraper generates scrapes that are then passed into a DataIterator to generate the time series values.
