#!/usr/bin/env node

const { execSync, spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const ROOT = __dirname;
const FRONTEND_PORT = 3571;
const BACKEND_PORT = 3591;

const isWindows = os.platform() === 'win32';

const colors = {
  frontend: '\x1b[43m \x1b[0m',
  backend: '\x1b[44m \x1b[0m',
  system: '\x1b[46m \x1b[0m'
};

const run = (command, cwd = ROOT) => {
  execSync(command, { cwd, stdio: 'inherit', shell: true });
};

const exists = (...parts) => fs.existsSync(path.join(...parts));

const ensureDependencies = () => {
  if (!exists(ROOT, 'node_modules')) {
    console.log(`${colors.system} Installing root Node dependencies...`);
    run('npm install');
  }

  const frontendPath = path.join(ROOT, 'frontend');
  if (!exists(frontendPath, 'node_modules')) {
    console.log(`${colors.system} Installing frontend Node dependencies...`);
    run('npm install', frontendPath);
  }

  console.log(`${colors.system} Preparing Go modules...`);
  run('go mod tidy', path.join(ROOT, 'backend'));
};

const killPortIfInUse = (port) => {
  try {
    if (isWindows) {
      const output = execSync(`netstat -ano | findstr :${port}`, { encoding: 'utf8', shell: true });
      const pids = new Set(
        output
          .split('\n')
          .map((line) => line.trim().split(/\s+/).pop())
          .filter(Boolean)
      );
      for (const pid of pids) {
        execSync(`taskkill /PID ${pid} /F`, { stdio: 'ignore', shell: true });
      }
      return;
    }

    const output = execSync(`lsof -i :${port} -t`, { encoding: 'utf8', shell: true });
    const pids = output.split('\n').map((line) => line.trim()).filter(Boolean);
    for (const pid of pids) {
      execSync(`kill -9 ${pid}`, { stdio: 'ignore', shell: true });
    }
  } catch {
    // Port is free or lsof/netstat is not available.
  }
};

const logWithPrefix = (prefix, data) => {
  const lines = String(data).split('\n').filter((line) => line.trim().length > 0);
  for (const line of lines) {
    process.stdout.write(`${prefix} ${line}\n`);
  }
};

const startProcess = ({ name, prefix, command, args, cwd }) => {
  console.log(`${colors.system} Starting ${name}: ${command} ${args.join(' ')}`);
  const child = spawn(command, args, {
    cwd,
    shell: isWindows,
    stdio: ['ignore', 'pipe', 'pipe']
  });

  child.stdout.on('data', (data) => logWithPrefix(prefix, data));
  child.stderr.on('data', (data) => logWithPrefix(prefix, data));
  child.on('exit', (code) => logWithPrefix(prefix, `${name} exited with code ${code}`));

  return child;
};

const main = () => {
  ensureDependencies();
  [FRONTEND_PORT, BACKEND_PORT].forEach(killPortIfInUse);

  const children = [
    startProcess({
      name: 'frontend',
      prefix: colors.frontend,
      command: 'npm',
      args: ['run', 'dev'],
      cwd: path.join(ROOT, 'frontend')
    }),
    startProcess({
      name: 'backend',
      prefix: colors.backend,
      command: 'go',
      args: ['run', '.'],
      cwd: path.join(ROOT, 'backend')
    })
  ];

  const shutdown = () => {
    for (const child of children) {
      if (!child.killed) {
        child.kill('SIGTERM');
      }
    }
    process.exit(0);
  };

  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);
};

main();

