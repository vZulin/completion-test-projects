/**
 * Postfix templates, language injection, and injected-to-host transition scenarios.
 *
 * Covers: TC-114, TC-115, TC-116, TC-117, TC-118, TC-119.
 */
package completion.templates

/**
 * TC-114: Postfix completion `.if` — wraps boolean expression in if-block.
 * TC-115: Conflict between basic completion and postfix suggestions.
 *
 * To test: delete the if-block below, type `ok.if` and accept the postfix template.
 * Expected result: `if (ok) { }` generated automatically.
 */
fun postfixIf() {
    val ok = true
    // <caret> TC-114: Delete the if-block below, type 'ok.if', accept postfix;
    //   expect: if (ok) { }
    // <caret> TC-115: On 'ok.if' invoke completion and verify list shows postfix `.if`
    //   together with regular completion items (if available for current context).
    if (ok) {
        println("postfix expansion result")
    }
}

/** TC-116: Live template `sout` expansion in Kotlin editor context. */
fun liveTemplateSout() {
    // <caret> TC-116-KT: Delete println below, type 'sout', press Tab;
    //   expect println(...) template expansion with caret inside parentheses
    println("template expansion result")
}

/** TC-117: SQL language injection — completion inside SQL string. */
fun sqlInjection() {
    // <caret> TC-117: Place caret after 'SELECT ' inside the string;
    //   invoke completion — expect SQL column/keyword suggestions if injection is active
    @Suppress("unused")
    // language=SQL
    val query = "SELECT id FROM users WHERE id = 1"
}

/** TC-118: Regex language injection — completion inside regex string. */
fun regexInjection() {
    // <caret> TC-118: Place caret inside regex string;
    //   invoke completion — expect regex tokens if injection is active
    @Suppress("unused")
    // language=RegExp
    val pattern = "\\d+"
}

/**
 * TC-119: Transition from injected language back to host (Kotlin).
 * After editing inside an injected SQL/regex fragment, move caret outside the string
 * and verify that Kotlin completion resumes normally.
 */
fun injectedToHostTransition() {
    // language=SQL
    @Suppress("unused")
    // <caret> TC-119: Place caret after "SELECT " inside SQL string, invoke completion,
    //   then move caret to host Kotlin line below and verify popup closes/updates correctly.
    val q = "SELECT id FROM users WHERE id = 1"
    val x = 42
    val tc86 = x.toLong()
}
