/**
 * Refactoring-aware completion scenarios: rename and change-signature smoke tests.
 *
 * Covers: TC-87, TC-88.
 *
 * These are manual/semi-automated tests:
 * 1. Perform the refactoring (Rename / Change Signature).
 * 2. Invoke completion at the marked caret position.
 * 3. Verify that completion reflects the refactored state.
 */
package completion.refactoring

fun renamedTarget() {}

/** TC-87: Rename refactoring — completion should reflect the new name. */
fun renameScenario() {
    // <caret> TC-87: First, rename 'renamedTarget' to 'newTargetName' via Refactor > Rename.
    //   Then delete 'renamedTarget()' below, type 'new' and invoke completion;
    //   expect 'newTargetName' in suggestions
    renamedTarget()
}

fun foo(a: Int) {}

/**
 * TC-88: Change Signature refactoring — completion should reflect new parameter list.
 * Steps: use Change Signature to add `b: String` parameter to `foo`,
 * then invoke completion inside the call.
 */
fun changeSignatureScenario() {
    // <caret> TC-88: After adding 'b: String' to foo() via Change Signature,
    //   place caret inside foo(); expect parameter info shows 'a: Int, b: String'
    foo(42)
}
