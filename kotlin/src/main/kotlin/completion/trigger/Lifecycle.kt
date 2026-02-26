/**
 * Completion trigger, auto-popup, and race/cancel scenarios.
 *
 * Covers: TC-1, TC-2, TC-3, TC-4, TC-5, TC-6, TC-7, TC-8, TC-9, TC-11, TC-12, TC-13, TC-14.
 *
 * How to use:
 * - Place the caret at positions described in `// <caret>` comments.
 * - Invoke completion (Ctrl+Space) or wait for auto-popup where noted.
 * - The compilable code below shows the "after completion" state.
 *   To test, delete the part after the dot/paren and re-type.
 */
package completion.trigger

import completion.model.User
import completion.model.buildUser
import java.io.File

fun triggerScenarios() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    // --- Explicit invocation (Ctrl+Space) ---

    // <caret> TC-1: Place caret after "user." and invoke basic completion; expect name, age, toString, etc.
    val tc1 = user.name

    // <caret> TC-2: Delete 'user' below after '=', place caret after "val u1: User =",
    //   then invoke smart completion; expect user/User(...)/buildUser(...)
    val u1: User = user

    // <caret> TC-3: Place caret after "user.", type 'na' then Backspace to 'n'; list should re-expand
    val tc3 = user.name

    // <caret> TC-4: Place caret after "user.", type 'xyz'; expect empty list / "No suggestions"
    val tc4 = user.name

    // <caret> TC-5: Invoke completion after "user.", press Escape to dismiss popup; verify it closes
    val tc5 = user.name

    // <caret> TC-6: Press Escape then re-invoke Ctrl+Space; popup reappears
    val tc6 = user.name

    // --- Auto-popup triggers ---

    // <caret> TC-7: Place caret after "user" (no dot), then type '.'; auto-popup should appear
    val tc7 = user.name

    // <caret> TC-8: Place caret after '(' in buildUser(; auto-popup / parameter info should appear
    val tc8 = buildUser(userName, userAge)

    // <caret> TC-9: Place caret after ',' in buildUser(userName,; auto-popup for next arg
    val tc9 = buildUser(userName, userAge)

    // --- Path completion ---

    // <caret> TC-11: Place caret inside quotes in File(""); invoke completion for file paths
    val tc12 = File("")

    // --- Race / cancel / fast-typing scenarios (manual steps) ---

    // <caret> TC-12: Fast typing — type "user." quickly followed by "na" before popup renders.
    //        Verify: popup appears with filtered results containing "name".
    val tc13 = user.name

    // <caret> TC-13: Backspace cancel — invoke completion on "user.n", then press Backspace.
    //        Verify: popup either updates or closes without errors.
    val tc14 = user.name

    // <caret> TC-14: Caret movement — invoke completion, then press Left/Right arrow.
    //        Verify: popup dismisses; no residual UI artifacts.
    val tc15 = user.name
}
