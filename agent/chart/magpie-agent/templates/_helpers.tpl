{{- define "broker.connection" -}}
{{- with .Values.agent.remoteWrite.logs -}}
- "--remoteWriteBrokerUrl"
- {{ .url | quote }}
- "--remoteWriteNodeTopic"
- "nodes"
- "--remoteWriteApplicationTopic"
- "pods"
- "--remoteWriteBatchSize"
- "0"
- "--remoteWriteBrokerUsername"
- {{ .username | quote }}
- "--remoteWriteBrokerPassword"
- {{ .password | quote }}
{{- end -}}
{{- end -}}

{{- define "metadata.scrapeIntervals" -}}
- "--logScrapeIntervalSeconds"
- "60"
- "--metadataScrapeIntervalSeconds"
- "60"
{{- end -}}

{{- define "redis.connection" -}}
{{- with .Values.agent.redis -}}
{{- if .enabled -}}
- "--redisUrl"
- "magpie-agent-redis.{{ $.Release.Namespace }}.svc.cluster.local:6379"
- "--redisDatabase"
- "0"
- "--redisPassword"
- {{ .password | quote }}
{{- else -}}
- "--redisUrl"
- {{ .url | quote }}"
- "--redisDatabase"
- {{ .database | quote }}
- "--redisPassword"
- {{ .password | quote }}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "agent.pod.config" -}}
- "--clusterFriendlyName"
- {{ .Values.agent.friendlyName | quote }}
- "--scrape"
- "pods"
- "--podRemoteWriteMetadataUrl"
- {{ .Values.agent.remoteWrite.metadata.url | quote }}
{{- end -}}

{{- define "agent.pod.excludedNamespaces" }}
{{- range $namespace := .Values.agent.application.excludedNamespaces }}
- "--excludedNamespace"
- {{ $namespace | quote }}
{{- end -}}
{{- end -}}

{{- define "agent.node.config" -}}
- "--clusterFriendlyName"
- {{ .Values.agent.friendlyName | quote }}
- "--scrape"
- "nodes"
- "--nodeRemoteWriteMetadataUrl"
- {{ .Values.agent.remoteWrite.metadata.url | quote }}
{{- end -}}

{{- define "agent.node.watchedFiles" }}
{{- range $file := .Values.agent.node.files }}
- "--file"
- {{ $file | quote }}
{{- end -}}
{{- end -}}

{{- define "agent.node.volumeMounts" }}
{{- range $file := .Values.agent.node.files }}
- name: {{ splitList "/" $file | last | replace "." "-" }}
  mountPath: {{ $file }}
  readOnly: true
{{- end -}}
{{- end -}}

{{- define "agent.node.volumes" }}
{{- range $file := .Values.agent.node.files }}
- name: {{ splitList "/" $file | last | replace "." "-" }}
  hostPath:
    path: {{ $file }}
{{- end -}}
{{- end -}}