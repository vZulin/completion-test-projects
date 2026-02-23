Иерархический чеклист тестирования Completion (после рефакторинга)
Нотация:

[P0/P1/P2] — приоритет.
[EX-…] — ссылка на пример кода в конце (можно копировать в файл и ставить caret в указанные места).
Во всех пунктах: дополнительно проверяем отсутствие ошибок в IDE log (exceptions/assertions) при открытии/обновлении completion.
1) Триггеры и жизненный цикл completion popup
1.1 Manual trigger (вручную)
[P0] Нажать Ctrl+Space в месте <caret> → открывается список completion. [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Нажать Ctrl+Shift+Space в месте <caret> → открывается smart completion. [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Открыть completion, не выбирая элемент, повторно нажать Ctrl+Space → список обновляется/расширяется (если в IDE есть 2-й уровень) или остаётся консистентным без “мигания”. [EX-KT-1]
[P0] Открыть completion и нажать Esc → popup закрывается, текст не меняется. [EX-KT-1]
[P0] Открыть completion и кликнуть мышью вне popup → popup закрывается. [EX-KT-1]
[P0] Открыть completion и продолжить вводить текст → список фильтруется без зависаний. [EX-KT-1]
1.2 Auto-popup (автоматическое появление)
(Проверка: popup появляется сам, без Ctrl+Space)

[P0] Ввести . после выражения/переменной → появляется member completion. [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Ввести ( после имени функции/метода → появляется completion в аргументах (если ожидается). [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Ввести , внутри списка аргументов → появляется completion для следующего аргумента. [EX-KT-1], [EX-JV-1], [EX-TS-1]
[P1] Ввести @ (аннотация/декоратор) → появляется completion по аннотациям/декораторам. [EX-KT-2], [EX-JV-2], [EX-PY-2]
[P0] В import ... начать вводить путь "./ → появляется completion путей. [EX-TS-2]
[P0] В строке пути File(" / Path(" начать ввод → появляется completion файлов/папок. [EX-KT-1], [EX-PY-1], [EX-CS-1]
1.3 Обновление/отмена запросов (race/cancel)
[P0] Открыть completion, быстро набрать 5–10 символов → список меняется только по актуальному вводу, нет “отката” к старым результатам. [EX-KT-1]
[P0] Открыть completion, затем быстро удалить введённый префикс Backspace → список корректно возвращается к общему набору, без пустых/битых состояний. [EX-KT-1]
[P0] Открыть completion и сразу переместить caret стрелками в другое место (без выбора) → popup закрывается или корректно переоткрывается по новому месту (без артефактов). [EX-KT-1]
2) Basic completion: корректность списка и ранжирование
2.1 Локальные символы / scope
[P0] Внутри функции вызвать completion на префиксе локальной переменной → локальная переменная присутствует в списке. [EX-KT-1]
[P0] Внутри функции вызвать completion на префиксе параметра функции → параметр присутствует в списке. [EX-KT-1], [EX-JV-1], [EX-TS-1]
[P0] Внутри nested scope (например if/for) completion показывает символы из внешнего scope. [EX-KT-3]
[P0] Проверить, что символы вне области видимости не предлагаются (например private в другом классе, где недоступно). [EX-JV-3]
2.2 Ключевые слова
[P0] Набрать ret и вызвать completion → в списке есть return, выбор вставляет return корректно. [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Набрать cla и вызвать completion в top-level → есть class. [EX-KT-4], [EX-JV-4], [EX-TS-3], [EX-PY-3]
2.3 Ранжирование (relevance)
[P0] При одинаковом префиксе локальная переменная должна быть выше глобальной/импортированной. [EX-KT-5]
[P0] При вводе точного префикса userN элемент userName выше fuzzy совпадений. [EX-KT-6]
[P0] Если IDE поддерживает MRU: выбрать один и тот же элемент 3 раза, затем снова вызвать completion → выбранный ранее элемент поднимается, но не ломает релевантность (не становится №1 в нерелевантном контексте). [EX-KT-1]
2.4 Негативные контексты (устойчивость)
[P0] Вызвать completion в заведомо “сломанном” месте (например посреди числа/токена) → IDE не падает, нет исключений; popup либо не показывается, либо показывает предсказуемый минимум. [EX-KT-7]
3) Member completion (после доступа к членам)
3.1 После 
.
[P0] user. → в списке есть свойства/методы типа User. [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Ввести после user. префикс na → список фильтруется до name/getName (по языку). [EX-KT-1], [EX-JV-1]
[P0] Выбрать метод из списка → вставляются скобки () если метод требует вызова; caret оказывается внутри (). [EX-JV-1], [EX-TS-1]
[P0] Выбрать поле/свойство → вставляется только имя без (). [EX-KT-1], [EX-TS-1]
3.2 Статические/companion (где применимо)
[P0] Math./User. → предлагаются доступные статические члены/companion. [EX-JV-5], [EX-KT-8]
[P0] Проверить, что нестатические члены не появляются в ClassName. контексте. [EX-JV-5]
3.3 Nullable/safe access (Kotlin)
[P1] nullableUser?. → показываются члены типа, применимые к safe-call. [EX-KT-9]
3.4 Extension methods (Kotlin)
[P1] На типе, для которого есть extension, str. → extension функция в списке. [EX-KT-10]
4) Smart completion (тип-ориентированная)
4.1 Присваивание
[P0] В месте val u: User = <caret> вызвать smart completion → в списке есть User(...), buildUser(...) и переменные типа User. [EX-KT-1], [EX-JV-1], [EX-TS-1], [EX-PY-1]
[P0] Выбрать buildUser(...) → вставляется вызов, caret внутри () и/или подсказка параметров появляется. [EX-KT-1]
4.2 Return
[P0] В функции с возвращаемым типом User на return <caret> вызвать smart completion → предлагаются выражения типа User. [EX-KT-11], [EX-JV-6], [EX-TS-4]
4.3 Аргументы функции (по ожидаемому типу)
[P0] Вызов consumeUser(<caret>) вызвать smart completion → предлагаются переменные типа User. [EX-KT-12], [EX-JV-7], [EX-TS-5]
4.4 Fallback при неизвестном типе
[P1] В динамическом/неразрешимом месте вызвать smart completion → нет ошибок, список разумный (может быть как basic). [EX-PY-4], [EX-TS-6]
5) Completion в аргументах/параметрах и сигнатурах
5.1 Подстановка аргументов
[P0] Внутри buildUser(<caret>) вызвать completion → предлагаются переменные String для name и Int для age (проверять по очереди на каждом параметре). [EX-KT-1]
[P0] После вставки первого аргумента набрать , → появляется completion для второго аргумента. [EX-KT-1], [EX-JV-1], [EX-TS-1]
[P0] Вызвать completion сразу после ( без ввода → список не пустой (если есть подходящие переменные), нет мусора. [EX-KT-1]
5.2 Parameter info интеграция (если включена)
[P1] При выборе функции из completion проверить, что показывается parameter info и подсветка текущего параметра при вводе. [EX-KT-1], [EX-JV-1], [EX-TS-1]
5.3 Named arguments (Kotlin, Python partial)
[P1] В Kotlin внутри buildUser(<caret>) вызвать completion → есть варианты name = и age =. [EX-KT-13]
[P1] В Kotlin после уже введённого name = ... completion не предлагает повторно name = или помечает его как уже использованный (зависит от поведения IDE). [EX-KT-13]
[P1] В Python для функции с именованными параметрами completion предлагает имена параметров при вводе func(<caret>) (если поддерживается). [EX-PY-5]
6) Accept/commit: вставка, замены, commit characters
6.1 Enter / Tab
[P0] Открыть completion, выбрать элемент стрелками, нажать Enter → элемент вставляется. [EX-KT-1]
[P0] Открыть completion, выбрать элемент, нажать Tab → элемент вставляется (или выполняется “replace common prefix” — согласно настройкам), поведение соответствует настройке IDE. [EX-KT-1]
[P0] Выделить текст в редакторе, вызвать completion и принять элемент → выделенный текст заменён полностью. [EX-KT-14]
6.2 Commit characters
[P0] Ввести user.na и принять completion нажатием . → вставляется name и затем . остаётся/добавляется корректно (в итоге user.name.). [EX-KT-15]
[P0] Вызвать completion на имени метода и принять нажатием ( → вставляется метод + () и открывающая ( не дублируется. [EX-JV-8], [EX-TS-7]
[P1] Принять completion нажатием , внутри списка аргументов → вставка корректная, запятая не дублируется, пробелы по code style. [EX-KT-1]
6.3 Caret placement
[P0] При выборе функции из completion caret становится внутри () для ввода аргументов. [EX-KT-1], [EX-JV-1], [EX-TS-1]
[P0] При выборе конструктора/инициализации caret становится в ожидаемом месте (внутри скобок/после имени). [EX-JV-1], [EX-KT-1]
7) Auto-import и символы из зависимостей
7.1 Импорт при выборе из completion
[P0] В файле без import начать ввод ArrayList и принять completion → добавляется import java.util.ArrayList; (или используется fully qualified — по настройке), код компилируется. [EX-JV-9]
[P0] В Kotlin без import набрать File и принять completion → добавляется import java.io.File. [EX-KT-16]
[P0] В TS принять completion по неимпортированному символу из другого модуля → добавляется import строка (если IDE так умеет в данном проекте). [EX-TS-8]
7.2 Конфликт имён (две сущности с одним именем)
[P0] Ввести имя типа, которое есть в двух пакетах/модулях → completion показывает различия (package/module) и при выборе импортирует правильный. [EX-JV-10], [EX-KT-17]
[P1] После выбора одного варианта второй не подмешивается в imports автоматически. [EX-JV-10]
7.3 Поведение в зависимости от настроек
[P1] Выключить auto-import в настройках → принять элемент из completion, import не добавляется (или добавляется fully qualified) согласно настройке. [EX-JV-9]
8) Completion в строках (paths/resources)
8.1 File system paths
[P0] В File("<caret>") начать ввод src/ → показываются папки/файлы проекта. [EX-KT-1], [EX-CS-1]
[P0] Выбрать файл из списка → строка дополняется корректно, кавычки не ломаются, caret остаётся в строке. [EX-KT-1]
[P0] Проверить completion относительных путей ./ и ../ (если поддерживается) → список соответствует файловой структуре. [EX-TS-2]
8.2 Негатив: обычная строка
[P1] В строковом литерале, не похожем на путь, вызвать completion → не должно навязывать file path completion (если продукт так задуман), либо поведение стабильное и повторяемое. [EX-KT-18]
9) Documentation popup / QuickDoc / детали элемента
[P0] Открыть completion, выбрать элемент, нажать Ctrl+Q → показывается QuickDoc для выбранного элемента. [EX-KT-1], [EX-JV-1]
[P0] В completion списке корректно показывается тип/сигнатура (например (String, Int) -> User). [EX-KT-1]
[P1] Deprecated элементы помечены как deprecated (визуально/иконкой/стилем). [EX-JV-11], [EX-KT-19]
10) Generics / type parameters
[P1] В List<<caret>> вызвать completion → предлагаются типы из scope/пакетов. [EX-JV-12], [EX-KT-20], [EX-TS-9]
[P1] В Map<String, <caret>> → предлагаются типы, корректная вставка > и запятых. [EX-JV-12]
[P1] Выбор типа из completion добавляет import при необходимости. [EX-JV-12], [EX-KT-20]
11) Аннотации / атрибуты / декораторы
[P1] В Java набрать @Depr<caret> → completion предлагает @Deprecated. [EX-JV-2]
[P1] В Kotlin набрать @Depre<caret> → completion предлагает @Deprecated. [EX-KT-2]
[P1] В Python набрать @datac<caret> → completion предлагает @dataclass. [EX-PY-2]
[P1] В аннотации с параметрами вызвать completion внутри скобок → предлагаются имена параметров/значения по типу (если применимо). [EX-KT-21], [EX-JV-12]
12) Doc/Comments completion
[P1] В Javadoc после /** вызвать completion → предлагаются теги @param, @return, @throws. [EX-JV-13]
[P1] В Javadoc выбрать @param → вставляется шаблон с именем параметра (или предлагает выбрать параметр). [EX-JV-13]
[P1] В KDoc после /** completion предлагает @param/@return. [EX-KT-22]
[P1] В Python docstring completion предлагает секции/параметры (если IDE поддерживает). [EX-PY-6]
13) Build/config completion (минимальный smoke)
[P1] В package.json (если проект web) completion по ключам (например "dependencies") работает согласно schema. [EX-CONFIG-1]
[P1] В tsconfig.json completion по ключам/enum значениям работает. [EX-CONFIG-2]
[P1] В build.gradle.kts completion по DSL (plugins/dependencies) работает. [EX-CONFIG-3]
14) Templates: postfix / live templates (если включено в продукт)
[P2] Набрать выражение и postfix .if → появляется postfix completion, принятие превращает в if (...) {}. [EX-KT-23], [EX-JV-14]
[P2] При конфликте между basic completion и postfix completion список содержит оба типа, приоритет корректный. [EX-KT-23]
[P2] Live template: набрать сокращение (например sout) → expand по Tab, caret прыгает по плейсхолдерам. [EX-JV-15]
15) Injected languages / SQL / regex (если поддерживается в IDE)
[P2] В строке с SQL (с включённой инъекцией) completion предлагает SQL keywords. [EX-INJECT-1]
[P2] В regex строке completion предлагает классы/escape последовательности (если фича есть). [EX-INJECT-2]
[P2] Переместить caret из injected в host → popup корректно закрывается/обновляется. [EX-INJECT-1]
16) Refactoring-aware проверки
[P2] Выполнить Rename для класса/функции, затем вызвать completion в месте использования → предлагается новое имя, старое не предлагается. [EX-JV-16], [EX-KT-24]
[P2] После Change Signature (добавить параметр) вызвать completion в аргументах вызова → предлагаются корректные аргументы/имена параметров. [EX-KT-25]
17) Dumb mode / индексация / устойчивость
[P0] Во время индексации (Dumb mode) вызвать completion → IDE не падает; либо показывает ограниченный список, либо информирует, что результаты ограничены. [EX-KT-1]
[P0] Во время индексации вызвать completion 5 раз подряд → нет зависаний UI, нет очереди “старых” попапов. [EX-KT-1]
[P1] После завершения индексации повторить completion → результаты становятся полными и корректными. [EX-KT-1]
18) Производительность (smoke, без точных метрик)
[P0] В среднем проекте открыть completion в “горячей” точке (класс с 50+ символами в scope) → popup появляется визуально быстро, без фриза UI. [EX-JV-17]
[P1] 30 раз подряд открыть/закрыть completion → не наблюдается прогрессирующего замедления (признак утечки/накопления). [EX-KT-1]
Примеры кода (в конец). В примерах указаны места для 
<caret>
[EX-KT-1] Kotlin: базовый комбайн (member/smart/args/paths/UI)
package demo

import java.io.File

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)
fun consumeUser(u: User) {}

fun main() {
    val userName = "Ann"
    val userAge = 21
    val user = User(userName, userAge)

    user.<caret>                 // member completion + filtering + accept

    val u1: User = <caret>       // smart completion: user / User(...) / buildUser(...)
    val u2: User = buildUser(<caret>) // args completion (String), then after comma -> Int
    consumeUser(<caret>)         // smart completion by expected type

    File("<caret>")              // path completion inside string
}

[EX-KT-2] Kotlin: аннотации
package demo

@Depre<caret>
class A

[EX-KT-3] Kotlin: nested scope
fun main() {
    val outerValue = 10
    if (true) {
        out<caret> // should suggest outerValue
    }
}

[EX-KT-4] Kotlin: top-level keywords
cla<caret>

[EX-KT-5] Kotlin: ранжирование local vs global
package demo

val globalUser = User("Global", 0)

data class User(val name: String, val age: Int)

fun main() {
    val globalUserName = "LocalShadow"
    glo<caret> // should prefer local globalUserName / locals first (depending on exact prefix)
}

[EX-KT-6] Kotlin: точный префикс vs fuzzy
fun main() {
    val userName = "Ann"
    val usernameLegacy = "legacy"
    userN<caret> // userName should rank above usernameLegacy / fuzzy results
}

[EX-KT-7] Kotlin: негативный контекст (устойчивость)
fun main() {
    val x = 1234
    val y = 12<caret>34 // invoke completion here: should not crash
}

[EX-KT-8] Kotlin: companion/static-like
class Util {
    companion object {
        fun make(): Util = Util()
    }
}

fun main() {
    Util.<caret> // should suggest make()
}

[EX-KT-9] Kotlin: nullable safe-call
data class User(val name: String, val age: Int)

fun main() {
    val u: User? = null
    u?.<caret> // should suggest name/age/etc.
}

[EX-KT-10] Kotlin: extension function
fun String.extHello(): String = "Hello, $this"

fun main() {
    val s = "Ann"
    s.<caret> // should suggest extHello()
}

[EX-KT-11] Kotlin: smart completion for return
data class User(val name: String, val age: Int)

fun buildUser(): User = User("Ann", 21)

fun f(): User {
    return <caret> // smart completion should suggest buildUser(), User(...)
}

[EX-KT-12] Kotlin: expected type in args
data class User(val name: String, val age: Int)
fun consumeUser(u: User) {}

fun main() {
    val u = User("Ann", 21)
    consumeUser(<caret>) // should suggest u
}

[EX-KT-13] Kotlin: named arguments
data class User(val name: String, val age: Int)
fun buildUser(name: String, age: Int): User = User(name, age)

fun main() {
    buildUser(<caret>) // should suggest "name =" and "age ="
    buildUser(name = "Ann", <caret>) // should not re-suggest name=
}

[EX-KT-14] Kotlin: replace selection
fun main() {
    val userName = "Ann"
    val x = "REPLACE_ME"
    // select REPLACE_ME then completion to insert userName
    val y = "REPLACE_ME"<caret>
}

[EX-KT-15] Kotlin: commit by dot
data class User(val name: String, val age: Int)

fun main() {
    val user = User("Ann", 21)
    user.na<caret> // accept completion by pressing '.' -> user.name.
}

[EX-KT-16] Kotlin: auto-import File
package demo

fun main() {
    val f = File("<caret>") // File is not imported initially; accept completion should add import java.io.File
}

[EX-KT-17] Kotlin: import conflict (нужно подготовить зависимости/пакеты)
package demo

fun main() {
    Date<caret> // if both java.util.Date and java.sql.Date are available, completion should show both with qualifiers
}

[EX-KT-18] Kotlin: “обычная строка” (негатив для path completion)
fun main() {
    val s = "hello <caret> world" // completion here should be stable; path completion should not aggressively appear
}

[EX-KT-19] Kotlin: deprecated marker
@Deprecated("old")
fun oldFun() {}

fun main() {
    old<caret> // completion should show deprecated style for oldFun
}

[EX-KT-20] Kotlin: generics
package demo

import java.util.ArrayList

fun main() {
    val x: ArrayList<<caret>> = arrayListOf()
}

[EX-KT-21] Kotlin: annotation params completion
annotation class Ann(val value: String, val enabled: Boolean)

@Ann(<caret>) // completion should suggest value = , enabled =
class A

[EX-KT-22] Kotlin: KDoc tags
/**
 * <caret>
 */
