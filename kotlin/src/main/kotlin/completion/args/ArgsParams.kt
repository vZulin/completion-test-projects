/**
 * Argument / parameter completion scenarios: type matching, parameter info, named args.
 *
 * Covers: TC-33, TC-34, TC-35, TC-63, TC-64, TC-65, TC-66.
 *
 * To test: delete the argument(s) inside the function call, place caret,
 * and invoke completion.
 */
package completion.args

import completion.model.User
import completion.model.buildUser

/** TC-33: Argument substitution — first parameter (String). */
fun argSubstitutionFirst() {
    val name = "Ann"
    @Suppress("unused")
    val age = 21
    // <caret> TC-33: Delete 'name' below (first arg), place caret inside buildUser();
    //   expect 'name' (String) suggested by type
    buildUser(name, 21)
}

/** TC-34: Argument substitution — second parameter after comma (Int). */
fun argSubstitutionSecond() {
    @Suppress("unused")
    val name = "Ann"
    val age = 21
    // <caret> TC-34: Delete 'age' below (second arg after comma), place caret;
    //   expect 'age' (Int) suggested by type
    buildUser("Ann", age)
}

/** TC-35: Completion immediately after opening parenthesis. */
fun argByExpectedType() {
    val userName = "Ann"
    val userAge = 21
    // <caret> TC-35: Delete both args below, place caret right after '(' in buildUser( ... );
    //   invoke completion with empty prefix; expect non-empty context-aware suggestions
    buildUser(userName, userAge)
}

/**
 * TC-63: Parameter info integration.
 * After typing `buildUser(`, parameter info popup should display:
 * `name: String, age: Int`
 */
fun parameterInfo() {
    // <caret> TC-63: Place caret after '(' in buildUser(; parameter info should show
    //   'name: String, age: Int'
    buildUser("Ann", 21)
}

/** TC-64: Named arguments completion. */
fun namedArguments() {
    // <caret> TC-64: Delete named args below, place caret inside buildUser();
    //   expect 'name =' and 'age =' as named arg suggestions
    buildUser(name = "Ann", age = 21)
}

/** TC-65: Already used named argument should not be re-suggested. */
fun namedArgNotDuplicated() {
    // <caret> TC-65: Delete 'age = 21' below, place caret after comma;
    //   expect 'age =' but NOT 'name =' again
    buildUser(name = "Ann", age = 21)
}

/** TC-66: Named arguments — partial prefix filtering. */
fun namedArgFiltering() {
    // <caret> TC-66: Delete 'name' below, type 'na' inside buildUser();
    //   expect 'name =' filtered from named args
    buildUser(name = "Ann", age = 21)
}
