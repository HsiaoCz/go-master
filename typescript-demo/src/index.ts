// Basic Types
let userName: string = "John";
let age: number = 30;
let isStudent: boolean = true;
let numbers: number[] = [1, 2, 3, 4, 5];
// tell me what is tuple?
let tuple: [string, number] = ["hello", 42];

// Type Inference
let inferredString = "This type is inferred";  // TypeScript knows this is a string

// Interface
interface Person {
    name: string;
    age: number;
    email?: string;  // Optional property
}

// Class
class Student implements Person {
    constructor(
        public name: string,
        public age: number,
        private studentId: string
    ) {}

    // Method
    study(): void {
        console.log(`${this.name} is studying`);
    }
}

// Function with Type Annotations
function greet(person: Person): string {
    return `Hello, ${person.name}! You are ${person.age} years old.`;
}

// Enum
enum Role {
    ADMIN = "ADMIN",
    USER = "USER",
    GUEST = "GUEST"
}

// Union Types
type Status = "active" | "inactive" | "pending";

// Generic Function
function firstElement<T>(arr: T[]): T | undefined {
    return arr[0];
}

// Usage Examples
const student = new Student("Alice", 20, "12345");
student.study();

const greeting = greet({ name: "Bob", age: 25 });
console.log(greeting);

const numbers2 = [1, 2, 3, 4, 5];
const firstNumber = firstElement(numbers2);
console.log("First number:", firstNumber);

const userRole: Role = Role.ADMIN;
console.log("User role:", userRole);
