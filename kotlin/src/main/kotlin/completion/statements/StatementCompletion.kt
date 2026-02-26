/**
 * Statement-level completion scenarios: if-else, when.
 *
 * Covers: TC-92 (if-else completion), TC-93 (when expression/statement).
 */
package completion.statements

import completion.model.User

enum class Status { ACTIVE, INACTIVE, BANNED }

/** TC-92: if-else statement completion. */
fun checkAge(user: User) {
    // Scenario A: type 'if' keyword and complete the statement structure
    // <caret> TC-92-A: type 'if' here, invoke completion — expect 'if' keyword/postfix
    if (user.age > 18) {
        println("adult")
    }
    // <caret> TC-92-B: place caret here after '}', type 'el' — expect 'else' / 'else if'
    else {
        println("minor")
    }
}

/** TC-92 continued: if-else if-else chain completion. */
fun classifyAge(age: Int): String {
    // <caret> TC-92-C: after first '}', type 'else' — expect 'else if' and 'else'
    if (age < 13) {
        return "child"
    } else if (age < 18) {
        return "teen"
    } else {
        return "adult"
    }
}

/** TC-93: when expression/statement completion. */
fun describeStatus(status: Status): String {
    // <caret> TC-93-A: type 'wh' inside function body, invoke completion — expect 'when' keyword
    // <caret> TC-93-B: place caret inside 'when' block after last '->' — expect remaining enum entries
    return when (status) {
        Status.ACTIVE -> "active user"
        Status.INACTIVE -> "inactive user"
        Status.BANNED -> "banned user"
    }
}

/** TC-93 continued: when without argument — boolean branches. */
fun describeAge(age: Int): String {
    // <caret> TC-93-C: inside when{} block, invoke completion for branch conditions
    return when {
        age < 13 -> "child"
        age < 18 -> "teen"
        else -> "adult"
    }
}
