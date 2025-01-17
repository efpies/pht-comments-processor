terraform {
  required_version = ">= 1.10.4"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "aws_region" {
  default = "eu-north-1"
}

variable "aws_profile" {
  default = "pht-deployer"
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile
}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "AWSLambdaTrustPolicy" {
  statement {
    actions = ["sts:AssumeRole"]
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "lambda_access_policy" {
  statement {
    actions = [
      "ssm:GetParameters",
      "ssm:GetParametersByPath",
      "ssm:PutParameter",
    ]
    effect = "Allow"
    resources = [
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.id}:parameter/pht-comments-processor",
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.id}:parameter/pht-comments-processor/*",
    ]
  }
}

resource "aws_iam_policy" "lambda_access_policy" {
  name   = "PhtCommentsProcessorPolicy"
  policy = data.aws_iam_policy_document.lambda_access_policy.json
}

resource "aws_iam_role" "lambda_exec_role" {
  name               = "pht-comments-processor-exec-role"
  assume_role_policy = data.aws_iam_policy_document.AWSLambdaTrustPolicy.json
}

resource "aws_iam_policy_attachment" "lambda_exec_policy" {
  name       = "${aws_iam_role.lambda_exec_role.name}-AWSLambdaBasicExecutionRole"
  roles = [aws_iam_role.lambda_exec_role.name]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy_attachment" "lambda_access_policy" {
  name       = "${aws_iam_role.lambda_exec_role.name}-${basename(aws_iam_policy.lambda_access_policy.arn)}"
  roles = [aws_iam_role.lambda_exec_role.name]
  policy_arn = aws_iam_policy.lambda_access_policy.arn
}

resource "aws_ecr_repository" "repo" {
  name         = "pht-comments-processor"
  force_delete = true
}

locals {
  dockerfile_hash = filesha256("../Dockerfile")
  src_hash = sha256(join("", [for f in fileset(".", "../src/**") : filesha256(f) if basename(f) != ".DS_Store" && !strcontains(basename(f), "_test")]))
  dir_hash = sha256(join("", [local.dockerfile_hash, local.src_hash]))
}

resource "null_resource" "docker_build_and_push" {
  triggers = {
    dir_hash = local.dir_hash
  }

  provisioner "local-exec" {
    command = <<EOT
      docker build --platform linux/amd64 -t ${aws_ecr_repository.repo.name}:latest -f ../Dockerfile ..
      docker tag ${aws_ecr_repository.repo.name}:latest ${aws_ecr_repository.repo.repository_url}:latest

      aws ecr get-login-password --region ${var.aws_region} --profile ${var.aws_profile} | docker login --username AWS --password-stdin ${aws_ecr_repository.repo.repository_url}
      docker push ${aws_ecr_repository.repo.repository_url}:latest
    EOT
  }

  depends_on = [aws_ecr_repository.repo]
}

resource "null_resource" "update_lambda_image" {
  triggers = {
    image_uri = "${aws_ecr_repository.repo.repository_url}:latest"
  }

  provisioner "local-exec" {
    command = <<EOT
      aws lambda update-function-code \
        --region ${var.aws_region} \
        --profile ${var.aws_profile} \
        --function-name ${aws_lambda_function.lambda_function.function_name} \
        --image-uri ${aws_ecr_repository.repo.repository_url}:latest
    EOT
  }

  depends_on = [null_resource.docker_build_and_push]

  lifecycle {
    replace_triggered_by = [
      null_resource.docker_build_and_push,
    ]
  }
}

resource "aws_lambda_function" "lambda_function" {
  function_name = "pht-comments-processor"
  role          = aws_iam_role.lambda_exec_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.repo.repository_url}:latest"
  timeout       = 30
  memory_size   = 128

  lifecycle {
    ignore_changes = [
      image_uri,
    ]
  }
}

resource "aws_lambda_function_url" "lambda_function_url" {
  function_name      = aws_lambda_function.lambda_function.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "FunctionURLAllowPublicAccess" {
  statement_id           = "FunctionURLAllowPublicAccess"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.lambda_function.function_name
  principal              = "*"
  function_url_auth_type = "NONE"
}