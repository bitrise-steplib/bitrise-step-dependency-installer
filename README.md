# Dependency Installer

[![Step changelog](https://shields.io/github/v/release/bitrise-io/bitrise-step-dependency-installer?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-io/bitrise-step-dependency-installer/releases)

Installs dependencies from tool version files in your repo.

<details>
<summary>Description</summary>

This wrapper for the bitrise CLI tool installs dependencies
defined in the tool version file(s) in your repo.

To learn more about defining tool versions as code, read our
[documentation](https://bitrise.io/stacks/tips/tool-versions).

Tool version files supported by this step include:
* .tool-versions
* .ruby-version
* .node-version
* .python-version
* .java-version
* .go-version
* .terraform-version
* .kubectl-version
* bitrise.yml (with a `tools:` section)

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `tool_version_file` | The tool version file contains the version(s) of the tool(s) for the step to install. See the step description for the list of supported tool version files.  | required | `.tool-versions` |
| `verbose` | If enabled, the step will output additional debug information during execution. Set to "true" for verbose logging.  |  | `false` |
| `bitrise_workflow_id` | This input is only used when passing a bitrise.yml file as a tool version file and rather than using the global `tools:` block, you want to use the `tools:` block defined under specific workflows. The step accepts a workflow ID and installs the tools defined under that workflow's `tools:` block.  |  |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-io/bitrise-step-dependency-installer/pulls) and [issues](https://github.com/bitrise-io/bitrise-step-dependency-installer/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
