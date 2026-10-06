class Animal {
    age: number = 1;
    speak() { console.log("Animal", this.age); }
    grow() { this.age++; this.speak(); }
}
class Person extends Animal {
    age: number = 2;
    speak() { console.log("Person", this.age); super.speak(); }
}
class Employee extends Person {
    age: number = 3;
    speak() { console.log("Employee", this.age); super.speak(); }
}
class Manager extends Employee {
    age: number = 4;
    speak() { console.log("Manager", this.age); super.speak(); }
}
const manager = new Manager();
manager.grow();
manager["age"] = 10;
console.log(manager.age);
const animals: Animal[] = [manager, new Animal()];
animals[0].speak();
animals[1].speak();
