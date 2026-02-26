// Test cases covered: TC-55
// Source: EX-TS-BR-1

type User = { name: string; age: number };

const user: User = { name: "Ann", age: 21 };
// <caret> TC-55-A: Delete '"name"' below inside user[...], type 'na', invoke completion;
//   expect key completion for known property 'name'
const tc105a = user["name"];

const names = ["Ann", "Bob"];
const idx = 0;
// <caret> TC-55-B: Delete 'idx' below inside names[...], invoke completion;
//   expect numeric index variables from scope (e.g., idx)
const tc105b = names[idx];

void tc105a;
void tc105b;

export {};
