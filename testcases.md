# Сгенерированные тест-кейсы: Completion

## Метаданные
- **Функциональность:** Completion (Code Completion)
- **Глобальные поля:**
  - Подсистема: Completion
  - Приоритеты: P0 / P1 / P2

---

## План 1: Acceptance (P0)

Тесты этого плана выполняются в порядке `TC-1` → `TC-56`.

### Источник: Раздел 1: Триггеры и жизненный цикл completion popup

#### Функциональность: Ручной вызов (Manual trigger)

#### [ ] TC-1: Вызвать basic completion через Ctrl+Space

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что нажатие Ctrl+Space в позиции каретки открывает список автодополнения.

**Предусловие:**
- IDE открыта с проектом, содержащим файл с приведённым кодом.
- Файл открыт в редакторе.

```kotlin
package demo

import java.io.File

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)
fun consumeUser(u: User) {}

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.` (позиция `<caret>`).
2. Нажать Ctrl+Space.
3. Убедиться, что отображается popup со списком completion.

**Ожидаемый результат:** Открывается popup с предложениями автодополнения (свойства и методы объекта `user`). Ошибок в IDE log нет.

---

#### [ ] TC-2: Вызвать smart completion через Ctrl+Shift+Space

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что нажатие Ctrl+Shift+Space открывает smart completion с тип-ориентированными предложениями.

**Предусловие:**
- Файл с приведённым кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    val u1: User = <caret>
}
```

**Шаги:**
1. Поставить каретку в позицию `<caret>` после `val u1: User = `.
2. Нажать Ctrl+Shift+Space.
3. Убедиться, что отображается popup smart completion.

**Ожидаемый результат:** Открывается popup smart completion с предложениями, соответствующими типу `User` (переменная `user`, вызов `User(...)`, `buildUser(...)`).

---

#### [ ] TC-3: Повторное нажатие Ctrl+Space при открытом popup

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что повторное нажатие Ctrl+Space при открытом popup обновляет/расширяет список без визуальных артефактов.

**Предусловие:**
- Файл с приведённым кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val userName = "Ann"
    val user = User(userName, 21)

    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup открывается.
3. Не выбирая элемент, повторно нажать Ctrl+Space.
4. Убедиться, что popup обновляется или остаётся стабильным.

**Ожидаемый результат:** Список обновляется/расширяется (если в IDE есть второй уровень) или остаётся консистентным без «мигания». Нет визуальных артефактов.

---

#### [ ] TC-4: Закрытие popup по Esc

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что нажатие Esc закрывает popup completion без изменения текста.

**Предусловие:**
- Файл с кодом открыт, каретка после `user.`.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup появляется.
3. Нажать Esc.
4. Убедиться, что popup закрылся.
5. Убедиться, что текст в редакторе не изменился.

**Ожидаемый результат:** Popup закрывается. Текст в редакторе остаётся прежним (`user.`).

---

#### [ ] TC-5: Закрытие popup кликом мыши вне popup

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что клик мышью вне popup completion закрывает его.

**Предусловие:**
- Файл с кодом открыт, каретка после `user.`.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup появляется.
3. Кликнуть мышью в любое место редактора вне popup.
4. Убедиться, что popup закрылся.

**Ожидаемый результат:** Popup закрывается при клике вне его области.

---

#### [ ] TC-6: Фильтрация списка при продолжении ввода текста

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что при вводе текста во время отображения popup completion список фильтруется без зависаний.

**Предусловие:**
- Файл с кодом открыт, каретка после `user.`.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup появляется.
3. Набрать `na`.
4. Убедиться, что список фильтруется (например, `name` остаётся, лишние варианты исчезают).

**Ожидаемый результат:** Список completion фильтруется по введённому префиксу без зависаний и задержек UI.

---

#### Функциональность: Автоматическое появление (Auto-popup)

#### [ ] TC-7: Auto-popup после точки (member completion)

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что popup completion автоматически появляется при вводе `.` после переменной.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user<caret>
}
```

**Шаги:**
1. Поставить каретку сразу после `user`.
2. Ввести символ `.`.
3. Убедиться, что popup completion появляется автоматически (без Ctrl+Space).

**Ожидаемый результат:** Popup появляется автоматически с предложениями членов типа `User` (`name`, `age` и т.д.).

---

#### [ ] TC-8: Auto-popup после открывающей скобки (аргументы)

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что popup completion автоматически появляется при вводе `(` после имени функции.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    buildUser<caret>
}
```

**Шаги:**
1. Поставить каретку после `buildUser`.
2. Ввести символ `(`.
3. Убедиться, что popup completion появляется с предложениями для первого аргумента.

**Ожидаемый результат:** Popup появляется автоматически с предложениями переменных типа `String` для параметра `name`.

---

#### [ ] TC-9: Auto-popup после запятой в списке аргументов

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что popup completion появляется при вводе запятой внутри списка аргументов.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    buildUser(userName<caret>)
}
```

**Шаги:**
1. Поставить каретку после `userName` внутри скобок.
2. Ввести символ `,`.
3. Убедиться, что popup completion появляется для следующего аргумента.

**Ожидаемый результат:** Popup появляется с предложениями переменных типа `Int` для параметра `age` (например, `userAge`).

---

#### [ ] TC-10: Auto-popup при вводе пути в import

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/typescript/src/importPath.ts`

**Описание:** Проверить, что popup completion появляется при вводе пути в директиве import.

**Предусловие:**
- TypeScript/JavaScript проект открыт в редакторе.

```typescript
import x from "./<caret>"
```

**Шаги:**
1. Поставить каретку внутри кавычек после `./`.
2. Начать ввод имени файла или директории.
3. Убедиться, что popup completion появляется с предложениями файлов/папок проекта.

**Ожидаемый результат:** Popup появляется с файлами и папками, соответствующими относительному пути проекта.

---

#### [ ] TC-11: Auto-popup при вводе пути в строке File/Path

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что popup completion появляется при вводе файлового пути внутри строки File("...") или Path("...").

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

import java.io.File

fun main() {
    File("<caret>")
}
```

**Шаги:**
1. Поставить каретку внутри кавычек в `File("")`.
2. Начать ввод пути (например, `src/`).
3. Убедиться, что popup completion показывает файлы и папки проекта.

**Ожидаемый результат:** Popup появляется с предложениями файлов и папок, соответствующих файловой системе проекта.

---

#### Функциональность: Обновление/отмена запросов (race/cancel)

#### [ ] TC-12: Быстрый ввод символов — актуальность списка

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что при быстром наборе 5–10 символов список completion отображает только актуальные результаты без «отката» к старым.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup появляется.
3. Быстро набрать 5–10 символов (например, `toStri`).
4. Убедиться, что список отображает только результаты, соответствующие текущему вводу.

**Ожидаемый результат:** Список completion отображает только актуальные результаты по текущему вводу. Нет «отката» к старым результатам, нет мерцания.

---

#### [ ] TC-13: Удаление префикса Backspace — возврат к полному списку

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что после удаления введённого префикса клавишей Backspace список корректно возвращается к общему набору.

**Предусловие:**
- Файл с кодом открыт, каретка после `user.`.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup появляется.
3. Набрать `na` — список фильтруется.
4. Удалить введённый текст `na` с помощью Backspace.
5. Убедиться, что popup показывает полный набор предложений.

**Ожидаемый результат:** Список корректно возвращается к общему набору. Нет пустых или битых состояний.

---

#### [ ] TC-14: Перемещение каретки стрелками при открытом popup

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/trigger/Lifecycle.kt`

**Описание:** Проверить, что перемещение каретки стрелками при открытом popup корректно закрывает или переоткрывает его.

**Предусловие:**
- Файл с кодом открыт, каретка после `user.`.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
    val x = 42
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup появляется.
3. Нажать стрелку «вверх» или «влево» для перемещения каретки в другое место (без выбора элемента из списка).
4. Убедиться, что popup закрывается или корректно переоткрывается по новому месту.

**Ожидаемый результат:** Popup закрывается или корректно переоткрывается по новому контексту. Нет визуальных артефактов.

---

### Источник: Раздел 2: Basic completion — корректность списка и ранжирование

#### Функциональность: Локальные символы / scope

#### [ ] TC-15: Completion для локальной переменной по префиксу

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`

**Описание:** Проверить, что при вызове completion внутри функции по префиксу локальной переменной она присутствует в списке.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    user<caret>
}
```

**Шаги:**
1. Поставить каретку после `user` (без точки).
2. Нажать Ctrl+Space.
3. Убедиться, что `userName`, `userAge`, `user` присутствуют в списке.

**Ожидаемый результат:** Все локальные переменные с префиксом `user` присутствуют в списке completion.

---

#### [ ] TC-16: Completion для параметра функции

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что параметры функции доступны в списке completion по префиксу.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

fun greet(userName: String, userAge: Int) {
    user<caret>
}
```

**Шаги:**
1. Поставить каретку после `user` внутри тела функции `greet`.
2. Нажать Ctrl+Space.
3. Убедиться, что `userName` и `userAge` присутствуют в списке.

**Ожидаемый результат:** Параметры функции `userName` и `userAge` присутствуют в списке completion.

---

#### [ ] TC-17: Completion показывает символы из внешнего scope

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`

**Описание:** Проверить, что внутри вложенного scope (if/for) completion показывает символы из внешнего scope.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val outerValue = 10
    if (true) {
        out<caret>
    }
}
```

**Шаги:**
1. Поставить каретку после `out` внутри блока `if`.
2. Нажать Ctrl+Space.
3. Убедиться, что `outerValue` присутствует в списке.

**Ожидаемый результат:** Переменная `outerValue` из внешнего scope присутствует в списке completion.

---

#### [ ] TC-18: Символы вне области видимости не предлагаются

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/Visibility.java`

**Описание:** Проверить, что private-члены другого класса не предлагаются в completion, где они недоступны.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

class A { private void secret() {} }

public class B {
  public static void main(String[] args) {
    A a = new A();
    a.sec<caret>
  }
}
```

**Шаги:**
1. Поставить каретку после `a.sec`.
2. Нажать Ctrl+Space.
3. Убедиться, что метод `secret` **не** предлагается в списке (если он недоступен в данном контексте).

**Ожидаемый результат:** Метод `secret()` не отображается в списке completion (или помечен как недоступный).

---

#### Функциональность: Ключевые слова

#### [ ] TC-19: Completion для ключевого слова return

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что при наборе `ret` и вызове completion в списке есть ключевое слово `return`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    ret<caret>
}
```

**Шаги:**
1. Поставить каретку после `ret` внутри функции.
2. Нажать Ctrl+Space.
3. Убедиться, что `return` присутствует в списке.
4. Выбрать `return` из списка.
5. Убедиться, что `return` вставлено корректно.

**Ожидаемый результат:** Ключевое слово `return` есть в списке. При выборе вставляется корректно.

---

#### [ ] TC-20: Completion для ключевого слова class на top-level

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`
- Java: `completion-test-projects/java/src/main/java/completion/TopLevelKeyword.java`
- TypeScript: `completion-test-projects/typescript/src/keyword.ts`
- Python: `completion-test-projects/python/decorator_keyword.py`

**Описание:** Проверить, что при наборе `cla` на верхнем уровне файла в completion есть ключевое слово `class`.

**Предусловие:**
- Новый файл открыт в редакторе.

```kotlin
cla<caret>
```

**Шаги:**
1. Поставить каретку после `cla` на верхнем уровне файла.
2. Нажать Ctrl+Space.
3. Убедиться, что `class` присутствует в списке.

**Ожидаемый результат:** Ключевое слово `class` есть в списке completion.

---

#### Функциональность: Негативные контексты (устойчивость)

