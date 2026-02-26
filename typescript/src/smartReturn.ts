// Test cases covered: TC-30 (TS variant)
// Source: EX-TS-4

type User = { name: string; age: number };

function buildUser(): User {
  return { name: "Ann", age: 21 };
}

function f(): User {
  // <caret> TC-30: Delete 'buildUser()' below after 'return', invoke smart completion;
  //   should suggest buildUser() or object literal matching User
  return buildUser();
}

const tc36 = f();
void tc36;

export {};
