# godo - Directory-Level Todo Manager

A simple, directory-specific todo manager for your terminal. When you `cd` into a project directory, your todos automatically appear!

## Features

- **Directory-level todos**: Each directory has its own `.godo.json` file
- **Auto-display**: Todos automatically show when you enter a directory
- **Simple commands**: Add, list, complete, and remove todos easily
- **Hide completed**: Completed todos are hidden by default, but viewable with `--all`
- **Colorful output**: Easy-to-read colored terminal output
- **Lightweight**: Single binary with no dependencies

## Quick Reference

| Command | Alias | Description |
|---------|-------|-------------|
| `godo add <text>` | `a` | Add a new todo |
| `godo list` | `l` | List active todos |
| `godo list --all` | `l --all` | List all todos including completed |
| `godo done <id>` | `d` | Mark todo as complete |
| `godo remove <id>` | `rm` | Remove an active todo |
| `godo clean` | `c` | Remove all completed todos |
| `godo clean <id>` | `c` | Remove a specific completed todo |
| `godo show` | - | Show active todos (auto-display) |
| `godo help` | - | Show help message |

## Installation

### Build from source

```bash
# Clone the repository
git clone https://github.com/KhasarMunkh/godo
cd godo

# Build the binary
go build -o godo

# Install to your PATH (recommended)
sudo mv godo /usr/local/bin/godo

# Or create a symlink instead
sudo ln -s $(pwd)/godo /usr/local/bin/godo
```

### Set up shell integration

To enable automatic todo display when changing directories, add the following to your shell configuration file:

(add to `~/.bashrc` or `~/.zshrc`):
```bash
source /path/to/godo/godo.sh
```

Replace `/path/to/godo` with the actual path where you cloned the repository.

Then reload your shell:
```bash
source ~/.bashrc  # or ~/.zshrc for zsh
```

## Usage

### Add a todo
```bash
godo add "Implement user authentication"
godo a "Write unit tests"  # Short form
```

### List todos
```bash
# List active (incomplete) todos
godo list
godo l  

# List all todos including completed ones
godo list --all
godo l --all  
```

### Mark a todo as complete
```bash
godo done 1
godo d 1  
```

### Remove a todo
```bash
godo remove 2
godo rm 2  
```

### Clean completed todos
```bash
# Remove all completed todos
godo clean
godo c  

# Remove a specific completed todo (by its position in the completed list)
godo clean 1
godo c 1  
```

### Show todos (auto-display)
```bash
godo show
```
This command is used internally by the shell integration to display todos when you `cd` into a directory.

### Get help
```bash
godo help
```

## How It Works

- Each directory stores its todos in a `.godo.json` file
- When you run `godo` commands, it operates on the `.godo.json` in your current directory
- The shell integration automatically runs `godo show` when you change directories
- Completed todos are hidden from `godo list` and `godo show` by default, but can be viewed with `godo list --all`
- **Position-based IDs**: Todo IDs are based on position in the active list and automatically renumber when todos are completed or removed. This means you'll always have simple, sequential IDs like 1, 2, 3 instead of having to remember large numbers

## Example Workflow

```bash
# Navigate to your project
cd ~/projects/myapp

# Add some todos 
godo a "Fix login bug"
godo a "Add email validation"
godo a "Update README"

# List todos 
godo l
# Output:
# Active Todos:
#   [1] Fix login bug
#   [2] Add email validation
#   [3] Update README

# Complete a todo
godo d 1
# [Completed] Todo #1: Fix login bug

# List again (completed todo is hidden, IDs automatically renumber!)
godo l
# Output:
# Active Todos:
#   [1] Add email validation
#   [2] Update README

# View all todos including completed
godo l --all
# Output:
# Active Todos:
#   [1] Add email validation
#   [2] Update README
#
# Completed Todos:
#   [1] Fix login bug

# When you cd to another directory and back, todos show automatically
cd ..
cd myapp
# Output:
# Todos:
#   [1] Add email validation
#   [2] Update README

# Clean up completed todos when you're done
godo c
# [Cleaned] Removed 1 completed todo(s)
```

## Tips

- Add `.godo.json` to your global `.gitignore` if you don't want to commit todos to version control
- Or commit `.godo.json` to share project todos with your team

## License

MIT
