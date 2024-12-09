# Reports Service API v1 documentation

AsyncAPI document for the Reports Service that handles report requests, generates reports, and handles report request failures.

## Table of Contents

* [Servers](#servers)
* [Channels](#channels)

## Servers

### **kafkaBroker** Server

| URL | Protocol | Description |
|---|---|---|
| kafka-broker.example.com:9092 | kafka | Central Kafka broker used by the Reports Service. |

#### Security Requirements

| Type | Description | security.protocol | sasl.mechanism |
|---|---|---|---|
| - | - | PLAINTEXT | - |

## Channels

### **ReportRequested** Channel

Topic where the Reports Service listens for report generation requests.

#### `subscribe` Operation

*Reports Service listens for report requests*

##### Message `ReportRequestedMessage`

Message sent to request a report generation.

###### Payload

| Name | Type | Description | Value | Constraints | Notes |
|---|---|---|---|---|---|
| (root) | object | - | - | - | **additional properties are allowed** |
| correlationId | string | - | - | - | - |
| reportRequest | object | - | - | - | **additional properties are allowed** |
| reportRequest.clusterId | string | - | - | - | - |
| reportRequest.sinceMs | integer | - | - | format (`int64`) | - |
| reportRequest.toMs | integer | - | - | format (`int64`) | - |
| reportRequest.applicationConfiguration | array<object> | - | - | - | - |
| reportRequest.applicationConfiguration.applicationName | string | - | - | - | - |
| reportRequest.applicationConfiguration.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| reportRequest.applicationConfiguration.customPrompt | string | - | - | - | - |
| reportRequest.nodeConfiguration | array<object> | - | - | - | - |
| reportRequest.nodeConfiguration.nodeName | string | - | - | - | - |
| reportRequest.nodeConfiguration.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| reportRequest.nodeConfiguration.customPrompt | string | - | - | - | - |

> Examples of payload _(generated)_

```json
{
  "correlationId": "string",
  "reportRequest": {
    "clusterId": "string",
    "sinceMs": 0,
    "toMs": 0,
    "applicationConfiguration": [
      {
        "applicationName": "string",
        "accuracy": "HIGH",
        "customPrompt": "string"
      }
    ],
    "nodeConfiguration": [
      {
        "nodeName": "string",
        "accuracy": "HIGH",
        "customPrompt": "string"
      }
    ]
  }
}
```




### **ReportGenerated** Channel

Topic where the Reports Service publishes generated reports.

#### `publish` Operation

*Reports Service publishes generated reports*

##### Message `ReportGeneratedMessage`

Message sent when a report is successfully generated.

###### Payload

| Name | Type | Description | Value | Constraints | Notes |
|---|---|---|---|---|---|
| (root) | object | - | - | - | **additional properties are allowed** |
| correlationId | string | - | - | - | - |
| report | object | - | - | - | **additional properties are allowed** |
| report.id | string | - | - | - | - |
| report.status | string | The current state of the report | allowed (`"failed_to_generate"`, `"awaiting_generation"`, `"generated"`) | - | - |
| report.clusterId | string | - | - | - | - |
| report.sinceMs | integer | - | - | format (`int64`) | - |
| report.toMs | integer | - | - | format (`int64`) | - |
| report.requestedAtNs | integer | - | - | format (`int64`) | - |
| report.scheduledGenerationAtMs | integer | - | - | format (`int64`) | - |
| report.title | string | - | - | - | - |
| report.nodeReports | array<object> | - | - | - | - |
| report.nodeReports.node | string | - | - | - | - |
| report.nodeReports.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| report.nodeReports.customPrompt | string | - | - | - | - |
| report.nodeReports.incidents | array<object> | - | - | - | - |
| report.nodeReports.incidents.id | string | - | - | - | - |
| report.nodeReports.incidents.title | string | - | - | - | - |
| report.nodeReports.incidents.category | string | - | - | - | - |
| report.nodeReports.incidents.customPrompt | string | - | - | - | - |
| report.nodeReports.incidents.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| report.nodeReports.incidents.clusterId | string | - | - | - | - |
| report.nodeReports.incidents.nodeName | string | - | - | - | - |
| report.nodeReports.incidents.summary | string | - | - | - | - |
| report.nodeReports.incidents.recommendation | string | - | - | - | - |
| report.nodeReports.incidents.urgency | string | - | allowed (`"LOW"`, `"MEDIUM"`, `"HIGH"`) | - | - |
| report.nodeReports.incidents.sources | array<object> | - | - | - | - |
| report.nodeReports.incidents.sources.timestamp | integer | - | - | format (`int64`) | - |
| report.nodeReports.incidents.sources.content | string | - | - | - | - |
| report.nodeReports.incidents.sources.filename | string | - | - | - | - |
| report.applicationReports | array<object> | - | - | - | - |
| report.applicationReports.applicationName | string | - | - | - | - |
| report.applicationReports.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| report.applicationReports.customPrompt | string | - | - | - | - |
| report.applicationReports.incidents | array<object> | - | - | - | - |
| report.applicationReports.incidents.id | string | - | - | - | - |
| report.applicationReports.incidents.title | string | - | - | - | - |
| report.applicationReports.incidents.customPrompt | string | - | - | - | - |
| report.applicationReports.incidents.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| report.applicationReports.incidents.clusterId | string | - | - | - | - |
| report.applicationReports.incidents.applicationName | string | - | - | - | - |
| report.applicationReports.incidents.category | string | - | - | - | - |
| report.applicationReports.incidents.summary | string | - | - | - | - |
| report.applicationReports.incidents.recommendation | string | - | - | - | - |
| report.applicationReports.incidents.urgency | string | - | allowed (`"LOW"`, `"MEDIUM"`, `"HIGH"`) | - | - |
| report.applicationReports.incidents.sources | array<object> | - | - | - | - |
| report.applicationReports.incidents.sources.timestamp | integer | - | - | format (`int64`) | - |
| report.applicationReports.incidents.sources.podName | string | - | - | - | - |
| report.applicationReports.incidents.sources.containerName | string | - | - | - | - |
| report.applicationReports.incidents.sources.image | string | - | - | - | - |
| report.applicationReports.incidents.sources.content | string | - | - | - | - |
| report.totalApplicationEntries | integer | - | - | - | - |
| report.totalNodeEntries | integer | - | - | - | - |
| report.analyzedApplications | integer | - | - | - | - |
| report.analyzedNodes | integer | - | - | - | - |
| report.urgency | string | - | allowed (`"LOW"`, `"MEDIUM"`, `"HIGH"`) | - | - |
| report.scheduledApplicationInsights | object | - | - | - | **additional properties are allowed** |
| report.scheduledApplicationInsights.scheduledJobIds | array<string> | - | - | - | - |
| report.scheduledApplicationInsights.scheduledJobIds (single item) | string | - | - | - | - |
| report.scheduledApplicationInsights.sinceMs | integer | - | - | format (`int64`) | - |
| report.scheduledApplicationInsights.toMs | integer | - | - | format (`int64`) | - |
| report.scheduledApplicationInsights.clusterId | string | - | - | - | - |
| report.scheduledApplicationInsights.applicationConfiguration | array<object> | - | - | - | - |
| report.scheduledApplicationInsights.applicationConfiguration.applicationName | string | - | - | - | - |
| report.scheduledApplicationInsights.applicationConfiguration.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| report.scheduledApplicationInsights.applicationConfiguration.customPrompt | string | - | - | - | - |
| report.scheduledNodeInsights | object | - | - | - | **additional properties are allowed** |
| report.scheduledNodeInsights.scheduledJobIds | array<string> | - | - | - | - |
| report.scheduledNodeInsights.scheduledJobIds (single item) | string | - | - | - | - |
| report.scheduledNodeInsights.sinceMs | integer | - | - | format (`int64`) | - |
| report.scheduledNodeInsights.toMs | integer | - | - | format (`int64`) | - |
| report.scheduledNodeInsights.clusterId | string | - | - | - | - |
| report.scheduledNodeInsights.nodeConfiguration | array<object> | - | - | - | - |
| report.scheduledNodeInsights.nodeConfiguration.nodeName | string | - | - | - | - |
| report.scheduledNodeInsights.nodeConfiguration.accuracy | string | - | allowed (`"HIGH"`, `"MEDIUM"`, `"LOW"`) | - | - |
| report.scheduledNodeInsights.nodeConfiguration.customPrompt | string | - | - | - | - |
| timestampMs | integer | - | - | format (`int64`) | - |

> Examples of payload _(generated)_

```json
{
  "correlationId": "string",
  "report": {
    "id": "string",
    "status": "failed_to_generate",
    "clusterId": "string",
    "sinceMs": 0,
    "toMs": 0,
    "requestedAtNs": 0,
    "scheduledGenerationAtMs": 0,
    "title": "string",
    "nodeReports": [
      {
        "node": "string",
        "accuracy": "HIGH",
        "customPrompt": "string",
        "incidents": [
          {
            "id": "string",
            "title": "string",
            "category": "string",
            "customPrompt": "string",
            "accuracy": "HIGH",
            "clusterId": "string",
            "nodeName": "string",
            "summary": "string",
            "recommendation": "string",
            "urgency": "LOW",
            "sources": [
              {
                "timestamp": 0,
                "content": "string",
                "filename": "string"
              }
            ]
          }
        ]
      }
    ],
    "applicationReports": [
      {
        "applicationName": "string",
        "accuracy": "HIGH",
        "customPrompt": "string",
        "incidents": [
          {
            "id": "string",
            "title": "string",
            "customPrompt": "string",
            "accuracy": "HIGH",
            "clusterId": "string",
            "applicationName": "string",
            "category": "string",
            "summary": "string",
            "recommendation": "string",
            "urgency": "LOW",
            "sources": [
              {
                "timestamp": 0,
                "podName": "string",
                "containerName": "string",
                "image": "string",
                "content": "string"
              }
            ]
          }
        ]
      }
    ],
    "totalApplicationEntries": 0,
    "totalNodeEntries": 0,
    "analyzedApplications": 0,
    "analyzedNodes": 0,
    "urgency": "LOW",
    "scheduledApplicationInsights": {
      "scheduledJobIds": [
        "string"
      ],
      "sinceMs": 0,
      "toMs": 0,
      "clusterId": "string",
      "applicationConfiguration": [
        {
          "applicationName": "string",
          "accuracy": "HIGH",
          "customPrompt": "string"
        }
      ]
    },
    "scheduledNodeInsights": {
      "scheduledJobIds": [
        "string"
      ],
      "sinceMs": 0,
      "toMs": 0,
      "clusterId": "string",
      "nodeConfiguration": [
        {
          "nodeName": "string",
          "accuracy": "HIGH",
          "customPrompt": "string"
        }
      ]
    }
  },
  "timestampMs": 0
}
```




### **ReportRequestFailed** Channel

Topic where the Reports Service publishes failed report requests.

#### `publish` Operation

*Reports Service publishes failed report requests*

##### Message `ReportRequestFailedMessage`

Message sent when a report request fails.

###### Payload

| Name | Type | Description | Value | Constraints | Notes |
|---|---|---|---|---|---|
| (root) | object | - | - | - | **additional properties are allowed** |
| correlationId | string | - | - | - | - |
| errorType | string | - | allowed (`"VALIDATION_ERROR"`, `"TIMEOUT"`, `"INTERNAL_ERROR"`) | - | - |
| errorMessage | string | - | - | - | - |
| timestampMs | integer | - | - | format (`int64`) | - |

> Examples of payload _(generated)_

```json
{
  "correlationId": "string",
  "errorType": "VALIDATION_ERROR",
  "errorMessage": "string",
  "timestampMs": 0
}
```




