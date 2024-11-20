const express = require("express");
const multer = require("multer");

// 设置存储位置和文件名
const storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, "uploads/"); // 存储路径
  },
  filename: function (req, file, cb) {
    cb(
      null,
      file.fieldname + "-" + Date.now() + path.extname(file.originalname)
    ); // 文件名
  },
});

// 初始化multer
const upload = multer({ storage: storage });

const app = express();

// 创建上传路由
app.post("/upload", upload.single("image"), (req, res, next) => {
  try {
    console.log(req.file); // 打印上传文件的信息
    res.send("文件上传成功！");
  } catch (err) {
    console.error(err);
    res.status(500).send("文件上传失败！");
  }
});

// 启动服务器
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
