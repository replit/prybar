// This import requires a bunch of polyfill so there's a pretty significant line offset
import { assertEqual, getCallSite } from './util';

const callSite = getCallSite();
assertEqual(callSite.line, 4, 'resolves correct line number');
assertEqual(
  // this should be at the start of the identifier (column 18)
  // instead, swc maps it to after the function call - (column 31)
  callSite.column,
  31,
  'resolves correct column number'
);
assertEqual(callSite.file, __filename, 'resolves correct filename');
