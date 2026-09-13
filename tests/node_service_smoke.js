const { spawn } = require('node:child_process');
const path = require('node:path');

const [serviceDirectory, entrypoint, readyMessage, port] = process.argv.slice(2);

if (!serviceDirectory || !entrypoint || !readyMessage || !port) {
  console.error('usage: node_service_smoke.js <service-dir> <entrypoint> <ready-message> <port>');
  process.exit(2);
}

const child = spawn(process.execPath, [entrypoint], {
  cwd: path.resolve(serviceDirectory),
  env: {
    ...process.env,
    DISABLE_PROFILER: '1',
    ENABLE_TRACING: '0',
    PORT: port,
  },
  stdio: ['ignore', 'pipe', 'pipe'],
});

let output = '';
let ready = false;

const timer = setTimeout(() => {
  child.kill('SIGTERM');
  console.error(`service did not become ready:\n${output}`);
  process.exitCode = 1;
}, 10_000);

function capture(chunk) {
  output += chunk.toString();
  if (!ready && output.includes(readyMessage)) {
    ready = true;
    clearTimeout(timer);
    child.kill('SIGTERM');
  }
}

child.stdout.on('data', capture);
child.stderr.on('data', capture);

child.on('error', (error) => {
  clearTimeout(timer);
  console.error(error);
  process.exitCode = 1;
});

child.on('exit', (code, signal) => {
  clearTimeout(timer);
  if (!ready) {
    console.error(`service exited before readiness (code=${code}, signal=${signal}):\n${output}`);
    process.exitCode = 1;
  }
});
