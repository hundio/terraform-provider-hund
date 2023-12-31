---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hund_metric_provider Resource - terraform-provider-hund"
subcategory: ""
description: |-
  MetricProviders gather metrics from a configured service for viewing on the Status Page.
---

# hund_metric_provider (Resource)

MetricProviders gather metrics from a configured service for viewing on the Status Page.

## Example Usage

```terraform
resource "hund_component" "component0" {
  name        = "Terraform Component"
  description = "Created by Terraform"

  group = "63f6c4938fbb652be74a9ae6"

  watchdog = {
    service = {
      webhook = { webhook_key = "KEY" }
    }
  }
}

resource "hund_metric_provider" "default" {
  watchdog = hund_component.component0.watchdog.id
  service  = { webhook = hund_component.component0.watchdog.service.webhook }

  # Managing a Watchdog's default MetricProvider requires importation.
  default = true

  instances = {}
}

resource "hund_metric_provider" "builtin" {
  watchdog = hund_component.component0.watchdog.id

  service = { builtin = {} }

  instances = {
    percent_uptime     = { enabled = true },
    incidents_reported = { enabled = false }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `watchdog` (String) The Watchdog that owns this MetricProvider.

### Optional

- `default` (Boolean) When true, denotes that this MetricProvider is the default MetricProvider of the Watchdog. This implies that they share the same service configuration, which the MetricProvider inherits from the Watchdog. This MetricProvider is created **automatically**, depending on the Watchdog, and cannot be deleted without also deleting the Watchdog.

~> Default MetricProviders cannot be created directly, and must be imported to be managed by Terraform. Deleting a default MetricProvider from your Terraform configuration will only remove the resource from the state.
- `instances` (Attributes Map) A Map of MetricInstances, which describe each Metric that the MetricProvider provides. The keys of this Map define the slugs of each provided metric. (see [below for nested schema](#nestedatt--instances))
- `service` (Attributes) The service configuration for this MetricProvider, which describes how the given `instances` are provided. (see [below for nested schema](#nestedatt--service))

### Read-Only

- `id` (String) The ObjectId of this MetricProvider.

<a id="nestedatt--instances"></a>
### Nested Schema for `instances`

Optional:

- `aggregation` (String) The kind of aggregation method to use in case multiple displayed data points share the same time-axis value (depending on the axis configured for time, by default x).

-> this field does not have any effect on the underlying data; it is purely cosmetic, and applied only when viewing the data on the status page.
- `enabled` (Boolean) Whether or not to show this metric on the Component that uses it (through the Watchdog).
- `interpolation` (String) The kind of interpolation to use between points displayed in the graph (line plots only). One of `linear`, `step`, `basis`, `bundle`, or `cardinal`.
- `plot_type` (String) The kind of visualization to display the metric with. One of `line` or `bar`.
- `title` (String) The title of the metric, displayed above its graph on the status page, in the default translation.
- `title_translations` (Map of String) The title of the metric, displayed above its graph on the status page, translated into multiple languages. Map keys express the language each string value is to be interpreted in. The `original` field of this map denotes the language used for the non-`_translations` version of this attribute.
- `top_level_enabled` (Boolean) Whether or not to show this metric on the status page home.
- `x_title` (String) The title of the x-axis of this metric, in the default translation.
- `x_title_translations` (Map of String) The title of the x-axis of this metric, translated into multiple languages. Map keys express the language each string value is to be interpreted in. The `original` field of this map denotes the language used for the non-`_translations` version of this attribute.
- `x_type` (String) The type of quantity represented by the x-axis. One of `time` or `measure`.
- `y_supremum` (Number) The least upper bound to display the y-axis on. The metric will always display up to at least this value on the y-axis regardless of the graphed data. If the graph exceeds this value, then the bound will be raised as much as necessary to accommodate the data.
- `y_title` (String) The title of the y-axis of this metric, in the default translation.
- `y_title_translations` (Map of String) The title of the y-axis of this metric, translated into multiple languages. Map keys express the language each string value is to be interpreted in. The `original` field of this map denotes the language used for the non-`_translations` version of this attribute.
- `y_type` (String) The type of quantity represented by the y-axis. One of `time` or `measure`.

Read-Only:

- `definition_slug` (String) A descriptive string that identifies the metric definition this instance derives from (e.g. `http.tcp_connection_time`, `apdex`, etc.).
- `id` (String) The ObjectId of this MetricInstance.
- `slug` (String) A string that uniquely identifies this MetricInstance by referencing the `definition_slug` and MetricProvider `id`.


<a id="nestedatt--service"></a>
### Nested Schema for `service`

Optional:

- `builtin` (Attributes) The builtin Hund metric provider, which provides metrics based on recorded uptime and incidents. (see [below for nested schema](#nestedatt--service--builtin))
- `dns` (Attributes) A Hund Native Monitoring DNS Check. (see [below for nested schema](#nestedatt--service--dns))
- `http` (Attributes) A Hund Native Monitoring HTTP Check. (see [below for nested schema](#nestedatt--service--http))
- `icmp` (Attributes) A Hund Native Monitoring ICMP Check. (see [below for nested schema](#nestedatt--service--icmp))
- `pingdom` (Attributes) A [pingdom](https://www.pingdom.com) service. (see [below for nested schema](#nestedatt--service--pingdom))
- `tcp` (Attributes) A Hund Native Monitoring TCP Check. (see [below for nested schema](#nestedatt--service--tcp))
- `udp` (Attributes) A Hund Native Monitoring UDP Check. (see [below for nested schema](#nestedatt--service--udp))
- `updown` (Attributes) An [Updown.io](https://updown.io) service. (see [below for nested schema](#nestedatt--service--updown))
- `uptimerobot` (Attributes) An [Uptime Robot](https://uptimerobot.com) service. (see [below for nested schema](#nestedatt--service--uptimerobot))
- `webhook` (Attributes) A [webhook](https://hund.io/help/documentation/incoming-webhook-metrics) service. (see [below for nested schema](#nestedatt--service--webhook))

<a id="nestedatt--service--builtin"></a>
### Nested Schema for `service.builtin`


<a id="nestedatt--service--dns"></a>
### Nested Schema for `service.dns`

Required:

- `record_type` (String) The type of DNS record to query for on the target.
- `regions` (Set of String) The regions you would like the target to be checked from. All regions are
weighted equally when calculating the outcome of a check. Currently, a single
check can use up to 8 regions simultaneously. Using at least two regions for a
single check is recommended in order to confirm failures across regions.

-> Each check may use up to **three** regions at no extra cost. Each region added to this check beyond the base three will incur an additional cost. For specific pricing information, please visit the [pricing](https://hund.io/pricing) page.
- `target` (String) The domain/IP address that will be queried. IP addresses do not need to be
converted to the `z.y.x.w.in-addr.arpa` format, as this will be done
automatically; however, both formats are accepted.

Optional:

- `consecutive_check_degraded_threshold` (Number) The number of consecutive failed checks required before posting a "degraded"
status.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When `null`, denotes that this check will not use a "degraded" stage
when encountering check failures.

  When 0, denotes that this check will post "degraded" upon the first check failure.
- `consecutive_check_outage_threshold` (Number) The number of consecutive failed checks required before posting an "outage"
status. If `consecutive_check_degraded_threshold` is non-null, then the outage
will only be posted after degraded has posted according to its own threshold.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When 0, denotes that this check will post "outage" upon the first check failure
(or the first check failure after "degraded" has been posted in case
`consecutive_check_degraded_threshold` is set).
- `frequency` (Number) The frequency of the check in milliseconds. The maximum frequency is every 30
seconds.

-> Any frequency greater than every 60 seconds will force the component
to become High-Frequency, at an additional cost. For specific pricing
information, please visit the [pricing](https://hund.io/pricing) page.
- `nameservers` (List of String) An optional list of nameservers to make DNS queries with. This field is
ignored by SOA queries since they use the nameservers yielded by querying NS
on the target.
- `percentage_regions_failed_threshold` (Number) The percentage of regions that must report a failed check before the entire
check can be considered failed. Requiring at least two regions for this
threshold is recommended in order to confirm failures across regions.
- `response_containment` (String) Whether `all` of the assertions in `responses_must_contain` must match the DNS response,
or rather just `any` of them (i.e. at least one).
- `responses_must_contain` (Set of String) A set of assertions to make against the records yielded by the query. The
format of these assertions is *similar* to DNS record syntax, but is
slightly simplified and allows for only asserting parts of a record's RDATA,
rather than the entire thing. The check will fail depending on the value of
`response_containment`.

  This field is ignored by the SOA check, as it does not use assertions to
determine the validity of SOA records. Instead, we ensure that every
nameserver reported by querying NS on the target reports the same SOA serial.
If your target's nameservers report conflicting SOA serials, we consider the
check failed.

  **Example Assertions (for MX record type):**
```json
[
  "10 mail.example.com",
  "spool.example.com",
  "mail2.example.com"
]
```

  Note above how we can assert both the priority and domain (*without* the
terminating period required by canonical DNS) of an MX record, or instead
simply the domain.
- `timeout` (Number) The maximum number of milliseconds the check should wait on the host before
failing.


<a id="nestedatt--service--http"></a>
### Nested Schema for `service.http`

Required:

- `regions` (Set of String) The regions you would like the target to be checked from. All regions are
weighted equally when calculating the outcome of a check. Currently, a single
check can use up to 8 regions simultaneously. Using at least two regions for a
single check is recommended in order to confirm failures across regions.

-> Each check may use up to **three** regions at no extra cost. Each region added to this check beyond the base three will incur an additional cost. For specific pricing information, please visit the [pricing](https://hund.io/pricing) page.
- `target` (String) The host the check will make calls to.

Optional:

- `consecutive_check_degraded_threshold` (Number) The number of consecutive failed checks required before posting a "degraded"
status.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When `null`, denotes that this check will not use a "degraded" stage
when encountering check failures.

  When 0, denotes that this check will post "degraded" upon the first check failure.
- `consecutive_check_outage_threshold` (Number) The number of consecutive failed checks required before posting an "outage"
status. If `consecutive_check_degraded_threshold` is non-null, then the outage
will only be posted after degraded has posted according to its own threshold.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When 0, denotes that this check will post "outage" upon the first check failure
(or the first check failure after "degraded" has been posted in case
`consecutive_check_degraded_threshold` is set).
- `follow_redirects` (Boolean) Follow any HTTP redirects given by the requested target. Please note that this check will only follow up to 9 redirects.
- `frequency` (Number) The frequency of the check in milliseconds. The maximum frequency is every 30
seconds.

-> Any frequency greater than every 60 seconds will force the component
to become High-Frequency, at an additional cost. For specific pricing
information, please visit the [pricing](https://hund.io/pricing) page.
- `headers` (Map of String) A list of additional HTTP headers to send to the target. The following list of
header names are reserved and cannot be set by a check:

```
Accept-Charset
Accept-Encoding
Authentication
Connection
Content-Length
Date
Host
Keep-Alive
Origin
Proxy-.*
Sec-.*
Referer
TE
Trailer
Transfer-Encoding
User-Agent
Via
```
- `password` (String, Sensitive) An optional HTTP Basic Authentication password.
- `percentage_regions_failed_threshold` (Number) The percentage of regions that must report a failed check before the entire
check can be considered failed. Requiring at least two regions for this
threshold is recommended in order to confirm failures across regions.
- `response_body_must_contain` (String) This field supports two different matching modes (given by
`response_body_must_contain_mode`):

  `exact`: If the requested page does not contain this exact (case-sensitive)
string, then the check will fail.

  `regex`: If the requested page does not match against the given regex, then
the check will fail. [Click here](https://hund.io/help/documentation/regular-expressions) for
more information on the use and supported syntax of Hund regexes.
- `response_body_must_contain_mode` (String) The response containment mode; either `exact` or `regex`. The modes are discussed
under `response_body_must_contain`.
- `response_code_must_be` (Number) If the requested page does not return this response code, then the check will
fail.
- `ssl_verify_peer` (Boolean) Require the target's TLS certificate to be valid.
- `timeout` (Number) The maximum number of milliseconds the check should wait on the host before
failing.
- `username` (String) An optional HTTP Basic Authentication username.


<a id="nestedatt--service--icmp"></a>
### Nested Schema for `service.icmp`

Required:

- `regions` (Set of String) The regions you would like the target to be checked from. All regions are
weighted equally when calculating the outcome of a check. Currently, a single
check can use up to 8 regions simultaneously. Using at least two regions for a
single check is recommended in order to confirm failures across regions.

-> Each check may use up to **three** regions at no extra cost. Each region added to this check beyond the base three will incur an additional cost. For specific pricing information, please visit the [pricing](https://hund.io/pricing) page.
- `target` (String) The host the check will make calls to.

Optional:

- `consecutive_check_degraded_threshold` (Number) The number of consecutive failed checks required before posting a "degraded"
status.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When `null`, denotes that this check will not use a "degraded" stage
when encountering check failures.

  When 0, denotes that this check will post "degraded" upon the first check failure.
- `consecutive_check_outage_threshold` (Number) The number of consecutive failed checks required before posting an "outage"
status. If `consecutive_check_degraded_threshold` is non-null, then the outage
will only be posted after degraded has posted according to its own threshold.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When 0, denotes that this check will post "outage" upon the first check failure
(or the first check failure after "degraded" has been posted in case
`consecutive_check_degraded_threshold` is set).
- `frequency` (Number) The frequency of the check in milliseconds. The maximum frequency is every 30
seconds.

-> Any frequency greater than every 60 seconds will force the component
to become High-Frequency, at an additional cost. For specific pricing
information, please visit the [pricing](https://hund.io/pricing) page.
- `ip_version` (String) The IP version to use when pinging.
- `percentage_failed_threshold` (Number) The percentage of addresses at the given target that must fail for a region to be counted as failed. This option only matters when there are multiple IP addresses behind the target when the target is a domain.
- `percentage_regions_failed_threshold` (Number) The percentage of regions that must report a failed check before the entire
check can be considered failed. Requiring at least two regions for this
threshold is recommended in order to confirm failures across regions.
- `timeout` (Number) The maximum number of milliseconds the check should wait on the host before
failing.


<a id="nestedatt--service--pingdom"></a>
### Nested Schema for `service.pingdom`

Required:

- `api_token` (String, Sensitive) The Pingdom API v3 key.
- `check_id` (String) The ID of the check to pull status from on Pingdom.

Optional:

- `check_type` (String) The type of the Pingdom check. `check` denotes a normal Pingdom uptime check, and `transactional` denotes a Pingdom TMS check.


<a id="nestedatt--service--tcp"></a>
### Nested Schema for `service.tcp`

Required:

- `port` (Number) The port at the target to connect to.
- `regions` (Set of String) The regions you would like the target to be checked from. All regions are
weighted equally when calculating the outcome of a check. Currently, a single
check can use up to 8 regions simultaneously. Using at least two regions for a
single check is recommended in order to confirm failures across regions.

-> Each check may use up to **three** regions at no extra cost. Each region added to this check beyond the base three will incur an additional cost. For specific pricing information, please visit the [pricing](https://hund.io/pricing) page.
- `target` (String) The host the check will make calls to.

Optional:

- `consecutive_check_degraded_threshold` (Number) The number of consecutive failed checks required before posting a "degraded"
status.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When `null`, denotes that this check will not use a "degraded" stage
when encountering check failures.

  When 0, denotes that this check will post "degraded" upon the first check failure.
- `consecutive_check_outage_threshold` (Number) The number of consecutive failed checks required before posting an "outage"
status. If `consecutive_check_degraded_threshold` is non-null, then the outage
will only be posted after degraded has posted according to its own threshold.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When 0, denotes that this check will post "outage" upon the first check failure
(or the first check failure after "degraded" has been posted in case
`consecutive_check_degraded_threshold` is set).
- `frequency` (Number) The frequency of the check in milliseconds. The maximum frequency is every 30
seconds.

-> Any frequency greater than every 60 seconds will force the component
to become High-Frequency, at an additional cost. For specific pricing
information, please visit the [pricing](https://hund.io/pricing) page.
- `ip_version` (String) The IP version to use when calling the target.
- `percentage_regions_failed_threshold` (Number) The percentage of regions that must report a failed check before the entire
check can be considered failed. Requiring at least two regions for this
threshold is recommended in order to confirm failures across regions.
- `response_must_contain` (String) This field supports two different matching modes (given by `response_must_contain_mode`):

  `exact`: Text that the response from the target must contain exactly
(case-sensitive). In exact match mode, this field supports
[escape codes](https://hund.io/help/documentation/text-field-escape-codes).

  `regex`: A regex that the response from the target must match against.
[Click here](https://hund.io/help/documentation/regular-expressions) for more information on
the use and supported syntax of Hund regexes.

  If you send data and expect the target to reply, you must populate this field.
Leaving this field blank will prevent the check from receiving data from the
target unless forced to wait for an initial response.

  The "response" from the target that this text is asserted against will be the
response from the target *after* sending data. If data is not sent to the
target, this text is asserted against the *initial* response.
- `response_must_contain_mode` (String) The response containment mode; either `exact` or `regex`. The modes are discussed
under `response_must_contain`.
- `send_data` (String) Optional data to send to the target after connecting. If this field is left
blank, nothing is sent to the target after connecting. This field supports [escape codes](https://hund.io/help/documentation/text-field-escape-codes).
- `timeout` (Number) The maximum number of milliseconds the check should wait on the host before
failing.
- `wait_for_initial_response` (Boolean) Whether or not to wait for an initial response from the target before sending
data or closing the connection.


<a id="nestedatt--service--udp"></a>
### Nested Schema for `service.udp`

Required:

- `port` (Number) The port at the target to connect to.
- `regions` (Set of String) The regions you would like the target to be checked from. All regions are
weighted equally when calculating the outcome of a check. Currently, a single
check can use up to 8 regions simultaneously. Using at least two regions for a
single check is recommended in order to confirm failures across regions.

-> Each check may use up to **three** regions at no extra cost. Each region added to this check beyond the base three will incur an additional cost. For specific pricing information, please visit the [pricing](https://hund.io/pricing) page.
- `send_data` (String) Data to send to the target after connecting. Unlike in `tcp`, this
field is required. This field supports [escape codes](https://hund.io/help/documentation/text-field-escape-codes).
- `target` (String) The host the check will make calls to.

Optional:

- `consecutive_check_degraded_threshold` (Number) The number of consecutive failed checks required before posting a "degraded"
status.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When `null`, denotes that this check will not use a "degraded" stage
when encountering check failures.

  When 0, denotes that this check will post "degraded" upon the first check failure.
- `consecutive_check_outage_threshold` (Number) The number of consecutive failed checks required before posting an "outage"
status. If `consecutive_check_degraded_threshold` is non-null, then the outage
will only be posted after degraded has posted according to its own threshold.

  Note that regardless of threshold settings, a component will post "operational"
whenever a check succeeds, thus resetting the consecutive check failure count.

  When 0, denotes that this check will post "outage" upon the first check failure
(or the first check failure after "degraded" has been posted in case
`consecutive_check_degraded_threshold` is set).
- `frequency` (Number) The frequency of the check in milliseconds. The maximum frequency is every 30
seconds.

-> Any frequency greater than every 60 seconds will force the component
to become High-Frequency, at an additional cost. For specific pricing
information, please visit the [pricing](https://hund.io/pricing) page.
- `ip_version` (String) The IP version to use when calling the target.
- `percentage_regions_failed_threshold` (Number) The percentage of regions that must report a failed check before the entire
check can be considered failed. Requiring at least two regions for this
threshold is recommended in order to confirm failures across regions.
- `response_must_contain` (String) This field supports two different matching modes (given by `response_must_contain_mode`):

  `exact`: Text that the response from the target must contain exactly
(case-sensitive). In exact match mode, this field supports
[escape codes](https://hund.io/help/documentation/text-field-escape-codes).

  `regex`: A regex that the response from the target must match against.
[Click here](https://hund.io/help/documentation/regular-expressions) for more information on
the use and supported syntax of Hund regexes.

  Leaving this field blank will still cause the check to wait for a response
from the target after sending data, though no assertions will be made about
the payload of the response.
- `response_must_contain_mode` (String) The response containment mode; either `exact` or `regex`. The modes are discussed
under `response_must_contain`.
- `timeout` (Number) The maximum number of milliseconds the check should wait on the host before
failing.


<a id="nestedatt--service--updown"></a>
### Nested Schema for `service.updown`

Required:

- `monitor_api_key` (String, Sensitive) An Updown.io monitor API key. This API key can be read-only.
- `monitor_token` (String) An Updown.io monitor token to retrieve status from.


<a id="nestedatt--service--uptimerobot"></a>
### Nested Schema for `service.uptimerobot`

Required:

- `monitor_api_key` (String, Sensitive) An Uptime Robot monitor API key to retrieve status from.


<a id="nestedatt--service--webhook"></a>
### Nested Schema for `service.webhook`

Optional:

- `webhook_key` (String, Sensitive) The key to use for this webhook, expected in request headers.
