/**
 * Auto-import for extension functions from another package.
 *
 * Covers: TC-107.
 */
package completion.advanced

import completion.advanced.external.RemoteSettings
import completion.advanced.external.syncNow

fun autoImportExtensionFunction() {
    val settings = RemoteSettings("machine-A")
    // <caret> TC-107: Delete import `syncNow` above and delete 'syncNow()' below;
    //   type 'sync' after settings., invoke completion; expect extension function with auto-import.
    val syncResult = settings.syncNow()
    println(syncResult)
}