#### [ ] TC-21: Completion в «сломанном» месте — устойчивость

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`

**Описание:** Проверить, что вызов completion в заведомо некорректном месте (например, посреди числа) не приводит к падению IDE.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val x = 1234
    val y = 12<caret>34
}
```

**Шаги:**
1. Поставить каретку между `12` и `34` в числовом литерале.
2. Нажать Ctrl+Space.
3. Убедиться, что IDE не падает и нет исключений в логе.

**Ожидаемый результат:** IDE не падает. Popup либо не показывается, либо показывает предсказуемый минимум. Нет исключений в IDE log.

---

### Источник: Раздел 3: Member completion (после доступа к членам)

#### Функциональность: Completion после точки

#### [ ] TC-22: Member completion — свойства и методы типа

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/member/MemberAccess.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что после `.` у объекта в списке completion есть свойства и методы его типа.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space (или дождаться auto-popup).
3. Убедиться, что в списке есть `name`, `age` и другие методы типа `User`.

**Ожидаемый результат:** В списке completion присутствуют свойства `name`, `age` и стандартные методы (`toString`, `hashCode`, `copy` и т.д.).

---

#### [ ] TC-23: Фильтрация member completion по префиксу

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/member/MemberAccess.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`

**Описание:** Проверить, что ввод префикса после `.` фильтрует список до соответствующих членов.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.na<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.na`.
2. Нажать Ctrl+Space (или дождаться auto-popup).
3. Убедиться, что список фильтруется до `name` (и, возможно, других совпадений).

**Ожидаемый результат:** Список отфильтрован до `name` (и/или `getName` — зависит от языка).

---

#### [ ] TC-24: Выбор метода — вставка скобок ()

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что при выборе метода из completion вставляются скобки `()` и каретка оказывается внутри.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

public class Main {
  static class User {
    String name;
    String getName() { return name; }
  }

  public static void main(String[] args) {
    User user = new User();
    user.<caret>
  }
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Вызвать completion.
3. Выбрать метод `getName` из списка и нажать Enter.
4. Убедиться, что вставлено `getName()`.
5. Убедиться, что каретка стоит внутри `()`.

**Ожидаемый результат:** Вставлено `getName()`, каретка — внутри скобок.

---

#### [ ] TC-25: Выбор поля/свойства — вставка без скобок

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/member/MemberAccess.kt`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что при выборе поля/свойства из completion вставляется только имя без скобок.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Вызвать completion.
3. Выбрать свойство `name` из списка и нажать Enter.
4. Убедиться, что вставлено `name` без `()`.

**Ожидаемый результат:** Вставлено только `name`, скобки `()` не добавляются.

---

#### Функциональность: Статические/companion члены

#### [ ] TC-26: Completion для статических/companion членов

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/member/MemberAccess.kt`
- Java: `completion-test-projects/java/src/main/java/completion/StaticMembers.java`

**Описание:** Проверить, что после `ClassName.` предлагаются статические члены или companion-объект.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
class Util {
    companion object {
        fun make(): Util = Util()
    }
}

fun main() {
    Util.<caret>
}
```

**Шаги:**
1. Поставить каретку после `Util.`.
2. Нажать Ctrl+Space.
3. Убедиться, что `make()` присутствует в списке.

**Ожидаемый результат:** В списке completion есть `make()` из companion object.

---

#### [ ] TC-27: Нестатические члены не отображаются в контексте ClassName.

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/StaticMembers.java`

**Описание:** Проверить, что нестатические (instance) члены не появляются при обращении через имя класса.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

public class StaticTest {
  int instanceField = 42;
  static int staticField = 0;

  public static void main(String[] args) {
    Math.<caret>
  }
}
```

**Шаги:**
1. Поставить каретку после `Math.`.
2. Нажать Ctrl+Space.
3. Убедиться, что в списке только статические члены (`abs`, `max`, `PI` и т.д.).
4. Убедиться, что нестатические члены (instance methods) отсутствуют.

**Ожидаемый результат:** В контексте `Math.` предлагаются только статические члены. Нестатические члены отсутствуют.

---

### Источник: Раздел 4: Smart completion (тип-ориентированная)

#### Функциональность: Присваивание

#### [ ] TC-28: Smart completion при присваивании переменной с указанным типом

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/smart/SmartCompletion.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`
- Python: `completion-test-projects/python/basic_combo.py`

**Описание:** Проверить, что smart completion в месте присваивания `val u: User = <caret>` предлагает выражения соответствующего типа.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    val u1: User = <caret>
}
```

**Шаги:**
1. Поставить каретку после `val u1: User = `.
2. Нажать Ctrl+Shift+Space (smart completion).
3. Убедиться, что в списке есть `user`, `User(...)`, `buildUser(...)`.

**Ожидаемый результат:** Smart completion предлагает: переменную `user`, конструктор `User(...)`, функцию `buildUser(...)` — всё типа `User`.

---

#### [ ] TC-29: Выбор функции из smart completion — вставка с кареткой в скобках

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/smart/SmartCompletion.kt`

**Описание:** Проверить, что выбор функции из smart completion вставляет вызов с кареткой внутри скобок.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val u1: User = <caret>
}
```

**Шаги:**
1. Поставить каретку после `val u1: User = `.
2. Нажать Ctrl+Shift+Space.
3. Выбрать `buildUser(...)` из списка и нажать Enter.
4. Убедиться, что вставлен вызов `buildUser()`.
5. Убедиться, что каретка стоит внутри `()` и/или показывается подсказка параметров.

**Ожидаемый результат:** Вставлен вызов `buildUser()`, каретка — внутри скобок. Возможно, отображается подсказка параметров.

---

#### Функциональность: Return

#### [ ] TC-30: Smart completion для return с указанным возвращаемым типом

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/smart/SmartCompletion.kt`
- Java: `completion-test-projects/java/src/main/java/completion/SmartCompletion.java`
- TypeScript: `completion-test-projects/typescript/src/smartReturn.ts`

**Описание:** Проверить, что smart completion после `return` предлагает выражения, соответствующие возвращаемому типу функции.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
data class User(val name: String, val age: Int)

fun buildUser(): User = User("Ann", 21)

fun f(): User {
    return <caret>
}
```

**Шаги:**
1. Поставить каретку после `return `.
2. Нажать Ctrl+Shift+Space.
3. Убедиться, что в списке есть `buildUser()`, `User(...)`.

**Ожидаемый результат:** Smart completion предлагает выражения типа `User`: `buildUser()`, `User(...)`.

---

#### Функциональность: Аргументы функции (по ожидаемому типу)

#### [ ] TC-31: Smart completion по ожидаемому типу аргумента

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/smart/SmartCompletion.kt`
- Java: `completion-test-projects/java/src/main/java/completion/SmartCompletion.java`
- TypeScript: `completion-test-projects/typescript/src/expectedType.ts`

**Описание:** Проверить, что smart completion внутри вызова функции предлагает переменные, соответствующие типу ожидаемого аргумента.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
data class User(val name: String, val age: Int)
fun consumeUser(u: User) {}

fun main() {
    val u = User("Ann", 21)
    consumeUser(<caret>)
}
```

**Шаги:**
1. Поставить каретку внутри `consumeUser()`.
2. Нажать Ctrl+Shift+Space.
3. Убедиться, что `u` присутствует в списке.

**Ожидаемый результат:** Smart completion предлагает переменную `u` типа `User`.

---

#### Функциональность: Fallback при неизвестном типе

#### [ ] TC-32: Smart completion при динамическом/неразрешимом типе

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Python: `completion-test-projects/python/dynamic_named.py`
- TypeScript: `completion-test-projects/typescript/src/unknownFallback.ts`

**Описание:** Проверить, что smart completion в контексте с неразрешимым типом не приводит к ошибкам.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```python
dyn = object()
dyn.<caret>
```

**Шаги:**
1. Поставить каретку после `dyn.`.
2. Нажать Ctrl+Shift+Space (или аналог smart completion).
3. Убедиться, что IDE не падает.
4. Убедиться, что список разумный (может быть как basic completion).

**Ожидаемый результат:** Нет ошибок. Список может содержать базовые методы `object`. IDE работает стабильно.

---

### Источник: Раздел 5: Completion в аргументах/параметрах и сигнатурах

#### Функциональность: Подстановка аргументов

#### [ ] TC-33: Completion аргументов по типу параметра

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/args/ArgsParams.kt`

**Описание:** Проверить, что completion внутри вызова функции предлагает переменные, соответствующие типу текущего параметра.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    buildUser(<caret>)
}
```

**Шаги:**
1. Поставить каретку внутри `buildUser()`.
2. Нажать Ctrl+Space.
3. Убедиться, что `userName` (тип `String`) предлагается для параметра `name`.

**Ожидаемый результат:** Completion предлагает переменные типа `String` для первого параметра `name` (в частности `userName`).

---

#### [ ] TC-34: Completion для второго аргумента после запятой

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/args/ArgsParams.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что после ввода первого аргумента и запятой completion предлагает варианты для следующего параметра.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    buildUser(userName, <caret>)
}
```

**Шаги:**
1. Поставить каретку после запятой внутри `buildUser(userName, )`.
2. Нажать Ctrl+Space.
3. Убедиться, что `userAge` (тип `Int`) предлагается для параметра `age`.

**Ожидаемый результат:** Completion предлагает переменные типа `Int` для второго параметра `age` (в частности `userAge`).

---

