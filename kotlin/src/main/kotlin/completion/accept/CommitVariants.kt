/**
 * Additional commit-character scenarios.
 *
 * Covers: TC-121.
 */
package completion.accept

fun commitBySemicolonAndSpace() {
    val text = "Ann"
    // <caret> TC-121-A: Delete 'uppercase()' below so only 'text.up' remains,
    //   then press ';' to commit selected completion item and insert semicolon.
    val a = text.uppercase()

    val num = 42
    // <caret> TC-121-B: Delete 'toString()' below so only 'num.toS' remains,
    //   then press Space to commit completion (if enabled in IDE settings).
    val b = num.toString()

    println(a + b)
}
