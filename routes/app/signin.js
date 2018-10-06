const uuid = require('node-uuid')

exports.route = {
    get() {
      throw 'POST ONLY'
    },

    async post({username, password}) {

        // 校验用户名密码是否有效
        let authCollection = await this.db('lock-auth')
        if ( (await authCollection.countDocuments({username, password})) === 0 ) {
            throw 401
        }

        // 生成token
        let token = uuid.v4()
        // 生成时间戳
        let timestamp = +moment()

        await authCollection.updateOne({username, password}, {$set:{token, timestamp}})

        return token

    }

  }