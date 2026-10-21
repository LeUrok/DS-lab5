{{- define "app.name" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "app.fullname" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "app.labels" -}}
app.kubernetes.io/name: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "app.selectorLabels" -}}
app.kubernetes.io/name: {{ .Release.Name }}
{{- end -}}