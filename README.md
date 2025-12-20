# terraform-ui

> `terraform-ui` is under development. A first release candidate will be released "soon".

Graphical user interface fo Terraform without, without dealing directly with Terraform or the underlying infrastructure.

More details can be found in the [idea description](./IDEA.md).

## Components

terraform-ui consists of two components.

### terraform-ui-core

**terraform-ui-core** is the backbone of the whole stack.  
Triggered via REST API, terraform-ui executes terraform commands to provision infrastructure.

Source code and detailed documentation is located under [terraform-ui-core](./terraform-ui-core)

### terraform-ui-web

Not developed yet.

# Development

Please make sure the read the following guidelines before contributing to terraform-ui.

## Commit messages

Since this is a "mono-repository" for more than one component of terraform-ui, consistent commit names are important. This allows commits to be grouped and assigned to the corresponding component.  

Commit naming:
- `gh-actions - ...`: Prefix for GitHub Action changes
- `terraform-ui-core - ..`: Prefix for backend changes
- `terraform-ui-web - ...`: Prefix for frontend changes

Changes that are not part of actions, backend or frontend doenst need a consistent prefix.

It is important to separate changes to the backend and frontend into separate commits.

# LICENSE

MIT © Tobias Bauriedel