fun f(a: Int): Int = a

[EX-KT-23] Kotlin: postfix completion
fun main() {
    val ok = true
    ok.if<caret> // accept postfix -> if (ok) { ... }
}

[EX-KT-24] Kotlin: refactoring rename smoke (manual step)
fun renamedTarget() {}

fun main() {
    ren<caret> // after renaming renamedTarget, completion should suggest new name
}

[EX-KT-25] Kotlin: change signature smoke (manual step)
fun foo(a: Int) {}

fun main() {
    foo(<caret>) // after Change Signature add param b: String, completion in args should reflect new parameter structure
}

[EX-JV-1] Java: базовый комбайн
package demo;

import java.util.ArrayList;
import java.util.List;

public class Main {
  static class User {
    String name;
    int age;
    String getName() { return name; }
  }

  static User buildUser(String name, int age) { return new User(); }
  static void consumeUser(User u) {}

  public static void main(String[] args) {
    String userName = "Ann";
    int userAge = 21;
    User user = new User();

    user.<caret>                 // member completion

    User u1 = <caret>            // smart completion: user / new User() / buildUser(...)
    User u2 = buildUser(<caret>); // args completion: String, then ',' -> int
    consumeUser(<caret>);        // smart completion by expected type

    List<String> list = new ArrayList<>();
    list.add(<caret>);           // argument completion
  }
}

