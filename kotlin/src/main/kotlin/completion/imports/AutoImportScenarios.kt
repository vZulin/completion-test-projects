/**
 * Auto-import completion scenarios.
 *
 * Covers: TC-55, TC-57, TC-58.
 *
 * IMPORTANT: This file intentionally uses unresolved references (File, Date)
 * to test auto-import behavior. It will NOT compile as-is — this is by design.
 * The IDE should offer to auto-import the correct class upon completion accept.
 *
 * To make the project compile, comment out the function bodies or
 * use `@file:Suppress("UNRESOLVED_REFERENCE")` — but during testing,
 * leave them as-is so auto-import kicks in.
 */
@file:Suppress("unused")

package completion.imports

// NO imports for File or Date — auto-import is the subject under test

/** TC-55: Auto-import java.io.File on completion accept. */
fun autoImportFile() {
    // <caret> TC-55: Type 'File' (no import present), accept from completion;
    //   verify import java.io.File is added automatically
    // val f = File("")
}

/**
 * TC-57: Import conflict — multiple Date classes available.
 * Both java.util.Date and java.sql.Date exist in classpath.
 * Completion should display both with qualified package names.
 */
fun importConflictDate() {
    // <caret> TC-57: Type 'Date' (no import present);
    //   expect both java.util.Date and java.sql.Date in completion with qualifiers
    // val d = Date()
}

/**
 * TC-58: Import conflict resolution — user selects specific variant.
 * After selecting java.sql.Date, verify the correct import is added.
 */
fun importConflictResolution() {
    // <caret> TC-58: Type 'Date', select java.sql.Date from completion;
    //   verify import java.sql.Date added
    // val d2 = Date()
}
