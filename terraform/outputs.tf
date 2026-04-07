output "service_name" {
  value = kubernetes_service.query_service.metadata[0].name
}

output "deployment_name" {
  value = kubernetes_deployment.query_service.metadata[0].name
}
