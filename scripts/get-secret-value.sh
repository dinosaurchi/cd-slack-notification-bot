#!/bin/sh

jq_path=$(which jq)

if [ "$jq_path" = "" ]; then
  echo "Please install 'jq' CLI tool"
  echo "
    brew install jq (macOS)

    sudo apt-get install jq (Ubuntu)

    pkg install jq (Fedora)
  "
  exit 1
fi

secret_id=$1
aws_region=$2
secret_path=$3

if [ "$secret_id" = "" ]; then
  echo "Missing secret_id as the 1st argument"
  exit 1
elif [ "$aws_region" = "" ]; then
  echo "Missing aws_region as the 2nd argument (such as us-west-2, ap-southeast-1,...)"
  exit 1
elif [ "$secret_path" = "" ]; then
  echo "Missing secret_path as the 3rd argument (such as .MyModule.field1)"
  exit 1
fi

aws secretsmanager get-secret-value --secret-id $secret_id --region $aws_region --query SecretString | jq --raw-output | jq -r $secret_path
