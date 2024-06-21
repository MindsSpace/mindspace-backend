provider "google" {
  project = "mindspace-capstone-project"
}

resource "google_cloud_run_v2_service" "default" {
  name     = "mindspace-backend"
  location = "asia-southeast2"
  client   = "terraform"

  template {
    containers {
      image = "asia-southeast2-docker.pkg.dev/mindspace-capstone-project/mindspace/mindspace-backend:latest"

      env {
        name  = "DB_HOST"
        value = var.DB_HOST
      }

      env {
        name  = "DB_USER"
        value = var.DB_USER
      }

      env {
        name  = "DB_PASS"
        value = var.DB_PASS
      }

      env {
        name  = "DB_NAME"
        value = var.DB_NAME
      }

      env {
        name  = "DB_PORT"
        value = var.DB_PORT
      }

      env {
        name  = "JWT_SECRET"
        value = var.JWT_SECRET
      }

      env {
        name  = "MLAPI_URL"
        value = var.MLAPI_URL
      }

      liveness_probe {
        failure_threshold     = 5
        initial_delay_seconds = 10
        timeout_seconds       = 3
        period_seconds        = 3

        http_get {
          path = "/health"

          http_headers {
            name  = "Access-Control-Allow-Origin"
            value = "*"
          }
        }
      }
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
