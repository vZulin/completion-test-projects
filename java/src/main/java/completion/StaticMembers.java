package completion;

/**
 * Static members completion — only static members should appear.
 * Covers: TC-27.
 * Based on EX-JV-5.
 */
public class StaticMembers {

    public static void main(String[] args) {
        // <caret> TC-27: Delete 'abs(-42)' below so only 'Math.' remains,
        //   then invoke completion — expect static members only
        double tc31 = Math.abs(-42);
        System.out.println(tc31);
    }
}
