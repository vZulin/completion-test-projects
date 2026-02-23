package completion;

/**
 * Visibility negative test — private members must not leak.
 * Covers: TC-19.
 * Based on EX-JV-3.
 */
public class Visibility {

    public static void main(String[] args) {
        Secret a = new Secret();
        a.sec // <caret> TC-19: secret() should NOT be suggested (private access)
        ;
    }
}

class Secret {
    private void secret() {
    }
}
