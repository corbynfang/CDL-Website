resource "random_password" "db" {
  length  = 32
  special = false # avoid chars like @ # that break postgres connection strings
}
