const { assertEqual } = require('./util');
import get2B, { b, alias } from './somefile';

assertEqual(b, 32, 'correctly imports b');
assertEqual(b, alias, 'correctly imports aliased exports');
assertEqual(get2B(), 64, 'correctly imports default exports');
