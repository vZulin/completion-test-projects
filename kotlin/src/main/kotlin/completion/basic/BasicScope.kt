/**
 * Basic completion scope and ranking scenarios.
 *
 * Covers: TC-15, TC-16, TC-17, TC-18, TC-19, TC-20, TC-58, TC-59, TC-60, TC-21.
 *
 * Each function is self-contained. Place caret at positions described
 * in `// <caret>` comments and invoke completion to verify behavior.
 * The code shows the "after completion" state — to test, delete the
 * identifier and re-type the prefix mentioned in the comment.
 */
package completion.basic

import completion.model.User

val globalUser = User("Global", 0)

/** TC-15: Local variable completion by prefix. */
fun localVariableByPrefix() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    // <caret> TC-15: Delete 'userName' below, type 'user' and invoke completion;
    //   expect userName, userAge, user in list
    println(userName)
}

/** TC-16: Parameter completion. */
fun parameterCompletion(userName: String, userAge: Int) {
    // <caret> TC-16: Delete 'userName' below, type 'user' and invoke completion;
    //   expect userName, userAge from params
    println(userName)
}

/** TC-17: Nested scope — outer variable visible inside inner block. */
fun nestedScope() {
    val outerValue = 10
    if (true) {
        // <caret> TC-17: Delete 'outerValue' below, type 'out' inside if-block;
        //   expect outerValue from enclosing scope
        println(outerValue)
    }
}

/**
 * TC-18: Out-of-scope visibility (negative).
 * Primary test is in Java (private field across classes).
 * In Kotlin, verify that local variables from a different function are NOT suggested.
 */
fun outOfScopeSource() {
    @Suppress("unused")
    val secretLocal = "hidden"
}

fun outOfScopeTarget() {
    // <caret> TC-18: Type 'secret'; secretLocal from outOfScopeSource should NOT appear
    println("no secret here")
}

/** TC-19: Keyword completion — 'return'. */
fun keywordReturn(): Int {
    val result = 42
    // <caret> TC-19: Delete 'return' below, type 'ret'; expect 'return' keyword in suggestions
    return result
}

/**
 * TC-20: Keyword completion — 'class' at top level.
 * To test: create a new .kt file, type 'cla' at top level; expect 'class' keyword.
 */
// <caret> TC-20: In a new file, type 'cla' at top level; expect 'class' keyword

/** TC-58: Ranking — local variable should rank above global with same prefix. */
fun rankingLocalVsGlobal() {
    val globalUserName = "LocalShadow"
    // <caret> TC-58: Delete 'globalUserName' below, type 'glo';
    //   expect globalUserName (local) ranked above globalUser (global)
    println(globalUserName)
}

/** TC-59: Ranking — exact prefix match should rank above fuzzy match. */
fun rankingExactVsFuzzy() {
    val userName = "Ann"
    @Suppress("unused")
    val usernameLegacy = "legacy"
    // <caret> TC-59: Delete 'userName' below, type 'userN';
    //   expect userName (exact case) above usernameLegacy (fuzzy)
    println(userName)
}

/**
 * TC-60: Most Recently Used (MRU) ranking.
 * Manual test: select 'userName' from completion, then re-invoke completion.
 * Expect 'userName' to appear higher in the list on second invocation.
 */
fun rankingMRU() {
    val userName = "Ann"
    @Suppress("unused")
    val userAge = 21
    @Suppress("unused")
    val userEmail = "ann@test.com"
    // <caret> TC-60: Delete 'userName' below, type 'user', select userName, then
    //   re-invoke completion — previously selected item should rank higher
    println(userName)
}

/** TC-21: Negative context — caret inside number literal. */
fun negativeContextNumberLiteral() {
    @Suppress("unused")
    val x = 1234
    // <caret> TC-21: Place caret inside '1234' (between '12' and '34');
    //   invoke completion — should not crash, may show no results
}
