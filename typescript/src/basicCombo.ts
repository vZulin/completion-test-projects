// Test cases covered: TC-7 (TS), TC-9 (TS), TC-26, TC-27, TC-28, TC-29 (TS variants)
// Source: EX-TS-1

import { User, buildUser, consumeUser } from "./model";

const userName = "Ann";
const userAge = 21;
const user: User = buildUser(userName, userAge);

user. // <caret> TC-7: member completion + filtering (type "na" to filter to "name")

const u1: User = // <caret> TC-26: smart completion — should suggest buildUser, variables of type User

const u2: User = buildUser( // <caret> TC-27: args completion — should suggest userName; then comma for second arg (TC-28)
)

consumeUser( // <caret> TC-29: expected type completion — should suggest user, u1, u2
)
