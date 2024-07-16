package safefun

import "fmt"

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
 @Time    : 2024/7/13 -- 14:07
 @Author  : bishop ❤️ MONEY
 @Description: 回滚序列函数
 @TODO: panic 封装处理
*/

// RollbackOp 回滚操作序列
type RollbackOp struct {
	opList []func() error
}

func NewRollbackOp() *RollbackOp {
	return &RollbackOp{}
}

func (r *RollbackOp) Add(f func() error) {
	r.opList = append(r.opList, f)
}

func (r *RollbackOp) Rollback() (err error) {
	// 倒叙执行
	for i := len(r.opList) - 1; i >= 0; i-- {
		nErr := r.opList[i]()
		if nErr != nil {
			err = fmt.Errorf("%s: %s", err, nErr)
		}
	}
	return
}
