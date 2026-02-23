// Shared model types used across completion test files.

export type User = { name: string; age: number };

export function buildUser(name: string, age: number): User {
  return { name, age };
}

export function consumeUser(u: User): void {}
