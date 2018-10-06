exports.route = {

    // 获取所有管理的/可以解锁的终端
    async get() {

        let locks = []
        let userInfo = this.userInfo()
        let userId = userInfo._id
        let lockUserCollection = await this.db('lock_user')
        let ownList = await lockUserCollection.find({userId, own:true, activate:true}).toArray()
        ownList.forEach( k => {
            locks.push({
                lockId:k.lockId,
                admin:'own'
            })
        })
        let accessList = await lockUserCollection.find({userId, access:true, activate:true}).toArray()
        accessList.forEach( k => {
            locks.push({
                lockId:k.lockId,
                admin:'access'
            })
        })
        return { locks }

    },

    //为一个终端添加可以访问的用户
    async post({username, lockId}) {
        let userId = this.userInfo()._id
        
        // 检查用户是否是该终端的拥有者
        let lockUserCollection = await this.db('lock_user')
        if ((await lockUserCollection.countDocuments({userId, lockId, own:true})) <= 0) {
            throw '无权操作'
        }

        const authCollection = await this.db('lock-auth')
        let accessUser = await authCollection.findOne({ username })

        await lockUserCollection.insertOne({userId:accessUser._id.toString(), lockId, own:false, access:true, activate:true})
        return 'ok'
    },

    async put({lockId, command}) {
        let userId = this.userInfo()._id
        
        // 检查用户是否是该终端的拥有者
        let lockUserCollection = await this.db('lock_user')
        if ((await lockUserCollection.countDocuments({userId, lockId, own:true})) <= 0) {
            throw '无权操作'
        }
        let lockCollection = await this.db('lock')
        switch(command) {
            case 'disable':
                await lockCollection.updateOne({lockId}, {$set:{enable:false}})
                return 'ok'
            case 'enable':
                await lockCollection.updateOne({lockId}, {$set:{enable:true}})
                return 'ok'
            default:
                throw '不知道你要干什么'
        }

    },

    async delete({lockId, username}) {
        let userId = this.userInfo()._id
        
        // 检查用户是否是该终端的拥有者
        let lockUserCollection = await this.db('lock_user')
        if ((await lockUserCollection.countDocuments({userId, lockId, own:true})) <= 0) {
            throw '无权操作'
        }

        const authCollection = await this.db('lock-auth')
        let accessUser = await authCollection.findOne({ username })

        await lockUserCollection.deleteMany({userId:accessUser._id.toString(), lockId})

        return 'ok'
    }
}