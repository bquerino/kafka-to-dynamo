###############################################################################
# Provider
###############################################################################
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "aws" {
  # Usando LocalStack
  region                  = "us-east-1"
  access_key              = "test"
  secret_key              = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_requesting_account_id  = true

  # Endpoint para o DynamoDB no LocalStack
  endpoints {
    dynamodb = "http://localhost:4566"
  }
}

###############################################################################
# Resources
###############################################################################
resource "aws_dynamodb_table" "pagamentos" {
  name         = "Payments"
  billing_mode = "PAY_PER_REQUEST"

  # Definindo as chaves explicitamente:
  hash_key  = "CustomerID#PaymentID"
  range_key = "PaymentEventDate"

  # Declarando os atributos
  attribute {
    name = "CustomerID#PaymentID"
    type = "S"
  }

  attribute {
    name = "PaymentEventDate"
    type = "S"
  }

  tags = {
    Environment = "dev"
    Application = "GoLangKafkaConsumer"
  }
}
