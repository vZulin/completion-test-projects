package completion;

/**
 * Commit completion by pressing '(' character.
 * Covers: TC-40.
 * Based on EX-JV-8.
 */
public class AcceptCommit {

    static void doWork() {
    }

    public static void main(String[] args) {
        // <caret> TC-40: Delete 'ork()' below so only 'doW' remains,
        //   then press '(' to accept completion; expect doWork()
        doWork();
    }
}
