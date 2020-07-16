package db

import (
	"fmt"
	"strings"
)

/*
 * 插入数据
 * tableName:表名 fields:键值对
 */
func (c *Client) Set(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (__key__,__version__,%s)
VALUES (%s);`

	columns, values := []string{}, []string{"$1", "$2"}
	args := []interface{}{key, 1}
	var i = 3
	for k, v := range fields {
		columns = append(columns, k)
		values = append(values, fmt.Sprintf("$%d", i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), strings.Join(values, ","))
	//fmt.Println(sqlStatement)
	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

/*
 * 更新数据
 * tableName:表名 key:key fields:键值对
 */
func (c *Client) Update(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
UPDATE "%s" 
SET %s
WHERE __key__ = '%s';`

	columns := []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		if k != "__key__" {
			columns = append(columns, fmt.Sprintf(`%s = $%d`, k, i))
			i++
			args = append(args, v)
		}
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), key)
	//fmt.Println(sqlStatement)
	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err

}

/*
 * set
 * 没有数据插入，有则更改。
 * tableName:表名 key:主键 fields:键值对，包含主键
 */
func (c *Client) SetNx(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (__key__,__version__,%s)
VALUES(%s) 
ON conflict(%s) DO 
UPDATE SET %s;`

	columns, values, sets := []string{}, []string{"$1", "$2"}, []string{}
	args := []interface{}{key, 1}
	var i = 3
	for k, v := range fields {
		if "__key__" != k {
			sets = append(sets, fmt.Sprintf(`%s = $%d`, k, i))
			columns = append(columns, k)
			values = append(values, fmt.Sprintf("$%d", i))
			i++
			args = append(args, v)
		}
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), strings.Join(values, ","), "__key__", strings.Join(sets, ","))
	//fmt.Println(sqlStatement, args)
	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

/*
 * 读取数据。
 * tableName:表名 key:主键
 * ret 返回键值对
 */
func (c *Client) Get(tableName, key string) (ret map[string]interface{}, err error) {
	sqlStr := `
SELECT * FROM "%s" 
WHERE __key__ = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, tableName, key)
	//fmt.Println(sqlStatement)
	rows, err := c.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	var columns []string
	ret = map[string]interface{}{}

	defer rows.Close()
	for rows.Next() {
		columns, err = rows.Columns()
		if err != nil {
			return nil, err
		}

		columnsLen := len(columns)
		values := make([]interface{}, 0, columnsLen)
		for i := 0; i < columnsLen; i++ {
			values = append(values, new(interface{}))
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		for i, k := range columns {
			ret[k] = *(values[i].(*interface{}))
		}
		break
	}

	return ret, nil
}

/*
 * 读取所有数据。
 * tableName:表名
 * ret 返回键值对的slice
 * limit 限制每次读取 1000 行数据.
 */
func (c *Client) GetAll(tableName string, callback func([]map[string]interface{}) error) error {
	start := 0
	total, err := c.Count(tableName)
	if err != nil {
		return err
	}

	ret := make([]map[string]interface{}, 0, 1000)
	for start < total {
		sqlStatement := fmt.Sprintf(`SELECT * FROM "%s" LIMIT %d OFFSET %d;`, tableName, 1000, start)
		//fmt.Println(sqlStatement)
		rows, err := c.db.Query(sqlStatement)
		if err != nil {
			return err
		}

		var columns []string
		var values []interface{}

		for rows.Next() {
			start++
			if len(columns) == 0 || len(values) != len(columns) {
				columns, err = rows.Columns()
				if err != nil {
					return err
				}

				columnsLen := len(columns)
				values = make([]interface{}, 0, columnsLen)
				for i := 0; i < columnsLen; i++ {
					values = append(values, new(interface{}))
				}
			}
			//fmt.Println(columns, values)
			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			mid := map[string]interface{}{}
			for i, k := range columns {
				mid[k] = *(values[i].(*interface{}))
			}
			ret = append(ret, mid)
		}

		e := callback(ret)
		ret = make([]map[string]interface{}, 0, 1000)
		_ = rows.Close()

		if e != nil {
			return e
		}
	}

	return nil
}

func (c *Client) Delete(tableName, key string) error {
	sqlStr := `
DELETE FROM "%s" 
WHERE __key__ = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, tableName, key)
	//fmt.Println(sqlStatement)

	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 清空表
func (this *Client) Truncate(tableName string) error {
	sqlStr := `
TRUNCATE TABLE %s;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	//fmt.Println(sqlStatement)
	smt, err := this.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 表行数
func (this *Client) Count(tableName string) (int, error) {
	sqlStr := `
select count(*) from %s;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	smt, err := this.db.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}
	row := smt.QueryRow()
	var count int
	err = row.Scan(&count)
	return count, err
}

// 执行 sql
func (this *Client) ExecSql(sqlStr string) error {
	smt, err := this.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}
