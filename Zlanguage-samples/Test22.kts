data class BuildTarget(val name: String, val enabled: Boolean)

val targets = listOf(
    BuildTarget("IntelliJ IDEA", enabled = true),
    BuildTarget("Kotlin DSL", enabled = true),
    BuildTarget("Disabled target", enabled = false),
)

val enabledTargetNames = targets
    .filter(BuildTarget::enabled)
    .map(BuildTarget::name)

println(enabledTargetNames.joinToString())
