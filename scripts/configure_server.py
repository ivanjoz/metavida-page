#!/usr/bin/env python3

import json
import os
import pwd
import re
import shutil
import stat
import subprocess
import sys
from pathlib import Path
from urllib.parse import urlparse

SYSTEMD_DIRECTORY = Path("/etc/systemd/system")
SERVICE_NAME = "metavida.service"
RESTART_SERVICE_NAME = "metavida-restart.service"
RESTART_PATH_NAME = "metavida-restart.path"
NGINX_CONFIGURATION_DIRECTORY = Path("/etc/nginx/conf.d")
LETSENCRYPT_DIRECTORY = Path("/etc/letsencrypt/live")
CERTBOT_TLS_DIRECTIVES = {
    "ssl_certificate",
    "ssl_certificate_key",
    "ssl_trusted_certificate",
    "ssl_dhparam",
}


def print_debug(message_text):
    print(f"[*] {message_text}")


def fail_with_error(error_message):
    print(f"[!] {error_message}", file=sys.stderr)
    sys.exit(1)


def run_command(command_arguments):
    print_debug(f"Running command: {' '.join(command_arguments)}")
    command_result = subprocess.run(command_arguments, text=True, capture_output=True)
    print_debug(f"Exit code: {command_result.returncode}")

    if command_result.stdout.strip():
        print_debug("stdout:")
        for stdout_line in command_result.stdout.rstrip().splitlines():
            print(f"    {stdout_line}")

    if command_result.stderr.strip():
        print_debug("stderr:")
        for stderr_line in command_result.stderr.rstrip().splitlines():
            print(f"    {stderr_line}")

    if command_result.returncode != 0:
        fail_with_error(f"Command failed: {' '.join(command_arguments)}")

    return command_result


def require_root_execution():
    if os.geteuid() != 0:
        fail_with_error("This script must be executed as root.")


def detect_repository_credentials_path():
    # Prefer the deployed layout: configure.py, credentials.json, and server.bin in one directory.
    script_path = Path(__file__).resolve()
    script_directory_credentials_path = script_path.parent / "credentials.json"
    if script_directory_credentials_path.exists():
        print_debug(f"Using script-directory credentials path: {script_directory_credentials_path}")
        return script_directory_credentials_path

    # Fall back to local repository usage while developing the server setup script.
    for parent_path in script_path.parents:
        if (parent_path / "credentials.json").exists() and (parent_path / "backend").exists():
            repository_credentials_path = parent_path / "credentials.json"
            print_debug(f"Using repository credentials path: {repository_credentials_path}")
            return repository_credentials_path

    fail_with_error("Could not detect credentials.json next to configure_server.py or in the repository root.")


def load_project_credentials(repository_credentials_path):
    print_debug(f"Loading project credentials from {repository_credentials_path}")
    try:
        credentials_content = repository_credentials_path.read_text(encoding="utf-8")
    except OSError as read_error:
        fail_with_error(f"Could not read credentials.json: {read_error}")

    try:
        return json.loads(credentials_content)
    except json.JSONDecodeError as parse_error:
        fail_with_error(f"Could not parse credentials.json: {parse_error}")


def extract_server_configuration(project_credentials):
    raw_server_configuration = project_credentials.get("SERVER")
    if not isinstance(raw_server_configuration, dict):
        fail_with_error("credentials.json must contain a SERVER object.")

    remote_binary_path = str(raw_server_configuration.get("bin", "")).strip()
    if not remote_binary_path:
        fail_with_error("credentials.json SERVER.bin is required.")

    return {
        "binary_path": Path(remote_binary_path),
        "http_addr": str(project_credentials.get("HTTP_ADDR", ":3591")).strip() or ":3591",
    }


def extract_endpoint_configuration(project_credentials):
    endpoint_route = str(project_credentials.get("ENPOINT", "")).strip()
    endpoint_hostname = urlparse(endpoint_route).hostname or ""
    if not endpoint_route or not endpoint_hostname:
        fail_with_error("credentials.json ENPOINT must be a valid URL.")

    return {
        "route": endpoint_route,
        "hostname": endpoint_hostname,
    }


