# Datadog Log Proxy
This proxy sends [Monolog](https://github.com/Seldaek/monolog) messages in batches to the Datadog API. The log messages will be send *non blocking* over **UDP** to this proxy and than over **TCP** to the Datadog API. To reduce the amount of requests, messages will be batched and send in one request.

## PHP Monolog Handler
This proxy is build to use in conjunction with the [Monolog Datadog UDP Handler](https://github.com/YourSurpriseCom/monolog-dd-udp-handler).


```PHP
<?php

use Monolog\Logger;
use YourSurpriseCom\Monolog\DatadogUdp\Handler\DataDogUdpHandler;

$logger = new Logger('my_logger');
$handler = new DataDogUdpHandler("<proxy host>",1053);
$logger->pushHandler($handler);

$logger->info("This is log message is send non blocking over UDP to datadog!");
```

More information about the Monolog Datadog UDP Handler can be found [here](https://github.com/YourSurpriseCom/monolog-dd-udp-handler).

## Data flow
This proxy is meant to place between your PHP application and Datadog. The data flow is as following:


```
+-----+       +-------------------------+             +--------------+             +-------------+ 
| PHP |  ==>  | Monolog Datadog Handler |  ==> (UDP)  | dd-log-proxy |  ==> (TCP)  | Datadog API |
+-----+       +-------------------------+             +--------------+             +-------------+ 
```

## Usage
The proxy is available as a docker container and can be started as following:

```shell
foo@bar:~$ docker run --name datadog-log-proxy \
    -e DD_SITE=datadoghq.eu \
    -e DD_API_KEY="<DD_API_KEY>" \
    yoursurprise/dd-log-proxy:latest
```

### Configuration
The proxy will use environment variable for configuration, there are at least 2 variables that are required:

1. `DD_SITE` (The regional site for Datadog customers.)
2. `DD_API_KEY="<DD_API_KEY>"` (API key, see https://docs.datadoghq.com/account_management/api-app-keys/#api-keys)

#### Datadog Regional Site
By default this logger will send its data to the EU region.
This can be changed by setting de environment variable `DD_SITE`, the following values are possible:

* `datadoghq.com`
* `us3.datadoghq.com`
* `us5.datadoghq.com`
* `ap1.datadoghq.com`
* `datadoghq.eu`
* `ddog-gov.com`

#### Batch Configuration
This logger will create log batches to lower the requests to Datadog. The size of the batch and the wait time to send incomplete batches can be changed by setting the following environment variables:

```
BATCH_SIZE=50
BATCH_WAIT_IN_SECONDS=5
```

#### Debug logger
To enable the debug log, set the following environment variable:

`DEBUG_LEVEL=debug`

## Development
Install the required tooling, and run `make` for a list of development commands.

### Required tooling
```
go install github.com/kisielk/errcheck@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/oligot/go-mod-upgrade@latest
```