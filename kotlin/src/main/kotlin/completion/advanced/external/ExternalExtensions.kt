package completion.advanced.external

class RemoteSettings(val id: String)

// <caret> TC-107: Declaration source for extension auto-import scenario.
fun RemoteSettings.syncNow(): String = "sync:$id"
