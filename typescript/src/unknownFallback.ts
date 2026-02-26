// Test cases covered: TC-32 (TS variant)
// Source: EX-TS-6

const dyn: any = {};

// <caret> TC-32: Delete 'toString()' below so only 'dyn.' remains,
//   invoke completion — should not crash; results may be limited
const tc38 = dyn.toString();
void tc38;

export {};
