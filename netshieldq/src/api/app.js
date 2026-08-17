const express = require('express');
const morgan = require('morgan');

const app = express();
const port = 3000;
app.use(morgan('dev'));

// 路由定义
app.get('/', (req, res) => {
    res.send('欢迎访问首页！');
});

app.get('/about', (req, res) => {
    res.send('关于我们页面');
});

// 监听端口
app.listen(port, () => {
    console.log(`服务器运行在 http://localhost:${port}`);
});
