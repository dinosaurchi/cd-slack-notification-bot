#!/bin/sh

toml_path=$(which toml || which toml-tool)

if [ "$toml_path" = "" ]; then
  echo "Please install 'toml-cli' tool with one of the following commands"
  echo "
    npm install -g toml-tool (recommended)

    pip install toml-cli

    pip3 install toml-cli
  "
  exit 1
fi

toml_file_path=$1
value_path=$2

if [ "$toml_file_path" = "" ]; then
  echo "Missing toml_file_path as the 1st argument"
  exit 1
elif [ "$value_path" = "" ]; then
  echo "Missing value_path as the 2nd argument (such as MyModule.field1)"
  exit 1
fi

toml_cli_path=$(which toml)
toml_tool_path=$(which toml-tool)

if [ "$toml_tool_path" != "" ]; then
  # Prioritize using toml-tool as it is more stable than the toml-cli
  toml-tool get $toml_file_path $value_path
elif [ "$toml_cli_path" != "" ]; then
  # https://pypi.org/project/toml-cli/
  toml get --toml-path $toml_file_path $value_path
fi
