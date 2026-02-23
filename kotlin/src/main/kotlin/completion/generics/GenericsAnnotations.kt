/**
 * Generics and annotation completion scenarios.
 *
 * Covers: TC-67, TC-69, TC-71, TC-73.
 */
package completion.generics

import java.util.ArrayList

annotation class TestAnnotation(val value: String, val enabled: Boolean = true)

/** TC-67: Generic type parameter completion inside angle brackets. */
fun genericTypeParam() {
    // <caret> TC-67: Delete 'String' below inside '<>', place caret and invoke completion;
    //   expect type suggestions (String, Int, User, etc.)
    val list: ArrayList<String> = arrayListOf()
    println(list)
}

/**
 * TC-69: Auto-import inside generic type parameter.
 * Type a class name not yet imported inside angle brackets.
 */
fun autoImportInGenerics() {
    // <caret> TC-69: Delete 'Int' below inside nested '<>', type unimported class name;
    //   expect auto-import on accept
    val map: HashMap<String, ArrayList<Int>> = hashMapOf()
    println(map)
}

/**
 * TC-71: Annotation completion by prefix.
 * To test: create a new class/function, type '@Depre' before it;
 * expect @Deprecated in suggestions.
 */
// <caret> TC-71: Type '@Depre' before this function; expect @Deprecated in suggestions
@Suppress("unused")
fun annotationTarget() {}

/** TC-73: Annotation parameter completion. */
@TestAnnotation(
    // <caret> TC-73: Delete params below, place caret inside annotation parens;
    //   expect 'value =' and 'enabled =' suggestions
    value = "test", enabled = false
)
fun annotatedFunction() {}
