/**
 * Documentation, type info, and deprecated marker scenarios.
 *
 * Covers: TC-64, TC-65, TC-66, TC-76.
 */
package completion.doc

import completion.model.User

@Deprecated("Use newFun() instead")
fun oldFun(): String = "legacy"

fun newFun(): String = "modern"

/**
 * Processes user data.
 *
 * @param user the user to process
 * @param verbose whether to enable verbose logging
 * @return processed user name
 */
fun processUser(user: User, verbose: Boolean): String = user.name

/** TC-64: QuickDoc — invoke Ctrl+Q on a completion item. */
fun quickDocScenario() {
    val user = User("Ann", 21)
    // <caret> TC-64: Place caret after "user.", select 'name' from completion,
    //   press Ctrl+Q; expect KDoc/type info for 'name: String'
    val tc64 = user.name
}

/** TC-65: Type/signature display in completion popup. */
fun typeSignatureDisplay() {
    val user = User("Ann", 21)
    // <caret> TC-65: Place caret after "user." and invoke completion;
    //   verify each item shows return type (e.g., 'name: String', 'age: Int')
    val tc65 = user.age
}

/** TC-66: Deprecated function marker in completion list. */
@Suppress("DEPRECATION")
fun deprecatedMarker() {
    // <caret> TC-66: Delete 'oldFun()' below, type 'old'; expect 'oldFun' shown
    //   with strikethrough / deprecated styling
    val tc66 = oldFun()
}

/**
 * TC-76: KDoc tag completion inside documentation comment.
 *
 * To test: place caret inside a KDoc block and type '@'.
 * Expect: @param, @return, @throws, @see, etc.
 */

/**
 * // <caret> TC-76: Place caret on this line, type '@';
 *   expect @param, @return, @throws, @see
 */
fun documentedFunction(a: Int, b: String): Boolean = a > 0