#### [ ] TC-35: Completion сразу после открывающей скобки без ввода

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/args/ArgsParams.kt`

**Описание:** Проверить, что вызов completion сразу после `(` без ввода текста показывает непустой список.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    buildUser(<caret>)
}
```

**Шаги:**
1. Поставить каретку внутри `buildUser()` сразу после `(`.
2. Нажать Ctrl+Space.
3. Убедиться, что список не пустой.
4. Убедиться, что нет «мусорных» предложений.

**Ожидаемый результат:** Список не пустой, содержит подходящие переменные (`userName` и т.д.). Нет мусорных предложений.

---

### Источник: Раздел 6: Accept/commit — вставка, замены, commit characters

#### Функциональность: Enter / Tab

#### [ ] TC-36: Принятие элемента из completion клавишей Enter

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`

**Описание:** Проверить, что выбор элемента из completion и нажатие Enter вставляет его в код.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space.
3. Выбрать элемент `name` с помощью стрелок.
4. Нажать Enter.
5. Убедиться, что `name` вставлено в код.

**Ожидаемый результат:** Элемент `name` вставлен в код: `user.name`.

---

#### [ ] TC-37: Принятие элемента из completion клавишей Tab

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`

**Описание:** Проверить, что выбор элемента клавишей Tab вставляет элемент (или выполняет замену — согласно настройкам IDE).

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space.
3. Выбрать элемент `name` с помощью стрелок.
4. Нажать Tab.
5. Убедиться, что `name` вставлено корректно.

**Ожидаемый результат:** Элемент `name` вставлен в код. Поведение соответствует настройке IDE (Insert/Replace).

---

#### [ ] TC-38: Принятие completion при выделенном тексте — замена

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`

**Описание:** Проверить, что при наличии выделенного текста принятие элемента из completion заменяет выделение.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val userName = "Ann"
    val x = "REPLACE_ME"<caret>
}
```

**Шаги:**
1. Выделить текст `"REPLACE_ME"` в редакторе.
2. Нажать Ctrl+Space для вызова completion.
3. Выбрать `userName` из списка.
4. Нажать Enter.
5. Убедиться, что `"REPLACE_ME"` заменён на `userName`.

**Ожидаемый результат:** Выделенный текст полностью заменён на `userName`.

---

#### Функциональность: Commit characters

#### [ ] TC-39: Принятие completion нажатием точки (commit by `.`)

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`

**Описание:** Проверить, что принятие completion нажатием `.` вставляет элемент и добавляет точку.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.na<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.na`.
2. Нажать Ctrl+Space — popup с `name` появляется.
3. Нажать `.` (точку) для принятия completion.
4. Убедиться, что результат: `user.name.`.

**Ожидаемый результат:** Вставлено `name`, за ним точка: `user.name.`. Точка не дублируется.

---

#### [ ] TC-40: Принятие completion нажатием открывающей скобки (commit by `(`)

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Java: `completion-test-projects/java/src/main/java/completion/AcceptCommit.java`
- TypeScript: `completion-test-projects/typescript/src/commitByParen.ts`

**Описание:** Проверить, что принятие completion нажатием `(` вставляет метод и скобки без дублирования.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

public class Commit {
  static void doWork() {}
  public static void main(String[] args) {
    doW<caret>
  }
}
```

**Шаги:**
1. Поставить каретку после `doW`.
2. Нажать Ctrl+Space — popup с `doWork` появляется.
3. Нажать `(` для принятия completion.
4. Убедиться, что результат: `doWork()`.
5. Убедиться, что `(` не дублируется.

**Ожидаемый результат:** Вставлено `doWork()`. Открывающая скобка не дублируется.

---

#### Функциональность: Позиция каретки (Caret placement)

#### [ ] TC-41: Каретка внутри скобок после выбора функции

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что при выборе функции из completion каретка ставится внутри `()`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    build<caret>
}
```

**Шаги:**
1. Поставить каретку после `build`.
2. Нажать Ctrl+Space.
3. Выбрать `buildUser` и нажать Enter.
4. Убедиться, что вставлен `buildUser()`.
5. Убедиться, что каретка стоит внутри `()`.

**Ожидаемый результат:** Каретка стоит между `(` и `)` для ввода аргументов.

---

#### [ ] TC-42: Каретка в ожидаемом месте после выбора конструктора

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`

**Описание:** Проверить, что при выборе конструктора/инициализации каретка оказывается в ожидаемом месте.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

public class Main {
  static class User {
    String name;
    int age;
  }

  public static void main(String[] args) {
    User u = new <caret>
  }
}
```

**Шаги:**
1. Поставить каретку после `new `.
2. Нажать Ctrl+Space.
3. Выбрать `User()` из списка и нажать Enter.
4. Убедиться, что вставлен `User()`.
5. Убедиться, что каретка стоит в ожидаемом месте (внутри скобок или после имени).

**Ожидаемый результат:** Конструктор вставлен, каретка — в ожидаемом месте (внутри скобок для ввода аргументов).

---

### Источник: Раздел 7: Auto-import и символы из зависимостей

#### Функциональность: Импорт при выборе из completion

#### [ ] TC-43: Java — auto-import ArrayList при принятии completion

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/ImportScenarios.java`

**Описание:** Проверить, что принятие completion для неимпортированного класса добавляет соответствующий import.

**Предусловие:**
- Файл с кодом открыт в редакторе. Импорт `java.util.ArrayList` отсутствует.

```java
package demo;

public class Imports {
  public static void main(String[] args) {
    ArrayLis<caret> list;
  }
}
```

**Шаги:**
1. Поставить каретку после `ArrayLis`.
2. Нажать Ctrl+Space.
3. Выбрать `ArrayList` из списка и нажать Enter.
4. Убедиться, что вставлено `ArrayList`.
5. Убедиться, что `import java.util.ArrayList;` добавлен в начало файла.
6. Убедиться, что код компилируется.

**Ожидаемый результат:** `ArrayList` вставлен. `import java.util.ArrayList;` добавлен автоматически. Код компилируется.

---

#### [ ] TC-44: Kotlin — auto-import File при принятии completion

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/imports/AutoImportScenarios.kt`

**Описание:** Проверить, что в Kotlin принятие completion для неимпортированного класса добавляет import.

**Предусловие:**
- Файл с кодом открыт в редакторе. Импорт `java.io.File` отсутствует.

```kotlin
package demo

fun main() {
    val f = File("<caret>")
}
```

**Шаги:**
1. Набрать `Fil` и вызвать Ctrl+Space.
2. Выбрать `File` из списка и нажать Enter.
3. Убедиться, что `import java.io.File` добавлен.

**Ожидаемый результат:** `import java.io.File` добавлен автоматически.

---

#### [ ] TC-45: TypeScript — auto-import при принятии completion

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/typescript/src/autoImport.ts`

**Описание:** Проверить, что в TypeScript принятие completion по неимпортированному символу добавляет import.

**Предусловие:**
- Проект содержит файл `a.ts` с экспортом:

```typescript
// a.ts
export function utilFn() {}
```

- Файл `b.ts` открыт в редакторе:

```typescript
// b.ts
utilF<caret>
```

**Шаги:**
1. Поставить каретку после `utilF`.
2. Нажать Ctrl+Space.
3. Выбрать `utilFn` из списка и нажать Enter.
4. Убедиться, что добавлена строка импорта `import { utilFn } from "./a";`.

**Ожидаемый результат:** Добавлен `import { utilFn } from "./a";` автоматически.

---

#### Функциональность: Конфликт имён

#### [ ] TC-46: Completion при конфликте имён — два класса из разных пакетов

**Приоритет:** P0
**План:** Acceptance
**Тестовые файлы:**
- Java: `completion-test-projects/java/src/main/java/completion/ImportScenarios.java`
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/imports/AutoImportScenarios.kt`

**Описание:** Проверить, что при наличии двух классов с одинаковым именем из разных пакетов completion показывает различия.

**Предусловие:**
- Файл с кодом открыт. Оба пакета `java.util.Date` и `java.sql.Date` доступны в classpath.

```java
package demo;

public class Conflict {
  public static void main(String[] args) {
    Date<caret> d;
  }
}
```

**Шаги:**
1. Поставить каретку после `Date`.
2. Нажать Ctrl+Space.
3. Убедиться, что completion показывает оба варианта с указанием пакета: `java.util.Date` и `java.sql.Date`.
4. Выбрать `java.util.Date`.
5. Убедиться, что импортирован правильный класс.

**Ожидаемый результат:** Оба варианта отображаются с квалификаторами пакетов. При выборе импортируется правильный.

---

### Источник: Раздел 8: Completion в строках (paths/resources)

#### Функциональность: File system paths

#### [ ] TC-47: Path completion внутри File("...")

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/strings/StringPaths.kt`

**Описание:** Проверить, что в строке `File("...")` completion показывает папки и файлы проекта.

**Предусловие:**
- Файл с кодом открыт в редакторе. В проекте существует директория `src/`.

```kotlin
package demo

import java.io.File

fun main() {
    File("<caret>")
}
```

**Шаги:**
1. Поставить каретку внутри кавычек `File("")`.
2. Начать ввод `src/`.
3. Убедиться, что popup показывает папки и файлы проекта.

**Ожидаемый результат:** Popup отображает файлы и папки, соответствующие файловой структуре проекта.

---

#### [ ] TC-48: Выбор файла из path completion — корректность строки

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/strings/StringPaths.kt`

**Описание:** Проверить, что выбор файла из path completion корректно дополняет строку.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

import java.io.File

fun main() {
    File("src/<caret>")
}
```

**Шаги:**
1. Поставить каретку внутри кавычек после `src/`.
2. Вызвать completion.
3. Выбрать файл из списка и нажать Enter.
4. Убедиться, что строка дополнена корректно.
5. Убедиться, что кавычки не сломаны.
6. Убедиться, что каретка остаётся внутри строки.

**Ожидаемый результат:** Строка дополнена корректно. Кавычки целы. Каретка — внутри строки.

---

#### [ ] TC-49: Completion относительных путей `./` и `../`

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/typescript/src/importPath.ts`

**Описание:** Проверить, что completion относительных путей `./` и `../` работает корректно.

**Предусловие:**
- TypeScript-проект открыт в редакторе.

```typescript
import x from "./<caret>"
```

**Шаги:**
1. Поставить каретку внутри кавычек после `./`.
2. Вызвать completion.
3. Убедиться, что popup показывает файлы и папки текущей директории.
4. Изменить путь на `../`.
5. Убедиться, что popup показывает файлы и папки родительской директории.

**Ожидаемый результат:** Список соответствует файловой структуре для `./` и `../`.

---

### Источник: Раздел 17: Dumb mode / индексация / устойчивость

#### Функциональность: Dumb mode / индексация

#### [ ] TC-50: Completion в Dumb mode — устойчивость

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/stability/StabilityPerf.kt`

**Описание:** Проверить, что вызов completion во время индексации (Dumb mode) не приводит к падению IDE.

**Предусловие:**
- IDE открыта, идёт индексация проекта (Dumb mode).
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Открыть крупный проект (или инвалидировать кэш) для запуска индексации.
2. Во время индексации (видна полоса прогресса) поставить каретку после `user.`.
3. Нажать Ctrl+Space.
4. Убедиться, что IDE не падает.
5. Убедиться, что показывается ограниченный список или информация об ограниченных результатах.

**Ожидаемый результат:** IDE не падает. Показывается ограниченный список или уведомление об ограничениях.

---

#### [ ] TC-51: Множественный вызов completion в Dumb mode — нет зависаний

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/stability/StabilityPerf.kt`

**Описание:** Проверить, что 5-кратный вызов completion во время индексации не приводит к зависаниям UI.

**Предусловие:**
- IDE в Dumb mode (идёт индексация).
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Во время индексации поставить каретку после `user.`.
2. Нажать Ctrl+Space 5 раз подряд.
3. Убедиться, что нет зависаний UI.
4. Убедиться, что нет очереди «старых» попапов (накладывающихся друг на друга).

**Ожидаемый результат:** Нет зависаний UI. Нет очереди старых попапов. IDE остаётся отзывчивой.

---

#### [ ] TC-52: Completion после завершения индексации — полные результаты

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/stability/StabilityPerf.kt`

**Описание:** Проверить, что после завершения индексации completion возвращает полные и корректные результаты.

**Предусловие:**
- IDE только что завершила индексацию.
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Дождаться завершения индексации (полоса прогресса исчезла).
2. Поставить каретку после `user.`.
3. Нажать Ctrl+Space.
4. Убедиться, что результаты полные и корректные (все свойства и методы `User` доступны).

**Ожидаемый результат:** Completion возвращает полный набор членов типа `User`. Результаты корректные.

---

### Источник: Раздел 21: Completion внутри индексаторов `[]`

#### Функциональность: Completion индексного выражения в списках/массивах

#### [ ] TC-53: Kotlin — completion Int-индекса внутри `list[ ... ]`

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/brackets/BracketIndexAccess.kt`

**Описание:** Проверить, что completion внутри `[]` для списка предлагает переменные подходящего типа индекса (`Int`).

**Предусловие:**
- Файл `BracketIndexAccess.kt` открыт в редакторе.

```kotlin
val users = listOf(User("Ann", 21), User("Bob", 22))
val selectedIndex = 0
val target = users[<caret>]
```

**Шаги:**
1. Поставить каретку внутри `users[ ]`.
2. Вызвать completion (Ctrl+Space).
3. Проверить, что предлагается `selectedIndex`.
4. Выбрать `selectedIndex` и убедиться, что выражение остаётся валидным.

**Ожидаемый результат:** Внутри `[]` предлагаются релевантные Int-выражения для индексации. Синтаксис не ломается.

---

#### [ ] TC-54: Java — completion индекса в `array[ ... ]` и `matrix[ ... ][ ... ]`

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:**
- Java: `completion-test-projects/java/src/main/java/completion/BracketIndexCompletion.java`

**Описание:** Проверить completion внутри `[]` для одномерного и двумерного массива в Java.

**Предусловие:**
- Файл `BracketIndexCompletion.java` открыт в редакторе.

```java
int[] ages = {21, 22, 23};
int position = 1;
int a = ages[<caret>];

int[][] matrix = {{1, 2}, {3, 4}};
int row = 0;
int col = 1;
int b = matrix[<caret>][<caret>];
```

**Шаги:**
1. Поставить каретку внутри `ages[ ]` и вызвать completion.
2. Выбрать `position`.
3. Повторить для `matrix[ ][ ]`: сначала выбрать `row`, затем `col`.
4. Убедиться, что выражения компилируемы.

**Ожидаемый результат:** Completion в Java корректно работает внутри `[]` и предлагает локальные переменные подходящего типа.

---

#### [ ] TC-55: TypeScript — completion в `obj["..."]` и `arr[ ... ]`

**Приоритет:** P0
**План:** Acceptance
**Тестовый файл:**
- TypeScript: `completion-test-projects/typescript/src/bracketIndex.ts`

**Описание:** Проверить completion внутри `[]` для строкового ключа объекта и числового индекса массива.

**Предусловие:**
- Файл `bracketIndex.ts` открыт в редакторе.

```typescript
type User = { name: string; age: number };
const user: User = { name: "Ann", age: 21 };
const a = user["<caret>"];

const names = ["Ann", "Bob"];
const idx = 0;
const b = names[<caret>];
```

**Шаги:**
1. Поставить каретку внутри `user[""]`, ввести `na`, вызвать completion.
2. Убедиться, что предлагается ключ `name`.
3. Поставить каретку внутри `names[ ]`, вызвать completion.
4. Убедиться, что предлагается `idx`.

**Ожидаемый результат:** Completion работает в обоих bracket-контекстах: ключи объекта и индекс массива.

---

### Источник: Раздел 22: Продвинутые контексты completion

#### Функциональность: Type-flow, DSL, refactoring, устойчивость

#### [ ] TC-56: Completion в синтаксически «битом» коде (resilience)

**Приоритет:** P0  
**План:** Acceptance
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/AdvancedCompletionContracts.java`

**Описание:** Проверить устойчивость completion при временно невалидном синтаксисе.

**Шаги:**
1. В `brokenCodeResilience` временно удалить закрывающую `)` у `println`.
2. Поставить каретку после `user.` и вызвать completion.
3. Проверить, что IDE не падает и показывает релевантные предложения.

**Ожидаемый результат:** Completion остаётся работоспособным в broken-коде; ошибок в IDE log нет.

---

## План 2: Regression (P1)

Тесты этого плана выполняются в порядке `TC-57` → `TC-113`.

### Источник: Раздел 1: Триггеры и жизненный цикл completion popup

#### Функциональность: Автоматическое появление (Auto-popup)

#### [ ] TC-57: Auto-popup при вводе @ (аннотация/декоратор)

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/generics/GenericsAnnotations.kt`
- Java: `completion-test-projects/java/src/main/java/completion/DocDeprecated.java`
- Python: `completion-test-projects/python/decorator_keyword.py`

**Описание:** Проверить, что popup completion появляется при вводе символа `@` для аннотаций/декораторов.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

<caret>
class A
```

**Шаги:**
1. Поставить каретку перед `class A`.
2. Ввести символ `@`.
3. Убедиться, что popup completion появляется с предложениями аннотаций.

**Ожидаемый результат:** Popup появляется с доступными аннотациями (`Deprecated`, `Suppress`, `JvmStatic` и т.д.).

---

### Источник: Раздел 2: Basic completion — корректность списка и ранжирование

#### Функциональность: Ранжирование (relevance)

#### [ ] TC-58: Локальная переменная ранжируется выше глобальной

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`

**Описание:** Проверить, что при одинаковом префиксе локальная переменная ранжируется выше глобальной/импортированной.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

val globalUser = User("Global", 0)

data class User(val name: String, val age: Int)

fun main() {
    val globalUserName = "LocalShadow"
    glo<caret>
}
```

**Шаги:**
1. Поставить каретку после `glo` внутри функции `main`.
2. Нажать Ctrl+Space.
3. Убедиться, что `globalUserName` (локальная) стоит выше `globalUser` (глобальная) в списке.

**Ожидаемый результат:** Локальная переменная `globalUserName` ранжируется выше глобальной `globalUser`.

---

#### [ ] TC-59: Точный префикс ранжируется выше fuzzy-совпадений

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`

**Описание:** Проверить, что при вводе точного префикса элемент с точным совпадением отображается выше fuzzy-совпадений.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val userName = "Ann"
    val usernameLegacy = "legacy"
    userN<caret>
}
```

**Шаги:**
1. Поставить каретку после `userN`.
2. Нажать Ctrl+Space.
3. Убедиться, что `userName` стоит выше `usernameLegacy` и других fuzzy-совпадений.

**Ожидаемый результат:** `userName` (точное совпадение по регистру) ранжируется выше `usernameLegacy` (fuzzy-совпадение).

---

#### [ ] TC-60: MRU — ранее выбранный элемент поднимается, но не ломает релевантность

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/basic/BasicScope.kt`

**Описание:** Проверить поведение MRU (Most Recently Used): ранее выбранный элемент поднимается в списке, но не становится №1 в нерелевантном контексте.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Вызвать completion и выбрать `name`.
3. Отменить действие (Ctrl+Z).
4. Повторить выбор `name` из completion ещё 2 раза (всего 3 выбора).
5. Снова вызвать completion после `user.`.
6. Убедиться, что `name` поднимается в списке.
7. Перейти в другой контекст (например, `user.toString()`) и вызвать completion.
8. Убедиться, что MRU не ломает релевантность в новом контексте.

**Ожидаемый результат:** Ранее выбранный элемент `name` поднимается в контексте `user.`, но не становится №1 в нерелевантном контексте.

---

### Источник: Раздел 3: Member completion (после доступа к членам)

#### Функциональность: Nullable / safe access (Kotlin)

#### [ ] TC-61: Safe-call completion для nullable типа

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/member/MemberAccess.kt`

**Описание:** Проверить, что `?.` для nullable-типа показывает члены типа, применимые к safe-call.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
data class User(val name: String, val age: Int)

fun main() {
    val u: User? = null
    u?.<caret>
}
```

**Шаги:**
1. Поставить каретку после `u?.`.
2. Нажать Ctrl+Space (или дождаться auto-popup).
3. Убедиться, что в списке присутствуют `name`, `age` и другие члены типа `User`.

**Ожидаемый результат:** Completion показывает члены типа `User` (`name`, `age`, `toString` и т.д.), применимые к safe-call.

---

#### Функциональность: Extension methods (Kotlin)

#### [ ] TC-62: Extension function отображается в completion

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/member/MemberAccess.kt`

**Описание:** Проверить, что extension-функция отображается в completion для типа, к которому она определена.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun String.extHello(): String = "Hello, $this"

fun main() {
    val s = "Ann"
    s.<caret>
}
```

**Шаги:**
1. Поставить каретку после `s.`.
2. Нажать Ctrl+Space.
3. Убедиться, что `extHello()` присутствует в списке completion.

**Ожидаемый результат:** Extension-функция `extHello()` отображается в списке completion для строковой переменной `s`.

---

### Источник: Раздел 5: Completion в аргументах/параметрах и сигнатурах

#### Функциональность: Интеграция с Parameter Info

#### [ ] TC-63: Parameter info при выборе функции из completion

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/args/ArgsParams.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`
- TypeScript: `completion-test-projects/typescript/src/basicCombo.ts`

**Описание:** Проверить, что при выборе функции из completion отображается parameter info с подсветкой текущего параметра.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    build<caret>
}
```

**Шаги:**
1. Поставить каретку после `build`.
2. Нажать Ctrl+Space.
3. Выбрать `buildUser` из списка и нажать Enter.
4. Убедиться, что вставлен `buildUser()` с кареткой внутри скобок.
5. Убедиться, что отображается parameter info (подсказка с параметрами).
6. Начать ввод первого аргумента.
7. Убедиться, что текущий параметр подсвечен в parameter info.

**Ожидаемый результат:** Отображается parameter info с подсветкой текущего параметра `name: String`.

---

#### Функциональность: Именованные аргументы (Named arguments)

#### [ ] TC-64: Kotlin — completion именованных аргументов

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/args/ArgsParams.kt`

**Описание:** Проверить, что completion внутри вызова функции в Kotlin предлагает именованные аргументы `name =` и `age =`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
data class User(val name: String, val age: Int)
fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    buildUser(<caret>)
}
```

**Шаги:**
1. Поставить каретку внутри `buildUser()`.
2. Нажать Ctrl+Space.
3. Убедиться, что в списке есть варианты `name = ` и `age = `.

**Ожидаемый результат:** Completion предлагает именованные аргументы `name = ` и `age = `.

---

#### [ ] TC-65: Kotlin — уже использованный named argument не предлагается повторно

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/args/ArgsParams.kt`

**Описание:** Проверить, что после указания именованного аргумента он не предлагается повторно.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
data class User(val name: String, val age: Int)
fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    buildUser(name = "Ann", <caret>)
}
```

**Шаги:**
1. Поставить каретку после `name = "Ann", `.
2. Нажать Ctrl+Space.
3. Убедиться, что `name = ` **не** предлагается повторно.
4. Убедиться, что `age = ` предлагается.

**Ожидаемый результат:** Уже использованный `name = ` не предлагается. Предлагается только `age = `.

---

#### [ ] TC-66: Python — completion именованных параметров

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/python/dynamic_named.py`

**Описание:** Проверить, что completion для функции с именованными параметрами в Python предлагает имена параметров.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```python
def f(name: str, age: int) -> None:
    pass

f(<caret>)
```

**Шаги:**
1. Поставить каретку внутри `f()`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются `name=` и `age=` (если IDE поддерживает).

**Ожидаемый результат:** Completion предлагает именованные параметры `name=` и `age=`.

---

### Источник: Раздел 6: Accept/commit — вставка, замены, commit characters

#### Функциональность: Commit characters

#### [ ] TC-67: Принятие completion нажатием запятой в аргументах

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/accept/AcceptCommit.kt`

**Описание:** Проверить, что принятие completion нажатием `,` в списке аргументов вставляет элемент и запятую корректно.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val userName = "Ann"
    val userAge = 21
    buildUser(user<caret>)
}
```

**Шаги:**
1. Поставить каретку после `user` внутри `buildUser()`.
2. Нажать Ctrl+Space — `userName` появляется в popup.
3. Нажать `,` для принятия completion.
4. Убедиться, что вставлено `userName,` (с запятой и пробелом по code style).
5. Убедиться, что запятая не дублируется.

**Ожидаемый результат:** Вставлено `userName,` корректно. Запятая не дублируется, пробелы по code style.

---

### Источник: Раздел 7: Auto-import и символы из зависимостей

#### Функциональность: Конфликт имён

#### [ ] TC-68: После выбора одного варианта второй не подмешивается

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/ImportScenarios.java`

