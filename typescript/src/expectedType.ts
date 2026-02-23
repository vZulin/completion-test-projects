// Test cases covered: TC-37 (TS variant)
// Source: EX-TS-5

type User = { name: string; age: number };

function consume(u: User): void {}

const u: User = { name: "Ann", age: 21 };

consume( // <caret> TC-37: expected type in args — should suggest variable "u" of type User
)
