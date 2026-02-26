package completion.advanced.external

class RemoteSettings(val id: String)

fun RemoteSettings.syncNow(): String = "sync:$id"
