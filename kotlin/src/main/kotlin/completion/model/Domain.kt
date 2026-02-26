/**
 * Shared domain model used across all completion test scenarios.
 *
 * Provides [User] data class, builder, and consumer functions
 * referenced by TC-1 through TC-91.
 */
package completion.model

data class User(val name: String, val age: Int)

fun buildUser(name: String, age: Int): User = User(name, age)

fun consumeUser(u: User) {}
