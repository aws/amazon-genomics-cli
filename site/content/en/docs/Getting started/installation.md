---
title: "Installation"
date: 2021-09-07T13:42:31-04:00
draft: false
weight: 20
---

## Download and install Amazon Genomics CLI

Download the Amazon Genomics CLI zip, unzip its contents, and run the `install.sh` script:

To download a specific release, see [releases page](https://github.com/aws/amazon-genomics-cli/releases) of our Github repo.

```
curl -OLs https://github.com/aws/amazon-genomics-cli/releases/latest/download/amazon-genomics-cli.zip
unzip amazon-genomics-cli.zip -d agc
./agc/install.sh
```

This will place the `agc` command in `$HOME/bin`.

The Amazon Genomics CLI is a statically compiled Go binary. It should run in your environment natively without any additional setup. Test the CLI with:

```
$ agc --help

🧬 Launch and manage genomics workloads on AWS.

Commands
  Getting Started 🌱
    account     Commands for AWS account setup.
                Install or remove Amazon Genomics CLI from your account.

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

If this doesn't work immediately, try:

* start a new terminal shell
* modifying your `$HOME/.bashrc` (or equivalent file) appending the following line and restarting your shell:

```
export PATH=$HOME/bin:$PATH
```

Verify that you have the latest version of Amazon Genomics CLI with:

```
agc --version
```

If you do not, you may need to uninstall any previous versions of Amazon Genomics CLI and reinstall the latest.

## Command Completion

Amazon Genomics CLI can generate shell completion scripts that enable 'Tab' completion of commands. 
Command completion is optional and not required to use Amazon Genomics CLI. To generate a completion script you can use:

```shell
 agc generate <shell>
``` 

where "shell" is one of:

### Bash

```shell
source <(agc completion bash)
```

To load completions for each session, execute once:
#### Linux:
```shell
agc completion bash > /etc/bash_completion.d/agc
```

#### macOS:

If you haven't already installed `bash-completion`, execute the following once

```shell
brew install bash-completion
```

and then, add the following line to your ~/.bash_profile:

```shell
[[ -r "/usr/local/etc/profile.d/bash_completion.sh" ]] && . "/usr/local/etc/profile.d/bash_completion.sh"
```

Once bash completion is installed

```shell
agc completion bash > /usr/local/etc/bash_completion.d/agc
```



### Zsh:

If shell completion is not already enabled in your environment, you will need to enable it.  You can execute the following once:

```shell
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for each session, execute once:

```shell
agc completion zsh > "${fpath[1]}/_agc"
```

You will need to start a new shell for this setup to take effect.

### fish:

```shell
agc completion fish | source
```

To load completions for each session, execute once:
```shell
agc completion fish > ~/.config/fish/completions/agc.fish
```
PowerShell:

```shell
agc completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, run:

```shell
agc completion powershell > agc.ps1
```

and source this file from your PowerShell profile.