def build_backend_proxy_url(http_addr):
    normalized_http_addr = http_addr.strip()
    if normalized_http_addr.startswith(":"):
        return f"http://127.0.0.1{normalized_http_addr}"
    if normalized_http_addr.startswith("http://") or normalized_http_addr.startswith("https://"):
        return normalized_http_addr
    return f"http://{normalized_http_addr}"


def detect_runtime_username():
    ubuntu_account_exists = shutil.which("id") is not None and subprocess.run(
        ["id", "ubuntu"],
        text=True,
        capture_output=True,
    ).returncode == 0
    if ubuntu_account_exists:
        print_debug("Using existing 'ubuntu' account as the service runtime user.")
        return "ubuntu"

    sudo_username = os.environ.get("SUDO_USER", "").strip()
    if sudo_username and sudo_username != "root":
        print_debug(f"Using SUDO_USER '{sudo_username}' as the service runtime user.")
        return sudo_username

    fail_with_error(
        "Could not detect a non-root runtime user. Create 'ubuntu' or run via sudo from a non-root user."
    )


def resolve_runtime_user(runtime_username):
    try:
        return pwd.getpwnam(runtime_username)
    except KeyError as user_error:
        fail_with_error(f"Runtime user '{runtime_username}' does not exist: {user_error}")


def ensure_binary_placeholder(binary_path, runtime_user_entry):
    print_debug(f"Ensuring binary directory exists: {binary_path.parent}")
    binary_path.parent.mkdir(parents=True, exist_ok=True)

    if binary_path.exists():
        print_debug(f"Binary already exists at {binary_path}. Preserving it.")
    else:
        # The deploy script overwrites this placeholder with the compiled backend.
        print_debug(f"Creating placeholder binary at {binary_path}.")
        binary_path.touch()

    os.chown(binary_path, runtime_user_entry.pw_uid, runtime_user_entry.pw_gid)
    current_mode = stat.S_IMODE(binary_path.stat().st_mode)
    executable_mode = current_mode | 0o750
    os.chmod(binary_path, executable_mode)
    print_debug(f"Binary permissions set to {oct(executable_mode)}.")


def build_main_service_contents(runtime_username, binary_path, repository_credentials_path, http_addr):
    repository_root_path = repository_credentials_path.parent
    deployed_credentials_path = binary_path.parent / "credentials.json"
    return f"""[Unit]
Description=Metavida Backend Service
After=network.target

[Service]
Type=simple
User={runtime_username}
Group={runtime_username}
WorkingDirectory={binary_path.parent}
Environment=METAVIDA_CREDENTIALS={deployed_credentials_path}
Environment=METAVIDA_REPOSITORY_ROOT={repository_root_path}
Environment=HTTP_ADDR={http_addr}
ExecStart={binary_path}
Restart=always
RestartSec=5

# Security hardening keeps the process non-root and limits write access to the binary directory only.
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=read-only
ProtectControlGroups=yes
ProtectKernelModules=yes
ProtectKernelTunables=yes
RestrictRealtime=yes
CapabilityBoundingSet=
ReadWritePaths={binary_path.parent}

[Install]
WantedBy=multi-user.target
"""


def build_restart_path_contents(binary_path):
    return f"""[Unit]
Description=Watch for changes to metavida backend binary

[Path]
PathChanged={binary_path}

[Install]
WantedBy=multi-user.target
"""


def build_restart_service_contents():
    return f"""[Unit]
Description=Restart Metavida Service

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl restart {SERVICE_NAME}
"""


def write_unit_file(unit_file_path, unit_contents):
    existing_unit_contents = None
    if unit_file_path.exists():
        existing_unit_contents = unit_file_path.read_text(encoding="utf-8")

    if existing_unit_contents == unit_contents:
        print_debug(f"Configuration unchanged: {unit_file_path}")
        return False

    print_debug(f"Writing systemd unit: {unit_file_path}")
    unit_file_path.write_text(unit_contents, encoding="utf-8")
    os.chmod(unit_file_path, 0o644)
    return True


