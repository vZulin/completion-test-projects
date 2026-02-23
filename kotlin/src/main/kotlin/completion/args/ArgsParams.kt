/**
 * Argument / parameter completion scenarios: type matching, parameter info, named args.
 *
 * Covers: TC-39, TC-40, TC-41, TC-42, TC-43, TC-44, TC-45.
 *
 * To test: delete the argument(s) inside the function call, place caret,
 * and invoke completion.
 */
package completion.args

import completion.model.User
import completion.model.buildUser
import completion.model.consumeUser

/** TC-39: Argument substitution — first parameter (String). */
fun argSubstitutionFirst() {
    val name = "Ann"
    @Suppress("unused")
    val age = 21
    // <caret> TC-39: Delete 'name' below (first arg), place caret inside buildUser();
    //   expect 'name' (String) suggested by type
    buildUser(name, 21)
}

/** TC-40: Argument substitution — second parameter after comma (Int). */
fun argSubstitutionSecond() {
    @Suppress("unused")
    val name = "Ann"
    val age = 21
    // <caret> TC-40: Delete 'age' below (second arg after comma), place caret;
    //   expect 'age' (Int) suggested by type
    buildUser("Ann", age)
}

/** TC-41: Argument by expected type for consumer function. */
fun argByExpectedType() {
    val u = User("Ann", 21)
    // <caret> TC-41: Delete 'u' below inside consumeUser(), invoke completion;
    //   expect 'u' matching type User
    consumeUser(u)
}

/**
 * TC-42: Parameter info integration.
 * After typing `buildUser(`, parameter info popup should display:
 * `name: String, age: Int`
 */
fun parameterInfo() {
    // <caret> TC-42: Place caret after '(' in buildUser(; parameter info should show
    //   'name: String, age: Int'
    buildUser("Ann", 21)
}

/** TC-43: Named arguments completion. */
fun namedArguments() {
    // <caret> TC-43: Delete named args below, place caret inside buildUser();
    //   expect 'name =' and 'age =' as named arg suggestions
    buildUser(name = "Ann", age = 21)
}

/** TC-44: Already used named argument should not be re-suggested. */
fun namedArgNotDuplicated() {
    // <caret> TC-44: Delete 'age = 21' below, place caret after comma;
    //   expect 'age =' but NOT 'name =' again
    buildUser(name = "Ann", age = 21)
}

/** TC-45: Named arguments — partial prefix filtering. */
fun namedArgFiltering() {
    // <caret> TC-45: Delete 'name' below, type 'na' inside buildUser();
    //   expect 'name =' filtered from named args
    buildUser(name = "Ann", age = 21)
}
