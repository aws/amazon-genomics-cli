---
title: "Installation"
date: 2021-09-07T13:42:31-04:00
draft: false
weight: 20
---

## Download and install Amazon Genomics CLI

Download the Amazon Genomics CLI zip, unzip its contents, and run the `install.sh` script:

```
aws s3 cp s3://healthai-public-assets-us-east-1/amazon-genomics-cli/0.9.0/amazon-genomics-cli.zip .
unzip amazon-genomics-cli.zip -d agc
./agc/install.sh
```

This will place the `agc` command in `~/bin` and attempt to add the command to your `$PATH`

The Amazon Genomics CLI is a statically compiled Go binary. It should run in your environment natively without any additional setup. Test the CLI with:

```
$ agc --help
👩🧬 Launch and manage genomics workloads on AWS.

Commands
  Getting Started 🌱
    account     Commands for AWS account setup.
                Install or remove AGC from your account.

  Contexts
    context     Commands for contexts.
                Contexts specify workflow engines and computational fleets to use when running a workflow.

  Logs
    logs        Commands for various logs.

  Projects
    project     Commands to interact with projects.

  Workflows
    workflow    Commands for workflows.
                Workflows are potentially-dynamic graphs of computational tasks to execute.

  Settings ⚙️
    version     Print the version number.
    configure   Commands for configuration.
                Configuration is stored per user.

Flags
  -h, --help      help for agc
  -v, --verbose   Display verbose diagnostic information.
      --version   version for agc
Examples
  Displays the help menu for the specified sub-command.
  `$ agc account --help`
```

If this doesn’t work immediately, try:

* start a new terminal shell
* modifying your `~/.bashrc` (or equivalent file) appending the following line and restarting your shell:

```
export PATH=~/bin:$PATH
```

Verify that you have the latest version of Amazon Genomics CLI with:

```
agc --version
```

If you do not, you may need to uninstall any previous versions of Amazon Genomics CLI and reinstall the latest.