[EX-JV-2] Java: annotations completion
package demo;

@Depr<caret>
class A {}

[EX-JV-3] Java: visibility negative
package demo;

class A { private void secret() {} }

public class B {
  public static void main(String[] args) {
    A a = new A();
    a.sec<caret> // secret should NOT be suggested if not accessible (depends on context; here it's same package but different class)
  }
}

[EX-JV-4] Java: top-level keyword (в новом файле)
cla<caret>

[EX-JV-5] Java: static members only
package demo;

public class StaticTest {
  public static void main(String[] args) {
    Math.<caret> // should show static members like abs, max, etc.
  }
}

[EX-JV-6] Java: smart completion for return
package demo;

public class R {
  static class User {}
  static User build() { return new User(); }

  static User f() {
    return <caret>; // smart completion should suggest build(), new User()
  }
}

[EX-JV-7] Java: expected type in args
package demo;

public class C {
  static class User {}
  static void consume(User u) {}

  public static void main(String[] args) {
    User u = new User();
    consume(<caret>); // should suggest u
  }
}

[EX-JV-8] Java: commit by '('
package demo;

public class Commit {
  static void doWork() {}
  public static void main(String[] args) {
    doW<caret> // accept by pressing '(' -> doWork()
  }
}

[EX-JV-9] Java: auto-import ArrayList
package demo;

