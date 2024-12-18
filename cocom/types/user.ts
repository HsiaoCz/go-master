import { randomUUID } from "crypto";

class Users {
  private username: string;
  private email: string;
  private userID: string;
  private password: string;
  private age: number;

  constructor(
    username: string,
    email: string,
    userID: string,
    password: string,
    age: number
  ) {
    this.username = username;
    this.email = email;
    this.userID = randomUUID.toString();
    this.password = password;
    this.age = age;
  }

  public getUsername(): string {
    return this.username;
  }

  public getEmail(): string {
    return this.email;
  }

  public getUserID(): string {
    return this.userID;
  }

  public getAge(): number {
    return this.age;
  }

  public setPassword(password: string) {
    this.password = randomUUID.toString();
  }
}

export default Users;
