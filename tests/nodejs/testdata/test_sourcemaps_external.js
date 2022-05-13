var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __markAsModule = (target) => __defProp(target, "__esModule", { value: true });
var __reExport = (target, module2, desc) => {
  if (module2 && typeof module2 === "object" || typeof module2 === "function") {
    for (let key of __getOwnPropNames(module2))
      if (!__hasOwnProp.call(target, key) && key !== "default")
        __defProp(target, key, { get: () => module2[key], enumerable: !(desc = __getOwnPropDesc(module2, key)) || desc.enumerable });
  }
  return target;
};
var __toModule = (module2) => {
  return __reExport(__markAsModule(__defProp(module2 != null ? __create(__getProtoOf(module2)) : {}, "default", module2 && module2.__esModule && "default" in module2 ? { get: () => module2.default, enumerable: true } : { value: module2, enumerable: true })), module2);
};
var import_util = __toModule(require("./util"));
var import_path = __toModule(require("path"));
/*! Generated from https://replit.com/@AllAwesome497/WorseMiserlyGenericsoftware#sourcemaps.ts */
/*! line 8 in src */
const callSite = (0, import_util.getCallSite)();
if (process.env.ENABLE_SOURCE_MAPS) {
  (0, import_util.assertEqual)(callSite.line, 7, "maps to original line number");
  (0, import_util.assertEqual)(callSite.column, 18, "resolves correct column number");
  (0, import_util.assertEqual)(callSite.file, import_path.default.join(__dirname, "../sourcemaps.ts"), "maps to original file");
} else {
  (0, import_util.assertEqual)(callSite.line, 23, "maps to interim line number");
  (0, import_util.assertEqual)(callSite.column, 46, "resolves interim column number");
  (0, import_util.assertEqual)(callSite.file, __filename, "maps to interim file name");
}
//# sourceMappingURL=test_sourcemaps_external.js.map