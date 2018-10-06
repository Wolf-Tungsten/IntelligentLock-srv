/**
 * 
 * 数据持久化中间件
 * 
 * 使用 MongoDB 作为数据库
 * 该中间件维持到 MongoDB 的单例连接并暴露 ctx.db
 *
 */

const MongoClient = require('mongodb').MongoClient;
const config = require('../secret.json').mongodb

const user = encodeURIComponent(config.user);
const password = encodeURIComponent(config.pwd);
const authMechanism = 'DEFAULT';

// Connection URL
const url = `mongodb://${user}:${password}@${config.host}:${config.port}/intelligent-lock?authMechanism=${authMechanism}`
let mongodb = null

const getCollection = async(col) => {
  if (mongodb) {
    return mongodb.collection(col)
  } else {
    mongodb = await MongoClient.connect(url, { useNewUrlParser: true })
    mongodb = mongodb.db("intelligent-lock")
    return mongodb.collection(col)
  }
}

module.exports = async (ctx, next) => {

   ctx.db = getCollection
   await next()

}