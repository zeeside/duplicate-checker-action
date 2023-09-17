terraform {
  required_version = "~> 1.1"

  backend "s3" {
    region         = "us-east-1"
    dynamodb_table = "terraform-state-locks"
    bucket         = "my-tf-bucket"
    key            = "context/address-service/service.tfstate"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}