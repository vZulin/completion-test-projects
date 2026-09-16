type User = {
    name: string;
};

const greeting = (user: User): string => `Hello, ${user.name}!`;

console.log(greeting({ name: "TypeScript" }));
