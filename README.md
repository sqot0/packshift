<p align="center">
  <img src="https://i.imgur.com/sFOMABG.png" alt="Packshift - Deploy modpack easy" />
</p>

# Packshift

Packshift is a simple command-line tool for uploading and deploying Minecraft modpacks (and other files) to remote FTP/SFTP servers. It is designed to be compatible with the [Packsmith](https://github.com/sqot0/packsmith) ecosystem — Packshift is a complementary CLI utility that helps you publish and sync modpack server files reliably.

Key capabilities:
- Upload local directories/files to a remote FTP or SFTP server.
- Configure multiple local -> remote path mappings.
- Concurrent uploads with safe remote directory creation.
- Passwords are stored encrypted in the local config file.
- Designed to work alongside Packsmith as an upload/deploy helper.

> [!WARNING]
> Packshift is under active development. Expect minor bugs or incomplete features. Test on non-critical data first.

## Requirements
- Windows 10 or later (development targeted on Windows; Unix systems should work but paths may differ).
- Network access to your FTP/SFTP server.
- Go toolchain (only required if building from source).

## Quick Start

1. Download the latest Packshift release (or build from source, see below).
2. Open a terminal in your project root (folder where you want `packshift.json` located).
3. Run:
   - `packshift init` — interactively create `packshift.json`.
   - `packshift deploy` — upload configured files to the remote server.

After `init` completes you will have a `packshift.json` that contains your connection info and path mappings.

## Usage Summary

- `packshift init` — Create or overwrite `packshift.json` by answering prompts (host, port, username, password, protocol, path mappings).
- `packshift deploy` — Read `packshift.json` and upload the mapped files to the remote server.

Example:
1. In your modpack project directory:
   - `packshift init`
     - FTP Host: example.com
     - FTP Port: 22
     - FTP Username: myuser
     - FTP Password: (your password)
     - Choose protocol: SFTP
     - Local path (relative): ./server_extras
     - Remote path on server: ./
     - Add another mapping? No
2. `packshift deploy` — uploads contents of `server_extras` (recursively) to the remote path.

## init: interactive guide and tips

When you run `packshift init`, you'll be asked to provide:
- FTP Host — hostname or IP of your FTP/SFTP server.
- FTP Port — common ports: 21 (FTP), 22 (SFTP). Default prompts may offer 21.
- FTP Username — server account user.
- FTP Password — entered securely. Packshift will encrypt this before saving.
- Protocol — choose between FTP and SFTP. SFTP uses SSH and is recommended when available.
- Path mappings — define one or more mappings:
  - Local path (relative to current directory) — e.g. `server_extras` or `config`.
  - Remote path — e.g. `.` or `config`
  - Repeat until you finish adding mappings.

Notes:
- Local paths are interpreted relative to where you run `packshift init` (project root).
- Remote paths should use forward slashes `/`. Packshift will normalize Windows paths when forming remote paths.
- Passwords are encrypted before saving to `packshift.json`, but keep the file private.

## Example `packshift.json`

A typical `packshift.json` will look like:

```json
{
  "ftpConfig": {
    "host": "example.com",
    "port": 22,
    "username": "myuser",
    "password": "ENCRYPTED_VALUE",
    "ssl": true
  },
  "pathMappings": {
    "server_extras": ".",
    "config": "config"
  }
}
```

- `ftpConfig.ssl` indicates whether the SFTP (SSL) client should be used.
- `pathMappings` keys are local relative paths, values are remote paths.

## Deploy behavior

- `packshift deploy` reads `packshift.json` and for each mapping uploads files found in the local path to the corresponding remote path.
- Nested directories are created on the server when necessary.
- Uploads are performed concurrently with a limited number of worker goroutines for throughput and to avoid overloading the server.
- The tool logs progress to stdout/stderr. Inspect output for any upload errors.

## Building from source

If you want to build Packshift yourself:

1. Install Go (1.18+ recommended).
2. Clone repository:
   - `git clone https://github.com/sqot0/packshift.git`
3. Build:
   - `cd packshift`
   - `go build`

This produces an executable named `packshift` (or `packshift.exe` on Windows).

## Integration with Packsmith

Packshift is intended as a CLI utility for quickly publishing modpack server files to FTP/SFTP servers. In a Packsmith workflow:
- Use Packsmith to manage mods, versions, and local structure.
- Use Packshift to push server-side files (mods, configs, extras) from the project to a hosted server.

Packshift does not replace Packsmith; it complements it by handling network deployment.

## Troubleshooting

- Remote file saved with backslashes in the name on Windows:
  - Packshift normalizes remote paths to use forward slashes. If you still see backslashes, ensure your mapping values are proper remote paths (use `/`), and update `packshift.json` if needed.
- Permission denied / authentication errors:
  - Verify host, port, username, password, and protocol. Test with an SFTP/FTP client (e.g., `sftp`, FileZilla).
- If uploads fail for nested files:
  - Ensure the server user has permission to create directories. Packshift attempts to create directories before writing files (SFTP) but may fail if permissions are restricted.

## Contributing

Contributions are welcome:
- Open issues for bugs and feature requests.
- Send pull requests for code changes.
- Help improve docs, examples, and tests.

When contributing, please include clear descriptions and reproducible steps.
