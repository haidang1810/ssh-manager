# SSH Connection Manager

![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)

## Introduction

**SSH Connection Manager** is a command-line interface (CLI) tool developed in Go (Golang) designed to help you manage SSH connections efficiently. This tool allows you to store, organize, and quickly connect to remote servers without having to remember complex SSH commands and configurations.

With SSH Connection Manager, you can:

*   Add, list, edit, and remove SSH connections.
*   Encrypt passwords for secure storage.
*   Manage and generate SSH key pairs.
*   Export and import connection configurations.
*   Use shorthand commands for quick connections.

## Installation

SSH Connection Manager is a standalone command-line tool. Once built, you do not need Go installed on your machine to run it.

### 1. Download Executable (Recommended)

You can download the pre-built executable for your operating system from the project's [Releases](https://github.com/haidang1810/ssh-manager/releases) page. Choose the version compatible with your operating system and CPU architecture (e.g., `ssh-manager-v1.0.0-linux-amd64.tar.gz` for Linux 64-bit, or `ssh-manager-v1.0.0-windows-amd64.zip` for Windows 64-bit).

After downloading, follow the instructions below depending on your operating system:

#### a. On Linux and macOS

1.  **Download:** Download the appropriate `.tar.gz` file for your system (e.g., `ssh-manager-<version>-linux-amd64.tar.gz`). You can use `curl` to download it directly:

    ```bash
    curl -LO https://github.com/haidang1810/ssh-manager/releases/download/v1.0.0/ssh-manager-v1.0.0-<os>-<arch>.tar.gz
    ```
2.  **Extract:** Open Terminal and extract the downloaded file, and rename it to `ssh-manager`:

    ```bash
    tar -xzf ssh-manager-<version>-<os>-<arch>.tar.gz -O ssh-manager-<os>-<arch> > ssh-manager
    chmod +x ssh-manager
    ```
3.  **Move to PATH:** Move the executable to a directory included in your system's `PATH` environment variable (e.g., `/usr/local/bin/`). This allows you to run the `ssh-manager` command from anywhere in the Terminal.

    ```bash
    sudo mv ssh-manager /usr/local/bin/
    ```

4.  **Verify Installation:** Open a new Terminal and run the following command to verify.

    ```bash
    ssh-manager --version
    ```

#### b. On Windows

1.  **Download:** Download the appropriate `.zip` file for your system (e.g., `ssh-manager-<version>-windows-amd64.zip`).
2.  **Extract:** Extract the downloaded `.zip` file to a directory where you want to store the tool (e.g., `C:\Program Files\ssh-manager`).
3.  **Add to PATH Environment Variable:** To be able to run `ssh-manager` from anywhere in Command Prompt or PowerShell, you need to add the path to the directory containing the executable to your system's `PATH` environment variable.
    *   Search for "Edit the system environment variables" in the Start Menu.
    *   Click on "Environment Variables...".
    *   In the "System variables" section, find and select the `Path` variable, then click "Edit...".
    *   Click "New" and add the path to the directory where you extracted `ssh-manager` (e.g., `C:\Program Files\ssh-manager`).
    *   Click "OK" on all windows to save changes.
4.  **Verify Installation:** Open a new Command Prompt or PowerShell (it's important to open a new window for the environment variables to take effect) and run the following command to verify.

    ```cmd
    ssh-manager --version
    ```

### 2. Build from Source (For Developers)

If you want to build the tool from source, you need to have Go (version 1.21 or later) installed on your system.

1.  **Clone repository:**

    ```bash
    git clone https://github.com/haidang1810/ssh-manager.git # Replace with actual repo URL
    cd ssh-manager
    ```

2.  **Build project:**

    ```bash
    go mod tidy
    go build -o ssh-manager .
    ```

    The `ssh-manager` executable will be created in the current directory. You can move it to a directory in your `PATH` as instructed above.

## Usage Guide

SSH Connection Manager provides a set of intuitive commands to manage your SSH connections.

### Configuration

The tool will automatically create a configuration file at `$HOME/.ssh-manager/config.yaml` or `$HOME/.config/ssh-manager/config.yaml` when you use it for the first time.

### Main Commands

#### 1. `ssh-manager add` - Add a new connection

Add a new SSH connection to your configuration. If you do not provide required flags, the tool will prompt you for information.

```bash
ssh-manager add <connection_name> --host <host> --user <user> [--port <port>] [--key <key_path>] [--pass <password>]

# Examples:
ssh-manager add my_server --host example.com --user admin --port 2222
ssh-manager add dev_machine --host 192.168.1.100 --user dev --key ~/.ssh/id_rsa
ssh-manager add web_server --host web.com --user ubuntu --pass mysecretpassword
```

#### 2. `ssh-manager list` - List connections

Display a list of all saved SSH connections.

```bash
ssh-manager list

# List in JSON format:
ssh-manager list --format json
```

#### 3. `ssh-manager connect` - Connect to server

Establish an interactive SSH session with the saved server. You can also use the shorthand command.

```bash
ssh-manager connect <connection_name>

# Shorthand command:
ssh-manager <connection_name>

# Examples:
ssh-manager connect my_server
ssh-manager dev_machine
```

#### 4. `ssh-manager edit` - Edit connection

Edit the details of an existing SSH connection. Only the provided flags will be updated.

```bash
ssh-manager edit <connection_name> [--host <new_host>] [--user <new_user>] [--port <new_port>] [--key <new_key_path>] [--pass <new_password>]

# Examples:
ssh-manager edit my_server --port 22
ssh-manager edit dev_machine --user new_dev_user
```

#### 5. `ssh-manager remove` - Remove connection

Remove a saved SSH connection. The tool will ask for confirmation before removal.

```bash
ssh-manager remove <connection_name>

# Example: 
ssh-manager remove my_server
```

#### 6. `ssh-manager export` - Export configuration

Export all saved connections and configurations to a YAML file or print to standard output.

```bash
ssh-manager export --output backup.yaml

# Export to standard output:
ssh-manager export
```

#### 7. `ssh-manager import` - Import configuration

Import connections from a YAML file into the existing configuration. Connections with duplicate names will be skipped.

```bash
ssh-manager import --input backup.yaml

# Example:
ssh-manager import -i my_connections_backup.yaml
```

#### 8. `ssh-manager keys` - Manage SSH keys

Command group to manage SSH keys used by ssh-manager.

##### `ssh-manager keys list` - List SSH keys

List all managed SSH keys.

```bash
ssh-manager keys list
```

##### `ssh-manager keys add` - Add existing SSH key

Add an existing SSH key to the managed list.

```bash
ssh-manager keys add --name <key_name> --path <key_path>

# Example:
ssh-manager keys add --name work_key --path ~/.ssh/id_rsa_work
```

##### `ssh-manager keys generate` - Generate new SSH key pair

Generate a new SSH key pair (private and public key).

```bash
ssh-manager keys generate --name <key_name> [--type <key_type>] [--bits <bits_number>]

# Examples:
ssh-manager keys generate --name my_new_rsa_key --type rsa --bits 4096
ssh-manager keys generate --name my_ed25519_key --type ed25519
```

## Contributing

If you wish to contribute to the project, please refer to the `docs/DEVELOPMENT.md` file.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