**Описание:** Проверить, что после выбора одного класса из конфликтующих пакетов второй не добавляется в imports автоматически.

**Предусловие:**
- Файл с конфликтом имён открыт (см. TC-46).

```java
package demo;

public class Conflict {
  public static void main(String[] args) {
    Date<caret> d;
  }
}
```

**Шаги:**
1. Выбрать `java.util.Date` из completion (как в TC-46).
2. Убедиться, что добавлен только `import java.util.Date;`.
3. Убедиться, что `import java.sql.Date;` **не** добавлен.

**Ожидаемый результат:** В файле только один import — `java.util.Date;`. Второй пакет не подмешивается.

---

#### Функциональность: Поведение в зависимости от настроек

#### [ ] TC-69: Отключённый auto-import — поведение при completion

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/ImportScenarios.java`

**Описание:** Проверить, что при выключенном auto-import принятие completion не добавляет import (или добавляет fully qualified name).

**Предусловие:**
- Auto-import отключён в настройках IDE.
- Файл с кодом открыт в редакторе.

```java
package demo;

public class Imports {
  public static void main(String[] args) {
    ArrayLis<caret> list;
  }
}
```

**Шаги:**
1. Выключить auto-import в настройках (Settings → Editor → General → Auto Import).
2. Поставить каретку после `ArrayLis`.
3. Нажать Ctrl+Space.
4. Выбрать `ArrayList` и нажать Enter.
5. Убедиться, что import **не** добавлен (или добавлено fully qualified name `java.util.ArrayList`).

**Ожидаемый результат:** Поведение соответствует настройке: import не добавлен или вставлен fully qualified name.

---

### Источник: Раздел 8: Completion в строках (paths/resources)

#### Функциональность: Негатив — обычная строка

#### [ ] TC-70: Completion в обычной строке — без навязывания path completion

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/strings/StringPaths.kt`

