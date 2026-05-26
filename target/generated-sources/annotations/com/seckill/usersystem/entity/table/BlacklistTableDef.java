package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class BlacklistTableDef extends TableDef {

    public static final BlacklistTableDef BLACKLIST = new BlacklistTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn REASON = new QueryColumn(this, "reason");

    /**
     * 状态：0-已解封，1-封禁中
     */
    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn DEVICE_ID = new QueryColumn(this, "device_id");

    public final QueryColumn IP_ADDRESS = new QueryColumn(this, "ip_address");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn EXPIRE_TIME = new QueryColumn(this, "expire_time");

    public final QueryColumn OPERATOR_ID = new QueryColumn(this, "operator_id");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    /**
     * 黑名单类型：1-用户，2-IP，3-设备
     */
    public final QueryColumn BLACKLIST_TYPE = new QueryColumn(this, "blacklist_type");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, REASON, STATUS, USER_ID, DEVICE_ID, IP_ADDRESS, CREATE_TIME, EXPIRE_TIME, OPERATOR_ID, UPDATE_TIME, BLACKLIST_TYPE};

    public BlacklistTableDef() {
        super("", "sys_blacklist");
    }

    private BlacklistTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public BlacklistTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new BlacklistTableDef("", "sys_blacklist", alias));
    }

}
