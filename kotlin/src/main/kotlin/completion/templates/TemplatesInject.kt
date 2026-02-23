/**
 * Postfix templates, language injection, and injected-to-host transition scenarios.
 *
 * Covers: TC-81, TC-82, TC-84, TC-85, TC-86.
 */
package completion.templates

/**
 * TC-81: Postfix completion `.if` — wraps boolean expression in if-block.
 * TC-82: Postfix completion result verification.
 *
 * To test: delete the if-block below, type `ok.if` and accept the postfix template.
 * Expected result: `if (ok) { }` generated automatically.
 */
fun postfixIf() {
    val ok = true
    // <caret> TC-81: Delete the if-block below, type 'ok.if', accept postfix;
    //   expect: if (ok) { }
    // <caret> TC-82: Verify the generated code is syntactically correct
    if (ok) {
        println("postfix expansion result")
    }
}

/** TC-84: SQL language injection — completion inside SQL string. */
fun sqlInjection() {
    // <caret> TC-84: Place caret after 'SELECT ' inside the string;
    //   invoke completion — expect SQL column/keyword suggestions if injection is active
    @Suppress("unused")
    // language=SQL
    val query = "SELECT id FROM users WHERE id = 1"
}

/** TC-85: Regex language injection — completion inside regex string. */
fun regexInjection() {
    // <caret> TC-85: Place caret inside regex string;
    //   invoke completion — expect regex tokens if injection is active
    @Suppress("unused")
    // language=RegExp
    val pattern = "\\d+"
}

/**
 * TC-86: Transition from injected language back to host (Kotlin).
 * After editing inside an injected SQL/regex fragment, move caret outside the string
 * and verify that Kotlin completion resumes normally.
 */
fun injectedToHostTransition() {
    // language=SQL
    @Suppress("unused")
    val q = "SELECT 1"
    val x = 42
    // <caret> TC-86: After editing inside SQL string above, place caret after "x."
    //   below; expect normal Kotlin Int members
    val tc86 = x.toLong()
}
