package zstorage

/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/13 -- 15:47
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: zstorage.go
*/

// Errors of not found include redis/mongo/mysql
const (
	RedisNil           = "redis: nil"
	MongoNotFound      = "not found"
	MysqlNoRows        = "sql: no rows in result set"
	RPCRespNil         = "rpc: nil"
	XsqlErrNilRows     = "[scanner]: rows can't be nil"
	XsqlErrEmptyResult = `[scanner]: empty result`
)

// ContainStr check if string target in string array
func ContainStr(source []string, target string) bool {
	for _, elem := range source {
		if elem == target {
			return true
		}
	}
	return false
}

// ContainInt check if int target in int array
func ContainInt(source []int, target int) bool {
	for _, elem := range source {
		if elem == target {
			return true
		}
	}
	return false
}

// IsNotFound check if nil in redis/not found in mongo/no rows in mysql
// If sdk has change error msg, this will return wrong result, rarely
func IsNotFound(err error) bool {
	if err != nil {
		errMsg := err.Error()
		if errMsg == RedisNil ||
			errMsg == MongoNotFound ||
			errMsg == MysqlNoRows ||
			errMsg == XsqlErrNilRows ||
			errMsg == XsqlErrEmptyResult ||
			errMsg == RPCRespNil {
			return true
		}
		return false
	}
	return false
}
