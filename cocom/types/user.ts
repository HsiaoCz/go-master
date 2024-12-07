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
}

export default Users;
