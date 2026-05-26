package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class UserProfileTableDef extends TableDef {

    public static final UserProfileTableDef USER_PROFILE = new UserProfileTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn BIO = new QueryColumn(this, "bio");

    public final QueryColumn AVATAR = new QueryColumn(this, "avatar");

    public final QueryColumn GENDER = new QueryColumn(this, "gender");

    public final QueryColumn ID_CARD = new QueryColumn(this, "id_card");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn ADDRESS = new QueryColumn(this, "address");

    public final QueryColumn BIRTHDAY = new QueryColumn(this, "birthday");

    public final QueryColumn NICKNAME = new QueryColumn(this, "nickname");

    public final QueryColumn REAL_NAME = new QueryColumn(this, "real_name");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, BIO, AVATAR, GENDER, ID_CARD, USER_ID, ADDRESS, BIRTHDAY, NICKNAME, REAL_NAME, CREATE_TIME, UPDATE_TIME};

    public UserProfileTableDef() {
        super("", "sys_user_profile");
    }

    private UserProfileTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public UserProfileTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new UserProfileTableDef("", "sys_user_profile", alias));
    }

}
