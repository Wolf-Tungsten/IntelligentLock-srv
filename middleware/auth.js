
/**
 * 
 * 身份认证中间件
 * 
 * 当请求头中包含 token 时，根据相关策略解析用户身份并向路由暴露用户信息
 *
 */

module.exports = async(ctx, next) => {

    const authCollection = await ctx.db('lock-auth')

    let _userInfo = null
    ctx.userInfo = () => {
        if (_userInfo) {
            _userInfo._id = _userInfo._id.toString()
            return _userInfo
        } else {
            throw 401
        }
    }

    if (ctx.request.headers.token) {
        _userInfo = await authCollection.findOne({token:ctx.request.headers.token})
    }

    await next()

}
