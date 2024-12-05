// app.ts
import express, { Request, Response } from "express";
const app = express();
const PORT = 3000;

// 定义一个简单的路由
app.get("/", (req: Request, res: Response) => {
  res.send("Hello, World!");
});

// 启动服务器
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