public class Imports {
  public static void main(String[] args) {
    ArrayLis<caret> list; // accept completion -> import java.util.ArrayList;
  }
}

[EX-JV-10] Java: import conflict (нужно наличие обоих классов в classpath)
package demo;

public class Conflict {
  public static void main(String[] args) {
    Date<caret> d; // should show java.util.Date and java.sql.Date with qualifiers
  }
}

[EX-JV-11] Java: deprecated marker
package demo;

public class Dep {
  /** @deprecated */
  @Deprecated
  static void oldMethod() {}

  public static void main(String[] args) {
    old<caret> // completion should show deprecated style for oldMethod
  }
}

[EX-JV-12] Java: generics + annotation params
package demo;

import java.util.List;
import java.util.Map;

@interface Ann { String value(); boolean enabled(); }

public class G {
  @Ann(<caret>) // suggest value= , enabled=
  void f(List<<caret>> a, Map<String, <caret>> b) {}
}

[EX-JV-13] Java: Javadoc completion
package demo;

public class Doc {
  /**
   * <caret>
   */
  int f(int a) { return a; }
}

[EX-JV-14] Java: postfix (если доступно)
package demo;

public class Postfix {
  public static void main(String[] args) {
    boolean ok = true;
    ok.if<caret>
  }
}

