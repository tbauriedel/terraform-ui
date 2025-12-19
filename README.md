# Terraform UI

Graphical user interface fo Terraform without, without dealing directly with Terraform or the underlying infrastructure.

More details can be found in the [idea description](./IDEA.md).

# Development

Please make sure the read the following guidelines before contributing to terraform-ui.

## Commit messages

Since this is a "mono-repository" for more than one component of terraform-ui, consistent commit names are important. This allows commits to be grouped and assigned to the corresponding component.  

Commit naming:
- `gh-actions - ...`: Prefix for GitHub Action changes
- `terraform-ui - ..`: Prefix for backend changes
- `terraform-ui-web - ...`: Prefix for frontend changes

Changes that are not part of actions, backend or frontend doenst need a consistent prefix.

It is important to separate changes to the backend and frontend into separate commits.

# LICENSE

MIT © Tobias Bauriedel