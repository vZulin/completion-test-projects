package completion;

/**
 * Visibility negative test — private members must not leak.
 * Covers: TC-18.
 * Based on EX-JV-3.
 */
public class Visibility {

    public static void main(String[] args) {
        Secret a = new Secret();
        // <caret> TC-18: Delete 'toString()' below so only 'a.sec' remains after typing;
        //   secret() should NOT be suggested (private access)
        a.toString();
    }
}

class Secret {
    private void secret() {
    }
}
