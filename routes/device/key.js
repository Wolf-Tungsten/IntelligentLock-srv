exports.route = {

    async get() {

        throw 'POST ONLY'

    },

    async post({lockId, keys}) {

        let lockKeyCollection = await this.db('lock_key')
        keys = keys.map( k => {return {lockId, key:k}} )
        await lockKeyCollection.insertMany(keys)
        return 'ok'

    }

}