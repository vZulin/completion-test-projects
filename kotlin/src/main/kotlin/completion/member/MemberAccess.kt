/**
 * Member access completion scenarios: dot access, companion, nullable, extensions.
 *
 * Covers: TC-26, TC-27, TC-29, TC-30, TC-32, TC-33.
 *
 * To test: place caret after the dot (delete the member), invoke completion.
 */
package completion.member

import completion.model.User

class Util {
    companion object {
        fun make(): Util = Util()
        val VERSION = "1.0"
    }
}

fun String.extHello(): String = "Hello, $this"

/** TC-26: Dot access — properties. */
fun dotAccessProperties() {
    val user = User("Ann", 21)
    // <caret> TC-26: Delete 'name' below, place caret after "user." and invoke completion;
    //   expect 'name', 'age' properties in list
    val tc26 = user.name
}

/** TC-27: Dot access — methods. */
fun dotAccessMethods() {
    val user = User("Ann", 21)
    // <caret> TC-27: Delete 'toString()' below, place caret after "user." and invoke completion;
    //   expect toString(), hashCode(), copy(), etc.
    val tc27 = user.toString()
}

/** TC-29: Field vs method — field selection. */
fun fieldSelection() {
    val user = User("Ann", 21)
    // <caret> TC-29: Delete 'name' below, type 'na' after "user.";
    //   expect 'name' (property) at top, not a method
    val tc29 = user.name
}

/** TC-30: Companion object members via class name. */
fun companionAccess() {
    // <caret> TC-30: Delete 'make()' below, place caret after "Util." and invoke completion;
    //   expect make(), VERSION from companion
    val tc30 = Util.make()
}

/** TC-32: Nullable safe-call completion. */
fun nullableSafeCall() {
    val u: User? = null
    // <caret> TC-32: Delete 'name' below, place caret after "u?." and invoke completion;
    //   expect name, age, toString(), etc.
    val tc32 = u?.name
}

/** TC-33: Extension function visible through dot. */
fun extensionFunction() {
    val s = "Ann"
    // <caret> TC-33: Delete 'extHello()' below, place caret after "s." and invoke completion;
    //   expect extHello() alongside standard String members
    val tc33 = s.extHello()
}
