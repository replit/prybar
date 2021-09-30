const { readSync, writeSync, openSync } = require("fs");
const { isatty } = require("tty");
const buf = Buffer.alloc(1);
const isTTY = isatty(process.stdin.fd);

/**
 * Reads a single byte from stdin to buf.
 * @return {boolean} Whether a byte was read or not (false on EOF)
 */
function readByteSync() {
  const read = readSync(stdinFd, buf);

  return read > 0;
}

/**
 * The file descriptor of our stdin stream.
 * @type {number}
 */
const stdinFd = isTTY
  ? // We can't just use process.stdin.fd here since node has some getter shenanigans
    // which cause sync reads to throw
    openSync("/dev/tty", "r")
  : openSync("/dev/stdin", "r");

/**
 * Writes output to stdout.
 *
 * @param {NodeJS.ArrayBufferView} b The outpus which should be written to stdout.
 */
function writeOutput(b) {
  writeSync(process.stdout.fd, b);
}

/**
 * Writes a buffer or string to stdout, only if stdin is a tty.
 * @param {NodeJS.ArrayBufferView} b The outpus which should be written to stdout (if stdin is a tty).
 */
function writeTTYOutput(b) {
  if (isTTY) {
    writeOutput(b);
  }
}

/**
 * If stdin is a tty, ensure raw mode is on prior to executing cb,
 * then reset it to its previous state afterwards.  If stdin isn't a tty,
 * cb will just be executed normally.
 *
 * @template T
 * @param {() => T} cb The function which is being executed.
 * @return {T} The return value of cb
 */
function ensureRawMode(cb) {
  if (!isTTY) {
    return cb();
  }

  const previousRawMode = process.stdin.isRaw;

  let ret;

  try {
    process.stdin.setRawMode(true);
    ret = cb();
  } finally {
    process.stdin.setRawMode(previousRawMode);
  }

  return ret;
}

const delChar = "\x7f";
const escapeClearLineRight = "\x1b[K";

/**
 * Checks to see if the input character is what we get in raw mode for
 * Ctrl+C or Ctrl+D, and if so sends the proc SIGINT>
 *
 * @param {string} char The character read.
 */
function checkForSigs(char) {
  if (isTTY && (char === "\x03" || char === "\x04")) {
    process.exit();
  }
}

/**
 * Synchronously reads from stdin until `\n` or `\r`
 *
 * @param {string} prompt The prompt to be displayed
 * @return {string} The input read (excluding newlines)
 */
function question(prompt) {
  return ensureRawMode(() => {
    let str = "";

    writeOutput(prompt);

    for (;;) {
      const didRead = readByteSync();

      if (!didRead) {
        return str;
      }

      const char = buf.toString("binary");
      checkForSigs(char);

      if (char === "\n" || char === "\r") {
        writeTTYOutput("\r\n");

        return str;
      } else if (isTTY && char === delChar) {
        if (str.length > 0) {
          writeOutput("\b" + escapeClearLineRight);
          str = str.slice(0, str.length - 1);
        }
      } else {
        writeTTYOutput(buf);
        str += char;
      }
    }
  });
}

/**
 * Synchronously waits for 'y' or 'n' (case insensitive) and returns a boolean
 * where 'y' is true and 'n' is false.
 *
 * @param {string | undefined | null} prompt The prompt.
 * @returns {boolean} True if the input was y; false otherwise.
 * @throws {Error} If EOF is received before y/n
 */
function keyInYNStrict(prompt) {
  return ensureRawMode(() => {
    writeOutput(`${prompt == null ? "Are you sure?" : prompt} [y/n]: `);

    for (;;) {
      const didRead = readByteSync();

      if (!didRead) {
        throw new Error("Unexpected EOF / end of input.  Expected y/n.");
      }

      const char = buf.toString("binary");
      checkForSigs(char);

      if (char.match(/[yn]/i)) {
        writeTTYOutput(`${char}\r\n`);

        return char === "y" || char === "Y";
      }
    }
  });
}

module.exports = {
  question,
  keyInYNStrict,
};
