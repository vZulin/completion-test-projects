package completion;

/**
 * Bracket-index completion scenarios for Java arrays.
 *
 * Covers: TC-54.
 */
public class BracketIndexCompletion {

    public static void main(String[] args) {
        int[] ages = {21, 22, 23};
        int position = 1;
        // <caret> TC-54-A: Delete 'position' below inside ages[...], invoke completion;
        //   expect Int variables from scope (e.g., position)
        int tc104a = ages[position];

        int[][] matrix = {{1, 2}, {3, 4}};
        int row = 0;
        int col = 1;
        // <caret> TC-54-B: Delete 'row' in first [] and 'col' in second [] below;
        //   invoke completion in each bracket, expect Int variables row/col
        int tc104b = matrix[row][col];

        System.out.println(tc104a + tc104b);
    }
}
