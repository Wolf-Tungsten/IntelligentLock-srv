const koa = require('koa')
const app = new koa()
const kf = require('kf-router')
const fs = require('fs')

// 调试时打开，设置所有 Promise 超时 20 秒，超过 20 秒自动 reject 并输出超时 Promise 所在位置
// require('./promise-timeout')(20000)

// 将 moment 导出到全局作用域
global.moment = require('moment')

// 解析 YAML 配置文件
const config = require('js-yaml').load(fs.readFileSync('./config.yml'))
exports.config = config

// 为 Moment 设置默认语言
moment.locale('zh-cn')

// 出错输出
process.on('unhandledRejection', e => { throw e })
process.on('uncaughtException', console.trace)

// 监听两个结束进程事件，将它们绑定至 exit 事件，有两个作用：
// 1. 使用 child_process 运行子进程时，可直接监听主进程 exit 事件来杀掉子进程；
// 2. 防止按 Ctrl+C 时程序变为后台僵尸进程。
process.on('SIGTERM', () => process.exit())
process.on('SIGINT', () => process.exit())

/**
  # WS3 框架中间件
  以下中间件共分为六层，每层内部、层与层之间都严格按照依赖关系排序。

  ## A. 超监控层
  不受监控、不受格式变换的一些高级操作
*/
// 1. 跨域中间件，定义允许访问本服务的第三方前端页面
app.use(require('./middleware/cors'))

/**
  ## B. 监控层
  负责对服务运行状况进行监控，便于后台分析和交互，对服务本身不存在影响的中间件。
*/
// 1. 如果是生产环境，显示请求计数器；此中间件在 module load 时，会对 console 的方法做修改
app.use(require('./middleware/counter'))
// 2. 日志输出，需要依赖返回格式中间件中返回出来的 JSON 格式
app.use(require('./middleware/logger'))

/**
  ## C. 接口层
  为了方便双方通信，负责对服务收到的请求和发出的返回值做变换的中间件。
*/
// 2. 参数格式化，对上游传入的 URL 参数和请求体参数进行合并
app.use(require('./middleware/params'))
// 3. 返回格式化，将下游返回内容包装一层JSON
app.use(require('./middleware/return'))

/**
  ## D. API 层
  负责为路由处理程序提供 API 以便路由处理程序使用的中间件。
*/
// 1. 接口之间相互介绍的 API
app.use(require('./middleware/related'))
// 2. 网络请求，为身份认证和路由处理程序提供了网络请求 API
app.use(require('./middleware/axios'))

/**
  ## E. 数据持久化层
  维持数据库连接并暴露接口
*/
// 1. MongoDB数据库接口
app.use(require('./middleware/presistence'))
app.use(require('./middleware/auth'))

/**
  ## F. 路由层
  负责调用路由处理程序执行处理的中间件。
*/
app.use(kf())
app.listen(config.port)

// 开发环境下，启动 REPL
if (process.env.NODE_ENV === 'development') {
  require('./repl').start()
}
