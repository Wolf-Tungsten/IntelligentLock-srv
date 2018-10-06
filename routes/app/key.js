exports.route = {
    async get({lockId}) {

        let userInfo = this.userInfo()
        let userId = userInfo._id

        // 检查用户关系，判断是否可以解锁
        let lockUserCollection = await this.db('lock_user')
        if((await lockUserCollection.countDocuments({userId, lockId})) === 0) {
            throw '无权解锁'
        }

        let lockCollection = await this.db('lock')
        let lockInfo = await lockCollection.findOne({lockId})
        
        if (lockInfo && lockInfo.enable) {
            
            let lockKeyCollection = await this.db('lock_key')
            let key = await lockKeyCollection.findOne({lockId})
            if ( key ) {
                await lockKeyCollection.deleteOne({lockId, key:key.key})
                return key.key
            }
             
            // 执行到此处说明该锁还没有提交过密钥
            throw '未生成密钥'
            
        }
        // 执行到此处说明该锁目前处于停用状态
        throw '该锁停止使用'

    }

  }