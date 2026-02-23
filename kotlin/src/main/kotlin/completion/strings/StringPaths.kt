/**
 * String literal and path completion scenarios.
 *
 * Covers: TC-60, TC-61, TC-62, TC-63.
 */
package completion.strings

import java.io.File

/** TC-60: Path completion inside File("") — absolute path root. */
fun pathCompletionAbsolute() {
    val f = File("/") // <caret> TC-60: Place caret after '/'; invoke completion — expect filesystem root entries
}

/** TC-61: Path completion inside File("") — nested directory. */
fun pathCompletionNested() {
    val f = File("/tmp/") // <caret> TC-61: Place caret after '/tmp/'; invoke completion — expect contents of /tmp
}

/**
 * TC-62: Relative path completion.
 * Primary test is in TypeScript (import path).
 * In Kotlin, relative paths inside File("") may also trigger path completion.
 */
fun pathCompletionRelative() {
    val f = File("src/") // <caret> TC-62: Place caret after 'src/'; invoke completion — expect project-relative paths
}

/** TC-63: Plain string — no path completion expected. */
fun plainStringNegative() {
    val s = "hello world" // <caret> TC-63: Place caret inside 'hello world'; invoke completion — should not offer file paths
}
