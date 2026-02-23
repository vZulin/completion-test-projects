# Completion Test Projects

Test fixture projects for manual IDE Completion testing. Each project provides
code files with marked caret positions (`<caret>`) for systematically testing
completion scenarios.

## Project Structure

### Kotlin (`kotlin/`)

Gradle-based Kotlin project covering TC-1 through TC-93.

| File | Description |
|------|-------------|
| `model/Domain.kt` | Shared data classes and helper functions |
| `trigger/Lifecycle.kt` | Trigger and lifecycle scenarios (TC-1..TC-15) |
| `basic/BasicScope.kt` | Local scope, keywords, ranking (TC-16..TC-25) |
| `member/MemberAccess.kt` | Member, companion, nullable, extension (TC-26..TC-33) |
| `smart/SmartCompletion.kt` | Smart completion (TC-34..TC-38) |
| `args/ArgsParams.kt` | Arguments and named params (TC-39..TC-45) |
| `accept/AcceptCommit.kt` | Accept, commit characters, caret (TC-46..TC-53) |
| `imports/AutoImportScenarios.kt` | Auto-import and conflicts (TC-55..TC-58) |
| `strings/StringPaths.kt` | Path and string completion (TC-60..TC-63) |
| `doc/DocInfo.kt` | QuickDoc, deprecated, KDoc (TC-64..TC-66, TC-76) |
| `generics/GenericsAnnotations.kt` | Generics and annotations (TC-67..TC-73) |
| `templates/TemplatesInject.kt` | Postfix, SQL, regex (TC-81..TC-86) |
| `refactoring/RefactoringAware.kt` | Rename, change signature (TC-87..TC-88) |
| `stability/StabilityPerf.kt` | Dumb mode, performance (TC-89..TC-93) |

### Java (`java/`)

Gradle-based Java project.

| File | Description |
|------|-------------|
| `model/User.java` | Shared model |
| `BasicCombo.java` | Basic member/smart/args completion |
| `Visibility.java` | Private access negative test |
| `StaticMembers.java` | Static vs instance members |
| `SmartCompletion.java` | Smart return and expected type |
| `AcceptCommit.java` | Commit by parenthesis |
| `ImportScenarios.java` | Auto-import and conflicts |
| `DocDeprecated.java` | Javadoc, deprecated, annotations |
| `GenericsAnnot.java` | Generics and annotation params |
| `TemplatesRefactor.java` | Postfix, live templates, rename |
| `PerfMany.java` | Performance (50+ methods) |

### TypeScript (`typescript/`)

| File | Description |
|------|-------------|
| `package.json` | Also serves as TC-78 fixture |
| `tsconfig.json` | Also serves as TC-79 fixture |
| `src/model.ts` | Shared types |
| `src/basicCombo.ts` | Member/smart/args completion |
| `src/importPath.ts` | Import path completion |
| `src/autoImport.ts` + `src/utilModule.ts` | Auto-import scenarios |
| `src/generics.ts` | Generic type params |

### Python (`python/`)

| File | Description |
|------|-------------|
| `basic_combo.py` | Member, smart, path completion (TC-7, TC-26, TC-34, TC-37, TC-60) |
| `decorator_keyword.py` | Decorators, keywords (TC-72, TC-21) |
| `dynamic_named.py` | Dynamic fallback, named args (TC-38, TC-45) |
| `docstring.py` | Docstring completion (TC-77) |

## How to Use

1. Open the desired project in IntelliJ IDEA (Kotlin/Java), WebStorm (TypeScript),
   or PyCharm (Python).
2. Navigate to the file for the test case you want to verify.
3. Find the `<caret>` marker comment with the TC number.
4. Place your cursor at the indicated position.
5. Follow the test case steps from `generated-testcases.md`.

## Caret Position Convention

All caret positions are marked with comments in the format:

```
// <caret> TC-N: brief description          (Kotlin / Java / TypeScript)
#  <caret> TC-N: brief description          (Python)
```

In some cases the caret position is embedded in the code line itself.
Look for the comment to identify where to place your cursor.

## Test Case Reference

See `generated-testcases.md` in the parent directory for the full list of 93 test
cases with detailed steps and expected results.
