import jwt from 'jsonwebtoken';
import User, { IUser } from '../models/user.model';

export class AuthService {
  static generateToken(user: IUser): string {
    return jwt.sign({ id: user._id }, process.env.JWT_SECRET as string, {
      expiresIn: process.env.JWT_EXPIRES_IN
    });
  }
static async getMe(userId: string) {
    // Find user by ID and return user data
    const user = await User.findById(userId).select('-password');
    if (!user) {
        throw new Error('User not found');
    }
    return user;
}
  static async register(email: string, password: string, name: string): Promise<IUser> {
    const userExists = await User.findOne({ email });
    if (userExists) {
      throw new Error('User already exists');
    }

    const user = await User.create({
      email,
      password,
      name
    });

    return user;
  }

  static async login(email: string, password: string): Promise<{ user: IUser; token: string }> {
    const user = await User.findOne({ email });
    if (!user) {
      throw new Error('Invalid credentials');
    }

    const isMatch = await user.comparePassword(password);
    if (!isMatch) {
      throw new Error('Invalid credentials');
    }

    const token = this.generateToken(user);
    return { user, token };
  }
}
