#!/usr/bin/env node
/**
 * Dev server launcher — starts the Go API backend and Angular dev server together.
 * Usage: node dev-server.mjs [ng serve args...]
 *
 * The Go backend is spawned as a child process and killed when this script exits.
 */
import { spawn } from 'node:child_process';
import { request as httpRequest } from 'node:http';

const API_PORT = 8080;
const API_URL = `http://localhost:${API_PORT}/api/exercises?pageSize=1`;

function waitForApi(url, timeoutMs = 15000) {
  const start = Date.now();
  return new Promise((resolve, reject) => {
    function poll() {
      if (Date.now() - start > timeoutMs) {
        return reject(new Error(`API did not become ready within ${timeoutMs}ms`));
      }
      const req = httpRequest(url, (res) => {
        if (res.statusCode === 200) return resolve();
        setTimeout(poll, 500);
      });
      req.on('error', () => setTimeout(poll, 500));
      req.end();
    }
    poll();
  });
}

// Start Go API
const api = spawn('go', ['run', '.'], {
  cwd: '..',
  stdio: 'inherit',
  env: { ...process.env, DEV: 'true', AUTH_FALLBACK_USER: 'anon' },
});

function cleanup() {
  if (!api.killed) {
    api.kill('SIGTERM');
  }
}
process.on('exit', cleanup);
process.on('SIGINT', () => process.exit());
process.on('SIGTERM', () => process.exit());

api.on('error', (err) => {
  console.error('Failed to start Go API:', err.message);
  process.exit(1);
});

// Wait for API, then start ng serve
console.log('Waiting for Go API...');
try {
  await waitForApi(API_URL);
} catch (err) {
  console.error(err.message);
  process.exit(1);
}
console.log('Go API ready, starting Angular dev server...');

const ngArgs = ['ng', 'serve', ...process.argv.slice(2)];
const ng = spawn('npx', ngArgs, {
  stdio: 'inherit',
  env: { ...process.env, CI: 'true' },
});

ng.on('exit', (code) => {
  process.exit(code ?? 0);
});