**Описание:** Проверить, что в строковом литерале, не похожем на путь, completion не навязывает file path completion.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val s = "hello <caret> world"
}
```

**Шаги:**
1. Поставить каретку внутри строки `"hello <caret> world"`.
2. Нажать Ctrl+Space.
3. Убедиться, что path completion не навязывается агрессивно.
4. Убедиться, что поведение стабильное и повторяемое.

**Ожидаемый результат:** File path completion не появляется агрессивно. Поведение стабильное. IDE не падает.

---

### Источник: Раздел 9: Documentation popup / QuickDoc / детали элемента

#### Функциональность: QuickDoc и детали элемента

#### [ ] TC-71: QuickDoc для элемента из completion (Ctrl+Q)

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/doc/DocInfo.kt`
- Java: `completion-test-projects/java/src/main/java/completion/BasicCombo.java`

**Описание:** Проверить, что при открытом completion нажатие Ctrl+Q показывает QuickDoc для выбранного элемента.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space — popup completion появляется.
3. Выбрать `name` стрелками.
4. Нажать Ctrl+Q.
5. Убедиться, что показывается QuickDoc для свойства `name`.

**Ожидаемый результат:** Отображается QuickDoc с информацией о свойстве `name` (тип, описание).

---

#### [ ] TC-72: Отображение типа/сигнатуры в списке completion

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/doc/DocInfo.kt`

**Описание:** Проверить, что в списке completion корректно отображается тип или сигнатура каждого элемента.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Нажать Ctrl+Space.
3. Убедиться, что для каждого элемента отображается тип (например, `name: String`, `age: Int`).
4. Убедиться, что для функций отображается сигнатура (например, `(String, Int) -> User`).

**Ожидаемый результат:** Типы и сигнатуры корректно отображаются для каждого элемента в списке completion.

---

#### [ ] TC-73: Deprecated элементы помечены визуально

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/doc/DocInfo.kt`
- Java: `completion-test-projects/java/src/main/java/completion/DocDeprecated.java`

**Описание:** Проверить, что deprecated элементы отображаются с визуальной пометкой (зачёркивание, иконка).

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
@Deprecated("old")
fun oldFun() {}

fun main() {
    old<caret>
}
```

**Шаги:**
1. Поставить каретку после `old`.
2. Нажать Ctrl+Space.
3. Убедиться, что `oldFun` присутствует в списке.
4. Убедиться, что `oldFun` помечен как deprecated (зачёркивание текста, специальная иконка или стиль).

**Ожидаемый результат:** `oldFun` отображается в completion с визуальной пометкой deprecated.

---

### Источник: Раздел 10: Generics / type parameters

#### Функциональность: Generics / type parameters

#### [ ] TC-74: Completion типов внутри generic-параметров

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/generics/GenericsAnnotations.kt`
- Java: `completion-test-projects/java/src/main/java/completion/GenericsAnnot.java`
- TypeScript: `completion-test-projects/typescript/src/generics.ts`

**Описание:** Проверить, что внутри `<>` completion предлагает типы из scope и пакетов.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

import java.util.ArrayList

fun main() {
    val x: ArrayList<<caret>> = arrayListOf()
}
```

**Шаги:**
1. Поставить каретку внутри `<>`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются типы: `String`, `Int`, `User` и т.д.

**Ожидаемый результат:** Completion предлагает доступные типы из scope и импортированных пакетов.

---

#### [ ] TC-75: Completion второго type parameter в Map

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/GenericsAnnot.java`

**Описание:** Проверить, что completion второго type parameter в `Map<String, <caret>>` работает корректно.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

import java.util.Map;

public class G {
  void f(Map<String, <caret>> b) {}
}
```

**Шаги:**
1. Поставить каретку после запятой внутри `Map<String, >`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются типы.
4. Выбрать тип.
5. Убедиться, что `>` и запятые не ломаются.

**Ожидаемый результат:** Типы предлагаются корректно. Вставка типа не ломает `>` и запятые.

---

#### [ ] TC-76: Auto-import при выборе типа в generic-параметре

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/generics/GenericsAnnotations.kt`
- Java: `completion-test-projects/java/src/main/java/completion/GenericsAnnot.java`

**Описание:** Проверить, что выбор типа из completion внутри generic-параметра добавляет import при необходимости.

**Предусловие:**
- Файл с кодом открыт в редакторе. Импорт нужного типа отсутствует.

```kotlin
package demo

import java.util.ArrayList

fun main() {
    val x: ArrayList<<caret>> = arrayListOf()
}
```

**Шаги:**
1. Поставить каретку внутри `<>`.
2. Нажать Ctrl+Space.
3. Начать набирать имя неимпортированного типа (например, `BigDec`).
4. Выбрать `BigDecimal`.
5. Убедиться, что `import java.math.BigDecimal` добавлен.

**Ожидаемый результат:** Import добавлен автоматически при выборе типа из completion.

---

### Источник: Раздел 11: Аннотации / атрибуты / декораторы

#### Функциональность: Аннотации / атрибуты / декораторы

#### [ ] TC-77: Java — completion аннотации @Deprecated

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/DocDeprecated.java`

**Описание:** Проверить, что при наборе `@Depr` в Java completion предлагает `@Deprecated`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

@Depr<caret>
class A {}
```

**Шаги:**
1. Поставить каретку после `@Depr`.
2. Нажать Ctrl+Space.
3. Убедиться, что `@Deprecated` присутствует в списке.
4. Выбрать `@Deprecated`.
5. Убедиться, что аннотация вставлена корректно.

**Ожидаемый результат:** `@Deprecated` предлагается и вставляется корректно.

---

#### [ ] TC-78: Kotlin — completion аннотации @Deprecated

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/generics/GenericsAnnotations.kt`

**Описание:** Проверить, что при наборе `@Depre` в Kotlin completion предлагает `@Deprecated`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

@Depre<caret>
class A
```

**Шаги:**
1. Поставить каретку после `@Depre`.
2. Нажать Ctrl+Space.
3. Убедиться, что `@Deprecated` присутствует в списке.
4. Выбрать `@Deprecated`.
5. Убедиться, что аннотация вставлена корректно.

**Ожидаемый результат:** `@Deprecated` предлагается и вставляется корректно.

---

#### [ ] TC-79: Python — completion декоратора @dataclass

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/python/decorator_keyword.py`

**Описание:** Проверить, что при наборе `@datac` в Python completion предлагает `@dataclass`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```python
from dataclasses import dataclass

@datac<caret>
class A:
    pass
```

**Шаги:**
1. Поставить каретку после `@datac`.
2. Нажать Ctrl+Space.
3. Убедиться, что `@dataclass` присутствует в списке.
4. Выбрать `@dataclass`.
5. Убедиться, что декоратор вставлен корректно.

**Ожидаемый результат:** `@dataclass` предлагается и вставляется корректно.

---

#### [ ] TC-80: Completion параметров аннотации

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/generics/GenericsAnnotations.kt`
- Java: `completion-test-projects/java/src/main/java/completion/GenericsAnnot.java`

**Описание:** Проверить, что completion внутри скобок аннотации предлагает имена параметров и значения по типу.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
annotation class Ann(val value: String, val enabled: Boolean)

@Ann(<caret>)
class A
```

**Шаги:**
1. Поставить каретку внутри `@Ann()`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются `value = ` и `enabled = `.

**Ожидаемый результат:** Completion предлагает имена параметров аннотации: `value = `, `enabled = `.

---

### Источник: Раздел 12: Doc/Comments completion

#### Функциональность: Doc/Comments completion

#### [ ] TC-81: Javadoc — completion тегов @param, @return, @throws

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/DocDeprecated.java`

**Описание:** Проверить, что в Javadoc-комментарии completion предлагает теги @param, @return, @throws.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

public class Doc {
  /**
   * <caret>
   */
  int f(int a) { return a; }
}
```

**Шаги:**
1. Поставить каретку после `* ` внутри Javadoc.
2. Набрать `@`.
3. Нажать Ctrl+Space.
4. Убедиться, что предлагаются `@param`, `@return`, `@throws`.

**Ожидаемый результат:** Completion предлагает теги Javadoc: `@param`, `@return`, `@throws`.

---

#### [ ] TC-82: Javadoc — @param вставляется с именем параметра

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/DocDeprecated.java`

**Описание:** Проверить, что выбор `@param` из completion вставляет шаблон с именем параметра.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```java
package demo;

public class Doc {
  /**
   * <caret>
   */
  int f(int a) { return a; }
}
```

**Шаги:**
1. Поставить каретку внутри Javadoc.
2. Набрать `@param` (или выбрать из completion).
3. Убедиться, что вставляется шаблон с именем параметра `a` (или предлагается выбрать параметр).

**Ожидаемый результат:** Вставлен `@param a ` или предлагается выбор параметра из списка.

---

#### [ ] TC-83: KDoc — completion тегов @param/@return

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/doc/DocInfo.kt`

**Описание:** Проверить, что в KDoc completion предлагает теги @param/@return.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
/**
 * <caret>
 */
fun f(a: Int): Int = a
```

**Шаги:**
1. Поставить каретку внутри KDoc.
2. Набрать `@`.
3. Нажать Ctrl+Space.
4. Убедиться, что предлагаются `@param`, `@return`.

**Ожидаемый результат:** Completion предлагает теги KDoc: `@param`, `@return`.

---

#### [ ] TC-84: Python — completion секций в docstring

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/python/docstring.py`

**Описание:** Проверить, что в Python docstring completion предлагает секции/параметры (если IDE поддерживает).

**Предусловие:**
- Файл с кодом открыт в редакторе.

```python
def f(a: int) -> int:
    """
    <caret>
    """
    return a
```

**Шаги:**
1. Поставить каретку внутри docstring.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются секции (Args, Returns, и т.д.) или параметры.

**Ожидаемый результат:** Completion предлагает секции docstring или параметры функции (зависит от поддержки IDE).

---

### Источник: Раздел 13: Build/config completion (минимальный smoke)

#### Функциональность: Build/config completion

#### [ ] TC-85: package.json — completion по ключам

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/typescript/package.json`

**Описание:** Проверить, что в `package.json` completion предлагает ключи согласно JSON-схеме.

**Предусловие:**
- Файл `package.json` открыт в редакторе.

```json
{
  "<caret>": ""
}
```