[EX-JV-15] Java: live template (если есть, например sout)
package demo;

public class T {
  public static void main(String[] args) {
    sout<caret>
  }
}

[EX-JV-16] Java: refactoring rename smoke
package demo;

public class Ref {
  static void targetToRename() {}
  public static void main(String[] args) {
    tar<caret>
  }
}

[EX-JV-17] Java: perf-ish scope (сгенерировать много методов)
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
  // (при необходимости добавить до 50-200)

  public static void main(String[] args) {
    Many m = new Many();
    m.met<caret> // completion should appear quickly
  }
}

[EX-TS-1] TypeScript: базовый комбайн
export type User = { name: string; age: number };

export function buildUser(name: string, age: number): User {
  return { name, age };
}

export function consumeUser(u: User) {}

const userName = "Ann";
const userAge = 21;
const user: User = buildUser(userName, userAge);

user.<caret>                   // member completion + filtering
const u1: User = <caret>       // smart completion
const u2: User = buildUser(<caret>) // args completion + comma -> second arg
consumeUser(<caret>)           // expected type completion

[EX-TS-2] TypeScript: import path completion
import x from "./<caret>"

[EX-TS-3] TypeScript: keyword
cla<caret>

[EX-TS-4] TypeScript: smart return
type User = { name: string; age: number };
function buildUser(): User { return { name: "Ann", age: 21 }; }

