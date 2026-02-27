// Test cases covered:
// TC-1, TC-2, TC-7, TC-8, TC-9, TC-16, TC-19, TC-22, TC-23,
// TC-24, TC-25, TC-28, TC-34, TC-41, TC-63, TC-71 (TS variants)
// Source: EX-TS-1

import { User, buildUser, consumeUser } from "./model";

const userName = "Ann";
const userAge = 21;
const user: User = buildUser(userName, userAge);

// <caret> TC-1/TC-7/TC-22/TC-23/TC-71:
//   - place caret after "user." for basic/member completion
//   - type 'na' to validate prefix filtering
//   - with popup open, trigger QuickDoc for selected member
const memberName = user.name;

// <caret> TC-24: Delete 'toString()' below after "user.", invoke completion and choose method;
//   verify method call inserted with parentheses.
const methodCallResult = user.toString();

// <caret> TC-25: Delete 'name' below after "user.", type 'na', invoke completion;
//   verify property insertion without adding "()".
const propertyName = user.name;

// <caret> TC-2/TC-28: Delete 'user' below after '=', invoke smart completion;
//   should suggest buildUser(...) and variables of type User
const u1: User = user;

// <caret> TC-8/TC-9/TC-16/TC-34/TC-63:
//   - test popup after '(' and after ','
//   - verify second argument suggestions and parameter info
const u2: User = buildUser(userName, userAge);

// <caret> TC-41: Delete 'u2' below inside consumeUser(...), invoke completion;
//   selecting consumeUser(...) should place caret inside function-call parentheses.
consumeUser(u2);

function keywordReturnCompletion(value: number): number {
  // <caret> TC-19: Inside function body type 'ret' and invoke completion; expect 'return'.
  return value;
}

void memberName;
void methodCallResult;
void propertyName;
void keywordReturnCompletion;
