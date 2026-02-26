// Test cases covered: TC-7 (TS), TC-9 (TS), TC-22, TC-23, TC-24, TC-25 (TS variants)
// Source: EX-TS-1

import { User, buildUser, consumeUser } from "./model";

const userName = "Ann";
const userAge = 21;
const user: User = buildUser(userName, userAge);

// <caret> TC-7/TC-22: Delete 'name' below, place caret after "user." and invoke completion;
//   expect member suggestions and filtering by prefix (e.g. "na" -> "name")
const tc26 = user.name;

// <caret> TC-23: Delete 'user' below after '=', invoke smart completion;
//   should suggest buildUser(...) and variables of type User
const u1: User = user;

// <caret> TC-24: Delete arguments below in buildUser(...), invoke completion;
//   expect userName for first arg and userAge for second arg
const u2: User = buildUser(userName, userAge);

// <caret> TC-25: Delete 'u2' below inside consumeUser(...), invoke completion;
//   should suggest user, u1, u2 by expected type
consumeUser(u2);

void tc26;
