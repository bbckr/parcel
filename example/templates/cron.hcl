job "cron-{{ .Values.name }}" {
  type        = "batch"
  datacenters = {{ .Values.datacenters | toJSON }}

  periodic {
    cron             = "{{ .Values.periodic.schedule }}"
    time_zone        = "{{ .Values.periodic.time_zone }}"
    prohibit_overlap = {{ .Values.periodic.prohibit_overlap }}
  }

  meta {
    description    = "{{ .Values.description }}"
    parcel.version = "{{ .Meta.version }}"
    parcel.name    = "{{ .Meta.name }}"
    parcel.owner   = "{{ .Meta.owner }}"
    {{- range $k, $v := .Values.meta }}
    {{ $k }} = {{ $v }}
    {{- end }}
  }

  group "cron" {
    count = 1

    restart {
      mode = "{{ .Values.restart_mode }}"
    }

    meta {
      maxscale.enabled = {{ .Values.maxscale }}
    }

    task "cron" {
      driver = "docker"

      resources {
        cpu    = {{ .Values.resources.cpu }}
        memory = {{ .Values.resources.memory }}

        network {
          mbits = {{ .Values.resources.network.mbits }}
        }
      }

    {{- if .Values.env }}
    
      env {
        {{- range $k, $v := .Values.env }}
        {{ $k }} = {{ $v | quote }}
        {{- end }}
      }
    {{- end }}

      config {
        image   = "{{ .Values.image.name }}:{{ .Values.image.tag }}"
        command = "{{ .Values.command }}"
        args    = {{ .Values.args | toJSON }}
      }
    }
  }
}
