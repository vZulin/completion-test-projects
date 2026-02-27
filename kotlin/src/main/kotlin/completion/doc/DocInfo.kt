/**
 * Documentation, type info, and deprecated marker scenarios.
 *
 * Covers: TC-71, TC-72, TC-73, TC-83.
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

/** TC-71: QuickDoc — invoke Ctrl+Q on a completion item. */
fun quickDocScenario() {
    val user = User("Ann", 21)
    // <caret> TC-71: Place caret after "user.", select 'name' from completion,
    //   press Ctrl+Q; expect KDoc/type info for 'name: String'
    val quickDocName = user.name
}

/** TC-72: Type/signature display in completion popup. */
fun typeSignatureDisplay() {
    val user = User("Ann", 21)
    // <caret> TC-72: Place caret after "user." and invoke completion;
    //   verify each item shows return type (e.g., 'name: String', 'age: Int')
    val typedAge = user.age
}

/** TC-73: Deprecated function marker in completion list. */
@Suppress("DEPRECATION")
fun deprecatedMarker() {
    // <caret> TC-73: Delete 'oldFun()' below, type 'old'; expect 'oldFun' shown
    //   with strikethrough / deprecated styling
    val deprecatedCall = oldFun()
}

/**
 * TC-83: KDoc tag completion inside documentation comment.
 *
 * To test: place caret inside a KDoc block and type '@'.
 * Expect: @param, @return, @throws, @see, etc.
 */

/**
 * // <caret> TC-83: Place caret on this line, type '@';
 *   expect @param, @return, @throws, @see
 */
fun documentedFunction(a: Int, b: String): Boolean = a > 0
