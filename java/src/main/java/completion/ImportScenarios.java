package completion;


import static java.lang.Math.max;

// NOTE: No regular class imports here; qualified names are intentional for auto-import testing.
// <caret> TC-123: In the static import above, delete 'h.max;' so it ends with 'java.lang.Mat';
//   invoke completion and accept Math. It should insert exactly 'java.lang.Math.'.

/**
 * Auto-import and import-conflict completion scenarios.
 * Covers: TC-43 (auto-import ArrayList), TC-46 (import conflict Date), TC-68, TC-69 (settings),
 * TC-122 (constructor diamond completion), TC-123 (static import class completion),
 * TC-128 (Remote Dev frontend insertion import sync), TC-129 (FQN constructor completion),
 * TC-130 (var with short class constructor completion), TC-131 (var with FQN constructor completion),
 * TC-134 (known-type assignment constructor completion), TC-135 (raw expected type constructor completion).
 * Based on EX-JV-9, EX-JV-10.
 */
public class ImportScenarios {

    public static void main(String[] args) {
        // <caret> TC-43: Delete both 'java.util.' qualifiers below and leave 'ArrayLi';
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



    void varShortClass() {
        var names = new ArrayList<>();
        //System.out.println(names.size());

    }

    void varFullyQualified() {
        var names = new java.util.ArrayList<>();
        //System.out.println(names.size());
    }

    void constructorCompletion() {
        java.util.List<String> names = new ArrayList<>();
    }


    void constructorDiamondCompletion() {
        // <caret> TC-122: Remove 'java.util.' from the constructor call below and delete 't<>();'
        //   so the expression ends with 'new ArrayLis'. Accept ArrayList and verify 'new ArrayList<>()'
        //   plus import java.util.ArrayList.
        java.util.List<String> names = new java.util.ArrayList<>();

        System.out.println(names.size());
    }

    void remoteFrontendImportCompletion() {
        // <caret> TC-128: In Remote Dev with frontend completion enabled, delete 'java.util.'
        //   and the final 't' below so the declaration starts with 'ArrayLis list;'.
        //   Accept ArrayList and verify frontend/backend documents converge with the import added.
        java.util.ArrayList list;

        System.out.println(max(1, 0));
    }

    void fullyQualifiedKnownTargetConstructorCompletion() {
        // <caret> TC-129: Delete 't<>();' below so the expression ends with 'new java.util.ArrayLis'.
        //   Accept ArrayList and verify the qualifier is not duplicated and constructor insertion stays valid.
        java.util.List<String> names = new java.util.ArrayList<>();

        System.out.println(names.size());
    }

    void varShortClassConstructorCompletion() {
        // <caret> TC-130: In Remote Dev with frontend completion enabled, remove 'java.util.'
        //   from the constructor call below and delete 't<>();' so the expression ends with 'new ArrayLis'.
        //   Accept ArrayList and verify frontend insertion adds 'new ArrayList<>()' and backend sync adds import.
        var names = new java.util.ArrayList<>();

        System.out.println(names.size());
    }

    void varFullyQualifiedConstructorCompletion() {
        // <caret> TC-131: Delete 't<>();' below so the expression ends with 'new java.util.ArrayLis'.
        //   Accept ArrayList and verify the FQN is completed without duplicate 'java.util' or a new import.
        var names = new java.util.ArrayList<>();

        System.out.println(names.size());
    }

    void knownTypeAssignmentConstructorCompletion() {
        java.util.List<String> names;

        // <caret> TC-134: Remove 'java.util.' from the constructor call below and delete 't<>();'
        //   so the assignment ends with 'new ArrayLis'. Accept ArrayList and verify constructor insertion
        //   plus import java.util.ArrayList.
        names = new java.util.ArrayList<>();

        System.out.println(names.size());
    }

    void rawExpectedTypeConstructorCompletion() {
        // <caret> TC-135: Remove 'java.util.' from the constructor call below and delete 't();'
        //   so the expression ends with 'new ArrayLis'. Accept ArrayList and verify raw expected type
        //   does not produce an invalid diamond.
        java.util.List names = new java.util.ArrayList();

        System.out.println(names.size());
    }
}
