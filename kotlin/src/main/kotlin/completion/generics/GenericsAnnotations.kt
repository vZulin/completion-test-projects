/**
 * Generics and annotation completion scenarios.
 *
 * Covers: TC-57, TC-74, TC-76, TC-78, TC-80.
 */
package completion.generics

import java.util.ArrayList

annotation class TestAnnotation(val value: String, val enabled: Boolean = true)

// <caret> TC-57: Place caret on empty line before declaration below, type '@',
//   verify annotation auto-popup appears.
@Suppress("unused")
class AnnotationAutoPopupTarget

/** TC-74: Generic type parameter completion inside angle brackets. */
fun genericTypeParam() {
    // <caret> TC-74: Delete 'String' below inside '<>', place caret and invoke completion;
    //   expect type suggestions (String, Int, User, etc.)
    val list: ArrayList<String> = arrayListOf()
    println(list)
}

/**
 * TC-76: Auto-import inside generic type parameter.
 * Type a class name not yet imported inside angle brackets.
 */
fun autoImportInGenerics() {
    // <caret> TC-76: Delete 'Int' below inside nested '<>', type unimported class name;
    //   expect auto-import on accept
    val map: HashMap<String, ArrayList<Int>> = hashMapOf()
    println(map)
}

/**
 * TC-78: Annotation completion by prefix.
 * To test: create a new class/function, type '@Depre' before it;
 * expect @Deprecated in suggestions.
 */
// <caret> TC-78: Type '@Depre' before this function; expect @Deprecated in suggestions
@Suppress("unused")
fun annotationTarget() {}

/** TC-80: Annotation parameter completion. */
@TestAnnotation(
    // <caret> TC-80: Delete params below, place caret inside annotation parens;
    //   expect 'value =' and 'enabled =' suggestions
    value = "test", enabled = false
)
fun annotatedFunction() {}
