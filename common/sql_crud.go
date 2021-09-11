package common

// 这个下面都是 数据库相关，是不是应该放到另一个 package 下
import (
	"errors"
	"reflect"

	"github.com/jinzhu/gorm"
)

type QueryChain struct {
	query [][]interface{}
}
type RawQuery struct {
	Predicate string
	Args      []interface{}
}

// reference for predicate and expression definition: https://en.wikipedia.org/wiki/SQL_syntax#Language_elements
func (q *QueryChain) Where(predicate interface{}, expressions ...interface{}) {
	q.query = append(q.query, append([]interface{}{"where", predicate}, expressions...))
}

func (q *QueryChain) And(predicate interface{}, expressions ...interface{}) {
	q.query = append(q.query, append([]interface{}{"and", predicate}, expressions...))
}

func (q *QueryChain) Or(predicate interface{}, expressions ...interface{}) {
	q.query = append(q.query, append([]interface{}{"or", predicate}, expressions...))
}

func (q *QueryChain) Not(predicate interface{}, expressions ...interface{}) {
	q.query = append(q.query, append([]interface{}{"not", predicate}, expressions...))
}

func (q *QueryChain) Select(predicate interface{}, expressions ...interface{}) {
	q.query = append(q.query, append([]interface{}{"select", predicate}, expressions...))
}

func ComposeQuery(db *gorm.DB, q *QueryChain) *gorm.DB {
	for _, clause := range q.query {
		switch keyword := clause[0].(string); keyword {
		case "select":
			if len(clause) > 2 {
				db = db.Select(clause[1], clause[2:]...)
			} else {
				db = db.Select(clause[1])
			}
		case "where":
			if len(clause) > 2 {
				db = db.Where(clause[1], clause[2:]...)
			} else {
				db = db.Where(clause[1])
			}
		case "and":
			if len(clause) > 2 {
				db = db.Where(clause[1], clause[2:]...)
			} else {
				db = db.Where(clause[1])
			}
		case "or":
			if len(clause) > 2 {
				db = db.Or(clause[1], clause[2:]...)
			} else {
				db = db.Or(clause[1])
			}
		case "not":
			if len(clause) > 2 {
				db = db.Not(clause[1], clause[2:]...)
			} else {
				db = db.Not(clause[1])
			}

		}
	}
	return db
}

func getTableName(record interface{}) (string, error) {
	var tableName string
	//switch aRecord := record.(type) {
	//default:
	//	return "", errors.New("Unrecognized db model check sql crud")
	//case *model.AccessKey:
	//	tableName = aRecord.TableName()
	//case *model.User:
	//	tableName = aRecord.TableName()
	//case *model.Group:
	//	tableName = aRecord.TableName()
	//case *model.MeshService:
	//	tableName = aRecord.TableName()
	//case *model.Admin:
	//	tableName = aRecord.TableName()
	//case *model.Policy:
	//	tableName = aRecord.TableName()
	//case *model.ServiceMethod:
	//	tableName = aRecord.TableName()
	//case *model.UserGroup:
	//	tableName = aRecord.TableName()
	//case *model.PublicResource:
	//	tableName = aRecord.TableName()
	//case *model.OPAgent:
	//	tableName = aRecord.TableName()
	//}

	return tableName, nil
}

func AddOne(MysqlCli *gorm.DB, aRecord interface{}) error {
	tableName, err := getTableName(aRecord)
	if err != nil {
		return err
	}

	var db = MysqlCli.Table(tableName)

	return db.Create(aRecord).Error
}

// DeleteOne deletes a single record matching the keys provided. Returns error if multiple records found.
func DeleteOne(MysqlCli *gorm.DB, aRecord interface{}) error {
	tableName, err := getTableName(aRecord)
	if err != nil {
		return err
	}

	var cnt int
	record := reflect.ValueOf(aRecord).Interface()
	var db = MysqlCli.Table(tableName)
	db.Where(aRecord).Find(record).Count(&cnt)

	if cnt > 1 {
		return errors.New("multiple records found for the keywords provided")
	}
	return db.Delete(record).Error
}

