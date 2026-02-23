package completion;

/**
 * Commit completion by pressing '(' character.
 * Covers: TC-50.
 * Based on EX-JV-8.
 */
public class AcceptCommit {

    static void doWork() {
    }

    public static void main(String[] args) {
        doW // <caret> TC-50: accept by pressing '(' — expect doWork()
        ;
    }
}
