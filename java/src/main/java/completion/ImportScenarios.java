package completion;

// NOTE: No imports here — that is intentional for auto-import testing.

/**
 * Auto-import and import-conflict completion scenarios.
 * Covers: TC-54 (auto-import ArrayList), TC-57 (import conflict Date), TC-59 (settings).
 * Based on EX-JV-9, EX-JV-10.
 */
public class ImportScenarios {

    public static void main(String[] args) {
        // TC-54: type "ArrayLis" and accept — should auto-import java.util.ArrayList
        ArrayLis // <caret> TC-54: auto-import ArrayList
        ;

        // TC-57: type "Date" — should show java.util.Date and java.sql.Date with qualifiers
        Date // <caret> TC-57: import conflict — java.util.Date vs java.sql.Date
        ;

        // TC-59: verify IDE settings for auto-import behavior
        // <caret> TC-59: check "Add unambiguous imports on the fly" setting
    }
}
