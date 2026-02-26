package completion;

/**
 * Postfix completion, live templates, and rename refactoring smoke test.
 * Covers: TC-114 (postfix .if), TC-116 (live template sout), TC-88 (rename).
 * Based on EX-JV-14, EX-JV-15, EX-JV-16.
 */
public class TemplatesRefactor {

    static void targetToRename() {
    }

    public static void main(String[] args) {
        boolean ok = true;

        // <caret> TC-114: Delete if-block below, type 'ok.if', accept postfix completion;
        //   expect if (ok) { ... }
        if (ok) {
            System.out.println("postfix");
        }

        // <caret> TC-116: Delete call below, type 'sout', accept template;
        //   expect System.out.println(...)
        System.out.println("template");

        // <caret> TC-88: Delete 'targetToRename()' below so only 'tar' remains,
        //   then invoke completion — expect targetToRename()
        targetToRename();
    }
}
