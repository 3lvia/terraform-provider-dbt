variable "account_id" {
  description = "The account id for DBT cloud"
  type        = number
  nullable    = false
  sensitive   = false
}

variable "service_token" {
  description = "The service token for api-requests to DBT"
  type        = string
  nullable    = false
  sensitive   = true
}
