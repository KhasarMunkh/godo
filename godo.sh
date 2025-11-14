#!/bin/bash
# godo shell integration
# Source this file in your ~/.bashrc or ~/.zshrc to enable auto-display of todos

# Function to display todos when changing directory
_godo_auto_display() {
    # Check if godo is installed and executable
    if command -v godo &> /dev/null; then
        godo show 2>/dev/null
    fi
}

# For Bash
if [ -n "$BASH_VERSION" ]; then
    # Hook into cd command
    _godo_cd() {
        builtin cd "$@" && _godo_auto_display
    }
    alias cd='_godo_cd'

    # Also run on new shell
    _godo_auto_display
fi

# For Zsh
if [ -n "$ZSH_VERSION" ]; then
    # Use chpwd hook for zsh
    autoload -U add-zsh-hook
    add-zsh-hook chpwd _godo_auto_display

    # Also run on new shell
    _godo_auto_display
fi
