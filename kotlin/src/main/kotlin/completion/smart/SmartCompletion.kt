/**
 * Smart (type-aware) completion scenarios.
 *
 * Covers: TC-34, TC-35, TC-36, TC-37, TC-38.
 *
 * To test: delete the RHS expression, place caret after '=' or 'return',
 * then invoke smart completion (Ctrl+Shift+Space).
 */
package completion.smart

import completion.model.User
import completion.model.buildUser
import completion.model.consumeUser

/** TC-34: Smart completion on assignment — variable matching expected type. */
fun assignmentVariable() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    // <caret> TC-34: Delete 'user' below after '=', invoke smart completion;
    //   expect 'user' variable matching type User
    val u1: User = user
}

/** TC-35: Smart completion on assignment — constructor / factory. */
fun assignmentFactory() {
    val userName = "Ann"
    val userAge = 21

    // <caret> TC-35: Delete expression below after '=', invoke smart completion;
    //   expect buildUser(...), User(...) constructor
    val u2: User = buildUser(userName, userAge)
}

/** TC-36: Smart completion in return position. */
fun returnCompletion(): User {
    val existing = User("Ann", 21)
    // <caret> TC-36: Delete 'existing' below after 'return', invoke smart completion;
    //   expect existing, buildUser(), User(...)
    return existing
}

/** TC-37: Smart completion for function argument by expected type. */
fun argsByExpectedType() {
    val u = User("Ann", 21)
    // <caret> TC-37: Delete 'u' below inside consumeUser(), invoke smart completion;
    //   expect 'u' matching type User
    consumeUser(u)
}

/**
 * TC-38: Smart completion — no matching type variable.
 * Verify that smart completion shows constructors/factory even when no local matches the type.
 */
fun noLocalMatch() {
    val name = "Ann"
    val age = 21
    // <caret> TC-38: Delete expression below after '=', invoke smart completion;
    //   no local User variable; expect User(...), buildUser(...) constructors
    val target: User = User(name, age)
}
