# Completion Test Projects

Fixture projects for manual IDE completion testing.

Primary specification: [`testcases.md`](./testcases.md).

## Test Plans

The test list is grouped by execution plans and ordered by priority:

| Plan | Priority | TC range | Purpose |
|------|----------|----------|---------|
| Acceptance | P0 | TC-1..TC-56 | First-pass checks before deeper runs |
| Regression | P1 | TC-57..TC-113 | Main stability and feature regression |
| Full | P2 | TC-114..TC-121 | Extended/non-critical scenarios |

Run tests in ascending TC order inside each plan.

## Repository Structure

### Kotlin (`kotlin/`)

`kotlin/src/main/kotlin/completion`:
- `trigger/Lifecycle.kt` - popup lifecycle and trigger scenarios.
- `basic/BasicScope.kt` - scope, keywords, ranking.
- `member/MemberAccess.kt` - member/static/nullable/extension completion.
- `smart/SmartCompletion.kt` - smart completion by expected type.
- `args/ArgsParams.kt` - arguments, parameter info, named args.
- `accept/AcceptCommit.kt`, `accept/CommitVariants.kt` - accept/commit behavior.
- `imports/AutoImportScenarios.kt` - auto-import and conflicts.
- `strings/StringPaths.kt` - path and string completion.
- `doc/DocInfo.kt` - QuickDoc/deprecated/doc tags.
- `generics/GenericsAnnotations.kt` - generics/annotations.
- `templates/TemplatesInject.kt` - postfix/live templates/injections.
- `refactoring/RefactoringAware.kt` - rename and signature changes.
- `stability/StabilityPerf.kt`, `stability/LargeFileScope.kt` - dumb mode/performance/large file.
- `statements/StatementCompletion.kt` - statement completion.
- `complex/ComplexDslCompletion.kt` - chained/DSL scenarios.
- `brackets/BracketIndexAccess.kt` - completion in `[]` contexts.
- `advanced/AdvancedCompletionScenarios.kt`, `advanced/ExtensionAutoImportScenario.kt`, `advanced/external/ExternalExtensions.kt` - advanced completion contracts.
- `model/Domain.kt` - shared model/helpers.

### Java (`java/`)

`java/src/main/java/completion`:
- `BasicCombo.java`, `Visibility.java`, `StaticMembers.java`, `SmartCompletion.java`.
- `AcceptCommit.java`, `ImportScenarios.java`.
- `DocDeprecated.java`, `GenericsAnnot.java`, `TopLevelKeyword.java`.
- `StatementCompletion.java`, `PerfMany.java`.
- `BracketIndexCompletion.java`, `AdvancedCompletionContracts.java`.
- `TemplatesRefactor.java`.
- `model/User.java`.

### TypeScript (`typescript/`)

- `typescript/src/basicCombo.ts`, `expectedType.ts`, `smartReturn.ts`, `unknownFallback.ts`.
- `typescript/src/importPath.ts`, `autoImport.ts`, `utilModule.ts`.
- `typescript/src/bracketIndex.ts`, `commitByParen.ts`, `generics.ts`, `keyword.ts`, `model.ts`.
- `typescript/package.json`, `typescript/tsconfig.json` are also fixtures for config completion scenarios.

### Python (`python/`)

- `python/basic_combo.py`
- `python/dynamic_named.py`
- `python/decorator_keyword.py`
- `python/docstring.py`
- `python/bracket_index.py`

## How to Execute a Test Case

1. Open the target project in the relevant IDE.
2. Find the TC in [`testcases.md`](./testcases.md).
3. Open the fixture file referenced in that TC.
4. Locate the marker/comment with the matching `TC-<N>`.
5. Place caret at `<caret>` and apply only the edit described in that marker (if required).
6. Execute steps from the TC and verify expected result.

## Marker Convention

Caret markers in fixtures:

```text
// <caret> TC-N: ...
#  <caret> TC-N: ...
```

Many scenarios intentionally keep a compilable/default line and describe the minimal change in comments
(for example, remove or uncomment a specific token) before invoking completion.
Follow the inline instruction next to the corresponding `TC-N`.
