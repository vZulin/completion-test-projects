// Test cases covered: TC-11, TC-62
// Source: EX-TS-2
//
// TC-11: import path completion — typing "./" should suggest sibling modules
// TC-62: relative path variants — test "./" and "../" paths

import x from "./" // <caret> TC-11: import path completion — should list model, basicCombo, etc.
