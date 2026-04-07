terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "kind-query-service"
}

resource "kubernetes_deployment" "query_service" {
  metadata {
    name = "query-service"
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "query-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "query-service"
        }
      }

      spec {
        container {
          name              = "query-service"
          image             = "query-service:${var.image_tag}"
          image_pull_policy = "Never"

          port {
            container_port = 8080
          }

          env {
            name  = "LLM_BASE_URL"
              value = "http://192.168.6.78:8081"
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "256Mi"
            }
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 5
            period_seconds        = 10
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "query_service" {
  metadata {
    name = "query-service"
  }

  spec {
    selector = {
      app = "query-service"
    }

    port {
      port        = 80
      target_port = 8080
    }

    type = "ClusterIP"
  }
}
