job "cron-{{ .Values.name }}" {
  type        = "batch"
  datacenters = ["us-east-1"]

  periodic {
    cron             = {{ .Values.schedule }}
    time_zone        = {{ .Values.periodic.time_zone }}
    prohibit_overlap = {{ .Values.periodic.prohibit_overlap }}
  }

  meta {
    description = "{{ .Values.description }}"
  }

  group "cron" {
    count = 1

    restart {
      mode = "{{ Values.restart_mode }}"
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
          mbits = {{ .Values.network.mbits }}
        }
      }

      env {
        PWD        = "/app"
        LOG_FORMAT = "json"
      }

      config {
        image   = "{{ .Values.image.name }}:{{ .Values.image.tag }}"
        command = "{{ .Values.command }}"
        args    = {{ .Values.args }}
      }
    }
  }
}
