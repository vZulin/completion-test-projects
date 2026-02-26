/**
 * String literal and path completion scenarios.
 *
 * Covers: TC-47, TC-48, TC-49, TC-70.
 */
package completion.strings

import java.io.File

/** TC-47: Path completion inside File("") — absolute path root. */
fun pathCompletionAbsolute() {
    val f = File("/") // <caret> TC-47: Place caret after '/'; invoke completion — expect filesystem root entries
}

/** TC-48: Path completion inside File("") — nested directory. */
fun pathCompletionNested() {
    val f = File("/tmp/") // <caret> TC-48: Place caret after '/tmp/'; invoke completion — expect contents of /tmp
}

/**
 * TC-49: Relative path completion.
 * Primary test is in TypeScript (import path).
 * In Kotlin, relative paths inside File("") may also trigger path completion.
 */
fun pathCompletionRelative() {
    val f = File("src/") // <caret> TC-49: Place caret after 'src/'; invoke completion — expect project-relative paths
}

/** TC-70: Plain string — no path completion expected. */
fun plainStringNegative() {
    val s = "hello world" // <caret> TC-70: Place caret inside 'hello world'; invoke completion — should not offer file paths
}
