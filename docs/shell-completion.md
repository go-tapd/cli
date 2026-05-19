# Shell Completion

Use `tapd completion` to generate shell completion scripts for the installed
`tapd` binary.

## Supported shells

```bash
tapd completion bash
tapd completion zsh
tapd completion fish
tapd completion powershell
```

Run shell-specific help for the exact generated script behavior:

```bash
tapd completion zsh --help
tapd completion bash --help
tapd completion fish --help
tapd completion powershell --help
```

## zsh

Load completion in the current shell session:

```zsh
source <(tapd completion zsh)
```

Install completion for future zsh sessions on macOS with Homebrew:

```zsh
tapd completion zsh > "$(brew --prefix)/share/zsh/site-functions/_tapd"
```

If zsh completion is not enabled yet, enable it once:

```zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Start a new shell after installing the completion file.

## bash

Load completion in the current shell session:

```bash
source <(tapd completion bash)
```

Install completion for future bash sessions on macOS with Homebrew:

```bash
tapd completion bash > "$(brew --prefix)/etc/bash_completion.d/tapd"
```

The generated bash completion script depends on the `bash-completion` package.
Install it first if completion does not load in bash.

## fish

Load completion in the current shell session:

```fish
tapd completion fish | source
```

Install completion for future fish sessions:

```fish
tapd completion fish > ~/.config/fish/completions/tapd.fish
```

Start a new shell after installing the completion file.

## PowerShell

Load completion in the current PowerShell session:

```powershell
tapd completion powershell | Out-String | Invoke-Expression
```

To load completion for future sessions, add the same command to your PowerShell
profile.

## Updating completions after adding commands

Completion scripts are generated from the current CLI command tree. After adding
or renaming commands, regenerate the installed completion file with the same
command used during installation.

For example, after adding a new zsh command:

```zsh
tapd completion zsh > "$(brew --prefix)/share/zsh/site-functions/_tapd"
```
