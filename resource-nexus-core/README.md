# resource-nexus-core

`resource-nexus-core` is the backbone of the whole stack.  
Triggered via REST API, `resource-nexus-core` executes terraform commands to provision infrastructure.

# How it works

`resource-nexus-core` simply works as a "wrapper" for [terraform](https://github.com/hashicorp/terraform). Using the terraform CLI, `resource-nexus-core` builds necessary configuration files and applies them.  
The main magic and logic comes from terraform itself. `resource-nexus-core` takes your configuration as JSON, stores them in a database and builds terraform files to execute.

# Documentation

Documentation can be found under the [docs](./docs) section.

# Development

For local testing, run `make`.  
A [Makefile](./Makefile) is provided to simplify the development process.

# Disclaimer

`resource-nexus` is an OSS project that uses and builds on Terraform. It is not affiliated with HashiCorp or Terraform.

# LICENSE

MIT © Tobias Bauriedel