// Test cases covered: TC-31 (TS variant)
// Source: EX-TS-5

type User = { name: string; age: number };

function consume(u: User): void {}

const u: User = { name: "Ann", age: 21 };

// <caret> TC-31: Delete 'u' below inside consume(...), invoke completion;
//   should suggest variable "u" of type User
consume(u);

export {};
