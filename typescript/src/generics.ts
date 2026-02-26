// Test cases covered: TC-74 (TS variant)
// Source: EX-TS-9

type Box<T> = { v: T };

// <caret> TC-74: Delete 'string' below inside Box<>, invoke completion;
//   completion inside angle brackets should suggest type arguments
let b: Box<string> = { v: "demo" };
void b;

export {};
