# SM

![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)

## Introduction

**SM** is a command-line interface (CLI) tool developed in Go (Golang) designed to help you manage SSH connections efficiently. This tool allows you to store, organize, and quickly connect to remote servers without having to remember complex SSH commands and configurations.

With SM, you can:

*   Add, list, edit, and remove SSH connections.
*   Encrypt passwords for secure storage.
*   Manage and generate SSH key pairs.
*   Export and import connection configurations.
*   Use shorthand commands for quick connections.

## Installation

SM is a standalone command-line tool. Once built, you do not need Go installed on your machine to run it.

### 1. Download Executable (Recommended)

You can download the pre-built executable for your operating system from the project's [Releases](https://github.com/haidang1810/sm/releases) page. Choose the version compatible with your operating system and CPU architecture (e.g., `sm-v1.0.0-linux-amd64.tar.gz` for Linux 64-bit, or `sm-v1.0.0-windows-amd64.zip` for Windows 64-bit).

After downloading, follow the instructions below depending on your operating system:

#### a. On Linux and macOS

1.  **Download:** Download the appropriate `.tar.gz` file for your system (e.g., `sm-<version>-linux-amd64.tar.gz`). You can use `curl` to download it directly:

    ```bash
    curl -LO https://github.com/haidang1810/sm/releases/download/v1.0.0/sm-v1.0.0-<os>-<arch>.tar.gz
    ```
2.  **Extract:** Open Terminal and extract the downloaded file, and rename it to `sm`:

    ```bash
    tar -xzf sm-<version>-<os>-<arch>.tar.gz -O sm-<os>-<arch> > sm
    chmod +x sm
    ```
3.  **Move to PATH:** Move the executable to a directory included in your system's `PATH` environment variable (e.g., `/usr/local/bin/`). This allows you to run the `sm` command from anywhere in the Terminal.

    ```bash
    sudo mv sm /usr/local/bin/
    ```

4.  **Verify Installation:** Open a new Terminal and run the following command to verify.

    ```bash
    sm
    ```

#### b. On Windows

1.  **Download:** Download the appropriate `.zip` file for your system (e.g., `sm-<version>-windows-amd64.zip`).
2.  **Extract:** Extract the downloaded `.zip` file to a directory where you want to store the tool (e.g., `C:\Program Files\sm`).
3.  **Add to PATH Environment Variable:** To be able to run `sm` from anywhere in Command Prompt or PowerShell, you need to add the path to the directory containing the executable to your system's `PATH` environment variable.
    *   Search for "Edit the system environment variables" in the Start Menu.
    *   Click on "Environment Variables...".
    *   In the "System variables" section, find and select the `Path` variable, then click "Edit...".
    *   Click "New" and add the path to the directory where you extracted `sm` (e.g., `C:\Program Files\sm`).
    *   Click "OK" on all windows to save changes.
4.  **Verify Installation:** Open a new Command Prompt or PowerShell (it's important to open a new window for the environment variables to take effect) and run the following command to verify.

    ```cmd
    sm
    ```

### 2. Build from Source (For Developers)

If you want to build the tool from source, you need to have Go (version 1.21 or later) installed on your system.

1.  **Clone repository:**

    ```bash
    git clone https://github.com/haidang1810/sm.git # Replace with actual repo URL
    cd sm
    ```

2.  **Build project:**

    ```bash
    go mod tidy
    go build -o sm .
    ```

    The `sm` executable will be created in the current directory. You can move it to a directory in your `PATH` as instructed above.

## Usage Guide

SM provides a set of intuitive commands to manage your SSH connections.

### Configuration

The tool will automatically create a configuration file at `$HOME/.sm/config.yaml` or `$HOME/.config/sm/config.yaml` when you use it for the first time.

### Main Commands

#### 1. `sm add` - Add a new connection

Add a new SSH connection to your configuration. If you do not provide required flags, the tool will prompt you for information.

```bash
sm add <connection_name> --host <host> --user <user> [--port <port>] [--key <key_path>] [--pass <password>]

# Examples:
sm add my_server --host example.com --user admin --port 2222
sm add dev_machine --host 192.168.1.100 --user dev --key ~/.ssh/id_rsa
sm add web_server --host web.com --user ubuntu --pass mysecretpassword
```

#### 2. `sm list` - List connections

Display a list of all saved SSH connections.

```bash
sm list

# List in JSON format:
sm list --format json
```

#### 3. `sm connect` - Connect to server

Establish an interactive SSH session with the saved server. You can also use the shorthand command.

```bash
sm connect <connection_name>

# Shorthand command:
sm <connection_name>

# Examples:
sm connect my_server
sm dev_machine
```

#### 4. `sm edit` - Edit connection

Edit the details of an existing SSH connection. Only the provided flags will be updated.

```bash
sm edit <connection_name> [--host <new_host>] [--user <new_user>] [--port <new_port>] [--key <new_key_path>] [--pass <new_password>]

# Examples:
sm edit my_server --port 22
sm edit dev_machine --user new_dev_user
```

#### 5. `sm remove` - Remove connection

Remove a saved SSH connection. The tool will ask for confirmation before removal.

```bash
sm remove <connection_name>

# Example: 
sm remove my_server
```

#### 6. `sm export` - Export configuration

Export all saved connections and configurations to a YAML file or print to standard output.

```bash
sm export --output backup.yaml

# Export to standard output:
sm export
```

#### 7. `sm import` - Import configuration

Import connections from a YAML file into the existing configuration. Connections with duplicate names will be skipped.

```bash
sm import --input backup.yaml

# Example:
sm import -i my_connections_backup.yaml
```

#### 8. `sm keys` - Manage SSH keys

Command group to manage SSH keys used by sm.

##### `sm keys list` - List SSH keys

List all managed SSH keys.

```bash
sm keys list
```

##### `sm keys add` - Add existing SSH key

Add an existing SSH key to the managed list.

```bash
sm keys add --name <key_name> --path <key_path>

# Example:
sm keys add --name work_key --path ~/.ssh/id_rsa_work
```

##### `sm keys generate` - Generate new SSH key pair

Generate a new SSH key pair (private and public key).

```bash
sm keys generate --name <key_name> [--type <key_type>] [--bits <bits_number>]

# Examples:
sm keys generate --name my_new_rsa_key --type rsa --bits 4096
sm keys generate --name my_ed25519_key --type ed25519
```

## Contributing

If you wish to contribute to the project, please refer to the `docs/DEVELOPMENT.md` file.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