function f(): User {
  return <caret>
}

[EX-TS-5] TypeScript: expected type in args
type User = { name: string; age: number };
function consume(u: User) {}
const u: User = { name: "Ann", age: 21 };
consume(<caret>)

[EX-TS-6] TypeScript: unknown type fallback
declare const dyn: any;
dyn.<caret> // smart completion should not crash; results may be limited

[EX-TS-7] TypeScript: commit by '('
function doWork() {}
doW<caret> // accept by '(' -> doWork()

[EX-TS-8] TypeScript: auto import (needs multi-file project)
// In file a.ts export function utilFn() {}
// Here:
utilF<caret> // accept should add: import { utilFn } from "./a";

[EX-TS-9] TypeScript: generics
type Box<T> = { v: T };
let b: Box<<caret>>

[EX-PY-1] Python: базовый комбайн
from dataclasses import dataclass
from pathlib import Path

@dataclass
class User:
    name: str
    age: int

def build_user(name: str, age: int) -> User:
    return User(name, age)

def consume_user(u: User) -> None:
    pass

user_name = "Ann"
user_age = 21
user = User(user_name, user_age)

user.<caret>                # member completion
u1: User = <caret>          # smart-ish completion (if supported)
consume_user(<caret>)       # expected type completion (if supported)
Path("<caret>")             # path completion

