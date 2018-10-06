exports.route = {
    get() {
      throw 'POST ONLY'
    },

    async post({username, password}) {

        if (!/^[0-9A-Za-z_]{8,}$/.test(username)) {
            throw '用户名不合法'
        }

        if (!/^.{6,}$/.test(password)) {
            throw '密码不合法'
        }

        // 用户名需要全局唯一，首先检查
        let authCollection = await this.db('lock-auth')
        if ( (await authCollection.countDocuments({username})) > 0 ) {
            throw '用户名已被占用'
        }

        await authCollection.insertOne({ username, password })

        return 'success'

    }

  }