def configure_systemd_units(runtime_username, binary_path, repository_credentials_path, http_addr):
    main_service_configuration_changed = write_unit_file(
        SYSTEMD_DIRECTORY / SERVICE_NAME,
        build_main_service_contents(
            runtime_username,
            binary_path,
            repository_credentials_path,
            http_addr,
        ),
    )
    restart_service_configuration_changed = write_unit_file(
        SYSTEMD_DIRECTORY / RESTART_SERVICE_NAME,
        build_restart_service_contents(),
    )
    restart_path_configuration_changed = write_unit_file(
        SYSTEMD_DIRECTORY / RESTART_PATH_NAME,
        build_restart_path_contents(binary_path),
    )
    return (
        main_service_configuration_changed
        or restart_service_configuration_changed
        or restart_path_configuration_changed
    )


def ensure_nginx_is_installed():
    nginx_binary_path = shutil.which("nginx")
    if not nginx_binary_path:
        fail_with_error("Nginx is not installed or not available in PATH.")

    if not NGINX_CONFIGURATION_DIRECTORY.exists():
        fail_with_error(f"Nginx configuration directory not found: {NGINX_CONFIGURATION_DIRECTORY}")

    print_debug(f"Detected Nginx binary at {nginx_binary_path}")


def extract_existing_certbot_tls_lines(existing_nginx_configuration_contents):
    preserved_tls_lines = []
    seen_directives = set()

    for raw_line in existing_nginx_configuration_contents.splitlines():
        stripped_line = raw_line.strip()
        if not stripped_line or stripped_line.startswith("#"):
            continue

        directive_match = re.match(r"^([A-Za-z0-9_]+)\s+(.+?);(?:\s*#.*)?$", stripped_line)
        if not directive_match:
            continue

        directive_name = directive_match.group(1)
        directive_value = directive_match.group(2)
        is_certbot_include = (
            directive_name == "include"
            and ("letsencrypt" in directive_value or "certbot" in directive_value.lower())
        )
        should_preserve_directive = directive_name in CERTBOT_TLS_DIRECTIVES or is_certbot_include
        if not should_preserve_directive:
            continue

        # Keep the first copy of each TLS directive so repeated runs stay stable.
        directive_key = f"{directive_name}:{directive_value}"
        if directive_key in seen_directives:
            continue
        seen_directives.add(directive_key)
        preserved_tls_lines.append(f"    {stripped_line}")

    return preserved_tls_lines


def build_tls_directive_lines(endpoint_hostname, existing_nginx_configuration_contents):
    preserved_tls_lines = extract_existing_certbot_tls_lines(existing_nginx_configuration_contents)
    has_preserved_certificate = any(line.strip().startswith("ssl_certificate ") for line in preserved_tls_lines)
    has_preserved_certificate_key = any(line.strip().startswith("ssl_certificate_key ") for line in preserved_tls_lines)
    if has_preserved_certificate and has_preserved_certificate_key:
        return preserved_tls_lines

    certificate_directory = LETSENCRYPT_DIRECTORY / endpoint_hostname
    certificate_fullchain_path = certificate_directory / "fullchain.pem"
    certificate_private_key_path = certificate_directory / "privkey.pem"
    if certificate_fullchain_path.exists() and certificate_private_key_path.exists():
        return [
            f"    ssl_certificate {certificate_fullchain_path};",
            f"    ssl_certificate_key {certificate_private_key_path};",
        ]

    return []


def build_nginx_configuration(endpoint_hostname, backend_proxy_url, existing_nginx_configuration_contents):
    tls_directive_lines = build_tls_directive_lines(
        endpoint_hostname,
        existing_nginx_configuration_contents,
    )

    common_proxy_location = f"""    location / {{
        # Forward API traffic to the local Metavida backend process.
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_pass {backend_proxy_url};

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 3600s;
        proxy_buffering on;
        proxy_buffer_size 16k;
        proxy_buffers 4 16k;
    }}"""

    if tls_directive_lines:
        tls_directives = "\n".join(tls_directive_lines)
        return f"""server {{
    listen 443 ssl;
    listen [::]:443 ssl;

    server_name {endpoint_hostname};

{tls_directives}

{common_proxy_location}
}}
"""

    return f"""server {{
    listen 80;
    listen [::]:80;

    server_name {endpoint_hostname};

{common_proxy_location}
}}
"""


