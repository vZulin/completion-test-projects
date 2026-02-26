/**
 * Complex completion scenarios: fluent API chains, nested constants, enum options.
 *
 * Covers: TC-94, TC-95, TC-96, TC-97, TC-98.
 *
 * To test: delete the completed suffix described in `// <caret>` comments,
 * place caret at the target position, and invoke completion.
 */
package completion.complex

object TestCases {
    object IU {
        const val IntelliJHelloWorld = "IU.IntelliJHelloWorld"
    }

    object WS {
        const val Empty = "WS.Empty"
    }
}

enum class EnableSettingSyncOptions { PUSH, GET }

class CommandChain {
    fun enableSettingsSync(action: EnableSettingSyncOptions = EnableSettingSyncOptions.PUSH): CommandChain = this
    fun exitApp(): CommandChain = this
}

fun commandChainSaveAll(): CommandChain = CommandChain()

class RunResult {
    fun assertSettingsEqualTo(snapshotId: String): RunResult = this
}

class IdeContext(private val contextId: String, private val template: String) {
    fun copySettingsToIde(snapshotId: String): IdeContext = this
    fun runIDE(commands: CommandChain): RunResult = RunResult()
}

fun initializeTestContext(contextId: String, template: String): IdeContext = IdeContext(contextId, template)

fun testSyncOnNewMachine() {
    // <caret> TC-96-A: Delete 'IntelliJHelloWorld' below after "TestCases.IU.";
    //   invoke completion; expect IntelliJHelloWorld
    val compA = initializeTestContext("testSyncOnNewMachine/compA", TestCases.IU.IntelliJHelloWorld)
        // <caret> TC-94: Delete 'copySettingsToIde(...)' below after initializeTestContext(...).;
        //   expect fluent methods copySettingsToIde(), runIDE()
        .copySettingsToIde("testSyncOnNewMachine/compA")
    // <caret> TC-98: Delete 'enableSettingsSync().exitApp()' below after commandChainSaveAll().;
    //   expect chain methods enableSettingsSync(), exitApp()
    compA.runIDE(commands = commandChainSaveAll().enableSettingsSync().exitApp())

    // <caret> TC-96-B: Delete 'Empty' below after "TestCases.WS.";
    //   invoke completion; expect Empty
    val compB = initializeTestContext("testSyncOnNewMachine/compB", TestCases.WS.Empty)
        .copySettingsToIde("testSyncOnNewMachine/compB")
    compB.runIDE(commands = commandChainSaveAll().enableSettingsSync().exitApp())

    val compC = initializeTestContext("testSyncOnNewMachine/compC", TestCases.IU.IntelliJHelloWorld)
    // <caret> TC-95: Delete 'GET' below after "EnableSettingSyncOptions.";
    //   invoke completion in named arg; expect GET, PUSH
    // <caret> TC-97: Delete 'assertSettingsEqualTo(...)' below after runIDE(...).;
    //   expect assertSettingsEqualTo(...)
    compC.runIDE(commands = commandChainSaveAll().enableSettingsSync(action = EnableSettingSyncOptions.GET).exitApp())
        .assertSettingsEqualTo("testSyncOnNewMachine/compA")
}
