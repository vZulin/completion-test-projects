/**
 * Completion accept / commit character scenarios.
 *
 * Covers: TC-46, TC-47, TC-48, TC-49, TC-51, TC-52, TC-53.
 *
 * To test: delete the completed identifier, type the prefix shown
 * in the comment, then use the specified accept key.
 */
package completion.accept

import completion.model.User
import completion.model.buildUser

/** TC-46: Accept completion with Enter key — inserts item. */
fun acceptWithEnter() {
    val userName = "Ann"
    @Suppress("unused")
    val userAge = 21
    // <caret> TC-46: Delete 'userName' below, type 'user', select 'userName' from list,
    //   press Enter; expect 'userName' inserted
    println(userName)
}

/** TC-47: Accept completion with Tab key — replaces current token. */
fun acceptWithTab() {
    @Suppress("unused")
    val userName = "Ann"
    val userAge = 21
    // <caret> TC-47: Delete 'userAge' below, type 'user', select 'userAge' from list,
    //   press Tab; expect token replaced with 'userAge'
    println(userAge)
}

/** TC-48: Replace selection with completion item. */
fun replaceSelection() {
    val userName = "Ann"
    @Suppress("UNUSED_VARIABLE")
    val x = "REPLACE_ME"
    // <caret> TC-48: Select "REPLACE_ME" above, invoke completion, choose 'userName';
    //   expect replacement
    println(userName)
}

/** TC-49: Commit by dot — accept completion and immediately chain dot access. */
fun commitByDot() {
    val user = User("Ann", 21)
    // <caret> TC-49: Delete 'name.length' below, type 'na' after "user.", then press '.';
    //   expect 'name' accepted and new dot appended -> user.name.
    val tc49 = user.name.length
}

/** TC-51: Commit by comma in function arguments. */
fun commitByComma() {
    val name = "Ann"
    val age = 21
    // <caret> TC-51: Delete 'name' below (first arg), type 'na', press ',';
    //   expect 'name' accepted and comma inserted
    buildUser(name, age)
}

/** TC-52: Caret placement after selecting a function — caret inside parentheses. */
fun caretAfterFunctionSelect() {
    // <caret> TC-52: Delete 'buildUser("Ann", 21)' below, type 'build', select 'buildUser'
    //   from list, accept; expect 'buildUser()' with caret between parens
    buildUser("Ann", 21)
}

/** TC-53: Caret placement after selecting a no-arg function. */
fun caretAfterNoArgFunction() {
    val user = User("Ann", 21)
    // <caret> TC-53: Delete 'toString()' below, type 'toStr', select 'toString()';
    //   expect caret after closing paren since no args needed
    val tc53 = user.toString()
}
