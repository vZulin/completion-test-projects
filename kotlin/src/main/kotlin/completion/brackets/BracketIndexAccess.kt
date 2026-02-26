/**
 * Bracket-index completion scenarios for [] access.
 *
 * Covers: TC-53, TC-99, TC-100.
 */
package completion.brackets

import completion.model.User

/** TC-53: Completion for Int index inside list brackets. */
fun listIndexByInt() {
    val users = listOf(User("Ann", 21), User("Bob", 22))
    val selectedIndex = 0
    // <caret> TC-53: Delete 'selectedIndex' below inside users[...], invoke completion;
    //   expect Int variables from scope (e.g., selectedIndex)
    val tc101 = users[selectedIndex]
    println(tc101.name)
}

/** TC-99: Completion in nested brackets matrix[row][col]. */
fun nestedBracketIndex() {
    val matrix = arrayOf(intArrayOf(1, 2), intArrayOf(3, 4))
    val row = 1
    val col = 0
    // <caret> TC-99: Delete 'row' in first [] and 'col' in second [] below;
    //   invoke completion in each bracket, expect Int variables row/col
    val tc102 = matrix[row][col]
    println(tc102)
}

/** TC-100: Completion for map key expression inside brackets. */
fun mapKeyInsideBrackets() {
    val userByKey = mapOf("name" to "Ann", "age" to "21")
    val nameKey = "name"
    // <caret> TC-100: Delete 'nameKey' below inside userByKey[...], invoke completion;
    //   expect String key variables from scope (e.g., nameKey)
    val tc103 = userByKey[nameKey]
    println(tc103)
}
