// Test cases covered: TC-45
// Source: EX-TS-8
//
// This file is intentionally compilable.
// For TC-45, remove the import below before testing auto-import behavior.

import { utilFn } from "./utilModule";

// <caret> TC-45: Delete import above, then delete 'n()' below so only 'utilF' remains;
//   accept completion and verify import { utilFn } from "./utilModule" is added automatically.
const importedUtilResult = utilFn();
void importedUtilResult;
