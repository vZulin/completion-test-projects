package completion.model;

/**
 * Shared model class used across all completion test scenarios.
 */
public class User {

    String name;
    int age;

    public User() {
    }

    public User(String name, int age) {
        this.name = name;
        this.age = age;

    }

    public String getName() {
        return name;
    }

    public int getAge() {
        return age;
    }

    public static User buildUser(String name, int age) {
        return new User(name, age);
    }

    public static void consumeUser(User u) {
        // no-op for completion testing
    }
}
