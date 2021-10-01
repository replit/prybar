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

/**
 * Handles ANSI escapes from stdin.
 *
 * @return {string | -1 | 1} String if the escape isn't a left or right arrow .
 * Otherwise, -1 on left arrow and 1 on right arrow
 */
function handleArrowKey() {
  if (!readByteSync()) {
    return '^'
  }

  let str = buf.toString("binary");

  if (str !== "[") {
    return `^${str}`;
  }

  if (!readByteSync()) {
    return `^${str}`;
  }

  str += buf.toString("binary");

  switch (str) {
    case "[C": // \x1b[C -> right arrow key
      return 1;
    case "[D": // \x1b[D -> left arrow key
      return -1;
    default:
      return `^${str}`;
  }
}

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
 * Writes a string at an index of another string, appending as needed.
 *
 * @param {string} str The base string
 * @param {string} other The new string which is being written
 * @param {number} index The index at which the new string should start at.
 * @return {string} str with other written at index.
 */
function insertAt(str, other, index) {
  return [str.slice(0, index) + other + str.slice(index), index + other.length];
}

/**
 * Sets the current line to our promt + string w/ the cursor at the right index.
 *
 * @param {string} prompt The question's prompt
 * @param {string} current The current string (what the user has input so far)
 * @param {number} index The index that
 */
function displayPromptAndStr(prompt, current, index) {
  writeTTYOutput(
    // reset cursor position
    "\r" +
      // clear the rest of the line
      // EL (Erase in Line ): in this case, as no number is speciifed,
      // erases everything to the right of the cursor.
      "\x1b[K" +
      // write the prompt
      prompt +
      // write the string
      current +
      // CHA (Cursor Horizontal Absolute): move cursor to column
      // (starts at 1 for some reason)
      `\x1b[${prompt.length + index + 1}G`
  );
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
    let index = 0;

    if (!isTTY) {
      writeOutput(prompt);
    }

    for (;;) {
      displayPromptAndStr(prompt, str, index);
      const didRead = readByteSync();

      if (!didRead) {
        return str;
      }

      const char = buf.toString("binary");
      checkForSigs(char);

      if (char === "\n" || char === "\r") {
        writeTTYOutput("\r\n");

        return str;
      } else if (isTTY && char === "\x1b") {
        const ret = handleArrowKey();

        // if ret is a number, its the difference for the index
        if (typeof ret === "number") {
          // Only move the cursor if it will be in a valid position.
          const newIndex = index + ret;
          // the index can be equal to the strs length, if that's the case we're appending to the string.
          if (newIndex >= 0 && newIndex <= str.length) {
            index = newIndex;
          }

          // otherwise, the escape wasn't a left or right arrow key,
          // meaning we got an escaped version of the code.
        } else {
          [str, index] = insertAt(str, ret, index);
        }
        // \x7f: DEL
      } else if (isTTY && char === "\x7f") {
        if (index > 0) {
          index--;
          // remove the character at the old index
          str = str.slice(0, index) + str.slice(index + 1);
        }
      } else {
        [str, index] = insertAt(str, char, index);
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
