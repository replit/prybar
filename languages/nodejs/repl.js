const util = require('util');
const repl = require("repl");
const path = require("path");
const fs = require("fs");
const vm = require("vm");
const rl = require(path.join(
  process.cwd(),
  "prybar_assets",
  "nodejs",
  "readline-sync.js"
));
const tty = require("tty");
const Module = require("module");

let r;
if (!process.env.PRYBAR_QUIET) {
  console.log("Node " + process.version + " on " + process.platform);
}

// Red errors.
function logError(msg) {
  process.stdout.write("\u001b[0m\u001b[31m" + msg + "\u001b[0m");
}

// The nodejs repl operates in raw mode and does some funky stuff to
// the terminal. This ns the repl and forces non-raw mode.
function pauseRepl() {
  if (!r) return;

  r.pause();
  process.stdin.setRawMode(false);
}

// Forces raw mode and resumes the repl.
function resumeRepl() {
  if (!r) return;

  process.stdin.setRawMode(true);
  r.resume();
}

// Clear the line if it has anything on it.
function clearLine() {
  if (r && r.line) r.clearLine();
}

// Adapted from the internal node repl code just a lot simpler and adds
// red errors (see https://bit.ly/2FRM86S)
function handleError(e) {
  if (r) {
    r.lastError = e;
  }

  if (e && typeof e === "object" && e.stack && e.name) {
    if (e.name === "SyntaxError") {
      e.stack = e.stack
        .replace(/^repl:\d+\r?\n/, "")
        .replace(/^\s+at\s.*\n?/gm, "");
    }

    logError(e.stack);
  } else {
    // For some reason needs a newline to flush.
    logError("Thrown: " + r.writer(e) + "\n");
  }

  if (r) {
    r.clearBufferedCommand();
    r.lines.level = [];
    r.displayPrompt();
  }
}

function start(context) {
  r = repl.start({
    prompt: process.env.PRYBAR_PS1,
  });
  if (context) r.context = context;

  // remove the internal error and ours for red etc.
  r._domain.removeListener("error", r._domain.listeners("error")[0]);
  r._domain.on("error", handleError);
  process.on("uncaughtException", handleError);
}

global.alert = console.log;
global.prompt = (p) => {
  pauseRepl();
  clearLine();

  let ret = rl.question(`${p}> `, {
    hideEchoBack: false,
  });

  resumeRepl();

  // Display prompt on the next turn.
  if (r) setImmediate(() => r.displayPrompt());

  return ret;
};

global.confirm = (q) => {
  pauseRepl();
  clearLine();

  const ret = rl.keyInYNStrict(q);

  resumeRepl();

  // Display prompt on the next turn.
  if (r) setImmediate(() => r.displayPrompt());
  return ret;
};

if (process.env.PRYBAR_CODE) {
  vm.runInThisContext(process.env.PRYBAR_CODE);
  if (process.env.PRYBAR_INTERACTIVE) {
    start();
  }
} else if (process.env.PRYBAR_EXP) {
  console.log(vm.runInThisContext(process.env.PRYBAR_EXP));
  if (process.env.PRYBAR_INTERACTIVE) {
    start();
  }
} else if (process.env.PRYBAR_FILE) {
  const mainPath = path.resolve(process.env.PRYBAR_FILE);
  const main = fs.readFileSync(mainPath, "utf-8");
  const module = new Module(mainPath, null);

  module.id = ".";
  module.filename = mainPath;
  module.paths = Module._nodeModulePaths(path.dirname(mainPath));

  process.mainModule = module;
  const sandbox = {
    module,
    require: module.require.bind(module),
    __dirname: path.dirname(mainPath),
    __filename: mainPath,

    // These are deprecated properties and accessing them will trigger a warning.
    // We add them manually for backward compat.
    GLOBAL: global,
    root: global,
  };

  // These properties will show a warning. We can avoid them.
  for (const prop of Object.getOwnPropertyNames(global)) {
    if (sandbox.hasOwnProperty(prop)) {
      continue;
    }

    sandbox[prop] = global[prop];
  }

  console.log(
    "\u001b[0m\u001b[90mHint: hit control+c anytime to enter REPL.\u001b[0m"
  );
  const context = vm.createContext(sandbox);

  let script;
  try {
    script = vm.createScript(main, {
      filename: mainPath,
      displayErrors: false,
    });
  } catch (e) {
    handleError(e);
  }

  if (script) {
    let res;
    try {
      res = script.runInContext(context, {
        displayErrors: false,
      });
    } catch (e) {
      handleError(e);
    }

    module.loaded = true;

    if (typeof res !== "undefined") {
      console.log(util.inspect(res, { colors: true }));
    }
  }

  if (process.env.PRYBAR_INTERACTIVE) {
    process.once("SIGINT", () => start(context));
    process.once("beforeExit", () => start(context));
  }
} else if (process.env.PRYBAR_INTERACTIVE) {
  start();
}
