/**
 * Advanced context-aware completion scenarios.
 *
 * Covers: TC-102, TC-103, TC-104, TC-106, TC-109, TC-110, TC-111, TC-112, TC-113.
 */
package completion.advanced

import completion.model.User
import completion.model.buildUser

sealed interface SyncState {
    data object Idle : SyncState
    data class Running(val progress: Int) : SyncState
    data class Failed(val message: String) : SyncState
}

class SettingsDsl {
    var enabled: Boolean = false
    var profileName: String = "default"
    fun save(): String = "saved:$profileName"
}

fun configureSettings(block: SettingsDsl.() -> Unit): SettingsDsl {
    val dsl = SettingsDsl()
    dsl.block()
    return dsl
}

fun smartCastAfterIs(x: Any) {
    if (x is User) {
        // <caret> TC-102: Delete 'name' below after "x.", invoke completion;
        //   expect User members because of smart-cast in is-check
        val smartCastName = x.name
        println(smartCastName)
    }
}

fun lambdaItTypedCompletion() {
    val users = listOf(User("Ann", 21), User("Bob", 22))
    // <caret> TC-103: Delete 'name' below after "it." in lambda; invoke completion;
    //   expect User members for implicit lambda parameter
    val names = users.map { it.name }
    println(names)
}

fun completionInsideStringTemplate(user: User) {
    // <caret> TC-104: Delete 'name' below inside string template "${...}", invoke completion;
    //   expect members of user as in normal expression context
    val text = "User=${user.name}"
    println(text)
}

fun sealedWhenBranches(state: SyncState): String {
    // <caret> TC-106: Inside when(state), invoke completion for missing branches;
    //   expect Idle, Running, Failed (and no redundant else when exhaustive)
    return when (state) {
        SyncState.Idle -> "idle"
        is SyncState.Running -> "running:${state.progress}"
        is SyncState.Failed -> "failed:${state.message}"
    }
}

fun <T : CharSequence> boundedEcho(value: T): Int = value.length

fun genericBoundContext() {
    val text = "abc"
    val number = 42
    // <caret> TC-109: Delete 'text' below in boundedEcho(...), invoke completion;
    //   expect values compatible with T : CharSequence, and reject Int-only options
    val boundedLength = boundedEcho(text)
    println(number + boundedLength)
}

fun safeCallElvisExpectedType(u: User?): String {
    val fallback = "unknown"
    // <caret> TC-110: Delete 'fallback' below after Elvis (?:), invoke completion;
    //   expect String-compatible expressions by expected type
    return u?.name ?: fallback
}

fun dslReceiverContext() {
    val result = configureSettings {
        // <caret> TC-111: Delete 'profileName' below and invoke completion inside DSL block;
        //   expect receiver members enabled, profileName, save()
        profileName = "qa"
        enabled = true
    }
    println(result.save())
}

fun crossPackageDependencyLikeModule() {
    val fromDependency = buildUser("Ann", 21)
    // <caret> TC-112: Delete 'name' below after "fromDependency.", invoke completion;
    //   expect symbols resolved from dependency package/module-like boundary
    val dependencyName = fromDependency.name
    println(dependencyName)
}

fun renamedApiTarget() = "newApi"

fun refactorMoveRenameFollowup() {
    // <caret> TC-113: After Rename/Move of renamedApiTarget across files,
    //   delete call below and invoke completion; expect updated symbol only.
    val renamedCallResult = renamedApiTarget()
    println(renamedCallResult)
}