func UpdateOne(MysqlCli *gorm.DB, aRecord interface{}, changedFields map[string]interface{}) error {
	tableName, err := getTableName(aRecord)
	if err != nil {
		return err
	}

	db := MysqlCli
	record := reflect.ValueOf(aRecord).Interface()
	var cnt int
	db.Table(tableName).Where(aRecord).Find(record).Count(&cnt)

	if cnt > 1 {
		return errors.New("multiple records found for the keywords provided")
	}
	if cnt == 0 {
		return errors.New("no record found for the keywords provided")
	}

	return db.Model(record).Updates(changedFields).Error
}

func Find(MysqlCli *gorm.DB, dummyRecord interface{}, query map[string]interface{},
	limit, offset int32) (interface{}, int, error) {
	//if limit < -1 || limit == 0 || limit > PAGE_SIZE
	if limit < -1 || limit == 0  {
		return nil, 0, errors.New("invalid page size")
	}
	if offset < -1 {
		return nil, 0, errors.New("offset cannot be negative")
	}

	tableName, err := getTableName(dummyRecord)
	if err != nil {
		return nil, 0, err
	}

	var db = MysqlCli.Table(tableName)
	modelType := reflect.Indirect(reflect.ValueOf(dummyRecord)).Type()
	records := reflect.New(reflect.SliceOf(modelType))
	err = db.Where(query).Limit(limit).Offset(offset).Find(records.Interface()).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = nil
		} else {
			return nil, 0, err
		}

	}
	var total int
	err = db.Where(query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return reflect.Indirect(records).Interface(), total, nil
}

// FindLike support `LIKE` in query with many fields, and the relation is `AND`
func FindLike(MysqlCli *gorm.DB, dummyRecord interface{}, query map[string]interface{},
	search map[string]interface{}, limit, offset int32) (interface{}, int, error) {
	//|| limit > PAGE_SIZE
	if limit < -1 || limit == 0  {
		return nil, 0, errors.New("invalid page size")
	}
	if offset < -1 {
		return nil, 0, errors.New("offset cannot be negative")
	}

	tableName, err := getTableName(dummyRecord)
	if err != nil {
		return nil, 0, err
	}

	var db = MysqlCli.Table(tableName)
	modelType := reflect.Indirect(reflect.ValueOf(dummyRecord)).Type()
	records := reflect.New(reflect.SliceOf(modelType))

	Q := db.Where(query)
	for k, v := range search {
		Q = Q.Where(k, v)
	}
	err = Q.Limit(limit).Offset(offset).Find(records.Interface()).Error
	if err != nil {
		return nil, 0, err
	}
	var total int
	err = Q.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return reflect.Indirect(records).Interface(), total, nil
}

// ShadowFind support `LIKE` in query with many fields, and the relation is `AND`
func ShadowFind(MysqlCli *gorm.DB, dummyRecord interface{}, qc *QueryChain,
	limit, offset int32) (interface{}, int, error) {
	//|| limit > PAGE_SIZE
	if limit < -1 || limit == 0 {
		return nil, 0, errors.New("invalid page size")
	}
	if offset < -1 {
		return nil, 0, errors.New("offset cannot be negative")
	}

	tableName, err := getTableName(dummyRecord)
	if err != nil {
		return nil, 0, err
	}

	var db = MysqlCli.Table(tableName)
	modelType := reflect.Indirect(reflect.ValueOf(dummyRecord)).Type()
	records := reflect.New(reflect.SliceOf(modelType))

	Q := ComposeQuery(db, qc)
	err = Q.Limit(limit).Offset(offset).Find(records.Interface()).Error
	if err != nil {
		return nil, 0, err
	}
	var total int
	err = Q.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return reflect.Indirect(records).Interface(), total, nil
}