**Шаги:**
1. Поставить каретку внутри кавычек первого ключа.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются ключи: `name`, `version`, `dependencies`, `devDependencies` и т.д.

**Ожидаемый результат:** Completion предлагает ключи `package.json` согласно JSON-схеме.

---

#### [ ] TC-86: tsconfig.json — completion по ключам/значениям

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/typescript/tsconfig.json`

**Описание:** Проверить, что в `tsconfig.json` completion предлагает ключи и enum-значения.

**Предусловие:**
- Файл `tsconfig.json` открыт в редакторе.

```json
{
  "compilerOptions": {
    "<caret>": ""
  }
}
```

**Шаги:**
1. Поставить каретку внутри кавычек ключа в `compilerOptions`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются ключи: `target`, `module`, `strict`, `outDir` и т.д.

**Ожидаемый результат:** Completion предлагает ключи `compilerOptions` согласно JSON-схеме TypeScript.

---

#### [ ] TC-87: build.gradle.kts — completion по DSL

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/build.gradle.kts`

**Описание:** Проверить, что в `build.gradle.kts` completion предлагает DSL-элементы Gradle.

**Предусловие:**
- Файл `build.gradle.kts` открыт в редакторе.

```kotlin
plugins {
    <caret>
}

dependencies {
    <caret>
}
```

**Шаги:**
1. Поставить каретку внутри блока `plugins { }`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются элементы DSL (например, `id(...)`, `kotlin(...)`, `java`).
4. Переместить каретку внутри блока `dependencies { }`.
5. Нажать Ctrl+Space.
6. Убедиться, что предлагаются элементы DSL (`implementation(...)`, `testImplementation(...)` и т.д.).

**Ожидаемый результат:** Completion предлагает корректные DSL-элементы в блоках `plugins` и `dependencies`.

---

### Источник: Раздел 16: Refactoring-aware проверки

#### Функциональность: Refactoring-aware проверки

#### [ ] TC-88: Completion после Rename — предлагается новое имя

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/refactoring/RefactoringAware.kt`
- Java: `completion-test-projects/java/src/main/java/completion/TemplatesRefactor.java`

**Описание:** Проверить, что после переименования (Rename) completion предлагает новое имя и не предлагает старое.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun renamedTarget() {}

fun main() {
    ren<caret>
}
```

**Шаги:**
1. Переименовать функцию `renamedTarget` в `newTarget` через Shift+F6.
2. Подтвердить переименование.
3. Поставить каретку после `ren` (или `new`).
4. Нажать Ctrl+Space.
5. Убедиться, что предлагается `newTarget`.
6. Убедиться, что `renamedTarget` **не** предлагается.

**Ожидаемый результат:** Completion предлагает новое имя `newTarget`. Старое имя `renamedTarget` отсутствует.

---

#### [ ] TC-89: Completion после Change Signature — корректные аргументы

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/refactoring/RefactoringAware.kt`

**Описание:** Проверить, что после Change Signature (добавление параметра) completion в аргументах вызова отражает новую сигнатуру.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun foo(a: Int) {}

fun main() {
    foo(<caret>)
}
```

**Шаги:**
1. Выполнить Change Signature для `foo`: добавить параметр `b: String`.
2. Подтвердить изменение.
3. Поставить каретку внутри `foo()`.
4. Нажать Ctrl+Space.
5. Убедиться, что completion предлагает аргументы для обоих параметров (`Int` и `String`).

**Ожидаемый результат:** Completion отражает новую сигнатуру функции с двумя параметрами.

---

### Источник: Раздел 18: Производительность (smoke, без точных метрик)

#### Функциональность: Производительность

#### [ ] TC-90: Completion в «горячей» точке — быстрое появление popup

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/PerfMany.java`

**Описание:** Проверить, что в среднем проекте в классе с 50+ символами в scope popup появляется быстро без фриза UI.

**Предусловие:**
- Проект среднего размера открыт в IDE.
- Файл с классом, содержащим множество методов, открыт в редакторе.

```java
package demo;

public class Many {
  void method000() {}
  void method001() {}
  void method002() {}
  void method003() {}
  void method004() {}
  void method005() {}
  void method006() {}
  void method007() {}
  void method008() {}
  void method009() {}
  void method010() {}
  void method011() {}
  void method012() {}
  void method013() {}
  void method014() {}
  void method015() {}
  void method016() {}
  void method017() {}
  void method018() {}
  void method019() {}
  void method020() {}
  void method021() {}
  void method022() {}
  void method023() {}
  void method024() {}
  void method025() {}
  void method026() {}
  void method027() {}
  void method028() {}
  void method029() {}
  void method030() {}
  void method031() {}
  void method032() {}
  void method033() {}
  void method034() {}
  void method035() {}
  void method036() {}
  void method037() {}
  void method038() {}
  void method039() {}
  void method040() {}
  void method041() {}
  void method042() {}
  void method043() {}
  void method044() {}
  void method045() {}
  void method046() {}
  void method047() {}
  void method048() {}
  void method049() {}

  public static void main(String[] args) {
    Many m = new Many();
    m.met<caret>
  }
}
```

**Шаги:**
1. Поставить каретку после `m.met`.
2. Нажать Ctrl+Space.
3. Убедиться, что popup появляется визуально быстро (без заметного фриза UI).
4. Убедиться, что все методы `method000`–`method049` присутствуют в списке.

**Ожидаемый результат:** Popup появляется быстро, без фриза UI. Все 50 методов доступны в списке.

---

#### [ ] TC-91: 30-кратное открытие/закрытие completion — нет прогрессирующего замедления

**Приоритет:** P1
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/stability/StabilityPerf.kt`

**Описание:** Проверить, что многократное открытие и закрытие completion не приводит к прогрессирующему замедлению.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
package demo

data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.<caret>
}
```

**Шаги:**
1. Поставить каретку после `user.`.
2. Выполнить 30 раз цикл: Ctrl+Space → Esc (открыть/закрыть completion).
3. Убедиться, что не наблюдается прогрессирующего замедления.
4. Убедиться, что последнее открытие так же быстро, как первое.

**Ожидаемый результат:** Нет прогрессирующего замедления. 30-е открытие сопоставимо по скорости с первым.

---

### Источник: Раздел 19: Statement completion — if/else, switch/case/when

#### Функциональность: Completion конструкций if-else

#### [ ] TC-92: Completion для if-else statement — предложение else / else if

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/statements/StatementCompletion.kt`
- Java: `completion-test-projects/java/src/main/java/completion/StatementCompletion.java`

**Описание:** Проверить, что после закрывающей фигурной скобки `}` блока `if` completion предлагает ключевые слова `else` и `else if`.

**Предусловие:**
- IDE открыта с проектом, содержащим файл с приведённым кодом.
- Файл открыт в редакторе.

```kotlin
package completion.statements

import completion.model.User

fun checkAge(user: User) {
    if (user.age > 18) {
        println("adult")
    }
    <caret>
}
```

**Шаги:**
1. Написать блок `if (user.age > 18) { println("adult") }`.
2. Поставить каретку сразу после закрывающей `}` блока `if`.
3. Набрать `el` и вызвать Ctrl+Space.
4. Убедиться, что в списке completion есть `else` и `else if`.
5. Выбрать `else` и убедиться, что вставляется корректная конструкция `else { }`.
6. Повторить для `else if` — убедиться, что вставляется `else if () { }` с кареткой внутри скобок условия.

**Ожидаемый результат:** Completion предлагает `else` и `else if` после закрывающей `}` блока `if`. При выборе вставляется корректная конструкция. Ошибок в IDE log нет.

---

#### Функциональность: Completion конструкций switch/case/default и when

#### [ ] TC-93: Completion для switch/case/default (Java) и when (Kotlin) — предложение ветвей

**Приоритет:** P1
**План:** Regression
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/statements/StatementCompletion.kt`
- Java: `completion-test-projects/java/src/main/java/completion/StatementCompletion.java`

**Описание:** Проверить, что completion предлагает ключевые слова и ветви внутри конструкций `switch` (Java) и `when` (Kotlin), включая `case`, `default` и значения enum.

**Предусловие:**
- IDE открыта с проектом, содержащим файл с приведённым кодом.
- Файл открыт в редакторе.

**Java — switch/case/default:**

```java
package completion;

public class StatementCompletion {

    enum Status { ACTIVE, INACTIVE, BANNED }

    static String describeStatus(Status status) {
        switch (status) {
            case ACTIVE:
                return "active user";
            <caret>
        }
    }
}
```

**Kotlin — when с enum:**

```kotlin
package completion.statements

enum class Status { ACTIVE, INACTIVE, BANNED }

fun describeStatus(status: Status): String {
    return when (status) {
        Status.ACTIVE -> "active user"
        <caret>
    }
}
```

**Шаги:**
1. **(Java)** Открыть файл `StatementCompletion.java`. Внутри блока `switch (status)` после существующего `case ACTIVE:` поставить каретку на новой строке.
2. Набрать `ca` и вызвать Ctrl+Space.
3. Убедиться, что completion предлагает `case` с оставшимися значениями enum (`INACTIVE`, `BANNED`) и `default`.
4. Выбрать `case INACTIVE` — убедиться, что вставляется корректная ветвь.
5. **(Kotlin)** Открыть файл `StatementCompletion.kt`. Внутри блока `when (status)` после существующей ветви `Status.ACTIVE ->` поставить каретку на новой строке.
6. Начать набирать `Status.` или `IN` и вызвать Ctrl+Space.
7. Убедиться, что completion предлагает оставшиеся значения enum (`INACTIVE`, `BANNED`) и `else ->`.
8. Выбрать `Status.INACTIVE` — убедиться, что вставляется корректная ветвь `Status.INACTIVE ->`.

**Ожидаемый результат:** Completion корректно предлагает оставшиеся ветви `case`/enum-значения и `default`/`else`. Уже использованные значения enum занижены в списке или отсутствуют. Ошибок в IDE log нет.

---

### Источник: Раздел 20: Сложные chained/DSL completion сценарии

#### Функциональность: Fluent-chain completion после фабричного вызова

#### [ ] TC-94: Completion в длинной цепочке после `initializeTestContext(...)`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/complex/ComplexDslCompletion.kt`

**Описание:** Проверить, что completion корректно предлагает fluent-методы для продолжения цепочки сразу после результата фабричного вызова.

**Предусловие:**
- IDE открыта с проектом, содержащим файл `ComplexDslCompletion.kt`.
- Файл открыт в редакторе.

```kotlin
val compA = initializeTestContext("testSyncOnNewMachine/compA", TestCases.IU.IntelliJHelloWorld)
    <caret>
```

**Шаги:**
1. В строке с `initializeTestContext(...)` поставить каретку сразу после точки на новой строке.
2. Нажать Ctrl+Space.
3. Убедиться, что completion предлагает методы `copySettingsToIde(...)` и `runIDE(...)`.
4. Выбрать `copySettingsToIde(...)` и убедиться, что выражение собирается в корректную цепочку.

**Ожидаемый результат:** Для типа результата `IdeContext` предлагаются релевантные fluent-методы, вставка выполняется корректно, без поломки синтаксиса.

---

#### Функциональность: Named-argument completion с enum-значениями

#### [ ] TC-95: Completion enum-значений в именованном аргументе `action = ...`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/complex/ComplexDslCompletion.kt`

**Описание:** Проверить type-aware completion в именованном аргументе, где ожидается enum.

**Предусловие:**
- Файл `ComplexDslCompletion.kt` открыт в редакторе.

```kotlin
commandChainSaveAll().enableSettingsSync(action = EnableSettingSyncOptions.<caret>).exitApp()
```

