// Test cases covered: TC-36 (TS variant)
// Source: EX-TS-4

type User = { name: string; age: number };

function buildUser(): User {
  return { name: "Ann", age: 21 };
}

function f(): User {
  return // <caret> TC-36: smart return completion — should suggest buildUser() or object literal matching User
}
