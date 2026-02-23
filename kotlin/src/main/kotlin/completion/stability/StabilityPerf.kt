/**
 * Stability and performance scenarios: dumb mode, rapid invocation, large scope.
 *
 * Covers: TC-89, TC-90, TC-91, TC-93.
 */
package completion.stability

import completion.model.User

/**
 * TC-89: Dumb mode — basic completion during indexing.
 * Steps: invalidate caches (File > Invalidate Caches), reopen project,
 * and immediately invoke completion before indexing completes.
 * Expect: completion either works with partial results or shows "Indexing..." without crash.
 */
fun dumbModeBasic() {
    val user = User("Ann", 21)
    // <caret> TC-89: Place caret after "user." during indexing; invoke completion;
    //   expect graceful handling
    val tc89 = user.name
}

/**
 * TC-90: Dumb mode — smart completion during indexing.
 * Same setup as TC-89, but invoke smart completion (Ctrl+Shift+Space).
 */
fun dumbModeSmart() {
    // <caret> TC-90: Delete expression below after '=', invoke smart completion
    //   during indexing; expect no crash
    val u: User = User("Ann", 21)
    println(u)
}

/**
 * TC-91: Dumb mode — auto-popup during indexing.
 * Same setup as TC-89, but type a dot to trigger auto-popup.
 */
fun dumbModeAutoPopup() {
    val user = User("Ann", 21)
    // <caret> TC-91: Place caret after "user" (no dot), type '.' during indexing;
    //   expect no crash
    val tc91 = user.name
}

/**
 * TC-93: Performance — rapid open/close completion (30x).
 * Steps: invoke Ctrl+Space then Escape 30 times in quick succession.
 * Expect: no UI freeze, no memory leak, IDE remains responsive.
 */
fun performanceRapidInvocation() {
    val user = User("Ann", 21)
    // <caret> TC-93: Place caret after "user." and invoke/dismiss completion 30 times
    //   rapidly (Ctrl+Space, Escape); verify IDE stability
    val tc93 = user.name
}

/**
 * Performance reference: for large-scope completion (50+ methods),
 * see the Java project's PerfMany class (EX-JV-17).
 * That class contains 50+ generated methods to stress-test completion performance.
 */
