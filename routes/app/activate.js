const uuid = require('node-uuid')

exports.route = {
    async get() {
    
      let userInfo = this.userInfo()
      let lockId = uuid.v4()
    
      console.log(userInfo._id)
      let lockUserCollection = await this.db('lock_user')
      await lockUserCollection.insertOne({ userId:userInfo._id, lockId , own:true})
      return { userId:userInfo._id, lockId }

    }
  }