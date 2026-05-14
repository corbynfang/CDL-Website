terraform {
  backend "s3" {
    bucket  = "cdl-tf-state-233060639311-us-east-1"
    key     = "cdl-website/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}