[EX-PY-2] Python: decorator completion
from dataclasses import dataclass

@datac<caret>
class A:
    pass

[EX-PY-3] Python: keyword
cla<caret>

[EX-PY-4] Python: unknown/dynamic fallback
dyn = object()
dyn.<caret>  # should not crash, may show limited members

[EX-PY-5] Python: named args (if supported)
def f(name: str, age: int) -> None:
    pass

f(<caret>)  # should suggest name= and age= (if IDE supports)

[EX-PY-6] Python: docstring completion (if supported)
def f(a: int) -> int:
    """
    <caret>
    """
    return a

[EX-CS-1] C#: path completion
using System.IO;

class P {
  static void Main() {
    File.Open("<caret>", FileMode.Open);
  }
}

[EX-CONFIG-1] package.json (schema-based)
{
  "<caret>": ""
}

[EX-CONFIG-2] tsconfig.json (schema-based)
{
  "compilerOptions": {
    "<caret>": ""
  }
}

[EX-CONFIG-3] build.gradle.kts (Gradle Kotlin DSL)
plugins {
    <caret>
}

dependencies {
    <caret>
}

[EX-INJECT-1] SQL injection (требует включенной инъекции/поддержки)
fun main() {
    val q = "SELECT <caret> FROM users WHERE id = 1"
}

[EX-INJECT-2] Regex (если поддерживается)
fun main() {
    val r = "\\d<caret>"
}