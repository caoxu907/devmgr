job "devmgr" {
  datacenters = ["dc1"]
  type = "service"
  reschedule {
     delay          = "30s"
     delay_function = "constant"
     unlimited      = true
  }
  update {
      max_parallel = 1
      min_healthy_time = "10s"
      healthy_deadline = "3m"
      auto_revert = false
      auto_promote = true
      canary = 1
      health_check = "task_states"
  }
  group "devmgr" {
    count = 2
    spread{
        attribute = "${node.unique.id}"
        weight = 100
      }
    task "devmgr" {
      driver = "docker"
      template {
          data = <<EOF
          {{- with nomadVar "nomad/jobs/devmgr/devmgr/devmgr" -}}
              {{- range $k, $v := . }}
                  {{ $k }}={{ $v }}
              {{- end }}
          {{- end }}
          EOF
          destination = "secrets/env"
          env = true
      }
      config {
        network_mode = "cloud_net"
        image = "10.17.196.52:12888/devmgr:${IMAGE_VERSION}"
        extra_hosts = [
          "kafka-node1.cluster.local:10.17.196.182",
          "kafka-node2.cluster.local:10.17.196.183"
        ]
      }
      restart {
          interval         = "2m"
          attempts         = 2
          delay            = "15s"
          mode             = "fail"
          render_templates = false
      }

      resources {
        cpu    = 500 # 0.5 CPU
        memory = 512 # 256MB memory
      }
    }
  }
}