**Шаги:**
1. Поставить каретку после `EnableSettingSyncOptions.`.
2. Нажать Ctrl+Space.
3. Проверить, что предложены `GET` и `PUSH`.
4. Выбрать `GET` и убедиться, что выражение остаётся корректным.

**Ожидаемый результат:** Completion предлагает только значения ожидаемого enum-типа (`GET`, `PUSH`) и корректно вставляет выбранный элемент.

---

#### Функциональность: Completion в nested object пути

#### [ ] TC-96: Completion для констант в `TestCases.IU.*` и `TestCases.WS.*`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/complex/ComplexDslCompletion.kt`

**Описание:** Проверить completion в многоуровневом пути через вложенные объекты.

**Предусловие:**
- Файл `ComplexDslCompletion.kt` открыт в редакторе.

```kotlin
val compA = initializeTestContext("...", TestCases.IU.<caret>)
val compB = initializeTestContext("...", TestCases.WS.<caret>)
```

**Шаги:**
1. Поставить каретку после `TestCases.IU.` и вызвать completion.
2. Убедиться, что предлагается `IntelliJHelloWorld`.
3. Поставить каретку после `TestCases.WS.` и вызвать completion.
4. Убедиться, что предлагается `Empty`.

**Ожидаемый результат:** Completion корректно разрешает nested object контекст и предлагает только релевантные константы текущего уровня.

---

#### Функциональность: Completion после результата chain-вызова

#### [ ] TC-97: Completion методов результата после `runIDE(...)`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/complex/ComplexDslCompletion.kt`

**Описание:** Проверить completion на следующем шаге цепочки, где тип контекста уже сменился (`IdeContext -> RunResult`).

**Предусловие:**
- Файл `ComplexDslCompletion.kt` открыт в редакторе.

```kotlin
compC.runIDE(commands = commandChainSaveAll().enableSettingsSync(action = EnableSettingSyncOptions.GET).exitApp())
    .<caret>
```

**Шаги:**
1. Поставить каретку после точки на следующей строке после `runIDE(...)`.
2. Вызвать completion (Ctrl+Space).
3. Проверить, что в списке есть `assertSettingsEqualTo(...)`.
4. Выбрать `assertSettingsEqualTo(...)` и убедиться, что вызов вставляется корректно.

**Ожидаемый результат:** Completion учитывает изменившийся тип (`RunResult`) и предлагает его API (`assertSettingsEqualTo(...)`) без нерелевантных методов предыдущего контекста.

---

#### Функциональность: Completion внутри промежуточного звена command-chain

#### [ ] TC-98: Completion промежуточных шагов в `commandChainSaveAll().<...>`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/complex/ComplexDslCompletion.kt`

**Описание:** Проверить completion методов в середине DSL-цепочки команд.

**Предусловие:**
- Файл `ComplexDslCompletion.kt` открыт в редакторе.

```kotlin
commandChainSaveAll().<caret>
```

**Шаги:**
1. Поставить каретку после `commandChainSaveAll().`.
2. Нажать Ctrl+Space.
3. Убедиться, что предлагаются `enableSettingsSync(...)` и `exitApp()`.
4. Выбрать `enableSettingsSync(...)`, затем после него снова вызвать completion и выбрать `exitApp()`.

**Ожидаемый результат:** Completion стабильно поддерживает построение многошаговой DSL-цепочки и предлагает валидные следующие шаги на каждом этапе.

---

### Источник: Раздел 21: Completion внутри индексаторов `[]`

#### Функциональность: Completion индексного выражения в списках/массивах

#### [ ] TC-99: Kotlin — completion в вложенной индексации `matrix[row][col]`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/brackets/BracketIndexAccess.kt`

**Описание:** Проверить completion в двух независимых позициях внутри вложенных `[]`.

**Предусловие:**
- Файл `BracketIndexAccess.kt` открыт в редакторе.

```kotlin
val matrix = arrayOf(intArrayOf(1, 2), intArrayOf(3, 4))
val row = 1
val col = 0
val x = matrix[<caret>][<caret>]
```

**Шаги:**
1. Поставить каретку в первые `[]`, вызвать completion, выбрать `row`.
2. Поставить каретку во вторые `[]`, вызвать completion, выбрать `col`.
3. Убедиться, что итоговая конструкция валидна: `matrix[row][col]`.

**Ожидаемый результат:** Completion корректно работает на каждом уровне вложенной индексации и предлагает локальные Int-переменные.

---

#### [ ] TC-100: Kotlin — completion ключа внутри `map[ ... ]`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/brackets/BracketIndexAccess.kt`

**Описание:** Проверить completion выражения ключа внутри `[]` при доступе к `Map`.

**Предусловие:**
- Файл `BracketIndexAccess.kt` открыт в редакторе.

```kotlin
val userByKey = mapOf("name" to "Ann", "age" to "21")
val nameKey = "name"
val v = userByKey[<caret>]
```

**Шаги:**
1. Поставить каретку внутри `userByKey[ ]`.
2. Вызвать completion.
3. Проверить, что предлагается `nameKey`.
4. Выбрать `nameKey` и убедиться, что выражение валидно.

**Ожидаемый результат:** Completion внутри `map[ ]` предлагает релевантные выражения ключа из текущего scope.

---

#### [ ] TC-101: Python — completion внутри `list[ ... ]` и `dict[ ... ]`

**Приоритет:** P1
**План:** Regression
**Тестовый файл:**
- Python: `completion-test-projects/python/bracket_index.py`

**Описание:** Проверить устойчивость и релевантность completion внутри индексаторов `[]` в Python.

**Предусловие:**
- Файл `bracket_index.py` открыт в редакторе.

```python
users = [{"name": "Ann", "age": 21}, {"name": "Bob", "age": 22}]
idx = 0
a = users[<caret>]
b = users[idx]["<caret>"]
```

**Шаги:**
1. Поставить каретку внутри `users[ ]`, вызвать completion и выбрать `idx`.
2. Поставить каретку внутри `users[idx][""]`, ввести `na`, вызвать completion.
3. Проверить, что completion не приводит к ошибкам IDE и вставляет валидный синтаксис.

**Ожидаемый результат:** Completion стабильно работает внутри `[]` в Python-контекстах списка и словаря.

---

### Источник: Раздел 22: Продвинутые контексты completion

#### Функциональность: Type-flow, DSL, refactoring, устойчивость

#### [ ] TC-102: Smart-cast completion после `is`-проверки

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что после `if (x is User)` completion предлагает члены `User` для `x.`.

**Шаги:**
1. В функции `smartCastAfterIs` поставить каретку после `x.`.
2. Вызвать completion.
3. Убедиться, что есть `name`, `age` и другие члены `User`.

**Ожидаемый результат:** Completion учитывает smart-cast и предлагает API конкретного типа.

---

#### [ ] TC-103: Completion для implicit `it` в lambda

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что в `users.map { it.<caret> }` completion правильно типизирует `it`.

**Шаги:**
1. В `lambdaItTypedCompletion` поставить каретку после `it.`.
2. Вызвать completion.
3. Проверить, что предлагаются `name`, `age`.

**Ожидаемый результат:** `it` определяется как `User`; список completion релевантен типу `User`.

---

#### [ ] TC-104: Completion внутри string template `${...}`

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить completion внутри интерполяции строки.

**Шаги:**
1. В `completionInsideStringTemplate` поставить каретку после `user.` в `"${user.<caret>}"`.
2. Вызвать completion.
3. Убедиться, что доступны члены `User`.

**Ожидаемый результат:** Completion внутри `${...}` эквивалентен обычному expression-контексту.

---

#### [ ] TC-105: Override/implements completion (template methods)

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/AdvancedCompletionContracts.java`

**Описание:** Проверить, что completion/template по `over` предлагает корректные `@Override`-заготовки.

**Шаги:**
1. В `LocalWorker` поставить каретку в теле класса.
2. Набрать `over`, вызвать completion.
3. Выбрать override для `runSync()` или `workerId()`.

**Ожидаемый результат:** Генерируется корректный override-метод с валидной сигнатурой.

---

#### [ ] TC-106: Completion ветвей `when` для sealed-типа

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что `when (state)` предлагает/подсказывает все ветви sealed-иерархии.

**Шаги:**
1. В `sealedWhenBranches` удалить одну из ветвей.
2. Внутри `when` вызвать completion/quick-fix для missing branches.
3. Проверить, что предлагаются оставшиеся варианты `Idle`, `Running`, `Failed`.

**Ожидаемый результат:** Completion/IDE action корректно подсказывает недостающие sealed-ветви.

---

#### [ ] TC-107: Auto-import extension function из другого package

**Приоритет:** P1  
**План:** Regression
**Тестовые файлы:**
- `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/ExtensionAutoImportScenario.kt`
- `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/external/ExternalExtensions.kt`

**Описание:** Проверить, что extension-функция предлагается после `.` и может быть auto-import'нута.

**Шаги:**
1. В `ExtensionAutoImportScenario.kt` удалить import `syncNow`.
2. Поставить каретку после `settings.`.
3. Набрать `sync`, вызвать completion и выбрать `syncNow`.

**Ожидаемый результат:** Extension-функция доступна в списке; при выборе добавляется корректный import.

---

#### [ ] TC-108: Completion при перегрузках — релевантная сигнатура

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/java/src/main/java/completion/AdvancedCompletionContracts.java`

**Описание:** Проверить, что в вызовах перегруженной функции completion учитывает ожидаемую сигнатуру.

**Шаги:**
1. В `overloadResolution` поставить каретку внутри `pick(...)` для строкового вызова.
2. Вызвать completion и выбрать `str`.
3. Повторить для числового вызова и выбрать `num`.

**Ожидаемый результат:** Для `pick(String)` ранжируются string-аргументы, для `pick(int)` — numeric.

---

#### [ ] TC-109: Completion в generic-bound контексте (`T : CharSequence`)

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что completion в `boundedEcho(...)` подсказывает значения, совместимые с bound.

**Шаги:**
1. В `genericBoundContext` поставить каретку внутри `boundedEcho(...)`.
2. Вызвать completion.
3. Проверить, что релевантно предлагается `text`, а не только `Int`-переменные.

**Ожидаемый результат:** Completion учитывает ограничение `T : CharSequence`.

---

#### [ ] TC-110: Completion после safe-call + Elvis (`?:`)

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить expected-type completion после `?:`.

**Шаги:**
1. В `safeCallElvisExpectedType` поставить каретку после `?:`.
2. Вызвать completion.
3. Проверить, что предложены `String`-совместимые выражения (например, `fallback`).

**Ожидаемый результат:** Completion корректно использует ожидаемый тип RHS выражения.

---

#### [ ] TC-111: Completion в DSL receiver (`configureSettings { ... }`)

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что внутри DSL-блока предлагаются члены receiver-объекта.

**Шаги:**
1. В `dslReceiverContext` поставить каретку внутри блока `configureSettings { }`.
2. Вызвать completion.
3. Убедиться, что предлагаются `enabled`, `profileName`, `save()`.

**Ожидаемый результат:** Completion показывает API текущего DSL receiver.

---

#### [ ] TC-112: Межмодульный/межпакетный символ — completion и import

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что символы из внешнего пакета/модуля проекта корректно предлагаются и резолвятся.

**Шаги:**
1. В `crossPackageDependencyLikeModule` удалить использование `buildUser` и/или import.
2. Начать набирать `buildUs...`, вызвать completion.
3. Выбрать элемент и проверить добавление import/resolution.

**Ожидаемый результат:** Символы из «внешней» части проекта корректно находятся и подключаются.

---

#### [ ] TC-113: Completion после Move/Rename между файлами

