package completion;

/**
 * Static members completion — only static members should appear.
 * Covers: TC-26, TC-27.
 * Based on EX-JV-5.
 */
public class StaticMembers {

    static class StaticBox {
        static int version() {
            return 1;
        }

        int instanceValue() {
            return 42;
        }
    }

    public static void main(String[] args) {
        // <caret> TC-26: Delete 'version()' below so only 'StaticBox.' remains,
        //   invoke completion — expect static members via class reference.
        // <caret> TC-27: In the same popup for 'StaticBox.', verify instance members
        //   (like instanceValue()) are not suggested in class-name context.
        int tc31 = StaticBox.version();
        System.out.println(tc31);
    }
}
