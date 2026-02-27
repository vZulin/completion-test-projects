package completion;

// NOTE: No imports here — that is intentional for auto-import testing.

/**
 * Auto-import and import-conflict completion scenarios.
 * Covers: TC-43 (auto-import ArrayList), TC-46 (import conflict Date), TC-68, TC-69 (settings).
 * Based on EX-JV-9, EX-JV-10.
 */
public class ImportScenarios {

    public static void main(String[] args) {
        // <caret> TC-43: Delete both 'java.util.' qualifiers below and leave 'ArrayList';
        //   then type 'ArrayLis' and accept completion — should add import java.util.ArrayList
        java.util.ArrayList<String> listForImport = new java.util.ArrayList<>();

        // <caret> TC-46: Delete package qualifiers below and type 'Date';
        //   completion should show java.util.Date and java.sql.Date with qualifiers
        // <caret> TC-68: Choose java.util.Date and verify java.sql.Date is NOT auto-imported.
        java.util.Date utilDate = new java.util.Date();
        java.sql.Date sqlDate = new java.sql.Date(System.currentTimeMillis());

        // <caret> TC-69: Disable auto-import, then delete 'java.util.' and final 't' below so token is 'ArrayLis';
        //   invoke completion and accept ArrayList — import should stay absent or FQN should be inserted.
        java.util.ArrayList<String> list = new java.util.ArrayList<>();

        System.out.println(listForImport.size() + utilDate.getTime() + sqlDate.getTime() + list.size());
    }
}
