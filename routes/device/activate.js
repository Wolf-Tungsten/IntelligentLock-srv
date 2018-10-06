exports.route = {

    async get({ lockId }) {
    
      let lockCollection = await this.db('lock')
      if ( (await lockCollection.countDocuments({lockId})) === 1) {
          return {activate:true}
      } else {
          return {activate:false}
      }

    },

    async post({lockId, userId}) {

        // 先通过 lock_user 集合检查是否存在对应关系，存在则证明是正确的激活过程
        let lockUserCollection = await this.db('lock_user')
        if ( (await lockUserCollection.countDocuments({lockId, userId})) <= 0 ) {
            throw '非法激活过程-用户ID不正确'
        }

        // 检查 lockId 是否已被激活
        let lockCollection = await this.db('lock')
        if ( (await lockCollection.countDocuments({lockId})) > 0) {
            throw '非法激活过程-该lockId已激活'
        }

        // 激活
        await lockCollection.insertOne({lockId, enable:true})
        await lockUserCollection.updateOne({lockId, userId}, {$set:{activate:true}})

        return 'ok'
       
    }

}