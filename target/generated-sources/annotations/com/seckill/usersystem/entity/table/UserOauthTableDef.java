package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class UserOauthTableDef extends TableDef {

    public static final UserOauthTableDef USER_OAUTH = new UserOauthTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn AVATAR = new QueryColumn(this, "avatar");

    public final QueryColumn OPENID = new QueryColumn(this, "openid");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn UNIONID = new QueryColumn(this, "unionid");

    public final QueryColumn NICKNAME = new QueryColumn(this, "nickname");

    public final QueryColumn EXPIRES_IN = new QueryColumn(this, "expires_in");

    /**
     * OAuth类型：WECHAT/ALIPAY/QQ/GOOGLE/GITHUB
     */
    public final QueryColumn OAUTH_TYPE = new QueryColumn(this, "oauth_type");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn ACCESS_TOKEN = new QueryColumn(this, "access_token");

    public final QueryColumn REFRESH_TOKEN = new QueryColumn(this, "refresh_token");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, AVATAR, OPENID, USER_ID, UNIONID, NICKNAME, EXPIRES_IN, OAUTH_TYPE, CREATE_TIME, UPDATE_TIME, ACCESS_TOKEN, REFRESH_TOKEN};

    public UserOauthTableDef() {
        super("", "sys_user_oauth");
    }

    private UserOauthTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public UserOauthTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new UserOauthTableDef("", "sys_user_oauth", alias));
    }

}
