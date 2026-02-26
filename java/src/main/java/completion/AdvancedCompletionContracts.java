package completion;

/**
 * Advanced completion scenarios for override/overload and broken-code resilience.
 *
 * Covers: TC-105, TC-108, TC-56.
 */
public class AdvancedCompletionContracts {

    interface SyncWorker {
        void runSync();
        default String workerId() {
            return "sync-worker";
        }
    }

    static class LocalWorker implements SyncWorker {
        // <caret> TC-105: Type 'over' in class body and invoke completion;
        //   expect override templates for runSync()/workerId().
        @Override
        public void runSync() {
        }
    }

    static String pick(String value) {
        return "S:" + value;
    }

    static String pick(int value) {
        return "I:" + value;
    }

    static void overloadResolution() {
        String str = "Ann";
        int num = 7;
        // <caret> TC-108-A: Delete argument below in pick(...), invoke completion;
        //   with expected String overload, prefer String-compatible values.
        String tc113a = pick(str);
        // <caret> TC-108-B: Delete argument below in pick(...), invoke completion;
        //   with expected int overload, prefer numeric values.
        String tc113b = pick(num);
        System.out.println(tc113a + tc113b);
    }

    static void brokenCodeResilience() {
        String user = "Ann";
        // <caret> TC-56: Temporarily remove closing ')' in println below and invoke completion
        //   after 'user.'; IDE should stay stable and keep relevant suggestions.
        System.out.println(user.toUpperCase());
    }

    public static void main(String[] args) {
        overloadResolution();
        brokenCodeResilience();
    }
}
