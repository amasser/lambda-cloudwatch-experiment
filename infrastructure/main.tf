provider "aws" {
  region = "eu-central-1"
}

resource "aws_lambda_function" "lambda_cloudwatch_experiment" {
  function_name = "lambda-cloudwatch-experiment-dev"
  filename      = "../lambda.zip"
  handler       = "lambda"
  runtime       = "go1.x"
  role          = aws_iam_role.exec.arn
}

resource "aws_iam_role" "exec" {
  name = "lambda-cloudwatch-experiment-exec-dev"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "exec_to_lambda_basic_execution_attach" {
  role       = aws_iam_role.exec.name
  policy_arn = data.aws_iam_policy.lambda_basic_execution.arn
}

data "aws_iam_policy" "lambda_basic_execution" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_cloudwatch_log_group" "lambda_cloudwatch_experiment" {
  name              = "/aws/lambda/${aws_lambda_function.lambda_cloudwatch_experiment.function_name}"
  retention_in_days = 90
}