def configure_nginx_reverse_proxy(endpoint_configuration, http_addr):
    ensure_nginx_is_installed()

    endpoint_hostname = endpoint_configuration["hostname"]
    backend_proxy_url = build_backend_proxy_url(http_addr)
    nginx_configuration_path = NGINX_CONFIGURATION_DIRECTORY / f"{endpoint_hostname}.conf"

    existing_nginx_configuration_contents = ""
    if nginx_configuration_path.exists():
        existing_nginx_configuration_contents = nginx_configuration_path.read_text(encoding="utf-8")

    nginx_configuration_contents = build_nginx_configuration(
        endpoint_hostname,
        backend_proxy_url,
        existing_nginx_configuration_contents,
    )

    if existing_nginx_configuration_contents == nginx_configuration_contents:
        print_debug(f"Nginx configuration unchanged: {nginx_configuration_path}")
        return False

    print_debug(f"Writing Nginx reverse proxy config: {nginx_configuration_path}")
    nginx_configuration_path.write_text(nginx_configuration_contents, encoding="utf-8")
    os.chmod(nginx_configuration_path, 0o644)

    run_command(["nginx", "-t"])
    run_command(["systemctl", "enable", "nginx"])
    run_command(["systemctl", "restart", "nginx"])
    return True


def enable_units():
    run_command(["systemctl", "enable", SERVICE_NAME])
    run_command(["systemctl", "enable", RESTART_PATH_NAME])


def reload_systemd_if_configuration_changed(systemd_configuration_changed):
    if not systemd_configuration_changed:
        print_debug("Systemd configuration unchanged. Skipping daemon-reload and watcher restart.")
        return

    run_command(["systemctl", "daemon-reload"])
    run_command(["systemctl", "restart", RESTART_PATH_NAME])


def print_summary(runtime_username, binary_path, repository_credentials_path, http_addr, endpoint_configuration):
    deployed_credentials_path = binary_path.parent / "credentials.json"
    print_debug("Configuration completed.")
    print_debug(f"Runtime user: {runtime_username}")
    print_debug(f"Binary path: {binary_path}")
    print_debug(f"HTTP address: {http_addr}")
    print_debug(f"Repository credentials path: {repository_credentials_path}")
    print_debug(f"Deployed credentials path: {deployed_credentials_path}")
    print_debug(f"Nginx endpoint: {endpoint_configuration['route']}")
    print_debug(f"Main service unit: {SYSTEMD_DIRECTORY / SERVICE_NAME}")
    print_debug(f"Path watcher unit: {SYSTEMD_DIRECTORY / RESTART_PATH_NAME}")
    print_debug(f"Restart helper unit: {SYSTEMD_DIRECTORY / RESTART_SERVICE_NAME}")
    print_debug(f"Start service with: sudo systemctl start {SERVICE_NAME}")


def main():
    require_root_execution()
    repository_credentials_path = detect_repository_credentials_path()
    project_credentials = load_project_credentials(repository_credentials_path)
    server_configuration = extract_server_configuration(project_credentials)
    endpoint_configuration = extract_endpoint_configuration(project_credentials)
    runtime_username = detect_runtime_username()
    runtime_user_entry = resolve_runtime_user(runtime_username)
    ensure_binary_placeholder(server_configuration["binary_path"], runtime_user_entry)
    systemd_configuration_changed = configure_systemd_units(
        runtime_username,
        server_configuration["binary_path"],
        repository_credentials_path,
        server_configuration["http_addr"],
    )
    configure_nginx_reverse_proxy(endpoint_configuration, server_configuration["http_addr"])
    enable_units()
    reload_systemd_if_configuration_changed(systemd_configuration_changed)
    print_summary(
        runtime_username,
        server_configuration["binary_path"],
        repository_credentials_path,
        server_configuration["http_addr"],
        endpoint_configuration,
    )


if __name__ == "__main__":
    main()
