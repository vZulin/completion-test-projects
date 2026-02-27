/**
 * Member access completion scenarios: dot access, companion, nullable, extensions.
 *
 * Covers: TC-22, TC-23, TC-25, TC-26, TC-61, TC-62.
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

/** TC-22: Dot access — properties. */
fun dotAccessProperties() {
    val user = User("Ann", 21)
    // <caret> TC-22: Delete 'name' below, place caret after "user." and invoke completion;
    //   expect 'name', 'age' properties in list
    val selectedName = user.name
}

/** TC-23: Dot access — prefix filtering after dot. */
fun dotAccessMethods() {
    val user = User("Ann", 21)
    // <caret> TC-23: Delete 'name' below, type 'na' after "user.", invoke completion;
    //   expect list filtered to 'name' (or language-specific getter variant)
    val filteredName = user.name
}

/** TC-25: Field vs method — field selection. */
fun fieldSelection() {
    val user = User("Ann", 21)
    // <caret> TC-25: Delete 'name' below, type 'na' after "user.";
    //   expect 'name' (property) at top, not a method
    val fieldName = user.name
}

/** TC-26: Companion object members via class name. */
fun companionAccess() {
    // <caret> TC-26: Delete 'make()' below, place caret after "Util." and invoke completion;
    //   expect make(), VERSION from companion
    val utilInstance = Util.make()
}

/** TC-61: Nullable safe-call completion. */
fun nullableSafeCall() {
    val u: User? = null
    // <caret> TC-61: Delete 'name' below, place caret after "u?." and invoke completion;
    //   expect name, age, toString(), etc.
    val safeName = u?.name
}

/** TC-62: Extension function visible through dot. */
fun extensionFunction() {
    val s = "Ann"
    // <caret> TC-62: Delete 'extHello()' below, place caret after "s." and invoke completion;
    //   expect extHello() alongside standard String members
    val extensionGreeting = s.extHello()
}
