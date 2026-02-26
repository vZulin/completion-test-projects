// Test cases covered: TC-10, TC-49
// Source: EX-TS-2
//
// TC-10: import path completion — typing "./" should suggest sibling modules
// TC-49: relative path variants — test "./" and "../" paths

// <caret> TC-10/TC-49: In import below, delete 'model' and place caret after "./";
//   completion should list sibling modules and relative variants.
import { buildUser } from "./model";

const tc62 = buildUser("Ann", 21);
void tc62;
