package completion;

/**
 * Postfix completion, live templates, and rename refactoring smoke test.
 * Covers: TC-81 (postfix .if), TC-83 (live template sout), TC-87 (rename).
 * Based on EX-JV-14, EX-JV-15, EX-JV-16.
 */
public class TemplatesRefactor {

    static void targetToRename() {
    }

    public static void main(String[] args) {
        boolean ok = true;

        // TC-81: postfix completion — type "ok.if" and accept to wrap in if-statement
        ok.if // <caret> TC-81: postfix .if — expect if (ok) { }
        ;

        // TC-83: live template — type "sout" and accept to expand to System.out.println()
        sout // <caret> TC-83: live template sout — expect System.out.println()
        ;

        // TC-87: rename refactoring smoke — type "tar" and see targetToRename in completions
        tar // <caret> TC-87: rename smoke — expect targetToRename()
        ;
    }
}
