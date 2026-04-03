#!/usr/bin/env node

// Generates a demo .cast file (asciinema v2 format) and renders it to GIF via agg.

const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

const COLS = 90;
const ROWS = 30;
const PROMPT = "\u001b[32m❯\u001b[0m ";

// Asciinema v2 events
const events = [];
let t = 0;

function pause(sec) {
  t += sec;
}

function typeChar(ch) {
  const delay = 0.03 + Math.random() * 0.03; // 30-60ms jitter
  t += delay;
  events.push([t, "o", ch]);
}

function typeText(text) {
  for (const ch of text) {
    typeChar(ch);
  }
}

function output(text) {
  events.push([t, "o", text]);
}

function outputLine(text) {
  const delay = 0.06 + Math.random() * 0.06; // 60-120ms per line
  t += delay;
  events.push([t, "o", text + "\r\n"]);
}

function clearScreen() {
  t += 0.1;
  events.push([t, "o", "\u001b[2J\u001b[H"]);
}

function showPrompt() {
  output(PROMPT);
}

function comment(text) {
  // Dim italic comment
  outputLine("\u001b[3;90m# " + text + "\u001b[0m");
  pause(0.5);
}

function enter() {
  output("\r\n");
  pause(0.3);
}

// ─── Scene 1: Show config ───────────────────────────────────────────────────

clearScreen();
comment("gitwise — AI-powered commit messages from your terminal");
pause(1);

showPrompt();
typeText("cat ~/.config/gitwise/config.yaml");
enter();
pause(0.3);

const configOutput = [
  "\u001b[36mprovider\u001b[0m: ollama",
  "\u001b[36mmodel\u001b[0m: llama3:latest",
  "\u001b[36mlanguage\u001b[0m: english",
  "\u001b[36mconvention\u001b[0m: conventional",
  "\u001b[36mmax_diff_lines\u001b[0m: 5000",
  "\u001b[36mscope_from_path\u001b[0m: true",
  "\u001b[36memoji\u001b[0m: false",
  "\u001b[36mshow_cost\u001b[0m: true",
];
for (const line of configOutput) {
  outputLine(line);
}

pause(3.5);

// ─── Scene 2: Dry-run commit ────────────────────────────────────────────────

clearScreen();
comment("Stage some changes and preview the commit message");
pause(0.5);

showPrompt();
typeText("git add internal/auth/jwt.go internal/auth/session.go");
enter();
pause(0.5);

showPrompt();
typeText("gitwise commit --dry-run");
enter();
pause(0.8);

// Spinner simulation
const spinChars = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
for (let i = 0; i < 12; i++) {
  t += 0.15;
  events.push([t, "o", `\r\u001b[K\u001b[33m${spinChars[i % spinChars.length]} Generating commit message...\u001b[0m`]);
}
t += 0.1;
events.push([t, "o", "\r\u001b[K"]);

outputLine("");
outputLine("  \u001b[1;33mfeat(auth): add JWT refresh token rotation\u001b[0m");
outputLine("");
outputLine("  \u001b[37m- Implement automatic token refresh when access token expires\u001b[0m");
outputLine("  \u001b[37m- Add refresh token to session storage\u001b[0m");
outputLine("  \u001b[37m- Handle token rotation edge cases for concurrent requests\u001b[0m");
outputLine("");
outputLine("  \u001b[37mBREAKING CHANGE: auth middleware now requires refresh_token\u001b[0m");
outputLine("  \u001b[90mRefs: AUTH-234\u001b[0m");
outputLine("  \u001b[90m(~$0.0008, 623 input tokens)\u001b[0m");

pause(4);

// ─── Scene 3: AI code review ────────────────────────────────────────────────

clearScreen();
comment("Review your code before committing");
pause(0.5);

showPrompt();
typeText("gitwise review");
enter();
pause(0.8);

for (let i = 0; i < 10; i++) {
  t += 0.15;
  events.push([t, "o", `\r\u001b[K\u001b[33m${spinChars[i % spinChars.length]} Reviewing code...\u001b[0m`]);
}
t += 0.1;
events.push([t, "o", "\r\u001b[K"]);