**Приоритет:** P1  
**План:** Regression
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/advanced/AdvancedCompletionScenarios.kt`

**Описание:** Проверить, что completion отражает новое имя после rename/move API.

**Шаги:**
1. Переименовать `renamedApiTarget` через Refactor > Rename.
2. В `refactorMoveRenameFollowup` удалить вызов.
3. Начать набирать новое имя и вызвать completion.

**Ожидаемый результат:** Completion предлагает актуальное имя, старое не приоритизируется.

---

## План 3: Full (P2)

Тесты этого плана выполняются в порядке `TC-114` → `TC-121`.

### Источник: Раздел 14: Templates — postfix / live templates

#### Функциональность: Postfix и Live templates

#### [ ] TC-114: Postfix completion `.if` → if (...) { }

**Приоритет:** P2
**План:** Full
**Тестовые файлы:**
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/templates/TemplatesInject.kt`
- Java: `completion-test-projects/java/src/main/java/completion/TemplatesRefactor.java`

**Описание:** Проверить, что postfix completion `.if` превращает выражение в конструкцию `if (...) { }`.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val ok = true
    ok.if<caret>
}
```

**Шаги:**
1. Поставить каретку после `ok.if`.
2. Нажать Ctrl+Space или Tab (в зависимости от настроек).
3. Принять postfix completion.
4. Убедиться, что результат: `if (ok) { }`.

**Ожидаемый результат:** Выражение `ok.if` превращается в `if (ok) { }`. Каретка — внутри фигурных скобок.

---

#### [ ] TC-115: Конфликт basic completion и postfix — оба типа в списке

**Приоритет:** P2
**План:** Full
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/templates/TemplatesInject.kt`

**Описание:** Проверить, что при конфликте между basic completion и postfix completion оба типа присутствуют в списке с корректным приоритетом.

**Предусловие:**
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val ok = true
    ok.if<caret>
}
```

**Шаги:**
1. Поставить каретку после `ok.if`.
2. Нажать Ctrl+Space.
3. Убедиться, что в списке есть и postfix `.if`, и обычные completion-элементы (если есть совпадения).
4. Убедиться, что приоритет корректный (postfix `.if` — в ожидаемом месте списка).

**Ожидаемый результат:** Оба типа предложений присутствуют в списке. Приоритет корректный.

---

#### [ ] TC-116: Live template — sout + Tab → expand с плейсхолдерами

**Приоритет:** P2
**План:** Full
**Тестовые файлы:**
- Java: `completion-test-projects/java/src/main/java/completion/TemplatesRefactor.java`
- Kotlin: `completion-test-projects/kotlin/src/main/kotlin/completion/templates/TemplatesInject.kt`

**Описание:** Проверить, что live template `sout` раскрывается по Tab с корректными плейсхолдерами.

**Предусловие:**
- Файл с кодом открыт в редакторе (Java или Kotlin).

```java
package demo;

public class T {
  public static void main(String[] args) {
    sout<caret>
  }
}
```

**Шаги:**
1. Набрать `sout` внутри метода/функции.
2. Нажать Tab.
3. Убедиться, что `sout` раскрыт в `System.out.println()` (Java) или `println()` (Kotlin).
4. Убедиться, что каретка прыгает по плейсхолдерам (внутрь скобок).

**Ожидаемый результат:** `sout` раскрывается в корректный print-вызов для языка файла. Каретка — внутри скобок для ввода аргумента.

---

### Источник: Раздел 15: Injected languages / SQL / regex

#### Функциональность: Injected languages

#### [ ] TC-117: SQL completion в строке с инъекцией

**Приоритет:** P2
**План:** Full
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/templates/TemplatesInject.kt`

**Описание:** Проверить, что в строке с SQL-инъекцией (при включённой поддержке) completion предлагает SQL keywords.

**Предусловие:**
- IDE с включённой поддержкой language injection.
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val q = "SELECT <caret> FROM users WHERE id = 1"
}
```

**Шаги:**
1. Поставить каретку после `SELECT ` внутри строки.
2. Нажать Ctrl+Space.
3. Убедиться, что completion предлагает SQL keywords (`*`, имена столбцов, `DISTINCT` и т.д.).

**Ожидаемый результат:** Completion предлагает SQL keywords и элементы.

---

#### [ ] TC-118: Regex completion в строке

**Приоритет:** P2
**План:** Full
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/templates/TemplatesInject.kt`

**Описание:** Проверить, что в regex-строке completion предлагает классы символов и escape-последовательности (если поддерживается).

**Предусловие:**
- IDE с поддержкой regex completion.
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val r = "\\d<caret>"
}
```

**Шаги:**
1. Поставить каретку после `\\d` внутри строки.
2. Нажать Ctrl+Space.
3. Убедиться, что completion предлагает regex-классы/escape-последовательности (если фича поддерживается).

**Ожидаемый результат:** Completion предлагает regex-элементы (если поддерживается). IDE не падает.

---

#### [ ] TC-119: Переход из injected в host — корректное закрытие popup

**Приоритет:** P2
**План:** Full
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/templates/TemplatesInject.kt`

**Описание:** Проверить, что при перемещении каретки из injected language в host popup корректно закрывается/обновляется.

**Предусловие:**
- IDE с включённой language injection.
- Файл с кодом открыт в редакторе.

```kotlin
fun main() {
    val q = "SELECT <caret> FROM users WHERE id = 1"
    val x = 42
}
```

**Шаги:**
1. Поставить каретку внутри SQL-строки.
2. Нажать Ctrl+Space — появляется SQL completion.
3. Переместить каретку за пределы строки (на строку `val x = 42`).
4. Убедиться, что popup корректно закрывается или обновляется.

**Ожидаемый результат:** Popup корректно закрывается при выходе из injected-контекста. Нет артефактов.

---

### Источник: Раздел 22: Продвинутые контексты completion

#### Функциональность: Type-flow, DSL, refactoring, устойчивость

#### [ ] TC-120: Completion в большом файле (1000+ строк) и local scope

**Приоритет:** P2  
**План:** Full
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/stability/LargeFileScope.kt`

**Описание:** Проверить latency и релевантность completion в очень большом локальном контексте.

**Шаги:**
1. Открыть `LargeFileScope.kt`.
2. Поставить каретку на позицию `TC-120`, набрать `target`.
3. Вызвать completion и выбрать `targetUserName`.

**Ожидаемый результат:** Popup появляется без заметных лагов, локальная переменная ранжируется высоко.

---

#### [ ] TC-121: Commit behavior для дополнительных клавиш (`;`, `Space`)

**Приоритет:** P2  
**План:** Full
**Тестовый файл:** `completion-test-projects/kotlin/src/main/kotlin/completion/accept/CommitVariants.kt`

**Описание:** Проверить commit completion для менее типовых commit-символов.

**Шаги:**
1. В `commitBySemicolonAndSpace` в первом кейсе оставить `text.up`, вызвать completion.
2. Нажать `;` и проверить, что элемент принят.
3. Во втором кейсе оставить `num.toS`, вызвать completion.
4. Нажать `Space` и проверить поведение (зависит от IDE-настроек commit on space).

**Ожидаемый результат:** Для `;` commit работает предсказуемо; для `Space` поведение соответствует настройкам IDE и не ломает код.

---

## Проверка полноты покрытия

### Сводная таблица покрытия

| # Раздел чеклиста | Пункты в чеклисте | Тест-кейсы | Покрытие |
|---|---|---|---|
| 1 Триггеры и жизненный цикл completion popup | 15 | TC-1 — TC-14, TC-57 | 15/15 ✅ |
| 2 Basic completion — корректность списка и ранжирование | 10 | TC-15 — TC-21, TC-58 — TC-60 | 10/10 ✅ |
| 3 Member completion (после доступа к членам) | 8 | TC-22 — TC-27, TC-61 — TC-62 | 8/8 ✅ |
| 4 Smart completion (тип-ориентированная) | 5 | TC-28 — TC-32 | 5/5 ✅ |
| 5 Completion в аргументах/параметрах и сигнатурах | 7 | TC-33 — TC-35, TC-63 — TC-66 | 7/7 ✅ |
| 6 Accept/commit — вставка, замены, commit characters | 8 | TC-36 — TC-42, TC-67 | 8/8 ✅ |
| 7 Auto-import и символы из зависимостей | 6 | TC-43 — TC-46, TC-68 — TC-69 | 6/6 ✅ |
| 8 Completion в строках (paths/resources) | 4 | TC-47 — TC-49, TC-70 | 4/4 ✅ |
| 9 Documentation popup / QuickDoc / детали элемента | 3 | TC-71 — TC-73 | 3/3 ✅ |
| 10 Generics / type parameters | 3 | TC-74 — TC-76 | 3/3 ✅ |
| 11 Аннотации / атрибуты / декораторы | 4 | TC-77 — TC-80 | 4/4 ✅ |
| 12 Doc/Comments completion | 4 | TC-81 — TC-84 | 4/4 ✅ |
| 13 Build/config completion (минимальный smoke) | 3 | TC-85 — TC-87 | 3/3 ✅ |
| 14 Templates — postfix / live templates | 3 | TC-114 — TC-116 | 3/3 ✅ |
| 15 Injected languages / SQL / regex | 3 | TC-117 — TC-119 | 3/3 ✅ |
| 16 Refactoring-aware проверки | 2 | TC-88 — TC-89 | 2/2 ✅ |
| 17 Dumb mode / индексация / устойчивость | 3 | TC-50 — TC-52 | 3/3 ✅ |
| 18 Производительность (smoke, без точных метрик) | 2 | TC-90 — TC-91 | 2/2 ✅ |
| 19 Statement completion — if/else, switch/case/when | 2 | TC-92 — TC-93 | 2/2 ✅ |
| 20 Сложные chained/DSL completion сценарии | 5 | TC-94 — TC-98 | 5/5 ✅ |
| 21 Completion внутри индексаторов `[]` | 6 | TC-53 — TC-55, TC-99 — TC-101 | 6/6 ✅ |
| 22 Продвинутые контексты completion | 15 | TC-56, TC-102 — TC-113, TC-120 — TC-121 | 15/15 ✅ |
| **ИТОГО** | **121** | **TC-1 — TC-121** | **121/121 ✅ (100%)** |

### Дополнительная проверка полноты

**Общая проверка (из чеклиста):**
> «Во всех пунктах: дополнительно проверяем отсутствие ошибок в IDE log (exceptions/assertions) при открытии/обновлении completion.»

Это сквозное требование учтено в ожидаемых результатах всех тест-кейсов, где оно применимо (особенно TC-21, TC-32, TC-50, TC-51).

**Распределение по приоритетам:**
| Приоритет | Количество тест-кейсов |
|---|---|
| P0 | 56 |
| P1 | 57 |
| P2 | 8 |
| **Итого** | **121** |

**Распределение по планам выполнения:**
| План | Количество тест-кейсов | Диапазон TC |
|---|---|---|
| Acceptance | 56 | TC-1 — TC-56 |
| Regression | 57 | TC-57 — TC-113 |
| Full | 8 | TC-114 — TC-121 |
| **Итого** | **121** | **TC-1 — TC-121** |

**Распределение по языкам (основной код примера):**
| Язык | Тест-кейсы |
|---|---|
| Kotlin | 82 |
| Java | 25 |
| TypeScript | 6 |
| Python | 5 |
| JSON/Config | 3 |

**Покрытие мультиязычности:**
- Многие пункты чеклиста помечены несколькими примерами ([EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]). В тест-кейсах приведён основной пример на одном языке, но каждый тест-кейс следует выполнять на всех применимых языках.
- Рекомендуется параметризовать тест-кейсы по языкам для полного покрытия.
