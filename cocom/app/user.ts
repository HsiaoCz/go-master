import { Request, Response } from "express";
import user from "../types/user";
const handleUser = (req: Request, resp: Response) => {
  user.toString();
  resp.send("hello my man");
};

export default handleUser;