outputLine("");
outputLine("\u001b[1;37m- File:\u001b[0m internal/auth/jwt.go");
outputLine("\u001b[1;37m- Line:\u001b[0m 47");
outputLine("\u001b[1;33m- Severity: warning\u001b[0m");
outputLine("\u001b[37m- Issue: Token expiry not validated before refresh\u001b[0m");
outputLine("\u001b[32m- Fix: Add time.Now().After(token.ExpiresAt) check\u001b[0m");
outputLine("");
outputLine("\u001b[1;37m- File:\u001b[0m internal/auth/session.go");
outputLine("\u001b[1;37m- Line:\u001b[0m 23");
outputLine("\u001b[1;36m- Severity: suggestion\u001b[0m");
outputLine("\u001b[37m- Issue: Session cleanup runs synchronously\u001b[0m");
outputLine("\u001b[32m- Fix: Use a background goroutine with ticker\u001b[0m");
outputLine("");
outputLine("\u001b[32mNo security issues found. Overall code quality is good.\u001b[0m");

pause(4);

// ─── Scene 4: Generate PR description ───────────────────────────────────────

clearScreen();
comment("Generate a PR description with suggested labels");
pause(0.5);

showPrompt();
typeText("gitwise pr --labels");
enter();
pause(0.8);

for (let i = 0; i < 10; i++) {
  t += 0.15;
  events.push([t, "o", `\r\u001b[K\u001b[33m${spinChars[i % spinChars.length]} Generating PR description...\u001b[0m`]);
}
t += 0.1;
events.push([t, "o", "\r\u001b[K"]);

outputLine("");
outputLine("\u001b[1;35m## Summary\u001b[0m");
outputLine("\u001b[37mAdd JWT-based authentication with automatic token refresh,\u001b[0m");
outputLine("\u001b[37mreplacing the legacy session-based auth system.\u001b[0m");
outputLine("");
outputLine("\u001b[1;35m## Changes\u001b[0m");
outputLine("\u001b[37m- Add JWT access/refresh token generation (internal/auth/jwt.go)\u001b[0m");
outputLine("\u001b[37m- Implement token refresh middleware (internal/middleware/auth.go)\u001b[0m");
outputLine("\u001b[37m- Add login/logout API endpoints (cmd/api/auth_handler.go)\u001b[0m");
outputLine("");
outputLine("\u001b[1;35m## Breaking Changes\u001b[0m");
outputLine("\u001b[37m- Session-based auth endpoints removed\u001b[0m");
outputLine("");
outputLine("\u001b[90mSuggested labels: feature, backend, security, breaking-change\u001b[0m");

pause(3);

// ─── Final ──────────────────────────────────────────────────────────────────

clearScreen();
outputLine("");
outputLine("  \u001b[1;32m✓ Install:\u001b[0m go install github.com/aymenhmaidiwastaken/gitwise@latest");
outputLine("  \u001b[1;32m✓ Setup:\u001b[0m  gitwise setup");
outputLine("  \u001b[1;32m✓ Commit:\u001b[0m gitwise commit");
outputLine("  \u001b[1;32m✓ Review:\u001b[0m gitwise review");
outputLine("  \u001b[1;32m✓ PR:\u001b[0m     gitwise pr --create");
outputLine("");
outputLine("  \u001b[90mhttps://github.com/aymenhmaidiwastaken/gitwise\u001b[0m");
outputLine("");

pause(4);

// ─── Write .cast file ───────────────────────────────────────────────────────

const header = {
  version: 2,
  width: COLS,
  height: ROWS,
  timestamp: Math.floor(Date.now() / 1000),
  env: { SHELL: "/bin/bash", TERM: "xterm-256color" },
};

const castLines = [JSON.stringify(header)];
for (const [ts, type, data] of events) {
  castLines.push(JSON.stringify([parseFloat(ts.toFixed(4)), type, data]));
}

const castPath = path.join(__dirname, "demo.cast");
fs.writeFileSync(castPath, castLines.join("\n") + "\n");
console.log(`Written ${castPath} (${events.length} events, ${t.toFixed(1)}s)`);

// ─── Render GIF ─────────────────────────────────────────────────────────────

const gifPath = path.join(__dirname, "..", "demo.gif");
const aggBin = path.join(process.env.HOME || process.env.USERPROFILE, ".local", "bin", "agg.exe");

try {
  console.log("Rendering GIF...");
  execSync(
    `"${aggBin}" "${castPath}" "${gifPath}" --theme dracula --font-size 16`,
    { stdio: "inherit" }
  );
  console.log(`Written ${gifPath}`);
} catch (err) {
  console.error("agg failed:", err.message);
  console.log("You can render manually: agg demo/demo.cast demo.gif --theme dracula --font-size 16");
}
